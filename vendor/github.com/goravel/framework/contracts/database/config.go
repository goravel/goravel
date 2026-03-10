package database

import "gorm.io/gorm"

type Pool struct {
	Readers []Config
	Writers []Config
}

type Config struct {
	Dialector    gorm.Dialector
	NameReplacer Replacer
	Charset      string
	Connection   string
	Dsn          string
	Database     string
	Driver       string
	Host         string
	Password     string
	Prefix       string
	Schema       string
	Sslmode      string
	Timezone     string
	Username     string
	Port         int
	NoLowerCase  bool
	Singular     bool
}

// Replacer replacer interface like strings.Replacer
type Replacer interface {
	Replace(name string) string
}
