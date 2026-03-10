package inflector

import (
	"strings"
	"unicode"
)

const (
	styleFallback = iota
	styleAllLower
	styleAllUpper
	styleUcFirst
	styleUcWords
)

// MatchCase transforms s to match the casing style of pattern.
func MatchCase(s, pattern string) string {
	switch detectPatternStyle(pattern) {
	case styleAllLower:
		return strings.ToLower(s)
	case styleAllUpper:
		return strings.ToUpper(s)
	case styleUcFirst:
		return UcFirst(s)
	case styleUcWords:
		return UcWords(s)
	default:
		return s
	}
}

// detectPatternStyle scans pattern once to determine its case style.
func detectPatternStyle(p string) int {
	hasLetter := false
	isAllLower, isAllUpper := true, true
	firstLetter, firstLetterSeen := rune(0), false
	isUcFirst, isUcWords := true, true
	inWord := false

	for _, r := range p {
		if !unicode.IsLetter(r) {
			inWord = false
			continue
		}
		hasLetter = true

		if !firstLetterSeen {
			firstLetter = r
			firstLetterSeen = true
		} else if !unicode.IsLower(r) {
			isUcFirst = false
		}

		if !unicode.IsLower(r) {
			isAllLower = false
		}
		if !unicode.IsUpper(r) {
			isAllUpper = false
		}

		if inWord {
			if !unicode.IsLower(r) {
				isUcWords = false
			}
		} else {
			if !unicode.IsUpper(r) {
				isUcWords = false
			}
			inWord = true
		}
	}

	switch {
	case hasLetter && isAllLower:
		return styleAllLower
	case hasLetter && isAllUpper:
		return styleAllUpper
	case hasLetter && unicode.IsUpper(firstLetter) && isUcFirst:
		return styleUcFirst
	case hasLetter && isUcWords:
		return styleUcWords
	default:
		return styleFallback
	}
}

// UcFirst uppercases the first letter and lowercases the rest of the first word.
func UcFirst(s string) string {
	runes := []rune(s)
	start := -1
	for i, r := range runes {
		if unicode.IsLetter(r) {
			if start == -1 {
				runes[i] = unicode.ToUpper(r)
				start = i
			} else {
				runes[i] = unicode.ToLower(r)
			}
		} else if start != -1 {
			break
		}
	}
	return string(runes)
}

// UcWords uppercases the first letter of each word, lowercases the rest.
func UcWords(s string) string {
	runes := []rune(s)
	inWord := false
	for i, r := range runes {
		if unicode.IsLetter(r) {
			if !inWord {
				runes[i] = unicode.ToUpper(r)
				inWord = true
			} else {
				runes[i] = unicode.ToLower(r)
			}
		} else {
			inWord = false
		}
	}
	return string(runes)
}
