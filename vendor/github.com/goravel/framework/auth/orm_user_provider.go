package auth

import (
	"gorm.io/gorm/clause"

	contractsauth "github.com/goravel/framework/contracts/auth"
	"github.com/goravel/framework/contracts/database/orm"
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/errors"
	"github.com/goravel/framework/support/database"
)

var _ contractsauth.UserProviderFunc = NewOrmUserProvider

type OrmUserProvider struct {
	ctx http.Context
	orm orm.Orm
}

func NewOrmUserProvider(ctx http.Context) (contractsauth.UserProvider, error) {
	if ormFacade == nil {
		return nil, errors.OrmFacadeNotSet.SetModule(errors.ModuleAuth)
	}

	return &OrmUserProvider{
		ctx: ctx,
		orm: ormFacade,
	}, nil
}

// GetID implements auth.UserProvider.
func (r *OrmUserProvider) GetID(user any) (any, error) {
	return database.GetID(user), nil
}

// RetriveByID implements auth.UserProvider.
func (r *OrmUserProvider) RetriveByID(user any, id any) error {
	return r.orm.WithContext(r.ctx).Query().FindOrFail(user, clause.Eq{Column: clause.PrimaryColumn, Value: id})
}
