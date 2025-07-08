package utils

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"os"
)

func GetDynamoDBClient(awsCfg aws.Config) *dynamodb.Client {
	// Override endpoint for local development
	if os.Getenv("IS_LOCAL") == "true" {
		localEndpoint := os.Getenv("DYNAMODB_ENDPOINT")
		if localEndpoint == "" {
			localEndpoint = "http://localhost:8000"
		}

		// Use the modern BaseEndpoint approach
		return dynamodb.NewFromConfig(awsCfg, func(o *dynamodb.Options) {
			o.BaseEndpoint = aws.String(localEndpoint)
		})
	}

	return dynamodb.NewFromConfig(awsCfg)
}
