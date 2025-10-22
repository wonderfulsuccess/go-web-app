package webserver

import (
	"context"
	"log"
	"net/http"
	"time"

	"gorm.io/gorm"

	"github.com/wonderfulsuccess/go-web-app/back/config"
)

// Server bundles together the Gin engine, Gorm connection and websocket hub.
type Server struct {
	cfg        config.Config
	httpServer *http.Server
	hub        *Hub
}

func NewServer(cfg config.Config, db *gorm.DB) *Server {
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
	}

	go hub.Run()
	go server.logIncomingMessages()

	return server
}

func (s *Server) logIncomingMessages() {
	for msg := range s.hub.Incoming() {
		log.Printf("websocket message type=%s sender=%s receiver=%s", msg.Type, msg.Sender, msg.Receiver)
	}
}

// Start runs the HTTP server until the context is cancelled.
func (s *Server) Start(ctx context.Context) error {
	errCh := make(chan error, 1)

	go func() {
		log.Printf("starting webserver on %s", s.httpServer.Addr)
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
