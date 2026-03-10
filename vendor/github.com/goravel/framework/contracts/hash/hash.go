package hash

type Hash interface {
	// Make returns the hashed value of the given string.
	Make(value string) (string, error)
	// Check checks if the given string matches the given hash.
	Check(value string, hashedValue string) bool
	// NeedsRehash checks if the given hash needs to be rehashed.
	NeedsRehash(hashedValue string) bool
}
