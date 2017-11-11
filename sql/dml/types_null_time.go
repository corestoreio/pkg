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

package dml

import (
	"fmt"
	"time"

	"github.com/corestoreio/errors"
)

// TODO(cys): Remove GobEncoder, GobDecoder, MarshalJSON, UnmarshalJSON in Go 2.
// The same semantics will be provided by the generic MarshalBinary,
// MarshalText, UnmarshalBinary, UnmarshalText.

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

// String returns the string representation of the time or null.
func (nt NullTime) String() string {
	if !nt.Valid {
		return "null"
	}
	return nt.Time.String()
}

// GoString prints an optimized Go representation.
func (nt NullTime) GoString() string {
	if !nt.Valid {
		return "dml.NullTime{}"
	}
	return fmt.Sprintf("dml.MakeNullTime(time.Unix(%d,%d)", nt.Time.Unix(), nt.Time.Nanosecond())
}

// MarshalJSON implements json.Marshaler.
// It will encode null if this time is null.
func (nt NullTime) MarshalJSON() ([]byte, error) {
	if !nt.Valid {
		return []byte(sqlStrNullLC), nil
	}
	return nt.Time.MarshalJSON()
}

// UnmarshalJSON implements json.Unmarshaler.
// It supports string, object (e.g. pq.NullTime and friends)
// and null input.
func (nt *NullTime) UnmarshalJSON(data []byte) error {
	var err error
	var v interface{}
	if err = JSONUnMarshalFn(data, &v); err != nil {
		return err
	}
	switch x := v.(type) {
	case string:
		err = nt.Time.UnmarshalJSON(data)
	case map[string]interface{}:
		ti, tiOK := x["Time"].(string)
		valid, validOK := x["Valid"].(bool)
		if !tiOK || !validOK {
			return errors.NewNotValidf(`[dml] json: unmarshalling object into Go value of type dml.NullTime requires key "Time" to be of type string and key "Valid" to be of type bool; found %T and %T, respectively`, x["Time"], x["Valid"])
		}
		err = nt.Time.UnmarshalText([]byte(ti))
		nt.Valid = valid
		return err
	case nil:
		nt.Valid = false
		return nil
	default:
		err = errors.NewNotValidf("[dml] json: cannot unmarshal %#v into Go value of type dml.NullTime", v)
	}
	nt.Valid = err == nil
	return err
}

// MarshalText transforms the time type into a byte slice.
func (nt NullTime) MarshalText() ([]byte, error) {
	if !nt.Valid {
		return []byte(sqlStrNullLC), nil
	}
	return nt.Time.MarshalText()
}

// UnmarshalText parses the byte slice to create a time type.
func (nt *NullTime) UnmarshalText(text []byte) error {
	str := string(text)
	if str == "" || str == sqlStrNullLC {
		nt.Valid = false
		return nil
	}
	if err := nt.Time.UnmarshalText(text); err != nil {
		return err
	}
	nt.Valid = true
	return nil
}

// MarshalBinary transforms the time type into a byte slice.
func (nt NullTime) MarshalBinary() (data []byte, err error) {
	if !nt.Valid {
		return data, nil
	}
	return nt.Time.MarshalBinary()
}

// GobEncode implements the gob.GobEncoder interface for gob serialization.
func (nt NullTime) GobEncode() ([]byte, error) {
	if !nt.Valid {
		return nil, nil
	}
	return nt.Time.GobEncode()
}

// GobDecode implements the gob.GobDecoder interface for gob serialization.
func (nt *NullTime) GobDecode(data []byte) error {
	if len(data) == 0 {
		nt.Valid = false
		return nil
	}
	return nt.Time.GobDecode(data)
}

// UnmarshalBinary parses the byte slice to create a time type.
func (nt *NullTime) UnmarshalBinary(data []byte) error {
	if len(data) == 0 {
		nt.Valid = false
		return nil
	}
	err := nt.Time.UnmarshalBinary(data)
	nt.Valid = err == nil
	return err
}

// SetValid changes this Time's value and sets it to be non-null.
func (nt *NullTime) SetValid(v time.Time) *NullTime {
	nt.Time = v
	nt.Valid = true
	return nt
}

// Ptr returns a pointer to this Time's value, or a nil pointer if this Time is
// null.
func (nt NullTime) Ptr() *time.Time {
	if !nt.Valid {
		return nil
	}
	return &nt.Time
}

// Marshal binary encoder for protocol buffers. Implements proto.Marshaler.
func (a NullTime) Marshal() ([]byte, error) {
	return a.MarshalBinary()
}

// Marshal binary encoder for protocol buffers which writes into data.
func (a NullTime) MarshalTo(data []byte) (int, error) {
	if !a.Valid {
		return 0, nil
	}
	raw, err := a.Time.MarshalBinary()
	return copy(data, raw), err
}

// Unmarshal binary decoder for protocol buffers. Implements proto.Unmarshaler.
func (a *NullTime) Unmarshal(data []byte) error {
	return a.UnmarshalBinary(data)
}

// Size returns the size of the underlying type. If not valid, the size will be
// 0. Implements proto.Sizer.
func (m NullTime) Size() (n int) {
	if !m.Valid {
		return 0
	}
	secs := m.Time.Unix()
	nano := m.Time.Nanosecond()
	if secs != 0 {
		n += 1 + uint64Size(uint64(secs))
	}
	if nano != 0 {
		n += 1 + uint64Size(uint64(nano))
	}
	return n
}

func uint64Size(x uint64) (n int) {
	for {
		n++
		x >>= 7
		if x == 0 {
			break
		}
	}
	return n
}
