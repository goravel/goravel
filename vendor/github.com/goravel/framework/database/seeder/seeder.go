package seeder

import (
	"slices"

	"github.com/goravel/framework/contracts/database/seeder"
	"github.com/goravel/framework/support/color"
)

var _ seeder.Facade = (*SeederFacade)(nil)

type SeederFacade struct {
	Seeders []seeder.Seeder
	Called  []string
}

func NewSeederFacade() *SeederFacade {
	return &SeederFacade{}
}

func (s *SeederFacade) Register(seeders []seeder.Seeder) {
	existingSignatures := make(map[string]bool)

	for _, seeder := range seeders {
		signature := seeder.Signature()

		if existingSignatures[signature] {
			color.Errorf("Duplicate seeder signature: %s in %T\n", signature, seeder)
		} else {
			existingSignatures[signature] = true
			s.Seeders = append(s.Seeders, seeder)
		}
	}
}

func (s *SeederFacade) GetSeeder(name string) seeder.Seeder {
	var seeder seeder.Seeder
	for _, item := range s.Seeders {
		if item.Signature() == name {
			seeder = item
			break
		}
	}

	return seeder
}

func (s *SeederFacade) GetSeeders() []seeder.Seeder {
	return s.Seeders
}

// Call executes the specified seeder(s).
func (s *SeederFacade) Call(seeders []seeder.Seeder) error {
	for _, seeder := range seeders {
		signature := seeder.Signature()

		if err := seeder.Run(); err != nil {
			return err
		}

		if !slices.Contains(s.Called, signature) {
			s.Called = append(s.Called, signature)
		}
	}
	return nil
}

// CallOnce executes the specified seeder(s) only if they haven't been executed before.
func (s *SeederFacade) CallOnce(seeders []seeder.Seeder) error {
	for _, item := range seeders {
		signature := item.Signature()

		if slices.Contains(s.Called, signature) {
			continue
		}

		if err := s.Call([]seeder.Seeder{item}); err != nil {
			return err
		}
	}
	return nil
}
