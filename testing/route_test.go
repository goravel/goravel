package testing

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/goravel/framework/facades"
	"github.com/stretchr/testify/assert"
)

func TestRoute(t *testing.T) {
	tests := []struct {
		name       string
		method     string
		url        string
		expectCode int
		expectBody string
	}{
		{
			name:       "Get",
			method:     "GET",
			url:        "/input/1",
			expectCode: http.StatusOK,
			expectBody: "{\"id\":\"1\"}",
		},
		{
			name:       "Post",
			method:     "POST",
			url:        "/input/1",
			expectCode: http.StatusOK,
			expectBody: "{\"id\":\"1\"}",
		},
		{
			name:       "Put",
			method:     "PUT",
			url:        "/input/1",
			expectCode: http.StatusOK,
			expectBody: "{\"id\":\"1\"}",
		},
		{
			name:       "Delete",
			method:     "DELETE",
			url:        "/input/1",
			expectCode: http.StatusOK,
			expectBody: "{\"id\":\"1\"}",
		},
		{
			name:       "Options",
			method:     "OPTIONS",
			url:        "/input/1",
			expectCode: http.StatusOK,
			expectBody: "{\"id\":\"1\"}",
		},
		{
			name:       "Patch",
			method:     "PATCH",
			url:        "/input/1",
			expectCode: http.StatusOK,
			expectBody: "{\"id\":\"1\"}",
		},
		{
			name:       "Any Get",
			method:     "GET",
			url:        "/any/1",
			expectCode: http.StatusOK,
			expectBody: "{\"id\":\"1\"}",
		},
		{
			name:       "Any Post",
			method:     "POST",
			url:        "/any/1",
			expectCode: http.StatusOK,
			expectBody: "{\"id\":\"1\"}",
		},
		{
			name:       "Any Put",
			method:     "PUT",
			url:        "/any/1",
			expectCode: http.StatusOK,
			expectBody: "{\"id\":\"1\"}",
		},
		{
			name:       "Any Delete",
			method:     "DELETE",
			url:        "/any/1",
			expectCode: http.StatusOK,
			expectBody: "{\"id\":\"1\"}",
		},
		{
			name:       "Any Options",
			method:     "OPTIONS",
			url:        "/any/1",
			expectCode: http.StatusOK,
			expectBody: "{\"id\":\"1\"}",
		},
		{
			name:       "Any Patch",
			method:     "PATCH",
			url:        "/any/1",
			expectCode: http.StatusOK,
			expectBody: "{\"id\":\"1\"}",
		},
		{
			name:       "Static",
			method:     "GET",
			url:        "/static/logo.png",
			expectCode: http.StatusOK,
		},
		{
			name:       "StaticFile",
			method:     "GET",
			url:        "/static-file",
			expectCode: http.StatusOK,
		},
		{
			name:       "StaticFS",
			method:     "GET",
			url:        "/static-fs",
			expectCode: http.StatusMovedPermanently,
		},
		{
			name:       "Abort Middleware",
			method:     "GET",
			url:        "/middleware/1",
			expectCode: http.StatusNonAuthoritativeInfo,
		},
		{
			name:       "Multiple Middleware",
			method:     "GET",
			url:        "/middlewares/1",
			expectCode: http.StatusOK,
			expectBody: "{\"ctx\":\"Goravel\",\"ctx1\":\"Hello\",\"id\":\"1\"}",
		},
		{
			name:       "Multiple Prefix",
			method:     "GET",
			url:        "/prefix1/prefix2/input/1",
			expectCode: http.StatusOK,
			expectBody: "{\"id\":\"1\"}",
		},
		{
			name:       "Multiple Prefix Group Middleware",
			method:     "GET",
			url:        "/group1/group2/middleware/1",
			expectCode: http.StatusOK,
			expectBody: "{\"ctx\":\"Goravel\",\"ctx1\":\"Hello\",\"id\":\"1\"}",
		},
		{
			name:       "Multiple Group Middleware",
			method:     "GET",
			url:        "/group1/middleware/1",
			expectCode: http.StatusOK,
			expectBody: "{\"ctx\":\"Goravel\",\"ctx2\":\"World\",\"id\":\"1\"}",
		},

		//{
		//	name:   "Post+Form",
		//	method: "POST",
		//	url:    "/post-form",
		//	setup: func(method, url string) error {
		//		payload := &bytes.Buffer{}
		//		writer := multipart.NewWriter(payload)
		//		if err := writer.WriteField("name", "Goravel"); err != nil {
		//			return err
		//		}
		//		if err := writer.Close(); err != nil {
		//			return err
		//		}
		//
		//		req, _ = http.NewRequest(method, url, payload)
		//		req.Header.Set("Content-Type", writer.FormDataContentType())
		//
		//		return nil
		//	},
		//	expectCode: http.StatusOK,
		//	expectBody: "{\"name\":\"Goravel\"}",
		//},
		//{
		//	name:   "AbortMiddleware",
		//	method: "POST",
		//	url:    "/middleware/1",
		//	setup: func(method, url string) error {
		//		req, _ = http.NewRequest(method, url, nil)
		//
		//		return nil
		//	},
		//	expectCode: http.StatusNonAuthoritativeInfo,
		//	//expectBody: "{\"id\":\"1\"}",
		//},
		//{
		//	name:   "Prefix+Put+Query+String",
		//	method: "PUT",
		//	url:    "/prefix/put?id=2",
		//	setup: func(method, url string) error {
		//		req, _ = http.NewRequest(method, url, nil)
		//
		//		return nil
		//	},
		//	expectCode: http.StatusOK,
		//	expectBody: "2",
		//},
		//{
		//	name:   "Prefix+Prefix+Put+Query+String",
		//	method: "PUT",
		//	url:    "/prefix1/prefix2/put?id=2",
		//	setup: func(method, url string) error {
		//		req, _ = http.NewRequest(method, url, nil)
		//
		//		return nil
		//	},
		//	expectCode: http.StatusOK,
		//	expectBody: "2",
		//},
		//{
		//	name:   "Group+Delete",
		//	method: "DELETE",
		//	url:    "/group/1",
		//	setup: func(method, url string) error {
		//		req, _ = http.NewRequest(method, url, nil)
		//
		//		return nil
		//	},
		//	expectCode: http.StatusOK,
		//	expectBody: "1",
		//},
		//{
		//	name:   "Middleware+Context+Group",
		//	method: "GET",
		//	url:    "/middleware-group/1",
		//	setup: func(method, url string) error {
		//		req, _ = http.NewRequest(method, url, nil)
		//
		//		return nil
		//	},
		//	expectCode: http.StatusOK,
		//	expectBody: "{\"ctx\":\"Goravel\",\"id\":\"1\"}",
		//},
		//{
		//	name:   "Prefix+Middleware+Context+Group",
		//	method: "GET",
		//	url:    "/prefix/middleware-group/1",
		//	setup: func(method, url string) error {
		//		req, _ = http.NewRequest(method, url, nil)
		//
		//		return nil
		//	},
		//	expectCode: http.StatusOK,
		//	expectBody: "{\"ctx\":\"Goravel\",\"id\":\"1\"}",
		//},
		//{
		//	name:   "Prefix+Middleware+Context+Group",
		//	method: "GET",
		//	url:    "/prefix1/prefix2/middleware-group/1",
		//	setup: func(method, url string) error {
		//		req, _ = http.NewRequest(method, url, nil)
		//
		//		return nil
		//	},
		//	expectCode: http.StatusOK,
		//	expectBody: "{\"ctx\":\"Goravel\",\"ctx1\":\"Hello\",\"id\":\"1\"}",
		//},
	}

	for _, test := range tests {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(test.method, test.url, nil)
		facades.Route.ServeHTTP(w, req)

		if test.expectBody != "" {
			assert.Equal(t, test.expectBody, w.Body.String(), test.name)
		}
		assert.Equal(t, test.expectCode, w.Code, test.name)
	}
}
