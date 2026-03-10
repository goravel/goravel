package contracts

import (
	contractsconfig "github.com/goravel/framework/contracts/config"
)

type ConfigBuilder interface {
	Config() contractsconfig.Config
	Connection() string
	Readers() []FullConfig
	Writers() []FullConfig
}

// Replacer replacer interface like strings.Replacer
type Replacer interface {
	Replace(name string) string
}

// Config Used in config/database.go
type Config struct {
	Dsn      string
	Host     string
	Port     int
	Database string
	Username string
	Password string
	Schema   string
}

// FullConfig Fill the default value for Config
type FullConfig struct {
	Config
	Driver       string
	Connection   string
	Prefix       string
	Singular     bool
	Sslmode      string
	Timezone     string
	NoLowerCase  bool
	NameReplacer Replacer
}
