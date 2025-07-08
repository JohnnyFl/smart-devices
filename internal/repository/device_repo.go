package repository

import (
	"context"
	"example.com/smart-devices/internal/errors"
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
		return nil, errors.WrapError(errors.ErrorTypeDatabase, "failed to get device from database", err).
			WithOperation("GetDevice").
			WithLayer("repository").
			WithContext("device_id", id).
			WithContext("table", r.tableName)
	}

	if result.Item == nil {
		return nil, errors.ErrDomainDeviceNotFound.
			WithOperation("GetDevice").
			WithLayer("repository").
			WithContext("device_id", id)
	}

	var device models.Device

	if err := device.FromMap(result.Item); err != nil {
		r.logger.Error("failed to unmarshal device",
			zap.String("device_id", id),
			zap.Error(err),
		)
		return nil, errors.WrapError(errors.ErrorTypeDatabase, "failed to unmarshal device data", err).
			WithOperation("GetDevice").
			WithLayer("repository").
			WithContext("device_id", id)
	}

	return &device, nil
}

func (r *DeviceRepository) GetDevices(ctx context.Context) ([]models.Device, error) {
	r.logger.Debug("fetching devices")

	result, err := r.client.Scan(ctx, &dynamodb.ScanInput{
		TableName: &r.tableName,
	})

	if err != nil {
		r.logger.Error("database operation failed",
			zap.String("operation", "GetDevices"),
			zap.String("table", r.tableName),
			zap.Error(err),
		)
		return nil, errors.WrapError(errors.ErrorTypeDatabase, "failed to scan devices from database", err).
			WithOperation("GetDevices").
			WithLayer("repository").
			WithContext("table", r.tableName)
	}

	r.logger.Debug("fetched devices", zap.Int32("count", result.Count))

	if result.Count == 0 {
		return nil, errors.ErrDomainNoDevicesFound.
			WithOperation("GetDevices").
			WithLayer("repository")
	}

	var devices []models.Device

	for i, item := range result.Items {
		var device models.Device
		if err := attributevalue.UnmarshalMap(item, &device); err != nil {
			r.logger.Error("failed to unmarshal device",
				zap.Int("item_index", i),
				zap.Error(err))
			// Skip malformed items but continue processing
			continue
		}
		devices = append(devices, device)
	}

	// If no devices were successfully unmarshaled
	if len(devices) == 0 && len(result.Items) > 0 {
		return nil, errors.ErrUnmarshalDevice.
			WithOperation("GetDevices").
			WithLayer("repository").
			WithContext("items_count", len(result.Items))
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
		return errors.WrapError(errors.ErrorTypeDatabase, "failed to delete device from database", err).
			WithOperation("DeleteDevice").
			WithLayer("repository").
			WithContext("device_id", id).
			WithContext("table", r.tableName)
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
		return nil, errors.WrapError(errors.ErrorTypeDatabase, "failed to get device for update", err).
			WithOperation("UpdateDevice").
			WithLayer("repository").
			WithContext("device_id", id)
	}
	if result.Item == nil {
		return nil, errors.ErrDomainDeviceNotFound.
			WithOperation("UpdateDevice").
			WithLayer("repository").
			WithContext("device_id", id)
	}

	// Unmarshal the current device
	var currentDevice models.Device
	if err := attributevalue.UnmarshalMap(result.Item, &currentDevice); err != nil {
		r.logger.Error("failed to unmarshal current device",
			zap.String("device_id", id),
			zap.Error(err),
		)
		return nil, errors.WrapError(errors.ErrorTypeDatabase, "failed to unmarshal current device", err).
			WithOperation("UpdateDevice").
			WithLayer("repository").
			WithContext("device_id", id)
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
		return nil, errors.WrapError(errors.ErrorTypeDatabase, "failed to update device in database", err).
			WithOperation("UpdateDevice").
			WithLayer("repository").
			WithContext("device_id", id)
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
		return nil, errors.WrapError(errors.ErrorTypeDatabase, "failed to fetch updated device", err).
			WithOperation("UpdateDevice").
			WithLayer("repository").
			WithContext("device_id", id)
	}

	var updatedDevice models.Device
	if err := attributevalue.UnmarshalMap(updatedResult.Item, &updatedDevice); err != nil {
		r.logger.Error("failed to unmarshal updated device",
			zap.String("device_id", id),
			zap.Error(err),
		)
		return nil, errors.WrapError(errors.ErrorTypeDatabase, "failed to unmarshal updated device", err).
			WithOperation("UpdateDevice").
			WithLayer("repository").
			WithContext("device_id", id)
	}

	return &updatedDevice, nil
}

func (r *DeviceRepository) CreateDevice(ctx context.Context, device models.Device) (models.Device, error) {
	now := time.Now().UnixMilli()
	device.ID = uuid.New().String()
	device.CreatedAt = now
	device.ModifiedAt = now

	r.logger.Debug("creating device", zap.String("device_id", device.ID))

	item, err := attributevalue.MarshalMap(device)
	if err != nil {
		return device, errors.WrapError(errors.ErrorTypeDatabase, "failed to marshal device data", err).
			WithOperation("CreateDevice").
			WithLayer("repository").
			WithContext("device_id", device.ID)
	}

	_, err = r.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(r.tableName),
		Item:      item,
	})

	if err != nil {
		r.logger.Error("database operation failed",
			zap.String("operation", "CreateDevice"),
			zap.String("table", r.tableName),
			zap.String("device_id", device.ID),
			zap.Error(err),
		)
		return device, errors.WrapError(errors.ErrorTypeDatabase, "failed to create device in database", err).
			WithOperation("CreateDevice").
			WithLayer("repository").
			WithContext("device_id", device.ID).
			WithContext("table", r.tableName)
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
		r.logger.Error("failed to update device home ID",
			zap.String("device_id", id),
			zap.String("home_id", homeID),
			zap.Error(err),
		)
		return errors.WrapError(errors.ErrorTypeDatabase, "failed to update device home ID", err).
			WithOperation("UpdateDeviceHomeID").
			WithLayer("repository").
			WithContext("device_id", id).
			WithContext("home_id", homeID)
	}

	return nil
}
