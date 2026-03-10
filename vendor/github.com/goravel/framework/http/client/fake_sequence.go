package client

import (
	"bytes"
	"io"
	"net/http"
	"sync"

	"github.com/goravel/framework/contracts/foundation"
	"github.com/goravel/framework/contracts/http/client"
)

var _ client.FakeSequence = (*FakeSequence)(nil)

// responseSnapshot holds the static state of a response.
// We store this snapshot instead of the live http.Response to prevent
// "body drained" errors and to ensure thread safety during repeated test execution.
type responseSnapshot struct {
	body    string
	status  int
	headers http.Header
}

type FakeSequence struct {
	mu        sync.Mutex
	snapshots []responseSnapshot
	json      foundation.Json
	whenEmpty *responseSnapshot
	current   int
}

func NewFakeSequence(json foundation.Json) *FakeSequence {
	return &FakeSequence{
		json:      json,
		snapshots: make([]responseSnapshot, 0),
	}
}

func (r *FakeSequence) Push(response client.Response, count ...int) client.FakeSequence {
	// We consume the response body immediately to freeze its state.
	// We panic on error because if we can't read the mock body during test setup,
	// the test configuration is fundamentally invalid.
	body, err := response.Body()
	if err != nil {
		panic("Failed to read response body during FakeSequence setup: " + err.Error())
	}

	// Capture all state, including Cookies (which are just "Set-Cookie" headers).
	snapshot := responseSnapshot{
		body:    body,
		status:  response.Status(),
		headers: response.Headers(),
	}

	return r.pushSnapshot(snapshot, count...)
}

func (r *FakeSequence) PushStatus(status int, count ...int) client.FakeSequence {
	snapshot := responseSnapshot{
		status:  status,
		body:    "",
		headers: http.Header{},
	}

	return r.pushSnapshot(snapshot, count...)
}

func (r *FakeSequence) PushString(status int, body string, count ...int) client.FakeSequence {
	snapshot := responseSnapshot{
		body:    body,
		status:  status,
		headers: http.Header{},
	}

	return r.pushSnapshot(snapshot, count...)
}

func (r *FakeSequence) WhenEmpty(response client.Response) client.FakeSequence {
	body, err := response.Body()
	if err != nil {
		panic("Failed to read response body during FakeSequence WhenEmpty setup: " + err.Error())
	}

	snapshot := &responseSnapshot{
		body:    body,
		status:  response.Status(),
		headers: response.Headers(),
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	r.whenEmpty = snapshot

	return r
}

// getNext is an internal iterator used by the HTTP client to retrieve the next mock response.
// It reconstructs a fresh http.Response object from the stored snapshot, ensuring that
// the response body stream is fresh and readable for every request.
func (r *FakeSequence) getNext() client.Response {
	r.mu.Lock()
	defer r.mu.Unlock()

	var snapshot *responseSnapshot

	// Select the correct Snapshot based on the current index
	if r.current < len(r.snapshots) {
		snapshot = &r.snapshots[r.current]
		r.current++
	} else {
		snapshot = r.whenEmpty
	}

	// If sequence is exhausted and no WhenEmpty is set, return nil
	if snapshot == nil {
		return nil
	}

	return r.buildResponse(snapshot)
}

// buildResponse converts a static Snapshot back into a live, working Response.
func (r *FakeSequence) buildResponse(snapshot *responseSnapshot) client.Response {
	// We must clone headers so that if the user modifies the returned response,
	// it doesn't corrupt the snapshot for future iterations.
	newHeaders := snapshot.headers.Clone()

	httpResp := &http.Response{
		StatusCode: snapshot.status,
		// Status text is inferred from StatusCode if missing.
		Status: http.StatusText(snapshot.status),
		Header: newHeaders,
		// Create a fresh body stream
		// This creates a new Reader from the stored string, effectively "rewinding" the stream.
		Body:          io.NopCloser(bytes.NewBufferString(snapshot.body)),
		ContentLength: int64(len(snapshot.body)),
	}

	return NewResponse(httpResp, r.json)
}

func (r *FakeSequence) pushSnapshot(snapshot responseSnapshot, count ...int) client.FakeSequence {
	r.mu.Lock()
	defer r.mu.Unlock()

	times := 1
	if len(count) > 0 && count[0] > 0 {
		times = count[0]
	}

	for i := 0; i < times; i++ {
		r.snapshots = append(r.snapshots, snapshot)
	}

	return r
}
