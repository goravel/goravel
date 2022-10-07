package testing

import (
	"github.com/stretchr/testify/suite"
	"goravel/bootstrap"
	"testing"

	"github.com/goravel/framework/facades"
	"github.com/stretchr/testify/assert"
)

type ConfigTestSuite struct {
	suite.Suite
}

func TestConfigTestSuite(t *testing.T) {
	bootstrap.Boot()

	suite.Run(t, new(ConfigTestSuite))
}

func (s *ConfigTestSuite) SetupTest() {

}

func (s *ConfigTestSuite) TestConfig() {
	t := s.T()
	assert.Equal(t, "Goravel", facades.Config.GetString("app.name", "laravel"))
	assert.Equal(t, true, facades.Config.GetBool("app.debug", false))
	assert.Equal(t, 587, facades.Config.GetInt("mail.port", 123))
	assert.Equal(t, "Goravel", facades.Config.Env("APP_NAME", "laravel").(string))
	assert.Equal(t, "Goravel", facades.Config.Env("GRPC_HOST", "Goravel").(string))
}
