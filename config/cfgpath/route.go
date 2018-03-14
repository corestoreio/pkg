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

package cfgpath

import (
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/corestoreio/errors"
)

// Route consists of at least three parts each of them separated by a slash
// (See constant Separator). A route can be seen as a tree.
// Route example: catalog/product/scope or websites/1/catalog/product/scope
type Route struct {
	Data  string
	Valid bool
}

// SelfRouter represents a crappy work around because of omitempty in json.Marshal
// does not work with empty non-pointer structs.
// See isEmptyValue() in encoding/json/encode.go around line 278.
// Only used in element.Field.ConfigPath, at the moment.
type SelfRouter interface {
	SelfRoute() Route
}

// MakeRoute creates a new rout from a path.
func MakeRoute(path string) Route {
	return Route{
		Data:  path,
		Valid: true,
	}
}

// SelfRoute implements the SelfRouter interface. See description of the
// RouteSelfer type.
func (r Route) SelfRoute() Route {
	return r
}

func (r Route) String() string {
	if !r.Valid {
		return "<nil>"
	}
	return r.Data
}

// Validate checks if the route contains valid runes and is not empty.
// Error behaviour: Empty and NotValid.
func (r Route) Validate() error {
	if !r.Valid || len(r.Data) == 0 {
		return errRouteEmpty
	}

	if r.Separators() == len(r.Data) {
		return errors.NotValid.Newf(errIncorrectPathTpl, r.Data)
	}

	if !utf8.ValidString(r.Data) {
		return errors.NotValid.Newf(errRouteInvalidBytesTpl, r.Data)
	}
	for _, rn := range r.Data {
		switch {
		case rn == rune(Separator), rn == '_':
			// ok
		case unicode.IsDigit(rn), unicode.IsLetter(rn), unicode.IsNumber(rn):
			// ok
		default:
			return errors.NotValid.Newf("[cfgpath] Route %q contains invalid character %q.", r.Data, rn)
		}
	}
	return nil
}

// Equal compares the Sum32 of both routes
func (r Route) Equal(b Route) bool {
	return r.Valid == b.Valid && r.Data == b.Data
}

// Append adds other partial routes with a Separator between. After the partial
// routes have been added a validation check will NOT be done.
// Internally it creates a new byte slice.
//
//		a := cfgpath.Route(`catalog/product`)
//		b := cfgpath.Route(`enable_flat_tables`)
//		if err := a.Append(b); err != nil {
//			panic(err)
//		}
//		println(a.String())
//		// Should print: catalog/product/enable_flat_tables
func (r Route) Append(routes ...Route) Route {

	if i1, i2 := strings.LastIndexByte(r.Data, Separator), len(r.Data)-1; i1 > 0 && i2 > 0 && i1 == i2 {
		r.Data = r.Data[:len(r.Data)-1] // strip last Separator
	}

	var buf strings.Builder
	if r.Valid {
		buf.WriteString(r.Data)
	}

	i := 0
	for _, route := range routes {
		if route.Valid {
			if len(r.Data) > 0 && len(route.Data) > 0 && route.Data[0] != Separator {
				buf.WriteByte(Separator)
			}
			buf.WriteString(route.Data)
			i++
		}
	}

	return Route{
		Data:  buf.String(),
		Valid: true,
	}
}

// MarshalText implements interface encoding.TextMarshaler.
func (r Route) MarshalText() (text []byte, err error) {
	return []byte(r.Data), nil
}

// UnmarshalText transforms the text into a route with performed validation
// checks. Implements encoding.TextUnmarshaler.
// Error behaviour: NotValid, Empty.
func (r *Route) UnmarshalText(txt []byte) error {
	r.Data = string(txt)
	r.Valid = true
	err := r.Validate()
	r.Valid = err == nil
	return err
}

// Level returns a hierarchical based route depending on the depth.
// The depth argument defines the depth of levels to be returned.
// Depth 1 will return the first part like "a", Depth 2 returns "a/b"
// Depth 3 returns "a/b/c" and so on. Level -1 gives you all available levels.
// Does not generate a fully qualified path.
// The returned Route slice is owned by Route. For further modifications you must
// copy it via Route.Copy().
// Error behaviour: NotValid, Empty.
func (r Route) Level(depth int) (Route, error) {
	if err := r.Validate(); err != nil {
		return Route{}, errors.Wrap(err, "[cfgpath] Level.Validate")
	}

	lp := len(r.Data)
	switch {
	case depth < 0:
		return r, nil
	case depth == 0:
		return Route{}, nil
	case depth >= lp:
		return r, nil
	}

	pos := 0
	i := 1
	for pos <= lp {
		sc := strings.IndexByte(r.Data[pos:], Separator)
		if sc == -1 {
			return r, nil
		}
		pos += sc + 1

		if i == depth {
			return Route{Data: r.Data[:pos-1], Valid: true}, nil
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
// a []byte.
//
// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// Hash implements FNV-1 and FNV-1a, non-cryptographic hash functions
// created by Glenn Fowler, Landon Curt Noll, and Phong Vo.
// See
// http://en.wikipedia.org/wiki/Fowler%E2%80%93Noll%E2%80%93Vo_hash_function.
//
// Usage as map key.
func (r Route) Hash(depth int) (uint32, error) {
	r2, err := r.Level(depth)
	if err != nil {
		return 0, errors.Wrap(err, "[cfgpath] Hash.Level")
	}
	var hash uint32 = offset32
	for _, c := range r2.Data {
		hash ^= uint32(c)
		hash *= prime32
	}
	return hash, nil
}

// Hash32 returns a fnv32a value of the route. 32 has been chosen because
// routes consume not that much space as a []byte.
//
// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// Hash implements FNV-1 and FNV-1a, non-cryptographic hash functions
// created by Glenn Fowler, Landon Curt Noll, and Phong Vo.
// See
// http://en.wikipedia.org/wiki/Fowler%E2%80%93Noll%E2%80%93Vo_hash_function.
//
// Usage as map key.
func (r Route) Hash32() uint32 {
	var hash uint32 = offset32
	for _, c := range r.Data {
		hash ^= uint32(c)
		hash *= prime32
	}
	return hash
}

// Separators returns the number of separators
func (r Route) Separators() (count int) {
	for _, b := range r.Data {
		if b == rune(Separator) {
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
		return Route{}, errors.NotValid.Newf(errIncorrectPositionTpl, pos)
	}

	sepCount := r.Separators()
	if sepCount < 1 { // no separator found
		return r, nil
	}
	if pos > sepCount+1 {
		return Route{}, errors.NotValid.Newf(errIncorrectPositionTpl, pos)
	}

	var sepPos [maxLevels]int
	sp := 0
	for i, b := range r.Data {
		if b == rune(Separator) && sp < maxLevels {
			sepPos[sp] = i + 1 // positions of the separators in the slice
			sp++
		}
	}

	pos--
	min := 0
	for i := 0; i < maxLevels; i++ {
		if sepPos[i] == 0 { // no more separators found
			return Route{Data: r.Data[min:], Valid: true}, nil
		}
		max := sepPos[i]
		if i == pos {
			return Route{Data: r.Data[min : max-1], Valid: true}, nil
		}
		min = max
	}
	return Route{Data: r.Data[min:], Valid: true}, nil
}

// Split splits the route into its three parts. Does not run Validate()
// Example:
// 		routes := cfgpath.MakeRoute("aa/bb/cc")
//		rs, err := routes.Split()
//		rs[0].String() == "aa"
//		rs[1].String() == "bb"
//		rs[2].String() == "cc"
//
// Error behaviour: NotValid
func (r Route) Split() (ret [Levels]Route, err error) {

	const sepCount = Levels - 1 // only two separators supported
	var sepPos [sepCount]int
	sp := 0
	for i, b := range r.Data {
		if b == rune(Separator) && sp < sepCount {
			sepPos[sp] = i // positions of the separators in the slice
			sp++
		}
	}
	if sp < 1 {
		err = errors.NotValid.Newf(errIncorrectPathTpl, r.Data)
		return
	}

	min := 0
	for i := 0; i < Levels; i++ {
		var max int
		if i < sepCount && sepPos[i] > 0 {
			max = sepPos[i]
		} else {
			max = len(r.Data)
		}
		ret[i] = Route{Data: r.Data[min:max], Valid: true}
		if i < sepCount && sepPos[i] == 0 {
			return
		}
		min = max + 1
	}
	return
}
