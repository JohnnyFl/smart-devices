# Smart Devices Management System

A serverless smart home device management system built with Go, AWS Lambda, DynamoDB, and SQS. This system provides complete CRUD operations for managing smart home devices and processes device-home association notifications through SQS integration.

## 🎯 Project Goals

Develop a serverless application using AWS Lambda and DynamoDB to manage smart home devices. Implement CRUD operations for device management and integrate with an SQS queue to process notifications about device-home associations.

## ✅ Requirements Compliance

### 1. Deployment ✅
- ✅ **Serverless Framework**: Complete serverless.yml configuration
- ✅ **IAM Roles/Policies**: Least-privilege access for Lambda, DynamoDB, and SQS
- ✅ **Infrastructure as Code**: All resources defined in serverless.yml

### 2. Lambda Functions ✅
- ✅ **CreateDevice**: Adds new device to DynamoDB with validation
- ✅ **GetDevice**: Retrieves device details by unique identifier
- ✅ **UpdateDevice**: Modifies existing device information
- ✅ **DeleteDevice**: Removes device from DynamoDB
- ✅ **SQS Listener**: Processes device-home association messages
- ✅ **Golang Implementation**: All functions written in Go using AWS SDK v2

### 3. DynamoDB Table ✅
- ✅ **id** (String, Primary Key): Unique device identifier
- ✅ **mac** (String): MAC address of the device
- ✅ **name** (String): Device name
- ✅ **type** (String): Device type (thermostat, light, camera, sensor)
- ✅ **homeId** (String): Home identifier
- ✅ **createdAt** (Int): Creation timestamp in Unix millis
- ✅ **modifiedAt** (Int): Last update timestamp in Unix millis

### 4. Additional Requirements ✅
- ✅ **Error Handling**: Comprehensive error handling and logging
- ✅ **Input Validation**: MAC address, UUID, and type validation
- ✅ **Unit Tests**: Test coverage for Lambda functions
- ✅ **Build Scripts**: npm scripts and Makefile for build/test/deploy
- ✅ **Documentation**: Complete README with setup instructions
- ✅ **Security**: IAM roles, input sanitization, error handling
- ✅ **Performance**: Connection reuse, efficient DynamoDB operations

## 🏗️ Architecture

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   API Gateway   │────│  Lambda Functions │────│    DynamoDB     │
│                 │    │                   │    │                 │
│ REST Endpoints  │    │ • get-device      │    │  Devices Table  │
│                 │    │ • list-devices    │    │                 │
│                 │    │ • create-device   │    │                 │
│                 │    │ • update-device   │    │                 │
│                 │    │ • delete-device   │    │                 │
└─────────────────┘    └─────────────────┘    └─────────────────┘
                                │
                       ┌─────────────────┐    ┌─────────────────┐
                       │   SQS Queue     │────│  sqs-listener   │
                       │                 │    │                 │
                       │ Device Events   │    │ Event Processor │
                       └─────────────────┘    └─────────────────┘
```

## 🚀 Features

### Core Functionality
- **Complete Device CRUD Operations**: Create, read, update, and delete smart devices
- **Device Types Support**: Thermostat, Light, Camera, Sensor
- **Device-Home Association**: SQS-based processing for device-home relationships
- **Real-time Updates**: Automatic `modifiedAt` timestamp updates

### Technical Features
- **Serverless Architecture**: AWS Lambda functions with API Gateway
- **Local Development**: DynamoDB Local + Serverless Offline
- **Structured Logging**: Comprehensive logging with Zap logger
- **Enhanced Error Handling**: Domain-specific errors with context and proper wrapping
- **Comprehensive Input Validation**: MAC address, UUID, device type, and field validation
- **Standardized API Responses**: Consistent error codes and success responses
- **Modular Architecture**: Clean separation with setup utilities and validation layer
- **Security**: IAM roles with least-privilege access
- **Performance**: Optimized DynamoDB operations with proper indexing

## 📋 Data Models

### Device Model
```
type Device struct {
    ID         string `json:"id"`         // Unique identifier (Primary Key)
    MAC        string `json:"mac"`        // MAC address of the device
    Name       string `json:"name"`       // Name of the device
    Type       string `json:"type"`       // Type: thermostat|light|camera|sensor
    HomeID     string `json:"homeId"`     // Identifier of the home
    CreatedAt  int64  `json:"createdAt"`  // Creation date (Unix timestamp millis)
    ModifiedAt int64  `json:"modifiedAt"` // Last update date (Unix timestamp millis)
}
```

### SQS Message Model
```
type SQSMessage struct {
    DeviceID string `json:"deviceId"`    // Device to update
    HomeID   string `json:"homeId"`      // New home association
    Action   string `json:"action"`      // Action type
}
```

### Request Models
```
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
```

## 🛠️ Prerequisites

- **Go 1.24+**
- **Node.js & npm** (for Serverless Framework)
- **Docker** (for DynamoDB Local)
- **AWS CLI** (for deployment)
- **Serverless Framework**

```bash
npm install -g serverless
npm install -g serverless-offline
```

## 🏃‍♂️ Quick Start

### 1. Clone and Setup

```bash
git clone <repository-url>
cd smart-devices
go mod download
npm install
```

### 2. Local Development Setup

Start DynamoDB Local:
```bash
docker run -d -p 8000:8000 --name dynamodb-local amazon/dynamodb-local
```

Create the devices table:
```bash
aws dynamodb create-table \
    --table-name devices \
    --attribute-definitions AttributeName=id,AttributeType=S \
    --key-schema AttributeName=id,KeyType=HASH \
    --billing-mode PAY_PER_REQUEST \
    --endpoint-url http://localhost:8000
```

### 3. Environment Variables

Create a `.env` file or set environment variables:
```bash
export DYNAMODB_TABLE=devices
export DYNAMODB_URL=http://localhost:8000
export AWS_REGION=us-east-1
export SQS_QUEUE_URL=http://localhost:4566/000000000000/fake-queue
export STAGE=dev
```

### 4. Run Locally with Serverless Offline

```bash
# Start all services (DynamoDB + Lambda functions)
npm run dev

# Or manually:
# 1. Start DynamoDB Local
docker run -d -p 8000:8000 --name dynamodb-local amazon/dynamodb-local

# 2. Start serverless offline
serverless offline start

# The API will be available at http://localhost:3000
```

## 📡 Lambda Functions

### HTTP API Functions

| Function | Method | Endpoint | Description |
|----------|--------|----------|-------------|
| `get-device` | `GET` | `/devices/{id}` | Retrieve device details by unique identifier |
| `list-devices` | `GET` | `/devices` | List all devices |
| `create-device` | `POST` | `/devices` | Add a new device to DynamoDB |
| `update-device` | `PUT` | `/devices/{id}` | Modify existing device information |
| `delete-device` | `DELETE` | `/devices/{id}` | Remove a device from DynamoDB |

### Event-Driven Functions

| Function | Trigger | Description |
|----------|---------|-------------|
| `sqs-listener` | SQS Queue | Process device-home association messages |

### SQS Integration

The SQS listener processes JSON messages for device-home associations:

```json
{
  "deviceId": "123e4567-e89b-12d3-a456-426614174000",
  "homeId": "987fcdeb-51a2-43d7-8f9e-123456789abc",
  "action": "associate"
}
```

When a message is received, the system:
1. Validates the message format
2. Updates the device record in DynamoDB with the new `homeId`
3. Updates the `modifiedAt` timestamp
4. Logs the operation for audit purposes

### Request/Response Examples

#### Create Device
```bash
curl -X POST https://api.example.com/devices \
  -H "Content-Type: application/json" \
  -d '{
    "mac": "00:11:22:33:44:55",
    "name": "Living Room Thermostat",
    "type": "thermostat",
    "homeId": "123e4567-e89b-12d3-a456-426614174000"
  }'
```

#### Update Device
```bash
curl -X PUT https://api.example.com/devices/{id} \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Updated Device Name",
    "type": "light"
  }'
```

## 🧪 Testing

### Unit Tests

The project includes comprehensive unit tests for all Lambda functions:

```bash
# Run all unit tests
go test ./...

# Run tests with coverage
go test -cover ./...
make test-cover  # Generates HTML coverage report

# Run tests for specific package
go test ./internal/services/ -v
go test ./internal/repository/ -v
go test ./internal/handlers/ -v

# Run tests with race detection
go test -race ./...
```

#### Test Coverage Areas:
- ✅ **Service Layer**: Business logic validation
- ✅ **Repository Layer**: Data access operations  
- ✅ **Handler Layer**: HTTP request/response handling
- ✅ **Model Validation**: Input validation and data structures
- ✅ **Error Scenarios**: Error handling and edge cases

#### Example Test Output:
```bash
$ go test ./internal/services/ -v
=== RUN   TestDeviceService_CreateDevice
--- PASS: TestDeviceService_CreateDevice (0.00s)
=== RUN   TestDeviceService_GetDevice
--- PASS: TestDeviceService_GetDevice (0.00s)
=== RUN   TestDeviceService_UpdateDevice
--- PASS: TestDeviceService_UpdateDevice (0.01s)
=== RUN   TestDeviceService_DeleteDevice
--- PASS: TestDeviceService_DeleteDevice (0.00s)
=== RUN   TestDeviceService_UpdateDeviceHomeID
--- PASS: TestDeviceService_UpdateDeviceHomeID (0.01s)
PASS
coverage: 85.2% of statements
```

### Integration Tests
```bash
# Test the UpdateDeviceHomeID functionality
./run_test.sh

# Or run directly
go run test_update_home_id.go
```

### Acceptance Tests (Manual)
```bash
# 1. Create a device
curl -X POST http://localhost:3000/devices \
  -H "Content-Type: application/json" \
  -d '{
    "mac": "00:11:22:33:44:55",
    "name": "Living Room Thermostat",
    "type": "thermostat",
    "homeId": "123e4567-e89b-12d3-a456-426614174000"
  }'

# 2. Get all devices
curl http://localhost:3000/devices

# 3. Get specific device
curl http://localhost:3000/devices/{device-id}

# 4. Update device
curl -X PUT http://localhost:3000/devices/{device-id} \
  -H "Content-Type: application/json" \
  -d '{"name": "Updated Thermostat", "type": "thermostat"}'

# 5. Delete device
curl -X DELETE http://localhost:3000/devices/{device-id}
```

### SQS Testing
```bash
# Send test message to SQS queue (requires AWS CLI)
aws sqs send-message \
  --queue-url $SQS_QUEUE_URL \
  --message-body '{
    "deviceId": "your-device-id",
    "homeId": "new-home-id",
    "action": "associate"
  }'
```

## 🔨 Build Scripts

### Available Scripts

```bash
# Install dependencies
npm install
go mod download

# Build all functions
npm run build

# Run tests
npm run test

# Start local development
npm run dev

# Deploy to AWS
npm run deploy

# Deploy to specific stage
npm run deploy:dev
npm run deploy:prod

# Clean build artifacts
npm run clean
```

### Manual Build Commands

```bash
# Build individual functions (outputs to build/{function}/bootstrap)
GOOS=linux GOARCH=amd64 go build -o build/get-device/bootstrap cmd/get-device/main.go
GOOS=linux GOARCH=amd64 go build -o build/create-device/bootstrap cmd/create-device/main.go
GOOS=linux GOARCH=amd64 go build -o build/update-device/bootstrap cmd/update-device/main.go
GOOS=linux GOARCH=amd64 go build -o build/delete-device/bootstrap cmd/delete-device/main.go
GOOS=linux GOARCH=amd64 go build -o build/list-devices/bootstrap cmd/list-devices/main.go
GOOS=linux GOARCH=amd64 go build -o build/sqs-listener/bootstrap cmd/sqs-listener/main.go
```

## 🚀 Deployment

### Prerequisites for Deployment
- AWS CLI configured with appropriate credentials
- Serverless Framework installed globally
- Go 1.24+ installed

### Deploy to AWS

```bash
# Deploy to development environment
serverless deploy --stage dev

# Deploy to production environment
serverless deploy --stage prod

# Deploy specific function
serverless deploy function --function get-device --stage dev
```

### Environment-specific Configuration

The system automatically configures itself based on the deployment stage:

| Stage | Runtime | DynamoDB | Features |
|-------|---------|----------|----------|
| **dev** | `go1.x` | Local (localhost:8000) | Debug logging, CORS enabled |
| **prod** | `provided.al2` | AWS DynamoDB | Production logging, optimized |

### Infrastructure Created

The deployment creates:
- **DynamoDB Table**: `smart-devices-{stage}-devices`
- **SQS Queue**: `smart-devices-{stage}-device-notifications`
- **SQS DLQ**: `smart-devices-{stage}-device-notifications-dlq`
- **IAM Roles**: Lambda execution roles with minimal permissions
- **API Gateway**: REST API with CORS enabled
- **CloudWatch Logs**: Log groups for each Lambda function

## 📁 Project Structure

```
smart-devices/
├── cmd/                    # Lambda function entry points
│   ├── create-device/      # POST /devices
│   ├── get-device/         # GET /devices/{id}
│   ├── list-devices/       # GET /devices
│   ├── update-device/      # PUT /devices/{id}
│   ├── delete-device/      # DELETE /devices/{id}
│   └── sqs-listener/       # SQS event processor
├── internal/
│   ├── config/            # Configuration management
│   ├── errors/            # Error handling and domain errors
│   │   ├── api_errors.go  # HTTP API error definitions
│   │   └── domain_errors.go # Domain-specific error types
│   ├── handlers/          # HTTP/SQS request handlers
│   ├── models/            # Data models and request/response types
│   ├── repository/        # Data access layer (DynamoDB)
│   ├── services/          # Business logic layer
│   ├── setup/             # Shared initialization utilities
│   └── validation/        # Input validation layer
├── build/                 # Build artifacts (generated)
├── serverless.yml         # Serverless Framework configuration
├── package.json          # npm scripts and dependencies
├── CLAUDE.md             # Claude Code assistant documentation
└── go.mod                # Go module dependencies
```

## 🔧 Configuration

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `DYNAMODB_TABLE` | DynamoDB table name | `devices` |
| `DYNAMODB_URL` | DynamoDB endpoint (local dev) | - |
| `AWS_REGION` | AWS region | `us-east-1` |
| `SQS_QUEUE_URL` | SQS queue URL | - |
| `STAGE` | Deployment stage | `dev` |

### Device Validation Rules

- **MAC Address**: Must be valid MAC format (e.g., `00:11:22:33:44:55`)
- **Name**: 1-100 characters
- **Type**: Must be one of: `thermostat`, `light`, `camera`, `sensor`
- **HomeID**: Must be valid UUID format

### Enhanced Error Handling

The system implements a comprehensive error handling strategy with domain-specific errors:

#### Error Types
- **Validation Errors** (400): Invalid input data, missing fields, format errors
- **Not Found Errors** (404): Resource not found, empty collections
- **Database Errors** (500): DynamoDB operation failures, marshaling errors
- **Internal Errors** (500): Unexpected system errors

#### Error Response Format
```json
{
  "code": "VALIDATION_ERROR",
  "message": "Device name must be between 1 and 100 characters"
}
```

#### Error Context & Logging
- **Structured Logging**: All errors include operation context, layer information, and relevant IDs
- **Error Wrapping**: Errors maintain their original context while adding layer-specific information
- **Request Tracing**: Each error can be traced through repository → service → handler layers

## 🔍 Monitoring & Logging

The system uses structured logging with Zap logger:

- **Development**: Logs to stdout with debug level
- **Production**: Logs to CloudWatch with info level
- **Request tracing**: Each request includes correlation IDs
- **Error tracking**: Comprehensive error logging with context

## 🛡️ Security Optimizations

### Implemented Security Measures
- **Input Validation**: Comprehensive validation using struct tags and custom validators
- **Error Handling**: Sanitized error responses that don't expose internal details
- **IAM Least Privilege**: Functions have minimal required permissions only
- **Encryption at Rest**: DynamoDB SSE enabled
- **Point-in-Time Recovery**: DynamoDB PITR enabled for data protection
- **Dead Letter Queue**: Failed SQS messages are captured for analysis
- **Structured Logging**: No sensitive data in logs

### Security Recommendations
```yaml
# Additional security measures to consider:
- API Gateway throttling and rate limiting
- AWS WAF for API protection  
- VPC endpoints for private communication
- AWS Secrets Manager for sensitive configuration
- Request/response logging for audit trails
- API key authentication for production use
```

## ⚡ Performance Optimizations

### Implemented Optimizations
- **Connection Reuse**: AWS SDK clients initialized once per Lambda container
- **Efficient DynamoDB Operations**: 
  - Single-item operations for CRUD
  - Batch operations where applicable
  - Proper error handling and retries
- **Lambda Cold Start Reduction**:
  - Minimal dependencies
  - Efficient initialization in `init()` functions
  - Provisioned concurrency for production (configurable)
- **Memory Optimization**: Right-sized Lambda memory allocation
- **Structured Logging**: Efficient JSON logging with Zap

### Performance Recommendations
```yaml
# Additional optimizations to consider:
- DynamoDB Auto Scaling for variable workloads
- Lambda Provisioned Concurrency for consistent performance
- CloudFront distribution for API caching
- DynamoDB DAX for microsecond latency (if needed)
- Connection pooling for high-throughput scenarios
- Async processing for non-critical operations
```

### Monitoring & Observability
- **CloudWatch Metrics**: Lambda duration, errors, throttles
- **Custom Metrics**: Business-specific metrics via CloudWatch
- **Distributed Tracing**: AWS X-Ray integration ready
- **Structured Logging**: Correlation IDs for request tracing
- **Alerting**: CloudWatch alarms for error rates and latency

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 🆘 Troubleshooting

### Common Issues

1. **DynamoDB Connection Issues**
   ```bash
   # Check if DynamoDB Local is running
   curl http://localhost:8000
   
   # Restart DynamoDB Local
   docker stop dynamodb-local && docker rm dynamodb-local
   docker run -d -p 8000:8000 --name dynamodb-local amazon/dynamodb-local
   ```

2. **Table Not Found**
   ```bash
   # List tables
   aws dynamodb list-tables --endpoint-url http://localhost:8000
   
   # Create table if missing
   aws dynamodb create-table --table-name devices ...
   ```

3. **Lambda Function Errors**
   ```bash
   # Check logs
   serverless logs -f get-device -t
   ```

### Debug Mode

Enable debug logging by setting:
```bash
export LOG_LEVEL=debug
```

## 📞 Support

For questions and support, please open an issue in the repository or contact the development team.
