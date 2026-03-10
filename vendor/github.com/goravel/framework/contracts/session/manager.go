package session

type Manager interface {
	// BuildSession constructs a new session with the given handler and session ID.
	BuildSession(handler Driver, sessionID ...string) (Session, error)
	// Driver retrieves the session driver by name.
	Driver(name ...string) (Driver, error)
	// ReleaseSession releases the session back to the pool.
	ReleaseSession(session Session)
}
