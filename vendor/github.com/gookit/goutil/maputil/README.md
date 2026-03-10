# Map Utils

`maputil` provide map data util functions. eg: convert, sub-value get, simple merge

- use `map[string]any` as Data
- deep get value by key path
- deep set value by key path

## Install

```bash
go get github.com/gookit/goutil/maputil
```

## Go docs

- [Go docs](https://pkg.go.dev/github.com/gookit/goutil/maputil)

## Usage

### Deep get value

```go
mp := map[string]any {
	"top1": "val1",
	"arr1": []string{"ab", "cd"}
	"map1": map[string]any{
	    "sub1": "val2",	
    },
}

fmt.Println(maputil.DeepGet(mp, "map1.sub1")) // Output: VAL3

// get value from slice.
fmt.Println(maputil.DeepGet(mp, "arr1.1")) // Output: cd
fmt.Println(maputil.DeepGet(mp, "arr1[1]")) // Output: cd
```

### Deep set value

```go
mp := map[string]any {
	"top1": "val1",
	"arr1": []string{"ab"}
	"map1": map[string]any{
	    "sub1": "val2",	
    },
}

err := maputil.SetByPath(&mp, "map1.newKey", "VAL3")

fmt.Println(maputil.DeepGet(mp, "map1.newKey")) // Output: VAL3
```

## Code Check & Testing

```bash
gofmt -w -l ./
golint ./...
```

**Testing**:

```shell
go test -v ./maputil/...
```

**Test limit by regexp**:

```shell
go test -v -run ^TestSetByKeys ./maputil/...
```
