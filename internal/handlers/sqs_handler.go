package handlers

import (
	"context"
	"example.com/smart-devices/internal/services"
	"github.com/aws/aws-lambda-go/events"
	"go.uber.org/zap"
)

type SQSHandler struct {
	svc    *services.SQSService
	logger *zap.Logger
}

func NewSQSHandler(sqsService *services.SQSService, logger *zap.Logger) *SQSHandler {
	return &SQSHandler{
		svc:    sqsService,
		logger: logger,
	}
}

func (h *SQSHandler) ProcessMessage(ctx context.Context, sqsEvent events.SQSEvent) error {
	for _, record := range sqsEvent.Records {
		if err := h.svc.ProcessMessage(ctx, record.Body); err != nil {
			h.logger.Error("Error processing message", zap.Error(err))
			return err
		}
	}
	return nil
}
