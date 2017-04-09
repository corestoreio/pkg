package null
//
//import (
//	"encoding/json"
//	"time"
//
//	"github.com/corestoreio/errors"
//)
//
//// NewTime creates a new Time.
//func MakeTime(t time.Time, valid bool) Time {
//	return Time{
//		Time:  t,
//		Valid: valid,
//	}
//}
//
//// TimeFrom creates a new Time that will always be valid.
//func TimeFrom(t time.Time) Time {
//	return MakeTime(t, true)
//}
//
//// TimeFromPtr creates a new Time that will be null if t is nil.
//func TimeFromPtr(t *time.Time) Time {
//	if t == nil {
//		return MakeTime(time.Time{}, false)
//	}
//	return MakeTime(*t, true)
//}
//
//// MarshalJSON implements json.Marshaler.
//// It will encode null if this time is null.
//func (t Time) MarshalJSON() ([]byte, error) {
//	if !t.Valid {
//		return []byte("null"), nil
//	}
//	return t.Time.MarshalJSON()
//}
//
//// UnmarshalJSON implements json.Unmarshaler.
//// It supports string, object (e.g. pq.NullTime and friends)
//// and null input.
//func (t *Time) UnmarshalJSON(data []byte) error {
//	var err error
//	var v interface{}
//	if err = json.Unmarshal(data, &v); err != nil {
//		return err
//	}
//	switch x := v.(type) {
//	case string:
//		err = t.Time.UnmarshalJSON(data)
//	case map[string]interface{}:
//		ti, tiOK := x["Time"].(string)
//		valid, validOK := x["Valid"].(bool)
//		if !tiOK || !validOK {
//			return errors.NewNotValidf(`[null] json: unmarshalling object into Go value of type null.Time requires key "Time" to be of type string and key "Valid" to be of type bool; found %T and %T, respectively`, x["Time"], x["Valid"])
//		}
//		err = t.Time.UnmarshalText([]byte(ti))
//		t.Valid = valid
//		return err
//	case nil:
//		t.Valid = false
//		return nil
//	default:
//		err = errors.NewNotValidf("[null] json: cannot unmarshal %#v into Go value of type null.Time", v)
//	}
//	t.Valid = err == nil
//	return err
//}
//
//func (t Time) MarshalText() ([]byte, error) {
//	if !t.Valid {
//		return []byte("null"), nil
//	}
//	return t.Time.MarshalText()
//}
//
//func (t *Time) UnmarshalText(text []byte) error {
//	str := string(text)
//	if str == "" || str == "null" {
//		t.Valid = false
//		return nil
//	}
//	if err := t.Time.UnmarshalText(text); err != nil {
//		return err
//	}
//	t.Valid = true
//	return nil
//}
//
//// SetValid changes this Time's value and sets it to be non-null.
//func (t *Time) SetValid(v time.Time) {
//	t.Time = v
//	t.Valid = true
//}
//
//// Ptr returns a pointer to this Time's value, or a nil pointer if this Time is null.
//func (t Time) Ptr() *time.Time {
//	if !t.Valid {
//		return nil
//	}
//	return &t.Time
//}
