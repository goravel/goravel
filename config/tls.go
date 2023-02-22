package config

import (
	"github.com/goravel/framework/facades"
)

func init() {
	config := facades.Config
	config.Add("tls", map[string]interface{}{
		"enabled": config.Env("TLS_ENABLED", false),

		// SSL certificate configuration
		"ssl": map[string]interface{}{
			// Certificate file (PEM format). Path: storage/app/private/ssl/cert/example.com.pe
			"cert_file": config.Env("TLS_SSL_CERT_FILE", "/private/ssl/cert/example.com.pem"),
			// Certificate KEY. Path: storage/app/private/ssl/cert/example.com.key
			"key_file": config.Env("TLS_SSL_KEY_FILE", "/private/ssl/cert/example.com.key"),
		},
	})
}
