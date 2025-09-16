package redis

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewRealClient(t *testing.T) {
	tests := []struct {
		name        string
		addr        string
		password    string
		db          int
		wantErr     bool
		expectedErr string
	}{
		{
			name:     "create client with valid config",
			addr:     "localhost:6379",
			password: "",
			db:       0,
			wantErr:  false,
		},
		{
			name:     "create client with invalid host",
			addr:     "invalid_host:6379",
			password: "",
			db:       0,
			wantErr:  false, // NewRealClient doesn't validate connection
		},
		{
			name:     "create client with invalid port",
			addr:     "localhost:9999",
			password: "",
			db:       0,
			wantErr:  false, // NewRealClient doesn't validate connection
		},
		{
			name:     "create client with password",
			addr:     "localhost:6379",
			password: "password",
			db:       0,
			wantErr:  false,
		},
		{
			name:     "create client with different database",
			addr:     "localhost:6379",
			password: "",
			db:       1,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewRealClient(tt.addr, tt.password, tt.db)

			require.NotNil(t, client)

			// Test that we can ping Redis
			ctx := context.Background()
			err := client.Ping(ctx)
			if err != nil {
				// If Redis is not available, skip the test
				t.Skip("Skipping test - Redis not available")
			}
			require.NoError(t, err)

			// Close the client
			err = client.Close()
			require.NoError(t, err)
		})
	}
}

func TestNewMockClient(t *testing.T) {
	client := NewMockClient()

	require.NotNil(t, client)

	// Test that we can call methods without error
	ctx := context.Background()
	err := client.Ping(ctx)
	require.NoError(t, err)

	err = client.Set(ctx, "test_key", "test_value", 0)
	require.NoError(t, err)

	val, err := client.Get(ctx, "test_key")
	require.NoError(t, err)
	assert.Empty(t, val) // Mock returns empty string

	err = client.Close()
	require.NoError(t, err)
}

func TestRedisClient_Operations(t *testing.T) {
	// Skip if Redis is not available
	client := NewRealClient("localhost:6379", "", 0)
	defer client.Close()

	// Test ping first to check if Redis is available
	ctx := context.Background()
	err := client.Ping(ctx)
	if err != nil {
		t.Skip("Skipping test - Redis not available")
	}

	// Test set and get
	err = client.Set(ctx, "test_key", "test_value", 0)
	require.NoError(t, err)

	val, err := client.Get(ctx, "test_key")
	require.NoError(t, err)
	assert.Equal(t, "test_value", val)

	// Test set with expiration
	err = client.Set(ctx, "test_key_expire", "test_value_expire", time.Second)
	require.NoError(t, err)

	val, err = client.Get(ctx, "test_key_expire")
	require.NoError(t, err)
	assert.Equal(t, "test_value_expire", val)

	// Wait for expiration
	time.Sleep(2 * time.Second)

	val, err = client.Get(ctx, "test_key_expire")
	require.Error(t, err)
	assert.Empty(t, val)

	// Test delete
	err = client.Del(ctx, "test_key")
	require.NoError(t, err)

	val, err = client.Get(ctx, "test_key")
	require.Error(t, err)
	assert.Empty(t, val)

	// Test exists
	exists, err := client.Exists(ctx, "test_key")
	require.NoError(t, err)
	assert.Equal(t, int64(0), exists)

	// Test set and exists
	err = client.Set(ctx, "test_key_exists", "test_value_exists", 0)
	require.NoError(t, err)

	exists, err = client.Exists(ctx, "test_key_exists")
	require.NoError(t, err)
	assert.Equal(t, int64(1), exists)

	// Clean up
	client.Del(ctx, "test_key_exists")
}

func TestRedisClient_SetOperations(t *testing.T) {
	// Skip if Redis is not available
	client := NewRealClient("localhost:6379", "", 0)
	defer client.Close()

	// Test ping first to check if Redis is available
	ctx := context.Background()
	err := client.Ping(ctx)
	if err != nil {
		t.Skip("Skipping test - Redis not available")
	}

	// Test sadd and smembers
	err = client.SAdd(ctx, "test_set", "member1", "member2", "member3")
	require.NoError(t, err)

	members, err := client.SMembers(ctx, "test_set")
	require.NoError(t, err)
	assert.Len(t, members, 3)
	assert.Contains(t, members, "member1")
	assert.Contains(t, members, "member2")
	assert.Contains(t, members, "member3")

	// Test sismember
	isMember, err := client.SIsMember(ctx, "test_set", "member1")
	require.NoError(t, err)
	assert.True(t, isMember)

	isMember, err = client.SIsMember(ctx, "test_set", "member4")
	require.NoError(t, err)
	assert.False(t, isMember)

	// Test srem
	err = client.SRem(ctx, "test_set", "member1")
	require.NoError(t, err)

	isMember, err = client.SIsMember(ctx, "test_set", "member1")
	require.NoError(t, err)
	assert.False(t, isMember)

	// Test scard
	count, err := client.SCard(ctx, "test_set")
	require.NoError(t, err)
	assert.Equal(t, int64(2), count)

	// Clean up
	client.Del(ctx, "test_set")
}

func TestRedisClient_MGet(t *testing.T) {
	// Skip if Redis is not available
	client := NewRealClient("localhost:6379", "", 0)
	defer client.Close()

	// Test ping first to check if Redis is available
	ctx := context.Background()
	err := client.Ping(ctx)
	if err != nil {
		t.Skip("Skipping test - Redis not available")
	}

	// Set multiple keys
	err = client.Set(ctx, "key1", "value1", 0)
	require.NoError(t, err)

	err = client.Set(ctx, "key2", "value2", 0)
	require.NoError(t, err)

	err = client.Set(ctx, "key3", "value3", 0)
	require.NoError(t, err)

	// Test MGet
	values, err := client.MGet(ctx, "key1", "key2", "key3", "nonexistent")
	require.NoError(t, err)
	assert.Len(t, values, 4)
	assert.Equal(t, "value1", values[0])
	assert.Equal(t, "value2", values[1])
	assert.Equal(t, "value3", values[2])
	assert.Equal(t, "", values[3]) // nonexistent key returns empty string

	// Clean up
	client.Del(ctx, "key1", "key2", "key3")
}
