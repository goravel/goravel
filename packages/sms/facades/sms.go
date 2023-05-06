package facades

import (
	"log"

	"github.com/goravel/framework/services/sms"
	"github.com/goravel/framework/services/sms/contracts"
)

func Sms() contracts.Sms {
	instance, err := sms.App.Make("sms")
	if err != nil {
		log.Fatalln(err)
		return nil
	}

	return instance.(contracts.Sms)
}
