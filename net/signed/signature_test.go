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

package signed_test

import (
	"encoding/hex"
	"net/http/httptest"
	"testing"

	"github.com/corestoreio/csfw/net"
	"github.com/corestoreio/csfw/net/signed"
)

func TestSetHTTPContentSignature(t *testing.T) {

	w := httptest.NewRecorder()
	signed.WriteHTTPContentSignature(w, hex.EncodeToString, "myKeyID", "hmac-sha1", []byte(`Hello Gophers`))

	const wantSig = `keyId="myKeyID",algorithm="hmac-sha1",signature="48656c6c6f20476f7068657273"`
	if have, want := w.Header().Get(net.ContentSignature), wantSig; have != want {
		t.Errorf("Have: %v Want: %v", have, want)
	}
}
