package session

// Session is the interface that defines the methods that should be implemented by a session.
type Session interface {
	// All returns all attributes of the session.
	All() map[string]any
	// Exists checks if a key exists in the session attributes.
	Exists(key string) bool
	// Flash sets a flash data value in the session attributes.
	Flash(key string, value any) Session
	// Flush clears all attributes from the session.
	Flush() Session
	// Forget removes specified keys from the session attributes.
	Forget(keys ...string) Session
	// Get retrieves the value of a key from the session attributes.
	Get(key string, defaultValue ...any) any
	// GetName returns the name of the session.
	GetName() string
	// GetID returns the ID of the session.
	GetID() string
	// Has checks if a key exists and is not nil in the session attributes.
	Has(key string) bool
	// Invalidate invalidates the session.
	Invalidate() error
	// Keep reflash a subset of the current flash data.
	Keep(keys ...string) Session
	// Missing checks if a key is missing in the session attributes.
	Missing(key string) bool
	// Now flashes a key and value for immediate use.
	Now(key string, value any) Session
	// Only retrieves the specified keys and their values from the session attributes.
	Only(keys []string) map[string]any
	// Pull retrieves and removes the value of a key from the session attributes.
	Pull(key string, defaultValue ...any) any
	// Put sets the value of a key in the session attributes.
	Put(key string, value any) Session
	// Reflash keeps all the flash data for an additional request.
	Reflash() Session
	// Regenerate regenerates the session.
	Regenerate(destroy ...bool) error
	// Remove removes the value of a key from the session attributes.
	Remove(key string) any
	// Save saves the session.
	Save() error
	// SetDriver sets the session driver
	SetDriver(driver Driver) Session
	// SetID sets the ID of the session.
	SetID(id string) Session
	// SetName sets the name of the session.
	SetName(name string) Session
	// Start initiates the session.
	Start() bool
	// Token retrieves the session token.
	Token() string
}
