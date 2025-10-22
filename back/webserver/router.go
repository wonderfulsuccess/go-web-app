package webserver

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/wonderfulsuccess/go-web-app/back/config"
	"github.com/wonderfulsuccess/go-web-app/back/controller"
)

// NewRouter wires the HTTP endpoints for API and static assets.
func NewRouter(cfg config.Config, db *gorm.DB, hub *Hub) *gin.Engine {
	router := gin.Default()

	api := router.Group("/api")
	{
		api.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		})

		userGroup := api.Group("/users")
		userController := controller.NewUserController(db)
		userController.RegisterRoutes(userGroup)

		api.GET("/ws", hub.HandleWebSocket)
	}

	// Serve the compiled front-end assets.
	assetDir := filepath.Join(cfg.StaticDir, "assets")
	if stat, err := os.Stat(assetDir); err == nil && stat.IsDir() {
		router.Static("/assets", assetDir)
	}
	bindStaticFile := func(route, filename string) {
		path := filepath.Join(cfg.StaticDir, filename)
		if _, err := os.Stat(path); err == nil {
			router.StaticFile(route, path)
		}
	}
	bindStaticFile("/favicon.ico", "favicon.ico")
	bindStaticFile("/manifest.webmanifest", "manifest.webmanifest")

	router.NoRoute(func(c *gin.Context) {
		indexPath := filepath.Join(cfg.StaticDir, "index.html")
		if _, err := os.Stat(indexPath); err != nil {
			c.Status(http.StatusNotFound)
			return
		}
		c.File(indexPath)
	})

	return router
}
