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

package signed

import "net/http"

const (
	ContentSignature = "Content-Signature"
	ContentHMAC      = "Content-Hmac"
)

type (
	// EncodeFn encodes a raw signature byte slice to a string. Useful types are
	// hex.EncodeToString or base64.StdEncoding.EncodeToString.
	EncodeFn func(src []byte) string
	// DecodeFn decodes a raw signature from the header to a byte slice. Useful
	// types are hex.DecodeString or base64.StdEncoding.DecodeString.
	DecodeFn func(s string) ([]byte, error)
)

// HTTPEncoder writes a signature to the HTTP header.
type HTTPWriter interface {
	HeaderKey() string
	Write(w http.ResponseWriter, signature []byte)
}

// HTTPDecoder reads from a response the necessary data to find the signature
// and returns the signatures byte slice for further validation.
type HTTPParser interface {
	HeaderKey() string
	Parse(r *http.Request) (signature []byte, err error)
}
