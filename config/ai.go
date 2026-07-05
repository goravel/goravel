package config

import (
	"github.com/goravel/framework/contracts/ai"
	openaifacades "github.com/goravel/openai/facades"
	"goravel/app/facades"
)

func init() {
	config := facades.Config()
	config.Add("ai", map[string]any{
		// Default AI Provider
		//
		// This option controls the default AI provider that will be used.
		"default": config.Env("AI_PROVIDER"),

		// AI Providers
		//
		// Here you may configure each AI provider used by your application.
		// A variety of drivers are available, and each provider may also
		// configure the models available to your application.
		"providers": map[string]any{
			"openai": map[string]any{
				"key": config.Env("OPENAI_API_KEY", ""),
				"models": map[string]any{
					"text": map[string]any{
						"default": "",
					},
					"audio": map[string]any{
						"default": "",
					},
					"transcription": map[string]any{
						"default": "",
					},
					"image": map[string]any{
						"default": "",
					},
				},
				"failover": map[string][]string{},
				"url":      config.Env("OPENAI_BASE_URL", ""),
				"via": func() (ai.Provider, error) {
					return openaifacades.OpenAI("openai")
				},
			},
		},
	})
}
