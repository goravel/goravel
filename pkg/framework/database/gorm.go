package database

import (
	"fmt"
	"github.com/goravel/framework/support/facades"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
	"log"
)

type Gorm struct {
}

func (g *Gorm) Init() *gorm.DB {
	var (
		host     = facades.Config.GetString("database.mysql.host")
		port     = facades.Config.GetString("database.mysql.port")
		database = facades.Config.GetString("database.mysql.database")
		username = facades.Config.GetString("database.mysql.username")
		password = facades.Config.GetString("database.mysql.password")
		charset  = facades.Config.GetString("database.mysql.charset")
	)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=%t&loc=%s",
		username, password, host, port, database, charset, true, "Local")

	gormConfig := mysql.New(mysql.Config{
		DSN: dsn,
	})

	var level gormLogger.LogLevel
	if facades.Config.GetBool("app.debug") {
		level = gormLogger.Info
	} else {
		level = gormLogger.Error
	}

	db, err := gorm.Open(gormConfig, &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
		Logger:                                   gormLogger.Default.LogMode(level),
	})

	if err != nil {
		log.Fatal("Init DB fail:" + err.Error())
	}

	return db
}
