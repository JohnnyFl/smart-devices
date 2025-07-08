package main

import (
	"example.com/smart-devices/internal/handlers"
	"example.com/smart-devices/internal/setup"
	"github.com/aws/aws-lambda-go/lambda"
	"go.uber.org/zap"
)

var (
	deviceHandler *handlers.DeviceHandler
	logger        *zap.Logger
)

func init() {
	deviceHandler, _, logger = setup.SetupComponents()
}

func main() {
	lambda.Start(deviceHandler.GetDevices)
}
