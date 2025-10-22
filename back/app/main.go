package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/wonderfulsuccess/go-web-app/back/config"
	"github.com/wonderfulsuccess/go-web-app/back/model"
	"github.com/wonderfulsuccess/go-web-app/back/utils"
	"github.com/wonderfulsuccess/go-web-app/back/webserver"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg := config.Load()

	db, err := utils.InitDatabase(cfg.Database)
	if err != nil {
		log.Fatalf("failed to initialise database: %v", err)
	}

	if err := model.AutoMigrate(db); err != nil {
		log.Fatalf("failed to migrate database schema: %v", err)
	}

	server := webserver.NewServer(cfg, db)

	if err := server.Start(ctx); err != nil {
		log.Fatalf("server exited with error: %v", err)
	}
}
