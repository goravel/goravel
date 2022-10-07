package testing

import (
	"fmt"
	ormcontract "github.com/goravel/framework/contracts/database/orm"
	"github.com/goravel/framework/contracts/mail"
	"github.com/goravel/framework/database/orm"
	"github.com/goravel/framework/facades"
	"goravel/app/models"
	"strconv"
	"time"
)

//func scope() func(db database.OrmTransaction) database.OrmTransaction {
//	return func(db database.OrmTransaction) database.OrmTransaction {
//		return db.Limit(3)
//	}
//}

func Cache() string {
	if err := facades.Cache.Put("name", "goravel", 1*time.Minute); err != nil {
		fmt.Println("cache.put.error", err)
	}

	return facades.Cache.Get("name", "test").(string)
}

func Config() string {
	return facades.Config.GetString("app.name", "test")
}

//func Artisan() {
//	facades.Artisan.Call("list")
//}

type Test struct {
	orm.Model
}

func Orm() error {
	if err := facades.Orm.Query().Create(&Test{}); err != nil {
		return err
	}

	var test Test
	return facades.Orm.Query().Where("id = ?", 1).Find(&test)
}

func Transaction() error {
	return facades.Orm.Transaction(func(tx ormcontract.Transaction) error {
		var test Test
		if err := tx.Create(&test); err != nil {
			return err
		}

		var test1 Test
		return tx.Where("id = ?", test.ID).Find(&test1)
	})
}

func Begin() error {
	tx, _ := facades.Orm.Query().Begin()
	user := models.User{Name: "Goravel"}
	if err := tx.Create(&user); err != nil {
		return tx.Rollback()
	} else {
		return tx.Commit()
	}
}

//func Event() error {
//	return facades.Event.Job(&events.Test{}, []eventcontract.Arg{
//		{Type: "string", Value: "abcc"},
//		{Type: "int", Value: 1234},
//	}).Dispatch()
//}

func Log() {
	facades.Log.Debug("test")
}

//func Queue() error {
//	return facades.Queue.Job(&jobs.TestJob{}, []queue.Arg{}).Dispatch()
//}

func Paginator(page string, limit string) func(methods ormcontract.Query) ormcontract.Query {
	return func(query ormcontract.Query) ormcontract.Query {
		page, _ := strconv.Atoi(page)
		limit, _ := strconv.Atoi(limit)
		offset := (page - 1) * limit

		return query.Offset(offset).Limit(limit)
	}
}

func Mail() {
	_ = facades.Mail.To([]string{"example@example.com"}).
		Cc([]string{"example@example.com"}).
		Bcc([]string{"example@example.com"}).
		Attach([]string{"file.png"}).
		Content(mail.Content{Subject: "Subject", Html: "<h1>Hello Goravel</h1>"}).
		Send()

	_ = facades.Mail.To([]string{"example@example.com"}).
		Cc([]string{"example@example.com"}).
		Bcc([]string{"example@example.com"}).
		Attach([]string{"file.png"}).
		Content(mail.Content{Subject: "Subject", Html: "<h1>Hello Goravel</h1>"}).
		Queue(&mail.Queue{Connection: "high", Queue: "mail"})
}
