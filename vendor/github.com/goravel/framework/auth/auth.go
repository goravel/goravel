package auth

import (
	"fmt"
	"sync"

	contractsauth "github.com/goravel/framework/contracts/auth"
	"github.com/goravel/framework/contracts/config"
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/log"
	"github.com/goravel/framework/errors"
)

var (
	guardFuncs     = sync.Map{}
	providersFuncs = sync.Map{}
)

type Auth struct {
	contractsauth.GuardDriver
	config           config.Config
	ctx              http.Context
	log              log.Log
	defaultGuardName string
}

func NewAuth(ctx http.Context, config config.Config, log log.Log) (*Auth, error) {
	auth := &Auth{
		config: config,
		ctx:    ctx,
		log:    log,
	}

	auth.Extend("jwt", NewJwtGuard)
	auth.Extend("session", NewSessionGuard)
	auth.Provider("orm", NewOrmUserProvider)

	defaultGuardName := config.GetString("auth.defaults.guard")
	auth.defaultGuardName = defaultGuardName

	if ctx != nil {
		defaultGuard := auth.guard(defaultGuardName)
		auth.GuardDriver = defaultGuard
	}

	return auth, nil
}

func (r *Auth) Extend(name string, fn contractsauth.GuardFunc) {
	guardFuncs.Store(name, fn)
}

func (r *Auth) Guard(name string) contractsauth.GuardDriver {
	if name == "" || name == r.defaultGuardName {
		return r.GuardDriver
	}

	return r.guard(name)
}

func (r *Auth) Provider(name string, fn contractsauth.UserProviderFunc) {
	providersFuncs.Store(name, fn)
}

func (r *Auth) createUserProvider(name string) (contractsauth.UserProvider, error) {
	driverName := r.config.GetString(fmt.Sprintf("auth.providers.%s.driver", name))

	providerFunc, ok := providersFuncs.Load(driverName)
	if !ok {
		return nil, errors.AuthProviderDriverNotFound.Args(driverName, name)

	}

	provider, err := providerFunc.(contractsauth.UserProviderFunc)(r.ctx)
	if err != nil {
		return nil, err
	}

	return provider, nil
}

func (r *Auth) guard(name string) contractsauth.GuardDriver {
	driverName := r.config.GetString(fmt.Sprintf("auth.guards.%s.driver", name))
	guardFunc, ok := guardFuncs.Load(driverName)
	if !ok {
		// http recover will catch the panic and return the error to the client,
		// to avoid print repeated log.
		panic(errors.AuthGuardDriverNotFound.Args(driverName, name))
	}

	userProviderName := r.config.GetString(fmt.Sprintf("auth.guards.%s.provider", name))
	userProvider, err := r.createUserProvider(userProviderName)
	if err != nil {
		panic(err)
	}

	guard, err := guardFunc.(contractsauth.GuardFunc)(r.ctx, name, userProvider)
	if err != nil {
		panic(err)
	}

	return guard
}
