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

package csdb

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"strings"
	"unicode"

	"github.com/corestoreio/errors"
)

// maxIdentifierLength see http://dev.mysql.com/doc/refman/5.7/en/identifiers.html
const maxIdentifierLength = 64

// IsValidIdentifier checks the permissible syntax for identifiers. Certain
// objects within MySQL, including database, table, index, column, alias, view,
// stored procedure, partition, tablespace, and other object names are known as
// identifiers. ASCII: [0-9,a-z,A-Z$_] (basic Latin letters, digits 0-9, dollar,
// underscore) Max length 63 characters. Returns errors.NotValid
//
// http://dev.mysql.com/doc/refman/5.7/en/identifiers.html
func IsValidIdentifier(names ...string) error {
	if len(names) == 0 {
		return errors.NewNotValidf("[csdb] No arguments provided")
	}
	for _, name := range names {
		if len(name) > maxIdentifierLength || name == "" {
			return errors.NewNotValidf("[csdb] Incorrect identifier. Too long or empty: %q", name)
		}

		for _, r := range name {
			if mapAlNum(r) != r {
				return errors.NewNotValidf("[csdb] Invalid character %q in name %q", string(r), name)
			}
		}
	}
	return nil
}

func mapAlNum(r rune) rune {
	var ok bool
	switch {
	case '0' <= r && r <= '9':
		ok = true
	case 'a' <= r && r <= 'z', 'A' <= r && r <= 'Z':
		ok = true
	case r == '$', r == '_':
		ok = true
	}
	if !ok {
		return -1
	}
	return r
}

func mapAlNumUpper(r rune) rune {
	r = mapAlNum(r)
	if r < 0 {
		return r
	}
	return unicode.ToUpper(r)
}

// cleanIdentifier removes all invalid characters
// https://dev.mysql.com/doc/refman/5.7/en/identifiers.html
func cleanIdentifier(upper bool, name []byte) string {
	fn := mapAlNum
	if upper {
		fn = mapAlNumUpper
	}
	return string(bytes.Map(fn, name))
}

// TableName generates a table name, shortens it, if necessary, and removes all
// invalid characters. First round of shortening goes by replacing common words
// with their abbreviations and in the second round creating a MD5 hash of the
// table name.
func TableName(prefix, name string, suffixes ...string) string {
	if prefix == "" && len(suffixes) == 0 && len(name) <= maxIdentifierLength {
		return strings.Map(mapAlNum, name)
	}

	var buf = make([]byte, 0, maxIdentifierLength)
	if !strings.HasPrefix(name, prefix) {
		buf = append(buf, prefix...)
	}
	buf = append(buf, name...)
	for _, s := range suffixes {
		buf = append(buf, '_')
		buf = append(buf, s...)
	}
	return cleanIdentifier(false, shortenEntityName(buf, "t_"))
}

// IndexName creates a new valid index name. IndexType can only be one of the
// three following enums: `index`, `unique` or `fulltext`. If empty or mismatch
// it falls back to `index`. The returned string represents a valid identifier
// within MySQL.
func IndexName(indexType, tableName string, fields ...string) string {
	var prefix = "IDX_"
	switch indexType {
	case "unique":
		prefix = "UNQ_"
	case "fulltext":
		prefix = "FTI_"
	}

	var buf = make([]byte, 0, maxIdentifierLength)
	buf = append(buf, tableName...)
	for i, f := range fields {
		if i == 0 {
			buf = append(buf, '_')
		}
		buf = append(buf, f...)
		if i < len(fields)-1 {
			buf = append(buf, '_')
		}
	}
	return cleanIdentifier(true, shortenEntityName(buf, prefix))
}

// TriggerName creates a new trigger name. The returned string represents a
// valid identifier within MySQL. Argument time should be either `before` or
// `after`. Event should be one of the following types: `insert`, `update` or
// `delete`
func TriggerName(tableName, time, event string) string {
	var buf = make([]byte, 0, maxIdentifierLength)
	buf = append(buf, tableName...)
	buf = append(buf, '_')
	buf = append(buf, time...)
	buf = append(buf, '_')
	buf = append(buf, event...)
	return cleanIdentifier(false, shortenEntityName(buf, "trg_"))
}

// ForeignKeyName creates a new foreign key name. The returned string represents
// a valid identifier within MySQL.
func ForeignKeyName(priTableName, priColumnName, refTableName, refColumnName string) string {
	var buf = make([]byte, 0, maxIdentifierLength)
	buf = append(buf, priTableName...)
	buf = append(buf, '_')
	buf = append(buf, priColumnName...)
	buf = append(buf, '_')
	buf = append(buf, refTableName...)
	buf = append(buf, '_')
	buf = append(buf, refColumnName...)
	return cleanIdentifier(true, shortenEntityName(buf, "FK_"))
}

// TODO: micro optimize later 8-) to reduce allocations
func shortenEntityName(name []byte, prefix string) []byte {
	if len(name) < maxIdentifierLength {
		return name
	}
	name2 := name[:0]
	name2 = append(name2, translatedAbbreviations.Replace(string(name))...)
	if len(name2) > maxIdentifierLength {
		return []byte(fmt.Sprintf("%s%x", prefix, md5.Sum(name2))) // worse worse case
	}
	return name2
}

// translatedAbbreviations contains a list of names which gets translated to their abbreviation if an MySQL identifier has more
var translatedAbbreviations = strings.NewReplacer(
	"address", "addr",
	"admin", "adm",
	"aggregat", "aggr",
	"agreement", "agrt",
	"attribute", "attr",
	"bundle", "bndl",
	"calculation", "calc",
	"catalog", "cat",
	"category", "ctgr",
	"checkout", "chkt",
	"compare", "cmp",
	"customer", "cstr",
	"datetime", "dtime",
	"decimal", "dec",
	"directory", "dir",
	"downloadable", "dl",
	"element", "elm",
	"enterprise", "ent",
	"entity", "entt",
	"fieldset", "fset",
	"gallery", "glr",
	"index", "idx",
	"inventory", "inv",
	"label", "lbl",
	"layout", "lyt",
	"link", "lnk",
	"media", "mda",
	"minimal", "min",
	"maximal", "max",
	"newsletter", "nlttr",
	"notification", "ntfc",
	"option", "opt",
	"product", "prd",
	"query", "qr",
	"resource", "res",
	"search", "srch",
	"session", "sess",
	"shipping", "shpp",
	"status", "sts",
	"super", "spr",
	"title", "ttl",
	"user", "usr",
	"value", "val",
	"varchar", "vchr",
	"website", "ws",
)
