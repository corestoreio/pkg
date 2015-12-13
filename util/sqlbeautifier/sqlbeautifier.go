// Copyright 2015, Cyrill @ Schumacher.fm and the CoreStore contributors
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

package sqlbeautifier

import (
	"bytes"
	"io"
	"io/ioutil"

	"github.com/corestoreio/csfw/util/sqlparser"
)

// FromReader reads data from the Reader interface and converts that data
// to a string and converts the SQL string into a beautiful human readable
// multi line string.
func FromReader(r io.Reader) (*bytes.Buffer, error) {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	return FromString(string(b))
}

// FromString converts a SQL string into a beautiful human readable multi line string.
// The returned buffer is unique for the formatted query.
func FromString(s string) (*bytes.Buffer, error) {

	stmt, err := sqlparser.Parse(s)
	if err != nil {
		return nil, err
	}

	buf := sqlparser.NewTrackedBuffer(nodeFormatter)
	stmt.Format(buf)

	return buf.Buffer, nil
}

// MustFromString same as FromString but panics on error
func MustFromString(s string) string {
	buf, err := FromString(s)
	if err != nil {
		panic(err)
	}
	return buf.String()
}

func nodeFormatter(buf *sqlparser.TrackedBuffer, node sqlparser.SQLNode) {
	// this needs more love but for now it is quite ok
	switch node := node.(type) {
	case sqlparser.SelectExprs:
		node.Format(buf)
		buf.WriteString(" \n")
	case *sqlparser.StarExpr:
		buf.WriteString("\n\t")
		node.Format(buf)
	case *sqlparser.NonStarExpr:
		buf.WriteString("\n\t")
		node.Format(buf)
	case *sqlparser.ColName:
		node.Format(buf)
	case *sqlparser.TableName:
		buf.WriteString("\n\t")
		node.Format(buf)
	case *sqlparser.AliasedTableExpr:
		node.Format(buf)
	case *sqlparser.JoinTableExpr:
		node.Format(buf)
		buf.WriteRune('\n')
	case *sqlparser.Where:
		buf.WriteRune('\n')
		node.Format(buf)
	case *sqlparser.ParenBoolExpr:
		buf.WriteString("\n\t\t")
		node.Format(buf)
	case sqlparser.GroupBy:
		if len(node) > 0 {
			buf.WriteRune('\n')
			node.Format(buf)
		}
	case sqlparser.OrderBy:
		if len(node) > 0 {
			buf.WriteRune('\n')
			node.Format(buf)
		}
	case *sqlparser.Order:
		buf.WriteString("\n\t")
		node.Format(buf)
	case *sqlparser.Limit:
		buf.WriteRune('\n')
		node.Format(buf)
	default:
		//		fmt.Printf("%#v\n", node)
		node.Format(buf)
	}
}
