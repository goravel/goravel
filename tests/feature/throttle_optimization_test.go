package feature

import (
	"context"
	"errors"
	"testing"
	"time"

	contractshttp "github.com/goravel/framework/contracts/http"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"goravel/app/http/contracts"
	"goravel/app/http/middleware"
	"goravel/tests"
	"goravel/tests/mocks"
)

type ThrottleOptimizationTestSuite struct {
	suite.Suite
	tests.TestCase
}

func TestThrottleOptimizationTestSuite(t *testing.T) {
	suite.Run(t, new(ThrottleOptimizationTestSuite))
}

func (s *ThrottleOptimizationTestSuite) TestThrottleService_Check_Success() {
	// Given
	service := middleware.NewThrottleService()
	mockStore := &mocks.MockStore{
		TakeFunc: func(ctx context.Context, key string) (tokens, remaining, reset uint64, ok bool, err error) {
			return 100, 99, uint64(time.Now().Add(time.Minute).UnixNano()), true, nil
		},
	}
	rule := mocks.NewMockRateLimitRule(mockStore)
	rule.SetKey("test-key")

	// When
	response, err := service.Check(nil, rule, "throttle:test:0:test-key")

	// Then
	assert.NoError(s.T(), err)
	assert.True(s.T(), response.Allowed)
	assert.Equal(s.T(), uint64(100), response.Tokens)
	assert.Equal(s.T(), uint64(99), response.Remaining)
	assert.Greater(s.T(), response.Reset, uint64(0))
}

func (s *ThrottleOptimizationTestSuite) TestThrottleService_Check_RateLimited() {
	// Given
	service := middleware.NewThrottleService()
	resetTime := uint64(time.Now().Add(time.Minute).UnixNano())
	mockStore := &mocks.MockStore{
		TakeFunc: func(ctx context.Context, key string) (tokens, remaining, reset uint64, ok bool, err error) {
			return 100, 0, resetTime, false, nil
		},
	}
	rule := mocks.NewMockRateLimitRule(mockStore)
	rule.SetKey("test-key")

	// When
	response, err := service.Check(nil, rule, "throttle:test:0:test-key")

	// Then
	assert.NoError(s.T(), err)
	assert.False(s.T(), response.Allowed)
	assert.Equal(s.T(), uint64(100), response.Tokens)
	assert.Equal(s.T(), uint64(0), response.Remaining)
	assert.Equal(s.T(), resetTime, response.Reset)
	assert.Greater(s.T(), response.RetryAfter, 0)
}

func (s *ThrottleOptimizationTestSuite) TestThrottleService_Check_StoreError() {
	// Given
	service := middleware.NewThrottleService()
	expectedError := errors.New("store error")
	mockStore := &mocks.MockStore{
		TakeFunc: func(ctx context.Context, key string) (tokens, remaining, reset uint64, ok bool, err error) {
			return 0, 0, 0, false, expectedError
		},
	}
	rule := mocks.NewMockRateLimitRule(mockStore)

	// When
	response, err := service.Check(nil, rule, "throttle:test:0:test-key")

	// Then
	assert.Error(s.T(), err)
	assert.Equal(s.T(), expectedError, err)
	assert.False(s.T(), response.Allowed)
}

func (s *ThrottleOptimizationTestSuite) TestThrottleService_Check_NilStore() {
	// Given
	service := middleware.NewThrottleService()
	rule := mocks.NewMockRateLimitRule(nil)

	// When
	response, err := service.Check(nil, rule, "throttle:test:0:test-key")

	// Then
	assert.NoError(s.T(), err)
	assert.True(s.T(), response.Allowed)
}

func (s *ThrottleOptimizationTestSuite) TestThrottleService_BuildKey_WithCustomKey() {
	// Given
	service := middleware.NewThrottleService()
	rule := mocks.NewMockRateLimitRule(nil)
	rule.SetKey("custom-key")

	// When
	key := service.BuildKey(nil, rule, "test-limiter", 0)

	// Then
	assert.Equal(s.T(), "throttle:test-limiter:0:custom-key", key)
}

func (s *ThrottleOptimizationTestSuite) TestThrottleService_BuildKey_WithEmptyKey() {
	// Given
	service := middleware.NewThrottleService()
	rule := mocks.NewMockRateLimitRule(nil)
	rule.SetKey("")

	// When
	key := service.BuildKey(nil, rule, "test-limiter", 0)

	// Then  
	// Should use empty key when no context or context has no request
	assert.Equal(s.T(), "throttle:test-limiter:0:", key)
}

func (s *ThrottleOptimizationTestSuite) TestRateLimitRule_GettersAndSetters() {
	// Given
	mockStore := &mocks.MockStore{}
	rule := middleware.NewRateLimitRule(mockStore)

	// Test initial state
	assert.Equal(s.T(), mockStore, rule.GetStore())
	assert.Equal(s.T(), "", rule.GetKey())
	assert.NotNil(s.T(), rule.GetResponseCallback())

	// Test SetKey
	result := rule.SetKey("test-key")
	assert.Equal(s.T(), rule, result)
	assert.Equal(s.T(), "test-key", rule.GetKey())

	// Test SetResponseCallback
	called := false
	callback := func(ctx contractshttp.Context) {
		called = true
	}
	result = rule.SetResponseCallback(callback)
	assert.Equal(s.T(), rule, result)
	
	// Test callback execution
	rule.GetResponseCallback()(nil)
	assert.True(s.T(), called)
}

func (s *ThrottleOptimizationTestSuite) TestMockThrottleService() {
	// Given
	mockService := &mocks.MockThrottleService{
		CheckFunc: func(ctx contractshttp.Context, rule contracts.RateLimitRule, key string) (contracts.ThrottleResponse, error) {
			return contracts.ThrottleResponse{
				Allowed:   false,
				Tokens:    100,
				Remaining: 0,
				Reset:     uint64(time.Now().Add(time.Minute).UnixNano()),
			}, nil
		},
		BuildKeyFunc: func(ctx contractshttp.Context, rule contracts.RateLimitRule, name string, index int) string {
			return "mocked-key"
		},
	}

	rule := mocks.NewMockRateLimitRule(&mocks.MockStore{})

	// When
	response, err := mockService.Check(nil, rule, "test-key")
	key := mockService.BuildKey(nil, rule, "test", 0)

	// Then
	assert.NoError(s.T(), err)
	assert.False(s.T(), response.Allowed)
	assert.Equal(s.T(), uint64(100), response.Tokens)
	assert.Equal(s.T(), uint64(0), response.Remaining)
	assert.Equal(s.T(), "mocked-key", key)
}

func (s *ThrottleOptimizationTestSuite) TestTestableThrottleRule() {
	// Given
	mockStore := &mocks.MockStore{}
	responseCallbackCalled := false
	
	responseCallback := func(ctx contractshttp.Context) {
		responseCallbackCalled = true
	}

	// When
	rule := middleware.TestableThrottleRule(mockStore, "test-key", responseCallback)

	// Then
	assert.Equal(s.T(), mockStore, rule.GetStore())
	assert.Equal(s.T(), "test-key", rule.GetKey())
	
	// Test response callback
	rule.GetResponseCallback()(nil)
	assert.True(s.T(), responseCallbackCalled)
}

func (s *ThrottleOptimizationTestSuite) TestTestableThrottleRule_DefaultValues() {
	// Given
	mockStore := &mocks.MockStore{}

	// When - passing empty values
	rule := middleware.TestableThrottleRule(mockStore, "", nil)

	// Then
	assert.Equal(s.T(), mockStore, rule.GetStore())
	assert.Equal(s.T(), "", rule.GetKey())
	assert.NotNil(s.T(), rule.GetResponseCallback()) // Should have default callback
}

func (s *ThrottleOptimizationTestSuite) TestPerMinuteTestable() {
	// Given
	mockStore := &mocks.MockStore{}

	// When
	rule := middleware.PerMinuteTestable(100, mockStore)

	// Then
	assert.Equal(s.T(), mockStore, rule.GetStore())
	assert.NotNil(s.T(), rule.GetResponseCallback())
}

// This test demonstrates how easy it is to test throttle behavior with the optimized approach
func (s *ThrottleOptimizationTestSuite) TestThrottleScenarios() {
	s.Run("scenario: allow request under limit", func() {
		service := middleware.NewThrottleService()
		store := &mocks.MockStore{
			TakeFunc: func(ctx context.Context, key string) (tokens, remaining, reset uint64, ok bool, err error) {
				return 60, 59, uint64(time.Now().Add(time.Minute).UnixNano()), true, nil
			},
		}
		rule := middleware.NewRateLimitRule(store).SetKey("user:123")
		
		response, err := service.Check(nil, rule, "throttle:api:0:user:123")
		
		assert.NoError(s.T(), err)
		assert.True(s.T(), response.Allowed)
		assert.Equal(s.T(), uint64(60), response.Tokens)
		assert.Equal(s.T(), uint64(59), response.Remaining)
	})

	s.Run("scenario: block request at limit", func() {
		service := middleware.NewThrottleService()
		store := &mocks.MockStore{
			TakeFunc: func(ctx context.Context, key string) (tokens, remaining, reset uint64, ok bool, err error) {
				return 60, 0, uint64(time.Now().Add(time.Minute).UnixNano()), false, nil
			},
		}
		rule := middleware.NewRateLimitRule(store).SetKey("user:456")
		
		response, err := service.Check(nil, rule, "throttle:api:0:user:456")
		
		assert.NoError(s.T(), err)
		assert.False(s.T(), response.Allowed)
		assert.Equal(s.T(), uint64(60), response.Tokens)
		assert.Equal(s.T(), uint64(0), response.Remaining)
		assert.Greater(s.T(), response.RetryAfter, 0)
	})

	s.Run("scenario: database error", func() {
		service := middleware.NewThrottleService()
		store := &mocks.MockStore{
			TakeFunc: func(ctx context.Context, key string) (tokens, remaining, reset uint64, ok bool, err error) {
				return 0, 0, 0, false, errors.New("database connection failed")
			},
		}
		rule := middleware.NewRateLimitRule(store).SetKey("user:789")
		
		response, err := service.Check(nil, rule, "throttle:api:0:user:789")
		
		assert.Error(s.T(), err)
		assert.Contains(s.T(), err.Error(), "database connection failed")
		assert.False(s.T(), response.Allowed)
	})
}