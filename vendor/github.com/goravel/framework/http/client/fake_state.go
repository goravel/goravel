package client

import (
	"net/http"
	"regexp"
	"sort"
	"strings"
	"sync"

	"github.com/goravel/framework/contracts/foundation"
	"github.com/goravel/framework/contracts/http/client"
)

type FakeState struct {
	mu                   sync.RWMutex
	recorded             []client.Request
	rules                []*FakeRule
	allowedStrayPatterns []*regexp.Regexp
	preventStrayRequests bool
}

func NewFakeState(json foundation.Json, mocks map[string]any) *FakeState {
	rules := make([]*FakeRule, 0, len(mocks))
	for p, v := range mocks {
		rules = append(rules, NewFakeRule(p, toHandler(json, v)))
	}

	// Sort rules to ensure the most specific pattern matches first.
	// 1. Exact matches (no wildcards) take precedence over wildcards.
	// 2. Longer patterns take precedence over shorter ones.
	// 3. Alphabetical order determines priority for identical specificity.
	sort.Slice(rules, func(i, j int) bool {
		pI, pJ := rules[i].pattern, rules[j].pattern

		noWildI := !strings.Contains(pI, "*")
		noWildJ := !strings.Contains(pJ, "*")

		if noWildI != noWildJ {
			return noWildI
		}

		if len(pI) != len(pJ) {
			return len(pI) > len(pJ)
		}

		return pI < pJ
	})

	return &FakeState{
		rules: rules,
	}
}

func (r *FakeState) Record(req client.Request) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.recorded = append(r.recorded, req)
}

func (r *FakeState) Match(req *http.Request, name string) func(client.Request) client.Response {
	for _, rule := range r.rules {
		if rule.Matches(req, name) {
			return rule.handler
		}
	}
	return nil
}

func (r *FakeState) ShouldPreventStray(url string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if !r.preventStrayRequests {
		return false
	}

	for _, p := range r.allowedStrayPatterns {
		if p.MatchString(url) {
			return false
		}
	}
	return true
}

func (r *FakeState) PreventStrayRequests() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.preventStrayRequests = true
}

func (r *FakeState) AllowStrayRequests(patterns []string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, p := range patterns {
		r.allowedStrayPatterns = append(r.allowedStrayPatterns, compileWildcard(p))
	}
}

func (r *FakeState) AssertSent(f func(client.Request) bool) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, recorded := range r.recorded {
		if f(recorded) {
			return true
		}
	}
	return false
}

func (r *FakeState) AssertSentCount(count int) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.recorded) == count
}

func toHandler(json foundation.Json, value any) func(client.Request) client.Response {
	if value == nil {
		return func(_ client.Request) client.Response {
			return NewFakeResponse(json).Status(200)
		}
	}

	switch v := value.(type) {
	case func(client.Request) client.Response:
		return v
	case client.Response:
		// Wrap single response in a Sequence to handle body snapshotting automatically.
		seq := NewFakeSequence(json)
		seq.WhenEmpty(v)

		return func(_ client.Request) client.Response {
			return seq.getNext()
		}
	case *FakeSequence:
		return func(_ client.Request) client.Response {
			return v.getNext()
		}
	case int:
		return func(_ client.Request) client.Response {
			return NewFakeResponse(json).Status(v)
		}
	case string:
		return func(_ client.Request) client.Response {
			return NewFakeResponse(json).String(200, v)
		}
	default:
		return func(_ client.Request) client.Response {
			return NewFakeResponse(json).Json(200, v)
		}
	}
}
