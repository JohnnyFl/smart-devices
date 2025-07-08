package main

import (
	"context"
	appConfig "example.com/smart-devices/internal/config"
	"example.com/smart-devices/internal/handlers"
	"example.com/smart-devices/internal/repository"
	"example.com/smart-devices/internal/services"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"go.uber.org/zap"
)

var (
	deviceHandler *handlers.DeviceHandler
	logger        *zap.Logger
)

func init() {
	cfg := appConfig.Load()

	loggerCfg := zap.NewProductionConfig()
	loggerCfg.OutputPaths = []string{"stdout"}
	logger, _ = loggerCfg.Build()

	// Default AWS DynamoDB
	awsCfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(cfg.AWSRegion))

	if err != nil {
		logger.Fatal("failed to load AWS config", zap.Error(err))
	}

	// Create DynamoDB client with custom endpoint for local development
	var dynamoClient *dynamodb.Client
	if cfg.DynamoDBURL != "" {
		logger.Info("Using custom DynamoDB endpoint", zap.String("url", cfg.DynamoDBURL))

		// Create the DynamoDB client with the custom endpoint
		dynamoClient = dynamodb.NewFromConfig(awsCfg, func(o *dynamodb.Options) {
			o.BaseEndpoint = &cfg.DynamoDBURL
			o.EndpointOptions = dynamodb.EndpointResolverOptions{
				DisableHTTPS: true,
			}
		})
	} else {
		dynamoClient = dynamodb.NewFromConfig(awsCfg)
	}

	logger.Info("DynamoDB client initialized",
		zap.String("table", cfg.DynamoDBTable),
		zap.String("region", cfg.AWSRegion),
	)

	deviceRepo := repository.NewDeviceRepository(dynamoClient, cfg.DynamoDBTable, logger)
	deviceService := services.NewDeviceService(deviceRepo, logger)
	deviceHandler = handlers.NewDeviceHandler(deviceService, logger)
}

func main() {

	lambda.Start(deviceHandler.CreateDevice)
}
