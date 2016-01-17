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
	"encoding"
	"errors"
	"fmt"
	"unicode/utf8"
)

// ErrRouteInvalidBytes whenever a non-rune is detected.
var ErrRouteInvalidBytes = errors.New("Route contains invalid bytes which are not runes.")

// Route consists of at least three parts each of them separated by a slash
// (See constant Separator). A route can be seen as a tree.
// Route example: catalog/product/scope or websites/1/catalog/product/scope
type Route []byte

// newRoute creates a new rout from sub paths resp. path parts.
// Parts gets merged via Separator
func newRoute(parts ...string) (Route, error) {
	l := 0
	for _, p := range parts {
		l += len(p)
		l += 1 // len(sSeparator)
	}
	l -= 1 // remove last slash
	if l < 1 {
		return nil, ErrRouteEmpty
	}

	r := make(Route, l)
	pos := 0
	for i, p := range parts {
		pos += copy(r[pos:pos+len(p)], p)
		if i < len(parts)-1 {
			pos += copy(r[pos:pos+1], sSeparator)
		}
	}
	return r, nil
}

func (r Route) String() string {
	return string(r)
}

func (r Route) GoString() string {
	return fmt.Sprintf("path.Route(`%s`)", r)
}

const rSeparator = rune(Separator)

// Validate checks if the route contains valid runes and is not empty.
func (r Route) Validate() error {

	if r.IsEmpty() {
		return ErrRouteEmpty
	}

	if r.Separators() == len(r) {
		return ErrIncorrectPath
	}

	if false == utf8.Valid(r) {
		return ErrRouteInvalidBytes
	}

	var sepCount, length int
	i := 0
	for i < len(r) {
		var ru rune
		if r[i] < utf8.RuneSelf {
			ru = rune(r[i])
			i++
		} else {
			dr, _ := utf8.DecodeRune(r[i:])
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
		length++
	}
	return nil
}

func (r Route) IsEmpty() bool {
	return r == nil || len(r) == 0
}

func (r Route) Copy() []byte {
	n := make([]byte, len(r))
	copy(n, r)
	return n
}

// Append adds another partial route with a Separator between. After the partial
// route has been added a validation check will be done.
//
//		a := path.Route(`catalog/product`)
//		b := path.Route(`enable_flat_tables`)
//		if err := a.Append(b); err != nil {
//			panic(err)
//		}
//		println(a.String())
//		// Should print: catalog/product/enable_flat_tables
func (r *Route) Append(a Route) error {
	*r = append(*r, Separator)
	*r = append(*r, a...)
	if err := r.Validate(); err != nil {
		return err
	}
	return nil
}

var _ encoding.TextMarshaler = (*Route)(nil)
var _ encoding.TextUnmarshaler = (*Route)(nil)

// MarshalText transforms the byte slice into a text slice.
func (r Route) MarshalText() (text []byte, err error) {
	// this is magic in combination with json.Marshal ;-)
	return r, nil
}

// UnmarshalText transforms the text into a route with performed validation
// checks.
func (r *Route) UnmarshalText(text []byte) error {
	*r = append(*r, text...)
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
func (r Route) Level(level int) (Route, error) {
	if err := r.Validate(); err != nil {
		return nil, err
	}

	lp := len(r)
	switch {
	case level < 0:
		return r, nil
	case level == 0:
		return r[:0], nil
	case level >= lp:
		return r, nil
	}

	pos := 0
	i := 1
	for pos <= lp {
		sc := bytes.IndexByte(r[pos:], Separator)
		if sc == -1 {
			return r, nil
		}
		pos += sc + 1

		if i == level {
			return r[:pos-1], nil
		}
		i++
	}
	return r, nil
}

const (
	offset64 = 14695981039346656037
	prime64  = 1099511628211
)

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
	r, err := r.Level(level)
	if err != nil {
		return 0, err
	}
	var hash uint64 = offset64
	for _, c := range r {
		hash ^= uint64(c)
		hash *= prime64
	}
	return hash, nil
}

// Separators returns the number of separators
func (r Route) Separators() (count int) {
	for _, b := range r {
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
func (r Route) Part(pos int) (Route, error) {

	if err := r.Validate(); err != nil {
		return nil, err
	}

	if pos < 1 {
		return nil, ErrIncorrectPosition
	}

	sepCount := r.Separators()
	if sepCount < 1 { // no separator found
		return r, nil
	}
	if pos > sepCount+1 {
		return nil, ErrIncorrectPosition
	}

	const realLevels = Levels - 1
	var sepPos [realLevels]int
	sp := 0
	for i, b := range r {
		if b == Separator && sp < realLevels {
			sepPos[sp] = i + 1 // positions of the separators in the slice
			sp++
		}
	}

	pos -= 1
	min := 0
	for i := 0; i < realLevels; i++ {
		max := sepPos[i]
		if i == pos {
			return r[min : max-1], nil
		}
		min = max
	}
	return r[min:], nil
}
