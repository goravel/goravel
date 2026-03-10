package carbon

import (
	"github.com/dromara/carbon/v2"
)

type (
	TimestampType      int64
	TimestampMicroType int64
	TimestampMilliType int64
	TimestampNanoType  int64

	DateTimeType      string
	DateTimeMicroType string
	DateTimeMilliType string
	DateTimeNanoType  string

	DateType      string
	DateMilliType string
	DateMicroType string
	DateNanoType  string

	TimeType      string
	TimeMilliType string
	TimeMicroType string
	TimeNanoType  string
)

type (
	Timestamp      = carbon.TimestampType[TimestampType]
	TimestampMilli = carbon.TimestampType[TimestampMilliType]
	TimestampMicro = carbon.TimestampType[TimestampMicroType]
	TimestampNano  = carbon.TimestampType[TimestampNanoType]

	DateTime      = carbon.LayoutType[DateTimeType]
	DateTimeMicro = carbon.LayoutType[DateTimeMicroType]
	DateTimeMilli = carbon.LayoutType[DateTimeMilliType]
	DateTimeNano  = carbon.LayoutType[DateTimeNanoType]

	Date      = carbon.LayoutType[DateType]
	DateMilli = carbon.LayoutType[DateMilliType]
	DateMicro = carbon.LayoutType[DateMicroType]
	DateNano  = carbon.LayoutType[DateNanoType]

	Time      = carbon.LayoutType[TimeType]
	TimeMilli = carbon.LayoutType[TimeMilliType]
	TimeMicro = carbon.LayoutType[TimeMicroType]
	TimeNano  = carbon.LayoutType[TimeNanoType]
)

func NewTimestamp(c *Carbon) *Timestamp {
	return carbon.NewTimestampType[TimestampType](c)
}
func NewTimestampMilli(c *Carbon) *TimestampMilli {
	return carbon.NewTimestampType[TimestampMilliType](c)
}
func NewTimestampMicro(c *Carbon) *TimestampMicro {
	return carbon.NewTimestampType[TimestampMicroType](c)
}
func NewTimestampNano(c *Carbon) *TimestampNano {
	return carbon.NewTimestampType[TimestampNanoType](c)
}

func NewDateTime(c *Carbon) *DateTime {
	return carbon.NewLayoutType[DateTimeType](c)
}
func NewDateTimeMilli(c *Carbon) *DateTimeMilli {
	return carbon.NewLayoutType[DateTimeMilliType](c)
}
func NewDateTimeMicro(c *Carbon) *DateTimeMicro {
	return carbon.NewLayoutType[DateTimeMicroType](c)
}
func NewDateTimeNano(c *Carbon) *DateTimeNano {
	return carbon.NewLayoutType[DateTimeNanoType](c)
}

func NewDate(c *Carbon) *Date {
	return carbon.NewLayoutType[DateType](c)
}
func NewDateMilli(c *Carbon) *DateMilli {
	return carbon.NewLayoutType[DateMilliType](c)
}
func NewDateMicro(c *Carbon) *DateMicro {
	return carbon.NewLayoutType[DateMicroType](c)
}
func NewDateNano(c *Carbon) *DateNano {
	return carbon.NewLayoutType[DateNanoType](c)
}

func NewTime(c *Carbon) *Time {
	return carbon.NewLayoutType[TimeType](c)
}
func NewTimeMilli(c *Carbon) *TimeMilli {
	return carbon.NewLayoutType[TimeMilliType](c)
}
func NewTimeMicro(c *Carbon) *TimeMicro {
	return carbon.NewLayoutType[TimeMicroType](c)
}
func NewTimeNano(c *Carbon) *TimeNano {
	return carbon.NewLayoutType[TimeNanoType](c)
}

func (t TimestampType) Precision() string {
	return carbon.PrecisionSecond
}

func (t TimestampMilliType) Precision() string {
	return carbon.PrecisionMillisecond
}

func (t TimestampMicroType) Precision() string {
	return carbon.PrecisionMicrosecond
}

func (t TimestampNanoType) Precision() string {
	return carbon.PrecisionNanosecond
}

func (t DateTimeType) Layout() string {
	return DateTimeLayout
}

func (t DateTimeMilliType) Layout() string {
	return DateTimeMilliLayout
}

func (t DateTimeMicroType) Layout() string {
	return DateTimeMicroLayout
}

func (t DateTimeNanoType) Layout() string {
	return DateTimeNanoLayout
}

func (t DateType) Layout() string {
	return DateLayout
}

func (t DateMilliType) Layout() string {
	return DateMilliLayout
}

func (t DateMicroType) Layout() string {
	return DateMicroLayout
}

func (t DateNanoType) Layout() string {
	return DateNanoLayout
}

func (t TimeType) Layout() string {
	return TimeLayout
}

func (t TimeMilliType) Layout() string {
	return TimeMilliLayout
}

func (t TimeMicroType) Layout() string {
	return TimeMicroLayout
}

func (t TimeNanoType) Layout() string {
	return TimeNanoLayout
}
