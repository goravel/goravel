package http

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"

	"github.com/stretchr/testify/assert"

	"github.com/goravel/framework/contracts/foundation"
	"github.com/goravel/framework/contracts/route"
	"github.com/goravel/framework/contracts/session"
	"github.com/goravel/framework/contracts/testing"
	contractshttp "github.com/goravel/framework/contracts/testing/http"
	"github.com/goravel/framework/errors"
	"github.com/goravel/framework/support/collect"
	"github.com/goravel/framework/support/maps"
	"github.com/goravel/framework/support/str"
)

type TestRequest struct {
	t                 testing.TestingT
	ctx               context.Context
	defaultHeaders    map[string]string
	defaultCookies    []*http.Cookie
	json              foundation.Json
	route             route.Route
	session           session.Manager
	sessionAttributes map[string]any
}

func NewTestRequest(t testing.TestingT, json foundation.Json, route route.Route, session session.Manager) contractshttp.Request {
	return &TestRequest{
		t:                 t,
		ctx:               context.Background(),
		defaultHeaders:    make(map[string]string),
		defaultCookies:    make([]*http.Cookie, 0),
		json:              json,
		route:             route,
		session:           session,
		sessionAttributes: make(map[string]any),
	}
}

func (r *TestRequest) Get(uri string) (contractshttp.Response, error) {
	return r.call(http.MethodGet, uri, nil)
}

func (r *TestRequest) Post(uri string, body io.Reader) (contractshttp.Response, error) {
	if r.defaultHeaders["Content-Type"] == "" {
		r.WithHeader("Content-Type", "application/json")
	}

	return r.call(http.MethodPost, uri, body)
}

func (r *TestRequest) Put(uri string, body io.Reader) (contractshttp.Response, error) {
	if r.defaultHeaders["Content-Type"] == "" {
		r.WithHeader("Content-Type", "application/json")
	}

	return r.call(http.MethodPut, uri, body)
}

func (r *TestRequest) Delete(uri string, body io.Reader) (contractshttp.Response, error) {
	return r.call(http.MethodDelete, uri, body)
}

func (r *TestRequest) Patch(uri string, body io.Reader) (contractshttp.Response, error) {
	return r.call(http.MethodPatch, uri, body)
}

func (r *TestRequest) Head(uri string) (contractshttp.Response, error) {
	return r.call(http.MethodHead, uri, nil)
}

func (r *TestRequest) Options(uri string) (contractshttp.Response, error) {
	return r.call(http.MethodOptions, uri, nil)
}

func (r *TestRequest) FlushHeaders() contractshttp.Request {
	r.defaultHeaders = make(map[string]string)
	return r
}

func (r *TestRequest) WithHeaders(headers map[string]string) contractshttp.Request {
	r.defaultHeaders = collect.Merge(r.defaultHeaders, headers)
	return r
}

func (r *TestRequest) WithHeader(key, value string) contractshttp.Request {
	maps.Set(r.defaultHeaders, key, value)
	return r
}

func (r *TestRequest) WithoutHeader(key string) contractshttp.Request {
	maps.Forget(r.defaultHeaders, key)
	return r
}

func (r *TestRequest) WithCookies(cookies []*http.Cookie) contractshttp.Request {
	r.defaultCookies = append(r.defaultCookies, cookies...)
	return r
}

func (r *TestRequest) WithCookie(cookie *http.Cookie) contractshttp.Request {
	r.defaultCookies = append(r.defaultCookies, cookie)
	return r
}

func (r *TestRequest) WithContext(ctx context.Context) contractshttp.Request {
	r.ctx = ctx
	return r
}

func (r *TestRequest) WithToken(token string, ttype ...string) contractshttp.Request {
	tt := "Bearer"
	if len(ttype) > 0 {
		tt = ttype[0]
	}
	return r.WithHeader("Authorization", fmt.Sprintf("%s %s", tt, token))
}

func (r *TestRequest) WithBasicAuth(username, password string) contractshttp.Request {
	encoded := base64.StdEncoding.EncodeToString(fmt.Appendf(nil, "%s:%s", username, password))
	return r.WithToken(encoded, "Basic")
}

func (r *TestRequest) WithoutToken() contractshttp.Request {
	return r.WithoutHeader("Authorization")
}

func (r *TestRequest) WithSession(attributes map[string]any) contractshttp.Request {
	r.sessionAttributes = collect.Merge(r.sessionAttributes, attributes)
	return r
}

func (r *TestRequest) call(method string, uri string, body io.Reader) (contractshttp.Response, error) {
	err := r.setSession()
	if err != nil {
		return nil, err
	}
	if !str.Of(uri).StartsWith("/", "http://", "https://") {
		uri = "/" + uri
	}

	req := httptest.NewRequest(method, uri, body).WithContext(r.ctx)

	for key, value := range r.defaultHeaders {
		req.Header.Set(key, value)
	}

	for _, cookie := range r.defaultCookies {
		req.AddCookie(cookie)
	}

	if r.route == nil {
		assert.FailNow(r.t, errors.RouteFacadeNotSet.SetModule(errors.ModuleTesting).Error())
		return nil, errors.RouteFacadeNotSet
	}

	response, err := r.route.Test(req)
	if err != nil {
		return nil, err
	}

	return NewTestResponse(r.t, response, r.json, r.session), nil
}

func (r *TestRequest) setSession() error {
	if len(r.sessionAttributes) == 0 {
		return nil
	}

	if r.session == nil {
		return errors.SessionFacadeNotSet
	}

	// Retrieve session driver
	driver, err := r.session.Driver()
	if err != nil {
		return err
	}

	// Build session
	sess, err := r.session.BuildSession(driver)
	if err != nil {
		return err
	}

	for key, value := range r.sessionAttributes {
		sess.Put(key, value)
	}

	r.WithCookie(&http.Cookie{Name: sess.GetName(), Value: sess.GetID()})

	if err = sess.Save(); err != nil {
		return err
	}

	// Release session
	r.session.ReleaseSession(sess)
	return nil
}
