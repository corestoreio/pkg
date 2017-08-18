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

package hashpool_test

import (
	"bytes"
	"crypto"
	"encoding/hex"
	"hash"
	"hash/crc64"
	"hash/fnv"
	"sync"
	"testing"

	"github.com/corestoreio/csfw/util/hashpool"
	"github.com/corestoreio/errors"
	"github.com/dchest/siphash"
	"github.com/pierrec/xxHash/xxHash64"
	"github.com/stretchr/testify/assert"
	_ "golang.org/x/crypto/blake2b"
)

var data = []byte(`“The most important property of a program is whether it accomplishes the intention of its user.” ― C.A.R. Hoare`)

const dataSHA256 = "cc7b14f207d3896a74ba4e4e965d49e6098af2191058edb9e9247caf0db8cd7b"

func TestNew64(t *testing.T) {
	p := hashpool.New64(fnv.New64a)
	const iterations = 10
	var wg sync.WaitGroup
	wg.Add(iterations)
	for i := 0; i < iterations; i++ {
		go func(wg *sync.WaitGroup) {
			defer wg.Done()

			fh := fnv.New64a()
			fh.Write(data)
			want := fh.Sum(nil)
			have := p.Sum(data, nil)
			assert.Exactly(t, want, have)
		}(&wg)
	}
	wg.Wait()
}

func TestTank_Equal(t *testing.T) {
	hp := hashpool.New(crypto.SHA256.New)
	mac, err := hex.DecodeString(dataSHA256)
	assert.NoError(t, err)
	assert.True(t, hp.Equal(data, mac))
}

func TestTank_EqualReader(t *testing.T) {
	hp := hashpool.New(crypto.SHA256.New)
	mac, err := hex.DecodeString(dataSHA256)
	assert.NoError(t, err)
	isEqual, err := hp.EqualReader(bytes.NewReader(data), mac)
	assert.NoError(t, err)
	assert.True(t, isEqual)
}

func TestTank_EqualReader_Error(t *testing.T) {
	hp := hashpool.New(crypto.SHA256.New)
	mac, err := hex.DecodeString(dataSHA256)
	assert.NoError(t, err)
	isEqual, err := hp.EqualReader(readerError{}, mac)
	assert.True(t, errors.IsAlreadyClosed(err), "%+v", err)
	assert.False(t, isEqual)
}

func TestTank_EqualPairs(t *testing.T) {
	u1 := []byte(`username1`)
	u2 := []byte(`username1`)
	p1 := []byte(`password1`)
	p2 := []byte(`password1`)

	hp := hashpool.New(crypto.SHA256.New)
	tests := []struct {
		pairs [][]byte
		want  bool
	}{
		{[][]byte{data, data, data, data}, true},
		{[][]byte{data, data, data}, false},
		{[][]byte{data, data}, true},
		{[][]byte{data}, false},
		{[][]byte{nil}, false},
		{[][]byte{nil, nil}, false},
		{[][]byte{}, false},
		{nil, false},
		{[][]byte{u1, u2}, true},
		{[][]byte{p1, p2}, true},
		{[][]byte{u1, u2, p1, p2}, true},
		{[][]byte{u2, u1, p1, p2}, true},
		{[][]byte{u2, p1, p1, p2}, false},
	}
	for i, test := range tests {
		assert.Exactly(t, test.want, hp.EqualPairs(test.pairs...), "Index %d", i)
	}
}

func BenchmarkTank_EqualPairs_SHA256_4args(b *testing.B) {
	u1 := []byte(`username1`)
	u2 := []byte(`username1`)
	p1 := []byte(`password1`)
	p2 := []byte(`password1`)
	hp := hashpool.New(crypto.SHA256.New)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if !hp.EqualPairs(u1, u2, p1, p2) {
			b.Fatal("Expecting to be true")
		}
	}
}

type readerError struct{}

func (readerError) Read(p []byte) (int, error) {
	return 0, errors.NewAlreadyClosedf("Reader already closed")
}

func TestTank_SumHex(t *testing.T) {
	hp := hashpool.New(crypto.SHA256.New)
	if have, want := hp.SumHex(data), dataSHA256; have != want {
		t.Errorf("Have: %v Want: %v", have, want)
	}
}

func BenchmarkTank_SumHex_SHA256(b *testing.B) {
	hp := hashpool.New(crypto.SHA256.New)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if have, want := hp.SumHex(data), dataSHA256; have != want {
			b.Errorf("Have: %v Want: %v", have, want)
		}
	}
}

func BenchmarkTank_SumHex_Blake2b256(b *testing.B) {
	const dataBlake2b256 = "00fad91702b9d9cfce8f6d3a7e2134283aa370b453e033ed6442dfef2a5c8089"

	hp := hashpool.New(crypto.BLAKE2b_256.New)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if have, want := hp.SumHex(data), dataBlake2b256; have != want {
			b.Fatalf("Have: %v Want: %v", have, want)
		}
	}
}

func benchmarkTank_Hash64(wantHash uint64, h func() hash.Hash64) func(b *testing.B) {
	return func(b *testing.B) {
		hp := hashpool.New64(h)
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			h := hp.Get()
			_, _ = h.Write(data)
			if have, want := h.Sum64(), wantHash; have != want {
				b.Errorf("Have: %v Want: %v", have, want)
			}
			hp.Put(h)
		}
	}
}

func BenchmarkTank_Hash64(b *testing.B) {
	b.Run("FNV64a", benchmarkTank_Hash64(207718596844850661, fnv.New64a))
	b.Run("xxHash", benchmarkTank_Hash64(11301805909362518010, func() hash.Hash64 { return xxHash64.New(uint64(201608090723)) }))
	b.Run("crc64", benchmarkTank_Hash64(3866349411325606150, func() hash.Hash64 { return crc64.New(crc64.MakeTable(crc64.ISO)) }))
	b.Run("siphash", benchmarkTank_Hash64(18240385100576365171, func() hash.Hash64 { return siphash.New(data) }))
}
