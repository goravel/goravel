package facades

import (
	"github.com/goravel/framework/contracts/database/db"
)

func DB() db.DB {
	return App().MakeDB()
}
