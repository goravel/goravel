package carbon

import (
	stdtime "time"

	"github.com/dromara/carbon/v2"
)

type Carbon = carbon.Carbon

// New returns a new Carbon object.
func New(time ...stdtime.Time) *Carbon {
	return carbon.NewCarbon(time...)
}

// ZeroValue returns the zero value of Carbon object.
func ZeroValue() *Carbon {
	return carbon.ZeroValue()
}

// EpochValue returns the unix epoch value of Carbon object.
func EpochValue() *Carbon {
	return carbon.EpochValue()
}

// DefaultTimezone gets the default timezone.
func DefaultTimezone() string {
	return carbon.DefaultTimezone
}

// SetTimezone sets timezone.
func SetTimezone(timezone string) {
	carbon.SetTimezone(timezone)
}

// SetLocale sets language locale.
func SetLocale(locale string) {
	carbon.SetLocale(locale)
}

// SetTestNow sets the test time, remember to clean after use.
func SetTestNow(c *Carbon) {
	carbon.SetTestNow(c)
}

// UnsetTestNow unsets the test time.
//
// Deprecated: it will be removed in the future, use `ClearTestNow` instead.
func UnsetTestNow() {
	ClearTestNow()
}

// ClearTestNow clears the test time.
func ClearTestNow() {
	carbon.ClearTestNow()
}

// IsTestNow determines if the test now time is set.
func IsTestNow() bool {
	return carbon.IsTestNow()
}

// Now returns a Carbon object of now.
func Now(timezone ...string) *Carbon {
	return carbon.Now(timezone...)
}

// Parse returns a Carbon object of given value.
func Parse(value string, timezone ...string) *Carbon {
	return carbon.Parse(value, timezone...)
}

// ParseByLayout returns a Carbon object by a confirmed layout.
func ParseByLayout[T string | []string](value string, layout T, timezone ...string) (c *Carbon) {
	switch v := any(layout).(type) {
	case string:
		c = carbon.ParseByLayout(value, v, timezone...)
	case []string:
		c = carbon.ParseByLayouts(value, v, timezone...)
	}
	return
}

// ParseByFormat returns a Carbon object by a confirmed format.
func ParseByFormat[T string | []string](value string, format T, timezone ...string) (c *Carbon) {
	switch v := any(format).(type) {
	case string:
		c = carbon.ParseByFormat(value, v, timezone...)
	case []string:
		c = carbon.ParseByFormats(value, v, timezone...)
	}
	return
}

// ParseWithLayouts returns a Carbon object with multiple fuzzy layouts.
//
// Deprecated: it will be removed in the future, use `ParseByLayout` instead.
func ParseWithLayouts(value string, layouts []string, timezone ...string) *Carbon {
	return ParseByLayout(value, layouts, timezone...)
}

// ParseWithFormats returns a Carbon object with multiple fuzzy formats.
//
// Deprecated: it will be removed in the future, use `ParseByFormat` instead.
func ParseWithFormats(value string, formats []string, timezone ...string) *Carbon {
	return ParseByFormat(value, formats, timezone...)
}

// FromTimestamp returns a Carbon object of given timestamp.
func FromTimestamp(timestamp int64, timezone ...string) *Carbon {
	return carbon.CreateFromTimestamp(timestamp, timezone...)
}

// FromTimestampMilli returns a Carbon object of given millisecond timestamp.
func FromTimestampMilli(timestamp int64, timezone ...string) *Carbon {
	return carbon.CreateFromTimestampMilli(timestamp, timezone...)
}

// FromTimestampMicro returns a Carbon object of given microsecond timestamp.
func FromTimestampMicro(timestamp int64, timezone ...string) *Carbon {
	return carbon.CreateFromTimestampMicro(timestamp, timezone...)
}

// FromTimestampNano returns a Carbon object of given nanosecond timestamp.
func FromTimestampNano(timestamp int64, timezone ...string) *Carbon {
	return carbon.CreateFromTimestampNano(timestamp, timezone...)
}

// FromDateTime returns a Carbon object of given date and time.
func FromDateTime(year int, month int, day int, hour int, minute int, second int, timezone ...string) *Carbon {
	return carbon.CreateFromDateTime(year, month, day, hour, minute, second, timezone...)
}

// FromDateTimeMilli returns a Carbon object of given date and millisecond time.
func FromDateTimeMilli(year int, month int, day int, hour int, minute int, second int, millisecond int, timezone ...string) *Carbon {
	return carbon.CreateFromDateTimeMilli(year, month, day, hour, minute, second, millisecond, timezone...)
}

// FromDateTimeMicro returns a Carbon object of given date and microsecond time.
func FromDateTimeMicro(year int, month int, day int, hour int, minute int, second int, microsecond int, timezone ...string) *Carbon {
	return carbon.CreateFromDateTimeMicro(year, month, day, hour, minute, second, microsecond, timezone...)
}

// FromDateTimeNano returns a Carbon object of given date and nanosecond time.
func FromDateTimeNano(year int, month int, day int, hour int, minute int, second int, nanosecond int, timezone ...string) *Carbon {
	return carbon.CreateFromDateTimeNano(year, month, day, hour, minute, second, nanosecond, timezone...)
}

// FromDate returns a Carbon object of given date.
func FromDate(year int, month int, day int, timezone ...string) *Carbon {
	return carbon.CreateFromDate(year, month, day, timezone...)
}

// FromDateMilli returns a Carbon object of given millisecond date.
func FromDateMilli(year int, month int, day int, millisecond int, timezone ...string) *Carbon {
	return carbon.CreateFromDateMilli(year, month, day, millisecond, timezone...)
}

// FromDateMicro returns a Carbon object of given microsecond date.
func FromDateMicro(year int, month int, day int, microsecond int, timezone ...string) *Carbon {
	return carbon.CreateFromDateMicro(year, month, day, microsecond, timezone...)
}

// FromDateNano returns a Carbon object of given nanosecond date.
func FromDateNano(year int, month int, day int, nanosecond int, timezone ...string) *Carbon {
	return carbon.CreateFromDateNano(year, month, day, nanosecond, timezone...)
}

// FromTime returns a Carbon object of given time.
func FromTime(hour int, minute int, second int, timezone ...string) *Carbon {
	return carbon.CreateFromTime(hour, minute, second, timezone...)
}

// FromTimeMilli returns a Carbon object of given millisecond time.
func FromTimeMilli(hour int, minute int, second int, millisecond int, timezone ...string) *Carbon {
	return carbon.CreateFromTimeMilli(hour, minute, second, millisecond, timezone...)
}

// FromTimeMicro returns a Carbon object of given microsecond time.
func FromTimeMicro(hour int, minute int, second int, microsecond int, timezone ...string) *Carbon {
	return carbon.CreateFromTimeMicro(hour, minute, second, microsecond, timezone...)
}

// FromTimeNano returns a Carbon object of given nanosecond time.
func FromTimeNano(hour int, minute int, second int, nanosecond int, timezone ...string) *Carbon {
	return carbon.CreateFromTimeNano(hour, minute, second, nanosecond, timezone...)
}

// FromStdTime returns a Carbon object of given time.Time object.
func FromStdTime(time stdtime.Time, timezone ...string) *Carbon {
	return carbon.CreateFromStdTime(time, timezone...)
}

// FromLunar returns a Carbon object from Lunar date.
func FromLunar(year, month, day int, isLeapMonth bool) *Carbon {
	return carbon.CreateFromLunar(year, month, day, isLeapMonth)
}

// FromJulian returns a Carbon object from Julian Day or Modified Julian Day.
func FromJulian(f float64) *Carbon {
	return carbon.CreateFromJulian(f)
}

// FromPersian returns a Carbon object from Persian date.
func FromPersian(year, month, day int) *Carbon {
	return carbon.CreateFromPersian(year, month, day)
}

// FromHebrew returns a Carbon object from Hebrew date.
func FromHebrew(year, month, day int) *Carbon {
	return carbon.CreateFromHebrew(year, month, day)
}
