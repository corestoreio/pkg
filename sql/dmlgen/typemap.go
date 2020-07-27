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

// TypeDef used in variable `mysqlTypeToGo` to map a MySQL/MariaDB type to its
// appropriate Go and serializer type. Those types are getting printed in the
// generated files.
type TypeDef struct {
	GoUNull    string // unsigned null
	GoUNotNull string // unsigned not null
	GoNull     string // A go type which can be null/nil
	GoNotNull  string // A usually native Go primitive type

	SerializerUNull    string
	SerializerUNotNull string
	SerializerNull     string
	SerializerNotNull  string
}

var typeMap = map[string]map[string]*TypeDef{ // immutable
	// Go native type, covers also unsigned and NULL
	"int64": {
		// implementation name. Default uses no serializer.
		"default": &TypeDef{
			GoUNull:    "null.Uint64", // Go unsigned null
			GoUNotNull: "uint64",      // Go unsigned not null
			GoNull:     "null.Int64",  // Go signed null
			GoNotNull:  "int64",       // Go signed not null
		},
		// proto uses protocol buffers or gogoproto as serializer. It requires a
		// mapping of the go type to the protobuf type. Some native Go types are
		// not supported in protobuf. For example if the DB column type is
		// smallint, Go and protobuf types must be minimum int32. This can cause
		// data loss or overflow when writing int32 into the DB. So you think it
		// runs correctly until you load the values back from the DB...
		"protobuf": {
			GoUNull:    "null.Uint64", // Go unsigned null
			GoUNotNull: "uint64",      // Go unsigned not null
			GoNull:     "null.Int64",  // Go signed null
			GoNotNull:  "int64",       // Go signed not null

			// native types of the protobuf implementation. pkg null refers to
			// storage/null/null.proto file
			SerializerUNull:    "null.Uint64", // proto unsigned null | no pun intended ;-)
			SerializerUNotNull: "uint64",      // proto unsigned not null
			SerializerNull:     "null.Int64",  // proto signed null
			SerializerNotNull:  "int64",       // proto signed not null
		},
		"fbs": {
			GoUNull:    "null.Uint64",
			GoUNotNull: "uint64",
			GoNull:     "null.Int64",
			GoNotNull:  "int64",

			// native types of the flatbuffers implementation. pkg null refers
			// to storage/null/null.fbs file
			SerializerUNull:    "null.Uint64", // fbs unsigned null
			SerializerUNotNull: "ulong",       // fbs unsigned not null
			SerializerNull:     "null.Int64",  // fbs signed null
			SerializerNotNull:  "long",        // fbs signed not null
		},
	},
	"int32": {
		"default": &TypeDef{
			GoUNull:    "null.Uint32",
			GoUNotNull: "uint32",
			GoNull:     "null.Int32",
			GoNotNull:  "int32",
		},
		"protobuf": {
			GoUNull:            "null.Uint32",
			GoUNotNull:         "uint32",
			GoNull:             "null.Int32",
			GoNotNull:          "int32",
			SerializerUNull:    "null.Uint32",
			SerializerUNotNull: "uint32",
			SerializerNull:     "null.Int32",
			SerializerNotNull:  "int32",
		},
		"fbs": {
			GoUNull:            "null.Uint32",
			GoUNotNull:         "uint32",
			GoNull:             "null.Int32",
			GoNotNull:          "int32",
			SerializerUNull:    "null.Uint32",
			SerializerUNotNull: "uint",
			SerializerNull:     "null.Int32",
			SerializerNotNull:  "int",
		},
	},
	"int16": {
		"default": &TypeDef{
			GoUNull:    "null.Uint16",
			GoUNotNull: "uint16",
			GoNull:     "null.Int16",
			GoNotNull:  "int16",
		},
		"protobuf": {
			GoUNull:            "null.Uint32",
			GoUNotNull:         "uint32",
			GoNull:             "null.Int32",
			GoNotNull:          "int32",
			SerializerUNull:    "null.Uint32",
			SerializerUNotNull: "uint32",
			SerializerNull:     "null.Int32",
			SerializerNotNull:  "int32",
		},
		"fbs": {
			GoUNull:            "null.Uint16",
			GoUNotNull:         "uint16",
			GoNull:             "null.Int16",
			GoNotNull:          "int16",
			SerializerUNull:    "null.Uint16",
			SerializerUNotNull: "ushort",
			SerializerNull:     "null.Int16",
			SerializerNotNull:  "short",
		},
	},
	"int8": {
		"default": &TypeDef{
			GoUNull:    "null.Uint8",
			GoUNotNull: "uint8",
			GoNull:     "null.Int8",
			GoNotNull:  "int8",
		},
		"protobuf": {
			GoUNull:            "null.Uint32",
			GoUNotNull:         "uint32",
			GoNull:             "null.Int32",
			GoNotNull:          "int32",
			SerializerUNull:    "null.Uint32",
			SerializerUNotNull: "uint32",
			SerializerNull:     "null.Int32",
			SerializerNotNull:  "int32",
		},
		"fbs": {
			GoUNull:            "null.Uint8",
			GoUNotNull:         "uint8",
			GoNull:             "null.Int8",
			GoNotNull:          "int8",
			SerializerUNull:    "null.Uint8",
			SerializerUNotNull: "ubyte",
			SerializerNull:     "null.Int8",
			SerializerNotNull:  "byte",
		},
	},
	"float64": {
		"default": &TypeDef{
			GoUNull:    "null.Uint8",
			GoUNotNull: "float64",
			GoNull:     "null.Int8",
			GoNotNull:  "float64",
		},
		"protobuf": {
			GoUNull:            "null.Float64",
			GoUNotNull:         "float64",
			GoNull:             "null.Float64",
			GoNotNull:          "float64",
			SerializerUNull:    "null.Float64",
			SerializerUNotNull: "double",
			SerializerNull:     "null.Float64",
			SerializerNotNull:  "double",
		},
		"fbs": {
			GoUNull:            "null.Float64",
			GoUNotNull:         "float64",
			GoNull:             "null.Float64",
			GoNotNull:          "float64",
			SerializerUNull:    "null.Double",
			SerializerUNotNull: "double",
			SerializerNull:     "null.Double",
			SerializerNotNull:  "double",
		},
	},
	"time": {
		"default": &TypeDef{
			GoUNull:    "null.Time",
			GoUNotNull: "time.Time",
			GoNull:     "null.Time",
			GoNotNull:  "time.Time",
		},
		"protobuf": {
			GoUNull:            "null.Time",
			GoUNotNull:         "time.Time",
			GoNull:             "null.Time",
			GoNotNull:          "time.Time",
			SerializerUNull:    "null.Time",
			SerializerUNotNull: "google.protobuf.Timestamp",
			SerializerNull:     "null.Time",
			SerializerNotNull:  "google.protobuf.Timestamp",
		},
		"fbs": {
			GoUNull:            "null.Time",
			GoUNotNull:         "time.Time",
			GoNull:             "null.Time",
			GoNotNull:          "time.Time",
			SerializerUNull:    "null.Time",
			SerializerUNotNull: "null.Time",
			SerializerNull:     "null.Time",
			SerializerNotNull:  "null.Time",
		},
	},
	"string": {
		"default": &TypeDef{
			GoUNull:    "null.String",
			GoUNotNull: "string",
			GoNull:     "null.String",
			GoNotNull:  "string",
		},
		"protobuf": {
			GoUNull:            "null.String",
			GoUNotNull:         "string",
			GoNull:             "null.String",
			GoNotNull:          "string",
			SerializerUNull:    "null.String",
			SerializerUNotNull: "string",
			SerializerNull:     "null.String",
			SerializerNotNull:  "string",
		},
		"fbs": {
			GoUNull:            "null.String",
			GoUNotNull:         "string",
			GoNull:             "null.String",
			GoNotNull:          "string",
			SerializerUNull:    "null.String",
			SerializerUNotNull: "string",
			SerializerNull:     "null.String",
			SerializerNotNull:  "string",
		},
	},
	"bool": {
		"default": &TypeDef{
			GoUNull:    "null.Bool",
			GoUNotNull: "bool",
			GoNull:     "null.Bool",
			GoNotNull:  "bool",
		},
		"protobuf": {
			GoUNull:            "null.Bool",
			GoUNotNull:         "bool",
			GoNull:             "null.Bool",
			GoNotNull:          "bool",
			SerializerUNull:    "null.Bool",
			SerializerUNotNull: "bool",
			SerializerNull:     "null.Bool",
			SerializerNotNull:  "bool",
		},
		"fbs": {
			GoUNull:            "null.Bool",
			GoUNotNull:         "bool",
			GoNull:             "null.Bool",
			GoNotNull:          "bool",
			SerializerUNull:    "null.Bool",
			SerializerUNotNull: "bool",
			SerializerNull:     "null.Bool",
			SerializerNotNull:  "bool",
		},
	},
	"decimal": {
		"default": &TypeDef{
			GoUNull:    "null.Decimal",
			GoUNotNull: "null.Decimal",
			GoNull:     "null.Decimal",
			GoNotNull:  "null.Decimal",
		},
		"protobuf": {
			GoUNull:            "null.Decimal",
			GoUNotNull:         "null.Decimal",
			GoNull:             "null.Decimal",
			GoNotNull:          "null.Decimal",
			SerializerUNull:    "null.Decimal",
			SerializerUNotNull: "null.Decimal",
			SerializerNull:     "null.Decimal",
			SerializerNotNull:  "null.Decimal",
		},
		"fbs": {
			GoUNull:            "null.Decimal",
			GoUNotNull:         "null.Decimal",
			GoNull:             "null.Decimal",
			GoNotNull:          "null.Decimal",
			SerializerUNull:    "null.Decimal",
			SerializerUNotNull: "null.Decimal",
			SerializerNull:     "null.Decimal",
			SerializerNotNull:  "null.Decimal",
		},
	},
	"byte": {
		"default": &TypeDef{
			GoUNull:    "[]byte",
			GoUNotNull: "[]byte",
			GoNull:     "[]byte",
			GoNotNull:  "[]byte",
		},
		"protobuf": {
			GoUNull:            "[]byte",
			GoUNotNull:         "[]byte",
			GoNull:             "[]byte",
			GoNotNull:          "[]byte",
			SerializerUNull:    "bytes",
			SerializerUNotNull: "bytes",
			SerializerNull:     "bytes",
			SerializerNotNull:  "bytes",
		},
		"fbs": {
			GoUNull:            "[]byte",
			GoUNotNull:         "[]byte",
			GoNull:             "[]byte",
			GoNotNull:          "[]byte",
			SerializerUNull:    "[ubyte]",
			SerializerUNotNull: "[ubyte]",
			SerializerNull:     "[ubyte]",
			SerializerNotNull:  "[ubyte]",
		},
	},
}

func mustTMK(key string) map[string]*TypeDef {
	if m, ok := typeMap[key]; ok {
		return m
	}
	panic(fmt.Sprintf("[dmlgen] Key %q not found in typeMap", key))
}

// mysqlTypeToGo maps the MySql/MariaDB field type to the correct Go/protobuf
// type. See the type TypeDef for more details. This exported variable allows to
// set custom types before code generation.
var mysqlTypeToGo = map[string]map[string]*TypeDef{ // immutable
	// key1=MySQL Type; key2=serializer;
	"int":        mustTMK("int32"),
	"bigint":     mustTMK("int64"),
	"smallint":   mustTMK("int16"),
	"tinyint":    mustTMK("int8"),
	"mediumint":  mustTMK("int32"),
	"double":     mustTMK("float64"),
	"float":      mustTMK("float64"),
	"decimal":    mustTMK("decimal"),
	"date":       mustTMK("time"),
	"datetime":   mustTMK("time"),
	"timestamp":  mustTMK("time"),
	"time":       mustTMK("time"),
	"char":       mustTMK("string"),
	"varchar":    mustTMK("string"),
	"enum":       mustTMK("string"),
	"set":        mustTMK("string"),
	"text":       mustTMK("string"),
	"longtext":   mustTMK("string"),
	"mediumtext": mustTMK("string"),
	"tinytext":   mustTMK("string"),
	"blob":       mustTMK("byte"),
	"longblob":   mustTMK("byte"),
	"mediumblob": mustTMK("byte"),
	"tinyblob":   mustTMK("byte"),
	"binary":     mustTMK("byte"),
	"varbinary":  mustTMK("byte"),
	"bit":        mustTMK("bool"),
	// TODO add more MySQL types like JSON or GEO
}

func mustGetTypeDef(mysqlDataType, serializer string) *TypeDef {
	if serializer == "" {
		serializer = "default"
	}
	myType, ok := mysqlTypeToGo[mysqlDataType] // readonly access so safe for concurrent access
	if !ok {
		panic(errors.NotFound.Newf("[dmlgen] MySQL type %q not found", mysqlDataType))
	}

	goType, ok := myType[serializer]
	if !ok {
		panic(errors.NotFound.Newf("[dmlgen] Serializer %q not found", serializer))
	}
	return goType
}

func (g *Generator) findType(c *ddl.Column) *TypeDef {
	goType := mustGetTypeDef(c.DataType, g.Serializer)

	// The switch block overwrites the already retrieved goType by checking for
	// bool columns and columns which contains a money unit.
	switch {
	case c.IsBool():
		goType = mustGetTypeDef("bit", g.Serializer)
	case c.IsFloat() && c.IsMoney():
		goType = mustGetTypeDef("decimal", g.Serializer)
	}
	return goType
}

func (g *Generator) goTypeNull(c *ddl.Column) string { return g.mySQLToGoType(c, true) }
func (g *Generator) goType(c *ddl.Column) string     { return g.mySQLToGoType(c, false) }
func (g *Generator) goFuncNull(c *ddl.Column) string { return g.mySQLToGoDmlColumnMap(c, true) }

// mySQLToGoType calculates the data type of the field DataType. For example
// bigint, smallint, tinyint will result in "int". If withNull is true the
// returned type can store a null value.
func (g *Generator) mySQLToGoType(c *ddl.Column, withNull bool) string {
	goType := g.findType(c)

	var t string
	switch {
	case c.IsUnsigned() && c.IsNull() && withNull:
		t = goType.GoUNull // unsigned null
	case c.IsUnsigned() && (!c.IsNull() || !withNull):
		t = goType.GoUNotNull // unsigned not null
	case !c.IsUnsigned() && c.IsNull() && withNull:
		t = goType.GoNull // signed null
	case !c.IsUnsigned() && (!c.IsNull() || !withNull):
		t = goType.GoNotNull // signed not null
	}
	return t
}

// toGoPrimitiveFromNull returns for a Go type or structure the final primitive:
// int->int but NullInt->.Int. Either it is the struct field name or final type
// of a composite nullable type.
func (g *Generator) toGoPrimitiveFromNull(c *ddl.Column) string {
	t := g.mySQLToGoType(c, true)
	field := strs.ToGoCamelCase(c.Field)
	if strings.HasPrefix(t, "null.") && t != "null.Decimal" {
		f := t[5:] // 5 == len("null.")
		if t == "null.String" {
			f = "Data" // null.String type has field name `Data` instead of `String`
		}
		t = field + "." + f
	} else {
		t = field
	}
	return t
}

func (g *Generator) mySQLToGoDmlColumnMap(c *ddl.Column, withNull bool) string {
	gt := g.mySQLToGoType(c, withNull)
	if gt == "[]byte" {
		return "Byte"
	}

	if dot := strings.IndexByte(gt, '.'); dot > 0 {
		pkg := gt[:dot]
		fnName := gt[dot+1:]
		if fnName == "Decimal" || (pkg == "time" && fnName == "Time") {
			return fnName
		}
		r, n := utf8.DecodeRuneInString(pkg)
		return string(unicode.ToUpper(r)) + pkg[n:] + fnName
	}
	r, n := utf8.DecodeRuneInString(gt)
	return string(unicode.ToUpper(r)) + gt[n:]
}

func (g *Generator) toSerializerType(c *ddl.Column, withNull bool) string {
	goType := g.findType(c)

	var t string
	switch {
	case c.IsUnsigned() && c.IsNull() && withNull:
		t = goType.SerializerUNull // unsigned null
	case c.IsUnsigned() && (!c.IsNull() || !withNull):
		t = goType.SerializerUNotNull // unsigned not null
	case !c.IsUnsigned() && c.IsNull() && withNull:
		t = goType.SerializerNull // signed null
	case !c.IsUnsigned() && (!c.IsNull() || !withNull):
		t = goType.SerializerNotNull // signed not null
	}

	return t
}

// serializerType converts the column type to the supported type of the current
// serializer. For now supports only protobuf.
func (g *Generator) serializerType(c *ddl.Column) string {
	pt := g.toSerializerType(c, true)
	if strings.IndexByte(pt, '/') > 0 { // slash identifies an import path
		return "bytes"
	}
	return pt
}

func mySQLType2GoComparisonOperator(c *ddl.Column) string {
	switch c.DataType {

	case "blob", "char", "longblob", "longtext", "mediumblob", "mediumtext",
		"text", "tinytext", "varbinary", "varchar", "enum", "set":
		return ` != ""`

	case "date", "datetime", "time", "timestamp":
		return ".IsZero() == false"

	case "decimal":
		return ".Valid "

	case "bigint", "double", "float", "int", "smallint", "tinyint":
		if c.IsNull() {
			return ".Valid"
		}
		if c.IsUnsigned() {
			return " > 0"
		}

		return " != 0"

	default:
		return " > 0 /*TODO find correct case*/"
	}
}
