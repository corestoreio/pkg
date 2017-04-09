// Copyright 2015-2017, Cyrill @ Schumacher.fm and the CoreStore contributors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package dbr

import (
	"time"

	"github.com/corestoreio/errors"
)

func (a NullTime) toIFace(args *[]interface{}) {
	if a.Valid {
		*args = append(*args, a.Time)
	} else {
		*args = append(*args, nil)
	}
}

func (a NullTime) writeTo(w queryWriter, _ int) error {
	if a.Valid {
		dialect.EscapeTime(w, a.Time)
		return nil
	}
	_, err := w.WriteString("NULL")
	return err
}

func (a NullTime) len() int { return 1 }
func (a NullTime) Operator(opt byte) Argument {
	a.opt = opt
	return a
}

func (a NullTime) operator() byte { return a.opt }

// MakeNullTime creates a new NullTime. Setting the second optional argument to
// false, the string will not be valid anymore, hence NULL. NullTime implements
// interface Argument.
func MakeNullTime(t time.Time, valid ...bool) NullTime {
	v := true
	if len(valid) == 1 {
		v = valid[0]
	}
	return NullTime{
		Time:  t,
		Valid: v,
	}
}

// MarshalJSON implements json.Marshaler.
// It will encode null if this time is null.
func (t NullTime) MarshalJSON() ([]byte, error) {
	if !t.Valid {
		return []byte("null"), nil
	}
	return t.Time.MarshalJSON()
}

// UnmarshalJSON implements json.Unmarshaler.
// It supports string, object (e.g. pq.NullTime and friends)
// and null input.
func (t *NullTime) UnmarshalJSON(data []byte) error {
	var err error
	var v interface{}
	if err = JSONUnMarshalFn(data, &v); err != nil {
		return err
	}
	switch x := v.(type) {
	case string:
		err = t.Time.UnmarshalJSON(data)
	case map[string]interface{}:
		ti, tiOK := x["Time"].(string)
		valid, validOK := x["Valid"].(bool)
		if !tiOK || !validOK {
			return errors.NewNotValidf(`[dbr] json: unmarshalling object into Go value of type dbr.NullTime requires key "Time" to be of type string and key "Valid" to be of type bool; found %T and %T, respectively`, x["Time"], x["Valid"])
		}
		err = t.Time.UnmarshalText([]byte(ti))
		t.Valid = valid
		return err
	case nil:
		t.Valid = false
		return nil
	default:
		err = errors.NewNotValidf("[dbr] json: cannot unmarshal %#v into Go value of type dbr.NullTime", v)
	}
	t.Valid = err == nil
	return err
}

func (t NullTime) MarshalText() ([]byte, error) {
	if !t.Valid {
		return []byte("null"), nil
	}
	return t.Time.MarshalText()
}

func (t *NullTime) UnmarshalText(text []byte) error {
	str := string(text)
	if str == "" || str == "null" {
		t.Valid = false
		return nil
	}
	if err := t.Time.UnmarshalText(text); err != nil {
		return err
	}
	t.Valid = true
	return nil
}

// SetValid changes this Time's value and sets it to be non-null.
func (t *NullTime) SetValid(v time.Time) {
	t.Time = v
	t.Valid = true
}

// Ptr returns a pointer to this Time's value, or a nil pointer if this Time is null.
func (t NullTime) Ptr() *time.Time {
	if !t.Valid {
		return nil
	}
	return &t.Time
}

type argNullTimes struct {
	opt  byte
	data []NullTime
}

func (a argNullTimes) toIFace(args *[]interface{}) {
	for _, s := range a.data {
		if s.Valid {
			*args = append(*args, s.Time)
		} else {
			*args = append(*args, nil)
		}
	}
}

func (a argNullTimes) writeTo(w queryWriter, pos int) error {
	if a.operator() != OperatorIn && a.operator() != OperatorNotIn {
		if s := a.data[pos]; s.Valid {
			dialect.EscapeTime(w, s.Time)
			return nil
		}
		_, err := w.WriteString("NULL")
		return err
	}
	l := len(a.data) - 1
	w.WriteRune('(')
	for i, v := range a.data {
		if v.Valid {
			dialect.EscapeTime(w, v.Time)
		} else {
			w.WriteString("NULL")
		}
		if i < l {
			w.WriteRune(',')
		}
	}
	_, err := w.WriteRune(')')
	return err
}

func (a argNullTimes) len() int {
	if isNotIn(a.operator()) {
		return len(a.data)
	}
	return 1
}

func (a argNullTimes) Operator(opt byte) Argument {
	a.opt = opt
	return a
}

func (a argNullTimes) operator() byte { return a.opt }

// ArgNullTime adds a nullable Time or a slice of nullable Timess to the
// argument list. Providing no arguments returns a NULL type.
func ArgNullTime(args ...NullTime) Argument {
	if len(args) == 1 {
		return args[0]
	}
	return argNullTimes{data: args}
}
