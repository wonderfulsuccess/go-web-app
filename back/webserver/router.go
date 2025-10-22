package webserver

import (
	"net/http"
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
	router.Static("/static", filepath.Join(cfg.StaticDir, "assets"))
	router.StaticFile("/favicon.ico", filepath.Join(cfg.StaticDir, "favicon.ico"))
	router.StaticFile("/manifest.webmanifest", filepath.Join(cfg.StaticDir, "manifest.webmanifest"))

	router.NoRoute(func(c *gin.Context) {
		c.File(filepath.Join(cfg.StaticDir, "index.html"))
	})

	return router
}
