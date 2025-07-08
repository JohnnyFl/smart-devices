# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Development Commands

### Build Commands
- `npm run build` - Build all Lambda functions for deployment
- `make build` - Alternative build command using Makefile
- `./build.sh` - Build and zip all Lambda functions for AWS deployment

### Test Commands
- `npm test` or `go test ./...` - Run all unit tests
- `npm run test:coverage` or `make test-cover` - Run tests with coverage report
- `npm run test:integration` or `./run_test.sh` - Run integration tests
- `make test-integration` - Run integration tests using Makefile

### Development Environment
- `npm run dev` - Start complete local development environment (DynamoDB + serverless offline)
- `make setup` - Set up local development environment from scratch
- `make dev` - Start serverless offline after setup
- `docker-compose up -d` - Start DynamoDB Local only

### Deployment Commands
- `npm run deploy` or `make deploy` - Deploy to AWS dev stage
- `npm run deploy:prod` or `make deploy-prod` - Deploy to AWS prod stage
- `serverless deploy --stage <stage>` - Deploy to specific stage

### Linting & Formatting
- `go fmt ./...` - Format Go code
- `go vet ./...` - Vet Go code for issues

## Architecture Overview

This is a serverless smart home device management system built with:

### Technology Stack
- **Backend**: Go 1.24+ with AWS Lambda
- **Database**: DynamoDB with single table design
- **Queue**: SQS for async processing
- **Framework**: Serverless Framework
- **Local Development**: DynamoDB Local + Docker

### Project Structure
```
smart-devices/
├── cmd/                    # Lambda function entry points (main.go files)
│   ├── create-device/
│   ├── get-device/
│   ├── list-devices/
│   ├── update-device/
│   ├── delete-device/
│   └── sqs-listener/
├── internal/               # Core business logic
│   ├── config/            # Configuration management
│   ├── handlers/          # HTTP/SQS request handlers
│   ├── models/            # Data models and structs
│   ├── repository/        # Data access layer (DynamoDB)
│   └── services/          # Business logic layer
├── utils/                 # Utility functions (DynamoDB client, responses)
├── serverless.yml         # Serverless Framework configuration
├── go.mod                 # Go module definition
└── package.json          # Node.js scripts and dependencies
```

### Key Design Patterns
- **Layered Architecture**: Handler → Service → Repository pattern
- **Dependency Injection**: Services accept repository interfaces for testability
- **Interface-based Design**: Repository interfaces enable mocking in tests
- **Structured Logging**: Zap logger with consistent log levels and context

### Lambda Functions
- **get-device**: GET /devices/{id} - Retrieve single device
- **list-devices**: GET /devices - List all devices  
- **create-device**: POST /devices - Create new device
- **update-device**: PUT /devices/{id} - Update device
- **delete-device**: DELETE /devices/{id} - Delete device
- **sqs-listener**: SQS trigger - Process device-home association events

### Data Models
- **Device**: Core entity with id, mac, name, type, homeId, createdAt, modifiedAt
- **CreateDeviceRequest**: Input validation for device creation
- **UpdateDeviceRequest**: Input validation for device updates with optional fields
- **SQSMessage**: Structure for SQS event processing

### Environment Configuration
- **DYNAMODB_TABLE**: DynamoDB table name (default: "devices")
- **DYNAMODB_URL**: Local DynamoDB endpoint for development
- **SQS_QUEUE_URL**: SQS queue URL for device notifications
- **AWS_REGION**: AWS region (default: "us-east-1")
- **STAGE**: Deployment stage (dev/prod)

## Local Development Setup

1. **Install dependencies**:
   ```bash
   go mod download
   npm install
   ```

2. **Start local services**:
   ```bash
   npm run dev:setup  # Starts DynamoDB Local and creates table
   npm run dev:start  # Starts serverless offline
   ```

3. **Test endpoints**:
   - API available at `http://localhost:3000`
   - DynamoDB Local at `http://localhost:8000`

## Testing Strategy

### Unit Tests
- Located in `internal/services/device_service_test.go`
- Uses mock repositories for isolated testing
- Tests business logic without external dependencies

### Integration Tests
- `test_update_home_id.go` - Tests SQS functionality
- `./run_test.sh` - Integration test runner

### Test Coverage
- Run `make test-cover` to generate HTML coverage reports
- Target coverage areas: services, repository, handlers, models

## Build & Deployment

### Build Process
1. **Local builds**: `npm run build` creates Linux binaries in `bin/`
2. **AWS deployment**: `./build.sh` creates zip packages in `build/`
3. **Serverless deployment**: Uses artifacts from `build/` directory

### Deployment Stages
- **dev**: Uses `go1.x` runtime for local development
- **prod**: Uses `provided.al2` runtime for optimized performance

## AWS Resources Created
- **DynamoDB Table**: `smart-devices-{stage}-devices`
- **SQS Queue**: `smart-devices-{stage}-device-notifications`
- **SQS DLQ**: `smart-devices-{stage}-device-notifications-dlq`
- **Lambda Functions**: 6 functions for CRUD operations and SQS processing
- **API Gateway**: REST API with CORS enabled
- **IAM Roles**: Least-privilege access for each function

## Common Development Patterns

### Error Handling
- Repository layer returns wrapped errors with context
- Service layer adds business logic validation
- Handler layer returns appropriate HTTP status codes

### Logging
- Use structured logging with Zap
- Include context like device_id, operation, layer
- Log at appropriate levels (Debug, Info, Warn, Error)

### DynamoDB Operations
- Use attributevalue package for marshaling/unmarshaling
- Implement proper error handling for not found cases
- Update ModifiedAt timestamps on all updates

### SQS Message Processing
- Validate message structure before processing
- Update device homeId via repository layer
- Handle failures gracefully with DLQ

## Code Quality Guidelines

### Go Best Practices
- Use interfaces for dependency injection
- Keep functions focused and testable
- Follow Go naming conventions
- Use context.Context for cancellation

### Testing Requirements
- Unit tests for all service methods
- Mock external dependencies
- Test error scenarios
- Integration tests for SQS flows

### Performance Considerations
- Reuse AWS SDK clients across Lambda invocations
- Use single-item DynamoDB operations where possible
- Implement proper connection pooling for high throughput