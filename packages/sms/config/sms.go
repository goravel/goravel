package config

import (
	"github.com/goravel/framework/facades"
)

func init() {
	config := facades.Config()
	config.Add("sms", map[string]any{
		"driver": "aliyun",
	})
}
