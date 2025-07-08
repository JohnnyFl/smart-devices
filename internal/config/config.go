package config

import (
	"os"
)

type Config struct {
	DynamoDBTable string
	SQSQueueURL   string
	AWSRegion     string
	Stage         string
	DynamoDBURL   string
}

func Load() *Config {
	return &Config{
		DynamoDBTable: getEnv("DYNAMODB_TABLE", "devices"),
		SQSQueueURL:   getEnv("SQS_QUEUE_URL", ""),
		AWSRegion:     getEnv("AWS_REGION", "us-east-1"),
		Stage:         getEnv("STAGE", "dev"),
		DynamoDBURL:   os.Getenv("DYNAMODB_URL"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
