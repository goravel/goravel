# slog-multi: Advanced Handler Composition for Go's Structured Logging (pipelining, fanout, routing, failover...)

[![tag](https://img.shields.io/github/tag/samber/slog-multi.svg)](https://github.com/samber/slog-multi/releases)
![Go Version](https://img.shields.io/badge/Go-%3E%3D%201.21-%23007d9c)
[![GoDoc](https://godoc.org/github.com/samber/slog-multi?status.svg)](https://pkg.go.dev/github.com/samber/slog-multi)
![Build Status](https://github.com/samber/slog-multi/actions/workflows/test.yml/badge.svg)
[![Go report](https://goreportcard.com/badge/github.com/samber/slog-multi)](https://goreportcard.com/report/github.com/samber/slog-multi)
[![Coverage](https://img.shields.io/codecov/c/github/samber/slog-multi)](https://codecov.io/gh/samber/slog-multi)
[![Contributors](https://img.shields.io/github/contributors/samber/slog-multi)](https://github.com/samber/slog-multi/graphs/contributors)
[![License](https://img.shields.io/github/license/samber/slog-multi)](./LICENSE)

**slog-multi** provides advanced composition patterns for Go's structured logging (`slog`). It enables you to build sophisticated logging workflows by combining multiple handlers with different strategies for distribution, routing, transformation, and error handling.

## 🎯 Features

- **🔄 Fanout**: Distribute logs to multiple handlers in parallel
- **🛣️ Router**: Conditionally route logs based on custom criteria
- **🎯 First Match**: Route logs to the first matching handler only
- **🔄 Failover**: High-availability logging with automatic fallback
- **⚖️ Load Balancing**: Distribute load across multiple handlers
- **🔗 Pipeline**: Transform and filter logs with middleware chains
- **🛡️ Error Recovery**: Graceful handling of logging failures

Middlewares:
- **⚡ Inline Handlers**: Quick implementation of custom handlers
- **🔧 Inline Middleware**: Rapid development of transformation logic

<div align="center">
  <hr>
  <sup><b>Sponsored by:</b></sup>
  <br>
  <a href="https://www.dash0.com?utm_campaign=148395251-samber%20github%20sponsorship&utm_source=github&utm_medium=sponsorship&utm_content=samber">
    <div>
      <img src="https://github.com/user-attachments/assets/b1f2e876-0954-4dc3-824d-935d29ba8f3f" width="200" alt="Dash0">
    </div>
    <div>
      100% OpenTelemetry-native observability platform<br>Simple to use, built on open standards, and designed for full cost control
    </div>
  </a>
  <hr>
</div>

**See also:**

- [slog-multi](https://github.com/samber/slog-multi): `slog.Handler` chaining, fanout, routing, failover, load balancing...
- [slog-formatter](https://github.com/samber/slog-formatter): `slog` attribute formatting
- [slog-sampling](https://github.com/samber/slog-sampling): `slog` sampling policy
- [slog-mock](https://github.com/samber/slog-mock): `slog.Handler` for test purposes

**HTTP middlewares:**

- [slog-gin](https://github.com/samber/slog-gin): Gin middleware for `slog` logger
- [slog-echo](https://github.com/samber/slog-echo): Echo middleware for `slog` logger
- [slog-fiber](https://github.com/samber/slog-fiber): Fiber middleware for `slog` logger
- [slog-chi](https://github.com/samber/slog-chi): Chi middleware for `slog` logger
- [slog-http](https://github.com/samber/slog-http): `net/http` middleware for `slog` logger

**Loggers:**

- [slog-zap](https://github.com/samber/slog-zap): A `slog` handler for `Zap`
- [slog-zerolog](https://github.com/samber/slog-zerolog): A `slog` handler for `Zerolog`
- [slog-logrus](https://github.com/samber/slog-logrus): A `slog` handler for `Logrus`

**Log sinks:**

- [slog-datadog](https://github.com/samber/slog-datadog): A `slog` handler for `Datadog`
- [slog-betterstack](https://github.com/samber/slog-betterstack): A `slog` handler for `Betterstack`
- [slog-rollbar](https://github.com/samber/slog-rollbar): A `slog` handler for `Rollbar`
- [slog-loki](https://github.com/samber/slog-loki): A `slog` handler for `Loki`
- [slog-sentry](https://github.com/samber/slog-sentry): A `slog` handler for `Sentry`
- [slog-syslog](https://github.com/samber/slog-syslog): A `slog` handler for `Syslog`
- [slog-logstash](https://github.com/samber/slog-logstash): A `slog` handler for `Logstash`
- [slog-fluentd](https://github.com/samber/slog-fluentd): A `slog` handler for `Fluentd`
- [slog-graylog](https://github.com/samber/slog-graylog): A `slog` handler for `Graylog`
- [slog-quickwit](https://github.com/samber/slog-quickwit): A `slog` handler for `Quickwit`
- [slog-slack](https://github.com/samber/slog-slack): A `slog` handler for `Slack`
- [slog-telegram](https://github.com/samber/slog-telegram): A `slog` handler for `Telegram`
- [slog-mattermost](https://github.com/samber/slog-mattermost): A `slog` handler for `Mattermost`
- [slog-microsoft-teams](https://github.com/samber/slog-microsoft-teams): A `slog` handler for `Microsoft Teams`
- [slog-webhook](https://github.com/samber/slog-webhook): A `slog` handler for `Webhook`
- [slog-kafka](https://github.com/samber/slog-kafka): A `slog` handler for `Kafka`
- [slog-nats](https://github.com/samber/slog-nats): A `slog` handler for `NATS`
- [slog-parquet](https://github.com/samber/slog-parquet): A `slog` handler for `Parquet` + `Object Storage`
- [slog-channel](https://github.com/samber/slog-channel): A `slog` handler for Go channels

## 🚀 Installation

```sh
go get github.com/samber/slog-multi
```

**Compatibility**: go >= 1.21

No breaking changes will be made to exported APIs before v2.0.0.

> [!WARNING]
> Use this library carefully, log processing can be very costly (!)
> 
> Excessive logging —with multiple processing steps and destinations— can introduce significant overhead, which is generally undesirable in performance-critical paths. Logging is always expensive, and sometimes, metrics or a sampling strategy are cheaper. The library itself does not generate extra load.

## 💡 Usage

GoDoc: [https://pkg.go.dev/github.com/samber/slog-multi](https://pkg.go.dev/github.com/samber/slog-multi)

### Broadcast: `slogmulti.Fanout()`

Distribute logs to multiple `slog.Handler` in parallel for maximum throughput and redundancy.

```go
import (
    "net"
    slogmulti "github.com/samber/slog-multi"
    "log/slog"
    "os"
    "time"
)

func main() {
    logstash, _ := net.Dial("tcp", "logstash.acme:4242")    // use github.com/netbrain/goautosocket for auto-reconnect
    datadogHandler := slogdatadog.NewDatadogHandler(slogdatadog.Option{
        APIKey: "your-api-key",
        Service: "my-service",
    })
    stderr := os.Stderr

    logger := slog.New(
        slogmulti.Fanout(
            slog.NewJSONHandler(logstash, &slog.HandlerOptions{}),  // pass to first handler: logstash over tcp
            slog.NewTextHandler(stderr, &slog.HandlerOptions{}),    // then to second handler: stderr
            datadogHandler,
            // ...
        ),
    )

    logger.
        With(
            slog.Group("user",
                slog.String("id", "user-123"),
                slog.Time("created_at", time.Now()),
            ),
        ).
        With("environment", "dev").
        With("error", fmt.Errorf("an error")).
        Error("A message")
}
```

Stderr output:

```
time=2023-04-10T14:00:0.000000+00:00 level=ERROR msg="A message" user.id=user-123 user.created_at=2023-04-10T14:00:0.000000+00:00 environment=dev error="an error"
```

Netcat output:

```json
{
	"time":"2023-04-10T14:00:0.000000+00:00",
	"level":"ERROR",
	"msg":"A message",
	"user":{
		"id":"user-123",
		"created_at":"2023-04-10T14:00:0.000000+00:00"
	},
	"environment":"dev",
	"error":"an error"
}
```

### Routing: `slogmulti.Router()`

Distribute logs to all matching `slog.Handler` based on custom criteria like log level, attributes, or business logic.

```go
import (
    "context"
    slogmulti "github.com/samber/slog-multi"
    slogslack "github.com/samber/slog-slack"
    "log/slog"
    "os"
)

func main() {
    slackChannelUS := slogslack.Option{Level: slog.LevelError, WebhookURL: "xxx", Channel: "supervision-us"}.NewSlackHandler()
    slackChannelEU := slogslack.Option{Level: slog.LevelError, WebhookURL: "xxx", Channel: "supervision-eu"}.NewSlackHandler()
    slackChannelAPAC := slogslack.Option{Level: slog.LevelError, WebhookURL: "xxx", Channel: "supervision-apac"}.NewSlackHandler()

    consoleHandler := slog.NewTextHandler(os.Stderr, nil)

    logger := slog.New(
        slogmulti.Router().
            Add(slackChannelUS, recordMatchRegion("us")).
            Add(slackChannelEU, recordMatchRegion("eu")).
            Add(slackChannelAPAC, recordMatchRegion("apac")).
            Add(consoleHandler, slogmulti.LevelIs(slog.LevelInfo, slog.LevelDebug)).
            Handler(),
    )

    logger.
        With("region", "us").
        With("pool", "us-east-1").
        Error("Server desynchronized")
}

func recordMatchRegion(region string) func(ctx context.Context, r slog.Record) bool {
    return func(ctx context.Context, r slog.Record) bool {
        ok := false

        r.Attrs(func(attr slog.Attr) bool {
            if attr.Key == "region" && attr.Value.Kind() == slog.KindString && attr.Value.String() == region {
                ok = true
                return false
            }

            return true
        })

        return ok
    }
}
```

**Use Cases:**
- Environment-specific logging (dev vs prod)
- Level-based routing (errors to Slack, info to console)
- Business logic routing (user actions vs system events)

### First Match Routing: `Router().FirstMatch()`

Route logs to the **first matching handler only**, unlike regular routing which sends to all matching handlers. Perfect for priority-based routing where you want exactly one handler to receive each log.

```go
import (
    slogmulti "github.com/samber/slog-multi"
    slogslack "github.com/samber/slog-slack"
    "log/slog"
)

func main() {
    queryChannel := slogslack.Option{Level: slog.LevelDebug, WebhookURL: "xxx", Channel: "db-queries"}.NewSlackHandler()
    requestChannel := slogslack.Option{Level: slog.LevelError, WebhookURL: "xxx", Channel: "service-requests"}.NewSlackHandler()
    influxdbChannel := slogslack.Option{Level: slog.LevelInfo, WebhookURL: "xxx", Channel: "influxdb-metrics"}.NewSlackHandler()
    fallbackChannel := slogslack.Option{Level: slog.LevelError, WebhookURL: "xxx", Channel: "logs"}.NewSlackHandler()

    logger := slog.New(
        slogmulti.Router().
            Add(queryChannel, slogmulti.AttrKindIs("query", slog.KindString, "args", slog.KindAny)).
            Add(requestChannel, slogmulti.AttrKindIs("method", slog.KindString, "body", slog.KindAny)).
            Add(influxdbChannel, slogmulti.AttrValueIs("scope", "influx")).
            Add(fallbackChannel).  // Catch-all for everything else
            FirstMatch().           // ← Enable first-match routing
            Handler(),
    )

    // Goes to queryChannel only (stops at first match)
    logger.Debug("Executing SQL query", "query", "SELECT * FROM users WHERE id = ?", "args", []int{1})

    // Goes to requestChannel only (stops at first match)
    logger.Error("Incoming request failed", "method", "POST", "body", "{'name':'test'}")

    // Goes to fallbackChannel (no other handlers matched)
    logger.Error("An unexpected error occurred")
}
```

#### Built-in Predicates

**Level predicates:**
- `LevelIs(levels ...slog.Level)` - Match specific log levels
- `LevelIsNot(levels ...slog.Level)` - Exclude specific log levels

**Message predicates:**
- `MessageIs(msg string)` - Exact message match
- `MessageIsNot(msg string)` - Message doesn't match
- `MessageContains(part string)` - Message contains substring
- `MessageNotContains(part string)` - Message doesn't contain substring

**Attribute predicates:**
- `AttrValueIs(key, value, ...)` - Check attributes have exact values
- `AttrKindIs(key, kind, ...)` - Check attributes have specific types

### Failover: `slogmulti.Failover()`

Ensure logging reliability by trying multiple handlers in order until one succeeds. Perfect for high-availability scenarios.

```go
import (
    "net"
    slogmulti "github.com/samber/slog-multi"
    "log/slog"
    "os"
    "time"
)


func main() {
    // Create connections to multiple log servers
    // ncat -l 1000 -k
    // ncat -l 1001 -k
    // ncat -l 1002 -k

    // List AZs - use github.com/netbrain/goautosocket for auto-reconnect
    logstash1, _ := net.Dial("tcp", "logstash.eu-west-3a.internal:1000")
    logstash2, _ := net.Dial("tcp", "logstash.eu-west-3b.internal:1000")
    logstash3, _ := net.Dial("tcp", "logstash.eu-west-3c.internal:1000")

    logger := slog.New(
        slogmulti.Failover()(
            slog.HandlerOptions{}.NewJSONHandler(logstash1, nil),    // Primary
            slog.HandlerOptions{}.NewJSONHandler(logstash2, nil),    // Secondary
            slog.HandlerOptions{}.NewJSONHandler(logstash3, nil),    // Tertiary
        ),
    )

    logger.
        With(
            slog.Group("user",
                slog.String("id", "user-123"),
                slog.Time("created_at", time.Now()),
            ),
        ).
        With("environment", "dev").
        With("error", fmt.Errorf("an error")).
        Error("A message")
}
```

**Use Cases:**
- High-availability logging infrastructure
- Disaster recovery scenarios
- Multi-region deployments

### Load balancing: `slogmulti.Pool()`

Distribute logging load across multiple handlers using round-robin with randomization to increase throughput and provide redundancy.

```go
import (
    "net"
    slogmulti "github.com/samber/slog-multi"
    "log/slog"
    "os"
    "time"
)

func main() {
    // Create multiple log servers
    // ncat -l 1000 -k
    // ncat -l 1001 -k
    // ncat -l 1002 -k

    // List AZs - use github.com/netbrain/goautosocket for auto-reconnect
    logstash1, _ := net.Dial("tcp", "logstash.eu-west-3a.internal:1000")
    logstash2, _ := net.Dial("tcp", "logstash.eu-west-3b.internal:1000")
    logstash3, _ := net.Dial("tcp", "logstash.eu-west-3c.internal:1000")

    logger := slog.New(
        slogmulti.Pool()(
            // A random handler will be picked for each log
            slog.HandlerOptions{}.NewJSONHandler(logstash1, nil),
            slog.HandlerOptions{}.NewJSONHandler(logstash2, nil),
            slog.HandlerOptions{}.NewJSONHandler(logstash3, nil),
        ),
    )

    // High-volume logging
    for i := 0; i < 1000; i++ {
        logger.
            With(
                slog.Group("user",
                    slog.String("id", "user-123"),
                    slog.Time("created_at", time.Now()),
                ),
            ).
            With("environment", "dev").
            With("error", fmt.Errorf("an error")).
            Error("A message")
    }
}
```

**Use Cases:**
- High-throughput logging scenarios
- Distributed logging infrastructure
- Performance optimization

### Recover errors: `slogmulti.RecoverHandlerError()`

Gracefully handle logging failures without crashing the application. Catches both panics and errors from handlers.

```go
import (
    "context"
    slogformatter "github.com/samber/slog-formatter"
    slogmulti "github.com/samber/slog-multi"
    "log/slog"
    "os"
)

recovery := slogmulti.RecoverHandlerError(
    func(ctx context.Context, record slog.Record, err error) {
        // will be called only if subsequent handlers fail or return an error
        log.Println(err.Error())
    },
)
sink := NewSinkHandler(...)

logger := slog.New(
    slogmulti.
        Pipe(recovery).
        Handler(sink),
)

err := fmt.Errorf("an error")
logger.Error("a message",
    slog.Any("very_private_data", "abcd"),
    slog.Any("user", user),
    slog.Any("err", err))

// outputs:
// time=2023-04-10T14:00:0.000000+00:00 level=ERROR msg="a message" error.message="an error" error.type="*errors.errorString" user="John doe" very_private_data="********"
```

### Pipelining: `slogmulti.Pipe()`

Transform and filter logs using middleware chains. Perfect for data privacy, formatting, and cross-cutting concerns.

```go
import (
    "context"
    slogmulti "github.com/samber/slog-multi"
    "log/slog"
    "os"
    "time"
)

func main() {
    // First middleware: format Go `error` type into an structured object {error: "*myCustomErrorType", message: "could not reach https://a.b/c"}
    errorFormattingMiddleware := slogmulti.NewHandleInlineMiddleware(func(ctx context.Context, record slog.Record, next func(context.Context, slog.Record) error) error {
        record.Attrs(func(attr slog.Attr) bool {
            if attr.Key == "error" && attr.Value.Kind() == slog.KindAny {
                if err, ok := attr.Value.Any().(error); ok {
                    record.AddAttrs(
                        slog.String("error_type", "error"),
                        slog.String("error_message", err.Error()),
                    )
                }
            }
            return true
        })
        return next(ctx, record)
    })

    // Second middleware: remove PII
    gdprMiddleware := slogmulti.NewHandleInlineMiddleware(func(ctx context.Context, record slog.Record, next func(context.Context, slog.Record) error) error {
        record.Attrs(func(attr slog.Attr) bool {
            if attr.Key == "email" || attr.Key == "phone" || attr.Key == "created_at" {
                record.AddAttrs(slog.String(attr.Key, "*********"))
            }
            return true
        })
        return next(ctx, record)
    })

    // Final handler
    sink := slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{})

    logger := slog.New(
        slogmulti.
            Pipe(errorFormattingMiddleware).
            Pipe(gdprMiddleware).
            // ...
            Handler(sink),
    )

    logger.
        With(
            slog.Group("user",
                slog.String("id", "user-123"),
                slog.String("email", "user-123"),
                slog.Time("created_at", time.Now()),
            ),
        ).
        With("environment", "dev").
        Error("A message",
            slog.String("foo", "bar"),
            slog.Any("error", fmt.Errorf("an error")),
        )
}
```

Stderr output:

```json
{
    "time":"2023-04-10T14:00:0.000000+00:00",
    "level":"ERROR",
    "msg":"A message",
    "user":{
        "email":"*******",
        "phone":"*******",
        "created_at":"*******"
    },
    "environment":"dev",
    "foo":"bar",
    "error":{
        "type":"*myCustomErrorType",
        "message":"an error"
    }
}
```

**Use Cases:**
- Data privacy and GDPR compliance
- Error formatting and standardization
- Log enrichment and transformation
- Performance monitoring and metrics

## 🔧 Advanced Patterns

### Custom middleware

Middleware must match the following prototype:

```go
type Middleware func(slog.Handler) slog.Handler
```

The example above uses:
- a custom middleware, [see here](./examples/pipe/gdpr.go)
- an inline middleware, [see here](./examples/pipe/errors.go)

> **Note**: `WithAttrs` and `WithGroup` methods of custom middleware must return a new instance, not `this`.

#### Inline handler

Inline handlers provide shortcuts to implement `slog.Handler` without creating full struct implementations.

```go
mdw := slogmulti.NewHandleInlineHandler(
    // simulate "Handle()" method
    func(ctx context.Context, groups []string, attrs []slog.Attr, record slog.Record) error {
        // Custom logic here
        // [...]
        return nil
    },
)
```

```go
mdw := slogmulti.NewInlineHandler(
    // simulate "Enabled()" method
    func(ctx context.Context, groups []string, attrs []slog.Attr, level slog.Level) bool {
        // Custom logic here
        // [...]
        return true
    },
    // simulate "Handle()" method
    func(ctx context.Context, groups []string, attrs []slog.Attr, record slog.Record) error {
        // Custom logic here
        // [...]
        return nil
    },
)
```

#### Inline middleware

Inline middleware provides shortcuts to implement middleware functions that hook specific methods.

#### Hook `Enabled()` Method

```go
middleware := slogmulti.NewEnabledInlineMiddleware(func(ctx context.Context, level slog.Level, next func(context.Context, slog.Level) bool) bool{
    // Custom logic before calling next
    if level == slog.LevelDebug {
        return false // Skip debug logs
    }
    return next(ctx, level)
})
```

#### Hook `Handle()` Method

```go
middleware := slogmulti.NewHandleInlineMiddleware(func(ctx context.Context, record slog.Record, next func(context.Context, slog.Record) error) error {
    // Add timestamp to all logs
    record.AddAttrs(slog.Time("logged_at", time.Now()))
    return next(ctx, record)
})
```

#### Hook `WithAttrs()` Method

```go
mdw := slogmulti.NewWithAttrsInlineMiddleware(func(attrs []slog.Attr, next func([]slog.Attr) slog.Handler) slog.Handler{
    // Filter out sensitive attributes
    filtered := make([]slog.Attr, 0, len(attrs))
    for _, attr := range attrs {
        if attr.Key != "password" && attr.Key != "token" {
            filtered = append(filtered, attr)
        }
    }
    return next(attrs)
})
```

#### Hook `WithGroup()` Method

```go
mdw := slogmulti.NewWithGroupInlineMiddleware(func(name string, next func(string) slog.Handler) slog.Handler{
    // Add prefix to group names
    prefixedName := "app." + name
    return next(name)
})
```

#### Complete Inline Middleware

> **Warning**: You should implement your own middleware for complex scenarios.

```go
mdw := slogmulti.NewInlineMiddleware(
    func(ctx context.Context, level slog.Level, next func(context.Context, slog.Level) bool) bool{
        // Custom logic here
        // [...]
        return next(ctx, level)
    },
    func(ctx context.Context, record slog.Record, next func(context.Context, slog.Record) error) error{
        // Custom logic here
        // [...]
        return next(ctx, record)
    },
    func(attrs []slog.Attr, next func([]slog.Attr) slog.Handler) slog.Handler{
        // Custom logic here
        // [...]
        return next(attrs)
    },
    func(name string, next func(string) slog.Handler) slog.Handler{
        // Custom logic here
        // [...]
        return next(name)
    },
)
```

## 💡 Best Practices

### Performance Considerations

- **Use Fanout sparingly**: Broadcasting to many handlers can impact performance
- **Implement sampling**: For high-volume logs, consider sampling strategies
- **Monitor handler performance**: Some handlers (like network-based ones) can be slow
- **Use buffering**: Consider buffering for network-based handlers

### Error Handling

- **Always use error recovery**: Wrap handlers with `RecoverHandlerError`
- **Implement fallbacks**: Use failover patterns for critical logging
- **Monitor logging failures**: Track when logging fails to identify issues

### Security and Privacy

- **Redact sensitive data**: Use middleware to remove PII and secrets
- **Validate log content**: Ensure logs don't contain sensitive information
- **Use secure connections**: For network-based handlers, use TLS

### Monitoring and Observability

- **Add correlation IDs**: Include request IDs in logs for tracing
- **Structured logging**: Use slog's structured logging features consistently
- **Log levels**: Use appropriate log levels for different types of information

## 🤝 Contributing

- Ping me on twitter [@samuelberthe](https://twitter.com/samuelberthe) (DMs, mentions, whatever :))
- Fork the [project](https://github.com/samber/slog-multi)
- Fix [open issues](https://github.com/samber/slog-multi/issues) or request new features

Don't hesitate ;)

```bash
# Install some dev dependencies
make tools

# Run tests
make test
# or
make watch-test
```

## 👤 Contributors

![Contributors](https://contrib.rocks/image?repo=samber/slog-multi)

## 💫 Show your support

If this project helped you, please give it a ⭐️ on GitHub!

[![GitHub Sponsors](https://img.shields.io/github/sponsors/samber?style=for-the-badge)](https://github.com/sponsors/samber)

## 📝 License

Copyright © 2023 [Samuel Berthe](https://github.com/samber).

This project is [MIT](./LICENSE) licensed.
