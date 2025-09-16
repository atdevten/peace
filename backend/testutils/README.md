# Test Utilities

Thư mục này chứa các công cụ hỗ trợ cho việc testing trong dự án.

## Cấu trúc thư mục

```
testutils/
├── mocks/                    # Mock files được generate tự động
│   ├── repositories/         # Mock cho repository interfaces
│   ├── usecases/            # Mock cho use case interfaces
│   └── services/            # Mock cho service interfaces (tương lai)
├── helpers/                 # Helper functions cho testing
│   └── test_helpers.go      # Các hàm tiện ích tạo test data
└── README.md               # File này
```

## Cách sử dụng Mock

### 1. Generate Mock Files

```bash
# Generate tất cả mock files
make generate-mocks

# Hoặc chạy script trực tiếp
./scripts/generate-mocks.sh
```

### 2. Import Mock trong Test

```go
import (
    "github.com/atdevten/peace/testutils/mocks/repositories"
    "github.com/atdevten/peace/testutils/mocks/usecases"
    "github.com/atdevten/peace/testutils/helpers"
)
```

### 3. Sử dụng Mock trong Test

```go
func TestUserUseCase_GetByID(t *testing.T) {
    // Tạo mock repository
    mockRepo := repositories.NewMockUserRepository(t)
    
    // Setup mock behavior
    expectedUser := helpers.CreateTestUser()
    mockRepo.On("GetByID", mock.Anything, mock.Anything).Return(expectedUser, nil)
    
    // Tạo use case với mock
    useCase := NewUserUseCase(mockRepo)
    
    // Test
    user, err := useCase.GetByID(context.Background(), "user-id")
    
    // Assertions
    require.NoError(t, err)
    assert.NotNil(t, user)
    
    // Verify mock was called
    mockRepo.AssertExpectations(t)
}
```

## Helper Functions

### Tạo Test Data

```go
// Tạo test user
user := helpers.CreateTestUser()

// Tạo test user với thông tin tùy chỉnh
user := helpers.CreateTestGoogleUser()

// Tạo test entities khác
record := helpers.CreateTestMentalHealthRecord()
quote := helpers.CreateTestQuote()
tag := helpers.CreateTestTag()

// Tạo pointers
firstName := helpers.StringPtr("John")
age := helpers.IntPtr(25)
```

### Value Objects

```go
// Tạo value objects
userID := helpers.CreateTestUserID()
email := helpers.CreateTestEmail("test@example.com")
username := helpers.CreateTestUsername("testuser")
password := helpers.CreateTestPassword("Password123")
```

## Best Practices

### 1. Mock Setup

```go
// ✅ Tốt: Setup mock behavior rõ ràng
mockRepo.On("GetByID", mock.Anything, mock.Anything).Return(user, nil)

// ❌ Tránh: Setup quá phức tạp
mockRepo.On("GetByID", mock.MatchedBy(func(id *value_objects.UserID) bool {
    return id.String() == "specific-id"
}), mock.Anything).Return(user, nil)
```

### 2. Assertions

```go
// ✅ Tốt: Verify mock expectations
mockRepo.AssertExpectations(t)

// ✅ Tốt: Verify specific calls
mockRepo.AssertCalled(t, "GetByID", mock.Anything, mock.Anything)
mockRepo.AssertNumberOfCalls(t, "GetByID", 1)
```

### 3. Test Structure

```go
func TestSomething(t *testing.T) {
    tests := []struct {
        name        string
        input       string
        mockSetup   func(*repositories.MockUserRepository)
        wantErr     bool
        expectedErr string
    }{
        {
            name:  "success case",
            input: "valid-input",
            mockSetup: func(mock *repositories.MockUserRepository) {
                mock.On("SomeMethod", mock.Anything).Return(nil)
            },
            wantErr: false,
        },
        // ... more test cases
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            mockRepo := repositories.NewMockUserRepository(t)
            tt.mockSetup(mockRepo)
            
            // Test implementation
            // ...
            
            mockRepo.AssertExpectations(t)
        })
    }
}
```

## Commands

```bash
# Generate mocks
make generate-mocks

# Clean mocks
make clean-mocks

# Run unit tests
make test-unit

# Run integration tests
make test-integration

# Run tests with coverage
make test-coverage
```

## Lưu ý

1. **Không edit mock files thủ công** - Chúng được generate tự động từ interfaces
2. **Regenerate mocks** khi interface thay đổi
3. **Sử dụng helpers** để tạo test data thay vì hardcode
4. **Verify mock expectations** trong mọi test case
5. **Tách biệt unit tests và integration tests** bằng tags
