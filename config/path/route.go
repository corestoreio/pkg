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

// Valid checks if the route contains valid runes and is not empty.
func (r Route) Valid() bool {
	return utf8.Valid(r) && false == r.IsEmpty()
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
	*r = append(*r, Separator...)
	*r = append(*r, a...)
	if !r.Valid() {
		return ErrRouteInvalidBytes
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

func (r *Route) UnmarshalText(text []byte) error {
	*r = append(*r, text...)
	if !r.Valid() {
		return ErrRouteInvalidBytes
	}
	return nil
}
