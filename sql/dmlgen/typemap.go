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
	"github.com/corestoreio/csfw/sql/ddl"
	"github.com/corestoreio/errors"
	"strings"
	"unicode"
	"unicode/utf8"
)

const goTypeOptions = 4

var (
	goTypeInt64 = [...]string{
		"dml.NullInt64", // unsigned null
		"uint64",        // unsigned not null
		"dml.NullInt64", // signed null
		"int64",         // signed not null
	}
	goTypeInt = [...]string{
		"dml.NullInt64", // unsigned null
		"uint",          // unsigned not null
		"dml.NullInt64", // signed null
		"int",           // signed not null
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
	goTypeMoney = [...]string{
		"money.Money", // unsigned null
		"money.Money", // unsigned not null
		"money.Money", // signed null
		"money.Money", // signed not null
	}
	goTypeByte = [...]string{
		"[]byte", // unsigned null
		"[]byte", // unsigned not null
		"[]byte", // signed null
		"[]byte", // signed not null
	}
)

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

// MySQLToGoType calculates the data type of the field DataType. The
// calculated result will be cached. For example bigint, smallint, tinyint will
// result in "int". The returned string guarantees to be lower case. Available
// returned types are: bool, bytes, date, float, int, money, string, time. Data
// type money is special for the database schema. This function is thread safe.
func MySQLToGoType(c *ddl.Column) string {

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
		goType = goTypeMoney
	}

	var t string
	switch {
	case c.IsUnsigned() && c.IsNull():
		t = goType[0] // unsigned null
	case c.IsUnsigned() && !c.IsNull():
		t = goType[1] // unsigned not null
	case !c.IsUnsigned() && c.IsNull():
		t = goType[2] // signed null
	case !c.IsUnsigned() && !c.IsNull():
		t = goType[3] // signed not null
	}

	return t
}

func GoTypeFuncName(c *ddl.Column) string {

	gt := MySQLToGoType(c)
	switch gt {
	case "money.Money":
		return "NullFloat64" // TODO
	case "[]byte":
		return "Byte"
	}

	if dot := strings.IndexByte(gt, '.'); dot > 0 {
		return gt[dot+1:]
	}
	r, n := utf8.DecodeRuneInString(gt)
	return string(unicode.ToUpper(r)) + gt[n:]
}
