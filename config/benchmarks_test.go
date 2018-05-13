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
	"strings"
	"testing"

	"github.com/corestoreio/pkg/config"
	"github.com/corestoreio/pkg/config/storage"
	"github.com/corestoreio/pkg/store/scope"
)

var benchmarkNewByParts *config.Path

// BenchmarkNewByParts-4	 5000000	       297 ns/op	      48 B/op	       1 allocs/op
func BenchmarkNewByParts(b *testing.B) {
	want := config.MustNewPath("general/single_store_mode/enabled")
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchmarkNewByParts = config.MustNewPath("general/single_store_mode/enabled")
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
		benchmarkPathFQ, err = config.MustNewPath(p).BindWebsite(scopeID).FQ()
		if err != nil {
			b.Error(err)
		}
	}
	if benchmarkPathFQ != want {
		b.Errorf("Want: %s; Have, %s", want, benchmarkPathFQ)
	}
}

var benchmarkPathHash uint64

// BenchmarkPathHashFull-4  	 3000000	       502 ns/op	       0 B/op	       0 allocs/op
func BenchmarkPathHashFull(b *testing.B) {
	const scopeID = 12
	const want uint64 = 18184461197473735898
	p := "system/dev/debug"
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchmarkPathHash = config.MustNewPath(p).BindWebsite(scopeID).Hash64ByLevel(-1)
	}
	if benchmarkPathHash != want {
		b.Errorf("Want: %d; Have, %d", want, benchmarkPathHash)
	}
}

func BenchmarkPathHashLevel2(b *testing.B) {
	const scopeID = 13
	const want uint64 = 13528445590332414707
	p := "system/dev/debug"
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchmarkPathHash = config.MustNewPath(p).BindWebsite(scopeID).Hash64ByLevel(2)
	}
	if benchmarkPathHash != want {
		b.Errorf("Want: %d; Have, %d", want, benchmarkPathHash)
	}
}

var benchmarkReverseFQPath = new(config.Path)

// BenchmarkSplitFQ-4  	10000000	       199 ns/op	      32 B/op	       1 allocs/op
func BenchmarkSplitFQ(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if err := benchmarkReverseFQPath.ParseFQ("stores/7475/catalog/frontend/list_allow_all"); err != nil {
			b.Error(err)
		}
	}
	ls, _ := benchmarkReverseFQPath.Level(4)
	if ls != "stores/7475/catalog/frontend" {
		b.Errorf("stores/7475/catalog/frontend not found in Level(): %q", ls)
	}
}

var benchmarkRouteLevel string

// BenchmarkPath_Level_One-4	 5000000	       297 ns/op	      16 B/op	       1 allocs/op
func BenchmarkPath_Level_One(b *testing.B) {
	benchmarkRouteLevelRun(b, 3, "system/dev/debug", "default/0/system")
}

// BenchmarkPath_Level_Two-4	 5000000	       332 ns/op	      16 B/op	       1 allocs/op
func BenchmarkPath_Level_Two(b *testing.B) {
	benchmarkRouteLevelRun(b, 4, "system/dev/debug", "default/0/system/dev")
}

// BenchmarkPath_Level_All-4	 5000000	       379 ns/op	      16 B/op	       1 allocs/op
func BenchmarkPath_Level_All(b *testing.B) {
	benchmarkRouteLevelRun(b, -1, "system/dev/debug", "default/0/system/dev/debug")
}

func benchmarkRouteLevelRun(b *testing.B, level int, have, want string) {
	hp := config.MustNewPath(have)

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

var benchmarkPath_Part string

// BenchmarkPath_Part-4	 5000000	       240 ns/op	       0 B/op	       0 allocs/op
func BenchmarkPath_Part(b *testing.B) {
	have := config.MustNewPath("general/single_store_mode/enabled")
	want := "enabled"

	var err error
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchmarkPath_Part, err = have.Part(3)
		if err != nil {
			b.Error(err)
		}
		if benchmarkPath_Part == "" {
			b.Error("benchmarkPath_Part is nil! Unexpected")
		}
	}
	if want != benchmarkPath_Part {
		b.Errorf("Want: %q; Have: %q", want, benchmarkPath_Part)
	}
}

var benchmarkPath_Validate error

func BenchmarkPath_Validate(b *testing.B) {
	have := config.MustNewPath("system/dEv/d3bug")
	want := "system/dev/debug"

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchmarkPath_Validate = have.IsValid()
		if nil != benchmarkPath_Validate {
			b.Errorf("Want: %s; Have: %v", want, have)
		}
	}
}

var benchmarkRouteSplit = make([]string, 0, 4)

// BenchmarkPath_Split-4    	 5000000	       286 ns/op	       0 B/op	       0 allocs/op
func BenchmarkPath_Split(b *testing.B) {
	have := config.MustNewPath("general/single_store_mode/enabled")
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var err error
		benchmarkRouteSplit, err = have.Split(benchmarkRouteSplit...)
		if err != nil {
			b.Error(err)
		}
		if benchmarkRouteSplit[1] == "" {
			b.Error("benchmarkRouteSplit[1] is nil! Unexpected")
		}
		benchmarkRouteSplit = benchmarkRouteSplit[:0]
	}
}

var benchmarkScopedServiceVal *config.Value

// BenchmarkScopedServiceStringStore-4	 1000000	      2218 ns/op	     320 B/op	       9 allocs/op => Go 1.5.2
// BenchmarkScopedServiceStringStore-4	  500000	      2939 ns/op	     672 B/op	      17 allocs/op => Go 1.5.3 strings
// BenchmarkScopedServiceStringStore-4    500000	      2732 ns/op	     912 B/op	      17 allocs/op => cfgpath.Path with []ArgFunc
// BenchmarkScopedServiceStringStore-4	 1000000	      1821 ns/op	     336 B/op	       3 allocs/op => cfgpath.Path without []ArgFunc
// BenchmarkScopedServiceStringStore-4   1000000	      1747 ns/op	       0 B/op	       0 allocs/op => Go 1.6 sync.Pool cfgpath.Path without []ArgFunc
// BenchmarkScopedServiceStringStore-4    500000	      2604 ns/op	       0 B/op	       0 allocs/op => Go 1.7 with ScopeHash
func BenchmarkScopedServiceStringStore(b *testing.B) {
	benchmarkScopedServiceStringRun(b, 1, 1)
}

func BenchmarkScopedServiceStringWebsite(b *testing.B) {
	benchmarkScopedServiceStringRun(b, 1, 0)
}

func BenchmarkScopedServiceStringDefault(b *testing.B) {
	benchmarkScopedServiceStringRun(b, 0, 0)
}

func benchmarkScopedServiceStringRun(b *testing.B, websiteID, storeID int64) {
	route := "aa/bb/cc"
	want := strings.Repeat("Gopher", 100)

	sg := config.NewFakeService(storage.NewMap(
		config.MustNewPath(route).String(), want,
	)).NewScoped(websiteID, storeID)

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		benchmarkScopedServiceVal = sg.Get(scope.Store, route)
		if !benchmarkScopedServiceVal.IsValid() {
			b.Fatal(benchmarkScopedServiceVal)
		}
		s, ok, err := benchmarkScopedServiceVal.Str()
		if !ok {
			b.Fatal("path must be valid")
		}
		if err != nil {
			b.Fatal(err)
		}
		if s != want {
			b.Errorf("Want %s Have %s", want, benchmarkScopedServiceVal)
		}
	}
}

// BenchmarkPathSlice_Sort-4	 1000000	      1987 ns/op	     480 B/op	       8 allocs/op
func BenchmarkPathSlice_Sort(b *testing.B) {
	// allocs are here uninteresting

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ps := config.PathSlice{
			config.MustNewPathWithScope(scope.Store.WithID(3), "bb/cc/dd"),
			config.MustNewPathWithScope(scope.Store.WithID(2), "bb/cc/dd"),
			config.MustNewPath("bb/cc/dd"),
			config.MustNewPathWithScope(scope.Website.WithID(3), "xx/yy/zz"),
			config.MustNewPathWithScope(scope.Website.WithID(1), "xx/yy/zz"),
			config.MustNewPathWithScope(scope.Website.WithID(2), "xx/yy/zz"),
			config.MustNewPathWithScope(scope.Store.WithID(4), "zz/aa/bb"),
			config.MustNewPathWithScope(scope.Website.WithID(1), "zz/aa/bb"),
			config.MustNewPathWithScope(scope.Website.WithID(2), "aa/bb/cc"),
			config.MustNewPath("aa/bb/cc"),
		}
		ps.Sort()
		if len(ps) != 10 {
			b.Fatalf("Incorrect length %d of ps variable after sorting", len(ps))
		}
	}
}

var benchmarkPath config.Path

func BenchmarkPath_Marshal(b *testing.B) {
	var data []byte
	var err error
	const path = "system/full_page_cache/varnish/backend_port"
	const want = "stores/123/" + path

	// BenchmarkPath_Marshal/MarshalText-4         	 3000000	       592 ns/op	     112 B/op	       1 allocs/op
	b.Run("MarshalText", func(b *testing.B) {
		p := config.MustNewPath(path).BindStore(123)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			if data, err = p.MarshalText(); err != nil {
				b.Fatal(err)
			}
			if len(data) != 54 {
				b.Fatalf("Invalid data length: %d", len(data))
			}
		}
	})
	// BenchmarkPath_Marshal/UnmarshalText-4       	 3000000	       546 ns/op	      48 B/op	       1 allocs/op
	b.Run("UnmarshalText", func(b *testing.B) {
		var bData = []byte(want)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			if err := benchmarkPath.UnmarshalText(bData); err != nil {
				b.Fatal(err)
			}
			//b.Fatalf(benchmarkPath.String())
			if benchmarkPath.ScopeID != scope.TypeID(67108987) {
				b.Fatalf("Invalid scope: %d", benchmarkPath.ScopeID)
			}
		}
	})
	// BenchmarkPath_Marshal/MarshalBinary-4       	20000000	        95.5 ns/op	     112 B/op	       1 allocs/op
	b.Run("MarshalBinary", func(b *testing.B) {
		p := config.MustNewPath(path).BindStore(123)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			if data, err = p.MarshalBinary(); err != nil {
				b.Fatal(err)
			}
			if len(data) != 51 {
				b.Fatalf("Invalid data length: %d", len(data))
			}
		}
	})
	// BenchmarkPath_Marshal/UnmarshalBinary-4     	 3000000	       500 ns/op	      48 B/op	       1 allocs/op
	b.Run("UnmarshalBinary", func(b *testing.B) {
		var bData = []byte("{\x00\x00\x04\x00\x00\x00\x00system/full_page_cache/varnish/backend_port")
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			if err := benchmarkPath.UnmarshalBinary(bData); err != nil {
				b.Fatal(err)
			}
			if benchmarkPath.ScopeID != scope.TypeID(67108987) {
				b.Fatalf("Invalid scope: %d", benchmarkPath.ScopeID)
			}
		}
	})
}
