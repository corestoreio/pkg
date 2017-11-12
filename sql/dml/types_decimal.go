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
	"bytes"
	"database/sql/driver"
	"encoding/binary"
	"math"
	"strconv"
	"strings"

	"github.com/corestoreio/pkg/util/bufferpool"
	"github.com/corestoreio/pkg/util/byteconv"
	"github.com/corestoreio/errors"
)

// TODO(cys): Remove GobEncoder, GobDecoder, MarshalJSON, UnmarshalJSON in Go 2.
// The same semantics will be provided by the generic MarshalBinary,
// MarshalText, UnmarshalBinary, UnmarshalText.

// Flags get binary encoded in the marshalers
const (
	decimalFlagNegative = 1 << iota
	decimalFlagValid
	decimalFlagQuote
	decimalBinaryVersion01
)

// Decimal defines a container type for any MySQL/MariaDB
// decimal/numeric/float/double data type and their representation in Go.
// Decimal does not perform any kind of calculations. Helpful packages for
// arbitrary precision calculations are github.com/ericlagergren/decimal or
// gopkg.in/inf.v0 or github.com/shopspring/decimal or a future new Go type.
// https://dev.mysql.com/doc/refman/5.7/en/precision-math-decimal-characteristics.html
// https://dev.mysql.com/doc/refman/5.7/en/floating-point-types.html
type Decimal struct {
	Precision uint64
	Scale     int32
	Negative  bool
	Valid     bool
	// Quote if true JSON marshaling will quote the returned number and creates
	// a string. JavaScript floats are only 53 bits.
	Quote bool
}

// MakeDecimalInt64 converts an int64 with the scale to a Decimal.
func MakeDecimalInt64(value int64, scale int32) Decimal {
	neg := false
	if value < 0 {
		neg = true
		value *= -1
	}
	return Decimal{
		Precision: uint64(value),
		Scale:     scale,
		Negative:  neg,
		Valid:     true,
	}
}

// MakeDecimalFloat64 converts a float64 to Decimal.
//
// Example:
//
//     MakeDecimalFloat64(123.45678901234567).String()  // output: "123.45678901234567"
//     MakeDecimalFloat64(123.456789012345678).String() // output: "123.45678901234568"
//     MakeDecimalFloat64(.00000000000000001).String() // output: "0.00000000000000001"
//
func MakeDecimalFloat64(value float64) (d Decimal, err error) {
	floor := math.Floor(value)

	// fast path, where float is an int
	if floor == value && value <= math.MaxInt64 && value >= math.MinInt64 {
		return MakeDecimalInt64(int64(value), 0), nil
	}

	// slow path
	// this the slow hacky way for now because the logic to
	// convert a base-2 float to base-10 properly is not trivial
	str := strconv.AppendFloat([]byte(``), value, 'f', -1, 64)
	return MakeDecimalBytes(str)
}

// MakeDecimalBytes parses b to create a new Decimal. b must contain ASCII
// numbers. If b contains null/NULL the returned object represents that value.
// This function can be used for big.Int.Bytes() or similar implementations.
func MakeDecimalBytes(b []byte) (d Decimal, err error) {
	// maybe use string comparison but run benchmarks
	if len(b) == 0 || bytes.Equal(b, bTextNullLC) || bytes.Equal(b, bTextNullUC) {
		return
	}

	d.Valid = true
	d.Negative = b[0] == '-'
	if d.Negative || b[0] == '+' {
		b = b[1:]
	}

	digits := b
	if dotPos := bytes.IndexByte(digits, '.'); dotPos > 0 { // 0.333 dotPos is min 1
		d.Scale = int32(len(b)-dotPos) - 1
		// remove dot 2363.7800 => 23637800 => Scale=4
		digits = append(digits[:dotPos], b[dotPos+1:]...)
	}
	d.Precision, err = byteconv.ParseUint(digits, 10, 64)
	return
}

// Int64 converts the underlying uint64 to an int64. Very useful for creating a
// new 3rd party package type/object. If the Precision field overflows
// math.MaxInt64 the return values are 0,0. If you want to aovid this use the
// String function and create the 3rd party type via the string.
func (d Decimal) Int64() (value int64, scale int32) {
	if d.Precision > math.MaxInt64 {
		return 0, 0 // Better solution instead of panicking?
	}
	value = int64(d.Precision)
	scale = d.Scale
	if d.Negative {
		value *= -1
	}
	return value, scale
}

// Float64 converts the precision and the scale to a float64 value including the
// usual float behaviour. Overflow will result in an undefined float64.
func (d Decimal) Float64() (value float64) {
	value = float64(d.Precision)
	value *= math.Pow10(-int(d.Scale))
	if d.Negative {
		value *= -1
	}
	return value
}

// String returns the string representation of the fixed with decimal. Returns
// the word `NULL` if the current value is not valid, for now.
func (d Decimal) String() string {
	buf := bufferpool.Get()
	defer bufferpool.Put(buf)
	d.string(buf)
	return buf.String()
}

//when needed
//func (d Decimal) AppendString(b []byte) []byte {
//	buf := bytes.NewBuffer(b)
//	d.string(buf)
//	return buf.Bytes()
//}

func (d Decimal) string(buf *bytes.Buffer) {
	if !d.Valid {
		buf.WriteString(sqlStrNullUC)
		return
	}
	prevLen := int32(buf.Len())
	if d.Negative {
		buf.WriteByte('-')
	}

	if d.Scale == 0 {
		raw := strconv.AppendUint(buf.Bytes(), d.Precision, 10)
		buf.Reset()
		buf.Write(raw)
		return
	}

	digits := int32(math.Log10(float64(d.Precision)) + 1)
	leadingZeros := d.Scale - digits + 1

	if leadingZeros > 0 {
		const zeroLen = 128 // zeros
		const zeros = "00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"
		if leadingZeros >= zeroLen {
			// slow path
			buf.WriteString(strings.Repeat("0", int(leadingZeros)))
		} else {
			buf.WriteString(zeros[:leadingZeros])
		}
		digits += leadingZeros
	}

	raw := strconv.AppendUint(buf.Bytes(), d.Precision, 10)
	buf.Reset()
	buf.Write(raw)

	pos := digits - d.Scale + prevLen
	if d.Negative {
		pos++
	}
	raw = buf.Bytes()
	newRaw := append(raw[:pos], append([]byte("."), raw[pos:]...)...)
	buf.Reset()
	buf.Write(newRaw)
}

// GoString returns an optimized version of the Go representation of Decimal.
func (d Decimal) GoString() string {
	if !d.Valid {
		return "dml.Decimal{}"
	}
	buf := bufferpool.Get()
	defer bufferpool.Put(buf)
	buf.WriteString("dml.Decimal{")
	if d.Precision > 0 {
		buf.WriteString("Precision:")
		buf2 := strconv.AppendUint(buf.Bytes(), d.Precision, 10)
		buf.Reset()
		buf.Write(buf2)
		buf.WriteByte(',')
	}
	if d.Scale != 0 {
		buf.WriteString("Scale:")
		buf2 := strconv.AppendInt(buf.Bytes(), int64(d.Scale), 10)
		buf.Reset()
		buf.Write(buf2)
		buf.WriteByte(',')
	}
	if d.Negative {
		writeLabeledBool(buf, "Negative")
	}
	if d.Valid {
		writeLabeledBool(buf, "Valid")
	}
	if d.Quote {
		writeLabeledBool(buf, "Quote")
	}
	buf.WriteByte('}')
	return buf.String()
}

func writeLabeledBool(buf *bytes.Buffer, label string) {
	buf.WriteString(label)
	buf.WriteString(":true,")
}

func unquoteIfQuoted(b []byte) (_ []byte, isQuoted bool) {
	// If the amount is quoted, strip the quotes
	if lb := len(b); lb > 2 && b[0] == '"' && b[lb-1] == '"' {
		b = b[1 : lb-1]
		isQuoted = true
	}
	return b, isQuoted
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (d *Decimal) UnmarshalJSON(b []byte) (err error) {
	b, isQuoted := unquoteIfQuoted(b)
	if *d, err = MakeDecimalBytes(b); err != nil {
		return errors.NewNotValidf("[dml] Decimal Decoding failed of %q with error: %s", b, err)
	}
	d.Quote = isQuoted
	return nil
}

// MarshalJSON implements the json.Marshaler interface.
func (d Decimal) MarshalJSON() ([]byte, error) {
	if !d.Valid {
		return []byte(sqlStrNullLC), nil
	}
	buf := new(bytes.Buffer)
	if d.Quote {
		buf.WriteByte('"')
	}
	d.string(buf)
	if d.Quote {
		buf.WriteByte('"')
	}

	return buf.Bytes(), nil
}

// UnmarshalBinary implements the encoding.BinaryUnmarshaler interface. As a string representation
// is already used when encoding to text, this method stores that string as []byte
func (d *Decimal) UnmarshalBinary(data []byte) error {
	const validLength = 14
	ld := len(data)
	if ld == 0 {
		*d = Decimal{}
		return nil
	}
	if ld != validLength {
		return errors.NewNotValidf("[dml] Decimal.UnmarshalBinary Invalid length of input data. Should be %d but have %d", validLength, len(data))
	}
	d.Precision = uint64(binary.BigEndian.Uint64(data[:8]))
	d.Scale = int32(binary.BigEndian.Uint32(data[8:12]))
	flags := uint16(binary.BigEndian.Uint16(data[12:14]))

	if flags&decimalFlagNegative != 0 {
		d.Negative = true
	}
	if flags&decimalFlagValid != 0 {
		d.Valid = true
	}
	if flags&decimalFlagQuote != 0 {
		d.Quote = true
	}
	if flags&decimalBinaryVersion01 == 0 {
		return errors.NewNotValidf("[dml] Decimal.UnmarshalBinary invalid binary version")
	}
	return nil
}

// MarshalBinary implements the encoding.BinaryMarshaler interface.
func (d Decimal) MarshalBinary() (data []byte, err error) {
	if !d.Valid {
		return nil, nil
	}
	var v0 [14]byte
	_, err = d.MarshalTo(v0[:])
	return v0[:], err
}

// Value implements the driver.Valuer interface for database serialization. It
// stores a string in driver.Value.
func (d Decimal) Value() (driver.Value, error) {
	return d.String(), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface for XML
// deserialization.
func (d *Decimal) UnmarshalText(text []byte) (err error) {
	if *d, err = MakeDecimalBytes(text); err != nil {
		return errors.NewNotValidf("[dml] Decimal Decoding failed of %q with error: %s", text, err)
	}
	return nil
}

// MarshalText implements the encoding.TextMarshaler interface for XML
// serialization. Does not support quoting. An invalid type returns an empty
// string.
func (d Decimal) MarshalText() (text []byte, err error) {
	buf := new(bytes.Buffer)
	d.string(buf)
	return buf.Bytes(), nil
}

// GobEncode implements the gob.GobEncoder interface for gob serialization.
func (d Decimal) GobEncode() ([]byte, error) {
	return d.MarshalBinary()
}

// GobDecode implements the gob.GobDecoder interface for gob serialization.
func (d *Decimal) GobDecode(data []byte) error {
	return d.UnmarshalBinary(data)
}

// Marshal binary encoder for protocol buffers. Implements proto.Marshaler.
func (d Decimal) Marshal() ([]byte, error) {
	return d.MarshalBinary()
}

// Marshal binary encoder for protocol buffers which writes into data.
func (d Decimal) MarshalTo(data []byte) (n int, err error) {
	if !d.Valid {
		return 0, nil
	}

	binary.BigEndian.PutUint64(data[:8], d.Precision)

	binary.BigEndian.PutUint32(data[8:12], uint32(d.Scale))

	var flags uint16
	flags |= decimalBinaryVersion01
	if d.Negative {
		flags |= decimalFlagNegative
	}
	if d.Valid {
		flags |= decimalFlagValid
	}
	if d.Quote {
		flags |= decimalFlagQuote
	}

	binary.BigEndian.PutUint16(data[12:14], flags)

	return 14, nil
}

// Unmarshal binary decoder for protocol buffers. Implements proto.Unmarshaler.
func (d *Decimal) Unmarshal(data []byte) error {
	return d.UnmarshalBinary(data)
}

// Size returns the size of the underlying type. If not valid, the size will be
// 0. Implements proto.Sizer.
func (d Decimal) Size() (s int) {
	if d.Valid {
		s = 14
	}
	return
}
