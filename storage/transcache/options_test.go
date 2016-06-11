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

package transcache_test

import (
	"bytes"
	"encoding/gob"
	"io"
	"testing"
)

func xxxTestGobPrimer(t *testing.T) {

	bufEnc := new(bytes.Buffer)
	enc := gob.NewEncoder(bufEnc)
	bufDec := new(bytes.Buffer)
	dec := gob.NewDecoder(bufDec)

	//var tc := Country{}
	//var tss := TableStoreSlice{}
	//var ts := TableStore{}

	var primeType interface{}
	primeType = TableStore{IsActive: true}

	if err := enc.Encode(primeType); err != nil {
		t.Fatal(err)
	}
	t.Logf("Enc1: %#v\n", bufEnc.String())

	if _, err := io.Copy(bufDec, bufEnc); err != nil {
		t.Fatal(err)
	}

	var ts TableStore
	if err := dec.Decode(&ts); err != nil {
		t.Fatal(err)
	}
	t.Logf("Dec1: %#v\n", bufDec.String())
	t.Logf("Dec2: %#v\n", ts)

}
