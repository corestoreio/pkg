// Copyright 2015-present, Cyrill @ Schumacher.fm and the CoreStore contributors
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

package null

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/corestoreio/errors"
)

// TODO(cys): Remove GobEncoder, GobDecoder, MarshalJSON, UnmarshalJSON in Go 2.
// The same semantics will be provided by the generic MarshalBinary,
// MarshalText, UnmarshalBinary, UnmarshalText.

// MakeTime creates a new Time. Setting the second optional argument to
// false, the string will not be valid anymore, hence NULL. Time implements
// interface Argument.
func MakeTime(t time.Time) Time {
	return Time{
		NullTime: sql.NullTime{Time: t, Valid: true},
	}
}

// String returns the string representation of the time or null.
func (a Time) String() string {
	if !a.Valid {
		return "null"
	}
	return a.Time.String()
}

// GoString prints an optimized Go representation.
func (a Time) GoString() string {
	if !a.Valid {
		return "null.Time{}"
	}
	return fmt.Sprintf("null.MakeTime(time.Unix(%d,%d)", a.Time.Unix(), a.Time.Nanosecond())
}

// MarshalJSON implements json.Marshaler.
// It will encode null if this time is null.
func (a Time) MarshalJSON() ([]byte, error) {
	if !a.Valid {
		return bTextNullLC, nil
	}
	return a.Time.MarshalJSON()
}

// UnmarshalJSON implements json.Unmarshaler.
// It supports string, object (e.g. pq.Time and friends)
// and null input.
func (a *Time) UnmarshalJSON(data []byte) error {
	if len(data) == 0 || bytes.Equal(bTextNullLC, data) {
		a.Valid = false
		a.Time = time.Time{}
		return nil
	}
	var err error
	var v interface{}
	if err = json.Unmarshal(data, &v); err != nil {
		return err
	}
	switch x := v.(type) {
	case string:
		err = a.Time.UnmarshalJSON(data)
	case map[string]interface{}:
		ti, tiOK := x["Time"].(string)
		valid, validOK := x["Valid"].(bool)
		if !tiOK || !validOK {
			return errors.NotValid.Newf(`[dml] json: unmarshalling object into Go value of type null.Time requires key "Time" to be of type string and key "Valid" to be of type bool; found %T and %T, respectively`, x["Time"], x["Valid"])
		}
		err = a.Time.UnmarshalText([]byte(ti))
		a.Valid = valid
		return err
	case nil:
		a.Valid = false
		return nil
	default:
		err = errors.NotValid.Newf("[dml] json: cannot unmarshal %#v into Go value of type null.Time", v)
	}
	a.Valid = err == nil
	return err
}

// MarshalText transforms the time type into a byte slice.
func (a Time) MarshalText() ([]byte, error) {
	if !a.Valid {
		return []byte(sqlStrNullLC), nil
	}
	return a.Time.MarshalText()
}

// UnmarshalText parses the byte slice to create a time type.
func (a *Time) UnmarshalText(text []byte) error {
	if len(text) == 0 || bytes.Equal(bTextNullLC, text) {
		a.Valid = false
		a.Time = time.Time{}
		return nil
	}
	if err := a.Time.UnmarshalText(text); err != nil {
		return err
	}
	a.Valid = true
	return nil
}

// MarshalBinary transforms the time type into a byte slice.
func (a Time) MarshalBinary() (data []byte, err error) {
	if !a.Valid {
		return data, nil
	}
	return a.NullTime.Time.MarshalBinary()
}

// UnmarshalBinary parses the byte slice to create a time type.
func (a *Time) UnmarshalBinary(data []byte) error {
	if len(data) == 0 {
		a.Valid = false
		return nil
	}
	err := a.Time.UnmarshalBinary(data)
	a.Valid = err == nil
	return err
}

// SetValid changes this Time's value and sets it to be non-null.
func (a *Time) SetValid(v time.Time) { a.Time = v; a.Valid = true }

// Reset sets the value to Go's default value and Valid to false.
func (a *Time) Reset() { *a = Time{} }

// Ptr returns a pointer to this Time's value, or a nil pointer if this Time is
// null.
func (a Time) Ptr() *time.Time {
	if !a.Valid {
		return nil
	}
	return &a.Time
}

// SetPtr sets v according to the rules.
func (a *Time) SetPtr(v *time.Time) {
	a.Valid = v != nil
	a.Time = time.Time{}
	if v != nil {
		a.Time = *v
	}
}

func (a *Time) SetProto(v *timestamppb.Timestamp) {
	a.Valid = v != nil
	a.Time = time.Time{}
	if v != nil {
		a.Time = v.AsTime()
	}
}

func (a *Time) Proto() *timestamppb.Timestamp {
	if a == nil || !a.Valid {
		return nil
	}
	return timestamppb.New(a.Time)
}

// WriteTo uses a special dialect to encode the value and write it into w. w
// cannot be replaced by io.Writer and shall not be replaced by an interface
// because of inlining features of the compiler.
func (a Time) WriteTo(d Dialecter, w *bytes.Buffer) (err error) {
	if a.Valid {
		d.EscapeTime(w, a.Time)
	} else {
		_, err = w.WriteString(sqlStrNullUC)
	}
	return
}

// Append appends the value or its nil type to the interface slice.
func (a Time) Append(args []interface{}) []interface{} {
	if a.Valid {
		return append(args, a.Time)
	}
	return append(args, nil)
}
