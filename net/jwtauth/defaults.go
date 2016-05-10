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

package jwtauth

import (
	"time"

	"github.com/corestoreio/csfw/util/csjwt"
)

// DefaultSigningMethod HMAC-SHA signing with 512 bits. Gets applied if the
// ConfigSigningMethod model returns an empty string.
const DefaultSigningMethod = csjwt.HS512

// DefaultExpire duration when a token expires
const DefaultExpire = time.Hour

// DefaultSkew duration of time skew we allow between signer and verifier.
const DefaultSkew = time.Minute * 2
