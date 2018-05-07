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

package cfgenv

import (
	"os"
	"testing"

	"github.com/corestoreio/pkg/config"
	"github.com/corestoreio/pkg/store/scope"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var _ config.Storager = (*Storage)(nil)

func TestToEnvVar(t *testing.T) {

	tests := []struct {
		scpID scope.TypeID
		route string
		want  string
	}{
		{scope.DefaultTypeID, "aa/bb/cc", "CONFIG__AA__BB__CC"},
		{scope.Website.Pack(1), "aa/bb/cc", "CONFIG__WEBSITES__1__AA__BB__CC"},
		{scope.Store.Pack(444), "aa/bb/cc", "CONFIG__STORES__444__AA__BB__CC"},
		{scope.Store.Pack(444), "aa/bb/cc/dd/ee", "CONFIG__STORES__444__AA__BB__CC__DD__EE"},
		{scope.Store.Pack(444), "aa/bb/cc_dd/ee", "CONFIG__STORES__444__AA__BB__CC_DD__EE"},
	}

	for i, test := range tests {
		assert.Exactly(t, test.want, ToEnvVar(test.scpID, test.route), "Index %d", i)
	}
}

var benchmarkToEnvVar string

// BenchmarkToEnvVar-4   	 3000000	       570 ns/op	     163 B/op	       6 allocs/op
func BenchmarkToEnvVar(b *testing.B) {
	scpID := scope.Store.Pack(543)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchmarkToEnvVar = ToEnvVar(scpID, "aa/bb/cc")
	}
}

// BenchmarkFromEnvVar-4   	 3000000	       501 ns/op	     184 B/op	       7 allocs/op
func BenchmarkFromEnvVar(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, benchmarkToEnvVar = FromEnvVar(Prefix, "CONFIG__WEBSITES__321__AA__BB__CC")
	}
}

func TestFromEnvVar(t *testing.T) {
	tests := []struct {
		envVar    string
		wantScpID scope.TypeID
		wantRoute string
	}{
		{"CONFIG__AA__BB__CC", scope.DefaultTypeID, "aa/bb/cc"},
		{"CONFIG__AA__BB__CC__DD", 0, ""},
		{"CONFIG__AA__BB__CC_DD", scope.DefaultTypeID, "aa/bb/cc_dd"},
		{"CONFIG__WEBSITES__321__AA__BB__CC", scope.Website.Pack(321), "aa/bb/cc"},
		{"CONFIG__STORES__1__AA__BB__CC", scope.Store.Pack(1), "aa/bb/cc"},
		{"CONFIG__STORES__AA__BB__CC", 0, ""},
		{"ONFIG__STORES__AA__BB__CC", 0, ""},
		{"CONFIG__", 0, ""},
		{"ONFIG__", 0, ""},
		{"", 0, ""},
	}
	for i, test := range tests {
		haveScpID, haveRoute := FromEnvVar(Prefix, test.envVar)
		assert.Exactly(t, test.wantScpID, haveScpID, "Index %d", i)
		assert.Exactly(t, test.wantRoute, haveRoute, "Index %d", i)
	}
}

func validateFoundGet(t *testing.T, s config.Storager, scp scope.TypeID, route string, want string) {
	data, ok, err := s.Get(scp, route)
	require.NoError(t, err)
	assert.True(t, ok, "env value must be found")
	assert.Exactly(t, []byte(want), data)
}

func TestStorage_No_Preload(t *testing.T) {
	s, err := NewStorage(Options{
		Preload: false,
	})
	if err != nil {
		t.Fatal(err)
	}

	t.Run("empty write returns nil error", func(t *testing.T) {
		assert.NoError(t, s.Set(0, "", nil))
	})

	runner := func(envVar string, scp scope.TypeID, route string) func(*testing.T) {
		return func(t *testing.T) {
			defer func() { assert.NoError(t, os.Unsetenv(envVar)) }()
			require.NoError(t, os.Setenv(envVar, "DATA from ENV"))

			validateFoundGet(t, s, scp, route, `DATA from ENV`)
		}
	}
	t.Run("default scope", runner("CONFIG__AA__BB__CC", scope.DefaultTypeID, "aa/bb/cc"))
	t.Run("website 123 scope", runner("CONFIG__WEBSITES__1__AA__BB__CC", scope.Website.Pack(1), "aa/bb/cc"))
	t.Run("store 444 scope", runner("CONFIG__STORES__444__AA__BB__CC_DD__EE", scope.Store.Pack(444), "aa/bb/cc_dd/ee"))
	t.Run("wrong path with special symbols", func(t *testing.T) {
		envVar := "CONFIG____€__∏"
		defer func() { assert.NoError(t, os.Unsetenv(envVar)) }()
		require.NoError(t, os.Setenv(envVar, "DATA from ENV"))
		data, ok, err := s.Get(scope.DefaultTypeID, "aa/bb/cc")
		require.NoError(t, err)
		assert.False(t, ok, "env value must be found")
		assert.Nil(t, data)
	})
}

var benchmarkStorage []byte

// BenchmarkStorage-4      	 1000000	      1084 ns/op	     259 B/op	       7 allocs/op
func BenchmarkStorage_No_Preload(b *testing.B) {
	s, err := NewStorage(Options{
		Preload: false,
	})
	if err != nil {
		b.Fatal(err)
	}

	if err := os.Setenv("CONFIG__STORES__444__AA__BB__CC_DD__EE", "DATA from ENV"); err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var ok bool
		var err error
		benchmarkStorage, ok, err = s.Get(scope.Store.Pack(444), "aa/bb/cc_dd/ee")
		if err != nil {
			b.Fatal(err)
		}
		if !ok {
			b.Fatal("value must be found")
		}
	}
	if err := os.Unsetenv("CONFIG__STORES__444__AA__BB__CC_DD__EE"); err != nil {
		b.Fatal(err)
	}
}

func TestStorage_No_Preload_UnsetEnvAfterRead_And_Cache(t *testing.T) {
	s, err := NewStorage(Options{
		UnsetEnvAfterRead: true,
		Preload:           false,
		CacheVariableFn:   func(scp scope.TypeID, route string) bool { return true }, // cache all
	})
	if err != nil {
		t.Fatal(err)
	}
	const wantValue = "Banana 🍌"
	os.Setenv("CONFIG__WEBSITES__159__AA__BB__CC", wantValue)

	validateFoundGet(t, s, scope.Website.Pack(159), "aa/bb/cc", wantValue)

	ev, eOK := os.LookupEnv("CONFIG__WEBSITES__159__AA__BB__CC")
	assert.False(t, eOK, "Env var must be unset and not found")
	assert.Empty(t, ev, "Env var must be empty")

	// Read from cache
	validateFoundGet(t, s, scope.Website.Pack(159), "aa/bb/cc", wantValue)
}

func TestStorage_With_Preload_UnsetEnvAfterRead(t *testing.T) {
	const wantValue = "Pear 🍐"
	os.Setenv("CONFIG__STORES__345__XX__BB__CC", wantValue)
	os.Setenv("CONFIG__STORES__345__XY__BB__CC", "")

	s, err := NewStorage(Options{
		UnsetEnvAfterRead: true,
		Preload:           true,
		CacheVariableFn:   func(scp scope.TypeID, route string) bool { panic("Should not get called") }, // cache all
	})
	if err != nil {
		t.Fatal(err)
	}
	validateFoundGet(t, s, scope.Store.Pack(345), "xx/bb/cc", wantValue)
	validateFoundGet(t, s, scope.Store.Pack(345), "xy/bb/cc", "")

	ev, eOK := os.LookupEnv("CONFIG__STORES__345__XX__BB__CC")
	assert.False(t, eOK, "Env var must be unset and not found")
	assert.Empty(t, ev, "Env var must be empty")

}
