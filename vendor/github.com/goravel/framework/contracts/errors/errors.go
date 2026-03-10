package errors

// Error is the interface that wraps the basic error methods
type Error interface {
	// Args allows setting arguments for the placeholders in the text
	Args(...any) Error
	// Error implements the error interface and formats the error string
	Error() string
	// SetModule explicitly sets the module in the error message
	SetModule(string) Error
}
