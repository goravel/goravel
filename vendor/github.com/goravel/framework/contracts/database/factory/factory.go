package factory

type Factory interface {
	// Definition defines the model's default state.
	Definition() map[string]any
}

type Model interface {
	// Factory creates a new factory instance for the model.
	Factory() Factory
}
