package urlvalues

import (
	"strconv"
	"time"
)

// Values same as net/url so use url.ParseQuery to get a parsed query string.
type Values map[string][]string

func (v Values) Has(name string) bool {
	_, ok := v[name]
	return ok
}

func (v Values) SetDefault(name string, values ...string) {
	if !v.Has(name) {
		v[name] = values
	}
}

func (v Values) Strings(name string) []string {
	return v[name]
}

func (v Values) String(name string) string {
	values := v.Strings(name)
	if len(values) == 0 {
		return ""
	}
	return values[0]
}

func (v Values) Bool(name string) (bool, error) {
	if !v.Has(name) {
		return false, nil
	}
	s := v.String(name)
	if s == "" {
		return true, nil
	}
	return strconv.ParseBool(s)
}

func (v Values) MaybeBool(name string) bool {
	flag, _ := v.Bool(name)
	return flag
}

func (v Values) Int(name string) (int, error) {
	s := v.String(name)
	if s == "" {
		return 0, nil
	}
	return strconv.Atoi(s)
}

func (v Values) MaybeInt(name string) int {
	n, _ := v.Int(name)
	return n
}

func (v Values) Int64(name string) (int64, error) {
	s := v.String(name)
	if s == "" {
		return 0, nil
	}
	return strconv.ParseInt(s, 10, 64)
}

func (v Values) MaybeInt64(name string) int64 {
	n, _ := v.Int64(name)
	return n
}

func (v Values) Float64(name string) (float64, error) {
	s := v.String(name)
	if s == "" {
		return 0, nil
	}
	return strconv.ParseFloat(s, 64)
}

func (v Values) MaybeFloat64(name string) float64 {
	n, _ := v.Float64(name)
	return n
}

func (v Values) Time(name string) (time.Time, error) {
	s := v.String(name)
	if s == "" {
		return time.Time{}, nil
	}
	return parseTimeString(s)
}

func (v Values) MaybeTime(name string) time.Time {
	tm, _ := v.Time(name)
	return tm
}

func (v Values) Duration(name string) (time.Duration, error) {
	s := v.String(name)
	if s == "" {
		return 0, nil
	}
	return time.ParseDuration(s)
}

func (v Values) MaybeDuration(name string) time.Duration {
	dur, _ := v.Duration(name)
	return dur
}

func (v Values) Pager() *Pager {
	return NewPager(v)
}

const (
	dateFormat         = "2006-01-02"
	timeFormat         = "15:04:05.999999999"
	timestampFormat    = "2006-01-02 15:04:05.999999999"
	timestamptzFormat  = "2006-01-02 15:04:05.999999999-07:00:00"
	timestamptzFormat2 = "2006-01-02 15:04:05.999999999-07:00"
	timestamptzFormat3 = "2006-01-02 15:04:05.999999999-07"
)

func parseTimeString(s string) (time.Time, error) {
	switch l := len(s); {
	case l <= len(timeFormat):
		if s[2] == ':' {
			return time.ParseInLocation(timeFormat, s, time.UTC)
		}
		return time.ParseInLocation(dateFormat, s, time.UTC)
	default:
		if s[10] == 'T' {
			return time.Parse(time.RFC3339Nano, s)
		}
		if c := s[l-9]; c == '+' || c == '-' {
			return time.Parse(timestamptzFormat, s)
		}
		if c := s[l-6]; c == '+' || c == '-' {
			return time.Parse(timestamptzFormat2, s)
		}
		if c := s[l-3]; c == '+' || c == '-' {
			return time.Parse(timestamptzFormat3, s)
		}
		return time.ParseInLocation(timestampFormat, s, time.UTC)
	}
}
