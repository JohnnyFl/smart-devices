package main

import (
	"example.com/smart-devices/internal/handlers"
	"example.com/smart-devices/internal/setup"
	"github.com/aws/aws-lambda-go/lambda"
	"go.uber.org/zap"
)

var (
	sqsHandler *handlers.SQSHandler
	logger     *zap.Logger
)

func init() {
	_, sqsHandler, logger = setup.SetupComponents()
}

func main() {
	lambda.Start(sqsHandler.ProcessMessage)
}
