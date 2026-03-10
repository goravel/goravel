package foundation

import (
	"context"
	"fmt"
	"sync"

	contractsauth "github.com/goravel/framework/contracts/auth"
	contractsaccess "github.com/goravel/framework/contracts/auth/access"
	"github.com/goravel/framework/contracts/binding"
	contractscache "github.com/goravel/framework/contracts/cache"
	contractsconfig "github.com/goravel/framework/contracts/config"
	contractsconsole "github.com/goravel/framework/contracts/console"
	contractscrypt "github.com/goravel/framework/contracts/crypt"
	contractsdb "github.com/goravel/framework/contracts/database/db"
	contractsorm "github.com/goravel/framework/contracts/database/orm"
	contractsmigration "github.com/goravel/framework/contracts/database/schema"
	contractsseerder "github.com/goravel/framework/contracts/database/seeder"
	contractsevent "github.com/goravel/framework/contracts/event"
	"github.com/goravel/framework/contracts/facades"
	contractsfilesystem "github.com/goravel/framework/contracts/filesystem"
	contractsfoundation "github.com/goravel/framework/contracts/foundation"
	contractsgrpc "github.com/goravel/framework/contracts/grpc"
	contractshash "github.com/goravel/framework/contracts/hash"
	contractshttp "github.com/goravel/framework/contracts/http"
	contractshttpclient "github.com/goravel/framework/contracts/http/client"
	contractslog "github.com/goravel/framework/contracts/log"
	contractsmail "github.com/goravel/framework/contracts/mail"
	contractsprocess "github.com/goravel/framework/contracts/process"
	contractsqueue "github.com/goravel/framework/contracts/queue"
	contractsroute "github.com/goravel/framework/contracts/route"
	contractsschedule "github.com/goravel/framework/contracts/schedule"
	contractsession "github.com/goravel/framework/contracts/session"
	contractstelemetry "github.com/goravel/framework/contracts/telemetry"
	contractstesting "github.com/goravel/framework/contracts/testing"
	contractstranslation "github.com/goravel/framework/contracts/translation"
	contractsvalidation "github.com/goravel/framework/contracts/validation"
	contractsview "github.com/goravel/framework/contracts/view"
	"github.com/goravel/framework/support/color"
)

type instance struct {
	concrete any
	shared   bool
}

type Container struct {
	bindings  sync.Map
	instances sync.Map
}

func NewContainer() *Container {
	return &Container{}
}

func (r *Container) Bind(key any, callback func(app contractsfoundation.Application) (any, error)) {
	r.bindings.Store(key, instance{concrete: callback, shared: false})
}

func (r *Container) Bindings() []any {
	var bindings []any
	r.bindings.Range(func(key, value any) bool {
		bindings = append(bindings, key)
		return true
	})
	return bindings
}

func (r *Container) BindWith(key any, callback func(app contractsfoundation.Application, parameters map[string]any) (any, error)) {
	r.bindings.Store(key, instance{concrete: callback, shared: false})
}

func (r *Container) Fresh(bindings ...any) {
	if len(bindings) == 0 {
		r.instances.Range(func(key, value any) bool {
			if key != binding.Config {
				r.instances.Delete(key)
			}

			return true
		})
	} else {
		for _, binding := range bindings {
			r.instances.Delete(binding)
		}
	}
}

func (r *Container) Instance(key any, ins any) {
	r.bindings.Store(key, instance{concrete: ins, shared: true})
}

func (r *Container) Make(key any) (any, error) {
	return r.make(key, nil)
}

func (r *Container) MakeArtisan() contractsconsole.Artisan {
	instance, err := r.Make(facades.FacadeToBinding[facades.Artisan])
	if err != nil {
		color.Errorln(err)
		return nil
	}

	return instance.(contractsconsole.Artisan)
}

func (r *Container) MakeAuth(ctx ...contractshttp.Context) contractsauth.Auth {
	parameters := map[string]any{}
	if len(ctx) > 0 {
		parameters["ctx"] = ctx[0]
	}

	instance, err := r.MakeWith(facades.FacadeToBinding[facades.Auth], parameters)
	if err != nil {
		color.Errorln(err)
		return nil
	}
	if instance == nil {
		return nil
	}

	return instance.(contractsauth.Auth)
}

func (r *Container) MakeCache() contractscache.Cache {
	instance, err := r.Make(facades.FacadeToBinding[facades.Cache])
	if err != nil {
		color.Errorln(err)
		return nil
	}

	return instance.(contractscache.Cache)
}

func (r *Container) MakeConfig() contractsconfig.Config {
	instance, err := r.Make(facades.FacadeToBinding[facades.Config])
	if err != nil {
		color.Errorln(err)
		return nil
	}

	return instance.(contractsconfig.Config)
}

func (r *Container) MakeCrypt() contractscrypt.Crypt {
	instance, err := r.Make(facades.FacadeToBinding[facades.Crypt])
	if err != nil {
		color.Errorln(err)
		return nil
	}

	return instance.(contractscrypt.Crypt)
}

func (r *Container) MakeDB() contractsdb.DB {
	instance, err := r.Make(facades.FacadeToBinding[facades.DB])
	if err != nil {
		color.Errorln(err)
		return nil
	}
	if instance == nil {
		return nil
	}

	return instance.(contractsdb.DB)
}

func (r *Container) MakeEvent() contractsevent.Instance {
	instance, err := r.Make(facades.FacadeToBinding[facades.Event])
	if err != nil {
		color.Errorln(err)
		return nil
	}

	return instance.(contractsevent.Instance)
}

func (r *Container) MakeGate() contractsaccess.Gate {
	instance, err := r.Make(facades.FacadeToBinding[facades.Gate])
	if err != nil {
		color.Errorln(err)
		return nil
	}

	return instance.(contractsaccess.Gate)
}

func (r *Container) MakeGrpc() contractsgrpc.Grpc {
	instance, err := r.Make(facades.FacadeToBinding[facades.Grpc])
	if err != nil {
		color.Errorln(err)
		return nil
	}

	return instance.(contractsgrpc.Grpc)
}

func (r *Container) MakeHash() contractshash.Hash {
	instance, err := r.Make(facades.FacadeToBinding[facades.Hash])
	if err != nil {
		color.Errorln(err)
		return nil
	}

	return instance.(contractshash.Hash)
}

func (r *Container) MakeHttp() contractshttpclient.Factory {
	instance, err := r.Make(facades.FacadeToBinding[facades.Http])
	if err != nil {
		color.Errorln(err)
		return nil
	}

	return instance.(contractshttpclient.Factory)
}

func (r *Container) MakeLang(ctx context.Context) contractstranslation.Translator {
	instance, err := r.MakeWith(facades.FacadeToBinding[facades.Lang], map[string]any{
		"ctx": ctx,
	})
	if err != nil {
		color.Errorln(err)
		return nil
	}

	return instance.(contractstranslation.Translator)
}

func (r *Container) MakeLog() contractslog.Log {
	instance, err := r.Make(facades.FacadeToBinding[facades.Log])
	if err != nil {
		color.Errorln(err)
		return nil
	}

	return instance.(contractslog.Log)
}

func (r *Container) MakeMail() contractsmail.Mail {
	instance, err := r.Make(facades.FacadeToBinding[facades.Mail])
	if err != nil {
		color.Errorln(err)
		return nil
	}

	return instance.(contractsmail.Mail)
}

func (r *Container) MakeOrm() contractsorm.Orm {
	instance, err := r.Make(facades.FacadeToBinding[facades.Orm])
	if err != nil {
		color.Errorln(err)
		return nil
	}
	if instance == nil {
		return nil
	}

	return instance.(contractsorm.Orm)
}

func (r *Container) MakeProcess() contractsprocess.Process {
	instance, err := r.Make(facades.FacadeToBinding[facades.Process])
	if err != nil {
		color.Errorln(err)
		return nil
	}
	if instance == nil {
		return nil
	}

	return instance.(contractsprocess.Process)
}

func (r *Container) MakeQueue() contractsqueue.Queue {
	instance, err := r.Make(facades.FacadeToBinding[facades.Queue])
	if err != nil {
		color.Errorln(err)
		return nil
	}

	return instance.(contractsqueue.Queue)
}

func (r *Container) MakeRateLimiter() contractshttp.RateLimiter {
	instance, err := r.Make(facades.FacadeToBinding[facades.RateLimiter])
	if err != nil {
		color.Errorln(err)
		return nil
	}

	return instance.(contractshttp.RateLimiter)
}

func (r *Container) MakeRoute() contractsroute.Route {
	instance, err := r.Make(facades.FacadeToBinding[facades.Route])
	if err != nil {
		color.Errorln(err)
		return nil
	}

	return instance.(contractsroute.Route)
}

func (r *Container) MakeSchedule() contractsschedule.Schedule {
	instance, err := r.Make(facades.FacadeToBinding[facades.Schedule])
	if err != nil {
		color.Errorln(err)
		return nil
	}

	return instance.(contractsschedule.Schedule)
}

func (r *Container) MakeSchema() contractsmigration.Schema {
	instance, err := r.Make(facades.FacadeToBinding[facades.Schema])
	if err != nil {
		color.Errorln(err)
		return nil
	}
	if instance == nil {
		return nil
	}

	return instance.(contractsmigration.Schema)
}

func (r *Container) MakeSeeder() contractsseerder.Facade {
	instance, err := r.Make(facades.FacadeToBinding[facades.Seeder])

	if err != nil {
		color.Errorln(err)
		return nil
	}

	return instance.(contractsseerder.Facade)
}

func (r *Container) MakeSession() contractsession.Manager {
	instance, err := r.Make(facades.FacadeToBinding[facades.Session])
	if err != nil {
		color.Errorln(err)
		return nil
	}

	return instance.(contractsession.Manager)
}

func (r *Container) MakeStorage() contractsfilesystem.Storage {
	instance, err := r.Make(facades.FacadeToBinding[facades.Storage])
	if err != nil {
		color.Errorln(err)
		return nil
	}

	return instance.(contractsfilesystem.Storage)
}

func (r *Container) MakeTelemetry() contractstelemetry.Telemetry {
	instance, err := r.Make(facades.FacadeToBinding[facades.Telemetry])
	if err != nil {
		color.Errorln(err)
		return nil
	}

	return instance.(contractstelemetry.Telemetry)
}

func (r *Container) MakeTesting() contractstesting.Testing {
	instance, err := r.Make(facades.FacadeToBinding[facades.Testing])
	if err != nil {
		color.Errorln(err)
		return nil
	}

	return instance.(contractstesting.Testing)
}

func (r *Container) MakeValidation() contractsvalidation.Validation {
	instance, err := r.Make(facades.FacadeToBinding[facades.Validation])
	if err != nil {
		color.Errorln(err)
		return nil
	}

	return instance.(contractsvalidation.Validation)
}

func (r *Container) MakeView() contractsview.View {
	instance, err := r.Make(facades.FacadeToBinding[facades.View])
	if err != nil {
		color.Errorln(err)
		return nil
	}

	return instance.(contractsview.View)
}

func (r *Container) MakeWith(key any, parameters map[string]any) (any, error) {
	return r.make(key, parameters)
}

func (r *Container) Singleton(key any, callback func(app contractsfoundation.Application) (any, error)) {
	r.bindings.Store(key, instance{concrete: callback, shared: true})
}

func (r *Container) make(key any, parameters map[string]any) (any, error) {
	binding, ok := r.bindings.Load(key)
	if !ok {
		return nil, fmt.Errorf("binding not found: %+v", key)
	}

	if parameters == nil {
		instance, ok := r.instances.Load(key)
		if ok {
			return instance, nil
		}
	}

	bindingImpl := binding.(instance)
	switch concrete := bindingImpl.concrete.(type) {
	case func(app contractsfoundation.Application) (any, error):
		concreteImpl, err := concrete(App)
		if err != nil {
			return nil, err
		}
		if bindingImpl.shared {
			r.instances.Store(key, concreteImpl)
		}

		return concreteImpl, nil
	case func(app contractsfoundation.Application, parameters map[string]any) (any, error):
		concreteImpl, err := concrete(App, parameters)
		if err != nil {
			return nil, err
		}

		return concreteImpl, nil
	default:
		r.instances.Store(key, concrete)

		return concrete, nil
	}
}
