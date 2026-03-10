package strutil

import (
	"fmt"
	"strings"
)

// PosFlag type
type PosFlag uint8

// Position for padding/resize string
const (
	PosLeft PosFlag = iota
	PosRight
	PosMiddle
)

/*************************************************************
 * String padding operation
 *************************************************************/

// Padding a string.
func Padding(s, pad string, length int, pos PosFlag) string {
	diff := len(s) - length
	if diff >= 0 { // do not need padding.
		return s
	}

	if pad == "" || pad == " " {
		mark := ""
		if pos == PosRight { // to right
			mark = "-"
		}

		// padding left: "%7s", padding right: "%-7s"
		tpl := fmt.Sprintf("%s%d", mark, length)
		return fmt.Sprintf(`%`+tpl+`s`, s)
	}

	if pos == PosRight { // to right
		return s + Repeat(pad, -diff)
	}
	return Repeat(pad, -diff) + s
}

// PadLeft a string.
func PadLeft(s, pad string, length int) string {
	return Padding(s, pad, length, PosLeft)
}

// PadRight a string.
func PadRight(s, pad string, length int) string {
	return Padding(s, pad, length, PosRight)
}

// Resize a string by given length and align settings. padding space.
func Resize(s string, length int, align PosFlag) string {
	diff := len(s) - length
	if diff >= 0 { // do not need padding.
		return s
	}

	if align == PosMiddle {
		padLn := (length - len(s)) / 2
		if diff := length - padLn*2; diff > 0 {
			s += " "
		}

		padStr := string(RepeatBytes(' ', padLn))
		return padStr + s + padStr
	}

	return Padding(s, " ", length, align)
}

// PadChars padding a rune/byte to want length and with position flag
func PadChars[T byte | rune](cs []T, pad T, length int, pos PosFlag) []T {
	ln := len(cs)
	if ln >= length {
		ns := make([]T, length)
		copy(ns, cs[:length])
		return ns
	}

	idx := length - ln
	ns := make([]T, length)
	if pos == PosRight {
		copy(ns, cs)
		for i := ln; i < length; i++ {
			ns[i] = pad
		}
		return ns
	}

	// to left
	for i := 0; i < idx; i++ {
		ns[i] = pad
	}
	copy(ns[idx:], cs)
	return ns
}

// PadBytes padding a byte to want length and with position flag
func PadBytes(bs []byte, pad byte, length int, pos PosFlag) []byte {
	return PadChars(bs, pad, length, pos)
}

// PadBytesLeft a byte to want length
func PadBytesLeft(bs []byte, pad byte, length int) []byte {
	return PadChars(bs, pad, length, PosLeft)
}

// PadBytesRight a byte to want length
func PadBytesRight(bs []byte, pad byte, length int) []byte {
	return PadChars(bs, pad, length, PosRight)
}

// PadRunes padding a rune to want length and with position flag
func PadRunes(rs []rune, pad rune, length int, pos PosFlag) []rune {
	return PadChars(rs, pad, length, pos)
}

// PadRunesLeft a rune to want length
func PadRunesLeft(rs []rune, pad rune, length int) []rune {
	return PadChars(rs, pad, length, PosLeft)
}

// PadRunesRight a rune to want length
func PadRunesRight(rs []rune, pad rune, length int) []rune {
	return PadChars(rs, pad, length, PosRight)
}

/*************************************************************
 * String repeat operation
 *************************************************************/

// Repeat a string by given times.
func Repeat(s string, times int) string {
	if times < 1 {
		return ""
	}
	if times == 1 {
		return s
	}

	var sb strings.Builder
	sb.Grow(len(s) * times)

	for i := 0; i < times; i++ {
		sb.WriteString(s)
	}
	return sb.String()
}

// RepeatRune repeat a rune char.
func RepeatRune(char rune, times int) []rune { return RepeatChars(char, times) }

// RepeatBytes repeat a byte char.
func RepeatBytes(char byte, times int) []byte { return RepeatChars(char, times) }

// RepeatChars repeat a byte char.
func RepeatChars[T byte | rune](char T, times int) []T {
	if times <= 0 {
		return make([]T, 0)
	}

	chars := make([]T, times)
	for i := 0; i < times; i++ {
		chars[i] = char
	}
	return chars
}
