package client

import (
	"net/http"
	"sync"

	"github.com/goravel/framework/contracts/foundation"
	"github.com/goravel/framework/contracts/http/client"
	"github.com/goravel/framework/errors"
)

var _ client.Factory = (*Factory)(nil)

type Factory struct {
	client.Request

	json    foundation.Json
	config  *FactoryConfig
	clients sync.Map
	mu      sync.RWMutex

	fakeState *FakeState
	strict    bool
	stray     []string
}

func NewFactory(config *FactoryConfig, json foundation.Json) (*Factory, error) {
	if config == nil {
		return nil, errors.HttpClientConfigNotSet
	}

	factory := &Factory{
		config: config,
		json:   json,
	}

	if err := factory.bindDefault(); err != nil {
		return nil, err
	}

	return factory, nil
}

func (r *Factory) AllowStrayRequests(patterns []string) client.Factory {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.stray = append(r.stray, patterns...)

	if r.fakeState != nil {
		r.fakeState.AllowStrayRequests(patterns)
	}

	return r
}

func (r *Factory) AssertNotSent(assertion func(client.Request) bool) bool {
	return !r.AssertSent(assertion)
}

func (r *Factory) AssertNothingSent() bool {
	return r.AssertSentCount(0)
}

func (r *Factory) AssertSent(assertion func(client.Request) bool) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.fakeState != nil && r.fakeState.AssertSent(assertion)
}

func (r *Factory) AssertSentCount(count int) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if r.fakeState != nil {
		return r.fakeState.AssertSentCount(count)
	}

	return count == 0
}

func (r *Factory) Client(names ...string) client.Request {
	name := r.config.Default
	if len(names) > 0 && names[0] != "" {
		name = names[0]
	}

	r.mu.RLock()
	state := r.fakeState
	r.mu.RUnlock()

	httpClient, err := r.resolveClient(name, state)
	if err != nil {
		return newRequestWithError(err)
	}

	cfg := r.config.Clients[name]
	return NewRequest(httpClient, r.json, cfg.BaseUrl, name)
}

func (r *Factory) Fake(mocks map[string]any) client.Factory {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.fakeState = NewFakeState(r.json, mocks)

	if r.strict {
		r.fakeState.PreventStrayRequests()
	}
	if len(r.stray) > 0 {
		r.fakeState.AllowStrayRequests(r.stray)
	}

	// Flush existing clients to force them to re-resolve with the new FakeTransport
	r.flushClients()

	return r
}

func (r *Factory) PreventStrayRequests() client.Factory {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.strict = true

	if r.fakeState != nil {
		r.fakeState.PreventStrayRequests()
	}

	return r
}

func (r *Factory) Reset() {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.fakeState = nil
	r.strict = false
	r.stray = nil

	r.flushClients()
}

func (r *Factory) Response() client.FakeResponse {
	return NewFakeResponse(r.json)
}

func (r *Factory) Sequence() client.FakeSequence {
	return NewFakeSequence(r.json)
}

func (r *Factory) bindDefault() error {
	name := r.config.Default
	c, err := r.resolveClient(name, r.fakeState)
	if err != nil {
		return err
	}

	// Bind the default client to the embedded Request implementation
	// so that methods like Http.Get() use the default configuration.
	r.Request = NewRequest(c, r.json, r.config.Clients[name].BaseUrl, name)

	return nil
}

func (r *Factory) flushClients() {
	r.clients.Range(func(key, value any) bool {
		r.clients.Delete(key)
		return true
	})

	if err := r.bindDefault(); err != nil {
		panic(err)
	}
}

func (r *Factory) resolveClient(name string, state *FakeState) (*http.Client, error) {
	if name == "" {
		return nil, errors.HttpClientDefaultNotSet
	}

	if val, ok := r.clients.Load(name); ok {
		return val.(*http.Client), nil
	}

	cfg, ok := r.config.Clients[name]
	if !ok {
		return nil, errors.HttpClientConnectionNotFound.Args(name)
	}

	// We clone the default transport to ensure we don't modify the global state
	// when applying client-specific timeouts or fake transports.
	baseTransport := http.DefaultTransport.(*http.Transport).Clone()
	baseTransport.MaxIdleConns = cfg.MaxIdleConns
	baseTransport.MaxIdleConnsPerHost = cfg.MaxIdleConnsPerHost
	baseTransport.MaxConnsPerHost = cfg.MaxConnsPerHost
	baseTransport.IdleConnTimeout = cfg.IdleConnTimeout

	var transport http.RoundTripper = baseTransport
	if state != nil {
		// If testing mode is active, wrap the real transport with our interceptor.
		transport = NewFakeTransport(state, baseTransport, r.json)
	}

	httpClient := &http.Client{
		Timeout:   cfg.Timeout,
		Transport: transport,
	}

	actual, _ := r.clients.LoadOrStore(name, httpClient)
	return actual.(*http.Client), nil
}
