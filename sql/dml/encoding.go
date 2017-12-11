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

package dml

// This file contains all interfaces and function signatures for the various
// encoders and decoders. You must set them in your package before using the
// encoding/decoding methods.

// JSONMarshalFn and JSONUnMarshalFn functions which must be set if you decide
// to use JSON. Otherwise it panics. In your package write somewhere:
//		dml.JSONMarshalFn = json.Marshal
//		dml.JSONUnMarshalFn = json.UnMarshal
var (
	JSONMarshalFn   func(v interface{}) ([]byte, error)
	JSONUnMarshalFn func(data []byte, v interface{}) error
)
