package schema

// MorphKeyType represents the type of key used for morph relationships
type MorphKeyType string

const (
	MorphKeyTypeInt  MorphKeyType = "int"
	MorphKeyTypeUuid MorphKeyType = "uuid"
	MorphKeyTypeUlid MorphKeyType = "ulid"
)

var defaultMorphKeyType MorphKeyType = MorphKeyTypeInt

// SetDefaultMorphKeyType sets the default morph key type
func SetDefaultMorphKeyType(keyType MorphKeyType) {
	defaultMorphKeyType = keyType
}

// GetDefaultMorphKeyType returns the current default morph key type
func GetDefaultMorphKeyType() MorphKeyType {
	return defaultMorphKeyType
}

// MorphUsingUuids sets the default morph key type to UUID
func MorphUsingUuids() {
	defaultMorphKeyType = MorphKeyTypeUuid
}

// MorphUsingUlids sets the default morph key type to ULID
func MorphUsingUlids() {
	defaultMorphKeyType = MorphKeyTypeUlid
}

// MorphUsingInts sets the default morph key type to int (default)
func MorphUsingInts() {
	defaultMorphKeyType = MorphKeyTypeInt
}
