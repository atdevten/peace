#!/bin/bash

# Script to generate all mock files for the project
# Usage: ./scripts/generate-mocks.sh

set -e  # Exit on any error

echo "🚀 Generating all mock files..."

# Create mock directories if they don't exist
mkdir -p testutils/mocks/repositories
mkdir -p testutils/mocks/usecases
mkdir -p testutils/mocks/services

# Check if mockgen is installed
if ! command -v mockgen &> /dev/null; then
    echo "❌ mockgen is not installed. Installing..."
    go install go.uber.org/mock/mockgen@latest
fi

echo "📁 Generating repository mocks..."

# Generate repository mocks
mockgen -source=internal/domain/repositories/user_repository.go -destination=testutils/mocks/repositories/user_repository_mock.go
echo "✅ Generated repositories/user_repository_mock.go"

mockgen -source=internal/domain/repositories/mental_health_record_repository.go -destination=testutils/mocks/repositories/mental_health_record_repository_mock.go
echo "✅ Generated repositories/mental_health_record_repository_mock.go"

mockgen -source=internal/domain/repositories/quote_repository.go -destination=testutils/mocks/repositories/quote_repository_mock.go
echo "✅ Generated repositories/quote_repository_mock.go"

mockgen -source=internal/domain/repositories/tag_repository.go -destination=testutils/mocks/repositories/tag_repository_mock.go
echo "✅ Generated repositories/tag_repository_mock.go"

mockgen -source=internal/domain/repositories/user_online_status_repository.go -destination=testutils/mocks/repositories/user_online_status_repository_mock.go
echo "✅ Generated repositories/user_online_status_repository_mock.go"

echo "📁 Generating use case mocks..."

# Generate use case mocks
mockgen -source=internal/application/usecases/user_usecase.go -destination=testutils/mocks/usecases/user_usecase_mock.go
echo "✅ Generated usecases/user_usecase_mock.go"

mockgen -source=internal/application/usecases/auth_usecase.go -destination=testutils/mocks/usecases/auth_usecase_mock.go
echo "✅ Generated usecases/auth_usecase_mock.go"

mockgen -source=internal/application/usecases/mental_health_record_usecase.go -destination=testutils/mocks/usecases/mental_health_record_usecase_mock.go
echo "✅ Generated usecases/mental_health_record_usecase_mock.go"

mockgen -source=internal/application/usecases/quote_usecase.go -destination=testutils/mocks/usecases/quote_usecase_mock.go
echo "✅ Generated usecases/quote_usecase_mock.go"

mockgen -source=internal/application/usecases/tag_usecase.go -destination=testutils/mocks/usecases/tag_usecase_mock.go
echo "✅ Generated usecases/tag_usecase_mock.go"

mockgen -source=internal/application/usecases/user_online_status_usecase.go -destination=testutils/mocks/usecases/user_online_status_usecase_mock.go
echo "✅ Generated usecases/user_online_status_usecase_mock.go"
