package client

import (
	"bytes"
	"io"
	"net/http"
	"os"

	"github.com/goravel/framework/contracts/foundation"
	"github.com/goravel/framework/contracts/http/client"
)

var _ client.FakeResponse = (*FakeResponse)(nil)

type FakeResponse struct {
	json foundation.Json
}

func NewFakeResponse(json foundation.Json) *FakeResponse {
	return &FakeResponse{
		json: json,
	}
}

func (r *FakeResponse) File(status int, path string) client.Response {
	content, err := os.ReadFile(path)
	if err != nil {
		return r.make(http.StatusInternalServerError, "Failed to read mock file "+path+": "+err.Error(), nil)
	}

	return r.make(status, string(content), nil)
}

func (r *FakeResponse) Json(status int, data any) client.Response {
	content, err := r.json.Marshal(data)
	if err != nil {
		return r.make(http.StatusInternalServerError, "Failed to marshal mock JSON: "+err.Error(), nil)
	}

	header := http.Header{}
	header.Set("Content-Type", "application/json")

	return r.make(status, string(content), header)
}

func (r *FakeResponse) Make(status int, body string, header http.Header) client.Response {
	return r.make(status, body, header)
}

func (r *FakeResponse) OK() client.Response {
	return r.Status(http.StatusOK)
}

func (r *FakeResponse) Status(status int) client.Response {
	return r.make(status, "", nil)
}

func (r *FakeResponse) String(status int, body string) client.Response {
	return r.make(status, body, nil)
}

func (r *FakeResponse) make(status int, body string, header http.Header) client.Response {
	resp := &http.Response{
		StatusCode: status,
		Header:     make(http.Header),
		Body:       io.NopCloser(bytes.NewBufferString(body)),
	}

	for key, values := range header {
		for _, value := range values {
			resp.Header.Add(key, value)
		}
	}

	return NewResponse(resp, r.json)
}
