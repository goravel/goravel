package config

import (
	"github.com/goravel/framework/facades"
	"goravel/app/models"
)

func init() {
	config := facades.Config
	config.Add("auth", map[string]interface{}{
		//Authentication Defaults
		//
		//This option controls the default authentication "guard"
		//reset options for your application. You may change these defaults
		//as required, but they're a perfect start for most applications.
		"defaults": map[string]interface{}{
			"guard": "user",
		},

		//Authentication Guards
		//
		//Next, you may define every authentication guard for your application.
		//Of course, a great default configuration has been defined for you
		//here which uses session storage and the Eloquent user provider.
		//
		//All authentication drivers have a user provider. This defines how the
		//users are actually retrieved out of your database or other storage
		//mechanisms used by this application to persist your user's data.
		//
		//Supported: "jwt"
		"guards": map[string]interface{}{
			"user": map[string]interface{}{
				"driver":   "jwt",
				"provider": "users",
			},
		},

		//User Providers
		//
		//All authentication drivers have a user provider. This defines how the
		//users are actually retrieved out of your database or other storage
		//mechanisms used by this application to persist your user's data.
		//
		//If you have multiple user tables or models you may configure multiple
		//sources which represent each model / table. These sources may then
		//be assigned to any extra authentication guards you have defined.
		//
		//Supported: "database", "orm"
		"providers": map[string]interface{}{
			"users": map[string]interface{}{
				"driver": "orm",
				"model":  models.User{},
			},
			//"users": map[string]interface{}{
			//	"driver": "database",
			//	"table":  "users",
			//},
		},
	})
}
