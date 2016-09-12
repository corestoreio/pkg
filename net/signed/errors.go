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

const (
	errScopedConfigNotValid         = `[signed] ScopedConfig %s is invalid. IsNil(HeaderParseWriter=%t) AllowedMethods: %v`
	errScopedConfigMethodNotAllowed = `[signed] ValidateBody HTTP Method %q not allowed in list: %q`
	errScopedConfigSignatureNoMatch = `[signed] ValidateBody. Signatures do not match. Have: %q Want: %q`
	errScopedConfigCacheNotFound    = `[signed] ValidateBody. Signature %q not found in cache`
	errSignatureParseNotFound       = `[signed] Signature not found or empty`
	errSignatureParseInvalidHeader  = `[signed] Invalid signature header: %q`
	errSignatureParseInvalidKeyID   = `[signed] KeyID %q does not match required %q in header: %q`
	errSignatureParseInvalidAlg     = `[signed] Algorithm %q does not match required %q in header: %q`
	errHMACParseNotFound            = `[signed] Signature not found or empty`
	errHMACParseNotValid            = `[signed] Signature %q not valid in header %q`
	errHMACParseInvalidAlg          = `[signed] Unknown algorithm %q in Header %q with signature %q`
)
