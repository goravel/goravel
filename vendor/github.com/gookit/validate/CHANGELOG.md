# CHANGE LOG

## V2 - TODO

- [ ] inner validators always use reflect.Value as param. 

**old:**

```go
// Gt check value greater dst value. only check for: int(X), uint(X), float(X)
func Gt(val any, dstVal int64) bool {
```

**v2 new:**

```go
// Gt check value greater dst value. only check for: int(X), uint(X), float(X)
func Gt(val, dstVal any) bool {
	return gt(reflect.ValueOf(val), reflect.ValueOf(dstVal))
}

// internal implements
func gt(val reflect.Value, dstVal reflect.Value) bool
```

- can register custom type
- use sync.Pool for optimize create Validation.

```go
// Validation definition
type Validation struct {
	// for optimize create instance. refer go-playground/validator
	v *Validation
	pool *sync.Pool
    
    // ...
}

	v.pool = &sync.Pool{
		New: func() any {
			return &Validation{
				v: v,
			}
		},
	}
```

- all data operate move to DataSource

```go
type DataFace interface {
	BindStruct() error
	SafeVal(field string) any

...
}

```