package client

import (
	"net/http"
	"regexp"
	"strings"

	"github.com/goravel/framework/contracts/http/client"
)

type FakeRule struct {
	pattern    string
	clientName string
	regex      *regexp.Regexp
	handler    func(client.Request) client.Response
}

func NewFakeRule(pattern string, handler func(client.Request) client.Response) *FakeRule {
	var (
		clientName string
		regex      *regexp.Regexp
	)

	// Scoped rules use the format "client#path" (e.g., "github#repos").
	// We split the string once: the left side is the client, the right is the path.
	if before, after, found := strings.Cut(pattern, "#"); found {
		clientName = before
		regex = compileWildcard(after)

		// Patterns containing URL-specific characters or wildcards are treated as URL rules.
		// We strictly check "localhost" since it lacks special characters but represents a host.
	} else if strings.ContainsAny(pattern, ".:/*") || pattern == "localhost" {
		regex = compileWildcard(pattern)

		// If no special characters are found, assume it is a simple client name.
	} else {
		clientName = pattern
	}

	return &FakeRule{
		pattern:    pattern,
		clientName: clientName,
		regex:      regex,
		handler:    handler,
	}
}

func (r *FakeRule) Matches(req *http.Request, clientName string) bool {
	if r.clientName != "" {
		if r.clientName != clientName {
			return false
		}
		if r.regex != nil {
			return r.regex.MatchString(req.URL.Path)
		}
		return true
	}

	return r.regex.MatchString(req.URL.String())
}

func compileWildcard(p string) *regexp.Regexp {
	if p == "*" {
		return regexp.MustCompile(".*")
	}

	quoted := regexp.QuoteMeta(p)

	expr := strings.ReplaceAll(quoted, "\\*", ".*")

	// If the user provided a full URL (starting with http), strictly anchor start/end.
	// If the user provided a domain (api.stripe.com), allow matching the implicit https:// prefix.
	if strings.HasPrefix(p, "http") {
		expr = "^" + expr + "$"
	} else {
		expr = "^(https?://)?" + expr + "$"
	}

	return regexp.MustCompile(expr)
}
