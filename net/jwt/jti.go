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

package jwt

import "github.com/corestoreio/csfw/util/shortid"

// JTI represents the interface to generate a new UUID aka JWT ID
type IDGenerator interface {
	NewID() (string, error)
}

// jti type to generate a JTI for a token, a unique ID
type jti struct{}

func (j jti) NewID() (string, error) {
	return shortid.Generate()
}
