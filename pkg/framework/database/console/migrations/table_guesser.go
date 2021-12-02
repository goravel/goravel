package migrations

import (
	"regexp"
)

var CreatePatterns = []string{
	`^create_(\w+)_table$`,
	`^create_(\w+)$`,
}

var ChangePatterns = []string{
	`_(to|from|in)_(\w+)_table$`,
	`_(to|from|in)_(\w+)$`,
}

type TableGuesser struct {
}

func (receiver TableGuesser) Guess(migration string) (string, bool) {
	for _, createPattern := range CreatePatterns {
		reg := regexp.MustCompile(createPattern)
		matches := reg.FindStringSubmatch(migration)

		if len(matches) > 0 {
			return matches[1], true
		}
	}

	for _, changePattern := range ChangePatterns {
		reg := regexp.MustCompile(changePattern)
		matches := reg.FindStringSubmatch(migration)

		if len(matches) > 0 {
			return matches[2], false
		}
	}

	return "", false
}
