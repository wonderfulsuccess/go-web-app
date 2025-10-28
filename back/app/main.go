package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/wonderfulsuccess/go-web-app/back/config"
	"github.com/wonderfulsuccess/go-web-app/back/database"
	"github.com/wonderfulsuccess/go-web-app/back/logger"
	"github.com/wonderfulsuccess/go-web-app/back/model"
	"github.com/wonderfulsuccess/go-web-app/back/utils"
	"github.com/wonderfulsuccess/go-web-app/back/webserver"
)

func main() {

	logger.Infof("Starting server...")

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg := config.Load()

	utils.MapPrint(cfg)

	db, err := database.InitDatabase(cfg.Database)
	if err != nil {
		logger.Errorf("failed to initialise database: %v", err)
	}

	if err := model.AutoMigrate(db); err != nil {
		logger.Errorf("failed to migrate database schema: %v", err)
	}

	server := webserver.NewServer(cfg, db)

	if err := server.Start(ctx); err != nil {
		logger.Errorf("server exited with error: %v", err)
	}
}
