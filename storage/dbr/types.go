package dbr

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"time"

	"fmt"
	"github.com/corestoreio/csfw/utils/log"
	"github.com/go-sql-driver/mysql"
	"github.com/ugorji/go/codec"
	"math"
)

//
// Your app can use these Null types instead of the defaults. The sole benefit you get is a MarshalJSON method that is not retarded.
//

var (
	_          codec.Selfer = (*NullString)(nil)
	nullString              = []byte("null")
)

// NullString is a type that can be null or a string
type NullString struct {
	sql.NullString
}

// NullFloat64 is a type that can be null or a float64
type NullFloat64 struct {
	sql.NullFloat64
}

// NullInt64 is a type that can be null or an int
type NullInt64 struct {
	sql.NullInt64
}

// NullTime is a type that can be null or a time
type NullTime struct {
	mysql.NullTime
}

// NullBool is a type that can be null or a bool
type NullBool struct {
	sql.NullBool
}

// NewNullString generates a new non-pointer type. Valid argument is optional
// and will be detected automatically if left off. If value is empty, valid is
// false which means database value is NULL.
func NewNullString(value string, valid ...bool) NullString {
	ok := value != ""
	if len(valid) > 0 && value == "" {
		ok = valid[0]
	}
	return NullString{
		sql.NullString{
			String: value,
			Valid:  ok,
		},
	}
}

// GoString satisfies the interface fmt.GoStringer when using %#v in Printf methods.
// Returns
// 		dbr.NewNullString(`...`,bool)
func (ns NullString) GoString() string {
	// @todo fix bug to escape back ticks properly
	return fmt.Sprintf("dbr.NewNullString(`%s`, %t)", ns.String, ns.Valid)
}

// CodecEncodeSelf for ugorji.go codec package
func (n *NullString) CodecEncodeSelf(e *codec.Encoder) {
	if err := e.Encode(n.String); err != nil {
		log.Error("dbr.NullString.CodecEncodeSelf", "err", err, "n", n)
	}
}

// CodecDecodeSelf  for ugorji.go codec package @todo write test ... not sure if ok
func (n *NullString) CodecDecodeSelf(d *codec.Decoder) {
	if err := d.Decode(&n.String); err != nil {
		log.Error("dbr.NullString.CodecEncodeSelf", "err", err, "n", n)
	}
	// think about empty string and Valid value ...
}

// MarshalJSON correctly serializes a NullString to JSON
func (n *NullString) MarshalJSON() ([]byte, error) {
	if n.Valid {
		j, e := json.Marshal(n.String)
		return j, e
	}
	return nullString, nil
}

// UnmarshalJSON correctly deserializes a NullString from JSON
func (n *NullString) UnmarshalJSON(b []byte) error {
	var s interface{}
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	return n.Scan(s)
}

// NewNullInt64 generates a new non-pointer type. Valid argument is optional
// and will be detected automatically if left off. If value is 0, valid is
// false which means database value is NULL.
func NewNullInt64(value int64, valid ...bool) NullInt64 {
	ok := value != 0
	if len(valid) > 0 && value == 0 {
		ok = valid[0]
	}
	return NullInt64{
		sql.NullInt64{
			Int64: value,
			Valid: ok,
		},
	}
}

// MarshalJSON correctly serializes a NullInt64 to JSON
func (n *NullInt64) MarshalJSON() ([]byte, error) {
	if n.Valid {
		j, e := json.Marshal(n.Int64)
		return j, e
	}
	return nullString, nil
}

// CodecEncodeSelf for ugorji.go codec package
func (n *NullInt64) CodecEncodeSelf(e *codec.Encoder) {
	if err := e.Encode(n.Int64); err != nil {
		log.Error("dbr.NullInt64.CodecEncodeSelf", "err", err, "n", n)
	}
}

// CodecDecodeSelf  for ugorji.go codec package @todo write test ... not sure if ok
func (n *NullInt64) CodecDecodeSelf(d *codec.Decoder) {
	if err := d.Decode(&n.Int64); err != nil {
		log.Error("dbr.NullInt64.CodecEncodeSelf", "err", err, "n", n)
	}
}

// UnmarshalJSON correctly deserializes a NullInt64 from JSON
func (n *NullInt64) UnmarshalJSON(b []byte) error {
	var s interface{}
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	return n.Scan(s)
}

// NewNullFloat64 generates a new non-pointer type. Valid argument is optional
// and will be detected automatically if left off. If value is 0, valid is
// false which means database value is NULL.
func NewNullFloat64(value float64, valid ...bool) NullFloat64 {
	ok := math.Abs(value) > 0.000000001
	if len(valid) > 0 && math.Abs(value) < 0.000000001 {
		ok = valid[0]
	}
	return NullFloat64{
		sql.NullFloat64{
			Float64: value,
			Valid:   ok,
		},
	}
}

// MarshalJSON correctly serializes a NullFloat64 to JSON
func (n *NullFloat64) MarshalJSON() ([]byte, error) {
	if n.Valid {
		j, e := json.Marshal(n.Float64)
		return j, e
	}
	return nullString, nil
}

// CodecEncodeSelf for ugorji.go codec package
func (n *NullFloat64) CodecEncodeSelf(e *codec.Encoder) {
	if err := e.Encode(n.Float64); err != nil {
		log.Error("dbr.NullFloat64.CodecEncodeSelf", "err", err, "n", n)
	}
}

// CodecDecodeSelf  for ugorji.go codec package @todo write test ... not sure if ok
func (n *NullFloat64) CodecDecodeSelf(d *codec.Decoder) {
	if err := d.Decode(&n.Float64); err != nil {
		log.Error("dbr.NullFloat64.CodecEncodeSelf", "err", err, "n", n)
	}
}

// UnmarshalJSON correctly deserializes a NullFloat64 from JSON
func (n *NullFloat64) UnmarshalJSON(b []byte) error {
	var s interface{}
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	return n.Scan(s)
}

// NewNullTime generates a new non-pointer type. Valid argument is optional
// and will be detected automatically if left off. If value is 0, valid is
// false which means database value is NULL.
func NewNullTime(value time.Time, valid ...bool) NullTime {
	ok := false == value.IsZero()
	if len(valid) > 0 && value.IsZero() {
		ok = valid[0]
	}
	return NullTime{
		mysql.NullTime{
			Time:  value,
			Valid: ok,
		},
	}
}

// MarshalJSON correctly serializes a NullTime to JSON
func (n *NullTime) MarshalJSON() ([]byte, error) {
	if n.Valid {
		j, e := json.Marshal(n.Time)
		return j, e
	}
	return nullString, nil
}

// CodecEncodeSelf for ugorji.go codec package
func (n *NullTime) CodecEncodeSelf(e *codec.Encoder) {
	if err := e.Encode(n.Time); err != nil {
		log.Error("dbr.NullTime.CodecEncodeSelf", "err", err, "n", n)
	}
}

// CodecDecodeSelf  for ugorji.go codec package @todo write test ... not sure if ok
func (n *NullTime) CodecDecodeSelf(d *codec.Decoder) {
	if err := d.Decode(&n.Time); err != nil {
		log.Error("dbr.NullTime.CodecEncodeSelf", "err", err, "n", n)
	}
}

// UnmarshalJSON correctly deserializes a NullTime from JSON
func (n *NullTime) UnmarshalJSON(b []byte) error {
	// scan for null
	if bytes.Equal(b, nullString) {
		return n.Scan(nil)
	}
	// scan for JSON timestamp
	var t time.Time
	if err := json.Unmarshal(b, &t); err != nil {
		return err
	}
	return n.Scan(t)
}

// NewNullBool generates a new non-pointer type. To allow NULL values pass a
// false to the valid argument.
func NewNullBool(value bool, valid bool) NullBool {
	if value {
		valid = true
	}
	return NullBool{
		sql.NullBool{
			Bool:  value,
			Valid: valid,
		},
	}
}

// MarshalJSON correctly serializes a NullBool to JSON
func (n *NullBool) MarshalJSON() ([]byte, error) {
	if n.Valid {
		j, e := json.Marshal(n.Bool)
		return j, e
	}
	return nullString, nil
}

// CodecEncodeSelf for ugorji.go codec package
func (n *NullBool) CodecEncodeSelf(e *codec.Encoder) {
	if err := e.Encode(n.Bool); err != nil {
		log.Error("dbr.NullBool.CodecEncodeSelf", "err", err, "n", n)
	}
}

// CodecDecodeSelf  for ugorji.go codec package @todo write test ... not sure if ok
func (n *NullBool) CodecDecodeSelf(d *codec.Decoder) {
	if err := d.Decode(&n.Bool); err != nil {
		log.Error("dbr.NullBool.CodecEncodeSelf", "err", err, "n", n)
	}
}

// UnmarshalJSON correctly deserializes a NullBool from JSON
func (n *NullBool) UnmarshalJSON(b []byte) error {
	var s interface{}
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	return n.Scan(s)
}
