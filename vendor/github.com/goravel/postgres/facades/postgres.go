package facades

import (
	"fmt"

	"github.com/goravel/framework/contracts/database/driver"

	"github.com/goravel/postgres"
)

func Postgres(connection string) (driver.Driver, error) {
	if postgres.App == nil {
		return nil, fmt.Errorf("please register postgres service provider")
	}

	instance, err := postgres.App.MakeWith(postgres.Binding, map[string]any{
		"connection": connection,
	})
	if err != nil {
		return nil, err
	}

	return instance.(*postgres.Postgres), nil
}
