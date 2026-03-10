package auth

import (
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/cast"

	contractsauth "github.com/goravel/framework/contracts/auth"
	"github.com/goravel/framework/contracts/cache"
	"github.com/goravel/framework/contracts/config"
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/errors"
	"github.com/goravel/framework/support/carbon"
	"github.com/goravel/framework/support/convert"
)

var _ contractsauth.GuardFunc = NewJwtGuard

type Claims struct {
	Key string `json:"key"`
	jwt.RegisteredClaims
}

const ctxJwtKey = "GoravelAuthJwt"

type Guards map[string]*JwtToken

type JwtToken struct {
	Claims *Claims
	Token  string
}

type JwtGuard struct {
	cache    cache.Cache
	config   config.Config
	ctx      http.Context
	guard    string
	provider contractsauth.UserProvider

	secret     string
	ttl        int
	refreshTtl int
}

func NewJwtGuard(ctx http.Context, name string, userProvider contractsauth.UserProvider) (contractsauth.GuardDriver, error) {
	if ctx == nil {
		return nil, errors.InvalidHttpContext.SetModule(errors.ModuleAuth)
	}
	if cacheFacade == nil {
		return nil, errors.CacheFacadeNotSet.SetModule(errors.ModuleAuth)
	}

	jwtSecret := configFacade.GetString(fmt.Sprintf("auth.guards.%s.secret", name))

	if jwtSecret == "" {
		// Get the secret from the jwt config if the guard specific was not set
		jwtSecret = configFacade.GetString("jwt.secret")
	}

	if jwtSecret == "" {
		return nil, errors.AuthEmptySecret
	}

	ttl := configFacade.GetInt(fmt.Sprintf("auth.guards.%s.ttl", name))
	if ttl == 0 {
		ttl = configFacade.GetInt("jwt.ttl")
	}

	if ttl == 0 {
		// 100 years
		ttl = 60 * 24 * 365 * 100
	}

	refreshTtl := configFacade.GetInt(fmt.Sprintf("auth.guards.%s.refresh_ttl", name))

	if refreshTtl == 0 {
		// Get the ttl from the jwt config if the guard specific was not set
		refreshTtl = configFacade.GetInt("jwt.refresh_ttl")
	}

	if refreshTtl == 0 {
		// 100 years
		refreshTtl = 60 * 24 * 365 * 100
	}

	return &JwtGuard{
		cache:    cacheFacade,
		config:   configFacade,
		ctx:      ctx,
		guard:    name,
		provider: userProvider,

		secret:     jwtSecret,
		ttl:        ttl,
		refreshTtl: refreshTtl,
	}, nil
}

func (r *JwtGuard) Check() bool {
	if _, err := r.ID(); err != nil {
		return false
	}

	return true
}

func (r *JwtGuard) GetJwtToken() (*JwtToken, error) {
	guards, ok := r.ctx.Value(ctxJwtKey).(Guards)
	if !ok {
		return nil, ErrorParseTokenFirst
	}

	return r.jwtToken(guards)
}

func (r *JwtGuard) Guest() bool {
	return !r.Check()
}

func (r *JwtGuard) ID() (string, error) {
	guard, err := r.GetJwtToken()
	if err != nil {
		return "", err
	}

	return guard.Claims.Key, nil
}

func (r *JwtGuard) Login(user any) (token string, err error) {
	id, err := r.provider.GetID(user)
	if err != nil {
		return "", err
	}
	if id == nil {
		return "", errors.AuthNoPrimaryKeyField
	}

	return r.LoginUsingID(id)
}

func (r *JwtGuard) LoginUsingID(id any) (token string, err error) {
	nowTime := carbon.Now()
	expireTime := nowTime.Copy().AddMinutes(r.ttl).StdTime()
	key := cast.ToString(id)
	if key == "" {
		return "", errors.AuthInvalidKey
	}
	claims := Claims{
		key,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expireTime),
			IssuedAt:  jwt.NewNumericDate(nowTime.StdTime()),
			Subject:   r.guard,
		},
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err = tokenClaims.SignedString(convert.UnsafeBytes(r.secret))
	if err != nil {
		return "", err
	}

	r.makeAuthContext(&claims, token)

	return
}

func (r *JwtGuard) Logout() error {
	guards, ok := r.ctx.Value(ctxJwtKey).(Guards)
	if !ok {
		return errors.AuthParseTokenFirst
	}

	guard, err := r.jwtToken(guards)
	if err != nil {
		return err
	}

	if err := r.cache.Put(getDisabledCacheKey(guard.Token),
		true,
		time.Duration(r.ttl)*time.Minute,
	); err != nil {
		return err
	}

	delete(guards, r.guard)
	r.ctx.WithValue(ctxJwtKey, guards)

	return nil
}

func (r *JwtGuard) Parse(token string) (*contractsauth.Payload, error) {
	token = strings.ReplaceAll(token, "Bearer ", "")
	if r.tokenIsDisabled(token) {
		return nil, errors.AuthTokenDisabled
	}

	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (any, error) {
		return convert.UnsafeBytes(r.secret), nil
	}, jwt.WithTimeFunc(func() time.Time {
		return carbon.Now().StdTime()
	}))
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) && tokenClaims != nil {
			claims, ok := tokenClaims.Claims.(*Claims)
			if !ok {
				return nil, errors.AuthInvalidClaims
			}

			r.makeAuthContext(claims, "")

			return &contractsauth.Payload{
				Guard:    claims.Subject,
				Key:      claims.Key,
				ExpireAt: claims.ExpiresAt.Local(),
				IssuedAt: claims.IssuedAt.Local(),
			}, errors.AuthTokenExpired
		}

		return nil, errors.AuthInvalidToken
	}
	if tokenClaims == nil || !tokenClaims.Valid {
		return nil, errors.AuthInvalidToken
	}

	claims, ok := tokenClaims.Claims.(*Claims)
	if !ok {
		return nil, errors.AuthInvalidClaims
	}

	r.makeAuthContext(claims, token)

	return &contractsauth.Payload{
		Guard:    claims.Subject,
		Key:      claims.Key,
		ExpireAt: claims.ExpiresAt.Time,
		IssuedAt: claims.IssuedAt.Time,
	}, nil
}

// Refresh need parse token first.
func (r *JwtGuard) Refresh() (token string, err error) {
	guards, ok := r.ctx.Value(ctxJwtKey).(Guards)
	if !ok || guards[r.guard] == nil {
		return "", errors.AuthParseTokenFirst
	}
	if guards[r.guard].Claims == nil {
		return "", errors.AuthParseTokenFirst
	}

	expireTime := carbon.FromStdTime(guards[r.guard].Claims.ExpiresAt.Time).AddMinutes(r.refreshTtl)
	if carbon.Now().Gt(expireTime) {
		return "", errors.AuthRefreshTimeExceeded
	}

	return r.LoginUsingID(guards[r.guard].Claims.Key)
}

// User need parse token first.
func (r *JwtGuard) User(user any) error {
	guard, err := r.GetJwtToken()

	if err != nil {
		return err
	}

	err = r.provider.RetriveByID(user, guard.Claims.Key)

	return err
}

func (r *JwtGuard) jwtToken(guards Guards) (*JwtToken, error) {
	jwtToken, ok := guards[r.guard]
	if !ok || jwtToken == nil {
		return nil, ErrorParseTokenFirst
	}

	if jwtToken.Claims == nil {
		return nil, errors.AuthParseTokenFirst
	}

	if jwtToken.Claims.Key == "" {
		return nil, errors.AuthInvalidKey
	}

	if jwtToken.Token == "" {
		return nil, errors.AuthTokenExpired
	}

	return jwtToken, nil
}

func (r *JwtGuard) makeAuthContext(claims *Claims, token string) {
	guards, ok := r.ctx.Value(ctxJwtKey).(Guards)
	if !ok {
		guards = make(Guards)
	}
	if guard, ok := guards[r.guard]; ok {
		guard.Claims = claims
		guard.Token = token
	} else {
		guards[r.guard] = &JwtToken{claims, token}
	}
	r.ctx.WithValue(ctxJwtKey, guards)
}

func (r *JwtGuard) tokenIsDisabled(token string) bool {
	return r.cache.GetBool(getDisabledCacheKey(token), false)
}

func getDisabledCacheKey(token string) string {
	return "jwt:disabled:" + token
}
