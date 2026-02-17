package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"

	"oneproxy-clientwebui/internal/config"
	"oneproxy-clientwebui/internal/handler"
)

func main() {
	gin.SetMode(gin.ReleaseMode)

	if err := config.Load(); err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	r := gin.Default()

	api := r.Group("/api")
	{
		api.GET("/usage/query", handler.QueryUsage)
		api.POST("/admin/login", handler.AdminLogin)

		admin := api.Group("/config")
		admin.Use(handler.AdminAuth)
		{
			admin.GET("", handler.GetConfig)
			admin.POST("", handler.UpdateConfig)
		}
	}

	// Reverse proxy for AI API
	proxy := handler.NewProxyHandler()
	proxyHandler := gin.WrapF(proxy)
	r.Any("/v1/*any", proxyHandler)

	// Frontend static files
	r.Static("/assets", "./web/dist/assets")
	r.StaticFile("/favicon.ico", "./web/dist/favicon.ico")
	r.StaticFile("/", "./web/dist/index.html")
	r.NoRoute(func(c *gin.Context) {
		c.File("./web/dist/index.html")
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server starting on :%s", port)
	r.Run(":" + port)
}
