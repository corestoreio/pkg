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

//go:generate go run ../internal/scopedservice/main_copy.go "$GOPACKAGE"

package signed

// Service creates a middleware that facilitates using a hash function to sign a
// HTTP body and validate the HTTP body of a request.
type Service struct {
	service
}

// New creates a new signing middleware for signature creation and validation.
// The scope.Default and any other scopes have these default settings: InTrailer
// activated, Content-HMAC header with sha256, allowed HTTP methods set to POST,
// PUT, PATCH and password for the HMAC SHA 256 from a cryptographically random
// source with a length of 64 bytes.
func New(opts ...Option) (*Service, error) {
	return newService(opts...)
}
