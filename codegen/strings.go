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

package codegen

import (
	"bytes"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"text/template"
	"time"
	"unicode"

	"go/format"

	"path/filepath"

	"github.com/juju/errgo"
)

var (
	logFatalln = log.Fatalln
	logFatalf  = log.Fatalf
	letters    = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	Copyright  = []byte(`// Copyright 2015, Cyrill @ Schumacher.fm and the CoreStore contributors
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
`)
)

// OFile is a OutputFile
type OFile string

func (of OFile) String() string {
	return string(of) + ".go"
}

func (of OFile) AppendDir(s ...string) OFile {
	sof := string(of)
	parts := strings.Split(sof, string(filepath.Separator))
	parts = append(parts, s...)
	nf := NewOFile(parts...)
	if filepath.IsAbs(sof) {
		nf = OFile(string(filepath.Separator)) + nf
	}
	return nf
}

func (of OFile) AppendName(s ...string) OFile {
	return NewOFile(string(of) + strings.Join(s, ""))
}

// NewOFile creates a new path from parts
func NewOFile(paths ...string) OFile {
	return OFile(filepath.Join(paths...))
}

// GenerateCode uses text/template for create Go code. package name pkg will also be used
// to remove stutter in variable names.
func GenerateCode(pkg, tplCode string, data interface{}, addFM template.FuncMap) ([]byte, error) {

	funcMap := template.FuncMap{
		"quote":           func(s string) string { return "`" + s + "`" },
		"prepareVar":      prepareVar(pkg),
		"toLowerFirst":    toLowerFirst,
		"prepareVarIndex": func(i int, s string) string { return fmt.Sprintf("%03d%s", i, prepareVar(pkg)(s)) },
		"sprintf":         fmt.Sprintf,
	}
	for k, v := range addFM {
		funcMap[k] = v
	}

	codeTpl := template.Must(template.New("tpl_code").Funcs(funcMap).Parse(tplCode))

	var buf = &bytes.Buffer{}
	err := codeTpl.Execute(buf, data)
	if err != nil {
		return nil, errgo.Mask(err)
	}

	fmt, err := format.Source(buf.Bytes())
	if err != nil {
		return buf.Bytes(), err
	}
	return fmt, nil
}

func toLowerFirst(s string) string {
	if len(s) == 0 {
		return ""
	}
	return strings.ToLower(s[:1]) + s[1:]
}

// prepareVar converts a string into a Go code variable. Removes the package name if this string
// starts with the package name. Replaces all illegal characters with an underscore.
func prepareVar(pkg string) func(s string) string {

	return func(str string) string {

		l := len(pkg) + 1
		if len(str) > l && str[:l] == pkg+TableNameSeparator {
			str = str[l:]
		}

		str = strings.Map(func(r rune) rune {
			switch {
			case r >= 'A' && r <= 'Z', r >= 'a' && r <= 'z', r >= '0' && r <= '9':
				return r
			}
			return '_'
		}, str)

		return Camelize(str)
	}
}

// Camelize transforms from snake case to camelCase e.g. catalog_product_id to CatalogProductID. Also removes quotes.
func Camelize(s string) string {
	s = strings.ToLower(strings.Replace(s, `"`, "", -1))
	parts := strings.Split(s, "_")
	ret := ""
	for _, p := range parts {
		if u := strings.ToUpper(p); commonInitialisms[u] {
			p = u
		}
		ret = ret + strings.Title(p)
	}
	return ret
}

// LogFatal logs an error as fatal with printed location and exists the program.
func LogFatal(err error, args ...interface{}) {
	if err == nil {
		return
	}
	s := "Error: " + err.Error()
	if err, ok := err.(errgo.Locationer); ok {
		s += " " + err.Location().String()
	}
	if len(args) > 0 {
		msg := args[0].(string)
		args = args[1:]
		logFatalf(s+"\n"+msg, args...)
		return
	}
	logFatalln(s)
}

// randSeq returns a random string with a defined length n.
func randSeq(n int) string {
	rand.Seed(time.Now().UTC().UnixNano())
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

// ReplaceTablePrefix replaces the {{tableprefix}} place holder with the configure real TablePrefix
// TablePrefix can be set via init() statement in config_user.go
func ReplaceTablePrefix(query string) string {
	return strings.Replace(query, "{{tableprefix}}", TablePrefix, -1)
}

func extractSplit(s string) (path string, tp []string, hasVersion bool, err error) {
	hasVersion = false

	pp := strings.Split(s, "/")
	if len(pp) == 1 {
		return path, tp, hasVersion, errgo.Newf("Path '%s' contains no slashes", s)
	}

	switch pp[0] { // add more domains here
	case "gopkg.in":
		hasVersion = true
		break
	}
	path = strings.Join(pp[:len(pp)-1], "/")
	tp = strings.Split(pp[len(pp)-1:][0], ".")

	if len(tp) == 1 {
		return path, tp, hasVersion, errgo.Newf("Missing . in package.Type")
	}
	return path, tp, hasVersion, nil
}

// ExtractImportPath extracts from an extended import path with a function or type call
// the import path.
// github.com/corestoreio/csfw/customer.Customer() would become
// github.com/corestoreio/csfw/customer
func ExtractImportPath(s string) (string, error) {

	if s == "" {
		return "", nil
	}
	path, tp, hasVersion, err := extractSplit(s)
	if err != nil {
		return "", err
	}

	path = path + "/" + tp[0]
	if hasVersion {
		path = path + "." + strings.Join(tp[1:len(tp)-1], ".")
	}

	return path, nil
}

// ExtractFuncType extracts from an extended import path with a function or type call
// the function or type call.
// github.com/corestoreio/csfw/customer.Customer() would become customer.Customer()
func ExtractFuncType(s string) (string, error) {
	if s == "" {
		return "", nil
	}

	_, tp, hasVersion, err := extractSplit(s)
	if err != nil {
		return "", err
	}

	if hasVersion {
		return strings.Join(append(tp[:1], tp[len(tp)-1:]...), "."), nil // cut the version out
	}
	return strings.Join(tp, "."), nil
}

// ParseString text/template for a string. Fails on error.
func ParseString(tpl string, data interface{}) string {
	codeTpl := template.Must(template.New("tpl_code").Parse(tpl))
	var buf = &bytes.Buffer{}
	err := codeTpl.Execute(buf, data)
	if err != nil {
		LogFatal(errgo.Mask(err))
	}
	return buf.String()
}

// Copyright (c) 2013 The Go Authors. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file or at
// https://developers.google.com/open-source/licenses/bsd.

// commonInitialisms is a set of common initialisms.
// Only add entries that are highly unlikely to be non-initialisms.
// For instance, "ID" is fine (Freudian code is rare), but "AND" is not.
var commonInitialisms = map[string]bool{
	"API":   true,
	"ASCII": true,
	"CPU":   true,
	"CSS":   true,
	"DNS":   true,
	"EOF":   true,
	"GUID":  true,
	"HTML":  true,
	"HTTP":  true,
	"HTTPS": true,
	"ID":    true,
	"IP":    true,
	"JSON":  true,
	"LHS":   true,
	"QPS":   true,
	"RAM":   true,
	"RHS":   true,
	"RPC":   true,
	"SLA":   true,
	"SMTP":  true,
	"SSH":   true,
	"TCP":   true,
	"TLS":   true,
	"TTL":   true,
	"UDP":   true,
	"UI":    true,
	"UID":   true,
	"UUID":  true,
	"URI":   true,
	"URL":   true,
	"UTF8":  true,
	"VM":    true,
	"XML":   true,
	"XSRF":  true,
	"XSS":   true,
	// CoreStore specific
	"CS":  true,
	"TMP": true,
	"IDX": true,
	"EAV": true,
}

// LintName returns a different name if it should be different.
// @see github.com/golang/lint/lint.go
func LintName(name string) (should string) {
	// Fast path for simple cases: "_" and all lowercase.
	if name == "_" {
		return name
	}
	allLower := true
	for _, r := range name {
		if !unicode.IsLower(r) {
			allLower = false
			break
		}
	}
	if allLower {
		return name
	}

	// Split camelCase at any lower->upper transition, and split on underscores.
	// Check each word for common initialisms.
	runes := []rune(name)
	w, i := 0, 0 // index of start of word, scan
	for i+1 <= len(runes) {
		eow := false // whether we hit the end of a word
		if i+1 == len(runes) {
			eow = true
		} else if runes[i+1] == '_' {
			// underscore; shift the remainder forward over any run of underscores
			eow = true
			n := 1
			for i+n+1 < len(runes) && runes[i+n+1] == '_' {
				n++
			}

			// Leave at most one underscore if the underscore is between two digits
			if i+n+1 < len(runes) && unicode.IsDigit(runes[i]) && unicode.IsDigit(runes[i+n+1]) {
				n--
			}

			copy(runes[i+1:], runes[i+n+1:])
			runes = runes[:len(runes)-n]
		} else if unicode.IsLower(runes[i]) && !unicode.IsLower(runes[i+1]) {
			// lower->non-lower
			eow = true
		}
		i++
		if !eow {
			continue
		}

		// [w,i) is a word.
		word := string(runes[w:i])
		if u := strings.ToUpper(word); commonInitialisms[u] {
			// Keep consistent case, which is lowercase only at the start.
			if w == 0 && unicode.IsLower(runes[w]) {
				u = strings.ToLower(u)
			}
			// All the common initialisms are ASCII,
			// so we can replace the bytes exactly.
			copy(runes[w:], []rune(u))
		} else if w > 0 && strings.ToLower(word) == word {
			// already all lowercase, and not the first word, so uppercase the first character.
			runes[w] = unicode.ToUpper(runes[w])
		}
		w = i
	}
	return string(runes)
}

// END Copyright (c) 2013 The Go Authors.
