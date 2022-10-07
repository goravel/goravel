package testing

import (
	"fmt"
	"github.com/stretchr/testify/suite"
	"goravel/bootstrap"
	"io/ioutil"

	"strings"
	"testing"
	"time"

	"github.com/goravel/framework/facades"
	"github.com/goravel/framework/support/file"
	"github.com/stretchr/testify/assert"
)

type ArtisanTestSuite struct {
	suite.Suite
}

func TestArtisanTestSuite(t *testing.T) {
	kernel := HttpKernel{}
	kernel.Create()

	bootstrap.Boot()

	suite.Run(t, new(ArtisanTestSuite))
}

func (s *ArtisanTestSuite) SetupTest() {

}

func (s *ArtisanTestSuite) TestKeyGenerate() {
	t := s.T()
	Equal(t, "key:generate", "Application key set successfully")
	Equal(t, "key:generate", "Exist application key")
}

func (s *ArtisanTestSuite) TestList() {
	t := s.T()
	NotEmpty(t, "list")
}

func (s *ArtisanTestSuite) TestHelp() {
	t := s.T()
	NotEmpty(t, "help migrate")
}

func (s *ArtisanTestSuite) TestMakeCommand() {
	t := s.T()
	Equal(t, "make:command SendEmails", "Console command created successfully")
	assert.True(t, file.Exist("./app/console/commands/send_emails.go"))
	assert.True(t, file.Remove("./app"))
}

func (s *ArtisanTestSuite) TestCommand() {
	t := s.T()
	expect := "Run test command success, argument_0: argument0, argument_1: argument1, option_name: Goravel, option_age: 18, arguments: argument0,argument1"

	facades.Artisan.Call("test --name Goravel argument0 argument1")
	log := fmt.Sprintf("storage/logs/goravel-%s.log", time.Now().Format("2006-01-02"))
	assert.True(t, file.Exist(log))
	data, err := ioutil.ReadFile(log)
	assert.Nil(t, err)
	assert.True(t, strings.Contains(string(data), expect))
	assert.True(t, file.Remove("./storage"))

	Equal(t, "test --name Goravel argument0 argument1", expect)
	Equal(t, "test -n Goravel argument0 argument1", expect)
	assert.True(t, file.Remove("storage"))
}

type HttpKernel struct {
}

func (r *HttpKernel) stub() string {
	return `package console

import (
	"github.com/goravel/framework/contracts/console"
	"github.com/goravel/framework/contracts/schedule"
	
	"goravel/app/console/commands"
)

type Kernel struct {
}

func (kernel *Kernel) Schedule() []schedule.Event {
	return []schedule.Event{}
}

func (kernel *Kernel) Commands() []console.Command {
	return []console.Command{
		&commands.Test{},
	}
}
`
}

func (r *HttpKernel) Create() {
	path := "../app/console/kernel.go"
	file.Remove(path)
	file.Create(path, r.stub())
}
