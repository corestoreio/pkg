// Copyright © 2014 Steve Francia <spf@spf13.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package conv

import (
	"fmt"
	"html/template"
	"math"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/corestoreio/csfw/config/cfgpath"
	"github.com/corestoreio/csfw/storage/text"
	"github.com/corestoreio/errors"
)

// ToTimeE casts an empty interface to time.Time.
// Supported types: time.Time, string, int64 and float64.
// float64 fractals gets applied as nsec.
func ToTimeE(i interface{}) (time.Time, error) {
	i = indirect(i)

	switch s := i.(type) {
	case time.Time:
		return s, nil
	case string:
		d, e := StringToDate(s, nil)
		if e == nil {
			return d, nil
		}
		return time.Time{}, errors.NewNotValidf("[conv] Could not parse Date/Time format: %v\n", e)
	case int64:
		return time.Unix(s, 0), nil
	case float64:
		fi, frac := math.Modf(s)
		return time.Unix(int64(fi), int64(frac)), nil
	default:
		return time.Time{}, errors.NewNotValidf("[conv] Unable to cast %#v to Time\n", i)
	}
}

// ToDurationE casts an empty interface to time.Duration.
func ToDurationE(i interface{}) (d time.Duration, err error) {
	i = indirect(i)

	switch s := i.(type) {
	case time.Duration:
		return s, nil
	case int64:
		d = time.Duration(s)
		return
	case float64:
		d = time.Duration(s)
		return
	case string:
		d, err = time.ParseDuration(s)
		return
	default:
		err = errors.NewNotValidf("[conv] Unable to cast %#v to Duration\n", i)
		return
	}
}

// ToBoolE casts an empty interface to a bool. If a type implements function
//		ToBool() bool
// this function will get called.
func ToBoolE(i interface{}) (bool, error) {

	// if a type implements an iFacer interface then we call this. we call it
	// only at the end of the type switch because it costs performance to assert
	// to the iFacer interface. TODO(CyS): Implement for all other functions.
	type iFacer interface {
		ToBool() bool
	}

	i = indirect(i)

	switch b := i.(type) {
	case bool:
		return b, nil
	case nil:
		return false, nil
	case int:
		if b != 0 {
			return true, nil
		}
		return false, nil
	case int8:
		if b != 0 {
			return true, nil
		}
		return false, nil
	case int16:
		if b != 0 {
			return true, nil
		}
		return false, nil
	case int32:
		if b != 0 {
			return true, nil
		}
		return false, nil
	case int64:
		if b != 0 {
			return true, nil
		}
		return false, nil
	case uint:
		return b > 0, nil
	case uint8:
		return b > 0, nil
	case uint16:
		return b > 0, nil
	case uint32:
		return b > 0, nil
	case uint64:
		return b > 0, nil
	case float64:
		return b > 0, nil
	case float32:
		return b > 0, nil
	case string:
		switch b {
		case "YES", "yes":
			return true, nil
		case "NO", "no":
			return false, nil
		}
		b2, err := strconv.ParseBool(b)
		if err != nil {
			return false, errors.NewNotValidf("[conv] Unable to cast %#v to bool", i)
		}
		return b2, nil
	case iFacer:
		return b.ToBool(), nil
	default:
		return false, errors.NewNotValidf("[conv] Unable to cast %#v to bool", i)
	}
}

// ToFloat64E casts an empty interface to a float64.
func ToFloat64E(i interface{}) (float64, error) {
	i = indirect(i)

	switch s := i.(type) {
	case float64:
		return s, nil
	case float32:
		return float64(s), nil
	case int64:
		return float64(s), nil
	case int32:
		return float64(s), nil
	case int16:
		return float64(s), nil
	case int8:
		return float64(s), nil
	case int:
		return float64(s), nil
	case uint64:
		return float64(s), nil
	case uint32:
		return float64(s), nil
	case uint16:
		return float64(s), nil
	case uint8:
		return float64(s), nil
	case uint:
		return float64(s), nil
	case string:
		v, err := strconv.ParseFloat(s, 64)
		if err == nil {
			return float64(v), nil
		}
		return 0.0, errors.NewNotValidf("[conv] Unable to cast %#v to float. %s", i, err)
	case []byte:
		// real byte encoded floats will fail here
		// @see https://github.com/golang/go/issues/2632
		v, err := strconv.ParseFloat(string(s), 64)
		if err == nil {
			return float64(v), nil
		}
		return 0.0, errors.NewNotValidf("[conv] Unable to cast %#v to float. %s", i, err)
	default:
		return 0.0, errors.NewNotValidf("[conv] Unable to cast %#v to float", i)
	}
}

// ToIntE casts an empty interface to an int.
func ToIntE(i interface{}) (int, error) {
	i = indirect(i)

	switch s := i.(type) {
	case int:
		return s, nil
	case int64:
		return int(s), nil
	case int32:
		return int(s), nil
	case int16:
		return int(s), nil
	case int8:
		return int(s), nil
	case string:
		v, err := strconv.ParseInt(s, 0, 0)
		if err == nil {
			return int(v), nil
		}
		return 0, errors.NewNotValidf("[conv] Unable to cast %#v to int. %s", i, err)
	case float64:
		return int(s), nil
	case bool:
		if bool(s) {
			return 1, nil
		}
		return 0, nil
	case nil:
		return 0, nil
	default:
		return 0, errors.NewNotValidf("[conv] Unable to cast %#v to int", i)
	}
}

// ToUintE casts an empty interface to an uint.
func ToUintE(i interface{}) (uint, error) {
	i = indirect(i)

	switch s := i.(type) {
	case int:
		if s > 0 {
			return uint(s), nil
		}
		return 0, errors.NewNotValidf("[conv] Unable to cast %#v to uint", i)
	case int64:
		if s > 0 {
			return uint(s), nil
		}
		return 0, errors.NewNotValidf("[conv] Unable to cast %#v to uint", i)
	case int32:
		if s > 0 {
			return uint(s), nil
		}
		return 0, errors.NewNotValidf("[conv] Unable to cast %#v to uint", i)
	case int16:
		if s > 0 {
			return uint(s), nil
		}
		return 0, errors.NewNotValidf("[conv] Unable to cast %#v to uint", i)
	case int8:
		if s > 0 {
			return uint(s), nil
		}
		return 0, errors.NewNotValidf("[conv] Unable to cast %#v to uint", i)
	case uint:
		return uint(s), nil
	case uint64:
		return uint(s), nil
	case uint32:
		return uint(s), nil
	case uint16:
		return uint(s), nil
	case uint8:
		return uint(s), nil
	case string:
		v, err := strconv.ParseUint(s, 10, strconv.IntSize)
		if err == nil {
			return uint(v), nil
		}
		return 0, errors.NewNotValidf("[conv] Unable to cast %#v to uint. %s", i, err)
	case float64:
		if s > 0 && s < math.MaxUint64 {
			return uint(s), nil
		}
		return 0, errors.NewNotValidf("[conv] Unable to cast %#v to uint", i)
	case float32:
		if s > 0 && s < math.MaxUint32 {
			return uint(s), nil
		}
		return 0, errors.NewNotValidf("[conv] Unable to cast %#v to uint", i)
	case bool:
		if s {
			return 1, nil
		}
		return 0, nil
	case nil:
		return 0, nil
	default:
		return 0, errors.NewNotValidf("[conv] Unable to cast %#v to uint", i)
	}
}

// ToInt64E casts an empty interface to an int64.
func ToInt64E(i interface{}) (int64, error) {
	i = indirect(i)

	switch s := i.(type) {
	case int:
		return int64(s), nil
	case int64:
		return s, nil
	case int32:
		return int64(s), nil
	case int16:
		return int64(s), nil
	case int8:
		return int64(s), nil
	case string:
		v, err := strconv.ParseInt(s, 0, 64)
		if err == nil {
			return v, nil
		}
		return 0, errors.NewNotValidf("[conv] Unable to cast %#v to int64. %s", i, err)
	case float64:
		return int64(s), nil
	case float32:
		return int64(s), nil
	case bool:
		if bool(s) {
			return 1, nil
		}
		return 0, nil
	case nil:
		return 0, nil
	default:
		return 0, errors.NewNotValidf("[conv] Unable to cast %#v to int64", i)
	}
}

// From html/template/content.go
// Copyright 2011 The Go Authors. All rights reserved.
// indirect returns the value, after dereferencing as many times
// as necessary to reach the base type (or nil).
func indirect(a interface{}) interface{} {
	if a == nil {
		return nil
	}
	if t := reflect.TypeOf(a); t.Kind() != reflect.Ptr {
		// Avoid creating a reflect.Value if it's not a pointer.
		return a
	}
	v := reflect.ValueOf(a)
	for v.Kind() == reflect.Ptr && !v.IsNil() {
		v = v.Elem()
	}
	return v.Interface()
}

// From html/template/content.go
// Copyright 2011 The Go Authors. All rights reserved.
// indirectToStringerOrError returns the value, after dereferencing as many times
// as necessary to reach the base type (or nil) or an implementation of fmt.Stringer
// or error,
var errorType = reflect.TypeOf((*error)(nil)).Elem()
var fmtStringerType = reflect.TypeOf((*fmt.Stringer)(nil)).Elem()

func indirectToStringerOrError(a interface{}) interface{} {
	if a == nil {
		return nil
	}

	v := reflect.ValueOf(a)
	for !v.Type().Implements(fmtStringerType) && !v.Type().Implements(errorType) && v.Kind() == reflect.Ptr && !v.IsNil() {
		v = v.Elem()
	}
	return v.Interface()
}

// ToStringE casts an empty interface to a string.
func ToStringE(i interface{}) (string, error) {
	i = indirectToStringerOrError(i) // does not cost neither B/op nor allocs/op

	switch s := i.(type) {
	case string:
		return s, nil
	case bool:
		return strconv.FormatBool(s), nil
	case float64:
		return strconv.FormatFloat(i.(float64), 'f', -1, 64), nil
	case int:
		return strconv.FormatInt(int64(i.(int)), 10), nil
	case []byte:
		return string(s), nil
	case text.Chars:
		return s.String(), nil
	case cfgpath.Route:
		return s.String(), nil
	case cfgpath.Path:
		sp, err := s.FQ()
		return sp.String(), err
	case template.HTML:
		return string(s), nil
	case template.URL:
		return string(s), nil
	case nil:
		return "", nil
	case fmt.Stringer:
		return s.String(), nil
	case error:
		return s.Error(), nil
	default:
		return "", errors.NewNotValidf("[conv] Unable to cast %#v to string", i)
	}
}

// ToByteE casts an empty interface to a byte slice. Faster than ToStringE because
// ToByteE avoids some copying of data. Use wisely.
func ToByteE(i interface{}) ([]byte, error) {

	switch s := i.(type) {
	case []byte:
		return s, nil
	case string:
		return []byte(s), nil
	case bool:
		var empty = make([]byte, 0, 4)
		return strconv.AppendBool(empty, s), nil
	case float64:
		var empty = make([]byte, 0, 8)
		return strconv.AppendFloat(empty, i.(float64), 'f', -1, 64), nil
	case int:
		var empty = make([]byte, 0, 8)
		return strconv.AppendInt(empty, int64(i.(int)), 10), nil
	case int64:
		var empty = make([]byte, 0, 8)
		return strconv.AppendInt(empty, i.(int64), 10), nil
	case text.Chars:
		return s.Bytes(), nil
	case cfgpath.Route:
		return s.Bytes(), nil
	case cfgpath.Path:
		sp, err := s.FQ()
		return sp.Bytes(), err
	case template.HTML:
		return []byte(s), nil
	case template.URL:
		return []byte(s), nil
	case nil:
		return nil, nil
	default:
		return nil, errors.NewNotValidf("[conv] Unable to cast %#v to []byte", i)
	}
}

// ToStringMapStringE casts an empty interface to a map[string]string.
func ToStringMapStringE(i interface{}) (map[string]string, error) {

	var m = map[string]string{}

	switch v := i.(type) {
	case map[string]string:
		return v, nil
	case map[string]interface{}:
		for k, val := range v {
			m[ToString(k)] = ToString(val)
		}
		return m, nil
	case map[interface{}]string:
		for k, val := range v {
			m[ToString(k)] = ToString(val)
		}
		return m, nil
	case map[interface{}]interface{}:
		for k, val := range v {
			m[ToString(k)] = ToString(val)
		}
		return m, nil
	default:
		return m, errors.NewNotValidf("[conv] Unable to cast %#v to map[string]string", i)
	}
}

// ToStringMapStringSliceE casts an empty interface to a map[string][]string.
func ToStringMapStringSliceE(i interface{}) (map[string][]string, error) {

	var m = map[string][]string{}

	switch v := i.(type) {
	case map[string][]string:
		return v, nil
	case map[string][]interface{}:
		for k, val := range v {
			m[ToString(k)] = ToStringSlice(val)
		}
		return m, nil
	case map[string]string:
		for k, val := range v {
			m[ToString(k)] = []string{val}
		}
	case map[string]interface{}:
		for k, val := range v {
			m[ToString(k)] = []string{ToString(val)}
		}
		return m, nil
	case map[interface{}][]string:
		for k, val := range v {
			m[ToString(k)] = ToStringSlice(val)
		}
		return m, nil
	case map[interface{}]string:
		for k, val := range v {
			m[ToString(k)] = ToStringSlice(val)
		}
		return m, nil
	case map[interface{}][]interface{}:
		for k, val := range v {
			m[ToString(k)] = ToStringSlice(val)
		}
		return m, nil
	case map[interface{}]interface{}:
		for k, val := range v {
			key, err := ToStringE(k)
			if err != nil {
				return m, errors.NewNotValidf("[conv] Unable to cast %#v to map[string][]string. %s", i, err)
			}
			value, err := ToStringSliceE(val)
			if err != nil {
				return m, errors.NewNotValidf("[conv] Unable to cast %#v to map[string][]string. %s", i, err)
			}
			m[key] = value
		}
	default:
		return m, errors.NewNotValidf("[conv] Unable to cast %#v to map[string][]string", i)
	}
	return m, nil
}

// ToStringMapBoolE casts an empty interface to a map[string]bool.
func ToStringMapBoolE(i interface{}) (map[string]bool, error) {

	var m = map[string]bool{}

	switch v := i.(type) {
	case map[interface{}]interface{}:
		for k, val := range v {
			m[ToString(k)] = ToBool(val)
		}
		return m, nil
	case map[string]interface{}:
		for k, val := range v {
			m[ToString(k)] = ToBool(val)
		}
		return m, nil
	case map[string]bool:
		return v, nil
	default:
		return m, errors.NewNotValidf("[conv] Unable to cast %#v to map[string]bool", i)
	}
}

// ToStringMapE casts an empty interface to a map[string]interface{}.
func ToStringMapE(i interface{}) (map[string]interface{}, error) {

	var m = map[string]interface{}{}

	switch v := i.(type) {
	case map[interface{}]interface{}:
		for k, val := range v {
			m[ToString(k)] = val
		}
		return m, nil
	case map[string]interface{}:
		return v, nil
	default:
		return m, errors.NewNotValidf("[conv] Unable to cast %#v to map[string]interface{}", i)
	}
}

// ToSliceE casts an empty interface to a []interface{}.
func ToSliceE(i interface{}) ([]interface{}, error) {

	var s []interface{}

	switch v := i.(type) {
	case []interface{}:
		for _, u := range v {
			s = append(s, u)
		}
		return s, nil
	case []map[string]interface{}:
		for _, u := range v {
			s = append(s, u)
		}
		return s, nil
	default:
		return s, errors.NewNotValidf("[conv] Unable to cast %#v of type %v to []interface{}", i, reflect.TypeOf(i))
	}
}

// ToStringSliceE casts an empty interface to a []string.
func ToStringSliceE(i interface{}) ([]string, error) {

	var a []string

	switch v := i.(type) {
	case []interface{}:
		for _, u := range v {
			a = append(a, ToString(u))
		}
		return a, nil
	case []string:
		return v, nil
	case string:
		return strings.Fields(v), nil
	case interface{}:
		str, err := ToStringE(v)
		if err != nil {
			return a, errors.NewNotValidf("[conv] Unable to cast %#v to []string. %s", i, err)
		}
		return []string{str}, nil
	default:
		return a, errors.NewNotValidf("[conv] Unable to cast %#v to []string", i)
	}
}

// ToIntSliceE casts an empty interface to a []int.
func ToIntSliceE(i interface{}) ([]int, error) {

	if i == nil {
		return []int{}, errors.NewNotValidf("[conv] Unable to cast %#v to []int", i)
	}

	switch v := i.(type) {
	case []int:
		return v, nil
	}

	kind := reflect.TypeOf(i).Kind()
	switch kind {
	case reflect.Slice, reflect.Array:
		s := reflect.ValueOf(i)
		a := make([]int, s.Len())
		for j := 0; j < s.Len(); j++ {
			val, err := ToIntE(s.Index(j).Interface())
			if err != nil {
				return []int{}, errors.NewNotValidf("[conv] Unable to cast %#v to []int. %s", i, err)
			}
			a[j] = val
		}
		return a, nil
	default:
		return []int{}, errors.NewNotValidf("[conv] Unable to cast %#v to []int", i)
	}
}

// TimeFormats available time format to parse
var TimeFormats = [...]string{
	time.RFC3339,
	"2006-01-02T15:04:05", // iso8601 without timezone
	"2006-01-02 15:04:05", // MySQL
	time.RFC1123Z,
	time.RFC1123,
	time.RFC822Z,
	time.RFC822,
	time.ANSIC,
	time.UnixDate,
	time.RubyDate,
	"2006-01-02 15:04:05Z07:00",
	"02 Jan 06 15:04 MST",
	"2006-01-02",
	"02 Jan 2006",
	"2006-01-02 15:04:05 -07:00",
	"2006-01-02 15:04:05 -0700",
}

// StringToDate casts an empty interface to a time.Time.
// Location can be nil, then time.Local is the default value.
func StringToDate(s string, loc *time.Location) (time.Time, error) {
	if loc == nil {
		loc = time.Local
	}
	return parseDateWith(s, TimeFormats[:], loc)
}

func parseDateWith(s string, dates []string, loc *time.Location) (d time.Time, e error) {
	for _, dateType := range dates {
		if d, e = time.ParseInLocation(dateType, s, loc); e == nil {
			return
		}
	}
	return d, errors.NewNotValidf("[conv] Unable to parse date: %s", s)
}
