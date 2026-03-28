package utils

import (
	"fmt"
	"time"
)

const (
	DateTimeFormat = "2006-01-02 15:04:05"
	DateFormat     = "2006-01-02"
	TimeFormat     = "15:04:05"
)

func Now() time.Time {
	return time.Now()
}

func NowUnix() int64 {
	return time.Now().Unix()
}

func NowUnixMilli() int64 {
	return time.Now().UnixMilli()
}

func FormatDateTime(t time.Time) string {
	return t.Format(DateTimeFormat)
}

func FormatDate(t time.Time) string {
	return t.Format(DateFormat)
}

func FormatTime(t time.Time) string {
	return t.Format(TimeFormat)
}

func ParseDateTime(s string) (time.Time, error) {
	return time.Parse(DateTimeFormat, s)
}

func ParseDate(s string) (time.Time, error) {
	return time.Parse(DateFormat, s)
}

func ParseTime(s string) (time.Time, error) {
	return time.Parse(TimeFormat, s)
}

func UnixToTime(unix int64) time.Time {
	return time.Unix(unix, 0)
}

func UnixMilliToTime(unixMilli int64) time.Time {
	return time.UnixMilli(unixMilli)
}

func ToUnix(t time.Time) int64 {
	return t.Unix()
}

func ToUnixMilli(t time.Time) int64 {
	return t.UnixMilli()
}

func StartOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

func EndOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 999999999, t.Location())
}

func StartOfWeek(t time.Time) time.Time {
	weekday := int(t.Weekday())
	if weekday == 0 {
		weekday = 7
	}
	return StartOfDay(t.AddDate(0, 0, -weekday+1))
}

func EndOfWeek(t time.Time) time.Time {
	weekday := int(t.Weekday())
	if weekday == 0 {
		weekday = 7
	}
	return EndOfDay(t.AddDate(0, 0, 7-weekday))
}

func StartOfMonth(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())
}

func EndOfMonth(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month()+1, 0, 23, 59, 59, 999999999, t.Location())
}

func StartOfYear(t time.Time) time.Time {
	return time.Date(t.Year(), 1, 1, 0, 0, 0, 0, t.Location())
}

func EndOfYear(t time.Time) time.Time {
	return time.Date(t.Year(), 12, 31, 23, 59, 59, 999999999, t.Location())
}

func IsSameDay(t1, t2 time.Time) bool {
	return t1.Year() == t2.Year() && t1.YearDay() == t2.YearDay()
}

func DaysBetween(start, end time.Time) int {
	return int(end.Sub(start).Hours() / 24)
}

func AddDays(t time.Time, days int) time.Time {
	return t.AddDate(0, 0, days)
}

func AddMonths(t time.Time, months int) time.Time {
	return t.AddDate(0, months, 0)
}

func AddYears(t time.Time, years int) time.Time {
	return t.AddDate(years, 0, 0)
}

func ParseDuration(s string) (time.Duration, error) {
	return time.ParseDuration(s)
}

func FormatDuration(d time.Duration) string {
	if d < time.Minute {
		return fmt.Sprintf("%.0fs", d.Seconds())
	}
	if d < time.Hour {
		return fmt.Sprintf("%.0fm", d.Minutes())
	}
	if d < 24*time.Hour {
		return fmt.Sprintf("%.1fh", d.Hours())
	}
	return fmt.Sprintf("%.1fd", d.Hours()/24)
}

func InTimeRange(t, start, end time.Time) bool {
	return (t.Equal(start) || t.After(start)) && (t.Equal(end) || t.Before(end))
}

func FormatUnix(unix int64) string {
	return FormatDateTime(UnixToTime(unix))
}