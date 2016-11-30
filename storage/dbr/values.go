// Copyright 2015-2016, Cyrill @ Schumacher.fm and the CoreStore contributors
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

// initial idea but .. overhead ...

//import (
//	"database/sql/driver"
//	"github.com/corestoreio/csfw/util/errors"
//	"github.com/corestoreio/csfw/util/null"
//	"time"
//)
//
//type Value struct {
//	col     string
//	int64   null.Int64
//	float64 null.Float64
//	bool    null.Bool
//	bytes   null.Bytes
//	string  null.String
//	time    null.Time
//	isNull  bool
//	// maybe add iface interface{} ;-)
//	error
//}
//
//// Map converts a map into a correct value slice. Key is the column name and
//// value one of the six supported data types.
//func Map(m map[string]interface{}) Values {
//	vals := make(Values, len(m))
//	i := 0
//	for k, v := range m {
//		vals[i] = getVal(k, v)
//		i++
//	}
//	return vals
//}
//
//func getVal(column string, v interface{}) (val Value) {
//	val.col = column
//	switch v.(type) {
//	case int64:
//		val.int64 = null.Int64From(v.(int64))
//	case float64:
//		val.float64 = null.Float64From(v.(float64))
//	case bool:
//		val.bool = null.BoolFrom(v.(bool))
//	case []byte:
//		val.bytes = null.BytesFrom(v.([]byte))
//	case string:
//		val.string = null.StringFrom(v.(string))
//	case time.Time:
//		val.time = null.TimeFrom(v.(time.Time))
//	case *time.Time:
//		val.time = null.TimeFromPtr(v.(*time.Time))
//	default:
//		val.error = errors.NewNotSupportedf("[dbr] For column %q the type %#v is not yet supported", column, v)
//	}
//	return val
//}
//
//// DriverValue ...
//// It is either nil or an instance of one of these types:
////
////   int64
////   float64
////   bool
////   []byte
////   string
////   time.Time
//func DriverValue(column string, dv driver.Valuer) Value {
//
//	v, err := dv.Value()
//	if err != nil {
//		return Value{
//			col:   column,
//			error: errors.NewFatalf("[dbr] For column %q driver.Valuer return error %s", err),
//		}
//	}
//
//	if v == nil {
//		return Value{
//			col:    column,
//			isNull: true,
//		}
//	}
//
//	return getVal(column, v)
//}
//
//func Int64(column string, i int64) Value {
//	return Value{
//		col:   column,
//		int64: null.Int64From(i),
//	}
//}
//
//func Float64(column string, f float64) Value {
//	return Value{
//		col:     column,
//		float64: null.Float64From(f),
//	}
//}
//
//func String(column string, s string) Value {
//	return Value{
//		col:    column,
//		string: null.StringFrom(s),
//	}
//}
//
//func Bool(column string, b bool) Value {
//	return Value{
//		col:  column,
//		bool: null.BoolFrom(b),
//	}
//}
//
//func Time(column string, t time.Time) Value {
//	return Value{
//		col:  column,
//		time: null.TimeFrom(t),
//	}
//}
//
//func Bytes(column string, b []byte) Value {
//	return Value{
//		col:   column,
//		bytes: null.BytesFrom(b),
//	}
//}
//
//func Null(column string) Value {
//	return Value{
//		col:    column,
//		isNull: true,
//	}
//}
//
//type Values []Value
//
////func (vs Values) HasError() bool {
////	for _, v := range vs {
////		if v.error != nil {
////			return true
////		}
////	}
////	return false
////}
////
////func (vs Values) Error() string {
////	for _, v := range vs {
////		if v.error != nil {
////			return fmt.Sprintf("%+v", v.error)
////		}
////	}
////	return ""
////}
//
//func (vs Values) Arguments() ([]interface{}, error) {
//	for _, v := range vs {
//		if v.error != nil {
//			return nil, errors.Wrap(v.error, "[dbr] Values.Arguments")
//		}
//	}
//
//	args := make([]interface{}, len(vs))
//	for i, v := range vs {
//		switch {
//		case v.isNull:
//			args[i] = nil
//		case v.int64.Valid:
//			args[i] = v.int64.Int64
//		case v.float64.Valid:
//			args[i] = v.float64.Float64
//		case v.bool.Valid:
//			args[i] = v.bool.Bool
//		case v.bytes.Valid:
//			args[i] = v.bytes.Bytes
//		case v.string.Valid:
//			args[i] = v.string.String
//		case v.time.Valid:
//			args[i] = v.time.Time
//		default:
//			return nil, errors.NewEmptyf("[dbr] Provided value is empty: Index %d => %#v", i, v)
//		}
//	}
//	return args, nil
//}
