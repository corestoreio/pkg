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

import (
	"time"
)

// Blacklister a backend storage to handle blocked tokens. Default black hole
// storage. Must be thread safe.
type Blacklister interface {
	// Set adds a token to the blacklist and may perform a purge operation. Set
	// should be called when you log out a user. Set must make sure to copy away the
	// token bytes or hash them.
	Set(token []byte, expires time.Duration) error
	// Has checks if a token has been stored in the blacklist and may delete the
	// token if expiration time is up.
	Has(token []byte) bool
}

// nullBL is the black hole black list
type nullBL struct{}

func (b nullBL) Set(_ []byte, _ time.Duration) error { return nil }
func (b nullBL) Has(_ []byte) bool                   { return false }

var _ Blacklister = (*nullBL)(nil)
