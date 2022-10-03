package testing

import (
	"testing"

	"github.com/goravel/framework/facades"
	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	assert.Equal(t, "Goravel", facades.Config.GetString("app.name", "laravel"))
	assert.Equal(t, true, facades.Config.GetBool("app.debug", false))
	assert.Equal(t, 587, facades.Config.GetInt("mail.port", 123))
	assert.Equal(t, "Goravel", facades.Config.Env("APP_NAME", "laravel").(string))
	assert.Equal(t, "Goravel", facades.Config.Env("MAIL_HOST", "Goravel").(string))
}
