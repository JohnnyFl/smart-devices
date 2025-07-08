package models

import (
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type Device struct {
	ID         string `json:"id" dynamodbav:"id"`
	MAC        string `json:"mac" dynamodbav:"mac"`
	Name       string `json:"name" dynamodbav:"name"`
	Type       string `json:"type" dynamodbav:"type"`
	HomeID     string `json:"homeId" dynamodbav:"homeId"`
	CreatedAt  int64  `json:"createdAt" dynamodbav:"createdAt"`
	ModifiedAt int64  `json:"modifiedAt" dynamodbav:"modifiedAt"`
}

type CreateDeviceRequest struct {
	MAC    string `json:"mac" validate:"required,mac"`
	Name   string `json:"name" validate:"required,min=1,max=100"`
	Type   string `json:"type" validate:"required,oneof=thermostat light camera sensor"`
	HomeID string `json:"homeId" validate:"required,uuid"`
}

type UpdateDeviceRequest struct {
	Name   *string `json:"name,omitempty" validate:"omitempty,min=1,max=100"`
	Type   *string `json:"type,omitempty" validate:"omitempty,oneof=thermostat light camera sensor"`
	HomeID *string `json:"homeId,omitempty" validate:"omitempty,uuid"`
}

type SQSMessage struct {
	DeviceID string `json:"deviceId"`
	HomeID   string `json:"homeId"`
	Action   string `json:"action"`
}

// ToMap converts Device to map[string]types.AttributeValue for DynamoDB
func (d *Device) ToMap() (map[string]types.AttributeValue, error) {
	return attributevalue.MarshalMap(d)
}

// FromMap converts map[string]types.AttributeValue to Device
func (d *Device) FromMap(item map[string]types.AttributeValue) error {
	return attributevalue.UnmarshalMap(item, d)
}
