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

package jwt

import (
	"time"
)

// Blocklister a backend storage to handle blocked tokens. Default noop
// storage. Must be thread safe.
type Blocklister interface {
	// Set adds the token ID (The jti (JWT ID) claim provides a unique
	// identifier for the JWT) to the blocklist and may perform a purge
	// operation. Set should be called when you log out a user. Set must make
	// sure to copy away the bytes or hash them.
	Set(id []byte, expires time.Duration) error
	// Has checks if an ID (jti) has been stored in the blockList and may delete
	// the ID if the expiration time is up.
	Has(id []byte) bool
}

// nullBL is a noop block list
type nullBL struct{}

func (b nullBL) Set(_ []byte, _ time.Duration) error { return nil }
func (b nullBL) Has(_ []byte) bool                   { return false }

var _ Blocklister = (*nullBL)(nil)
