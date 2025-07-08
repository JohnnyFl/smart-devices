package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"example.com/smart-devices/internal/models"
	"example.com/smart-devices/internal/repository"
	"example.com/smart-devices/internal/services"
	"example.com/smart-devices/utils"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// Simple test to verify UpdateDeviceHomeID updates both homeId and modifiedAt
func main() {
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()

	// Set environment variables for local DynamoDB
	os.Setenv("IS_LOCAL", "true")
	os.Setenv("DYNAMODB_ENDPOINT", "http://localhost:8000")

	// Set AWS region (required even for local)
	os.Setenv("AWS_REGION", "us-east-1")

	// Load AWS config
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("us-east-1"),
		// For local testing, we can use dummy credentials
		config.WithCredentialsProvider(aws.CredentialsProviderFunc(func(ctx context.Context) (aws.Credentials, error) {
			return aws.Credentials{
				AccessKeyID:     "dummy",
				SecretAccessKey: "dummy",
			}, nil
		})),
	)
	if err != nil {
		log.Fatal("Failed to load AWS config:", err)
	}

	// Initialize DynamoDB client
	client := utils.GetDynamoDBClient(cfg)

	// Use your actual table name - you might need to adjust this
	tableName := "devices" // Change this to your actual table name
	repo := repository.NewDeviceRepository(client, tableName, logger)
	service := services.NewDeviceService(repo, logger)

	ctx := context.Background()

	// First, let's create a test device
	testDevice := models.Device{
		MAC:    "00:11:22:33:44:55",
		Name:   "Test Device",
		Type:   "thermostat",
		HomeID: "original-home-id",
	}

	fmt.Println("Creating test device...")
	createdDevice, err := service.CreateDevice(ctx, testDevice)
	if err != nil {
		log.Printf("Failed to create test device: %v", err)
		return
	}

	deviceID := createdDevice.ID
	newHomeID := "new-home-id-" + uuid.New().String()[:8]

	fmt.Printf("Created device with ID: %s\n", deviceID)
	fmt.Printf("Before update - HomeID: %s, ModifiedAt: %d\n",
		createdDevice.HomeID, createdDevice.ModifiedAt)

	// Wait a second to ensure timestamp difference
	time.Sleep(2 * time.Second)

	// Update the device's home ID
	fmt.Println("Updating device home ID...")
	err = service.UpdateDeviceHomeID(ctx, deviceID, newHomeID)
	if err != nil {
		log.Printf("Failed to update device: %v", err)
		return
	}

	// Get device after update
	deviceAfter, err := service.GetDevice(ctx, deviceID)
	if err != nil {
		log.Printf("Could not get device after update: %v", err)
		return
	}

	fmt.Printf("After update - HomeID: %s, ModifiedAt: %d\n",
		deviceAfter.HomeID, deviceAfter.ModifiedAt)

	// Verify both fields were updated
	if deviceAfter.HomeID == newHomeID {
		fmt.Println("✅ HomeID was updated successfully")
	} else {
		fmt.Println("❌ HomeID was not updated")
	}

	if deviceAfter.ModifiedAt > createdDevice.ModifiedAt {
		fmt.Println("✅ ModifiedAt was updated successfully")
		fmt.Printf("   Time difference: %d seconds\n", deviceAfter.ModifiedAt-createdDevice.ModifiedAt)
	} else {
		fmt.Println("❌ ModifiedAt was not updated")
	}

	// Clean up - delete the test device
	fmt.Println("Cleaning up test device...")
	err = service.DeleteDevice(ctx, deviceID)
	if err != nil {
		log.Printf("Failed to delete test device: %v", err)
	} else {
		fmt.Println("✅ Test device cleaned up successfully")
	}
}
