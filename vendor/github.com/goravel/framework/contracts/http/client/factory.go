package client

type Factory interface {
	// Request embeds the Request interface, allowing direct usage like Http.Get().
	Request

	// AllowStrayRequests permits specific URL patterns to bypass the mock firewall.
	AllowStrayRequests(patterns []string) Factory

	// AssertNotSent verifies that no request matching the given assertion was sent.
	AssertNotSent(assertion func(req Request) bool) bool

	// AssertNothingSent verifies that no HTTP requests were sent at all.
	AssertNothingSent() bool

	// AssertSent verifies that at least one request matching the given assertion was sent.
	AssertSent(assertion func(req Request) bool) bool

	// AssertSentCount verifies that the specific number of requests matching the criteria were sent.
	AssertSentCount(count int) bool

	// Client returns a new request builder.
	// If name is provided, it returns the configuration for that specific client.
	// If no name is provided, it returns the default client.
	Client(name ...string) Request

	// Fake registers the mock rules for testing.
	Fake(mocks map[string]any) Factory

	// PreventStrayRequests enforces that all sent requests must match a defined mock rule.
	PreventStrayRequests() Factory

	// Reset restores the factory to its original state, clearing all mocks.
	Reset()

	// Response returns a builder for creating stubbed responses.
	Response() FakeResponse

	// Sequence returns a builder for defining ordered mock responses.
	Sequence() FakeSequence
}
