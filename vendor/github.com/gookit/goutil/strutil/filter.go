package strutil

import "strings"

/*************************************************************
 * String filtering
 *************************************************************/

// Trim string. if cutSet is empty, will trim SPACE.
func Trim(s string, cutSet ...string) string {
	if ln := len(cutSet); ln > 0 && cutSet[0] != "" {
		if ln == 1 {
			return strings.Trim(s, cutSet[0])
		}

		return strings.Trim(s, strings.Join(cutSet, ""))
	}

	return strings.TrimSpace(s)
}

// Ltrim alias of TrimLeft
func Ltrim(s string, cutSet ...string) string { return TrimLeft(s, cutSet...) }

// LTrim alias of TrimLeft
func LTrim(s string, cutSet ...string) string { return TrimLeft(s, cutSet...) }

// TrimLeft char in the string. if cutSet is empty, will trim SPACE.
func TrimLeft(s string, cutSet ...string) string {
	if ln := len(cutSet); ln > 0 && cutSet[0] != "" {
		if ln == 1 {
			return strings.TrimLeft(s, cutSet[0])
		}

		return strings.TrimLeft(s, strings.Join(cutSet, ""))
	}

	return strings.TrimLeft(s, " ")
}

// Rtrim alias of TrimRight
func Rtrim(s string, cutSet ...string) string { return TrimRight(s, cutSet...) }

// RTrim alias of TrimRight
func RTrim(s string, cutSet ...string) string { return TrimRight(s, cutSet...) }

// TrimRight char in the string. if cutSet is empty, will trim SPACE.
func TrimRight(s string, cutSet ...string) string {
	if ln := len(cutSet); ln > 0 && cutSet[0] != "" {
		if ln == 1 {
			return strings.TrimRight(s, cutSet[0])
		}
		return strings.TrimRight(s, strings.Join(cutSet, ""))
	}

	return strings.TrimRight(s, " ")
}

// FilterEmail filter email, clear invalid chars.
func FilterEmail(s string) string {
	s = strings.TrimSpace(s)
	i := strings.LastIndex(s, "@")
	if i == -1 {
		return s
	}

	// According to rfc5321, "The local-part of a mailbox MUST BE treated as case-sensitive"
	return s[0:i] + "@" + strings.ToLower(s[i+1:])
}

// func Filter(ss []string, fls ...comdef.StringMatchFunc) []string  {
// }
