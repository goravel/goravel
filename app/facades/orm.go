package facades

import (
	"github.com/goravel/framework/contracts/database/orm"
)

func Orm() orm.Orm {
	return App().MakeOrm()
}
