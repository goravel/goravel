package json

import (
	"github.com/goravel/framework/errors"
	"github.com/goravel/framework/foundation"
)

// Marshal serializes the given value to a JSON-encoded byte slice.
func Marshal(v any) ([]byte, error) {
	app := foundation.App
	if app == nil {
		return nil, errors.ApplicationNotSet
	}

	return app.GetJson().Marshal(v)
}

// Unmarshal deserializes the given JSON-encoded byte slice into the provided value.
func Unmarshal(data []byte, v any) error {
	app := foundation.App
	if app == nil {
		return errors.ApplicationNotSet
	}

	return app.GetJson().Unmarshal(data, v)
}

// MarshalString serializes the given value to a JSON-encoded string.
func MarshalString(v any) (string, error) {
	app := foundation.App
	if app == nil {
		return "", errors.ApplicationNotSet
	}

	return app.GetJson().MarshalString(v)
}

// UnmarshalString deserializes the given JSON-encoded string into the provided value.
func UnmarshalString(data string, v any) error {
	app := foundation.App
	if app == nil {
		return errors.ApplicationNotSet
	}

	return app.GetJson().UnmarshalString(data, v)
}
