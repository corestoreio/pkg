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

package strs

import (
	"math/rand"
	"strings"
	"time"
	"unicode"
)

// IsAlNum returns true if an alpha numeric string consists of characters a-zA-Z0-9_
func IsAlNum(s string) bool {
	c := 0
	for _, r := range s {
		switch {
		case '0' <= r && r <= '9':
			c++
			break
		case 'a' <= r && r <= 'z':
			c++
			break
		case 'A' <= r && r <= 'Z':
			c++
			break
		case r == '_':
			c++
			break
		}
	}
	return len(s) == c
}

// ToCamelCase converts from underscore separated form to camel case form.
// Eg.: catalog_product_entity -> CatalogProductEntity
func ToCamelCase(s string) string {
	return strings.Replace(strings.Title(strings.Replace(strings.ToLower(s), "_", " ", -1)), " ", "", -1)
}

// ToGoCamelCase transforms from snake case to CamelCase. Also removes
// quotes and takes care of special names.
// 		idx_eav_id 				=> IDXEAVID
// 		hello_gopher_id 		=> HelloGopherID
//		catalog_product_entity 	=> CatalogProductEntity
func ToGoCamelCase(s string) string {
	s = strings.Map(func(r rune) rune {
		switch {
		case r >= 'A' && r <= 'Z', r >= 'a' && r <= 'z', r >= '0' && r <= '9':
			return r
		}
		return '_'
	}, s)

	//	s = strings.ToLower(strings.Replace(s, `"`, "", -1))
	s = strings.ToLower(s)
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

// FromCamelCase converts from camel case form to underscore separated form.
// Eg.: CatalogProductEntity -> catalog_product_entity
func FromCamelCase(str string) string {
	var newstr []rune
	firstTime := true

	for _, chr := range str {
		isUpper := 'A' <= chr && chr <= 'Z'
		if false == isUpper {
			firstTime = false
		}

		if isUpper {
			if firstTime == true {
				firstTime = false
			} else {

				newstr = append(newstr, '_')

			}
			chr -= ('A' - 'a')
		}
		newstr = append(newstr, chr)
	}

	return string(newstr)
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

// RandAlnum returns a random string with a defined length n of alpha numerical characters.
// This function does not use the global rand variable from math/rand package.
func RandAlnum(n int) string {
	rnd := rand.New(rand.NewSource(time.Now().UTC().UnixNano()))
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rnd.Intn(len(letters))]
	}
	return string(b)
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

// CutToLength cuts a long string at maxLength. If shorter, nothing changes.
func CutToLength(data string, maxLength int) string {
	if len(data) > maxLength {
		return data[:maxLength]
	}
	return data
}
