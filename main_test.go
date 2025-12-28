package main

import (
	"os/exec"
	"testing"

	"github.com/goravel/framework/support/file"
	"github.com/goravel/framework/support/path"
	"github.com/stretchr/testify/suite"

	"goravel/app/facades"
	"goravel/bootstrap"
)

type MainTestSuite struct {
	suite.Suite
}

func TestMainTestSuite(t *testing.T) {
	suite.Run(t, new(MainTestSuite))
}

func (s *MainTestSuite) SetupSuite() {
	bootstrap.Boot()

	// There is no .env file in the testing environment, so we need to initialize the configuration first.
	facades.Config().Add("app.name", "Goravel")
	facades.Config().Add("app.env", "local")
	facades.Config().Add("app.debug", "true")
}

func (s *MainTestSuite) TearDownTest() {
	s.NoError(exec.Command("git", "checkout", ".").Run())
	s.NoError(exec.Command("git", "clean", "-fd").Run())
	s.NoError(exec.Command("go", "mod", "tidy").Run())
}

func (s *MainTestSuite) TestPackageInstall_All() {
	s.NoError(facades.Artisan().Call("package:install --all --default --dev"))

	s.NoError(exec.Command("go", "run", ".", "artisan").Run())

	s.NoError(facades.Artisan().Call("package:uninstall Auth"))
	s.NoError(facades.Artisan().Call("package:uninstall Telemetry"))
	s.NoError(facades.Artisan().Call("package:uninstall Testing"))
	s.NoError(facades.Artisan().Call("package:uninstall Grpc"))
	s.NoError(facades.Artisan().Call("package:uninstall Hash"))
	s.NoError(facades.Artisan().Call("package:uninstall Route"))
	s.NoError(facades.Artisan().Call("package:uninstall Http"))
	s.NoError(facades.Artisan().Call("package:uninstall View"))
	s.NoError(facades.Artisan().Call("package:uninstall Session"))
	s.NoError(facades.Artisan().Call("package:uninstall Storage"))
	s.NoError(facades.Artisan().Call("package:uninstall Validation"))
	s.NoError(facades.Artisan().Call("package:uninstall Lang"))
	s.NoError(facades.Artisan().Call("package:uninstall Mail"))
	s.NoError(facades.Artisan().Call("package:uninstall Process"))
	s.NoError(facades.Artisan().Call("package:uninstall Crypt"))
	s.NoError(facades.Artisan().Call("package:uninstall RateLimiter"))
	s.NoError(facades.Artisan().Call("package:uninstall Schedule"))
	s.NoError(facades.Artisan().Call("package:uninstall Gate"))
	s.NoError(facades.Artisan().Call("package:uninstall Cache"))
	s.NoError(facades.Artisan().Call("package:uninstall Event"))
	s.NoError(facades.Artisan().Call("package:uninstall Queue"))
	s.NoError(facades.Artisan().Call("package:uninstall Schema"))
	s.NoError(facades.Artisan().Call("package:uninstall Seeder"))
	s.NoError(facades.Artisan().Call("package:uninstall DB"))
	s.NoError(facades.Artisan().Call("package:uninstall Orm"))
	s.NoError(facades.Artisan().Call("package:uninstall Log"))
}

func (s *MainTestSuite) TestPackageInstall_Auth() {
	s.NoError(facades.Artisan().Call("package:install Auth --default --dev"))
	s.FileExists(path.Facade("auth.go"))
	s.FileExists(path.Config("auth.go"))
	s.True(file.Contains(path.Bootstrap("providers.go"), "&auth.ServiceProvider{},"))

	s.NoError(facades.Artisan().Call("package:uninstall Auth"))
	s.NoFileExists(path.Facade("auth.go"))
	s.NoFileExists(path.Config("auth.go"))
	s.False(file.Contains(path.Bootstrap("providers.go"), "&auth.ServiceProvider{},"))
}

func (s *MainTestSuite) TestPackageInstall_Cache() {
	s.NoError(facades.Artisan().Call("package:install Cache --default --dev"))
	s.FileExists(path.Facade("cache.go"))
	s.FileExists(path.Config("cache.go"))
	s.True(file.Contains(path.Bootstrap("providers.go"), "&cache.ServiceProvider{},"))

	s.NoError(facades.Artisan().Call("package:uninstall Cache"))
	s.NoFileExists(path.Facade("cache.go"))
	s.NoFileExists(path.Config("cache.go"))
	s.False(file.Contains(path.Bootstrap("providers.go"), "&cache.ServiceProvider{},"))
}

func (s *MainTestSuite) TestPackageInstall_Crypt() {
	s.NoError(facades.Artisan().Call("package:install Crypt --default --dev"))
	s.FileExists(path.Facade("crypt.go"))
	s.True(file.Contains(path.Bootstrap("providers.go"), "&crypt.ServiceProvider{},"))

	s.NoError(facades.Artisan().Call("package:uninstall Crypt"))
	s.NoFileExists(path.Facade("crypt.go"))
	s.False(file.Contains(path.Bootstrap("providers.go"), "&crypt.ServiceProvider{},"))
}

func (s *MainTestSuite) TestPackageInstall_DB() {
	s.NoError(facades.Artisan().Call("package:install DB --default --dev"))
	s.FileExists(path.Facade("db.go"))
	s.FileExists(path.Config("database.go"))
	s.True(file.Contains(path.Bootstrap("providers.go"), "&database.ServiceProvider{},"))
	s.True(file.Contains(path.Bootstrap("providers.go"), "&postgres.ServiceProvider{},"))
	s.True(file.Contains(path.Config("database.go"), "postgres"))

	s.NoError(facades.Artisan().Call("package:uninstall DB"))
	s.NoFileExists(path.Facade("db.go"))

	// The Schema facade still exists, so database.go and ServiceProvider should still exist.
	s.FileExists(path.Config("database.go"))
	s.True(file.Contains(path.Bootstrap("providers.go"), "&database.ServiceProvider{},"))
	s.True(file.Contains(path.Bootstrap("providers.go"), "&postgres.ServiceProvider{},"))
	s.True(file.Contains(path.Config("database.go"), "postgres"))
}

func (s *MainTestSuite) TestPackageInstall_Event() {
	s.NoError(facades.Artisan().Call("package:install Event --default --dev"))
	s.FileExists(path.Facade("event.go"))
	s.True(file.Contains(path.Bootstrap("providers.go"), "&event.ServiceProvider{},"))

	s.NoError(facades.Artisan().Call("package:uninstall Event"))
	s.NoFileExists(path.Facade("event.go"))
	s.False(file.Contains(path.Bootstrap("providers.go"), "&event.ServiceProvider{},"))
}

func (s *MainTestSuite) TestPackageInstall_Gate() {
	s.NoError(facades.Artisan().Call("package:install Gate --default --dev"))
	s.FileExists(path.Facade("gate.go"))
	s.FileExists(path.Config("auth.go"))
	s.True(file.Contains(path.Bootstrap("providers.go"), "&auth.ServiceProvider{},"))

	s.NoError(facades.Artisan().Call("package:uninstall Gate"))
	s.NoFileExists(path.Facade("gate.go"))
	s.NoFileExists(path.Config("auth.go"))
	s.False(file.Contains(path.Bootstrap("providers.go"), "&auth.ServiceProvider{},"))
}

func (s *MainTestSuite) TestPackageInstall_Grpc() {
	s.NoError(facades.Artisan().Call("package:install Grpc --default --dev"))
	s.FileExists(path.Facade("grpc.go"))
	s.FileExists(path.Config("grpc.go"))
	s.FileExists(path.Route("grpc.go"))
	s.True(file.Contains(path.Bootstrap("providers.go"), "&grpc.ServiceProvider{},"))
	s.True(file.Contains(path.Bootstrap("app.go"), ".Grpc"))
	s.True(file.Contains(path.Base(".env.example"), `
GRPC_HOST=
GRPC_PORT=
`))

	s.NoError(facades.Artisan().Call("package:uninstall Grpc"))
	s.NoFileExists(path.Facade("grpc.go"))
	s.NoFileExists(path.Config("grpc.go"))
	s.NoFileExists(path.Route("grpc.go"))
	s.False(file.Contains(path.Bootstrap("providers.go"), "&grpc.ServiceProvider{},"))
	s.False(file.Contains(path.Bootstrap("app.go"), ".Grpc"))
}

func (s *MainTestSuite) TestPackageInstall_Hash() {
	s.NoError(facades.Artisan().Call("package:install Hash --default --dev"))
	s.FileExists(path.Facade("hash.go"))
	s.FileExists(path.Config("hashing.go"))
	s.True(file.Contains(path.Bootstrap("providers.go"), "&hash.ServiceProvider{},"))

	s.NoError(facades.Artisan().Call("package:uninstall Hash"))
	s.NoFileExists(path.Facade("hash.go"))
	s.NoFileExists(path.Config("hashing.go"))
	s.False(file.Contains(path.Bootstrap("providers.go"), "&hash.ServiceProvider{},"))
}

func (s *MainTestSuite) TestPackageInstall_Http() {
	s.NoError(facades.Artisan().Call("package:install Http --default --dev"))
	s.FileExists(path.Facade("http.go"))
	s.FileExists(path.Config("http.go"))
	s.FileExists(path.Config("jwt.go"))
	s.FileExists(path.Config("cors.go"))
	s.True(file.Contains(path.Bootstrap("providers.go"), "&http.ServiceProvider{},"))

	s.NoError(facades.Artisan().Call("package:uninstall Http"))
	s.NoFileExists(path.Facade("http.go"))
	s.NoFileExists(path.Config("http.go"))
	s.NoFileExists(path.Config("jwt.go"))
	s.NoFileExists(path.Config("cors.go"))
	s.False(file.Contains(path.Bootstrap("providers.go"), "&http.ServiceProvider{},"))
}

func (s *MainTestSuite) TestPackageInstall_Lang() {
	s.NoError(facades.Artisan().Call("package:install Lang --default --dev"))
	s.FileExists(path.Facade("lang.go"))
	s.True(file.Contains(path.Bootstrap("providers.go"), "&translation.ServiceProvider{},"))

	s.NoError(facades.Artisan().Call("package:uninstall Lang"))
	s.NoFileExists(path.Facade("lang.go"))
	s.False(file.Contains(path.Bootstrap("providers.go"), "&translation.ServiceProvider{},"))
}

func (s *MainTestSuite) TestPackageInstall_Log() {
	s.NoError(facades.Artisan().Call("package:install Log --default --dev"))
	s.FileExists(path.Facade("log.go"))
	s.FileExists(path.Config("logging.go"))
	s.True(file.Contains(path.Bootstrap("providers.go"), "&log.ServiceProvider{},"))
	s.True(file.Contains(path.Base(".env.example"), `
LOG_CHANNEL=stack
LOG_LEVEL=debug
`))

	s.NoError(facades.Artisan().Call("package:uninstall Log"))
	s.NoFileExists(path.Facade("log.go"))
	s.NoFileExists(path.Config("logging.go"))
	s.NoFileExists(path.Route("logging.go"))
	s.False(file.Contains(path.Bootstrap("providers.go"), "&log.ServiceProvider{},"))
}

func (s *MainTestSuite) TestPackageInstall_Mail() {
	s.NoError(facades.Artisan().Call("package:install Mail --default --dev"))
	s.FileExists(path.Facade("mail.go"))
	s.FileExists(path.Config("mail.go"))
	s.True(file.Contains(path.Bootstrap("providers.go"), "&mail.ServiceProvider{},"))
	s.True(file.Contains(path.Base(".env.example"), `
MAIL_HOST=
MAIL_PORT=
MAIL_USERNAME=
MAIL_PASSWORD=
MAIL_FROM_ADDRESS=
MAIL_FROM_NAME=
`))

	s.NoError(facades.Artisan().Call("package:uninstall Mail"))
	s.NoFileExists(path.Facade("mail.go"))
	s.NoFileExists(path.Config("mail.go"))
	s.NoFileExists(path.Route("mail.go"))
	s.False(file.Contains(path.Bootstrap("providers.go"), "&mail.ServiceProvider{},"))
}

func (s *MainTestSuite) TestPackageInstall_Orm() {
	s.NoError(facades.Artisan().Call("package:install Orm --default --dev"))
	s.FileExists(path.Facade("orm.go"))
	s.FileExists(path.Config("database.go"))
	s.True(file.Contains(path.Bootstrap("providers.go"), "&database.ServiceProvider{},"))

	// Schema should be uninstalled if Orm wants to be uninstalled.
	s.NoError(facades.Artisan().Call("package:uninstall Schema"))
	s.NoError(facades.Artisan().Call("package:uninstall Orm"))
	s.NoFileExists(path.Facade("orm.go"))
	s.NoFileExists(path.Config("database.go"))
	s.False(file.Contains(path.Bootstrap("providers.go"), "&database.ServiceProvider{},"))
}

func (s *MainTestSuite) TestPackageInstall_Process() {
	s.NoError(facades.Artisan().Call("package:install Process --default --dev"))
	s.FileExists(path.Facade("process.go"))
	s.True(file.Contains(path.Bootstrap("providers.go"), "&process.ServiceProvider{},"))

	s.NoError(facades.Artisan().Call("package:uninstall Process"))
	s.NoFileExists(path.Facade("process.go"))
	s.False(file.Contains(path.Bootstrap("providers.go"), "&process.ServiceProvider{},"))
}

func (s *MainTestSuite) TestPackageInstall_Queue() {
	s.NoError(facades.Artisan().Call("package:install Queue --default --dev"))
	s.FileExists(path.Facade("queue.go"))
	s.FileExists(path.Config("queue.go"))
	s.FileExists(path.Migration("20210101000001_create_jobs_table.go"))
	s.True(file.Contains(path.Bootstrap("providers.go"), "&queue.ServiceProvider{},"))
	s.True(file.Contains(path.Bootstrap("app.go"), "Migrations()"))
	s.True(file.Contains(path.Bootstrap("migrations.go"), "M20210101000001CreateJobsTable{}"))

	s.NoError(facades.Artisan().Call("package:uninstall Queue"))
	s.NoFileExists(path.Facade("queue.go"))
	s.NoFileExists(path.Config("queue.go"))
	s.NoFileExists(path.Migration("20210101000001_create_jobs_table.go"))
	s.False(file.Contains(path.Bootstrap("providers.go"), "&queue.ServiceProvider{},"))
	s.True(file.Contains(path.Bootstrap("app.go"), "Migrations()"))
	s.False(file.Contains(path.Bootstrap("migrations.go"), "M20210101000001CreateJobsTable{}"))
}

func (s *MainTestSuite) TestPackageInstall_RateLimiter() {
	s.NoError(facades.Artisan().Call("package:install RateLimiter --default --dev"))
	s.FileExists(path.Facade("rate_limiter.go"))
	s.FileExists(path.Config("http.go"))
	s.FileExists(path.Config("jwt.go"))
	s.FileExists(path.Config("cors.go"))
	s.True(file.Contains(path.Bootstrap("providers.go"), "&http.ServiceProvider{},"))

	s.NoError(facades.Artisan().Call("package:uninstall RateLimiter"))
	s.NoFileExists(path.Facade("rate_limiter.go"))
	s.NoFileExists(path.Config("http.go"))
	s.NoFileExists(path.Config("jwt.go"))
	s.NoFileExists(path.Config("cors.go"))
	s.False(file.Contains(path.Bootstrap("providers.go"), "&http.ServiceProvider{},"))
}

func (s *MainTestSuite) TestPackageInstall_Route() {
	s.NoError(facades.Artisan().Call("package:install Route --default --dev"))
	s.True(file.Contains(path.Bootstrap("providers.go"), "&route.ServiceProvider{},"))
	s.True(file.Contains(path.Bootstrap("providers.go"), "&gin.ServiceProvider{},"))
	s.True(file.Contains(path.Bootstrap("app.go"), ".Web"))
	s.True(file.Contains(path.Config("http.go"), "gin"))
	s.FileExists(path.Resource("views", "welcome.tmpl"))
	s.FileExists(path.Route("web.go"))
	s.FileExists(path.Facade("route.go"))
	s.True(file.Contains(path.Base(".env.example"), `
APP_URL=http://localhost
APP_HOST=127.0.0.1
APP_PORT=3000

JWT_SECRET=
`))

	s.NoError(facades.Artisan().Call("package:uninstall Route"))
	s.False(file.Contains(path.Bootstrap("providers.go"), "&route.ServiceProvider{},"))
	s.False(file.Contains(path.Bootstrap("app.go"), ".Web"))
	// The Http facade still exists, so "gin" related content in http.go should still exist.
	s.True(file.Contains(path.Bootstrap("providers.go"), "&gin.ServiceProvider{},"))
	s.True(file.Contains(path.Config("http.go"), "gin"))
	s.NoFileExists(path.Resource("views", "welcome.tmpl"))
	s.NoFileExists(path.Route("web.go"))
	s.NoFileExists(path.Facade("route.go"))
}

func (s *MainTestSuite) TestPackageInstall_Schedule() {
	s.NoError(facades.Artisan().Call("package:install Schedule --default --dev"))
	s.FileExists(path.Facade("schedule.go"))
	s.True(file.Contains(path.Bootstrap("providers.go"), "&schedule.ServiceProvider{},"))

	s.NoError(facades.Artisan().Call("package:uninstall Schedule"))
	s.NoFileExists(path.Facade("schedule.go"))
	s.False(file.Contains(path.Bootstrap("providers.go"), "&schedule.ServiceProvider{},"))
}

func (s *MainTestSuite) TestPackageInstall_Schema() {
	s.NoError(facades.Artisan().Call("package:install Schema --default --dev"))
	s.FileExists(path.Facade("schema.go"))
	s.FileExists(path.Config("database.go"))
	s.True(file.Contains(path.Bootstrap("providers.go"), "&database.ServiceProvider{},"))

	s.NoError(facades.Artisan().Call("package:uninstall Schema"))
	s.NoFileExists(path.Facade("schema.go"))

	// The Orm facade still exists, so database.go and ServiceProvider should still exist.
	s.FileExists(path.Config("database.go"))
	s.True(file.Contains(path.Bootstrap("providers.go"), "&database.ServiceProvider{},"))
}

func (s *MainTestSuite) TestPackageInstall_Seeder() {
	s.NoError(facades.Artisan().Call("package:install Seeder --default --dev"))
	s.FileExists(path.Facade("seeder.go"))
	s.FileExists(path.Config("database.go"))
	s.True(file.Contains(path.Bootstrap("providers.go"), "&database.ServiceProvider{},"))
	s.True(file.Contains(path.Config("database.go"), "default"))
	s.True(file.Contains(path.Config("database.go"), "connections"))
	s.True(file.Contains(path.Config("database.go"), "pool"))
	s.True(file.Contains(path.Config("database.go"), "slow_threshold"))
	s.True(file.Contains(path.Config("database.go"), "migrations"))
	s.True(file.Contains(path.Base(".env.example"), `
DB_HOST=
DB_PORT=
DB_DATABASE=
DB_USERNAME=
DB_PASSWORD=
`))

	s.NoError(facades.Artisan().Call("package:uninstall Seeder"))
	s.NoFileExists(path.Facade("seeder.go"))
	s.NoFileExists(path.Config("database.go"))
	s.False(file.Contains(path.Bootstrap("providers.go"), "&database.ServiceProvider{},"))
	s.False(file.Contains(path.Config("database.go"), "default"))
	s.False(file.Contains(path.Config("database.go"), "connections"))
	s.False(file.Contains(path.Config("database.go"), "pool"))
	s.False(file.Contains(path.Config("database.go"), "slow_threshold"))
	s.False(file.Contains(path.Config("database.go"), "migrations"))
}

func (s *MainTestSuite) TestPackageInstall_Session() {
	s.NoError(facades.Artisan().Call("package:install Session --default --dev"))
	s.True(file.Contains(path.Bootstrap("providers.go"), "&session.ServiceProvider{},"))
	s.FileExists(path.Config("session.go"))
	s.FileExists(path.Facade("session.go"))
	s.True(file.Contains(path.Base(".env.example"), `
SESSION_DRIVER=file
SESSION_LIFETIME=120
`))

	s.NoError(facades.Artisan().Call("package:uninstall Session"))
	s.False(file.Contains(path.Bootstrap("providers.go"), "&session.ServiceProvider{},"))
	s.NoFileExists(path.Config("session.go"))
	s.NoFileExists(path.Facade("session.go"))
}

func (s *MainTestSuite) TestPackageInstall_Storage() {
	s.NoError(facades.Artisan().Call("package:install Storage --default --dev"))
	s.True(file.Contains(path.Bootstrap("providers.go"), "&filesystem.ServiceProvider{},"))
	s.FileExists(path.Config("filesystems.go"))
	s.FileExists(path.Facade("storage.go"))

	s.NoError(facades.Artisan().Call("package:uninstall Storage"))
	s.False(file.Contains(path.Bootstrap("providers.go"), "&filesystem.ServiceProvider{},"))
	s.NoFileExists(path.Config("filesystems.go"))
	s.NoFileExists(path.Facade("storage.go"))
}

func (s *MainTestSuite) TestPackageInstall_Telemetry() {
	s.NoError(facades.Artisan().Call("package:install Telemetry --default --dev"))
	s.True(file.Contains(path.Bootstrap("providers.go"), "&telemetry.ServiceProvider{},"))
	s.FileExists(path.Config("telemetry.go"))
	s.FileExists(path.Facade("telemetry.go"))
	s.True(file.Contains(path.Config("logging.go"), "otel"))

	s.NoError(facades.Artisan().Call("package:uninstall Telemetry"))
	s.False(file.Contains(path.Bootstrap("providers.go"), "&telemetry.ServiceProvider{},"))
	s.NoFileExists(path.Config("telemetry.go"))
	s.NoFileExists(path.Facade("telemetry.go"))
	s.False(file.Contains(path.Config("logging.go"), "otel"))
}

func (s *MainTestSuite) TestPackageInstall_Testing() {
	s.NoError(facades.Artisan().Call("package:install Testing --default --dev"))
	s.True(file.Contains(path.Bootstrap("providers.go"), "&testing.ServiceProvider{},"))
	s.FileExists(path.Test("test_case.go"))
	s.FileExists(path.Test("feature", "example_test.go"))
	s.FileExists(path.Facade("testing.go"))

	s.NoError(facades.Artisan().Call("package:uninstall Testing"))
	s.False(file.Contains(path.Bootstrap("providers.go"), "&testing.ServiceProvider{},"))
	s.NoFileExists(path.Test("test_case.go"))
	s.NoFileExists(path.Test("feature", "example_test.go"))
	s.NoFileExists(path.Facade("testing.go"))
}

func (s *MainTestSuite) TestPackageInstall_Validation() {
	s.NoError(facades.Artisan().Call("package:install Validation --default --dev"))
	s.True(file.Contains(path.Bootstrap("providers.go"), "&validation.ServiceProvider{},"))
	s.FileExists(path.Facade("validation.go"))

	s.NoError(facades.Artisan().Call("package:uninstall Validation"))
	s.False(file.Contains(path.Bootstrap("providers.go"), "&validation.ServiceProvider{},"))
	s.NoFileExists(path.Facade("validation.go"))
}

func (s *MainTestSuite) TestPackageInstall_View() {
	s.NoError(facades.Artisan().Call("package:install View --default --dev"))
	s.True(file.Contains(path.Bootstrap("providers.go"), "&view.ServiceProvider{},"))
	s.FileExists(path.Facade("view.go"))

	s.NoError(facades.Artisan().Call("package:uninstall View"))
	s.False(file.Contains(path.Bootstrap("providers.go"), "&view.ServiceProvider{},"))
	s.NoFileExists(path.Facade("view.go"))
}

func (s *MainTestSuite) TestPackageInstall_DBDrivers() {
	s.NoError(facades.Artisan().Call("package:install DB --default --dev"))

	s.NoError(facades.Artisan().Call("package:install github.com/goravel/mysql --dev"))
	s.True(file.Contains(path.Bootstrap("providers.go"), "&mysql.ServiceProvider{},"))
	s.True(file.Contains(path.Config("database.go"), "mysql"))

	s.NoError(facades.Artisan().Call("package:uninstall github.com/goravel/mysql"))
	s.False(file.Contains(path.Bootstrap("providers.go"), "&mysql.ServiceProvider{},"))
	s.False(file.Contains(path.Config("database.go"), "mysql"))

	s.NoError(facades.Artisan().Call("package:install github.com/goravel/sqlserver --dev"))
	s.True(file.Contains(path.Bootstrap("providers.go"), "&sqlserver.ServiceProvider{},"))
	s.True(file.Contains(path.Config("database.go"), "sqlserver"))

	s.NoError(facades.Artisan().Call("package:uninstall github.com/goravel/sqlserver"))
	s.False(file.Contains(path.Bootstrap("providers.go"), "&sqlserver.ServiceProvider{},"))
	s.False(file.Contains(path.Config("database.go"), "sqlserver"))

	s.NoError(facades.Artisan().Call("package:install github.com/goravel/sqlite --dev"))
	s.True(file.Contains(path.Bootstrap("providers.go"), "&sqlite.ServiceProvider{},"))
	s.True(file.Contains(path.Config("database.go"), "sqlite"))

	s.NoError(facades.Artisan().Call("package:uninstall github.com/goravel/sqlite"))
	s.False(file.Contains(path.Bootstrap("providers.go"), "&sqlite.ServiceProvider{},"))
	s.False(file.Contains(path.Config("database.go"), "sqlite"))
}

func (s *MainTestSuite) TestPackageInstall_FilesystemDrivers() {
	s.NoError(facades.Artisan().Call("package:install Storage --default --dev"))

	s.NoError(facades.Artisan().Call("package:install github.com/goravel/s3 --dev"))
	s.True(file.Contains(path.Bootstrap("providers.go"), "&s3.ServiceProvider{},"))
	s.True(file.Contains(path.Config("filesystems.go"), "s3"))

	s.NoError(facades.Artisan().Call("package:uninstall github.com/goravel/s3"))
	s.False(file.Contains(path.Bootstrap("providers.go"), "&s3.ServiceProvider{},"))
	s.False(file.Contains(path.Config("filesystems.go"), "s3"))

	s.NoError(facades.Artisan().Call("package:install github.com/goravel/oss --dev"))
	s.True(file.Contains(path.Bootstrap("providers.go"), "&oss.ServiceProvider{},"))
	s.True(file.Contains(path.Config("filesystems.go"), "oss"))

	s.NoError(facades.Artisan().Call("package:uninstall github.com/goravel/oss"))
	s.False(file.Contains(path.Bootstrap("providers.go"), "&oss.ServiceProvider{},"))
	s.False(file.Contains(path.Config("filesystems.go"), "oss"))

	s.NoError(facades.Artisan().Call("package:install github.com/goravel/cos --dev"))
	s.True(file.Contains(path.Bootstrap("providers.go"), "&cos.ServiceProvider{},"))
	s.True(file.Contains(path.Config("filesystems.go"), "cos"))

	s.NoError(facades.Artisan().Call("package:uninstall github.com/goravel/cos"))
	s.False(file.Contains(path.Bootstrap("providers.go"), "&cos.ServiceProvider{},"))
	s.False(file.Contains(path.Config("filesystems.go"), "cos"))

	s.NoError(facades.Artisan().Call("package:install github.com/goravel/minio --dev"))
	s.True(file.Contains(path.Bootstrap("providers.go"), "&minio.ServiceProvider{},"))
	s.True(file.Contains(path.Config("filesystems.go"), "minio"))

	s.NoError(facades.Artisan().Call("package:uninstall github.com/goravel/minio"))
	s.False(file.Contains(path.Bootstrap("providers.go"), "&minio.ServiceProvider{},"))
	s.False(file.Contains(path.Config("filesystems.go"), "minio"))
}

func (s *MainTestSuite) TestPackageInstall_CacheDrivers() {
	s.NoError(facades.Artisan().Call("package:install Cache --default --dev"))
	s.NoError(facades.Artisan().Call("package:install Session --default --dev"))
	s.NoError(facades.Artisan().Call("package:install Queue --default --dev"))

	s.NoError(facades.Artisan().Call("package:install github.com/goravel/redis --dev"))
	s.True(file.Contains(path.Bootstrap("providers.go"), "&redis.ServiceProvider{},"))
	s.True(file.Contains(path.Config("cache.go"), "redis"))
	s.True(file.Contains(path.Config("database.go"), "redis"))
	s.True(file.Contains(path.Config("queue.go"), "redis"))
	s.True(file.Contains(path.Config("session.go"), "redis"))

	s.NoError(facades.Artisan().Call("package:uninstall github.com/goravel/redis"))
	s.False(file.Contains(path.Bootstrap("providers.go"), "&redis.ServiceProvider{},"))
	s.False(file.Contains(path.Config("cache.go"), "redis"))
	s.False(file.Contains(path.Config("database.go"), "redis"))
	s.False(file.Contains(path.Config("queue.go"), "redis"))
	s.False(file.Contains(path.Config("session.go"), "redis"))
}
