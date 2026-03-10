package comdef

import (
	"errors"
	"strings"
)

// ErrConvType error
var ErrConvType = errors.New("convert value type error")

// Errors multi error list
type Errors []error

// Error string
func (es Errors) Error() string {
	var sb strings.Builder
	for _, err := range es {
		sb.WriteString(err.Error())
		sb.WriteByte('\n')
	}
	return sb.String()
}

// ErrOrNil error
func (es Errors) ErrOrNil() error {
	if len(es) == 0 {
		return nil
	}
	return es
}

// First error
func (es Errors) First() error {
	if len(es) > 0 {
		return es[0]
	}
	return nil
}
