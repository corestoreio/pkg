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
	"net/http"
	"net/url"
	"testing"

	"github.com/corestoreio/csfw/storage/text"
	"github.com/corestoreio/csfw/util/csjwt"
	"github.com/corestoreio/csfw/util/csjwt/jwtclaim"
	"github.com/stretchr/testify/assert"
)

func genParseTk(t *testing.T) (text.Chars, csjwt.Keyfunc) {
	hs256 := csjwt.NewSigningMethodHS256()
	key := csjwt.WithPasswordRandom()
	rawTK, err := csjwt.NewToken(&jwtclaim.Map{"extractMe": 3.14159}).SignedString(hs256, key)
	if err != nil {
		t.Fatalf("%+v", err)
	}
	return rawTK, csjwt.NewKeyFunc(hs256, key)
}

func compareParseTk(t *testing.T, haveTK csjwt.Token, err error) {
	if err != nil {
		t.Fatalf("%+v", err)
	}

	me, err := haveTK.Claims.Get("extractMe")
	if err != nil {
		t.Fatalf("%+v", err)
	}
	assert.Exactly(t, 3.14159, me)
}

func TestParse(t *testing.T) {

	rawTK, kf := genParseTk(t)
	haveTK := csjwt.NewToken(&jwtclaim.Map{})
	err := csjwt.Parse(&haveTK, rawTK, kf)
	compareParseTk(t, haveTK, err)
}

func TestParseFromRequest(t *testing.T) {

	rawTK, kf := genParseTk(t)

	r, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatalf("%+v", err)
	}
	r.Form = url.Values{
		csjwt.HTTPFormInputName: []string{rawTK.String()},
	}

	haveTK := csjwt.NewToken(&jwtclaim.Map{})
	err = csjwt.ParseFromRequest(&haveTK, kf, r)
	compareParseTk(t, haveTK, err)
}
