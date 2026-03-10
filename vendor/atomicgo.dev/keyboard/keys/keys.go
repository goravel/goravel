package keys

import "atomicgo.dev/keyboard/internal"

// Key contains information about a keypress.
type Key struct {
	Code       KeyCode
	Runes      []rune // Runes that the key produced. Most key pressed produce one single rune.
	AltPressed bool   // True when alt is pressed while the key is typed.
}

// String returns a string representation of the key.
// (e.g. "a", "B", "alt+a", "enter", "ctrl+c", "shift-down", etc.)
//
// Example:
//
//	k := keys.Key{Code: keys.Enter}
//	fmt.Println(k)
//	// Output: enter
func (k Key) String() (str string) {
	if k.AltPressed {
		str += "alt+"
	}
	if k.Code == RuneKey {
		str += string(k.Runes)

		return str
	} else if s, ok := keyNames[k.Code]; ok {
		str += s

		return str
	}

	return ""
}

// KeyCode is an integer representation of a non-rune key, such as Escape, Enter, etc.
// All other keys are represented by a rune and have the KeyCode: RuneKey.
//
// Example:
//
//	k := Key{Code: RuneKey, Runes: []rune{'x'}, AltPressed: true}
//	if k.Code == RuneKey {
//	    fmt.Println(k.Runes)
//	    // Output: x
//
//	    fmt.Println(k.String())
//	    // Output: alt+x
//	}
type KeyCode int

func (k KeyCode) String() (str string) {
	if s, ok := keyNames[k]; ok {
		return s
	}

	return ""
}

// All control keys.
const (
	Null      KeyCode = internal.KeyNull
	Break     KeyCode = internal.KeyExit
	Enter     KeyCode = internal.KeyCarriageReturn
	Backspace KeyCode = internal.KeyDelete
	Tab       KeyCode = internal.KeyHorizontalTabulation
	Esc       KeyCode = internal.KeyEscape
	Escape    KeyCode = internal.KeyEscape

	CtrlAt KeyCode = internal.KeyNull
	CtrlA  KeyCode = internal.KeyStartOfHeading
	CtrlB  KeyCode = internal.KeyStartOfText
	CtrlC  KeyCode = internal.KeyExit
	CtrlD  KeyCode = internal.KeyEndOfTransimission
	CtrlE  KeyCode = internal.KeyEnquiry
	CtrlF  KeyCode = internal.KeyAcknowledge
	CtrlG  KeyCode = internal.KeyBELL
	CtrlH  KeyCode = internal.KeyBackspace
	CtrlI  KeyCode = internal.KeyHorizontalTabulation
	CtrlJ  KeyCode = internal.KeyLineFeed
	CtrlK  KeyCode = internal.KeyVerticalTabulation
	CtrlL  KeyCode = internal.KeyFormFeed
	CtrlM  KeyCode = internal.KeyCarriageReturn
	CtrlN  KeyCode = internal.KeyShiftOut
	CtrlO  KeyCode = internal.KeyShiftIn
	CtrlP  KeyCode = internal.KeyDataLinkEscape
	CtrlQ  KeyCode = internal.KeyDeviceControl1
	CtrlR  KeyCode = internal.KeyDeviceControl2
	CtrlS  KeyCode = internal.KeyDeviceControl3
	CtrlT  KeyCode = internal.KeyDeviceControl4
	CtrlU  KeyCode = internal.KeyNegativeAcknowledge
	CtrlV  KeyCode = internal.KeySynchronousIdle
	CtrlW  KeyCode = internal.KeyEndOfTransmissionBlock
	CtrlX  KeyCode = internal.KeyCancel
	CtrlY  KeyCode = internal.KeyEndOfMedium
	CtrlZ  KeyCode = internal.KeySubstitution

	CtrlOpenBracket  KeyCode = internal.KeyEscape
	CtrlBackslash    KeyCode = internal.KeyFileSeparator
	CtrlCloseBracket KeyCode = internal.KeyGroupSeparator
	CtrlCaret        KeyCode = internal.KeyRecordSeparator
	CtrlUnderscore   KeyCode = internal.KeyUnitSeparator
	CtrlQuestionMark KeyCode = internal.KeyDelete
)

// Other keys.
const (
	RuneKey KeyCode = -(iota + 1)
	Up
	Down
	Right
	Left
	ShiftTab
	Home
	End
	PgUp
	PgDown
	Delete
	Space
	CtrlUp
	CtrlDown
	CtrlRight
	CtrlLeft
	ShiftUp
	ShiftDown
	ShiftRight
	ShiftLeft
	CtrlShiftUp
	CtrlShiftDown
	CtrlShiftLeft
	CtrlShiftRight
	F1
	F2
	F3
	F4
	F5
	F6
	F7
	F8
	F9
	F10
	F11
	F12
	F13
	F14
	F15
	F16
	F17
	F18
	F19
	F20
)

var keyNames = map[KeyCode]string{
	// Control keys.
	internal.KeyNull:                   "ctrl+@", // also ctrl+backtick
	internal.KeyStartOfHeading:         "ctrl+a",
	internal.KeyStartOfText:            "ctrl+b",
	internal.KeyExit:                   "ctrl+c",
	internal.KeyEndOfTransimission:     "ctrl+d",
	internal.KeyEnquiry:                "ctrl+e",
	internal.KeyAcknowledge:            "ctrl+f",
	internal.KeyBELL:                   "ctrl+g",
	internal.KeyBackspace:              "ctrl+h",
	internal.KeyHorizontalTabulation:   "tab", // also ctrl+i
	internal.KeyLineFeed:               "ctrl+j",
	internal.KeyVerticalTabulation:     "ctrl+k",
	internal.KeyFormFeed:               "ctrl+l",
	internal.KeyCarriageReturn:         "enter",
	internal.KeyShiftOut:               "ctrl+n",
	internal.KeyShiftIn:                "ctrl+o",
	internal.KeyDataLinkEscape:         "ctrl+p",
	internal.KeyDeviceControl1:         "ctrl+q",
	internal.KeyDeviceControl2:         "ctrl+r",
	internal.KeyDeviceControl3:         "ctrl+s",
	internal.KeyDeviceControl4:         "ctrl+t",
	internal.KeyNegativeAcknowledge:    "ctrl+u",
	internal.KeySynchronousIdle:        "ctrl+v",
	internal.KeyEndOfTransmissionBlock: "ctrl+w",
	internal.KeyCancel:                 "ctrl+x",
	internal.KeyEndOfMedium:            "ctrl+y",
	internal.KeySubstitution:           "ctrl+z",
	internal.KeyEscape:                 "esc",
	internal.KeyFileSeparator:          "ctrl+\\",
	internal.KeyGroupSeparator:         "ctrl+]",
	internal.KeyRecordSeparator:        "ctrl+^",
	internal.KeyUnitSeparator:          "ctrl+_",
	internal.KeyDelete:                 "backspace",

	// Other keys.
	RuneKey:        "runes",
	Up:             "up",
	Down:           "down",
	Right:          "right",
	Space:          "space",
	Left:           "left",
	ShiftTab:       "shift+tab",
	Home:           "home",
	End:            "end",
	PgUp:           "pgup",
	PgDown:         "pgdown",
	Delete:         "delete",
	CtrlUp:         "ctrl+up",
	CtrlDown:       "ctrl+down",
	CtrlRight:      "ctrl+right",
	CtrlLeft:       "ctrl+left",
	ShiftUp:        "shift+up",
	ShiftDown:      "shift+down",
	ShiftRight:     "shift+right",
	ShiftLeft:      "shift+left",
	CtrlShiftUp:    "ctrl+shift+up",
	CtrlShiftDown:  "ctrl+shift+down",
	CtrlShiftLeft:  "ctrl+shift+left",
	CtrlShiftRight: "ctrl+shift+right",
	F1:             "f1",
	F2:             "f2",
	F3:             "f3",
	F4:             "f4",
	F5:             "f5",
	F6:             "f6",
	F7:             "f7",
	F8:             "f8",
	F9:             "f9",
	F10:            "f10",
	F11:            "f11",
	F12:            "f12",
	F13:            "f13",
	F14:            "f14",
	F15:            "f15",
	F16:            "f16",
	F17:            "f17",
	F18:            "f18",
	F19:            "f19",
	F20:            "f20",
}
