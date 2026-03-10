package auth

import (
	"fmt"

	"github.com/spf13/cast"

	contractsauth "github.com/goravel/framework/contracts/auth"
	"github.com/goravel/framework/contracts/http"
	contractsession "github.com/goravel/framework/contracts/session"
	"github.com/goravel/framework/errors"
)

type SessionGuard struct {
	session  contractsession.Session
	ctx      http.Context
	provider contractsauth.UserProvider
	guard    string
}

func NewSessionGuard(ctx http.Context, name string, userProvider contractsauth.UserProvider) (contractsauth.GuardDriver, error) {
	if ctx == nil {
		return nil, errors.InvalidHttpContext.SetModule(errors.ModuleAuth)
	}
	session := ctx.Request().Session()
	if session == nil {
		return nil, errors.SessionDriverIsNotSet.SetModule(errors.ModuleAuth)
	}

	return &SessionGuard{
		session:  session,
		ctx:      ctx,
		guard:    name,
		provider: userProvider,
	}, nil
}

func (r *SessionGuard) Check() bool {
	_, err := r.ID()

	return err == nil
}

func (r *SessionGuard) Guest() bool {
	return !r.Check()
}

func (r *SessionGuard) ID() (token string, err error) {
	sessionName := r.getSessionName()
	userID := r.session.Get(sessionName, nil)

	if userID == nil {
		return "", errors.AuthInvalidKey
	}

	if id, ok := userID.(string); ok {
		return id, nil
	}

	return "", errors.AuthInvalidKey
}

func (r *SessionGuard) Login(user any) (token string, err error) {
	id, err := r.provider.GetID(user)
	if err != nil {
		return "", err
	}

	return r.LoginUsingID(id)
}

func (r *SessionGuard) LoginUsingID(id any) (token string, err error) {
	sessionName := r.getSessionName()
	key := cast.ToString(id)

	if key == "" {
		return "", errors.AuthInvalidKey
	}

	r.session.Put(sessionName, key)

	return "", nil
}

func (r *SessionGuard) Logout() error {
	sessionName := r.getSessionName()
	r.session.Forget(sessionName)

	return nil
}

func (r *SessionGuard) Parse(token string) (*contractsauth.Payload, error) {
	return nil, errors.AuthUnsupportedDriverMethod.Args("session")
}

func (r *SessionGuard) Refresh() (token string, err error) {
	return "", errors.AuthUnsupportedDriverMethod.Args("session")
}

func (r *SessionGuard) User(user any) error {
	id, err := r.ID()

	if err != nil {
		return err
	}

	return r.provider.RetriveByID(user, id)
}

func (r *SessionGuard) getSessionName() string {
	return fmt.Sprintf("auth_%s_id", r.guard)
}
