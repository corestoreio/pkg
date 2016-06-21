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

package csjwt_test

import (
	"encoding/gob"
	"regexp"
	"testing"
	"time"

	"github.com/corestoreio/csfw/util/csjwt"
	"github.com/corestoreio/csfw/util/csjwt/jwtclaim"
	"github.com/stretchr/testify/assert"
)

var _ csjwt.Deserializer = (*csjwt.JSONEncoding)(nil)
var _ csjwt.Serializer = (*csjwt.JSONEncoding)(nil)

func isBase64Token(str []byte) bool {
	r, _ := regexp.Compile(`^[A-Za-z0-9\-_\.]+$`)
	return r.Match(str)
}

func init() {
	gob.Register(csjwt.NewHead())
	gob.Register(jwtclaim.NewStore())
}

func TestGobEncoding(t *testing.T) {

	storeClaim := jwtclaim.NewStore()

	storeClaim.Store = "ch-en"
	storeClaim.ID = "2342-234345-234234-23435"
	storeClaim.ExpiresAt = time.Now().Add(time.Minute * 2).Unix()
	storeClaim.IssuedAt = time.Now().Unix()

	tk := csjwt.NewToken(storeClaim)
	tk.Serializer = csjwt.GobEncoding{}

	m := csjwt.NewSigningMethodHS512()
	pw := csjwt.WithPasswordRandom()
	tkChar, err := tk.SignedString(m, pw)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("gob", tkChar)

	// check if it is base64 encoded
	if !isBase64Token(tkChar) {
		t.Fatalf("Token is not base64 encoded! %q", tkChar)
	}

	vrf := csjwt.NewVerification(m)
	vrf.Deserializer = csjwt.GobEncoding{}

	newTk := csjwt.NewToken(jwtclaim.NewStore())

	if err := vrf.Parse(&newTk, tkChar, csjwt.NewKeyFunc(m, pw)); err != nil {
		t.Fatalf("%+v", err)
	}

	haveStoreClaim := newTk.Claims.(*jwtclaim.Store)
	assert.Exactly(t, "ch-en", haveStoreClaim.Store)
	assert.Exactly(t, "2342-234345-234234-23435", haveStoreClaim.ID)
}
