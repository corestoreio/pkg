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

package ddl

import (
	"strconv"
	"strings"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/pkg/sql/dml"
)

// Variables contains multiple MySQL configuration variables. Not threadsafe.
type Variables struct {
	Data map[string]string
	Show *dml.Show
}

// NewVariables creates a new variable collection. If the argument names gets
// passed, the SQL query will load the all variables matching the names.
// Empty argument loads all variables.
func NewVariables(names ...string) *Variables {
	vs := &Variables{
		Data: make(map[string]string),
		Show: dml.NewShow().Variable().Interpolate(),
	}
	vs.Show.IsBuildCache = true
	if len(names) > 1 {
		vs.Show.Where(dml.Column("Variable_name").In().Strs(names...))
	} else if len(names) == 1 {
		vs.Show.Where(dml.Column("Variable_name").Like().Str(names[0]))
	}
	return vs
}

// Float64 returns for a given key its float64 value. If the key does not exists
// or string parsing into float fails, it returns false.
func (vs *Variables) Float64(key string) (val float64, ok bool) {
	vals, ok := vs.Data[key]
	if !ok {
		return val, false
	}
	val, err := strconv.ParseFloat(vals, 64)
	return val, err == nil
}

// Int64 returns for a given key its int64 value. If the key does not exists
// or string parsing into int fails, it returns false.
func (vs *Variables) Int64(key string) (val int64, ok bool) {
	vals, ok := vs.Data[key]
	if !ok {
		return val, false
	}
	val, err := strconv.ParseInt(vals, 10, 64)
	return val, err == nil
}

// Uint64 returns for a given key its uint64 value. If the key does not exists
// or string parsing into uint fails, it returns false.
func (vs *Variables) Uint64(key string) (val uint64, ok bool) {
	vals, ok := vs.Data[key]
	if !ok {
		return val, false
	}
	val, err := strconv.ParseUint(vals, 10, 64)
	return val, err == nil
}

// Bool returns for a given key its bool value. If the key does not exists or
// string parsing into bool fails, it returns false. Only allowed bool values
// are YES, NO, ON, OFF and yes, no, on, off.
func (vs *Variables) Bool(key string) (val bool, ok bool) {
	vals, ok := vs.Data[key]
	if !ok {
		return val, false
	}
	switch vals {
	case "YES", "yes", "ON", "on":
		val, ok = true, true
	case "NO", "no", "OFF", "off":
		val, ok = false, true
	default:
		ok = false
	}
	return
}

// String returns for a given key its string value. If the key does not exists,
// it returns false.
func (vs *Variables) String(key string) (val string, ok bool) {
	val, ok = vs.Data[key]
	return
}

// EqualFold reports whether the value of key and `expected`, interpreted as
// UTF-8 strings, are equal under Unicode case-folding.
func (vs *Variables) EqualFold(key, expected string) bool {
	return strings.EqualFold(vs.Data[key], expected)
}

// Equal compares case sensitive the value of key with the `expected`.
func (vs *Variables) Equal(key, expected string) bool {
	return vs.Data[key] == expected
}

// Keys if argument keys has been provided the map keys will be appended to the
// slice otherwise a new slice gets returned. The returned slice has a random
// order. The keys argument is not for filtering.
func (vs *Variables) Keys(keys ...string) []string {
	if keys == nil {
		keys = make([]string, 0, len(vs.Data))
	}
	for k := range vs.Data {
		keys = append(keys, k)
	}
	return keys
}

// ToSQL implements dml.QueryBuilder interface to assemble a SQL string and its
// arguments for query execution.
func (vs *Variables) ToSQL() (string, []interface{}, error) {
	return vs.Show.ToSQL()
}

// MapColumns implements dml.ColumnMapper interface and scans a single row from
// the database query result.
func (vs *Variables) MapColumns(rc *dml.ColumnMap) error {
	name, value := "", ""
	for rc.Next() {
		switch col := rc.Column(); col {
		case "Variable_name":
			rc.String(&name)
		case "Value":
			rc.String(&value)
		default:
			return errors.NewNotFoundf("[ddl] Column %q not found in SHOW VARIABLES", col)
		}
	}
	vs.Data[name] = value
	return errors.WithStack(rc.Err())
}
