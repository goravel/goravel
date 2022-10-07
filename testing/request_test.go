package testing

import (
	"bytes"
	"goravel/bootstrap"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/goravel/framework/facades"
	"github.com/goravel/framework/support/file"
	"github.com/stretchr/testify/assert"
)

func TestRequest(t *testing.T) {
	bootstrap.Boot()

	var (
		req *http.Request
	)

	tests := []struct {
		name       string
		method     string
		url        string
		setup      func(method, url string) error
		expectCode int
		expectBody string
	}{
		{
			name:   "Methods",
			method: "GET",
			url:    "/request/get/1?name=Goravel",
			setup: func(method, url string) error {
				var err error
				req, err = http.NewRequest(method, url, nil)
				if err != nil {
					return err
				}
				req.Header.Set("Hello", "goravel")

				return nil
			},
			expectCode: http.StatusOK,
			expectBody: "{\"full_url\":\"\",\"header\":\"goravel\",\"id\":\"1\",\"ip\":\"\",\"method\":\"GET\",\"name\":\"Goravel\",\"path\":\"/request/get/1\",\"url\":\"\"}",
		},
		{
			name:   "Headers",
			method: "GET",
			url:    "/request/headers",
			setup: func(method, url string) error {
				var err error
				req, err = http.NewRequest(method, url, nil)
				if err != nil {
					return err
				}
				req.Header.Set("Hello", "Goravel")

				return nil
			},
			expectCode: http.StatusOK,
			expectBody: "{\"Hello\":[\"Goravel\"]}",
		},
		{
			name:   "Form",
			method: "POST",
			url:    "/request/post",
			setup: func(method, url string) error {
				payload := &bytes.Buffer{}
				writer := multipart.NewWriter(payload)
				if err := writer.WriteField("name", "Goravel"); err != nil {
					return err
				}
				if err := writer.Close(); err != nil {
					return err
				}

				req, _ = http.NewRequest(method, url, payload)
				req.Header.Set("Content-Type", writer.FormDataContentType())

				return nil
			},
			expectCode: http.StatusOK,
			expectBody: "{\"name\":\"Goravel\"}",
		},
		{
			name:   "Bind",
			method: "POST",
			url:    "/request/bind",
			setup: func(method, url string) error {
				payload := strings.NewReader(`{
					"Name": "Goravel"
				}`)
				req, _ = http.NewRequest(method, url, payload)
				req.Header.Set("Content-Type", "application/json")

				return nil
			},
			expectCode: http.StatusOK,
			expectBody: "{\"name\":\"Goravel\"}",
		},
		{
			name:   "File",
			method: "POST",
			url:    "/request/file",
			setup: func(method, url string) error {
				payload := &bytes.Buffer{}
				writer := multipart.NewWriter(payload)
				logo, errFile1 := os.Open("./resources/logo.png")
				defer logo.Close()
				part1, errFile1 := writer.CreateFormFile("file", filepath.Base("./resources/logo.png"))
				_, errFile1 = io.Copy(part1, logo)
				if errFile1 != nil {
					return errFile1
				}
				err := writer.Close()
				if err != nil {
					return err
				}

				req, _ = http.NewRequest(method, url, payload)
				req.Header.Set("Content-Type", writer.FormDataContentType())

				return nil
			},
			expectCode: http.StatusOK,
			expectBody: "{\"exist\":true}",
		},
	}

	for _, test := range tests {
		err := test.setup(test.method, test.url)
		assert.Nil(t, err)

		w := httptest.NewRecorder()
		facades.Route.ServeHTTP(w, req)

		if test.expectBody != "" {
			assert.Equal(t, test.expectBody, w.Body.String(), test.name)
		}
		assert.Equal(t, test.expectCode, w.Code, test.name)

		file.Remove("./resources/test.png")
	}
}
