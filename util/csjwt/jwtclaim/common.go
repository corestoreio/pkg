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

package jwtclaim

import (
	"crypto/subtle"
	"time"
)

// Key constants define the main claims used for Set() and Get() functions.
// Those constants are implemented in the StandardClaims type.
const (
	KeyAudience  = "aud"
	KeyExpiresAt = "exp"
	KeyID        = "jti"
	KeyIssuedAt  = "iat"
	KeyIssuer    = "iss"
	KeyNotBefore = "nbf"
	KeySubject   = "sub"
	KeyTimeSkew  = "skew" // not marshalled, internal usage.
)

// allKeys first seven entries only supported by Standard type and all entries
// supported by Store type.
var allKeys = [9]string{KeyAudience, KeyExpiresAt, KeyID, KeyIssuedAt, KeyIssuer, KeyNotBefore, KeySubject, KeyStore, KeyUserID}

// TimeFunc provides the current time when parsing token to validate "exp" claim (expiration time).
// You can override it to use another time value.  This is useful for testing or if your
// server uses a different time zone than your tokens.
var TimeFunc = time.Now

func verifyConstantTime(aud, cmp []byte, required bool) bool {
	if len(aud) == 0 {
		return !required
	}
	return subtle.ConstantTimeCompare(aud, cmp) == 1
}

func verifyExp(skew time.Duration, exp, now int64, required bool) bool {
	if exp == 0 {
		return !required
	}
	now -= int64(skew.Seconds())
	return now <= exp
}

func verifyIat(iat int64, now int64, required bool) bool {
	if iat == 0 {
		return !required
	}
	return now >= iat
}

func verifyNbf(skew time.Duration, nbf int64, now int64, required bool) bool {
	if nbf == 0 {
		return !required
	}
	now += int64(skew.Seconds())
	return now >= nbf
}
