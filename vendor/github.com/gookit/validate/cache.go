package validate

/*
// TODO: idea for optimize ...
// - use pool for Validation
// - cache struct reflect value

// Usage:
// // use a single instance of Validate, it caches struct info
// vf := validate.NewFactory()
// v := vf.New(data)
// v.Validate()

// global factory for create Validation
var gf = newFactory()

// factory for create Validation instances
type factory struct {
	// pool for Validation instances
	pool sync.Pool
	// m atomic.Value // map[reflect.Type]*cStruct
}

func newFactory() *factory {
	f := &factory{}
	f.pool.New = func() any {
		return newValidation(nil)
	}

	return f
}

func (f *factory) get() *Validation {
	v := f.pool.Get().(*Validation)
	// TODO clear something
	return v
}

func (f *factory) put(v *Validation) {
	f.pool.Put(v)
}
*/

// func test() {
// 	i32 := driver.Int32
// 	fmt.Println(i32)
// }

/*
// TODO cache struct reflect value, tags and more
type structMeta struct {
}

type cache struct {
	m sync.Map
	// m atomic.Value
	// map[reflect.Type]*cStruct
}

func (c *cache) get(rt reflect.Type) *structMeta {
	// key := rt.PkgPath() + rt.Name()
	return c.m.Load(rt)
}

func (c *cache) set(rt reflect.Type, meta structMeta)  {
	// key := rt.PkgPath() + rt.Name()
	c.m.Store(rt, data)
}
*/
