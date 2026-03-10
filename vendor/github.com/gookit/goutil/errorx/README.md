# ErrorX

`errorx` provide an enhanced error implements for go, allow with stacktraces and wrap another error.

## Install

```go
go get github.com/gookit/goutil/errorx
```

## Go docs

- [Go docs](https://pkg.go.dev/github.com/gookit/goutil/errorx)

## Usage

### Create error with call stack info

- use the `errorx.New` instead `errors.New`

```go
func doSomething() error {
    if false {
	    // return errors.New("a error happen")
	    return errorx.New("a error happen")
	}
}
```

- use the `errorx.Newf` or `errorx.Errorf` instead `fmt.Errorf`

```go
func doSomething() error {
    if false {
	    // return fmt.Errorf("a error %s", "happen")
	    return errorx.Newf("a error %s", "happen")
	}
}
```

### Wrap the previous error

used like this before:

```go
    if err := SomeFunc(); err != nil {
	    return err
	}
```

can be replaced with:

```go
    if err := SomeFunc(); err != nil {
	    return errors.Stacked(err)
	}
```

## Output details

error output details for use `errorx`

### Use errorx.New

`errorx` functions for new error:

```go
func New(msg string) error
func Newf(tpl string, vars ...interface{}) error
func Errorf(tpl string, vars ...interface{}) error
```

Examples:

```go
    err := errorx.New("the error message")

    fmt.Println(err)
    // fmt.Printf("%v\n", err)
    // fmt.Printf("%#v\n", err)
```

> from the test: `errorx_test.TestNew()`

**Output**:

```text
the error message
STACK:
github.com/gookit/goutil/errorx_test.returnXErr()
  /Users/inhere/Workspace/godev/gookit/goutil/errorx/errorx_test.go:21
github.com/gookit/goutil/errorx_test.returnXErrL2()
  /Users/inhere/Workspace/godev/gookit/goutil/errorx/errorx_test.go:25
github.com/gookit/goutil/errorx_test.TestNew()
  /Users/inhere/Workspace/godev/gookit/goutil/errorx/errorx_test.go:29
testing.tRunner()
  /usr/local/Cellar/go/1.18/libexec/src/testing/testing.go:1439
runtime.goexit()
  /usr/local/Cellar/go/1.18/libexec/src/runtime/asm_amd64.s:1571
```

### Use errorx.With

`errorx` functions for with another error:

```go
func With(err error, msg string) error
func Withf(err error, tpl string, vars ...interface{}) error
```

With a go raw error:

```go
	err1 := returnErr("first error message")

	err2 := errorx.With(err1, "second error message")
	fmt.Println(err2)
```

> from the test: `errorx_test.TestWith_goerr()`

**Output**:

```text
second error message
STACK:
github.com/gookit/goutil/errorx_test.TestWith_goerr()
  /Users/inhere/Workspace/godev/gookit/goutil/errorx/errorx_test.go:51
testing.tRunner()
  /usr/local/Cellar/go/1.18/libexec/src/testing/testing.go:1439
runtime.goexit()
  /usr/local/Cellar/go/1.18/libexec/src/runtime/asm_amd64.s:1571

Previous: first error message
```

With a `errorx` error:

```go
	err1 := returnXErr("first error message")
	err2 := errorx.With(err1, "second error message")
	fmt.Println(err2)
```

> from the test: `errorx_test.TestWith_errorx()`

**Output**:

```text
second error message
STACK:
github.com/gookit/goutil/errorx_test.TestWith_errorx()
  /Users/inhere/Workspace/godev/gookit/goutil/errorx/errorx_test.go:64
testing.tRunner()
  /usr/local/Cellar/go/1.18/libexec/src/testing/testing.go:1439
runtime.goexit()
  /usr/local/Cellar/go/1.18/libexec/src/runtime/asm_amd64.s:1571

Previous: first error message
STACK:
github.com/gookit/goutil/errorx_test.returnXErr()
  /Users/inhere/Workspace/godev/gookit/goutil/errorx/errorx_test.go:21
github.com/gookit/goutil/errorx_test.TestWith_errorx()
  /Users/inhere/Workspace/godev/gookit/goutil/errorx/errorx_test.go:61
testing.tRunner()
  /usr/local/Cellar/go/1.18/libexec/src/testing/testing.go:1439
runtime.goexit()
  /usr/local/Cellar/go/1.18/libexec/src/runtime/asm_amd64.s:1571

```

### Use errorx.Wrap

```go
err := errors.New("first error message")
err = errorx.Wrap(err, "second error message")
err = errorx.Wrap(err, "third error message")
// fmt.Println(err)
// fmt.Println(err.Error())
```

Direct print the `err`:

```text
third error message
Previous: second error message
Previous: first error message
```

Print the `err.Error()`:

```text
third error message; second error message; first error message
```

## Code Check & Testing

```bash
gofmt -w -l ./
golint ./...
```

**Testing**:

```shell
go test -v ./errorx/...
```

**Test limit by regexp**:

```shell
go test -v -run ^TestSetByKeys ./errorx/...
```

## Refers

- golang errors
- https://github.com/joomcode/errorx
- https://github.com/pkg/errors
- https://github.com/juju/errors
- https://github.com/go-errors/errors
