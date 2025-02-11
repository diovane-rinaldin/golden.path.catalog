package views

import (
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/diovane-rinaldin/golden.path.catalog/backend/handlers"
	"github.com/diovane-rinaldin/golden.path.catalog/backend/middleware"
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
	r.POST("/auth", auth.Authenticate)

	// Protected routes
	api := r.Group("/api")
	api.Use(middleware.AuthMiddleware())

	// Technology routes
	tech := handlers.NewTechnologyHandler(db)
	api.POST("/technology", tech.Create)
	api.GET("/technology", tech.List)
	api.GET("/technology/:name", tech.Get)
	api.PUT("/technology/:id", tech.Update)

	// Component routes
	comp := handlers.NewComponentHandler(db)
	api.POST("/component", comp.Create)
	api.GET("/component", comp.List)
	api.GET("/component/:name", comp.Get)
	api.GET("/component/technology/:id", comp.ListByTechnology)
	api.PUT("/component/:id", comp.Update)

	// Service routes
	svc := handlers.NewServiceHandler(db)
	api.POST("/service", svc.Create)
	api.GET("/service", svc.List)
	api.GET("/service/:name", svc.Get)
	api.GET("/service/component/:id", svc.ListByComponent)
	api.PUT("/service/:id", svc.Update)

	// Image upload route
	upload := handlers.NewUploadHandler(s3)
	api.POST("/upload", upload.UploadImage)
}
