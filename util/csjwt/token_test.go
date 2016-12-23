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
	"testing"

	"bytes"

	"github.com/corestoreio/csfw/util/csjwt"
	"github.com/corestoreio/csfw/util/csjwt/jwtclaim"
	"github.com/corestoreio/log"
	"github.com/corestoreio/log/logw"
	"github.com/stretchr/testify/assert"
)

func TestTokenAlg(t *testing.T) {
	tests := []struct {
		tok     csjwt.Token
		wantAlg string
	}{
		{csjwt.NewToken(nil), ""},
		{
			csjwt.Token{
				Header: jwtclaim.NewHeadSegments("3"),
			},
			"3",
		},
		{
			csjwt.Token{
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
	tk := csjwt.NewToken(jwtclaim.Map{"lang": "Golang", "extractMe": 3.14159})
	buf := bytes.Buffer{}
	lg := logw.NewLog(logw.WithWriter(&buf), logw.WithLevel(logw.LevelDebug))
	lg.Debug("tokenTest", log.Marshal("xtoken", tk))
	have := `tokenTest token: "eyJ0eXAiOiJKV1QifQo.eyJleHRyYWN0TWUiOjMuMTQxNTksImxhbmciOiJHb2xhbmcifQo"`
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
