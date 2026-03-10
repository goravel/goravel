# eris ![Logo][eris-logo]

[![GoDoc][doc-img]][doc] [![Build][ci-img]][ci] [![GoReport][report-img]][report] [![Coverage Status][cov-img]][cov]

Package `eris` is an error handling library with readable stack traces and JSON formatting support.

`go get github.com/rotisserie/eris`

<!-- toc -->

- [Why you should switch to eris](#why-you-should-switch-to-eris)
- [Using eris](#using-eris)
  * [Creating errors](#creating-errors)
  * [Wrapping errors](#wrapping-errors)
  * [Formatting and logging errors](#formatting-and-logging-errors)
  * [Interpreting eris stack traces](#interpreting-eris-stack-traces)
  * [Inverting the stack trace and error output](#inverting-the-stack-trace-and-error-output)
  * [Inspecting errors](#inspecting-errors)
  * [Formatting with custom separators](#formatting-with-custom-separators)
  * [Writing a custom output format](#writing-a-custom-output-format)
  * [Sending error traces to Sentry](#sending-error-traces-to-sentry)
- [Comparison to other packages (e.g. pkg/errors)](#comparison-to-other-packages-eg-pkgerrors)
  * [Error formatting and stack traces](#error-formatting-and-stack-traces)
- [Migrating to eris](#migrating-to-eris)
- [Contributing](#contributing)

<!-- tocstop -->

## Why you should switch to eris

This package was inspired by a simple question: what if you could fix a bug without wasting time replicating the issue or digging through the code? With that in mind, this package is designed to give you more control over error handling via error wrapping, stack tracing, and output formatting.

The [example](https://github.com/rotisserie/eris/blob/master/examples/logging/example.go) that generated the output below simulates a realistic error handling scenario and demonstrates how to wrap and log errors with minimal effort. This specific error occurred because a user tried to access a file that can't be located, and the output shows a clear path from the top of the call stack to the source.

```json
{
  "error":{
    "root":{
      "message":"error internal server",
      "stack":[
        "main.main:/Users/roti/go/src/github.com/rotisserie/eris/examples/logging/example.go:143",
        "main.ProcessResource:/Users/roti/go/src/github.com/rotisserie/eris/examples/logging/example.go:85",
        "main.ProcessResource:/Users/roti/go/src/github.com/rotisserie/eris/examples/logging/example.go:82",
        "main.GetRelPath:/Users/roti/go/src/github.com/rotisserie/eris/examples/logging/example.go:61"
      ]
    },
    "wrap":[
      {
        "message":"failed to get relative path for resource 'res2'",
        "stack":"main.ProcessResource:/Users/roti/go/src/github.com/rotisserie/eris/examples/logging/example.go:85"
      },
      {
        "message":"Rel: can't make ./some/malformed/absolute/path/data.json relative to /Users/roti/",
        "stack":"main.GetRelPath:/Users/roti/go/src/github.com/rotisserie/eris/examples/logging/example.go:61"
      }
    ]
  },
  "level":"error",
  "method":"ProcessResource",
  "msg":"method completed with error",
  "time":"2020-01-16T11:20:01-05:00"
}
```

Many of the methods in this package will look familiar if you've used [pkg/errors](https://github.com/pkg/errors) or [xerrors](https://github.com/golang/xerrors), but `eris` employs some additional tricks during error wrapping and unwrapping that greatly improve the readability of the stack trace. This package also takes a unique approach to formatting errors that allows you to write custom formats that conform to your error or log aggregator of choice. You can find more information on the differences between `eris` and `pkg/errors` [here](#comparison-to-other-packages-eg-pkgerrors).

## Using eris

### Creating errors

Creating errors is simple via [`eris.New`](https://pkg.go.dev/github.com/rotisserie/eris#New).

```golang
var (
  // global error values can be useful when wrapping errors or inspecting error types
  ErrInternalServer = eris.New("error internal server")
)

func (req *Request) Validate() error {
  if req.ID == "" {
    // or return a new error at the source if you prefer
    return eris.New("error bad request")
  }
  return nil
}
```

### Wrapping errors

[`eris.Wrap`](https://pkg.go.dev/github.com/rotisserie/eris#Wrap) adds context to an error while preserving the original error.

```golang
relPath, err := GetRelPath("/Users/roti/", resource.AbsPath)
if err != nil {
  // wrap the error if you want to add more context
  return nil, eris.Wrapf(err, "failed to get relative path for resource '%v'", resource.ID)
}
```

### Formatting and logging errors

[`eris.ToString`](https://pkg.go.dev/github.com/rotisserie/eris#ToString) and [`eris.ToJSON`](https://pkg.go.dev/github.com/rotisserie/eris#ToJSON) should be used to log errors with the default format (shown above). The JSON method returns a `map[string]interface{}` type for compatibility with Go's `encoding/json` package and many common JSON loggers (e.g. [logrus](https://github.com/sirupsen/logrus)).

```golang
// format the error to JSON with the default format and stack traces enabled
formattedJSON := eris.ToJSON(err, true)
fmt.Println(json.Marshal(formattedJSON)) // marshal to JSON and print
logger.WithField("error", formattedJSON).Error() // or ideally, pass it directly to a logger

// format the error to a string and print it
formattedStr := eris.ToString(err, true)
fmt.Println(formattedStr)
```

`eris` also enables control over the [default format's separators](#formatting-with-custom-separators) and allows advanced users to write their own [custom output format](#writing-a-custom-output-format).

### Interpreting eris stack traces

Errors created with this package contain stack traces that are managed automatically. They're currently mandatory when creating and wrapping errors but optional when printing or logging. By default, the stack trace and all wrapped layers follow the opposite order of Go's `runtime` package, which means that the original calling method is shown first and the root cause of the error is shown last.

```golang
{
  "root":{
    "message":"error bad request", // root cause
    "stack":[
      "main.main:/Users/roti/go/src/github.com/rotisserie/eris/examples/logging/example.go:143", // original calling method
      "main.ProcessResource:/Users/roti/go/src/github.com/rotisserie/eris/examples/logging/example.go:71",
      "main.(*Request).Validate:/Users/roti/go/src/github.com/rotisserie/eris/examples/logging/example.go:29", // location of Wrap call
      "main.(*Request).Validate:/Users/roti/go/src/github.com/rotisserie/eris/examples/logging/example.go:28" // location of the root
    ]
  },
  "wrap":[
    {
      "message":"received a request with no ID", // additional context
      "stack":"main.(*Request).Validate:/Users/roti/go/src/github.com/rotisserie/eris/examples/logging/example.go:29" // location of Wrap call
    }
  ]
}
```

### Inverting the stack trace and error output

If you prefer some other order than the default, `eris` supports inverting both the stack trace and the entire error output. When both are inverted, the root error is shown first and the original calling method is shown last.

```golang
// create a default format with error and stack inversion options
format := eris.NewDefaultStringFormat(eris.FormatOptions{
  InvertOutput: true, // flag that inverts the error output (wrap errors shown first)
  WithTrace: true,    // flag that enables stack trace output
  InvertTrace: true,  // flag that inverts the stack trace output (top of call stack shown first)
})

// format the error to a string and print it
formattedStr := eris.ToCustomString(err, format)
fmt.Println(formattedStr)

// example output:
// error not found
//   main.GetResource:/Users/roti/go/src/github.com/rotisserie/eris/examples/logging/example.go:52
//   main.ProcessResource:/Users/roti/go/src/github.com/rotisserie/eris/examples/logging/example.go:76
//   main.main:/Users/roti/go/src/github.com/rotisserie/eris/examples/logging/example.go:143
// failed to get resource 'res1'
//   main.GetResource:/Users/roti/go/src/github.com/rotisserie/eris/examples/logging/example.go:52
```

### Inspecting errors

The `eris` package provides a couple ways to inspect and compare error types. [`eris.Is`](https://pkg.go.dev/github.com/rotisserie/eris#Is) returns true if a particular error appears anywhere in the error chain. Currently, it works simply by comparing error messages with each other. If an error contains a particular message (e.g. `"error not found"`) anywhere in its chain, it's defined to be that error type.

```golang
ErrNotFound := eris.New("error not found")
_, err := db.Get(id)
// check if the resource was not found
if eris.Is(err, ErrNotFound) {
  // return the error with some useful context
  return eris.Wrapf(err, "error getting resource '%v'", id)
}
```

[`eris.As`](https://pkg.go.dev/github.com/rotisserie/eris#As) finds the first error in a chain that matches a given target. If there's a match, it sets the target to that error value and returns true.

```golang
var target *NotFoundError
_, err := db.Get(id)
// check if the error is a NotFoundError type
if errors.As(err, &target) {
    // err is a *NotFoundError and target is set to the error's value
    return target
}
```

[`eris.Cause`](https://pkg.go.dev/github.com/rotisserie/eris#Cause) unwraps an error until it reaches the cause, which is defined as the first (i.e. root) error in the chain.

```golang
ErrNotFound := eris.New("error not found")
_, err := db.Get(id)
// compare the cause to some sentinel value
if eris.Cause(err) == ErrNotFound {
  // return the error with some useful context
  return eris.Wrapf(err, "error getting resource '%v'", id)
}
```

### Formatting with custom separators

For users who need more control over the error output, `eris` allows for some control over the separators between each piece of the output via the [`eris.Format`](https://pkg.go.dev/github.com/rotisserie/eris#Format) type. If this isn't flexible enough for your needs, see the [custom output format](#writing-a-custom-output-format) section below. To format errors with custom separators, you can define and pass a format object to [`eris.ToCustomString`](https://pkg.go.dev/github.com/rotisserie/eris#ToCustomString) or [`eris.ToCustomJSON`](https://pkg.go.dev/github.com/rotisserie/eris#ToCustomJSON).

```golang
// format the error to a string with custom separators
formattedStr := eris.ToCustomString(err, Format{
  FormatOptions: eris.FormatOptions{
    WithTrace: true,   // flag that enables stack trace output
  },
  MsgStackSep: "\n",   // separator between error messages and stack frame data
  PreStackSep: "\t",   // separator at the beginning of each stack frame
  StackElemSep: " | ", // separator between elements of each stack frame
  ErrorSep: "\n",      // separator between each error in the chain
})
fmt.Println(formattedStr)

// example output:
// error reading file 'example.json'
//   main.readFile | .../example/main.go | 6
// unexpected EOF
//   main.main | .../example/main.go | 20
//   main.parseFile | .../example/main.go | 12
//   main.readFile | .../example/main.go | 6
```

### Writing a custom output format

`eris` also allows advanced users to construct custom error strings or objects in case the default error doesn't fit their requirements. The [`UnpackedError`](https://pkg.go.dev/github.com/rotisserie/eris#UnpackedError) object provides a convenient and developer friendly way to store and access existing error traces. The `ErrRoot` and `ErrChain` fields correspond to the root error and wrap error chain, respectively. If a root error wraps an external error, that error will be default formatted and assigned to the `ErrExternal` field. If any other error type is unpacked, it will appear in the `ErrExternal` field. You can access all of the information contained in an error via [`eris.Unpack`](https://pkg.go.dev/github.com/rotisserie/eris#Unpack).

```golang
// get the unpacked error object
uErr := eris.Unpack(err)
// send only the root error message to a logging server instead of the complete error trace
sentry.CaptureMessage(uErr.ErrRoot.Msg)
```

### Sending error traces to Sentry

`eris` supports sending your error traces to [Sentry](https://sentry.io/) using the Sentry Go [client SDK](https://github.com/getsentry/sentry-go). You can run the example that generated the following output on Sentry UI using the command `go run examples/sentry/example.go -dsn=<DSN>`.

```
*eris.wrapError: test: wrap 1: wrap 2: wrap 3
  File "main.go", line 19, in Example
    return eris.New("test")
  File "main.go", line 23, in WrapExample
    err := Example()
  File "main.go", line 25, in WrapExample
    return eris.Wrap(err, "wrap 1")
  File "main.go", line 31, in WrapSecondExample
    err := WrapExample()
  File "main.go", line 33, in WrapSecondExample
    return eris.Wrap(err, "wrap 2")
  File "main.go", line 44, in main
    err := WrapSecondExample()
  File "main.go", line 45, in main
    err = eris.Wrap(err, "wrap 3")
```

## Comparison to other packages (e.g. pkg/errors)

### Error formatting and stack traces

Readability is a major design requirement for `eris`. In addition to the JSON output shown above, `eris` also supports formatting errors to a simple string.

```
failed to get resource 'res1'
  main.GetResource:/Users/roti/go/src/github.com/rotisserie/eris/examples/logging/example.go:52
error not found
  main.main:/Users/roti/go/src/github.com/rotisserie/eris/examples/logging/example.go:143
  main.ProcessResource:/Users/roti/go/src/github.com/rotisserie/eris/examples/logging/example.go:76
  main.GetResource:/Users/roti/go/src/github.com/rotisserie/eris/examples/logging/example.go:52
```

The `eris` error stack is designed to be easier to interpret than other error handling packages, and it achieves this by omitting extraneous information and avoiding unnecessary repetition. The stack trace above omits calls from Go's `runtime` package and includes just a single frame for wrapped layers which are inserted into the root error stack trace in the correct order. `eris` also correctly handles and updates stack traces for global error values in a transparent way.

The output of `pkg/errors` for the same error is shown below. In this case, the root error stack trace is incorrect because it was declared as a global value, and it includes several extraneous lines from the `runtime` package. The output is also much more difficult to read and does not allow for custom formatting.

```
error not found
main.init
  /Users/roti/go/src/github.com/rotisserie/eris/examples/logging/example.go:18
runtime.doInit
  /usr/local/Cellar/go/1.13.6/libexec/src/runtime/proc.go:5222
runtime.main
  /usr/local/Cellar/go/1.13.6/libexec/src/runtime/proc.go:190
runtime.goexit
  /usr/local/Cellar/go/1.13.6/libexec/src/runtime/asm_amd64.s:1357
failed to get resource 'res1'
main.GetResource
  /Users/roti/go/src/github.com/rotisserie/eris/examples/logging/example.go:52
main.ProcessResource
  /Users/roti/go/src/github.com/rotisserie/eris/examples/logging/example.go:76
main.main
  /Users/roti/go/src/github.com/rotisserie/eris/examples/logging/example.go:143
runtime.main
  /usr/local/Cellar/go/1.13.6/libexec/src/runtime/proc.go:203
runtime.goexit
  /usr/local/Cellar/go/1.13.6/libexec/src/runtime/asm_amd64.s:1357
```

## Migrating to eris

Migrating to `eris` should be a very simple process. If it doesn't offer something that you currently use from existing error packages, feel free to submit an issue to us. If you don't want to refactor all of your error handling yet, `eris` should work relatively seamlessly with your existing error types. Please submit an issue if this isn't the case for some reason.

Many of your dependencies will likely still use [pkg/errors](https://github.com/pkg/errors) for error handling. When external error types are wrapped with additional context, `eris` creates a new root error that wraps the original external error. Because of this, error inspection should work seamlessly with other error libraries.

## Contributing

If you'd like to contribute to `eris`, we'd love your input! Please submit an issue first so we can discuss your proposal.

-------------------------------------------------------------------------------

Released under the [MIT License].

[MIT License]: LICENSE.txt
[eris-logo]: https://cdn.emojidex.com/emoji/hdpi/minecraft_golden_apple.png?1511637499
[doc-img]: https://pkg.go.dev/badge/github.com/rotisserie/eris
[doc]: https://pkg.go.dev/github.com/rotisserie/eris
[ci-img]: https://github.com/rotisserie/eris/workflows/build/badge.svg
[ci]: https://github.com/rotisserie/eris/actions
[report-img]: https://goreportcard.com/badge/github.com/rotisserie/eris
[report]: https://goreportcard.com/report/github.com/rotisserie/eris
[cov-img]: https://codecov.io/gh/rotisserie/eris/branch/master/graph/badge.svg
[cov]: https://codecov.io/gh/rotisserie/eris
