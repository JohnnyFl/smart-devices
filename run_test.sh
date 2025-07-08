#!/bin/bash

echo "🚀 Testing UpdateDeviceHomeID with ModifiedAt update..."
echo ""

# Check if DynamoDB Local is running
if ! curl -s http://localhost:8000 > /dev/null 2>&1; then
    echo "❌ DynamoDB Local is not running on localhost:8000"
    echo "Please start DynamoDB Local with Docker:"
    echo "docker run -p 8000:8000 amazon/dynamodb-local"
    exit 1
fi

echo "✅ DynamoDB Local is running"
echo ""

# Run the test
echo "Running test..."
go run test_update_home_id.go

echo ""
echo "Test completed!"
