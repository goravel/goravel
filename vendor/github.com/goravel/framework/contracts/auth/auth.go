package auth

import (
	"time"

	"github.com/goravel/framework/contracts/http"
)

type Auth interface {
	GuardDriver
	Guard(name string) GuardDriver
	Extend(name string, fn GuardFunc)
	Provider(name string, fn UserProviderFunc)
}

type GuardDriver interface {
	// Check whether user logged in or not
	Check() bool
	// Check whether user *not* logged in or not | !Check()
	Guest() bool
	// ID returns the current user id.
	ID() (token string, err error)
	// Login logs a user into the application.
	Login(user any) (token string, err error)
	// LoginUsingID logs the given user ID into the application.
	LoginUsingID(id any) (token string, err error)
	// Logout logs the user out of the application.
	Logout() error
	// Parse the given token.
	Parse(token string) (*Payload, error)
	// Refresh the token for the current user.
	Refresh() (token string, err error)
	// User returns the current authenticated user.
	User(user any) error
}

type UserProvider interface {
	// GetID returns the user id.
	GetID(user any) (any, error)
	// RetriveByID returns the user by id, user should be a pointer.
	RetriveByID(user any, id any) error
}

type Payload struct {
	ExpireAt time.Time
	IssuedAt time.Time
	Guard    string
	Key      string
}

type GuardFunc func(ctx http.Context, name string, userProvider UserProvider) (GuardDriver, error)

type UserProviderFunc func(ctx http.Context) (UserProvider, error)
