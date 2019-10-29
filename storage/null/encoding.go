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

package null

import (
	"sync/atomic"

	"github.com/corestoreio/errors"
)

// This file contains all interfaces and function signatures for the various
// encoders and decoders. You must set them in your package before using the
// encoding/decoding methods.

// jsonMarshalFn and jsonUnMarshalFn functions which must be set if you decide
// to use JSON. Otherwise it panics. In your package write somewhere:
//		null.jsonMarshalFn = json.Marshal
//		null.jsonUnMarshalFn = json.UnMarshal
var (
	jsonMarshalFn   func(v interface{}) ([]byte, error)
	jsonUnMarshalFn func(data []byte, v interface{}) error
	jsonFnApplied   = new(int32)
)

// MustSetJSONMarshaler applies the global JSON marshal functions for encoding and
// decoding. This function can only be called once. It panics on multiple calls.
func MustSetJSONMarshaler(marshalFn func(v interface{}) ([]byte, error), unMarshalFn func(data []byte, v interface{}) error) {
	if atomic.LoadInt32(jsonFnApplied) == 1 {
		panic(errors.AlreadyExists.Newf("[null] JSON marshal and unmarshal already exists."))
	}
	atomic.StoreInt32(jsonFnApplied, 1)
	jsonMarshalFn = marshalFn
	jsonUnMarshalFn = unMarshalFn
}
