package utils

import (
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
)

type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func JSONErrorResponse(status int, code, message string) events.APIGatewayProxyResponse {
	body, _ := json.Marshal(ErrorResponse{
		Code:    code,
		Message: message,
	})
	return events.APIGatewayProxyResponse{
		StatusCode: status,
		Body:       string(body),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}
}

func JSONSuccessResponse(status int, data interface{}) events.APIGatewayProxyResponse {
	body, _ := json.Marshal(data)
	return events.APIGatewayProxyResponse{
		StatusCode: status,
		Body:       string(body),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}
}
