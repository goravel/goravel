# Optimized Throttle Middleware

This directory demonstrates the optimization of the Throttle method to use contracts instead of direct instance access, making it much easier to test. This addresses the issue mentioned in [#629](https://github.com/goravel/goravel/issues/629).

## Problem

The original throttle middleware in the framework had the following issues:

1. **Hard to test**: The middleware used type assertions to access concrete implementation details
2. **Tight coupling**: Direct access to `*httplimit.Limit` fields made the code brittle
3. **Poor abstraction**: No interfaces for key components, making mocking difficult

```go
// Original problematic code
if instance, exist := limit.(*httplimit.Limit); exist {
    tokens, remaining, reset, ok, err := instance.Store.Take(ctx, key(ctx, instance, name, index))
    // ... directly accessing instance.Store, instance.Key, instance.ResponseCallback
}
```

## Solution

The optimized approach introduces proper interfaces and abstractions:

### 1. Enhanced Contracts (`app/http/contracts/rate_limiter.go`)

- `RateLimitRule`: Interface for rate limiting rules with getter/setter methods
- `ThrottleResponse`: Structured response from throttle checks
- `ThrottleService`: Interface for throttle operations

### 2. Testable Implementations

- `RateLimitRule`: Concrete implementation of the rule interface
- `ThrottleService`: Service that handles throttle logic with proper abstraction
- `OptimizedThrottle`: Middleware that uses contracts instead of concrete types

### 3. Comprehensive Test Suite

The new approach makes testing straightforward:

```go
// Easy to mock and test
mockService := &mocks.MockThrottleService{
    CheckFunc: func(ctx contractshttp.Context, rule contracts.RateLimitRule, key string) (contracts.ThrottleResponse, error) {
        return contracts.ThrottleResponse{Allowed: false}, nil
    },
}

middlewareFunc := middleware.OptimizedThrottle("test", mockService)
```

## Key Improvements

### 1. Dependency Injection
The optimized middleware accepts a `ThrottleService` parameter, allowing easy testing with mocks.

### 2. Interface-Based Design
All components use interfaces, making the code more flexible and testable.

### 3. Separated Concerns
- `ThrottleService`: Handles rate limit checking logic
- `RateLimitRule`: Encapsulates rule configuration
- `OptimizedThrottle`: Coordinates the middleware flow

### 4. Better Error Handling
Structured error handling with proper return types.

## Usage Examples

### Basic Usage
```go
// Use with default service
middleware := middleware.OptimizedThrottle("api", nil)

// Use with custom service  
service := middleware.NewThrottleService()
middleware := middleware.OptimizedThrottle("api", service)
```

### Custom Response Handling
```go
customHandler := func(ctx contractshttp.Context) {
    ctx.Response().Json(429, map[string]string{
        "message": "Rate limit exceeded. Please try again later.",
    })
}

middleware := middleware.ThrottleWithCustomResponse("api", customHandler)
```

### Testing
```go
func TestThrottle(t *testing.T) {
    mockService := &mocks.MockThrottleService{
        CheckFunc: func(ctx, rule, key) (response, error) {
            return contracts.ThrottleResponse{Allowed: false}, nil
        },
    }
    
    middleware := middleware.OptimizedThrottle("test", mockService)
    // Test middleware behavior easily
}
```

## Files Structure

```
app/http/
├── contracts/
│   └── rate_limiter.go          # Interfaces for rate limiting
├── middleware/
│   ├── optimized_throttle.go    # Main optimized middleware
│   ├── rate_limit_rule.go       # Rule implementation
│   ├── throttle_service.go      # Service implementation
│   └── throttle_examples.go     # Usage examples
tests/
├── mocks/
│   └── throttle_mocks.go        # Mock implementations
└── feature/
    ├── optimized_throttle_test.go    # Service tests
    ├── throttle_middleware_test.go   # Middleware tests
    └── throttle_examples_test.go     # Example tests
```

## Benefits

1. **Easy Testing**: Mock any component independently
2. **Better Maintainability**: Clear separation of concerns
3. **Flexibility**: Easy to extend with custom behavior
4. **Type Safety**: Proper interfaces prevent runtime errors
5. **Documentation**: Clear contracts make the API self-documenting

This optimization demonstrates how to make middleware more testable while maintaining backward compatibility and improving code quality.