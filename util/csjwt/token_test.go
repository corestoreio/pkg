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

package csjwt_test

import (
	"bytes"
	"encoding"
	"encoding/json"
	"testing"

	"github.com/corestoreio/log"
	"github.com/corestoreio/log/logw"
	"github.com/corestoreio/pkg/util/assert"
	"github.com/corestoreio/pkg/util/csjwt"
	"github.com/corestoreio/pkg/util/csjwt/jwtclaim"
)

func TestTokenAlg(t *testing.T) {
	tests := []struct {
		tok     *csjwt.Token
		wantAlg string
	}{
		{csjwt.NewToken(nil), ""},
		{
			&csjwt.Token{
				Header: jwtclaim.NewHeadSegments("3"),
			},
			"3",
		},
		{
			&csjwt.Token{
				Header: jwtclaim.NewHeadSegments("Gopher"),
			},
			"Gopher",
		},
	}
	for i, test := range tests {
		assert.Exactly(t, test.wantAlg, test.tok.Alg(), "Index %d", i)
	}
}

func TestToken_MarshalLog_Ok(t *testing.T) {
	tk := csjwt.NewToken(&jwtclaim.Store{Store: "Golang", UserID: "extractMe"})
	buf := bytes.Buffer{}
	lg := logw.NewLog(logw.WithWriter(&buf), logw.WithLevel(logw.LevelDebug))
	lg.Debug("tokenTest", log.Marshal("xtoken", tk))
	have := `tokenTest token: "eyJ0eXAiOiJKV1QifQ.eyJzdG9yZSI6IkdvbGFuZyIsInVzZXJpZCI6ImV4dHJhY3RNZSJ9Cg"`
	assert.Contains(t, buf.String(), have)
}

func TestToken_MarshalLog_Error(t *testing.T) {
	tk := csjwt.NewToken(jwtclaim.Map{"lang": "Golang", "extractMe": make(chan struct{})})
	buf := bytes.Buffer{}
	lg := logw.NewLog(logw.WithWriter(&buf), logw.WithLevel(logw.LevelDebug))
	lg.Debug("tokenTest", log.Marshal("xtoken", tk))
	have := `tokenTest token_error: "json: unsupported type: chan struct {}`
	assert.Contains(t, buf.String(), have)
}

var (
	_ encoding.BinaryMarshaler   = (*claimMapBinary)(nil)
	_ encoding.BinaryUnmarshaler = (*claimMapBinary)(nil)

	_ encoding.TextMarshaler   = (*claimMapText)(nil)
	_ encoding.TextUnmarshaler = (*claimMapText)(nil)

	_ json.Marshaler   = (*claimMapJSON)(nil)
	_ json.Unmarshaler = (*claimMapJSON)(nil)
)

type claimMapBinary struct {
	jwtclaim.Map
}

func (m *claimMapBinary) MarshalBinary() (data []byte, err error) {
	return []byte(m.Map["key1"].(string)), nil
}

func (m *claimMapBinary) UnmarshalBinary(data []byte) error {
	m.Map["key1"] = string(data)
	return nil
}

type claimMapText struct {
	jwtclaim.Map
}

func (m *claimMapText) MarshalText() (data []byte, err error) {
	return []byte(m.Map["key1"].(string)), nil
}

func (m *claimMapText) UnmarshalText(data []byte) error {
	m.Map["key1"] = string(data)
	return nil
}

type claimMapJSON struct {
	jwtclaim.Map
}

func (m *claimMapJSON) MarshalJSON() (data []byte, err error) {
	return []byte(m.Map["key1"].(string)), nil
}

func (m *claimMapJSON) UnmarshalJSON(data []byte) error {
	m.Map["key1"] = string(data)
	return nil
}

type claimMapProto struct {
	jwtclaim.Map
}

func (m *claimMapProto) Marshal() (data []byte, err error) {
	return []byte(m.Map["key1"].(string)), nil
}

func (m *claimMapProto) Unmarshal(data []byte) error {
	m.Map["key1"] = string(data)
	return nil
}

func TestToken_AllMarshalers(t *testing.T) {
	pwKey := csjwt.WithPasswordRandom()
	m := csjwt.NewSigningMethodHS256()
	v := csjwt.NewVerification(m)
	hs256KeyFn := csjwt.NewKeyFunc(m, pwKey)

	tests := []struct {
		name       string
		claim      csjwt.Claimer
		emptyClaim csjwt.Claimer
	}{
		{
			"Binary",
			&claimMapBinary{Map: jwtclaim.Map{}},
			&claimMapBinary{Map: jwtclaim.Map{}},
		},
		{
			"Text",
			&claimMapText{Map: jwtclaim.Map{}},
			&claimMapText{Map: jwtclaim.Map{}},
		},
		{
			"JSON",
			&claimMapJSON{Map: jwtclaim.Map{}},
			&claimMapJSON{Map: jwtclaim.Map{}},
		},
		{
			"Proto",
			&claimMapProto{Map: jwtclaim.Map{}},
			&claimMapProto{Map: jwtclaim.Map{}},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.NoError(t, test.claim.Set("key1", test.name))

			tk := csjwt.NewToken(test.claim)
			tkChar, err := tk.SignedString(m, pwKey)
			assert.NoError(t, err)
			t.Log(string(tkChar))

			tokenParts := bytes.Split(tkChar, []byte(`.`))
			decClaim, err := csjwt.DecodeSegment(tokenParts[1])
			assert.NoError(t, err)
			assert.Exactly(t, test.name, string(decClaim), "Token claim must match")

			tk2 := csjwt.NewToken(test.emptyClaim)
			err = v.Parse(tk2, tkChar, hs256KeyFn)
			assert.NoError(t, err)
			k1, err := test.emptyClaim.Get("key1")
			assert.NoError(t, err)
			assert.Exactly(t, test.name, k1.(string))
		})
	}
}
