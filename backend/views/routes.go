package views

import (
	"time"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/diovane-rinaldin/golden.path.catalog/backend/handlers"
	"github.com/diovane-rinaldin/golden.path.catalog/backend/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(
	r *gin.Engine,
	db *dynamodb.Client,
	kms *kms.Client,
	s3 *s3.Client,
) {
	// Auth handler
	auth := handlers.NewAuthHandler(kms)
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"}, // Substitua pelo domínio do frontend
		AllowMethods:     []string{"GET", "POST", "PUT"},
		AllowHeaders:     []string{"Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour, // Cache da política de CORS por 12 horas
	}))
	r.POST("/auth", auth.Authenticate)

	// Protected routes
	api := r.Group("/api")
	api.Use(middleware.AuthMiddleware())

	// Technology routes
	tech := handlers.NewTechnologyHandler(db)
	api.POST("/technology", tech.Create)
	api.GET("/technology", tech.List)
	api.GET("/technology/name/:name", tech.GetByName)
	api.GET("/technology/id/:id", tech.GetById)
	api.PUT("/technology", tech.Update)

	// Component routes
	comp := handlers.NewComponentHandler(db)
	api.POST("/component", comp.Create)
	api.GET("/component", comp.List)
	api.GET("/component/name/:name", comp.GetByName)
	api.GET("/component/id/:id", comp.GetById)
	api.GET("/component/technology/:id", comp.ListByTechnology)
	api.PUT("/component", comp.Update)

	// Service routes
	svc := handlers.NewServiceHandler(db)
	api.POST("/service", svc.Create)
	api.GET("/service", svc.List)
	api.GET("/service/:name", svc.Get)
	api.GET("/service/component/:id", svc.ListByComponent)
	api.PUT("/service", svc.Update)

	// Image upload route
	upload := handlers.NewUploadHandler(s3)
	api.POST("/upload", upload.UploadImage)
}
