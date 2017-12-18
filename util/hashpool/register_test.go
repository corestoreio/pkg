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

package hashpool

import (
	"crypto"
	"hash/fnv"
	"testing"

	"github.com/corestoreio/errors"
	"github.com/stretchr/testify/assert"
)

const dataStr = `Don't communicate by sharing memory, share memory by communicating. ∏`

var dataByte = []byte(`Design the architecture, name the components, document the details. €`)

func TestFnv64aWrite(t *testing.T) {

	f64 := fnv.New64a()
	_, _ = f64.Write(dataByte)
	_, _ = f64.Write([]byte(dataStr))
	wantSum := f64.Sum64()

	haveSum := hash64(0).writeBytes(dataByte)
	haveSum = haveSum.writeStr(dataStr)

	assert.Exactly(t, hash64(wantSum), haveSum)
}

var benchmarkFnv64aWriteStrSum hash64

// 20000000	        89.6 ns/op	       0 B/op	       0 allocs/op
func BenchmarkFnv64aWriteStr(b *testing.B) {
	f64 := fnv.New64a()
	_, _ = f64.Write([]byte(dataStr))
	wantSum := f64.Sum64()
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		benchmarkFnv64aWriteStrSum = hash64(0).writeStr(dataStr)
	}
	if benchmarkFnv64aWriteStrSum != hash64(wantSum) {
		b.Errorf("Have: %v Want: %v", benchmarkFnv64aWriteStrSum, wantSum)
	}
}

func TestMustFromRegister(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			if err, ok := r.(error); ok {
				assert.True(t, errors.NotFound.Match(err), "%+v", err)
			} else {
				t.Errorf("Expecting an error but got %#v", r)
			}
		} else {
			t.Error("Expecting a panic")
		}
	}()
	_ = MustFromRegistry("LoseLose")
}

func TestRegisterDeregister(t *testing.T) {
	assert.NoError(t, Register("sha256", crypto.SHA256.New))

	assert.True(t, errors.AlreadyExists.Match(Register("sha256", nil)))

	h, err := FromRegistry("sha256")
	assert.NoError(t, err)
	assert.Len(t, h.SumHex(dataByte), 64)

	Deregister("sha256")
	h, err = FromRegistry("sha256")
	assert.True(t, errors.NotFound.Match(err), "%+v", err)
	assert.Exactly(t, Tank{}, h)
}

func TestFromRegistryHMAC_Error(t *testing.T) {
	s512t, err := FromRegistryHMAC("sha512", []byte(`pw123456`))
	assert.True(t, errors.NotFound.Match(err), "%+v", err)
	assert.Exactly(t, Tank{}, s512t)
}

func TestFromRegistryHMAC(t *testing.T) {
	assert.NoError(t, Register("sha512", crypto.SHA512.New))

	s512t, err := FromRegistryHMAC("sha512", []byte(`pw123456`))
	assert.NoError(t, err)
	assert.Len(t, s512t.SumHex(dataByte), 128)

	s512t2, err2 := FromRegistryHMAC("sha512", []byte(`pwABCDEFGH`))
	assert.NoError(t, err2)
	assert.Len(t, s512t2.SumHex(dataByte), 128)
	assert.NotEqual(t, s512t.SumHex(dataByte), s512t2.SumHex(dataByte))

	assert.Len(t, db.ht, 3)

	Deregister("sha512")
	s512t, err = FromRegistryHMAC("sha512", []byte(`pw123456`))
	assert.True(t, errors.NotFound.Match(err), "%+v", err)
	assert.Exactly(t, Tank{}, s512t)

	s512t2, err2 = FromRegistryHMAC("sha512", []byte(`pwABCDEFGH`))
	assert.True(t, errors.NotFound.Match(err2), "%+v", err2)
	assert.Exactly(t, Tank{}, s512t2)
	assert.Len(t, db.ht, 0)
}
