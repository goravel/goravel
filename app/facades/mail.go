package facades

import (
	"github.com/goravel/framework/contracts/mail"
)

func Mail() mail.Mail {
	return App().MakeMail()
}
