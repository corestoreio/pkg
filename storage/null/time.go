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
	"fmt"
	"time"

	"github.com/corestoreio/errors"
)

// TODO(cys): Remove GobEncoder, GobDecoder, MarshalJSON, UnmarshalJSON in Go 2.
// The same semantics will be provided by the generic MarshalBinary,
// MarshalText, UnmarshalBinary, UnmarshalText.

// MakeTime creates a new Time. Setting the second optional argument to
// false, the string will not be valid anymore, hence NULL. Time implements
// interface Argument.
func MakeTime(t time.Time, valid ...bool) Time {
	v := true
	if len(valid) == 1 {
		v = valid[0]
	}
	return Time{
		Time:  t,
		Valid: v,
	}
}

// String returns the string representation of the time or null.
func (nt Time) String() string {
	if !nt.Valid {
		return "null"
	}
	return nt.Time.String()
}

// GoString prints an optimized Go representation.
func (nt Time) GoString() string {
	if !nt.Valid {
		return "null.Time{}"
	}
	return fmt.Sprintf("null.MakeTime(time.Unix(%d,%d)", nt.Time.Unix(), nt.Time.Nanosecond())
}

// MarshalJSON implements json.Marshaler.
// It will encode null if this time is null.
func (nt Time) MarshalJSON() ([]byte, error) {
	if !nt.Valid {
		return bTextNullLC, nil
	}
	return nt.Time.MarshalJSON()
}

// UnmarshalJSON implements json.Unmarshaler.
// It supports string, object (e.g. pq.Time and friends)
// and null input.
func (nt *Time) UnmarshalJSON(data []byte) error {
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
			return errors.NotValid.Newf(`[dml] json: unmarshalling object into Go value of type null.Time requires key "Time" to be of type string and key "Valid" to be of type bool; found %T and %T, respectively`, x["Time"], x["Valid"])
		}
		err = nt.Time.UnmarshalText([]byte(ti))
		nt.Valid = valid
		return err
	case nil:
		nt.Valid = false
		return nil
	default:
		err = errors.NotValid.Newf("[dml] json: cannot unmarshal %#v into Go value of type null.Time", v)
	}
	nt.Valid = err == nil
	return err
}

// MarshalText transforms the time type into a byte slice.
func (nt Time) MarshalText() ([]byte, error) {
	if !nt.Valid {
		return []byte(sqlStrNullLC), nil
	}
	return nt.Time.MarshalText()
}

// UnmarshalText parses the byte slice to create a time type.
func (nt *Time) UnmarshalText(text []byte) error {
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
func (nt Time) MarshalBinary() (data []byte, err error) {
	if !nt.Valid {
		return data, nil
	}
	return nt.Time.MarshalBinary()
}

// GobEncode implements the gob.GobEncoder interface for gob serialization.
func (nt Time) GobEncode() ([]byte, error) {
	if !nt.Valid {
		return nil, nil
	}
	return nt.Time.GobEncode()
}

// GobDecode implements the gob.GobDecoder interface for gob serialization.
func (nt *Time) GobDecode(data []byte) error {
	if len(data) == 0 {
		nt.Valid = false
		return nil
	}
	return nt.Time.GobDecode(data)
}

// UnmarshalBinary parses the byte slice to create a time type.
func (nt *Time) UnmarshalBinary(data []byte) error {
	if len(data) == 0 {
		nt.Valid = false
		return nil
	}
	err := nt.Time.UnmarshalBinary(data)
	nt.Valid = err == nil
	return err
}

// SetValid changes this Time's value and sets it to be non-null.
func (nt *Time) SetValid(v time.Time) *Time {
	nt.Time = v
	nt.Valid = true
	return nt
}

// Ptr returns a pointer to this Time's value, or a nil pointer if this Time is
// null.
func (nt Time) Ptr() *time.Time {
	if !nt.Valid {
		return nil
	}
	return &nt.Time
}

// Marshal binary encoder for protocol buffers. Implements proto.Marshaler.
func (nt Time) Marshal() ([]byte, error) {
	return nt.MarshalBinary()
}

// MarshalTo binary encoder for protocol buffers which writes into data.
func (nt Time) MarshalTo(data []byte) (int, error) {
	if !nt.Valid {
		return 0, nil
	}
	raw, err := nt.Time.MarshalBinary()
	return copy(data, raw), err
}

// Unmarshal binary decoder for protocol buffers. Implements proto.Unmarshaler.
func (nt *Time) Unmarshal(data []byte) error {
	return nt.UnmarshalBinary(data)
}

// Size returns the size of the underlying type. If not valid, the size will be
// 0. Implements proto.Sizer.
func (nt Time) Size() (n int) {
	if !nt.Valid {
		return 0
	}
	secs := nt.Time.Unix()
	nano := nt.Time.Nanosecond()
	if secs != 0 {
		n += 1 + uint64Size(uint64(secs))
	}
	if nano != 0 {
		n += 1 + uint64Size(uint64(nano))
	}
	return n
}

// WriteTo uses a special dialect to encode the value and write it into w. w
// cannot be replaced by io.Writer and shall not be replaced by an interface
// because of inlining features of the compiler.
func (nt Time) WriteTo(d Dialecter, w *bytes.Buffer) (err error) {
	if nt.Valid {
		d.EscapeTime(w, nt.Time)
	} else {
		_, err = w.WriteString(sqlStrNullUC)
	}
	return
}

// Append appends the value or its nil type to the interface slice.
func (nt Time) Append(args []interface{}) []interface{} {
	if nt.Valid {
		return append(args, nt.Time)
	}
	return append(args, nil)
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
