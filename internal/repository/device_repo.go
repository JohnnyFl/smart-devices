package repository

import (
	"context"
	"example.com/smart-devices/internal/models"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"strconv"
	"strings"
	"time"
)

type DeviceRepository struct {
	client    *dynamodb.Client
	tableName string
	logger    *zap.Logger
}

func NewDeviceRepository(client *dynamodb.Client, tableName string, logger *zap.Logger) *DeviceRepository {
	//func NewDeviceRepository(client *dynamodb.Client, tableName string) *DeviceRepository {
	return &DeviceRepository{
		client:    client,
		tableName: tableName,
		logger:    logger,
	}
}

func (r *DeviceRepository) GetDevice(ctx context.Context, id string) (*models.Device, error) {
	r.logger.Debug("fetching device", zap.String("device_id", id))

	result, err := r.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: &r.tableName,
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{
				Value: id,
			},
		},
	})

	if err != nil {
		r.logger.Error("database operation failed",
			zap.String("operation", "GetDevice"),
			zap.String("table", r.tableName),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to get device: %w", err)
	}

	if result.Item == nil {
		return nil, fmt.Errorf("device not found")
	}

	var device models.Device

	if err := device.FromMap(result.Item); err != nil {
		r.logger.Error("failed to unmarshal device",
			zap.String("device_id", id),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to unmarshal device: %w", err)
	}

	return &device, nil

}

func (r *DeviceRepository) GetDevices(ctx context.Context) ([]models.Device, error) {
	r.logger.Debug("fetching device")

	result, err := r.client.Scan(ctx, &dynamodb.ScanInput{
		TableName: &r.tableName,
	})

	r.logger.Debug("fetching device",
		zap.Any("count", result))

	if err != nil {
		r.logger.Error("database operation failed",
			zap.String("operation", "GetDevices"),
			zap.String("table", r.tableName),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to get devices: %w", err)
	}

	if result.Count == 0 {
		return nil, fmt.Errorf("no devices found")
	}

	var devices []models.Device

	for _, item := range result.Items {
		var device models.Device
		err := attributevalue.UnmarshalMap(item, &device)
		if err != nil {
			r.logger.Error("failed to unmarshal device",
				zap.Error(err))

		}
		devices = append(devices, device)
	}

	return devices, nil

}

func (r *DeviceRepository) DeleteDevice(ctx context.Context, id string) error {
	r.logger.Debug("deleting device", zap.String("device_id", id))

	_, err := r.client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: &r.tableName,
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{
				Value: id,
			},
		},
	})

	if err != nil {
		r.logger.Error("database operation failed",
			zap.String("operation", "DeleteDevice"),
			zap.String("table", r.tableName),
			zap.Error(err),
		)
		return fmt.Errorf("failed to delete device: %w", err)
	}

	return nil

}

func (r *DeviceRepository) UpdateDevice(ctx context.Context, id string, update models.Device) (*models.Device, error) {
	r.logger.Debug("updating device", zap.String("device_id", id))

	// First, get the current device to preserve existing fields
	result, err := r.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: &r.tableName,
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: id},
		},
	})
	if err != nil {
		r.logger.Error("failed to get device for update",
			zap.String("device_id", id),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to get device: %w", err)
	}
	if result.Item == nil {
		return nil, fmt.Errorf("device not found")
	}

	// Unmarshal the current device
	var currentDevice models.Device
	if err := attributevalue.UnmarshalMap(result.Item, &currentDevice); err != nil {
		r.logger.Error("failed to unmarshal current device",
			zap.String("device_id", id),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to unmarshal current device: %w", err)
	}

	// Create a map of fields to update
	updates := make(map[string]types.AttributeValue)

	// Only include fields that are not zero values
	if update.Type != "" {
		updates[":type"] = &types.AttributeValueMemberS{Value: update.Type}
	}
	if update.Name != "" {
		updates[":name"] = &types.AttributeValueMemberS{Value: update.Name}
	}
	if update.MAC != "" {
		updates[":mac"] = &types.AttributeValueMemberS{Value: update.MAC}
	}
	if update.HomeID != "" {
		updates[":homeId"] = &types.AttributeValueMemberS{Value: update.HomeID}
	}

	// Always update ModifiedAt
	now := time.Now().Unix()
	updates[":modifiedAt"] = &types.AttributeValueMemberN{Value: strconv.FormatInt(now, 10)}

	if len(updates) == 1 { // Only ModifiedAt was updated
		return &currentDevice, nil
	}

	// Build the update expression
	var updateExpr []string
	exprAttrNames := make(map[string]string)
	for k := range updates {
		field := strings.TrimPrefix(k, ":")
		exprAttrNames["#"+field] = field
		updateExpr = append(updateExpr, fmt.Sprintf("#%s = %s", field, k))
	}

	// Execute the update
	_, err = r.client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName: &r.tableName,
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: id},
		},
		UpdateExpression:          aws.String("SET " + strings.Join(updateExpr, ", ")),
		ExpressionAttributeNames:  exprAttrNames,
		ExpressionAttributeValues: updates,
		ReturnValues:              types.ReturnValueAllNew,
	})

	if err != nil {
		r.logger.Error("failed to update device",
			zap.String("device_id", id),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to update device: %w", err)
	}

	// Get the updated device
	updatedResult, err := r.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: &r.tableName,
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: id},
		},
	})

	if err != nil {
		r.logger.Error("failed to fetch updated device",
			zap.String("device_id", id),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to fetch updated device: %w", err)
	}

	var updatedDevice models.Device
	if err := attributevalue.UnmarshalMap(updatedResult.Item, &updatedDevice); err != nil {
		r.logger.Error("failed to unmarshal updated device",
			zap.String("device_id", id),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to unmarshal updated device: %w", err)
	}

	return &updatedDevice, nil
}

func (r *DeviceRepository) CreateDevice(ctx context.Context, device models.Device) (models.Device, error) {
	r.logger.Debug("creating device", zap.String("device_id", device.ID))

	now := time.Now().UnixMilli()
	device.ID = uuid.New().String()
	device.CreatedAt = now
	device.ModifiedAt = now

	item, err := attributevalue.MarshalMap(device)
	if err != nil {
		return device, fmt.Errorf("failed to marshal device: %w", err)
	}

	_, err = r.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(r.tableName),
		Item:      item,
	})

	if err != nil {
		r.logger.Error("database operation failed",
			zap.String("operation", "CreateDevice"),
			zap.String("table", r.tableName),
			zap.Error(err),
		)
		return device, fmt.Errorf("failed to get device: %w", err)
	}

	return device, nil

}

func (r *DeviceRepository) UpdateDeviceHomeID(ctx context.Context, id string, homeID string) error {
	r.logger.Debug("updating device", zap.String("device_id", id))

	// Get current timestamp for ModifiedAt
	now := time.Now().Unix()

	_, err := r.client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName: &r.tableName,
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: id}},
		UpdateExpression: aws.String("SET #homeId = :homeId, #modifiedAt = :modifiedAt"),
		ExpressionAttributeNames: map[string]string{
			"#homeId":     "homeId",
			"#modifiedAt": "modifiedAt",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":homeId":     &types.AttributeValueMemberS{Value: homeID},
			":modifiedAt": &types.AttributeValueMemberN{Value: strconv.FormatInt(now, 10)},
		},
		ReturnValues: types.ReturnValueAllNew,
	})

	if err != nil {
		r.logger.Error("failed to update device",
			zap.String("device_id", id),
			zap.String("home_id", homeID),
			zap.Error(err),
		)
		return fmt.Errorf("failed to update device: %w", err)
	}

	return nil
}
