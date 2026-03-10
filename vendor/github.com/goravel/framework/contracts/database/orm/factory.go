package orm

type Factory interface {
	// Count sets the number of models that should be generated.
	Count(count int) Factory
	// Create creates a model and persists it to the database.
	Create(value any, attributes ...map[string]any) error
	// CreateQuietly creates a model and persists it to the database without firing any model events.
	CreateQuietly(value any, attributes ...map[string]any) error
	// Make creates a model and returns it, but does not persist it to the database.
	Make(value any, attributes ...map[string]any) error
}
