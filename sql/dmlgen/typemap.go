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
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/pkg/sql/ddl"
	"github.com/corestoreio/pkg/util/strs"
)

// TypeDef used in variable `MysqlTypeToGo` to map a MySQL/MariaDB type to its
// appropriate Go or Protocol Buffer type. Those types are getting printed in
// the generated files.
type TypeDef struct {
	MysqlUnsignedNull    string
	MysqlUnsignedNotNull string
	MysqlSignedNull      string
	MysqlSignedNotNull   string

	ProtobufUnsignedNull    string
	ProtobufUnsignedNotNull string
	ProtobufSignedNull      string
	ProtobufSignedNotNull   string
}

// These variables are mapping the un/signed and null/not-null types to the
// appropriate Go type.
var (
	// TODO further optimize this mappings
	goTypeInt64 = &TypeDef{
		MysqlUnsignedNull:    "null.Int64",
		MysqlUnsignedNotNull: "uint64",
		MysqlSignedNull:      "null.Int64",
		MysqlSignedNotNull:   "int64",

		ProtobufUnsignedNull:    "null.Int64", // Proto package and its type, not the Go package!
		ProtobufUnsignedNotNull: "uint64",
		ProtobufSignedNull:      "null.Int64", // Proto package and its type, not the Go package!
		ProtobufSignedNotNull:   "int64",
	}
	goTypeInt = &TypeDef{
		MysqlUnsignedNull:    "null.Int64",
		MysqlUnsignedNotNull: "uint64",
		MysqlSignedNull:      "null.Int64",
		MysqlSignedNotNull:   "int64",

		ProtobufUnsignedNull:    "null.Int64", // Proto package and its type, not the Go package!
		ProtobufUnsignedNotNull: "uint64",
		ProtobufSignedNull:      "null.Int64", // Proto package and its type, not the Go package!
		ProtobufSignedNotNull:   "int64",
	}
	goTypeFloat64 = &TypeDef{
		MysqlUnsignedNull:    "null.Float64",
		MysqlUnsignedNotNull: "float64",
		MysqlSignedNull:      "null.Float64",
		MysqlSignedNotNull:   "float64",

		ProtobufUnsignedNull:    "null.Float64", // Proto package and its type, not the Go package!
		ProtobufUnsignedNotNull: "double",
		ProtobufSignedNull:      "null.Float64", // Proto package and its type, not the Go package!
		ProtobufSignedNotNull:   "double",
	}
	goTypeTime = &TypeDef{
		MysqlUnsignedNull:    "null.Time",
		MysqlUnsignedNotNull: "time.Time",
		MysqlSignedNull:      "null.Time",
		MysqlSignedNotNull:   "time.Time",

		ProtobufUnsignedNull:    "null.Time", // Proto package and its type, not the Go package!
		ProtobufUnsignedNotNull: "google.protobuf.Timestamp",
		ProtobufSignedNull:      "null.Time", // Proto package and its type, not the Go package!
		ProtobufSignedNotNull:   "google.protobuf.Timestamp",
	}
	goTypeString = &TypeDef{
		MysqlUnsignedNull:    "null.String",
		MysqlUnsignedNotNull: "string",
		MysqlSignedNull:      "null.String",
		MysqlSignedNotNull:   "string",

		ProtobufUnsignedNull:    "null.String", // Proto package and its type, not the Go package!
		ProtobufUnsignedNotNull: "string",
		ProtobufSignedNull:      "null.String", // Proto package and its type, not the Go package!
		ProtobufSignedNotNull:   "string",
	}
	goTypeBool = &TypeDef{
		MysqlUnsignedNull:    "null.Bool",
		MysqlUnsignedNotNull: "bool",
		MysqlSignedNull:      "null.Bool",
		MysqlSignedNotNull:   "bool",

		ProtobufUnsignedNull:    "null.Bool", // Proto package and its type, not the Go package!
		ProtobufUnsignedNotNull: "bool",
		ProtobufSignedNull:      "null.Bool", // Proto package and its type, not the Go package!
		ProtobufSignedNotNull:   "bool",
	}
	goTypeDecimal = &TypeDef{
		MysqlUnsignedNull:    "dml.Decimal",
		MysqlUnsignedNotNull: "dml.Decimal",
		MysqlSignedNull:      "dml.Decimal",
		MysqlSignedNotNull:   "dml.Decimal",

		ProtobufUnsignedNull:    "dml.Decimal", // Proto package and its type not the Go package!
		ProtobufUnsignedNotNull: "dml.Decimal", // Proto package and its type not the Go package!
		ProtobufSignedNull:      "dml.Decimal", // Proto package and its type not the Go package!
		ProtobufSignedNotNull:   "dml.Decimal", // Proto package and its type not the Go package!
	}
	goTypeByte = &TypeDef{
		MysqlUnsignedNull:    "[]byte",
		MysqlUnsignedNotNull: "[]byte",
		MysqlSignedNull:      "[]byte",
		MysqlSignedNotNull:   "[]byte",

		ProtobufUnsignedNull:    "bytes",
		ProtobufUnsignedNotNull: "bytes",
		ProtobufSignedNull:      "bytes",
		ProtobufSignedNotNull:   "bytes",
	}
)

// MysqlTypeToGo maps the MySql/MariaDB field type to the correct Go/protobuf
// type. See the type TypeDef for more details. This exported variable allows to
// set custom types before code generation.
var MysqlTypeToGo = map[string]*TypeDef{
	"int":        goTypeInt,
	"bigint":     goTypeInt64,
	"smallint":   goTypeInt,
	"tinyint":    goTypeInt,
	"mediumint":  goTypeInt,
	"double":     goTypeFloat64,
	"float":      goTypeFloat64,
	"decimal":    goTypeDecimal,
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

func findType(c *ddl.Column) *TypeDef {

	goType, ok := MysqlTypeToGo[c.DataType]
	if !ok {
		panic(errors.NotFound.Newf("[dmlgen] MySQL type %q not found", c.DataType))
	}

	// The switch block overwrites the already retrieved goType by checking for
	// bool columns and columns which contains a money unit.
	switch {
	case c.IsBool():
		goType = MysqlTypeToGo["bit"]
	case c.IsFloat() && c.IsMoney():
		goType = MysqlTypeToGo["decimal"]
	}
	return goType
}

// mySQLToGoType calculates the data type of the field DataType. For example
// bigint, smallint, tinyint will result in "int". If withNull is true the
// returned type can store a null value.
func mySQLToGoType(c *ddl.Column, withNull bool) string {

	goType := findType(c)

	var t string
	switch {
	case c.IsUnsigned() && c.IsNull() && withNull:
		t = goType.MysqlUnsignedNull // unsigned null
	case c.IsUnsigned() && (!c.IsNull() || !withNull):
		t = goType.MysqlUnsignedNotNull // unsigned not null
	case !c.IsUnsigned() && c.IsNull() && withNull:
		t = goType.MysqlSignedNull // signed null
	case !c.IsUnsigned() && (!c.IsNull() || !withNull):
		t = goType.MysqlSignedNotNull // signed not null
	}

	return t
}

// toGoPrimitive returns for Go type or structure the final primitive:
// int->int but NullInt->.Int
func toGoPrimitive(c *ddl.Column) string {
	t := mySQLToGoType(c, true)
	field := strs.ToGoCamelCase(c.Field)
	if strings.HasPrefix(t, "null.") {
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

	goType := findType(c)

	var t string
	switch {
	case c.IsUnsigned() && c.IsNull() && withNull:
		t = goType.ProtobufUnsignedNull // unsigned null
	case c.IsUnsigned() && (!c.IsNull() || !withNull):
		t = goType.ProtobufUnsignedNotNull // unsigned not null
	case !c.IsUnsigned() && c.IsNull() && withNull:
		t = goType.ProtobufSignedNull // signed null
	case !c.IsUnsigned() && (!c.IsNull() || !withNull):
		t = goType.ProtobufSignedNotNull // signed not null
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
	if pt == "google.protobuf.Timestamp" {
		fmt.Fprint(&buf, ",(gogoproto.stdtime)=true")
	}
	if c.IsNull() || strings.IndexByte(pt, '.') > 0 /*whenever it is a custom type like null. or google.proto.timestamp*/ {
		// Indeed nullable Go Types must be not-nullable in Protobuf because we
		// have a non-pointer struct type which contains the field Valid.
		// Protobuf treats nullable fields as pointer fields, but that is
		// ridiculous.
		fmt.Fprint(&buf, ",(gogoproto.nullable)=false")
	}
	return buf.String()
}
