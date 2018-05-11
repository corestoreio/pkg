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

package storage_test

import (
	"os"
	"testing"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/pkg/config"
	"github.com/corestoreio/pkg/config/storage"
	"github.com/corestoreio/pkg/store/scope"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
		assert.Exactly(t, test.want, storage.ToEnvVar(config.MustNewPathWithScope(test.scpID, test.route)), "Index %d", i)
	}
}

func TestFromEnvVar(t *testing.T) {
	tests := []struct {
		envVar   string
		wantPath string
		wantErr  errors.Kind
	}{
		{"CONFIG__AA__BB__CC", "default/0/aa/bb/cc", errors.NoKind},
		{"CONFIG__AA__BB__CC__DD", "", errors.NotValid}, // errors.NotValid
		{"CONFIG__AA__BB__CC_DD", "default/0/aa/bb/cc_dd", errors.NoKind},
		{"CONFIG__WEBSITES__321__AA__BB__CC", "websites/321/aa/bb/cc", errors.NoKind},
		{"CONFIG__STORES__1__AA__BB__CC", "stores/1/aa/bb/cc", errors.NoKind},
		{"CONFIG__STORES__AA__BB__CC", "", errors.NotValid},
		{"ONFIG__STORES__AA__BB__CC", "", errors.NotValid},
		{"CONFIG__", "", errors.NotValid},
		{"ONFIG__", "", errors.NotValid},
		{"", "", errors.NotValid},
	}
	for i, test := range tests {
		haveP, haveErr := storage.FromEnvVar(storage.Prefix, test.envVar)
		if test.wantErr > 0 {
			assert.Nil(t, haveP)
			assert.True(t, test.wantErr.Match(haveErr), "%d: Kind %q\n%+v", i, errors.UnwrapKind(haveErr).String(), haveErr)
		} else {
			require.NoError(t, haveErr, "Index %d. Kind %q", i, errors.UnwrapKind(haveErr).String())
			assert.Exactly(t, test.wantPath, haveP.String(), "Index %d", i)
		}
	}
}

func TestStorage_No_Preload(t *testing.T) {
	s, err := storage.NewEnvironment(storage.EnvOp{
		Preload: false,
	})
	if err != nil {
		t.Fatal(err)
	}

	t.Run("empty write returns nil error", func(t *testing.T) {
		assert.NoError(t, s.Set(new(config.Path), nil))
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

		validateNotFoundGet(t, s, scope.DefaultTypeID, "aa/bb/cc")
	})
}

func TestStorage_No_Preload_UnsetEnvAfterRead_And_Cache(t *testing.T) {
	s, err := storage.NewEnvironment(storage.EnvOp{
		UnsetEnvAfterRead: true,
		Preload:           false,
		CacheVariableFn:   func(*config.Path) bool { return true }, // cache all
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

	s, err := storage.NewEnvironment(storage.EnvOp{
		UnsetEnvAfterRead: true,
		Preload:           true,
		CacheVariableFn:   nil, // cache all
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
