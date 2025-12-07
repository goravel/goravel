package facades

import (
	"github.com/goravel/framework/contracts/validation"
)

func Validation() validation.Validation {
	return App().MakeValidation()
}
