package validation

import "github.com/gookit/validate"

type Errors struct {
	errors validate.Errors
}

func NewErrors(errors validate.Errors) *Errors {
	return &Errors{errors}
}

func (r *Errors) One(key ...string) string {
	if len(key) > 0 {
		errors := r.Get(key[0])
		for _, err := range errors {
			return err
		}
	}

	return r.errors.One()
}

func (r *Errors) Get(key string) map[string]string {
	return r.errors.Field(key)
}

func (r *Errors) All() map[string]map[string]string {
	return r.errors.All()
}

func (r *Errors) Has(key string) bool {
	return r.errors.HasField(key)
}
