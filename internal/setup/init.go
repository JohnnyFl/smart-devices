package setup

import (
	"context"
	appConfig "example.com/smart-devices/internal/config"
	"example.com/smart-devices/internal/handlers"
	"example.com/smart-devices/internal/repository"
	"example.com/smart-devices/internal/services"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"go.uber.org/zap"
)

// SetupComponents initializes all common components and returns handlers and logger
func SetupComponents() (*handlers.DeviceHandler, *handlers.SQSHandler, *zap.Logger) {
	cfg := appConfig.Load()

	// Initialize logger
	loggerCfg := zap.NewProductionConfig()
	loggerCfg.OutputPaths = []string{"stdout"}
	logger, err := loggerCfg.Build()
	if err != nil {
		panic(err)
	}

	// Load AWS configuration
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

	// Initialize repository, services, and handlers
	deviceRepo := repository.NewDeviceRepository(dynamoClient, cfg.DynamoDBTable, logger)
	deviceService := services.NewDeviceService(deviceRepo, logger)
	sqsService := services.NewSQSService(deviceService, logger)

	deviceHandler := handlers.NewDeviceHandler(deviceService, logger)
	sqsHandler := handlers.NewSQSHandler(sqsService, logger)

	return deviceHandler, sqsHandler, logger
}
