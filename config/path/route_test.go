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

package path_test

import (
	"encoding/json"
	"testing"

	"bytes"
	"hash/fnv"

	"errors"

	"github.com/corestoreio/csfw/config/path"
	"github.com/stretchr/testify/assert"
)

func TestRouteEqual(t *testing.T) {
	tests := []struct {
		a    path.Route
		b    path.Route
		want bool
	}{
		{nil, nil, true},
		{path.Route("a"), path.Route("a"), true},
		{path.Route("a"), path.Route("b"), false},
		{path.Route("a\x80"), path.Route("a"), false},
	}
	for i, test := range tests {
		assert.Exactly(t, test.want, test.a.Equal(test.b), "Index %d", i)
	}
}

func TestRouteAppend(t *testing.T) {
	t.Parallel()
	tests := []struct {
		a       path.Route
		b       path.Route
		want    string
		wantErr error
	}{
		{path.Route("aa"), path.Route("bb/cc"), "aa/bb/cc", nil},
		{path.Route("aa"), path.Route("bbcc"), "aa/bbcc", nil},
		{path.Route("aa"), path.Route("bb\x80cc"), "", path.ErrRouteInvalidBytes},
	}
	for i, test := range tests {
		haveErr := test.a.Append(test.b)
		if test.wantErr != nil {
			assert.EqualError(t, haveErr, test.wantErr.Error(), "Index %d", i)
			continue
		}
		assert.NoError(t, haveErr, "Index %d", i)
		assert.Exactly(t, test.want, test.a.String(), "Index %d", i)
	}
}

func TestRouteTextMarshal(t *testing.T) {
	r := path.Route("admin/security/password_lifetime")
	j, err := json.Marshal(r)
	assert.NoError(t, err)
	assert.Exactly(t, "\"admin/security/password_lifetime\"", string(j))
}

func TestRouteUnmarshalTextOk(t *testing.T) {
	var r path.Route
	err := json.Unmarshal([]byte(`"admin/security/password_lifetime"`), &r)
	assert.NoError(t, err)
	assert.Exactly(t, "admin/security/password_lifetime", r.String())
}

func TestRouteLevel(t *testing.T) {
	t.Parallel()
	tests := []struct {
		have  path.Route
		level int
		want  string
	}{
		{path.Route("general/single_store_mode/enabled"), 0, ""},
		{path.Route("general/single_store_mode/enabled"), 1, "general"},
		{path.Route("general/single_store_mode/enabled"), 2, "general/single_store_mode"},
		{path.Route("general/single_store_mode/enabled"), 3, "general/single_store_mode/enabled"},
		{path.Route("general/single_store_mode/enabled"), -1, "general/single_store_mode/enabled"},
		{path.Route("general/single_store_mode/enabled"), 5, "general/single_store_mode/enabled"},
		{path.Route("general/single_store_mode/enabled"), 4, "general/single_store_mode/enabled"},
	}
	for i, test := range tests {
		r, err := test.have.Level(test.level)
		assert.NoError(t, err)
		assert.Exactly(t, test.want, r.String(), "Index %d", i)
	}
}

var benchmarkRouteLevel path.Route

// BenchmarkRouteLevel_One-4	 5000000	       297 ns/op	      16 B/op	       1 allocs/op
func BenchmarkRouteLevel_One(b *testing.B) {
	benchmarkRouteLevelRun(b, 1, path.Route("system/dev/debug"), path.Route("system"))
}

// BenchmarkRouteLevel_Two-4	 5000000	       332 ns/op	      16 B/op	       1 allocs/op
func BenchmarkRouteLevel_Two(b *testing.B) {
	benchmarkRouteLevelRun(b, 2, path.Route("system/dev/debug"), path.Route("system/dev"))
}

// BenchmarkRouteLevel_All-4	 5000000	       379 ns/op	      16 B/op	       1 allocs/op
func BenchmarkRouteLevel_All(b *testing.B) {
	benchmarkRouteLevelRun(b, -1, path.Route("system/dev/debug"), path.Route("system/dev/debug"))
}

func benchmarkRouteLevelRun(b *testing.B, level int, have, want path.Route) {
	b.ReportAllocs()
	b.ResetTimer()
	var err error
	for i := 0; i < b.N; i++ {
		benchmarkRouteLevel, err = have.Level(level)
	}
	if err != nil {
		b.Error(err)
	}
	if bytes.Equal(benchmarkRouteLevel, want) == false {
		b.Errorf("Want: %s; Have, %s", want, benchmarkRouteLevel)
	}
}

func TestRouteHash(t *testing.T) {
	t.Parallel()
	tests := []struct {
		have      path.Route
		level     int
		wantHash  uint64
		wantErr   error
		wantLevel string
	}{
		{path.Route("general/single_\x80store_mode/enabled"), 0, 0, path.ErrRouteInvalidBytes, ""},
		{path.Route("general/single_store_mode/enabled"), 0, 14695981039346656037, nil, ""},
		{path.Route("general/single_store_mode/enabled"), 1, 11396173686539659531, nil, "general"},
		{path.Route("general/single_store_mode/enabled"), 2, 12184827311064960716, nil, "general/single_store_mode"},
		{path.Route("general/single_store_mode/enabled"), 3, 8238786573751400402, nil, "general/single_store_mode/enabled"},
		{path.Route("general/single_store_mode/enabled"), -1, 8238786573751400402, nil, "general/single_store_mode/enabled"},
		{path.Route("general/single_store_mode/enabled"), 5, 8238786573751400402, nil, "general/single_store_mode/enabled"},
		{path.Route("general/single_store_mode/enabled"), 4, 8238786573751400402, nil, "general/single_store_mode/enabled"},
	}
	for i, test := range tests {

		hv, err := test.have.Hash(test.level)
		if test.wantErr != nil {
			assert.EqualError(t, err, test.wantErr.Error(), "Index %d", i)
			assert.Empty(t, hv, "Index %d", i)
			continue
		}
		assert.NoError(t, err, "Index %d", i)

		check := fnv.New64a()
		_, cErr := check.Write([]byte(test.wantLevel))
		assert.NoError(t, cErr)
		assert.Exactly(t, check.Sum64(), hv, "Index %d", i)

		l, err := test.have.Level(test.level)
		assert.Exactly(t, test.wantLevel, l.String(), "Index %d", i)
		assert.Exactly(t, test.wantHash, hv, "Index %d", i)
	}
}

var benchmarkRouteHash uint64

// BenchmarkRouteHash-4	 5000000	       288 ns/op	       0 B/op	       0 allocs/op
func BenchmarkRouteHash(b *testing.B) {
	have := path.Route("general/single_store_mode/enabled")
	want := uint64(8238786573751400402)

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

func TestRoutePartPosition(t *testing.T) {
	t.Parallel()
	tests := []struct {
		have     path.Route
		level    int
		wantPart string
		wantErr  error
	}{
		{path.Route("general/single_\x80store_mode/enabled"), 0, "", path.ErrRouteInvalidBytes},
		{path.Route("general/single_store_mode/enabled"), 0, "", path.ErrIncorrectPosition},
		{path.Route("general/single_store_mode/enabled"), 1, "general", nil},
		{path.Route("general/single_store_mode/enabled"), 2, "single_store_mode", nil},
		{path.Route("general/single_store_mode/enabled"), 3, "enabled", nil},
		{path.Route("general/single_store_mode/enabled"), -1, "", path.ErrIncorrectPosition},
		{path.Route("general/single_store_mode/enabled"), 5, "", path.ErrIncorrectPosition},
		{path.Route("general/single/store/website/group/mode/enabled/disabled/default"), 5, "store/website/group/mode/enabled/disabled/default", nil},
	}
	for i, test := range tests {
		part, haveErr := test.have.Part(test.level)
		if test.wantErr != nil {
			assert.EqualError(t, haveErr, test.wantErr.Error(), "Index %d", i)
			assert.Nil(t, part, "Index %d", i)
			continue
		}
		assert.Exactly(t, test.wantPart, part.String(), "Index %d", i)
	}
}

var benchmarkRoutePart path.Route

// BenchmarkRoutePart-4	 5000000	       240 ns/op	       0 B/op	       0 allocs/op
func BenchmarkRoutePart(b *testing.B) {
	have := path.Route("general/single_store_mode/enabled")
	want := "enabled"

	var err error
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchmarkRoutePart, err = have.Part(3)
		if err != nil {
			b.Error(err)
		}
		if benchmarkRoutePart == nil {
			b.Error("benchmarkRoutePart is nil! Unexpected")
		}
	}
	if want != benchmarkRoutePart.String() {
		b.Errorf("Want: %d; Have: %d", want, benchmarkRoutePart.String())
	}
}

func TestRouteValidate(t *testing.T) {
	t.Parallel()
	tests := []struct {
		have path.Route
		want error
	}{
		{path.Route("//"), path.ErrIncorrectPath},
		{path.Route("general/store_information/city"), nil},
		{path.Route("general/store_information/city"), nil},
		{path.Route(""), path.ErrRouteEmpty},
		{path.Route("general/store_information"), nil},
		////{path.Route(path.MustNew("system/dev/debug").Bind(scope.WebsiteID, 22).String()), path.ErrIncorrectPath},
		{path.Route("groups/33/general/store_information/street"), nil},
		{path.Route("groups/33"), nil},
		{path.Route("system/dEv/inv˚lid"), errors.New("This character \"˚\" is not allowed in Route system/dEv/inv˚lid")},
		{path.Route("system/dEv/inv'lid"), errors.New("This character \"'\" is not allowed in Route system/dEv/inv'lid")},
		{path.Route("syst3m/dEv/invalid"), nil},
		{nil, path.ErrRouteEmpty},
	}
	for i, test := range tests {
		haveErr := test.have.Validate()
		if test.want != nil {
			assert.EqualError(t, haveErr, test.want.Error(), "Index %d", i)
		} else {
			assert.NoError(t, haveErr, "Index %d", i)
		}
	}
}

var benchmarkRouteValidate error

// BenchmarkRouteValidate-4	20000000	        83.5 ns/op	       0 B/op	       0 allocs/op
func BenchmarkRouteValidate(b *testing.B) {
	have := path.Route("system/dEv/d3bug")
	want := "system/dev/debug"

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchmarkRouteValidate = have.Validate()
		if nil != benchmarkRouteValidate {
			b.Errorf("Want: %s; Have: %v", want, have)
		}
	}
}
