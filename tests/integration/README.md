# Integration Tests

This directory contains integration tests for the Reolink Server. Integration tests verify that different components work together correctly.

## Overview

Integration tests differ from unit tests in that they:
- Test multiple components working together
- Use real database connections (test database)
- Make actual HTTP requests to the server
- Test end-to-end workflows

## Prerequisites

Before running integration tests, you need:

1. **PostgreSQL Test Database**
   ```bash
   createdb reolink_server_test
   psql -d reolink_server_test -c "CREATE EXTENSION IF NOT EXISTS timescaledb;"
   ```

2. **Redis Instance**
   - Use a separate Redis database (DB 1) for tests
   - Or run a separate Redis instance on a different port

3. **Environment Variables**
   ```bash
   export TEST_DB_HOST=localhost
   export TEST_DB_NAME=reolink_server_test
   export TEST_DB_USER=reolink
   export TEST_DB_PASSWORD=password
   export TEST_REDIS_HOST=localhost
   ```

## Running Integration Tests

### Run All Integration Tests

```bash
go test -v -tags=integration ./tests/integration/...
```

### Run Specific Test

```bash
go test -v -tags=integration ./tests/integration/... -run TestCameraLifecycle
```

### Run with Coverage

```bash
go test -v -tags=integration -coverprofile=coverage.out ./tests/integration/...
go tool cover -html=coverage.out
```

## Test Structure

### Test Files

- `camera_integration_test.go` - Camera management integration tests
- `event_integration_test.go` - Event processing integration tests (TODO)
- `stream_integration_test.go` - Stream proxy integration tests (TODO)
- `auth_integration_test.go` - Authentication integration tests (TODO)

### Helper Functions

Each test file includes helper functions:
- `setupTestServer(t)` - Creates a test server with all dependencies
- `TestServer.Request()` - Makes HTTP requests to the test server
- `TestServer.Login()` - Authenticates and returns JWT token
- `TestServer.Cleanup()` - Cleans up resources after tests

## Writing Integration Tests

### Example Test

```go
// +build integration

package integration

import (
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestMyFeature(t *testing.T) {
    // Setup
    testServer := setupTestServer(t)
    defer testServer.Cleanup()

    // Get auth token
    token := testServer.Login(t, "admin", "admin")

    // Make request
    resp := testServer.Request(t, "GET", "/api/v1/cameras", token, nil)
    
    // Assert response
    assert.Equal(t, http.StatusOK, resp.StatusCode)
    
    // Parse and verify response body
    var result map[string]interface{}
    err := json.NewDecoder(resp.Body).Decode(&result)
    require.NoError(t, err)
    
    // Additional assertions
    assert.NotNil(t, result["cameras"])
}
```

### Best Practices

1. **Use Build Tags**: Always include `// +build integration` at the top
2. **Clean Up**: Always defer `testServer.Cleanup()` to clean up resources
3. **Isolation**: Each test should be independent and not rely on other tests
4. **Test Data**: Create test data at the beginning of each test
5. **Assertions**: Use `require` for critical checks, `assert` for non-critical
6. **Descriptive Names**: Use clear, descriptive test names

## Test Coverage

Current integration test coverage:

- [x] Camera lifecycle (add, get, update, delete)
- [x] Event flow (create, list, acknowledge)
- [ ] WebSocket event streaming
- [ ] Authentication and authorization
- [ ] Stream proxy (FLV, HLS)
- [ ] Recording management
- [ ] PTZ control
- [ ] Camera configuration

## Continuous Integration

Integration tests can be run in CI/CD pipelines:

### GitHub Actions Example

```yaml
name: Integration Tests

on: [push, pull_request]

jobs:
  integration-test:
    runs-on: ubuntu-latest
    
    services:
      postgres:
        image: timescale/timescaledb:latest-pg16
        env:
          POSTGRES_DB: reolink_server_test
          POSTGRES_USER: reolink
          POSTGRES_PASSWORD: password
        options: >-
          --health-cmd pg_isready
          --health-interval 10s
          --health-timeout 5s
          --health-retries 5
        ports:
          - 5432:5432
      
      redis:
        image: redis:7-alpine
        options: >-
          --health-cmd "redis-cli ping"
          --health-interval 10s
          --health-timeout 5s
          --health-retries 5
        ports:
          - 6379:6379
    
    steps:
      - uses: actions/checkout@v3
      
      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.24'
      
      - name: Run integration tests
        env:
          TEST_DB_HOST: localhost
          TEST_DB_NAME: reolink_server_test
          TEST_DB_USER: reolink
          TEST_DB_PASSWORD: password
          TEST_REDIS_HOST: localhost
        run: go test -v -tags=integration ./tests/integration/...
```

## Troubleshooting

### Database Connection Issues

```bash
# Check PostgreSQL is running
pg_isready -h localhost -p 5432

# Test connection
psql -h localhost -U reolink -d reolink_server_test

# Check TimescaleDB extension
psql -d reolink_server_test -c "SELECT * FROM pg_extension WHERE extname = 'timescaledb';"
```

### Redis Connection Issues

```bash
# Check Redis is running
redis-cli ping

# Test connection
redis-cli -h localhost -p 6379 ping
```

### Test Failures

1. **Check logs**: Integration tests output detailed logs
2. **Database state**: Ensure test database is clean before running tests
3. **Port conflicts**: Make sure required ports are available
4. **Environment variables**: Verify all required env vars are set

### Cleaning Test Database

```bash
# Drop and recreate test database
dropdb reolink_server_test
createdb reolink_server_test
psql -d reolink_server_test -c "CREATE EXTENSION IF NOT EXISTS timescaledb;"
```

## Future Enhancements

- [ ] Add more comprehensive test scenarios
- [ ] Implement WebSocket integration tests
- [ ] Add performance/load testing
- [ ] Implement test fixtures for common scenarios
- [ ] Add database migration tests
- [ ] Implement end-to-end UI tests with Selenium/Playwright

