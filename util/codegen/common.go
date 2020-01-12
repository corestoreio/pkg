// Copyright 2015-present, Cyrill @ Schumacher.fm and the CoreStore contributors
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

package codegen

import (
	"bytes"
	"strings"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/pkg/util/conv"
	"github.com/corestoreio/pkg/util/strs"
)

// FormatError gets returned when generating the source code in case gofmt or
// any equal other program for a language can't run the formatting.
type FormatError struct {
	error
	Code string
}

func (*FormatError) ErrorKind() errors.Kind {
	return errors.NotAcceptable
}

type common struct {
	SecondLineComments []string
	packageName        string
	*bytes.Buffer
	packageNames map[string]string // Imported package names in the current file.
	indent       string
}

// AddImport adds a new import path. importPath required and packageName optional.
func (g *common) AddImport(importPath, packageName string) {
	g.packageNames[importPath] = packageName
}

// AddImports adds multiple import paths at once. They all must have unique base
// names.
func (g *common) AddImports(importPaths ...string) {
	for _, ip := range importPaths {
		g.packageNames[ip] = ""
	}
}

// Writes a multiline comment and formats it to a max width of 80 chars. It adds
// automatically the comment prefix `//`. It converts all types to string, if it
// can't it panics.
func (g *common) C(comments ...interface{}) {
	cs := make([]string, 0, len(comments))
	for _, cIF := range comments {
		s, err := conv.ToStringE(cIF)
		if err != nil {
			panic(err)
		}
		cs = append(cs, s)
	}
	comment(g.Buffer, cs...)
}

func comment(g *bytes.Buffer, comments ...string) {
	cLines := strings.Split(strs.WordWrap(strings.Join(comments, " "), 78), "\n")
	for _, c := range cLines {
		g.WriteString("// ")
		g.WriteString(c)
		g.WriteByte('\n')
	}
}

var emptyIString interface{} = ""

// SkipWS converts all arguments to type string, panics if it does not support a
// type, and merges all arguments to one single string without white space
// concatenation. This function can be used as argument to Pln or P or C.
func SkipWS(str ...interface{}) []byte {
	var buf bytes.Buffer
	for _, v := range str {
		s, err := conv.ToStringE(v)
		if err != nil {
			panic(err)
		}
		buf.WriteString(s)
	}
	return buf.Bytes()
}

// Pln prints the arguments to the generated output. It tries to convert all
// kind of types to a string. It adds a line break at the end IF there are strs
// to print.
func (g *common) Pln(str ...interface{}) {
	if ls := len(str); ls == 0 || (ls == 1 && str[0] == emptyIString) {
		return
	}
	_, _ = g.WriteString(g.indent)
	for _, v := range str {
		s, err := conv.ToStringE(v)
		if err != nil {
			panic(err)
		}
		_, _ = g.WriteString(s)
		g.WriteByte(' ')
	}
	_ = g.WriteByte('\n')
}

func (g *common) P(str ...interface{}) {
	_, _ = g.WriteString(g.indent)
	for _, v := range str {
		s, err := conv.ToStringE(v)
		if err != nil {
			panic(err)
		}
		_, _ = g.WriteString(s)
		g.WriteByte(' ')
	}
}

// In Indents the output one tab stop.
func (g *common) In() { g.indent += "\t" }

// Out unindents the output one tab stop.
func (g *common) Out() {
	if len(g.indent) > 0 {
		g.indent = g.indent[1:]
	}
}

// EncloseBT encloses the string s in backticks
func EncloseBT(s string) string {
	return "`" + s + "`"
}
