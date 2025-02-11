package main

import (
	"context"
	"log"

	"github.com/diovane-rinaldin/golden.path.catalog/backend/utils"
	"github.com/diovane-rinaldin/golden.path.catalog/backend/views"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

func initTracer() (*sdktrace.TracerProvider, error) {
	exporter, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
	if err != nil {
		return nil, err
	}
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithBatcher(exporter),
	)
	otel.SetTracerProvider(tp)
	return tp, nil
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	tp, err := initTracer()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Printf("Error shutting down tracer provider: %v", err)
		}
	}()

	// Initialize AWS services
	dynamoDBClient := utils.InitDynamoDB()
	kmsClient := utils.InitKMS()
	s3Client := utils.InitS3()

	router := gin.Default()

	// Initialize routes
	views.SetupRoutes(router, dynamoDBClient, kmsClient, s3Client)

	router.Run(":8080")
}
