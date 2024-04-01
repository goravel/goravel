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
		// to immediately expire when the browser is closed, then you may
		// indicate that via the expire_on_close configuration option.
		"lifetime": config.Env("SESSION_LIFETIME", 120),

		"expire_on_close": config.Env("SESSION_EXPIRE_ON_CLOSE", false),

		// Session File Location
		//
		// When using the file session driver, we need a location where the
		// session files may be stored. A default has been set for you, but a
		// different location may be specified. This is only needed for file sessions.
		"files": path.Storage("framework/sessions"),

		// Session Sweeping Lottery
		//
		// Some session drivers must manually sweep their storage location to get
		// rid of old sessions from storage. Here are the chances out of 100 that
		// the sweeper will sweep the storage location. The default is 2 out of 100.
		"lottery": []int{2, 100},

		// Session Cookie Name
		//
		// Here you may change the name of the cookie used to identify a session
		// in the application. The name specified here will get used every time
		// a new session cookie is created by the framework for every driver.
		"cookie": config.Env("SESSION_COOKIE", str.Of(config.GetString("app.name")).Snake().Lower().String()+"_session"),

		// Session Cookie Path
		//
		// The session cookie path determines the path for which the cookie will
		// be regarded as available.Typically, this will be the root path of
		// your application, but you are free to change this when necessary.
		"path": config.Env("SESSION_PATH", "/"),

		// Session Cookie Domain
		//
		// Here you may change the domain of the cookie used to identify a session
		// in your application.This will determine which domains the cookie is
		// available to in your application.A sensible default has been set.
		"domain": config.Env("SESSION_DOMAIN", ""),

		// HTTPS Only Cookies
		//
		// By setting this option to true, session cookies will only be sent back
		// to the server if the browser has an HTTPS connection. This will keep
		// the cookie from being sent to you if it cannot be done securely.
		"secure": config.Env("SESSION_SECURE", false),

		// HTTP Access Only
		//
		// Setting this to true will prevent JavaScript from accessing the value of
		// the cookie, and the cookie will only be accessible through the HTTP protocol.
		"http_only": config.Env("SESSION_HTTP_ONLY", true),

		// Same-Site Cookies
		//
		// This option determines how your cookies behave when cross-site requests
		// take place, and can be used to mitigate CSRF attacks.By default, we
		// will set this value to "lax" since this is a secure default value.
		"same_site": config.Env("SESSION_SAME_SITE", "lax"),
	})
}
