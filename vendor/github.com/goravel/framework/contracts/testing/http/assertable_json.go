package http

type AssertableJSON interface {
	// Json returns the underlying JSON data as a map.
	Json() map[string]any
	// Count asserts that the property at the given key is an array and has the expected number of items.
	Count(key string, length int) AssertableJSON
	// Each gets the array at the given key and enforces the callback's assertions on every item in that array.
	Each(key string, callback func(AssertableJSON)) AssertableJSON
	// First gets the array at the given key and enforces the callback's assertions on the first item only.
	First(key string, callback func(AssertableJSON)) AssertableJSON
	// Has asserts that the property at the given key exists.
	Has(key string) AssertableJSON
	// HasAll asserts that all the provided keys exist.
	HasAll(keys []string) AssertableJSON
	// HasAny asserts that at least one of the provided keys exists.
	HasAny(keys []string) AssertableJSON
	// HasWithScope asserts that the property at the given key is an array of the expected length,
	// and enforces the callback's assertions on the first item of that array.
	// This is useful for verifying the structure of the items within a list.
	HasWithScope(key string, length int, callback func(AssertableJSON)) AssertableJSON
	// Missing asserts that the property at the given key does not exist.
	Missing(key string) AssertableJSON
	// MissingAll asserts that none of the provided keys exist.
	MissingAll(keys []string) AssertableJSON
	// Where asserts that the property at the given key equals the expected value.
	Where(key string, value any) AssertableJSON
	// WhereNot asserts that the property at the given key does not equal the provided value.
	WhereNot(key string, value any) AssertableJSON
}
