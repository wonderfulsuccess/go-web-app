package webserver

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/wonderfulsuccess/go-web-app/back/config"
	"github.com/wonderfulsuccess/go-web-app/back/logger"
)

// Server bundles together the Gin engine, Gorm connection and websocket hub.
type Server struct {
	cfg        config.Config
	httpServer *http.Server
	hub        *Hub
	quit       chan struct{}
	runMu      sync.Mutex
	running    bool
}

func NewServer(cfg config.Config, db *gorm.DB) *Server {
	gin.SetMode(cfg.Mode)
	if cfg.Mode == gin.ReleaseMode {
		gin.DefaultWriter = io.Discard
		gin.DefaultErrorWriter = io.Discard
	}

	hub := NewHub()
	router := NewRouter(cfg, db, hub)

	srv := &http.Server{
		Addr:    cfg.Address(),
		Handler: router,
	}

	server := &Server{
		cfg:        cfg,
		httpServer: srv,
		hub:        hub,
		quit:       make(chan struct{}),
	}

	go hub.Run()
	go server.logIncomingMessages()

	return server
}

func (s *Server) logIncomingMessages() {
	for msg := range s.hub.Incoming() {
		payload := string(msg.Payload)
		if payload == "" || payload == "null" {
			payload = "<empty>"
		}
		logger.Infof("websocket message type=%s sender=%s receiver=%s payload=%s", msg.Type, msg.Sender, msg.Receiver, payload)

		if msg.Type == "demo-start" {
			logger.Infof("websocket demo start requested by %s", msg.Sender)
			s.ensureDemoBroadcast()
		}
	}
}

func (s *Server) ensureDemoBroadcast() {
	s.runMu.Lock()
	if s.running {
		s.runMu.Unlock()
		return
	}
	s.running = true
	s.runMu.Unlock()

	logger.Infof("starting websocket demo broadcast loop")
	go s.broadcastLoop()
}

func (s *Server) broadcastLoop() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	counter := 0
	for {
		select {
		case <-s.quit:
			return
		case t := <-ticker.C:
			counter++
			payload, err := json.Marshal(map[string]string{
				"message": fmt.Sprintf("server tick #%d", counter),
				"sentAt":  t.UTC().Format(time.RFC3339Nano),
			})
			if err != nil {
				logger.Errorf("failed to marshal websocket payload: %v", err)
				continue
			}

			s.hub.SendMessage(WSMessage{
				Sender:    "server",
				Receiver:  "*",
				Type:      "server-tick",
				Timestamp: t.UTC(),
				Payload:   payload,
			})
		}
	}
}

// Start runs the HTTP server until the context is cancelled.
func (s *Server) Start(ctx context.Context) error {
	defer close(s.quit)

	errCh := make(chan error, 1)

	go func() {
		logger.Infof("starting webserver on %s", s.httpServer.Addr)
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- err
		}
	}()

	select {
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return s.httpServer.Shutdown(shutdownCtx)
	case err := <-errCh:
		return err
	}
}

// Hub exposes the websocket hub so other packages can push messages.
func (s *Server) Hub() *Hub {
	return s.hub
}
