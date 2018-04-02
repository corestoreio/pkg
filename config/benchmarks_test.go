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

package config_test

import (
	"strconv"
	"testing"

	"github.com/corestoreio/pkg/config"
	"github.com/corestoreio/pkg/store/scope"
)

var benchmarkNewByParts config.Path

// BenchmarkNewByParts-4	 5000000	       297 ns/op	      48 B/op	       1 allocs/op
func BenchmarkNewByParts(b *testing.B) {
	want := config.MustMakePath("general/single_store_mode/enabled")
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchmarkNewByParts = config.MustMakePath("general/single_store_mode/enabled")
	}
	if !benchmarkNewByParts.Equal(want) {
		b.Errorf("Want: %s; Have, %s", want, benchmarkNewByParts)
	}
}

var benchmarkPathFQ string

// BenchmarkPathFQ-4     	 3000000	       401 ns/op	     112 B/op	       1 allocs/op
func BenchmarkPathFQ(b *testing.B) {
	var scopeID int64 = 11
	want := scope.StrWebsites.String() + "/" + strconv.FormatInt(scopeID, 10) + "/system/dev/debug"
	p := "system/dev/debug"
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var err error
		benchmarkPathFQ, err = config.MustMakePath(p).BindWebsite(scopeID).FQ()
		if err != nil {
			b.Error(err)
		}
	}
	if benchmarkPathFQ != want {
		b.Errorf("Want: %s; Have, %s", want, benchmarkPathFQ)
	}
}

var benchmarkPathHash uint32

// BenchmarkPathHashFull-4  	 3000000	       502 ns/op	       0 B/op	       0 allocs/op
func BenchmarkPathHashFull(b *testing.B) {
	const scopeID int64 = 12
	const want uint32 = 1479679325
	p := "system/dev/debug"
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var err error
		benchmarkPathHash, err = config.MustMakePath(p).BindWebsite(scopeID).Hash(-1)
		if err != nil {
			b.Error(err)
		}
	}
	if benchmarkPathHash != want {
		b.Errorf("Want: %d; Have, %d", want, benchmarkPathHash)
	}
}

func BenchmarkPathHashLevel2(b *testing.B) {
	const scopeID int64 = 13
	const want uint32 = 723768876
	p := "system/dev/debug"
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var err error
		benchmarkPathHash, err = config.MustMakePath(p).BindWebsite(scopeID).Hash(2)
		if err != nil {
			b.Error(err)
		}
	}
	if benchmarkPathHash != want {
		b.Errorf("Want: %d; Have, %d", want, benchmarkPathHash)
	}
}

var benchmarkReverseFQPath config.Path

// BenchmarkSplitFQ-4  	10000000	       199 ns/op	      32 B/op	       1 allocs/op
func BenchmarkSplitFQ(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var err error
		benchmarkReverseFQPath, err = config.SplitFQ("stores/7475/catalog/frontend/list_allow_all")
		if err != nil {
			b.Error(err)
		}
	}
	ls, _ := benchmarkReverseFQPath.Level(-1)
	if ls != "catalog/frontend/list_allow_all" {
		b.Error("catalog/frontend/list_allow_all not found in Level()")
	}
}

var benchmarkRouteLevel string

// BenchmarkRouteLevel_One-4	 5000000	       297 ns/op	      16 B/op	       1 allocs/op
func BenchmarkRouteLevel_One(b *testing.B) {
	benchmarkRouteLevelRun(b, 1, "system/dev/debug", "system")
}

// BenchmarkRouteLevel_Two-4	 5000000	       332 ns/op	      16 B/op	       1 allocs/op
func BenchmarkRouteLevel_Two(b *testing.B) {
	benchmarkRouteLevelRun(b, 2, "system/dev/debug", "system/dev")
}

// BenchmarkRouteLevel_All-4	 5000000	       379 ns/op	      16 B/op	       1 allocs/op
func BenchmarkRouteLevel_All(b *testing.B) {
	benchmarkRouteLevelRun(b, -1, "system/dev/debug", "system/dev/debug")
}

func benchmarkRouteLevelRun(b *testing.B, level int, have, want string) {
	hp := config.MustMakePath(have)

	b.ResetTimer()
	var err error
	for i := 0; i < b.N; i++ {
		benchmarkRouteLevel, err = hp.Level(level)
	}
	if err != nil {
		b.Error(err)
	}
	if benchmarkRouteLevel != want {
		b.Errorf("Want: %s; Have, %s", want, benchmarkRouteLevel)
	}
}

var benchmarkRouteHash uint32

// BenchmarkRouteHash-4     	 5000000	       287 ns/op	       0 B/op	       0 allocs/op
func BenchmarkRouteHash(b *testing.B) {
	have := config.MustMakePath("general/single_store_mode/enabled")
	want := uint32(1644245266)

	var err error
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchmarkRouteHash, err = have.Hash(3)
		if err != nil {
			b.Error(err)
		}
		if want != benchmarkRouteHash {
			b.Errorf("Want: %d; Have: %d", want, benchmarkRouteHash)
		}
	}
}

// BenchmarkRouteHash32-4   	50000000	        37.7 ns/op	       0 B/op	       0 allocs/op
func BenchmarkRouteHash32(b *testing.B) {
	have := config.MustMakePath("general/single_store_mode/enabled")
	want := uint32(1644245266)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchmarkRouteHash = have.Hash32()
		if want != benchmarkRouteHash {
			b.Errorf("Want: %d; Have: %d", want, benchmarkRouteHash)
		}
	}
}

var benchmarkRoutePart string

// BenchmarkRoutePart-4	 5000000	       240 ns/op	       0 B/op	       0 allocs/op
func BenchmarkRoutePart(b *testing.B) {
	have := config.MustMakePath("general/single_store_mode/enabled")
	want := "enabled"

	var err error
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchmarkRoutePart, err = have.Part(3)
		if err != nil {
			b.Error(err)
		}
		if benchmarkRoutePart == "" {
			b.Error("benchmarkRoutePart is nil! Unexpected")
		}
	}
	if want != benchmarkRoutePart {
		b.Errorf("Want: %q; Have: %q", want, benchmarkRoutePart)
	}
}

var benchmarkRouteValidate error

// BenchmarkRouteValidate-4	20000000	        83.5 ns/op	       0 B/op	       0 allocs/op
func BenchmarkRouteValidate(b *testing.B) {
	have := config.MustMakePath("system/dEv/d3bug")
	want := "system/dev/debug"

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchmarkRouteValidate = have.IsValid()
		if nil != benchmarkRouteValidate {
			b.Errorf("Want: %s; Have: %v", want, have)
		}
	}
}

var benchmarkRouteSplit []string

// BenchmarkRouteSplit-4    	 5000000	       286 ns/op	       0 B/op	       0 allocs/op
func BenchmarkRouteSplit(b *testing.B) {
	have := config.MustMakePath("general/single_store_mode/enabled")
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var err error
		benchmarkRouteSplit, err = have.Split()
		if err != nil {
			b.Error(err)
		}
		if benchmarkRouteSplit[1] == "" {
			b.Error("benchmarkRouteSplit[1] is nil! Unexpected")
		}
	}
}
