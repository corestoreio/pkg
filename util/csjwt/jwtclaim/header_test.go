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

package jwtclaim_test

import (
	"encoding/json"
	"testing"

	"github.com/corestoreio/csfw/util/csjwt"
	"github.com/corestoreio/csfw/util/csjwt/jwtclaim"
	"github.com/stretchr/testify/assert"
)

var _ csjwt.Header = (*jwtclaim.HeadSegments)(nil)

func TestHeadSegmentsParseJSON(t *testing.T) {
	var sc csjwt.Header
	sc = &jwtclaim.HeadSegments{
		Algorithm: `ES999`,
		Type:      jwtclaim.ContentTypeJWT,
	}
	rawJSON, err := json.Marshal(sc)
	if err != nil {
		t.Fatal(err)
	}
	assert.Len(t, rawJSON, 27, "JSON: %s", rawJSON)

	scNew := &jwtclaim.HeadSegments{}
	if err := json.Unmarshal(rawJSON, &scNew); err != nil {
		t.Fatal(err)
	}
	assert.Exactly(t, sc, scNew)
}

func TestHeadSegmentsAlgTyp(t *testing.T) {

	var sc csjwt.Header
	sc = jwtclaim.NewHeadSegments(`ES999`)
	assert.Exactly(t, "ES999", sc.Alg())
	assert.Exactly(t, jwtclaim.ContentTypeJWT, sc.Typ())
}

func TestHeadSegmentsGetSet(t *testing.T) {
	tests := []struct {
		sc         csjwt.Header
		key        string
		val        string
		wantSetErr error
		wantGetErr error
	}{
		{&jwtclaim.HeadSegments{}, jwtclaim.HeaderAlg, "", nil, nil},
		{&jwtclaim.HeadSegments{}, jwtclaim.HeaderTyp, "Go", nil, nil},
	}
	for i, test := range tests {

		haveSetErr := test.sc.Set(test.key, test.val)
		if test.wantSetErr != nil {
			assert.EqualError(t, haveSetErr, test.wantSetErr.Error(), "Index %d", i)
		} else {
			assert.NoError(t, haveSetErr, "Index %d", i)
		}

		haveVal, haveGetErr := test.sc.Get(test.key)
		if test.wantGetErr != nil {
			assert.EqualError(t, haveGetErr, test.wantGetErr.Error(), "Index %d", i)
			continue
		} else {
			assert.NoError(t, haveGetErr, "Index %d", i)
		}

		if test.wantSetErr == nil {
			assert.Exactly(t, test.val, haveVal, "Index %d", i)
		}
	}
}
