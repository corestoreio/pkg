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

import (
	"encoding/hex"
	"net/http"
	"strings"

	"github.com/corestoreio/errors"
)

// ContentHMAC writes a simple Content-HMAC header. It can additionally parse a request
// and return the raw signature.
type ContentHMAC struct {
	// Algorithm parameter is used if the client and server agree on a
	// non-standard digital signature algorithm.  The full list of supported
	// signature mechanisms is listed below. REQUIRED.
	Algorithm string
	// HeaderKey (optional) a field name in the HTTP header, defaults to
	// Content-HMAC.
	HeaderName string
	// EncodeFn (optional) defines the byte to string encoding function.
	// Defaults to hex.EncodeString.
	EncodeFn
	// DecodeFn (optional) defines the string to byte decoding function.
	// Defaults to hex.DecodeString.
	DecodeFn
}

// NewContentHMAC creates a new header HMAC object with default hex encoding/decoding
// to write and parse the Content-HMAC field.
func NewContentHMAC(algorithm string) *ContentHMAC {
	return &ContentHMAC{
		Algorithm: algorithm,
	}
}

// HeaderKey returns the name of the header key
func (h *ContentHMAC) HeaderKey() string {
	if h.HeaderName != "" {
		return h.HeaderName
	}
	return HeaderContentHMAC
}

// Writes writes the signature into the response.
// Content-HMAC: <hash mechanism> <encoded binary HMAC>
// Content-HMAC: sha1 f1wOnLLwcTexwCSRCNXEAKPDm+U=
func (h *ContentHMAC) Write(w http.ResponseWriter, signature []byte) {
	encFn := h.EncodeFn
	if encFn == nil {
		encFn = hex.EncodeToString
	}
	w.Header().Set(h.HeaderKey(), h.Algorithm+" "+encFn(signature))
}

// Parse looks up the header or trailer for the HeaderKey Content-HMAC in an
// HTTP request and extracts the raw decoded signature. Errors can have the
// behaviour: NotFound or NotValid.
func (h *ContentHMAC) Parse(r *http.Request) (signature []byte, _ error) {
	hk := h.HeaderKey()
	hv := r.Header.Get(hk)
	if hv == "" {
		hv = r.Trailer.Get(hk)
	}
	if hv == "" {
		return nil, errors.NewNotFoundf(errHMACParseNotFound)
	}

	firstWS := strings.IndexByte(hv, ' ') // first white space after algorithm name
	if hv == "" || firstWS != len(h.Algorithm) {
		return nil, errors.NewNotValidf(errHMACParseNotValid, hv, hk)
	}
	if h.Algorithm == "" || h.Algorithm != hv[:firstWS] {
		return nil, errors.NewNotValidf(errHMACParseInvalidAlg, hv[:firstWS], hk, hv)
	}

	decFn := h.DecodeFn
	if decFn == nil {
		decFn = hex.DecodeString
	}
	dec, err := decFn(hv[firstWS+1:])
	if err != nil {
		// micro optimization: skip argument building
		return nil, errors.NewNotValidf("[signed] HMAC failed to decode: %q in header %q. Error: %s", hv[firstWS+1:], hv, err)
	}
	return dec, nil
}
