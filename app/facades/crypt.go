package facades

import (
	"github.com/goravel/framework/contracts/crypt"
)

func Crypt() crypt.Crypt {
	return App().MakeCrypt()
}
