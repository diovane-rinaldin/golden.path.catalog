package utils

import (
	"context"
	"log"
	"os"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var (
	awsConfig aws.Config
	once      sync.Once
)

func LoadAWSConfig() aws.Config {
	once.Do(func() {
		awsAccessKey := os.Getenv("AWS_ACCESS_KEY_ID")
		awsSecretKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
		awsRegion := os.Getenv("AWS_REGION")

		if awsAccessKey == "" || awsSecretKey == "" || awsRegion == "" {
			log.Fatal("Credenciais ou região AWS ausentes. Defina AWS_ACCESS_KEY_ID, AWS_SECRET_ACCESS_KEY, e AWS_REGION no arquivo .env.")
		}

		cfg, err := config.LoadDefaultConfig(context.Background(),
			config.WithRegion(awsRegion),
			config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(awsAccessKey, awsSecretKey, "")),
		)
		if err != nil {
			log.Fatal("Erro ao carregar as credenciais:", err)
		}

		awsConfig = cfg
	})

	return awsConfig
}

func InitDynamoDB() *dynamodb.Client {
	return dynamodb.NewFromConfig(LoadAWSConfig())
}

func InitKMS() *kms.Client {
	return kms.NewFromConfig(LoadAWSConfig())
}

func InitS3() *s3.Client {
	return s3.NewFromConfig(LoadAWSConfig())
}
