package seeder

type Facade interface {
	// Register registers seeders.
	Register(seeders []Seeder)
	// GetSeeder gets a seeder instance from the seeders.
	GetSeeder(name string) Seeder
	// GetSeeders gets all the seeders
	GetSeeders() []Seeder
	// Call executes the specified seeder(s).
	Call(seeders []Seeder) error
	// CallOnce executes the specified seeder(s) only once.
	CallOnce(seeders []Seeder) error
}

type Seeder interface {
	// Signature the unique signature of the seeder.
	Signature() string
	// Run executes the seeder logic.
	Run() error
}
