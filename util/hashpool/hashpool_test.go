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
	"crypto/sha256"
	"hash"
	"hash/crc64"
	"hash/fnv"
	"sync"
	"testing"

	"github.com/corestoreio/csfw/util/hashpool"
	"github.com/pierrec/xxHash/xxHash64"
	"github.com/stretchr/testify/assert"
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

func TestTank_SumHex(t *testing.T) {
	hp := hashpool.New(sha256.New)
	if have, want := hp.SumHex(data), dataSHA256; have != want {
		t.Errorf("Have: %v Want: %v", have, want)
	}
}

func BenchmarkTank_SumHex_SHA256(b *testing.B) {
	hp := hashpool.New(sha256.New)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if have, want := hp.SumHex(data), dataSHA256; have != want {
			b.Errorf("Have: %v Want: %v", have, want)
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

}
