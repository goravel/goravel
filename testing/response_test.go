package testing

import (
	"github.com/goravel/framework/support/file"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/goravel/framework/facades"
	"github.com/stretchr/testify/assert"
)

func TestResponse(t *testing.T) {
	var (
		req *http.Request
	)

	beforeEach := func() {

	}

	tests := []struct {
		name         string
		method       string
		url          string
		setup        func(method, url string) error
		expectCode   int
		expectBody   string
		expectHeader string
	}{
		{
			name:   "Json",
			method: "GET",
			url:    "/response/json",
			setup: func(method, url string) error {
				var err error
				req, err = http.NewRequest(method, url, nil)
				if err != nil {
					return err
				}

				return nil
			},
			expectCode: http.StatusOK,
			expectBody: "{\"id\":\"1\"}",
		},
		{
			name:   "String",
			method: "GET",
			url:    "/response/string",
			setup: func(method, url string) error {
				var err error
				req, err = http.NewRequest(method, url, nil)
				if err != nil {
					return err
				}

				return nil
			},
			expectCode: http.StatusCreated,
			expectBody: "Goravel",
		},
		{
			name:   "Success Json",
			method: "GET",
			url:    "/response/success/json",
			setup: func(method, url string) error {
				var err error
				req, err = http.NewRequest(method, url, nil)
				if err != nil {
					return err
				}

				return nil
			},
			expectCode: http.StatusOK,
			expectBody: "{\"id\":\"1\"}",
		},
		{
			name:   "Success String",
			method: "GET",
			url:    "/response/success/string",
			setup: func(method, url string) error {
				var err error
				req, err = http.NewRequest(method, url, nil)
				if err != nil {
					return err
				}

				return nil
			},
			expectCode: http.StatusOK,
			expectBody: "Goravel",
		},
		{
			name:   "File",
			method: "GET",
			url:    "/response/file",
			setup: func(method, url string) error {
				var err error
				req, err = http.NewRequest(method, url, nil)
				if err != nil {
					return err
				}

				return nil
			},
			expectCode: http.StatusOK,
		},
		{
			name:   "Download",
			method: "GET",
			url:    "/response/download",
			setup: func(method, url string) error {
				var err error
				req, err = http.NewRequest(method, url, nil)
				if err != nil {
					return err
				}

				return nil
			},
			expectCode: http.StatusOK,
		},
		{
			name:   "Header",
			method: "GET",
			url:    "/response/header",
			setup: func(method, url string) error {
				var err error
				req, err = http.NewRequest(method, url, nil)
				if err != nil {
					return err
				}

				return nil
			},
			expectCode:   http.StatusOK,
			expectBody:   "Goravel",
			expectHeader: "goravel",
		},
	}

	for _, test := range tests {
		beforeEach()
		err := test.setup(test.method, test.url)
		assert.Nil(t, err)

		w := httptest.NewRecorder()
		facades.Route.ServeHTTP(w, req)

		if test.expectBody != "" {
			assert.Equal(t, test.expectBody, w.Body.String(), test.name)
		}
		if test.expectHeader != "" {
			assert.Equal(t, test.expectHeader, strings.Join(w.Header().Values("Hello"), ""), test.name)
		}
		assert.Equal(t, test.expectCode, w.Code, test.name)

		file.Remove("./public/test.png")
	}
}
