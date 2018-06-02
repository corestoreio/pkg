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
		{scope.Website.WithID(1), "aa/bb/cc", "CONFIG__WEBSITES__1__AA__BB__CC"},
		{scope.Store.WithID(444), "aa/bb/cc", "CONFIG__STORES__444__AA__BB__CC"},
		{scope.Store.WithID(444), "aa/bb/cc/dd/ee", "CONFIG__STORES__444__AA__BB__CC__DD__EE"},
		{scope.Store.WithID(444), "aa/bb/cc_dd/ee", "CONFIG__STORES__444__AA__BB__CC_DD__EE"},
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
		{"CONFIG__AA__BB__CC__DD", "default/0/aa/bb/cc/dd", errors.NoKind},
		{"CONFIG__AA__BB__CC_DD", "default/0/aa/bb/cc_dd", errors.NoKind},
		{"CONFIG__WEBSITES__321__AA__BB__CC", "websites/321/aa/bb/cc", errors.NoKind},
		{"CONFIG__STORES__1__AA__BB__CC", "stores/1/aa/bb/cc", errors.NoKind},
		{"CONFIG__STORES__AA__BB__CC", "", errors.NotValid},
		{"ONFIG__STORES__AA__BB__CC", "default/0/tores/aa/bb/cc", errors.NoKind},
		{"CONFIG__", "", errors.NotValid},
		{"ONFIG__", "", errors.NotValid},
		{"", "", errors.NotValid},
	}
	for i, test := range tests {
		haveP, haveErr := storage.FromEnvVar(storage.Prefix, test.envVar)
		if test.wantErr > 0 {
			assert.Nil(t, haveP, "Index %d", i)
			assert.True(t, test.wantErr.Match(haveErr), "%d: Kind %q\n%+v", i, errors.UnwrapKind(haveErr).String(), haveErr)
		} else {
			require.NoError(t, haveErr, "Index %d. Kind %q", i, errors.UnwrapKind(haveErr).String())
			assert.Exactly(t, test.wantPath, haveP.String(), "Index %d", i)
		}
	}
}

func TestWithLoadEnvironmentVariables(t *testing.T) {

	t.Run("success", func(t *testing.T) {
		const wantValue = "Pear 🍐"
		assert.NoError(t, os.Setenv("CONFIG__STORES__345__XX__BB__CC", wantValue))
		assert.NoError(t, os.Setenv("CONFIG__STORES__345__XY__BB__CC", ""))
		assert.NoError(t, os.Setenv("CONFIG__GENERAL__GPS__LAT", "42.1234"))
		defer func() {
			assert.NoError(t, os.Unsetenv("CONFIG__STORES__345__XX__BB__CC"))
			assert.NoError(t, os.Unsetenv("CONFIG__STORES__345__XY__BB__CC"))
			assert.NoError(t, os.Unsetenv("CONFIG__GENERAL__GPS__LAT"))
		}()

		cfgSrv, err := config.NewService(
			storage.NewMap(), config.Options{},
			storage.WithLoadEnvironmentVariables(storage.EnvOp{}),
		)
		if err != nil {
			t.Fatalf("%+v", err)
		}

		assert.Exactly(t, "\"Pear 🍐\"", cfgSrv.Get(config.MustNewPathWithScope(scope.Store.WithID(345), "xx/bb/cc")).String())
		assert.Exactly(t, `""`, cfgSrv.Get(config.MustNewPathWithScope(scope.Store.WithID(345), "xy/bb/cc")).String())
		assert.Exactly(t, 42.1234, cfgSrv.Get(config.MustNewPath("general/gps/lat")).UnsafeFloat64())
	})

	t.Run("malformed path", func(t *testing.T) {
		assert.NoError(t, os.Setenv("CONFIG__BB__CC", "x"))
		defer func() {
			assert.NoError(t, os.Unsetenv("CONFIG__BB__CC"))
		}()

		cfgSrv, err := config.NewService(
			storage.NewMap(), config.Options{},
			storage.WithLoadEnvironmentVariables(storage.EnvOp{}),
		)
		assert.Nil(t, cfgSrv)
		assert.True(t, errors.NotValid.Match(err))
		assert.EqualError(t, err, "[config] Expecting: `aa/bb/cc` or `strScope/ID/aa/bb/cc` but got \"bb/cc\"`")
	})
}
