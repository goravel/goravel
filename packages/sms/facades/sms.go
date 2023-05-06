package facades

import (
	"log"

	"goravel/packages/sms"
	"goravel/packages/sms/contracts"
)

func Sms() contracts.Sms {
	instance, err := sms.App.Make(sms.Binding)
	if err != nil {
		log.Fatalln(err)
		return nil
	}

	return instance.(contracts.Sms)
}
