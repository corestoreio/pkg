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

package element

import "encoding"

// LongText avoids storing long string values in Labels, Comments, Hints ...
// in Section, Group or Field elements. A LongText can contain HTML.
// The byte slice should reduce copying long strings because we're only
// copying the slice header.
type LongText []byte

func (lt LongText) String() string {
	return string(lt)
}

func (lt LongText) IsEmpty() bool {
	return lt == nil || len(lt) == 0
}

func (lt LongText) Copy() []byte {
	n := make([]byte, len(lt))
	copy(n, lt)
	return n
}

var _ encoding.TextMarshaler = (*LongText)(nil)
var _ encoding.TextUnmarshaler = (*LongText)(nil)

// MarshalText transforms the byte slice into a text slice.
func (lt LongText) MarshalText() (text []byte, err error) {
	// this is magic in combination with json.Marshal ;-)
	return lt, nil
}

func (lt *LongText) UnmarshalText(text []byte) error {
	*lt = append(*lt, text...)
	return nil
}
