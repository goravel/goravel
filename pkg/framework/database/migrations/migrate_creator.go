package migrations

import (
	"github.com/goravel/framework/support/facades"
	"log"
	"os"
	"strings"
	"time"
)

type MigrateCreator struct {
}

func (receiver MigrateCreator) Create(name string, table string, create bool) {
	upStub, downStub := receiver.getStub(table, create)
	receiver.createFile(receiver.getPath(name, "up"), receiver.populateStub(upStub, table))
	receiver.createFile(receiver.getPath(name, "down"), receiver.populateStub(downStub, table))
}

func (receiver MigrateCreator) getStub(table string, create bool) (string, string) {
	if table == "" {
		return "", ""
	}

	if create {
		return MigrateStubs{}.CreateUp(), MigrateStubs{}.CreateDown()
	}

	return MigrateStubs{}.UpdateUp(), MigrateStubs{}.UpdateDown()
}

func (receiver MigrateCreator) populateStub(stub string, table string) string {
	stub = strings.ReplaceAll(stub, "DummyDatabaseCharset", facades.Config.GetString("database.connections."+facades.Config.GetString("database.default")+".charset"))

	if table != "" {
		stub = strings.ReplaceAll(stub, "DummyTable", table)
	}

	return stub
}

func (receiver MigrateCreator) getPath(name string, category string) string {
	pwd, _ := os.Getwd()

	return pwd + "/database/migrations/" + time.Now().Format("20060102150405") + "_" + name + "." + category + ".sql"
}

func (receiver MigrateCreator) createFile(path string, content string) {
	file, err := os.Create(path)
	if err != nil {
		log.Fatalln(err.Error())
	}

	_, err = file.WriteString(content)
	if err != nil {
		log.Fatalln(err.Error())
	}
}
