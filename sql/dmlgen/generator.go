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

package dmlgen

import (
	"bytes"
	"fmt"
	"go/format"
	"io"
	"sort"
	"strings"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/pkg/util/conv"
	"github.com/corestoreio/pkg/util/strs"
)

type Generator struct {
	packageName string
	BuildTags   string
	*bytes.Buffer
	init            []string          // Lines to emit in the init function.
	packageNames    map[string]string // Imported package names in the current file.
	constantsString []string
	indent          string
}

// NewGenerator creates a new source code generator for a specific new package.
func NewGenerator(packageName string) *Generator {
	return &Generator{
		Buffer:       new(bytes.Buffer),
		packageName:  packageName,
		packageNames: map[string]string{},
	}
}

// AddImport adds a new import path. importPath required and packageName optional.
func (g *Generator) AddImport(importPath, packageName string) {
	g.packageNames[importPath] = packageName
}

func (g *Generator) AddConstString(name, value string) {
	g.constantsString = append(g.constantsString, fmt.Sprintf("\t%s = %q\n", name, value))
}

// Writes a multiline comment and formats it to a max width of 80 chars. It adds
// automatically the comment prefix `//`.
func (g *Generator) C(comments ...string) {
	cLines := strings.Split(strs.WordWrap(strings.Join(comments, " "), 78), "\n")
	for _, c := range cLines {
		g.WriteString("// ")
		g.WriteString(c)
		g.WriteByte('\n')
	}
}

// P prints the arguments to the generated output. It tries to convert all kind
// of types to a string.
func (g *Generator) P(str ...interface{}) {
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

// AddInitf stores the given statement to be printed inside the file's init
// function. The statement is given as a format specifier and arguments.
func (g *Generator) AddInitf(stmt string, a ...interface{}) {
	g.init = append(g.init, fmt.Sprintf(stmt, a...))
}

func (g *Generator) generateInitFunction() {
	if len(g.init) == 0 {
		return
	}
	g.P("func init() {")
	g.In()
	for _, l := range g.init {
		g.P(l)
	}
	g.Out()
	g.P("}")
	g.init = nil
}

func (g *Generator) generateImports(w io.Writer) {
	fmt.Fprintln(w, "import (")
	pkgSorted := make([]string, 0, len(g.packageNames))
	for key := range g.packageNames {
		pkgSorted = append(pkgSorted, key)
	}
	sort.Strings(pkgSorted)
	for _, p := range pkgSorted {
		fmt.Fprintf(w, "\t%s %q\n", g.packageNames[p], p)
	}
	fmt.Fprintln(w, ")")
}

func (g *Generator) generateConstants(w io.Writer) {
	fmt.Fprintln(w, "const (")
	sort.Strings(g.constantsString)
	for _, cs := range g.constantsString {
		fmt.Fprint(w, cs)
	}
	fmt.Fprintln(w, ")")
}

// In Indents the output one tab stop.
func (g *Generator) In() { g.indent += "\t" }

// Out unindents the output one tab stop.
func (g *Generator) Out() {
	if len(g.indent) > 0 {
		g.indent = g.indent[1:]
	}
}

func (g *Generator) GenerateFile(w io.Writer) error {

	var buf bytes.Buffer
	if g.BuildTags != "" {
		fmt.Fprintln(&buf, "// +build ", g.BuildTags)
		fmt.Fprint(&buf, "\n") // the extra line as required from the Go spec
	}
	fmt.Fprintf(&buf, "package %s\n", g.packageName)
	fmt.Fprintln(&buf, "// Auto generated via github.com/corestoreio/pkg/sql/dmlgen")
	g.generateImports(&buf)
	g.generateConstants(&buf)

	g.Buffer.WriteTo(&buf)

	fmted, err := format.Source(buf.Bytes())
	if err != nil {
		return errors.NotAcceptable.New(err, "\nSource Code:\n%s\n", buf.String())
	}
	_, err = w.Write(fmted)
	return err
}
