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
	"github.com/corestoreio/csfw/store/scope"
)

var benchmarkStr string

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
			var h scope.Hash
			benchmarkStr, h, err = p1.Get(sg)
			if err != nil {
				b.Error(err)
			}
			if benchmarkStr != want {
				b.Errorf("Have: %s\nWant: %s\n", benchmarkStr, want)
			}
			if h != scope.DefaultHash {
				b.Errorf("Have: %s\nWant: %s\n", h, scope.DefaultHash)
			}
		}
	})
}

func Benchmark_SingleStrGetDefault(b *testing.B) {
	const want = `Content-Type,X-CoreStore-ID`
	p1 := cfgmodel.NewStr("web/cors/exposed_headers", cfgmodel.WithFieldFromSectionSlice(configStructure))

	sg := cfgmock.NewService().NewScoped(1, 1)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var err error
		var h scope.Hash
		benchmarkStr, h, err = p1.Get(sg)
		if err != nil {
			b.Error(err)
		}
		if benchmarkStr != want {
			b.Errorf("Have: %s\nWant: %s\n", benchmarkStr, want)
		}
		if h != scope.DefaultHash {
			b.Errorf("Have: %s\nWant: %s\n", h, scope.DefaultHash)
		}
	}
}

func Benchmark_SingleStrGetWebsite(b *testing.B) {
	const want = `Content-Application`
	var wantHash = scope.NewHash(scope.Website, 2)
	p1 := cfgmodel.NewStr("web/cors/exposed_headers", cfgmodel.WithFieldFromSectionSlice(configStructure))

	sg := cfgmock.NewService(cfgmock.WithPV(cfgmock.PathValue{
		p1.MustFQ(scope.Website, 2): want,
	})).NewScoped(2, 4)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var err error
		var h scope.Hash
		benchmarkStr, h, err = p1.Get(sg)
		if err != nil {
			b.Error(err)
		}
		if benchmarkStr != want {
			b.Errorf("Have: %s\nWant: %s\n", benchmarkStr, want)
		}
		if h != wantHash {
			b.Errorf("Have: %s\nWant: %s\n", h, scope.DefaultHash)
		}
	}
}

var benchmarkByte []byte

func Benchmark_SingleByteGetDefault(b *testing.B) {
	var want = []byte(`Hello Dud€`)

	p1 := cfgmodel.NewByte("web/cors/byte", cfgmodel.WithFieldFromSectionSlice(configStructure))

	sg := cfgmock.NewService().NewScoped(1, 1)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var err error
		var h scope.Hash
		benchmarkByte, h, err = p1.Get(sg)
		if err != nil {
			b.Error(err)
		}
		if bytes.Compare(benchmarkByte, want) != 0 {
			b.Errorf("Have: %s\nWant: %s\n", string(benchmarkByte), string(want))
		}
		if h != scope.DefaultHash {
			b.Errorf("Have: %s\nWant: %s\n", h, scope.DefaultHash)
		}
	}
}

var benchmark_SingleFloat64GetStore float64

func Benchmark_SingleFloat64GetStore(b *testing.B) {
	const want float64 = 3.14159
	var wantHash = scope.NewHash(scope.Store, 4)
	p1 := cfgmodel.NewFloat64("web/cors/float64_store", cfgmodel.WithFieldFromSectionSlice(configStructure))
	if p1.LastError != nil {
		b.Fatalf("%+v", p1.LastError)
	}

	sg := cfgmock.NewService(cfgmock.WithPV(cfgmock.PathValue{
		p1.MustFQ(scope.Store, 4): want,
	})).NewScoped(2, 4)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var err error
		var h scope.Hash
		benchmark_SingleFloat64GetStore, h, err = p1.Get(sg)
		if err != nil {
			b.Error(err)
		}
		if benchmark_SingleFloat64GetStore != want {
			b.Errorf("Have: %s\nWant: %s\n", benchmarkStr, want)
		}
		if h != wantHash {
			b.Errorf("Have: %s\nWant: %s\n", h, scope.DefaultHash)
		}
	}
}
