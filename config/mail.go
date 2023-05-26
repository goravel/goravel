package config

import "github.com/goravel/framework/facades"

func init() {
	config := facades.Config()
	config.Add("mail", map[string]any{
		// SMTP Host Address
		//
		// Here you may provide the host address of the SMTP server used by your
		// applications. A default option is provided that is compatible with
		// the Mailgun mail service which will provide reliable deliveries.
		"host": config.Env("MAIL_HOST", ""),

		// SMTP Host Port
		//
		// This is the SMTP port used by your application to deliver e-mails to
		// users of the application. Like the host we have set this value to
		// stay compatible with the Mailgun e-mail application by default.
		"port": config.Env("MAIL_PORT", 587),

		// --------------------------------------------------------------------------
		// Global "From" Address
		// --------------------------------------------------------------------------
		//
		// You may wish for all e-mails sent by your application to be sent from
		// the same address. Here, you may specify a name and address that is
		// used globally for all e-mails that are sent by your application.
		"from": map[string]any{
			"address": config.Env("MAIL_FROM_ADDRESS", "hello@example.com"),
			"name":    config.Env("MAIL_FROM_NAME", "Example"),
		},

		// SMTP Server Username
		//
		// If your SMTP server requires a username for authentication, you should
		// set it here. This will get used to authenticate with your server on
		// connection. You may also set the "password" value below this one.
		"username": config.Env("MAIL_USERNAME"),

		"password": config.Env("MAIL_PASSWORD"),
	})
}
