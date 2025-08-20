package rest

import (
	"consumer/internal/rest/handlers"
	"consumer/internal/services"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	// _ "blog-posts-api/docs"
)

// @title Swap Stats API
// @version 1.0
// @description REST API to fetch swap piar stats

// @host localhost:8081
// @BasePath /api/v1
// @schemes http https

// @tag.name Stats
// @tag.description Operations related to stats

type RestApi struct {
	port         string
	statsService *services.StatsService
}

func New(port string, statsService *services.StatsService) *RestApi {
	return &RestApi{port, statsService}
}

func (s *RestApi) Run() error {
	r := gin.Default()

	// Add CORS middleware for Swagger UI
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// Swagger documentation route
	r.GET("/api/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Health check endpoint
	r.GET("/api/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	handler := handlers.NewStatsHandler(s.statsService)
	v1 := r.Group("/api/v1")
	{
		handler.RegisterRoutes(v1)
	}

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message":  "Welcome to Swap stats API",
			"version":  "1.0.0",
			"docs":     "/api/docs/index.html",
			"health":   "/api/health",
			"api_base": "/api/v1",
			"endpoints": map[string]string{
				"GET /api/v1/stats/:period/tokens/:token": "Get single token stats in a specific period",
				"GET /api/v1/stats/:period/pairs/:pair":   "Get swap pair stats in a specific period",
			},
		})
	})

	log.Println("Stats API is starting...")
	log.Printf("Health check available at: http://localhost:%s/api/health\n", s.port)
	log.Printf("API endpoints available at: http://localhost:%s/api/v1\n", s.port)
	log.Printf("Swagger docs available at: http://localhost:%s/api/docs/index.html\n", s.port)

	addr := fmt.Sprintf(":%s", s.port)
	if err := r.Run(addr); err != nil {
		return errors.Wrapf(err, "failed to start http server")
	}
	return nil
}
