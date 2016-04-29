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

package cfgmodel_test

import (
	"bytes"
	"testing"

	"github.com/corestoreio/csfw/config/cfgmock"
	"github.com/corestoreio/csfw/config/cfgmodel"
)

var benchmarkStr string

// Benchmark_ParallelStrGetDefault-4	 1000000	      2368 ns/op	      97 B/op	       4 allocs/op
func Benchmark_ParallelStrGetDefault(b *testing.B) {
	const want = `Content-Type,X-CoreStore-ID`
	const pathWebCorsHeaders = "web/cors/exposed_headers"
	p1 := cfgmodel.NewStr(pathWebCorsHeaders, cfgmodel.WithFieldFromSectionSlice(configStructure))

	sg := cfgmock.NewService().NewScoped(1, 1)

	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			var err error
			benchmarkStr, err = p1.Get(sg)
			if err != nil {
				b.Error(err)
			}
			if benchmarkStr != want {
				b.Errorf("Have: %s\nWant: %s\n", benchmarkStr, want)
			}
		}
	})
}

// Benchmark_SingleStrGetDefault-4  	 1000000	      1871 ns/op	      97 B/op	       4 allocs/op
func Benchmark_SingleStrGetDefault(b *testing.B) {
	const want = `Content-Type,X-CoreStore-ID`
	p1 := cfgmodel.NewStr("web/cors/exposed_headers", cfgmodel.WithFieldFromSectionSlice(configStructure))

	sg := cfgmock.NewService().NewScoped(1, 1)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var err error
		benchmarkStr, err = p1.Get(sg)
		if err != nil {
			b.Error(err)
		}
		if benchmarkStr != want {
			b.Errorf("Have: %s\nWant: %s\n", benchmarkStr, want)
		}
	}
}

var benchmarkByte []byte

// Benchmark_SingleByteGetDefault-4 	 1000000	      1363 ns/op	       1 B/op	       1 allocs/op
func Benchmark_SingleByteGetDefault(b *testing.B) {
	var want = []byte(`Hello Dud€`)

	p1 := cfgmodel.NewByte("web/cors/byte", cfgmodel.WithFieldFromSectionSlice(configStructure))

	sg := cfgmock.NewService().NewScoped(1, 1)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var err error
		benchmarkByte, err = p1.Get(sg)
		if err != nil {
			b.Error(err)
		}
		if bytes.Compare(benchmarkByte, want) != 0 {
			b.Errorf("Have: %s\nWant: %s\n", string(benchmarkByte), string(want))
		}
	}
}
