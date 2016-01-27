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

// Route consists of at least three parts each of them separated by a slash
// (See constant Separator). A route can be seen as a tree.
// Route example: catalog/product/scope or websites/1/catalog/product/scope
type Route struct {
	// Sum32 is a fnv 32a hash for comparison and maybe later integrity checks.
	// Sum32 will be automatically updated when using New*() functions.
	// If you set yourself Chars then update Sum32 on your own with Hash32().
	Sum32 uint32
	text.Chars
}

// RouteSelfer represents a crappy work around because of omitempty in json.Marshal
// does not work with empty non-pointer structs.
// See isEmptyValue() in encoding/json/encode.go around line 278.
// Only used in element.Field.ConfigPath, at the moment.
type RouteSelfer interface {
	Self() Route
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

	if l < 1 {
		return Route{}
	}

	c := make(text.Chars, l, l)
	pos := 0
	for i, p := range parts {
		pos += copy(c[pos:pos+len(p)], p)
		if i < len(parts)-1 {
			pos += copy(c[pos:pos+1], sSeparator)
		}
	}
	return newRoute(c)
}

func newRoute(b []byte) Route {
	r := Route{
		Chars: b,
	}
	r.Sum32 = r.Hash32()
	return r
}

// Route implements the RouteSelfer interface. See description of the
// RouteSelfer type.
func (r Route) Self() Route {
	return r
}

// GoString returns the Go type of the Route including the underlying bytes.
func (r Route) GoString() string {
	if r.IsEmpty() {
		return "path.Route{}"
	}
	return fmt.Sprintf("path.Route{Chars:[]byte(`%s`)}", r)
}

const rSeparator = rune(Separator)

// Validate checks if the route contains valid runes and is not empty.
func (r Route) Validate() error {
	if r.IsEmpty() {
		return ErrRouteEmpty
	}

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

// Equal compares the Sum32 of both routes
func (r Route) Equal(b Route) bool {
	return r.Sum32 == b.Sum32
	//if r.Sum32 == b.Sum32 {
	//	return true
	//}
	//return r.Chars.Equal(b.Chars) // takes longer
}

// Clone returns a new allocated route with copied data.
func (r Route) Clone() Route {
	return newRoute(r.Chars.Clone())
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

	if i1, i2 := bytes.LastIndexByte((*r).Chars, Separator), len((*r).Chars)-1; i1 > 0 && i2 > 0 && i1 == i1 {
		(*r).Chars = (*r).Chars[:len((*r).Chars)-1] // strip last Separator
	}

	// calculate new buffer size
	size := len((*r).Chars)

	rsLen := len(routes) - 1
	i := 0
	for _, route := range routes {
		if route.Chars.IsEmpty() {
			rsLen -= 1
			continue
		}
		if i == 0 {
			size++ // Separator
		}
		size += len(route.Chars)
		if len(route.Chars) > 0 && i < len(routes)-1 {
			size++ // Separator
		}
		i++
	}

	var buf = make([]byte, size, size)
	var pos int
	if len((*r).Chars) > 0 {
		pos += copy(buf[pos:], (*r).Chars)
	}

	i = 0
	for _, route := range routes {
		if route.Chars.IsEmpty() {
			continue
		}
		if i == 0 && len((*r).Chars) > 0 && route.Chars[0] != Separator {
			pos += copy(buf[pos:], bSeparator)
		}
		pos += copy(buf[pos:], route.Chars)
		if i < rsLen {
			pos += copy(buf[pos:], bSeparator)
		}
		i++
	}

	if pos := bytes.IndexByte(buf, 0x00); pos >= 1 {
		buf = buf[:pos] // strip everything after the null byte
	}

	(*r) = newRoute(buf)
	if err := r.Validate(); err != nil {
		return err
	}
	return nil
}

// UnmarshalText transforms the text into a route with performed validation
// checks.
func (r *Route) UnmarshalText(txt []byte) error {
	var c text.Chars
	if err := c.UnmarshalText(txt); err != nil {
		return err
	}
	(*r) = newRoute(c)
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
		return Route{}, err
	}

	lp := len(r.Chars)
	switch {
	case level < 0:
		return r, nil
	case level == 0:
		return Route{}, nil
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
			return newRoute(r.Chars[:pos-1]), nil
		}
		i++
	}
	return r, nil
}

const (
	offset32 = 2166136261
	prime32  = 16777619
)

// Hash same as Level() but returns a fnv32a value or an error if the route is
// invalid. 32 has been chosen because routes consume not that much space as
// a text.Chars.
//
// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// Hash implements FNV-1 and FNV-1a, non-cryptographic hash functions
// created by Glenn Fowler, Landon Curt Noll, and Phong Vo.
// See
// http://en.wikipedia.org/wiki/Fowler%E2%80%93Noll%E2%80%93Vo_hash_function.
func (r Route) Hash(level int) (uint32, error) {
	r2, err := r.Level(level)
	if err != nil {
		return 0, err
	}
	var hash uint32 = offset32
	for _, c := range r2.Chars {
		hash ^= uint32(c)
		hash *= prime32
	}
	return hash, nil
}

// Hash32 returns a fnv32a value of the route. 32 has been chosen because
// routes consume not that much space as a text.Chars.
//
// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// Hash implements FNV-1 and FNV-1a, non-cryptographic hash functions
// created by Glenn Fowler, Landon Curt Noll, and Phong Vo.
// See
// http://en.wikipedia.org/wiki/Fowler%E2%80%93Noll%E2%80%93Vo_hash_function.
func (r Route) Hash32() uint32 {
	var hash uint32 = offset32
	for _, c := range r.Chars {
		hash ^= uint32(c)
		hash *= prime32
	}
	return hash
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
// before extracting the part. Does not run Validate()
//		Have Route: general/single_store_mode/enabled
//		Pos<1 => ErrIncorrectPosition
//		Pos=1 => general
//		Pos=2 => single_store_mode
//		Pos=3 => enabled
//		Pos>3 => ErrIncorrectPosition
// The returned Route slice is owned by Path. For further modifications you must
// copy it via Route.Copy().
func (r Route) Part(pos int) (Route, error) {

	if pos < 1 {
		return Route{}, ErrIncorrectPosition
	}

	sepCount := r.Separators()
	if sepCount < 1 { // no separator found
		return r, nil
	}
	if pos > sepCount+1 {
		return Route{}, ErrIncorrectPosition
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
			return newRoute(r.Chars[min:]), nil
		}
		max := sepPos[i]
		if i == pos {
			return newRoute(r.Chars[min : max-1]), nil
		}
		min = max
	}
	return newRoute(r.Chars[min:]), nil
}

// Split splits the route into its three parts. Does not run Validate()
// Example:
// 		routes := path.NewRoute("aa/bb/cc")
//		rs, err := routes.Split()
//		rs[0].String() == "aa"
//		rs[1].String() == "bb"
//		rs[2].String() == "cc"
func (r Route) Split() (ret [Levels]Route, err error) {

	const sepCount = Levels - 1 // only two separators supported
	var sepPos [sepCount]int
	sp := 0
	for i, b := range r.Chars {
		if b == Separator && sp < sepCount {
			sepPos[sp] = i // positions of the separators in the slice
			sp++
		}
	}
	if sp < 1 {
		err = ErrIncorrectPath
		return
	}

	min := 0
	for i := 0; i < Levels; i++ {
		var max int
		if i < sepCount && sepPos[i] > 0 {
			max = sepPos[i]
		} else {
			max = len(r.Chars)
		}
		ret[i] = newRoute(r.Chars[min:max])
		if i < sepCount && sepPos[i] == 0 {
			return
		}
		min = max + 1
	}
	return
}
