package carbon

import (
	"github.com/dromara/carbon/v2"
)

// timezone constants
// 时区常量
const (
	Local = carbon.Local // 本地时间
	UTC   = carbon.UTC   // 世界协调时间
	GMT   = carbon.GMT   // 格林尼治标准时间
	EET   = carbon.EET   // 欧洲东部标准时间
	WET   = carbon.WET   // 欧洲西部标准时间
	CET   = carbon.CET   // 欧洲中部标准时间
	EST   = carbon.EST   // 美国东部标准时间
	MST   = carbon.MST   // 美国山地标准时间

	Cuba      = carbon.Cuba      // 古巴
	Egypt     = carbon.Egypt     // 埃及
	Eire      = carbon.Eire      // 爱尔兰
	Greenwich = carbon.Greenwich // 格林尼治
	Iceland   = carbon.Iceland   // 冰岛
	Iran      = carbon.Iran      // 伊朗
	Israel    = carbon.Israel    // 以色列
	Jamaica   = carbon.Jamaica   // 牙买加
	Japan     = carbon.Japan     // 日本
	Libya     = carbon.Libya     // 利比亚
	Poland    = carbon.Poland    // 波兰
	Portugal  = carbon.Portugal  // 葡萄牙
	PRC       = carbon.PRC       // 中国
	Singapore = carbon.Singapore // 新加坡
	Turkey    = carbon.Turkey    // 土耳其

	Shanghai   = carbon.Shanghai   // 上海
	Chongqing  = carbon.Chongqing  // 重庆
	Harbin     = carbon.Harbin     // 哈尔滨
	Urumqi     = carbon.Urumqi     // 乌鲁木齐
	HongKong   = carbon.HongKong   // 香港
	Macao      = carbon.Macao      // 澳门
	Taipei     = carbon.Taipei     // 台北
	Tokyo      = carbon.Tokyo      // 东京
	Saigon     = carbon.Saigon     // 西贡
	Seoul      = carbon.Seoul      // 首尔
	Bangkok    = carbon.Bangkok    // 曼谷
	Dubai      = carbon.Dubai      // 迪拜
	NewYork    = carbon.NewYork    // 纽约
	LosAngeles = carbon.LosAngeles // 洛杉矶
	Chicago    = carbon.Chicago    // 芝加哥
	Moscow     = carbon.Moscow     // 莫斯科
	London     = carbon.London     // 伦敦
	Berlin     = carbon.Berlin     // 柏林
	Paris      = carbon.Paris      // 巴黎
	Rome       = carbon.Rome       // 罗马
	Sydney     = carbon.Sydney     // 悉尼
	Melbourne  = carbon.Melbourne  // 墨尔本
	Darwin     = carbon.Darwin     // 达尔文
)

// month constants
// 月份常量
const (
	January   = carbon.January   // 一月
	February  = carbon.February  // 二月
	March     = carbon.March     // 三月
	April     = carbon.April     // 四月
	May       = carbon.May       // 五月
	June      = carbon.June      // 六月
	July      = carbon.July      // 七月
	August    = carbon.August    // 八月
	September = carbon.September // 九月
	October   = carbon.October   // 十月
	November  = carbon.November  // 十一月
	December  = carbon.December  // 十二月
)

// constellation constants
// 星座常量
const (
	Aries       = carbon.Aries       // 白羊座
	Taurus      = carbon.Taurus      // 金牛座
	Gemini      = carbon.Gemini      // 双子座
	Cancer      = carbon.Cancer      // 巨蟹座
	Leo         = carbon.Leo         // 狮子座
	Virgo       = carbon.Virgo       // 处女座
	Libra       = carbon.Libra       // 天秤座
	Scorpio     = carbon.Scorpio     // 天蝎座
	Sagittarius = carbon.Sagittarius // 射手座
	Capricorn   = carbon.Capricorn   // 摩羯座
	Aquarius    = carbon.Aquarius    // 水瓶座
	Pisces      = carbon.Pisces      // 双鱼座
)

// week constants
// 星期常量
const (
	Monday    = carbon.Monday    // 周一
	Tuesday   = carbon.Tuesday   // 周二
	Wednesday = carbon.Wednesday // 周三
	Thursday  = carbon.Thursday  // 周四
	Friday    = carbon.Friday    // 周五
	Saturday  = carbon.Saturday  // 周六
	Sunday    = carbon.Sunday    // 周日
)

// season constants
// 季节常量
const (
	Spring = carbon.Spring // 春季
	Summer = carbon.Summer // 夏季
	Autumn = carbon.Autumn // 秋季
	Winter = carbon.Winter // 冬季
)

// number constants
// 数字常量
const (
	YearsPerMillennium = carbon.YearsPerMillennium // 每千年1000年
	YearsPerCentury    = carbon.YearsPerCentury    // 每世纪100年
	YearsPerDecade     = carbon.YearsPerDecade     // 每十年10年
	QuartersPerYear    = carbon.QuartersPerYear    // 每年4个季度
	MonthsPerYear      = carbon.MonthsPerYear      // 每年12月
	MonthsPerQuarter   = carbon.MonthsPerQuarter   // 每季度3月
	WeeksPerNormalYear = carbon.WeeksPerNormalYear // 每常规年52周
	WeeksPerLongYear   = carbon.WeeksPerLongYear   // 每长年53周
	WeeksPerMonth      = carbon.WeeksPerMonth      // 每月4周
	DaysPerLeapYear    = carbon.DaysPerLeapYear    // 每闰年366天
	DaysPerNormalYear  = carbon.DaysPerNormalYear  // 每常规年365天
	DaysPerWeek        = carbon.DaysPerWeek        // 每周7天
	HoursPerWeek       = carbon.HoursPerWeek       // 每周168小时
	HoursPerDay        = carbon.HoursPerDay        // 每天24小时
	MinutesPerDay      = carbon.MinutesPerDay      // 每天1440分钟
	MinutesPerHour     = carbon.MinutesPerHour     // 每小时60分钟
	SecondsPerWeek     = carbon.SecondsPerWeek     // 每周604800秒
	SecondsPerDay      = carbon.SecondsPerDay      // 每天86400秒
	SecondsPerHour     = carbon.SecondsPerHour     // 每小时3600秒
	SecondsPerMinute   = carbon.SecondsPerMinute   // 每分钟60秒
)

// layout constants
// 布局模板常量
const (
	ANSICLayout              = carbon.ANSICLayout
	UnixDateLayout           = carbon.UnixDateLayout
	RubyDateLayout           = carbon.RubyDateLayout
	RFC822Layout             = carbon.RFC822Layout
	RFC822ZLayout            = carbon.RFC822ZLayout
	RFC850Layout             = carbon.RFC850Layout
	RFC1123Layout            = carbon.RFC1123Layout
	RFC1123ZLayout           = carbon.RFC1123ZLayout
	RssLayout                = carbon.RssLayout
	KitchenLayout            = carbon.KitchenLayout
	RFC2822Layout            = carbon.RFC2822Layout
	CookieLayout             = carbon.CookieLayout
	RFC3339Layout            = carbon.RFC3339Layout
	RFC3339MilliLayout       = carbon.RFC3339MilliLayout
	RFC3339MicroLayout       = carbon.RFC3339MicroLayout
	RFC3339NanoLayout        = carbon.RFC3339NanoLayout
	ISO8601Layout            = carbon.ISO8601Layout
	ISO8601MilliLayout       = carbon.ISO8601MilliLayout
	ISO8601MicroLayout       = carbon.ISO8601MicroLayout
	ISO8601NanoLayout        = carbon.ISO8601NanoLayout
	RFC1036Layout            = carbon.RFC1036Layout
	RFC7231Layout            = carbon.RFC7231Layout
	DayDateTimeLayout        = carbon.DayDateTimeLayout
	DateTimeLayout           = carbon.DateTimeLayout
	DateTimeMilliLayout      = carbon.DateTimeMilliLayout
	DateTimeMicroLayout      = carbon.DateTimeMicroLayout
	DateTimeNanoLayout       = carbon.DateTimeNanoLayout
	ShortDateTimeLayout      = carbon.ShortDateTimeLayout
	ShortDateTimeMilliLayout = carbon.ShortDateTimeMilliLayout
	ShortDateTimeMicroLayout = carbon.ShortDateTimeMicroLayout
	ShortDateTimeNanoLayout  = carbon.ShortDateTimeNanoLayout
	DateLayout               = carbon.DateLayout
	DateMilliLayout          = carbon.DateMilliLayout
	DateMicroLayout          = carbon.DateMicroLayout
	DateNanoLayout           = carbon.DateNanoLayout
	ShortDateLayout          = carbon.ShortDateLayout
	ShortDateMilliLayout     = carbon.ShortDateMilliLayout
	ShortDateMicroLayout     = carbon.ShortDateMicroLayout
	ShortDateNanoLayout      = carbon.ShortDateNanoLayout
	TimeLayout               = carbon.TimeLayout
	TimeMilliLayout          = carbon.TimeMilliLayout
	TimeMicroLayout          = carbon.TimeMicroLayout
	TimeNanoLayout           = carbon.TimeNanoLayout
	ShortTimeLayout          = carbon.ShortTimeLayout
	ShortTimeMilliLayout     = carbon.ShortTimeMilliLayout
	ShortTimeMicroLayout     = carbon.ShortTimeMicroLayout
	ShortTimeNanoLayout      = carbon.ShortTimeNanoLayout
)

// format constants
// 格式模板常量
const (
	AtomFormat     = carbon.AtomFormat
	ANSICFormat    = carbon.ANSICFormat
	CookieFormat   = carbon.CookieFormat
	KitchenFormat  = carbon.KitchenFormat
	RssFormat      = carbon.RssFormat
	RubyDateFormat = carbon.RubyDateFormat
	UnixDateFormat = carbon.UnixDateFormat

	RFC1036Format      = carbon.RFC1036Format
	RFC1123Format      = carbon.RFC1123Format
	RFC1123ZFormat     = carbon.RFC1123ZFormat
	RFC2822Format      = carbon.RFC2822Format
	RFC3339Format      = carbon.RFC3339Format
	RFC3339MilliFormat = carbon.RFC3339MilliFormat
	RFC3339MicroFormat = carbon.RFC3339MicroFormat
	RFC3339NanoFormat  = carbon.RFC3339NanoFormat
	RFC7231Format      = carbon.RFC7231Format
	RFC822Format       = carbon.RFC822Format
	RFC822ZFormat      = carbon.RFC822ZFormat
	RFC850Format       = carbon.RFC850Format

	ISO8601Format      = carbon.ISO8601Format
	ISO8601MilliFormat = carbon.ISO8601MilliFormat
	ISO8601MicroFormat = carbon.ISO8601MicroFormat
	ISO8601NanoFormat  = carbon.ISO8601NanoFormat

	ISO8601ZuluFormat      = carbon.ISO8601ZuluFormat
	ISO8601ZuluMilliFormat = carbon.ISO8601ZuluMilliFormat
	ISO8601ZuluMicroFormat = carbon.ISO8601ZuluMicroFormat
	ISO8601ZuluNanoFormat  = carbon.ISO8601ZuluNanoFormat

	FormattedDateFormat    = carbon.FormattedDateFormat
	FormattedDayDateFormat = carbon.FormattedDayDateFormat

	DayDateTimeFormat        = carbon.DayDateTimeFormat
	DateTimeFormat           = carbon.DateTimeFormat
	DateTimeMilliFormat      = carbon.DateTimeMilliFormat
	DateTimeMicroFormat      = carbon.DateTimeMicroFormat
	DateTimeNanoFormat       = carbon.DateTimeNanoFormat
	ShortDateTimeFormat      = carbon.ShortDateTimeFormat
	ShortDateTimeMilliFormat = carbon.ShortDateTimeMilliFormat
	ShortDateTimeMicroFormat = carbon.ShortDateTimeMicroFormat
	ShortDateTimeNanoFormat  = carbon.ShortDateTimeNanoFormat

	DateFormat           = carbon.DateFormat
	DateMilliFormat      = carbon.DateMilliFormat
	DateMicroFormat      = carbon.DateMicroFormat
	DateNanoFormat       = carbon.DateNanoFormat
	ShortDateFormat      = carbon.ShortDateFormat
	ShortDateMilliFormat = carbon.ShortDateMilliFormat
	ShortDateMicroFormat = carbon.ShortDateMicroFormat
	ShortDateNanoFormat  = carbon.ShortDateNanoFormat

	TimeFormat           = carbon.TimeFormat
	TimeMilliFormat      = carbon.TimeMilliFormat
	TimeMicroFormat      = carbon.TimeMicroFormat
	TimeNanoFormat       = carbon.TimeNanoFormat
	ShortTimeFormat      = carbon.ShortTimeFormat
	ShortTimeMilliFormat = carbon.ShortTimeMilliFormat
	ShortTimeMicroFormat = carbon.ShortTimeMicroFormat
	ShortTimeNanoFormat  = carbon.ShortTimeNanoFormat

	TimestampFormat      = carbon.TimestampFormat
	TimestampMilliFormat = carbon.TimestampMilliFormat
	TimestampMicroFormat = carbon.TimestampMicroFormat
	TimestampNanoFormat  = carbon.TimestampNanoFormat
)
