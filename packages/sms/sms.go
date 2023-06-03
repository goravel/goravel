package sms

import (
	"fmt"

	"github.com/goravel/framework/contracts/config"
)

type Sms struct {
	config config.Config
}

func NewSms(config config.Config) *Sms {
	return &Sms{config: config}
}

func (s *Sms) Send() {
	fmt.Println(s.config.Get("app.key"), s.config.Get("sms.driver"))
}
