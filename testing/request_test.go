package testing

import (
	"bytes"
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
	var (
		req *http.Request
	)

	beforeEach := func() {

	}

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
					"name": "Goravel"
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
				file, errFile1 := os.Open("./public/logo.png")
				defer file.Close()
				part1, errFile1 := writer.CreateFormFile("file", filepath.Base("./public/logo.png"))
				_, errFile1 = io.Copy(part1, file)
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
		beforeEach()
		err := test.setup(test.method, test.url)
		assert.Nil(t, err)

		w := httptest.NewRecorder()
		facades.Route.ServeHTTP(w, req)

		if test.expectBody != "" {
			assert.Equal(t, test.expectBody, w.Body.String(), test.name)
		}
		assert.Equal(t, test.expectCode, w.Code, test.name)

		file.Remove("./public/test.png")
	}
}
