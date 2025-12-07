package facades

import (
	"github.com/goravel/framework/contracts/database/schema"
)

func Schema() schema.Schema {
	return App().MakeSchema()
}
