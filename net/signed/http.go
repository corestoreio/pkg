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

package signed

import "net/http"

// Content* constants are used as HTTP header key names.
const (
	HeaderContentSignature = "Content-Signature"
	HeaderContentHMAC      = "Content-Hmac"
)

// EncodeFn encodes a raw signature byte slice to a string. Useful types are
// hex.EncodeToString or base64.StdEncoding.EncodeToString.
type EncodeFn func(src []byte) string

// DecodeFn decodes a raw signature from the header to a byte slice. Useful
// types are hex.DecodeString or base64.StdEncoding.DecodeString.
type DecodeFn func(s string) ([]byte, error)

// HeaderParseWriter knows how to read and write the HTTP header in regards to
// the hash.
type HeaderParseWriter interface {
	HeaderKey() string
	// Write writes a signature to the HTTP response header.
	Write(w http.ResponseWriter, signature []byte)
	// Parse parses from a request the necessary data to find the signature hash
	// and returns the raw signatures byte slice for further validation.
	Parse(r *http.Request) (signature []byte, err error)
}
