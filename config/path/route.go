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
)

// ErrRouteInvalidBytes whenever a non-rune is detected.
var ErrRouteInvalidBytes = errors.New("Route contains invalid bytes which are not runes.")

// Route consists of at least three parts each of them separated by a slash
// (See constant Separator). A route can be seen as a tree.
// Route example: catalog/product/scope or websites/1/catalog/product/scope
type Route []byte

func (r Route) String() string {
	return string(r)
}

func (r Route) IsEmpty() bool {
	return r == nil || len(r) == 0
}

func (r Route) Copy() []byte {
	n := make([]byte, len(r))
	copy(n, r)
	return n
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
	return nil
}
