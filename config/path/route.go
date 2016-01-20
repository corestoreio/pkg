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

package path

import (
	"bytes"
	"errors"
	"fmt"
	"unicode/utf8"

	"github.com/corestoreio/csfw/storage/text"
)

// ErrRouteInvalidBytes whenever a non-rune is detected.
var ErrRouteInvalidBytes = errors.New("Route contains invalid bytes which are not runes.")

// TODO(cs) immutable: maybe implement this
// var ErrRouteChanged = errors.New("Route bytes changed")

// Route consists of at least three parts each of them separated by a slash
// (See constant Separator). A route can be seen as a tree.
// Route example: catalog/product/scope or websites/1/catalog/product/scope
type Route struct {
	// TODO(cs) immutable: maybe implement this
	// org uint64 // fnv hash to check if the byte slice has changed
	text.Chars
}

// NewRoute creates a new rout from sub paths resp. path parts.
// Parts gets merged via Separator
func NewRoute(parts ...string) Route {
	l := 0
	for _, p := range parts {
		l += len(p)
		l += 1 // len(sSeparator)
	}
	l -= 1 // remove last slash

	r := Route{}
	if l < 1 {
		return r
	}

	r.Chars = make(text.Chars, l, l)
	pos := 0
	for i, p := range parts {
		pos += copy(r.Chars[pos:pos+len(p)], p)
		if i < len(parts)-1 {
			pos += copy(r.Chars[pos:pos+1], sSeparator)
		}
	}
	// TODO(cs) immutable: maybe implement this
	// r.org = r.Chars.Hash()
	return r
}

// GoString returns the Go type of the Route including the underlying bytes.
func (r Route) GoString() string {
	return fmt.Sprintf("path.Route{Chars:[]byte(`%s`)}", r)
}

const rSeparator = rune(Separator)

// Validate checks if the route contains valid runes and is not empty.
func (r Route) Validate() error {
	if r.IsEmpty() {
		return ErrRouteEmpty
	}
	// TODO(cs) immutable: maybe implement this
	//if r.org != r.Chars.Hash() {
	//	return ErrRouteChanged
	//}
	if r.Separators() == len(r.Chars) {
		return ErrIncorrectPath
	}

	if false == utf8.Valid(r.Chars) {
		return ErrRouteInvalidBytes
	}

	var sepCount int
	i := 0
	for i < len(r.Chars) {
		var ru rune
		if r.Chars[i] < utf8.RuneSelf {
			ru = rune(r.Chars[i])
			i++
		} else {
			dr, _ := utf8.DecodeRune(r.Chars[i:])
			return fmt.Errorf("This character %q is not allowed in Route %s", string(dr), r)
		}
		ok := false
		switch {
		case '0' <= ru && ru <= '9':
			ok = true
		case 'a' <= ru && ru <= 'z':
			ok = true
		case 'A' <= ru && ru <= 'Z':
			ok = true
		case ru == '_':
			ok = true
		case ru == rSeparator:
			sepCount++
			ok = true
		}
		if !ok {
			return fmt.Errorf("This character %q is not allowed in Route %s", string(ru), r)
		}
	}
	return nil
}

func (r Route) Copy() Route {
	return Route{Chars: r.Chars.Copy()}
}

// Append adds other partial routes with a Separator between. After the partial
// routes have been added a validation check will be done.
//
//		a := path.Route(`catalog/product`)
//		b := path.Route(`enable_flat_tables`)
//		if err := a.Append(b); err != nil {
//			panic(err)
//		}
//		println(a.String())
//		// Should print: catalog/product/enable_flat_tables
func (r *Route) Append(routes ...Route) error {
	if bytes.LastIndexByte((*r).Chars, Separator) == len((*r).Chars)-1 {
		(*r).Chars = (*r).Chars[:len((*r).Chars)-1] // strip last Separator
	}

	// calculate new buffer size
	size := len((*r).Chars)
	for i, route := range routes {
		if i == 0 {
			size++ // Separator
		}
		size += len(route.Chars)
		if i < len(routes)-1 {
			size++ // Separator
		}
	}
	var buf = make([]byte, size, size)
	var pos int
	pos += copy(buf[pos:], (*r).Chars)

	for i, route := range routes {
		if i == 0 && route.Chars[0] != Separator {
			pos += copy(buf[pos:], bSeparator)
		}

		pos += copy(buf[pos:], route.Chars)
		if i < len(routes)-1 {
			pos += copy(buf[pos:], bSeparator)
		}
	}
	if pos := bytes.IndexByte(buf, 0x00); pos > 1 {
		buf = buf[:pos] // strip everything after the null byte
	}
	(*r).Chars = buf
	if err := r.Validate(); err != nil {
		return err
	}
	return nil
}

// UnmarshalText transforms the text into a route with performed validation
// checks.
func (r *Route) UnmarshalText(text []byte) error {
	if err := (*r).Chars.UnmarshalText(text); err != nil {
		return err
	}
	if err := r.Validate(); err != nil {
		return err
	}
	return nil
}

// Level joins a configuration path parts by the path separator PS.
// The level argument defines the depth of the path parts to join.
// Level 1 will return the first part like "a", Level 2 returns "a/b"
// Level 3 returns "a/b/c" and so on. Level -1 joins all available path parts.
// Does not generate a fully qualified path.
// The returned Route slice is owned by Path. For further modifications you must
// copy it via Route.Copy().
func (r Route) Level(level int) (ret Route, err error) {
	if err = r.Validate(); err != nil {
		return
	}

	lp := len(r.Chars)
	switch {
	case level < 0:
		return r, nil
	case level == 0:
		return
	case level >= lp:
		return r, nil
	}

	pos := 0
	i := 1
	for pos <= lp {
		sc := bytes.IndexByte(r.Chars[pos:], Separator)
		if sc == -1 {
			return r, nil
		}
		pos += sc + 1

		if i == level {
			ret.Chars = r.Chars[:pos-1]
			return
		}
		i++
	}
	return r, nil
}

// Hash same as Level() but returns a fnv64a value or an error if the route is
// invalid.
//
// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// Hash implements FNV-1 and FNV-1a, non-cryptographic hash functions
// created by Glenn Fowler, Landon Curt Noll, and Phong Vo.
// See
// http://en.wikipedia.org/wiki/Fowler%E2%80%93Noll%E2%80%93Vo_hash_function.
func (r Route) Hash(level int) (uint64, error) {
	r2, err := r.Level(level)
	if err != nil {
		return 0, err
	}
	return r2.Chars.Hash(), nil
}

// Separators returns the number of separators
func (r Route) Separators() (count int) {
	for _, b := range r.Chars {
		if b == Separator {
			count++
		}
	}
	return
}

// Part returns the route part on the desired position. The Route gets validated
// before extracting the part.
//		Have Route: general/single_store_mode/enabled
//		Pos<1 => ErrIncorrectPosition
//		Pos=1 => general
//		Pos=2 => single_store_mode
//		Pos=3 => enabled
//		Pos>3 => ErrIncorrectPosition
// The returned Route slice is owned by Path. For further modifications you must
// copy it via Route.Copy().
func (r Route) Part(pos int) (ret Route, err error) {

	if err = r.Validate(); err != nil {
		return
	}

	if pos < 1 {
		err = ErrIncorrectPosition
		return
	}

	sepCount := r.Separators()
	if sepCount < 1 { // no separator found
		return r, nil
	}
	if pos > sepCount+1 {
		err = ErrIncorrectPosition
		return
	}

	var sepPos [maxLevels]int
	sp := 0
	for i, b := range r.Chars {
		if b == Separator && sp < maxLevels {
			sepPos[sp] = i + 1 // positions of the separators in the slice
			sp++
		}
	}

	pos -= 1
	min := 0
	for i := 0; i < maxLevels; i++ {
		if sepPos[i] == 0 { // no more separators found
			ret.Chars = r.Chars[min:]
			return
		}
		max := sepPos[i]
		if i == pos {
			ret.Chars = r.Chars[min : max-1]
			return
		}
		min = max
	}
	ret.Chars = r.Chars[min:]
	return
}
