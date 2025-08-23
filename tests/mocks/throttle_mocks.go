package mocks

import (
	"context"
	"time"

	contractshttp "github.com/goravel/framework/contracts/http"

	"goravel/app/http/contracts"
)

// MockStore implements the Store interface for testing
type MockStore struct {
	TakeFunc func(ctx context.Context, key string) (tokens, remaining, reset uint64, ok bool, err error)
	GetFunc  func(ctx context.Context, key string) (tokens, remaining uint64, err error)
	SetFunc  func(ctx context.Context, key string, tokens uint64, interval time.Duration) error
	BurstFunc func(ctx context.Context, key string, tokens uint64) error
}

func (m *MockStore) Take(ctx context.Context, key string) (tokens, remaining, reset uint64, ok bool, err error) {
	if m.TakeFunc != nil {
		return m.TakeFunc(ctx, key)
	}
	return 100, 99, uint64(time.Now().Add(time.Minute).UnixNano()), true, nil
}

func (m *MockStore) Get(ctx context.Context, key string) (tokens, remaining uint64, err error) {
	if m.GetFunc != nil {
		return m.GetFunc(ctx, key)
	}
	return 100, 99, nil
}

func (m *MockStore) Set(ctx context.Context, key string, tokens uint64, interval time.Duration) error {
	if m.SetFunc != nil {
		return m.SetFunc(ctx, key, tokens, interval)
	}
	return nil
}

func (m *MockStore) Burst(ctx context.Context, key string, tokens uint64) error {
	if m.BurstFunc != nil {
		return m.BurstFunc(ctx, key, tokens)
	}
	return nil
}

// MockRateLimitRule implements the RateLimitRule interface for testing
type MockRateLimitRule struct {
	store            contractshttp.Store
	key              string
	responseCallback func(ctx contractshttp.Context)
}

func NewMockRateLimitRule(store contractshttp.Store) *MockRateLimitRule {
	return &MockRateLimitRule{
		store: store,
		responseCallback: func(ctx contractshttp.Context) {
			ctx.Request().Abort(contractshttp.StatusTooManyRequests)
		},
	}
}

func (m *MockRateLimitRule) GetStore() contractshttp.Store {
	return m.store
}

func (m *MockRateLimitRule) GetKey() string {
	return m.key
}

func (m *MockRateLimitRule) GetResponseCallback() func(ctx contractshttp.Context) {
	return m.responseCallback
}

func (m *MockRateLimitRule) SetKey(key string) contracts.RateLimitRule {
	m.key = key
	return m
}

func (m *MockRateLimitRule) SetResponseCallback(callback func(ctx contractshttp.Context)) contracts.RateLimitRule {
	m.responseCallback = callback
	return m
}

// MockThrottleService implements the ThrottleService interface for testing
type MockThrottleService struct {
	CheckFunc    func(ctx contractshttp.Context, rule contracts.RateLimitRule, key string) (contracts.ThrottleResponse, error)
	BuildKeyFunc func(ctx contractshttp.Context, rule contracts.RateLimitRule, name string, index int) string
}

func (m *MockThrottleService) Check(ctx contractshttp.Context, rule contracts.RateLimitRule, key string) (contracts.ThrottleResponse, error) {
	if m.CheckFunc != nil {
		return m.CheckFunc(ctx, rule, key)
	}
	return contracts.ThrottleResponse{
		Allowed:    true,
		Tokens:     100,
		Remaining:  99,
		Reset:      uint64(time.Now().Add(time.Minute).UnixNano()),
		RetryAfter: 0,
	}, nil
}

func (m *MockThrottleService) BuildKey(ctx contractshttp.Context, rule contracts.RateLimitRule, name string, index int) string {
	if m.BuildKeyFunc != nil {
		return m.BuildKeyFunc(ctx, rule, name, index)
	}
	return "test-key"
}