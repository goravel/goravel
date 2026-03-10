package http

import (
	"fmt"
	"html"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/goravel/framework/contracts/foundation"
	contractssession "github.com/goravel/framework/contracts/session"
	"github.com/goravel/framework/contracts/testing"
	contractshttp "github.com/goravel/framework/contracts/testing/http"
	"github.com/goravel/framework/errors"
	"github.com/goravel/framework/support/carbon"
)

type TestResponseImpl struct {
	t                 testing.TestingT
	json              foundation.Json
	session           contractssession.Manager
	response          *http.Response
	sessionAttributes map[string]any
	content           string
	mu                sync.Mutex
}

func NewTestResponse(t testing.TestingT, response *http.Response, json foundation.Json, session contractssession.Manager) contractshttp.Response {
	return &TestResponseImpl{
		t:        t,
		response: response,
		json:     json,
		session:  session,
	}
}

func (r *TestResponseImpl) Bind(value any) error {
	content, err := r.getContent()
	if err != nil {
		return err
	}

	return r.json.UnmarshalString(content, value)
}

func (r *TestResponseImpl) Json() (map[string]any, error) {
	content, err := r.getContent()
	if err != nil {
		return nil, err
	}

	assertable, err := NewAssertableJSON(r.t, r.json, content)
	if err != nil {
		return nil, err
	}

	return assertable.Json(), nil
}

func (r *TestResponseImpl) Headers() http.Header {
	return r.response.Header
}

func (r *TestResponseImpl) Cookies() []*http.Cookie {
	return r.response.Cookies()
}

func (r *TestResponseImpl) Cookie(name string) *http.Cookie {
	return r.getCookie(name)
}

func (r *TestResponseImpl) Session() (map[string]any, error) {
	if r.sessionAttributes != nil {
		return r.sessionAttributes, nil
	}

	if r.session == nil {
		return nil, errors.SessionFacadeNotSet
	}

	// Retrieve session driver
	driver, err := r.session.Driver()
	if err != nil {
		return nil, err
	}

	// Build session
	session, err := r.session.BuildSession(driver)
	if err != nil {
		return nil, err
	}

	r.sessionAttributes = session.All()
	r.session.ReleaseSession(session)

	return r.sessionAttributes, nil
}

func (r *TestResponseImpl) IsSuccessful() bool {
	statusCode := r.getStatusCode()
	return statusCode >= 200 && statusCode < 300
}

func (r *TestResponseImpl) IsServerError() bool {
	statusCode := r.getStatusCode()
	return statusCode >= 500 && statusCode < 600
}

func (r *TestResponseImpl) Content() (string, error) {
	return r.getContent()
}

func (r *TestResponseImpl) AssertStatus(status int) contractshttp.Response {
	actual := r.getStatusCode()
	assert.Equal(r.t, status, actual, fmt.Sprintf("Expected response status code [%d] but received %d.", status, actual))
	return r
}

func (r *TestResponseImpl) AssertOk() contractshttp.Response {
	return r.AssertStatus(http.StatusOK)
}

func (r *TestResponseImpl) AssertCreated() contractshttp.Response {
	return r.AssertStatus(http.StatusCreated)
}

func (r *TestResponseImpl) AssertAccepted() contractshttp.Response {
	return r.AssertStatus(http.StatusAccepted)
}

func (r *TestResponseImpl) AssertNoContent(status ...int) contractshttp.Response {
	expectedStatus := http.StatusNoContent
	if len(status) > 0 {
		expectedStatus = status[0]
	}

	r.AssertStatus(expectedStatus)

	content, err := r.getContent()
	assert.Nil(r.t, err)
	assert.Empty(r.t, content)

	return r
}

func (r *TestResponseImpl) AssertMovedPermanently() contractshttp.Response {
	return r.AssertStatus(http.StatusMovedPermanently)
}

func (r *TestResponseImpl) AssertFound() contractshttp.Response {
	return r.AssertStatus(http.StatusFound)
}

func (r *TestResponseImpl) AssertNotModified() contractshttp.Response {
	return r.AssertStatus(http.StatusNotModified)
}

func (r *TestResponseImpl) AssertPartialContent() contractshttp.Response {
	return r.AssertStatus(http.StatusPartialContent)
}

func (r *TestResponseImpl) AssertTemporaryRedirect() contractshttp.Response {
	return r.AssertStatus(http.StatusTemporaryRedirect)
}

func (r *TestResponseImpl) AssertBadRequest() contractshttp.Response {
	return r.AssertStatus(http.StatusBadRequest)
}

func (r *TestResponseImpl) AssertUnauthorized() contractshttp.Response {
	return r.AssertStatus(http.StatusUnauthorized)
}

func (r *TestResponseImpl) AssertPaymentRequired() contractshttp.Response {
	return r.AssertStatus(http.StatusPaymentRequired)
}

func (r *TestResponseImpl) AssertForbidden() contractshttp.Response {
	return r.AssertStatus(http.StatusForbidden)
}

func (r *TestResponseImpl) AssertNotFound() contractshttp.Response {
	return r.AssertStatus(http.StatusNotFound)
}

func (r *TestResponseImpl) AssertMethodNotAllowed() contractshttp.Response {
	return r.AssertStatus(http.StatusMethodNotAllowed)
}

func (r *TestResponseImpl) AssertNotAcceptable() contractshttp.Response {
	return r.AssertStatus(http.StatusNotAcceptable)
}

func (r *TestResponseImpl) AssertConflict() contractshttp.Response {
	return r.AssertStatus(http.StatusConflict)
}

func (r *TestResponseImpl) AssertRequestTimeout() contractshttp.Response {
	return r.AssertStatus(http.StatusRequestTimeout)
}

func (r *TestResponseImpl) AssertGone() contractshttp.Response {
	return r.AssertStatus(http.StatusGone)
}

func (r *TestResponseImpl) AssertUnsupportedMediaType() contractshttp.Response {
	return r.AssertStatus(http.StatusUnsupportedMediaType)
}

func (r *TestResponseImpl) AssertUnprocessableEntity() contractshttp.Response {
	return r.AssertStatus(http.StatusUnprocessableEntity)
}

func (r *TestResponseImpl) AssertTooManyRequests() contractshttp.Response {
	return r.AssertStatus(http.StatusTooManyRequests)
}

func (r *TestResponseImpl) AssertInternalServerError() contractshttp.Response {
	return r.AssertStatus(http.StatusInternalServerError)
}

func (r *TestResponseImpl) AssertServiceUnavailable() contractshttp.Response {
	return r.AssertStatus(http.StatusServiceUnavailable)
}

func (r *TestResponseImpl) AssertHeader(headerName, value string) contractshttp.Response {
	got := r.getHeader(headerName)
	assert.NotEmpty(r.t, got, fmt.Sprintf("Header [%s] not present on response.", headerName))
	if got != "" {
		assert.Equal(r.t, value, got, fmt.Sprintf("Header [%s] was found, but value [%s] does not match [%s].", headerName, got, value))
	}
	return r
}

func (r *TestResponseImpl) AssertHeaderMissing(headerName string) contractshttp.Response {
	got := r.getHeader(headerName)
	assert.Empty(r.t, got, fmt.Sprintf("Unexpected header [%s] is present on response.", headerName))
	return r
}

func (r *TestResponseImpl) AssertCookie(name, value string) contractshttp.Response {
	cookie := r.getCookie(name)
	assert.NotNil(r.t, cookie, fmt.Sprintf("Cookie [%s] not present on response.", name))

	if cookie == nil {
		return r
	}

	assert.Equal(r.t, value, cookie.Value, fmt.Sprintf("Cookie [%s] was found, but value [%s] does not match [%s]", name, cookie.Value, value))

	return r
}

func (r *TestResponseImpl) AssertCookieExpired(name string) contractshttp.Response {
	cookie := r.getCookie(name)
	assert.NotNil(r.t, cookie, fmt.Sprintf("Cookie [%s] not present on response.", name))

	if cookie == nil {
		return r
	}

	expirationTime := carbon.FromStdTime(cookie.Expires)
	assert.True(r.t, r.isCookieExpired(cookie), fmt.Sprintf("Cookie [%s] is not expired; it expires at [%s].", name, expirationTime.ToString()))

	return r
}

func (r *TestResponseImpl) AssertCookieNotExpired(name string) contractshttp.Response {
	cookie := r.getCookie(name)
	assert.NotNil(r.t, cookie, fmt.Sprintf("Cookie [%s] not present on response.", name))

	if cookie == nil {
		return r
	}

	expirationTime := carbon.FromStdTime(cookie.Expires)
	assert.True(r.t, !r.isCookieExpired(cookie), fmt.Sprintf("Cookie [%s] is expired; it expired at [%s].", name, expirationTime))
	return r
}

func (r *TestResponseImpl) AssertCookieMissing(name string) contractshttp.Response {
	assert.Nil(r.t, r.getCookie(name), fmt.Sprintf("Cookie [%s] is present on response.", name))

	return r
}

func (r *TestResponseImpl) AssertSuccessful() contractshttp.Response {
	assert.True(r.t, r.IsSuccessful(), fmt.Sprintf("Expected response status code >=200, <300 but received %d.", r.getStatusCode()))

	return r
}

func (r *TestResponseImpl) AssertServerError() contractshttp.Response {
	assert.True(r.t, r.IsServerError(), fmt.Sprintf("Expected response status code >=500, <600 but received %d.", r.getStatusCode()))

	return r
}

func (r *TestResponseImpl) AssertDontSee(value []string, escaped ...bool) contractshttp.Response {
	content, err := r.getContent()
	assert.Nil(r.t, err)

	shouldEscape := true
	if len(escaped) > 0 {
		shouldEscape = escaped[0]
	}

	for _, v := range value {
		checkValue := v
		if shouldEscape {
			checkValue = html.EscapeString(v)
		}

		assert.NotContains(r.t, content, checkValue, fmt.Sprintf("Response should not contain '%s', but it was found.", checkValue))
	}

	return r
}

func (r *TestResponseImpl) AssertSee(value []string, escaped ...bool) contractshttp.Response {
	content, err := r.getContent()
	assert.Nil(r.t, err)

	shouldEscape := true
	if len(escaped) > 0 {
		shouldEscape = escaped[0]
	}

	for _, v := range value {
		checkValue := v
		if shouldEscape {
			checkValue = html.EscapeString(v)
		}

		assert.Contains(r.t, content, checkValue, fmt.Sprintf("Expected to see '%s' in response, but it was not found.", checkValue))
	}

	return r
}

func (r *TestResponseImpl) AssertSeeInOrder(value []string, escaped ...bool) contractshttp.Response {
	content, err := r.getContent()
	assert.Nil(r.t, err)

	shouldEscape := true
	if len(escaped) > 0 {
		shouldEscape = escaped[0]
	}

	previousIndex := -1
	for _, v := range value {
		checkValue := v
		if shouldEscape {
			checkValue = html.EscapeString(v)
		}

		currentIndex := strings.Index(content[previousIndex+1:], checkValue)
		assert.GreaterOrEqual(r.t, currentIndex, 0, fmt.Sprintf("Expected to see '%s' in response in the correct order, but it was not found.", checkValue))
		previousIndex += currentIndex + len(checkValue)
	}

	return r
}

func (r *TestResponseImpl) AssertJson(data map[string]any) contractshttp.Response {
	content, err := r.getContent()
	assert.Nil(r.t, err)

	assertableJson, err := NewAssertableJSON(r.t, r.json, content)
	assert.Nil(r.t, err)

	for key, value := range data {
		assertableJson.Where(key, value)
	}

	return r
}

func (r *TestResponseImpl) AssertExactJson(data map[string]any) contractshttp.Response {
	actual, err := r.Json()
	assert.Nil(r.t, err)
	assert.Equal(r.t, data, actual, "The JSON response does not match exactly with the expected content")
	return r
}

func (r *TestResponseImpl) AssertJsonMissing(data map[string]any) contractshttp.Response {
	actual, err := r.Json()
	assert.Nil(r.t, err)

	for key, expectedValue := range data {
		actualValue, found := actual[key]
		if found {
			assert.NotEqual(r.t, expectedValue, actualValue, "Found unexpected key-value pair in JSON response: key '%s' with value '%v'", key, actualValue)
		}
	}
	return r
}

func (r *TestResponseImpl) AssertFluentJson(callback func(json contractshttp.AssertableJSON)) contractshttp.Response {
	content, err := r.getContent()
	assert.Nil(r.t, err)

	assertableJson, err := NewAssertableJSON(r.t, r.json, content)
	assert.Nil(r.t, err)

	callback(assertableJson)

	return r
}

func (r *TestResponseImpl) getStatusCode() int {
	return r.response.StatusCode
}

func (r *TestResponseImpl) getContent() (string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.content != "" {
		return r.content, nil
	}

	defer errors.Ignore(r.response.Body.Close)

	content, err := io.ReadAll(r.response.Body)
	if err != nil {
		return "", err
	}

	r.content = string(content)
	return r.content, nil
}

func (r *TestResponseImpl) getCookie(name string) *http.Cookie {
	for _, c := range r.response.Cookies() {
		if c.Name == name {
			return c
		}
	}

	return nil
}

func (r *TestResponseImpl) getHeader(name string) string {
	return r.response.Header.Get(name)
}

func (r *TestResponseImpl) isCookieExpired(cookie *http.Cookie) bool {
	if cookie.MaxAge > 0 {
		return false
	}

	if cookie.MaxAge < 0 {
		return true
	}

	// MaxAge == 0 means no Max-Age specified; check Expires attribute
	if cookie.Expires.IsZero() {
		// Session cookie; consider not expired until the session ends
		return false
	}

	return cookie.Expires.Before(time.Now())
}
