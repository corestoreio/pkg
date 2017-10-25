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

	"github.com/corestoreio/csfw/sql/ddl"
	"github.com/corestoreio/csfw/util/strs"
	"github.com/corestoreio/errors"
)

const goTypeOptions = 4

// These variables are mapping the un/signed and null/not-null types to the
// appropriate Go type.
var (
	goTypeInt64 = [...]string{
		"dml.NullInt64", // unsigned null
		"uint64",        // unsigned not null
		"dml.NullInt64", // signed null
		"int64",         // signed not null
	}
	goTypeInt = [...]string{
		"dml.NullInt64", // unsigned null
		"uint64",        // unsigned not null
		"dml.NullInt64", // signed null
		"int64",         // signed not null
	}
	goTypeFloat64 = [...]string{
		"dml.NullFloat64", // unsigned null
		"float64",         // unsigned not null
		"dml.NullFloat64", // signed null
		"float64",         // signed not null
	}
	goTypeTime = [...]string{
		"dml.NullTime", // unsigned null
		"time.Time",    // unsigned not null
		"dml.NullTime", // signed null
		"time.Time",    // signed not null
	}
	goTypeString = [...]string{
		"dml.NullString", // unsigned null
		"string",         // unsigned not null
		"dml.NullString", // signed null
		"string",         // signed not null
	}
	goTypeBool = [...]string{
		"dml.NullBool", // unsigned null
		"bool",         // unsigned not null
		"dml.NullBool", // signed null
		"bool",         // signed not null
	}
	goTypeDecimal = [...]string{
		"dml.Decimal", // unsigned null
		"dml.Decimal", // unsigned not null
		"dml.Decimal", // signed null
		"dml.Decimal", // signed not null
	}
	goTypeByte = [...]string{
		"[]byte", // unsigned null
		"[]byte", // unsigned not null
		"[]byte", // signed null
		"[]byte", // signed not null
	}
)

// mysqlTypeToGo maps the MySql/MariaDB field type to the correct Go type.
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
		t = goType[0] // unsigned null
	case c.IsUnsigned() && (!c.IsNull() || !withNull):
		t = goType[1] // unsigned not null
	case !c.IsUnsigned() && c.IsNull() && withNull:
		t = goType[2] // signed null
	case !c.IsUnsigned() && (!c.IsNull() || !withNull):
		t = goType[3] // signed not null
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
