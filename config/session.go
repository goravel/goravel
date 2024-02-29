package config

import (
	"github.com/goravel/framework/facades"
	"github.com/goravel/framework/support/path"
	"github.com/goravel/framework/support/str"
)

func init() {
	config := facades.Config()
	config.Add("session", map[string]any{
		// Default Session Driver
		//
		// This option controls the default session "driver" that will be used on
		// requests. By default, we will use the lightweight file session driver, but you
		// may specify any of the other wonderful drivers provided here.
		//
		// Supported: "file"
		"driver": config.Env("SESSION_DRIVER", "file"),

		// Session Lifetime
		//
		// Here you may specify the number of minutes that you wish the session
		// to be allowed to remain idle before it expires. If you want them
		// to immediately expire on the browser closing, set that option.
		"lifetime": config.Env("SESSION_LIFETIME", 120),

		// Session File Location
		//
		// When using the file session driver, we need a location where the
		// session files may be stored. A default has been set for you, but a
		// different location may be specified. This is only needed for file sessions.
		"files": path.Storage("framework/sessions"),

		// Session Cookie Name
		//
		// Here you may change the name of the cookie used to identify a session
		// in the application. The name specified here will get used every time
		// a new session cookie is created by the framework for every driver.
		"cookie": config.Env("SESSION_COOKIE", str.Of(config.GetString("app.name")).Lower().Snake("_").String()+"_session"),
	})
}
