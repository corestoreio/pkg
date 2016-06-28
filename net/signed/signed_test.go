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
	"crypto/sha256"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/corestoreio/csfw/net"
	"github.com/corestoreio/csfw/net/mw"
	"github.com/corestoreio/csfw/net/signed"
	"github.com/corestoreio/csfw/util/cstesting"
)

var data = []byte(`“The most important property of a program is whether it accomplishes the intention of its user.” ― C.A.R. Hoare`)

const dataSHA256 = `keyId="test",algorithm="rot13",signature="cc7b14f207d3896a74ba4e4e965d49e6098af2191058edb9e9247caf0db8cd7b"`

func TestWithSignature(t *testing.T) {

	hpu := cstesting.NewHTTPParallelUsers(10, 5, 200, time.Millisecond)
	hpu.AssertResponse = func(rec *httptest.ResponseRecorder) {
		if have, want := rec.Header().Get(net.ContentSignature), dataSHA256; have != want {
			t.Errorf("Signature mismatch Have: %v Want: %v", have, want)
		}
	}

	r := httptest.NewRequest("GET", "http://corestore.io", nil)

	hpu.ServeHTTP(r, mw.Chain(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(net.ContentEncoding, "sHa256")
		if _, err := w.Write(data); err != nil {
			t.Fatal(err)
		}
	}), signed.WithSignature(sha256.New)))
}
