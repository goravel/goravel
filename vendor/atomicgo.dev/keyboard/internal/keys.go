package internal

// See: https://en.wikipedia.org/wiki/C0_and_C1_control_codes
const (
	KeyNull                   = 0
	KeyStartOfHeading         = 1
	KeyStartOfText            = 2
	KeyExit                   = 3 // ctrl-c
	KeyEndOfTransimission     = 4
	KeyEnquiry                = 5
	KeyAcknowledge            = 6
	KeyBELL                   = 7
	KeyBackspace              = 8
	KeyHorizontalTabulation   = 9
	KeyLineFeed               = 10
	KeyVerticalTabulation     = 11
	KeyFormFeed               = 12
	KeyCarriageReturn         = 13
	KeyShiftOut               = 14
	KeyShiftIn                = 15
	KeyDataLinkEscape         = 16
	KeyDeviceControl1         = 17
	KeyDeviceControl2         = 18
	KeyDeviceControl3         = 19
	KeyDeviceControl4         = 20
	KeyNegativeAcknowledge    = 21
	KeySynchronousIdle        = 22
	KeyEndOfTransmissionBlock = 23
	KeyCancel                 = 24
	KeyEndOfMedium            = 25
	KeySubstitution           = 26
	KeyEscape                 = 27
	KeyFileSeparator          = 28
	KeyGroupSeparator         = 29
	KeyRecordSeparator        = 30
	KeyUnitSeparator          = 31
	KeyDelete                 = 127
)
