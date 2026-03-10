package foundation

type Json interface {
	// Marshal serializes the given value to a JSON-encoded byte slice.
	Marshal(any) ([]byte, error)
	// Unmarshal deserializes the given JSON-encoded byte slice into the provided value.
	Unmarshal([]byte, any) error
	// MarshalString serializes the given value to a JSON-encoded string.
	MarshalString(any) (string, error)
	// UnmarshalString deserializes the given JSON-encoded string into the provided value.
	UnmarshalString(string, any) error
}
