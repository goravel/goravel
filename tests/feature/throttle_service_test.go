package feature

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"goravel/app/http/middleware"
	"goravel/tests"
	"goravel/tests/mocks"
)

type ThrottleServiceTestSuite struct {
	suite.Suite
	tests.TestCase
}

func TestThrottleServiceTestSuite(t *testing.T) {
	suite.Run(t, new(ThrottleServiceTestSuite))
}

func (s *ThrottleServiceTestSuite) TestThrottleService_Check_Success() {
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

func (s *ThrottleServiceTestSuite) TestThrottleService_Check_RateLimited() {
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

func (s *ThrottleServiceTestSuite) TestThrottleService_Check_NilStore() {
	// Given
	service := middleware.NewThrottleService()
	rule := mocks.NewMockRateLimitRule(nil)

	// When
	response, err := service.Check(nil, rule, "throttle:test:0:test-key")

	// Then
	assert.NoError(s.T(), err)
	assert.True(s.T(), response.Allowed)
}

func (s *ThrottleServiceTestSuite) TestThrottleService_BuildKey_WithCustomKey() {
	// Given
	service := middleware.NewThrottleService()
	rule := mocks.NewMockRateLimitRule(nil)
	rule.SetKey("custom-key")

	// When
	key := service.BuildKey(nil, rule, "test-limiter", 0)

	// Then
	assert.Equal(s.T(), "throttle:test-limiter:0:custom-key", key)
}

func (s *ThrottleServiceTestSuite) TestRateLimitRule_Methods() {
	// Given
	mockStore := &mocks.MockStore{}
	rule := middleware.NewRateLimitRule(mockStore)

	// Test SetKey
	result := rule.SetKey("test-key")
	assert.Equal(s.T(), rule, result)
	assert.Equal(s.T(), "test-key", rule.GetKey())

	// Test GetStore
	assert.Equal(s.T(), mockStore, rule.GetStore())

	// Test GetResponseCallback
	assert.NotNil(s.T(), rule.GetResponseCallback())
}