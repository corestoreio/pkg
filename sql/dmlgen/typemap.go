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

package dmlgen

import (
	"strings"
	"unicode"
	"unicode/utf8"

	"fmt"

	"github.com/corestoreio/csfw/sql/ddl"
	"github.com/corestoreio/csfw/util/strs"
	"github.com/corestoreio/errors"
)

const goTypeOptions = 8

// These constants are used as slice indexes. The constance must start with zero.
const (
	idxMysqlUnsignedNull = iota
	idxMysqlUnsignedNotNull
	idxMysqlSignedNull
	idxMysqlSignedNotNull
	idxProtobufUnsignedNull
	idxProtobufUnsignedNotNull
	idxProtobufSignedNull
	idxProtobufSignedNotNull
)

// These variables are mapping the un/signed and null/not-null types to the
// appropriate Go type.
var (
	goTypeInt64 = [...]string{
		idxMysqlUnsignedNull:    "dml.NullInt64",
		idxMysqlUnsignedNotNull: "uint64",
		idxMysqlSignedNull:      "dml.NullInt64",
		idxMysqlSignedNotNull:   "int64",

		idxProtobufUnsignedNull:    "github.com/corestoreio/csfw/sql/dml.NullInt64",
		idxProtobufUnsignedNotNull: "uint64",
		idxProtobufSignedNull:      "github.com/corestoreio/csfw/sql/dml.NullInt64",
		idxProtobufSignedNotNull:   "int64",
	}
	goTypeInt = [...]string{
		idxMysqlUnsignedNull:    "dml.NullInt64",
		idxMysqlUnsignedNotNull: "uint64",
		idxMysqlSignedNull:      "dml.NullInt64",
		idxMysqlSignedNotNull:   "int64",

		idxProtobufUnsignedNull:    "github.com/corestoreio/csfw/sql/dml.NullInt64",
		idxProtobufUnsignedNotNull: "uint32",
		idxProtobufSignedNull:      "github.com/corestoreio/csfw/sql/dml.NullInt64",
		idxProtobufSignedNotNull:   "int32",
	}
	goTypeFloat64 = [...]string{
		idxMysqlUnsignedNull:    "dml.NullFloat64",
		idxMysqlUnsignedNotNull: "float64",
		idxMysqlSignedNull:      "dml.NullFloat64",
		idxMysqlSignedNotNull:   "float64",

		idxProtobufUnsignedNull:    "github.com/corestoreio/csfw/sql/dml.NullFloat64",
		idxProtobufUnsignedNotNull: "double",
		idxProtobufSignedNull:      "github.com/corestoreio/csfw/sql/dml.NullFloat64",
		idxProtobufSignedNotNull:   "double",
	}
	goTypeTime = [...]string{
		idxMysqlUnsignedNull:    "dml.NullTime",
		idxMysqlUnsignedNotNull: "time.Time",
		idxMysqlSignedNull:      "dml.NullTime",
		idxMysqlSignedNotNull:   "time.Time",

		idxProtobufUnsignedNull:    "google.protobuf.Timestamp",
		idxProtobufUnsignedNotNull: "google.protobuf.Timestamp",
		idxProtobufSignedNull:      "google.protobuf.Timestamp",
		idxProtobufSignedNotNull:   "google.protobuf.Timestamp",
	}
	goTypeString = [...]string{
		idxMysqlUnsignedNull:    "dml.NullString",
		idxMysqlUnsignedNotNull: "string",
		idxMysqlSignedNull:      "dml.NullString",
		idxMysqlSignedNotNull:   "string",

		idxProtobufUnsignedNull:    "string",
		idxProtobufUnsignedNotNull: "string",
		idxProtobufSignedNull:      "string",
		idxProtobufSignedNotNull:   "string",
	}
	goTypeBool = [...]string{
		idxMysqlUnsignedNull:    "dml.NullBool",
		idxMysqlUnsignedNotNull: "bool",
		idxMysqlSignedNull:      "dml.NullBool",
		idxMysqlSignedNotNull:   "bool",

		idxProtobufUnsignedNull:    "github.com/corestoreio/csfw/sql/dml.NullBool",
		idxProtobufUnsignedNotNull: "bool",
		idxProtobufSignedNull:      "github.com/corestoreio/csfw/sql/dml.NullBool",
		idxProtobufSignedNotNull:   "bool",
	}
	goTypeDecimal = [...]string{
		idxMysqlUnsignedNull:    "dml.Decimal",
		idxMysqlUnsignedNotNull: "dml.Decimal",
		idxMysqlSignedNull:      "dml.Decimal",
		idxMysqlSignedNotNull:   "dml.Decimal",

		idxProtobufUnsignedNull:    "github.com/corestoreio/csfw/sql/dml.Decimal",
		idxProtobufUnsignedNotNull: "github.com/corestoreio/csfw/sql/dml.Decimal",
		idxProtobufSignedNull:      "github.com/corestoreio/csfw/sql/dml.Decimal",
		idxProtobufSignedNotNull:   "github.com/corestoreio/csfw/sql/dml.Decimal",
	}
	goTypeByte = [...]string{
		idxMysqlUnsignedNull:    "[]byte",
		idxMysqlUnsignedNotNull: "[]byte",
		idxMysqlSignedNull:      "[]byte",
		idxMysqlSignedNotNull:   "[]byte",

		idxProtobufUnsignedNull:    "bytes",
		idxProtobufUnsignedNotNull: "bytes",
		idxProtobufSignedNull:      "bytes",
		idxProtobufSignedNotNull:   "bytes",
	}
)

// mysqlTypeToGo maps the MySql/MariaDB field type to the correct Go/protobuf type.
var mysqlTypeToGo = map[string][goTypeOptions]string{
	"int":        goTypeInt,
	"bigint":     goTypeInt64,
	"smallint":   goTypeInt,
	"tinyint":    goTypeInt,
	"mediumint":  goTypeInt,
	"double":     goTypeFloat64,
	"float":      goTypeFloat64,
	"decimal":    goTypeFloat64,
	"date":       goTypeTime,
	"datetime":   goTypeTime,
	"timestamp":  goTypeTime,
	"char":       goTypeString,
	"varchar":    goTypeString,
	"enum":       goTypeString,
	"set":        goTypeString,
	"text":       goTypeString,
	"longtext":   goTypeString,
	"mediumtext": goTypeString,
	"tinytext":   goTypeString,
	"blob":       goTypeString,
	"longblob":   goTypeString,
	"mediumblob": goTypeString,
	"tinyblob":   goTypeString,
	"binary":     goTypeByte,
	"varbinary":  goTypeByte,
	"bit":        goTypeBool,
}

func toGoTypeNull(c *ddl.Column) string {
	return mySQLToGoType(c, true)
}

func toGoType(c *ddl.Column) string {
	return mySQLToGoType(c, false)
}

// mySQLToGoType calculates the data type of the field DataType. For example
// bigint, smallint, tinyint will result in "int". If withNull is true the
// returned type can store a null value.
func mySQLToGoType(c *ddl.Column, withNull bool) string {

	goType, ok := mysqlTypeToGo[c.DataType]
	if !ok {
		panic(errors.NewNotFoundf("[dmlgen] MySQL type %q not found", c.DataType))
	}

	// The switch block overwrites the already retrieved goType by checking for
	// bool columns and columns which contains a money unit.
	switch {
	case c.IsBool():
		goType = goTypeBool
	case c.IsFloat() && c.IsMoney():
		goType = goTypeDecimal
	}

	var t string
	switch {
	case c.IsUnsigned() && c.IsNull() && withNull:
		t = goType[idxMysqlUnsignedNull] // unsigned null
	case c.IsUnsigned() && (!c.IsNull() || !withNull):
		t = goType[idxMysqlUnsignedNotNull] // unsigned not null
	case !c.IsUnsigned() && c.IsNull() && withNull:
		t = goType[idxMysqlSignedNull] // signed null
	case !c.IsUnsigned() && (!c.IsNull() || !withNull):
		t = goType[idxMysqlSignedNotNull] // signed not null
	}

	return t
}

// toGoPrimitive returns for Go type or structure the final primitive:
// int->int but NullInt->.Int
func toGoPrimitive(c *ddl.Column) string {
	t := mySQLToGoType(c, true)
	field := strs.ToGoCamelCase(c.Field)
	if strings.HasPrefix(t, "dml.Null") {
		t = field + "." + t[8:]
	} else {
		t = field
	}
	return t
}

func toGoFuncNull(c *ddl.Column) string {
	return mySQLToGoFunc(c, true)
}

func toGoFunc(c *ddl.Column) string {
	return mySQLToGoFunc(c, false)
}

func mySQLToGoFunc(c *ddl.Column, withNull bool) string {

	gt := mySQLToGoType(c, withNull)
	switch gt {
	case "[]byte":
		return "Byte"
	}

	if dot := strings.IndexByte(gt, '.'); dot > 0 {
		return gt[dot+1:]
	}
	r, n := utf8.DecodeRuneInString(gt)
	return string(unicode.ToUpper(r)) + gt[n:]
}

func toProto(c *ddl.Column, withNull bool) string {

	goType, ok := mysqlTypeToGo[c.DataType]
	if !ok {
		panic(errors.NewNotFoundf("[dmlgen] MySQL type %q not found", c.DataType))
	}

	// The switch block overwrites the already retrieved goType by checking for
	// bool columns and columns which contains a money unit.
	switch {
	case c.IsBool():
		goType = goTypeBool
	case c.IsFloat() && c.IsMoney():
		goType = goTypeDecimal
	}

	var t string
	switch {
	case c.IsUnsigned() && c.IsNull() && withNull:
		t = goType[idxProtobufUnsignedNull] // unsigned null
	case c.IsUnsigned() && (!c.IsNull() || !withNull):
		t = goType[idxProtobufUnsignedNotNull] // unsigned not null
	case !c.IsUnsigned() && c.IsNull() && withNull:
		t = goType[idxProtobufSignedNull] // signed null
	case !c.IsUnsigned() && (!c.IsNull() || !withNull):
		t = goType[idxProtobufSignedNotNull] // signed not null
	}

	return t
}

func toProtoType(c *ddl.Column) string {
	pt := toProto(c, true)
	if strings.IndexByte(pt, '/') > 0 { // slash identifies an import path
		return "bytes"
	}
	return pt
}

func toProtoCustomType(c *ddl.Column) string {
	pt := toProto(c, true)
	var buf strings.Builder
	if strings.IndexByte(pt, '/') > 0 { // slash identifies an import path
		fmt.Fprintf(&buf, `,(gogoproto.customtype)=%q`, pt)
	}
	if pt == "google.protobuf.Timestamp" {
		fmt.Fprint(&buf, ",(gogoproto.stdtime)=true")
	}
	if c.IsNull() {
		fmt.Fprint(&buf, ",(gogoproto.nullable)=true")
	}
	return buf.String()
}
