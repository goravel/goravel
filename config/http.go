package config

import (
	"github.com/gin-gonic/gin/render"
	"github.com/goravel/framework/contracts/route"
	"github.com/goravel/framework/facades"
	"github.com/goravel/gin"
	ginfacades "github.com/goravel/gin/facades"
)

func init() {
	config := facades.Config()
	config.Add("http", map[string]any{
		// HTTP Driver
		"default": "gin",
		// HTTP Drivers
		"drivers": map[string]any{
			"gin": map[string]any{
				// Optional, default is 4096 KB
				"body_limit": 4096,
				"header_limit": 4096,
				"route": func() (route.Route, error) {
					return ginfacades.Route("gin"), nil
				},
				// Optional, default is http/template
			        "template": func() (render.HTMLRender, error) {
			            return gin.DefaultTemplate()
			        },
			},
		},
		// HTTP URL
		"url": config.Env("APP_URL", "http://localhost"),
		// HTTP Host
		"host": config.Env("APP_HOST", "127.0.0.1"),
		// HTTP Port
		"port": config.Env("APP_PORT", "3000"),
		// HTTPS Configuration
		"tls": map[string]any{
			// HTTPS Host
			"host": config.Env("APP_HOST", "127.0.0.1"),
			// HTTPS Port
			"port": config.Env("APP_PORT", "3000"),
			// SSL Certificate, you can put the certificate in /public folder
			"ssl": map[string]any{
				// ca.pem
				"cert": "",
				// ca.key
				"key": "",
			},
		},
	})
}
