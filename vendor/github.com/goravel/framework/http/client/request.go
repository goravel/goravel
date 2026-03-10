package client

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"maps"
	"net/http"
	"net/url"
	"strings"

	"github.com/goravel/framework/contracts/foundation"
	"github.com/goravel/framework/contracts/http/client"
	supportmaps "github.com/goravel/framework/support/maps"
)

var _ client.Request = (*Request)(nil)

type Request struct {
	client *http.Client
	json   foundation.Json

	baseUrl     string
	clientName  string
	ctx         context.Context
	headers     http.Header
	queryParams url.Values
	urlParams   map[string]string
	cookies     []*http.Cookie

	// clientErr stores any error that occurred during the creation of the parent Client.
	//
	// This allows the Factory to return a "zombie" Client when a configuration is missing,
	// preserving the fluent API chain (e.g., Http.Client("missing").Get("/")).
	// The error is checked and returned lazily when the request is executed in send().
	clientErr error

	// Inspection fields
	// These fields are populated only during "Hydration" within FakeTransport.
	// They are used solely for inspection/assertion purposes in tests.
	payloadBody []byte
	method      string
	fullUrl     string
}

func NewRequest(client *http.Client, json foundation.Json, baseUrl string, clientName string) *Request {
	return &Request{
		client:     client,
		json:       json,
		baseUrl:    baseUrl,
		clientName: clientName,

		ctx:         context.Background(),
		headers:     make(http.Header),
		cookies:     make([]*http.Cookie, 0),
		queryParams: make(url.Values),
		urlParams:   make(map[string]string),
	}
}

// newRequestWithError creates a request instance that will fail immediately when executed.
//
// This is used internally by the Factory when a requested client configuration (e.g., "github")
// is not found. Instead of panicking, we return this "zombie" request which allows
// method chaining to continue, returning the error lazily only when the request is sent.
func newRequestWithError(err error) *Request {
	return &Request{
		clientErr:   err,
		ctx:         context.Background(),
		headers:     make(http.Header),
		cookies:     make([]*http.Cookie, 0),
		queryParams: make(url.Values),
		urlParams:   make(map[string]string),
	}
}

func (r *Request) HttpClient() *http.Client {
	return r.client
}

func (r *Request) Clone() client.Request {
	return r.clone()
}

func (r *Request) Get(uri string) (client.Response, error) {
	return r.send(http.MethodGet, uri, nil)
}

func (r *Request) Post(uri string, body io.Reader) (client.Response, error) {
	return r.send(http.MethodPost, uri, body)
}

func (r *Request) Put(uri string, body io.Reader) (client.Response, error) {
	return r.send(http.MethodPut, uri, body)
}

func (r *Request) Delete(uri string, body io.Reader) (client.Response, error) {
	return r.send(http.MethodDelete, uri, body)
}

func (r *Request) Patch(uri string, body io.Reader) (client.Response, error) {
	return r.send(http.MethodPatch, uri, body)
}

func (r *Request) Head(uri string) (client.Response, error) {
	return r.send(http.MethodHead, uri, nil)
}

func (r *Request) Options(uri string) (client.Response, error) {
	return r.send(http.MethodOptions, uri, nil)
}

func (r *Request) Accept(contentType string) client.Request {
	return r.WithHeader("Accept", contentType)
}

func (r *Request) AcceptJSON() client.Request {
	return r.Accept("application/json")
}

func (r *Request) AsForm() client.Request {
	return r.WithHeader("Content-Type", "application/x-www-form-urlencoded")
}

func (r *Request) BaseUrl(url string) client.Request {
	n := r.clone()
	n.baseUrl = url
	return n
}

func (r *Request) ClientName() string {
	return r.clientName
}

func (r *Request) FlushHeaders() client.Request {
	n := r.clone()
	n.headers = make(http.Header)
	return n
}

func (r *Request) ReplaceHeaders(headers map[string]string) client.Request {
	return r.WithHeaders(headers)
}

func (r *Request) WithBasicAuth(username, password string) client.Request {
	encoded := base64.StdEncoding.EncodeToString(fmt.Appendf(nil, "%s:%s", username, password))
	return r.WithToken(encoded, "Basic")
}

func (r *Request) WithContext(ctx context.Context) client.Request {
	n := r.clone()
	n.ctx = ctx
	return n
}

func (r *Request) WithCookies(cookies []*http.Cookie) client.Request {
	n := r.clone()
	n.cookies = append(n.cookies, cookies...)
	return n
}

func (r *Request) WithCookie(cookie *http.Cookie) client.Request {
	n := r.clone()
	n.cookies = append(n.cookies, cookie)
	return n
}

func (r *Request) WithHeader(key, value string) client.Request {
	n := r.clone()
	n.headers.Set(key, value)
	return n
}

func (r *Request) WithHeaders(headers map[string]string) client.Request {
	n := r.clone()
	for k, v := range headers {
		n.headers.Set(k, v)
	}
	return n
}

func (r *Request) WithQueryParameter(key, value string) client.Request {
	n := r.clone()
	n.queryParams.Set(key, value)
	return n
}

func (r *Request) WithQueryParameters(params map[string]string) client.Request {
	n := r.clone()
	for k, v := range params {
		n.queryParams.Set(k, v)
	}
	return n
}

func (r *Request) WithQueryString(query string) client.Request {
	params, err := url.ParseQuery(strings.TrimSpace(query))
	if err != nil {
		return r.clone()
	}

	n := r.clone()
	for k, v := range params {
		for _, vv := range v {
			n.queryParams.Add(k, vv)
		}
	}
	return n
}

func (r *Request) WithoutHeader(key string) client.Request {
	n := r.clone()
	n.headers.Del(key)
	return n
}

func (r *Request) WithToken(token string, ttype ...string) client.Request {
	tt := "Bearer"
	if len(ttype) > 0 {
		tt = ttype[0]
	}
	return r.WithHeader("Authorization", fmt.Sprintf("%s %s", tt, token))
}

func (r *Request) WithoutToken() client.Request {
	return r.WithoutHeader("Authorization")
}

func (r *Request) WithUrlParameter(key, value string) client.Request {
	n := r.clone()
	supportmaps.Set(n.urlParams, key, url.PathEscape(value))
	return n
}

func (r *Request) WithUrlParameters(params map[string]string) client.Request {
	n := r.clone()
	for k, v := range params {
		supportmaps.Set(n.urlParams, k, url.PathEscape(v))
	}
	return n
}

func (r *Request) Method() string {
	return r.method
}

func (r *Request) Url() string {
	return r.fullUrl
}

func (r *Request) Body() string {
	if len(r.payloadBody) > 0 {
		return string(r.payloadBody)
	}
	return ""
}

func (r *Request) Header(key string) string {
	return r.headers.Get(key)
}

func (r *Request) Headers() http.Header {
	return r.headers
}

func (r *Request) Input(key string) any {
	if len(r.payloadBody) > 0 {
		var data map[string]any
		if err := r.json.Unmarshal(r.payloadBody, &data); err == nil {
			if val, ok := data[key]; ok {
				return val
			}
		}
	}

	if r.queryParams.Has(key) {
		return r.queryParams.Get(key)
	}

	return nil
}

func (r *Request) clone() *Request {
	n := *r
	n.headers = r.headers.Clone()

	if len(r.cookies) > 0 {
		n.cookies = make([]*http.Cookie, len(r.cookies))
		copy(n.cookies, r.cookies)
	}

	if len(r.queryParams) > 0 {
		n.queryParams = make(url.Values, len(r.queryParams))
		for k, v := range r.queryParams {
			dst := make([]string, len(v))
			copy(dst, v)
			n.queryParams[k] = dst
		}
	} else {
		n.queryParams = make(url.Values)
	}

	if len(r.urlParams) > 0 {
		n.urlParams = make(map[string]string, len(r.urlParams))
		maps.Copy(n.urlParams, r.urlParams)
	} else {
		n.urlParams = make(map[string]string)
	}

	n.payloadBody = nil
	n.method = ""
	n.fullUrl = ""

	return &n
}

func (r *Request) parseRequestURL(uri string) (string, error) {
	baseURL := r.baseUrl

	if !strings.HasPrefix(uri, "http://") && !strings.HasPrefix(uri, "https://") {
		uri = strings.TrimSuffix(baseURL, "/") + "/" + strings.TrimPrefix(uri, "/")
	}

	var buf strings.Builder
	buf.Grow(len(uri) + 10)

	n := len(uri)
	i := 0
	for i < n {
		if uri[i] == '{' {
			j := i + 1
			for j < n && uri[j] != '}' {
				j++
			}

			if j == n {
				buf.WriteString(uri[i:])
				break
			}

			key := uri[i+1 : j]
			if value, found := r.urlParams[key]; found {
				buf.WriteString(value)
			} else {
				buf.WriteString(uri[i : j+1])
			}

			i = j + 1
		} else {
			start := i
			for i < n && uri[i] != '{' {
				i++
			}
			buf.WriteString(uri[start:i])
		}
	}

	reqURL, err := url.Parse(buf.String())
	if err != nil {
		return "", err
	}

	if len(r.queryParams) > 0 {
		if len(strings.TrimSpace(reqURL.RawQuery)) == 0 {
			reqURL.RawQuery = r.queryParams.Encode()
		} else {
			reqURL.RawQuery = reqURL.RawQuery + "&" + r.queryParams.Encode()
		}
	}

	return reqURL.String(), nil
}

func (r *Request) send(method, uri string, body io.Reader) (client.Response, error) {
	if r.clientErr != nil {
		return nil, r.clientErr
	}

	parsedURL, err := r.parseRequestURL(uri)
	if err != nil {
		return nil, err
	}

	ctx := context.WithValue(r.ctx, clientNameKey, r.clientName)
	req, err := http.NewRequestWithContext(ctx, method, parsedURL, body)
	if err != nil {
		return nil, err
	}

	req.Header = r.headers

	for _, value := range r.cookies {
		req.AddCookie(value)
	}

	res, err := r.client.Do(req)
	if err != nil {
		return nil, err
	}

	return NewResponse(res, r.json), nil
}
