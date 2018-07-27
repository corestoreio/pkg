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

package cfgfile_test

import (
	"testing"

	"github.com/alecthomas/assert"
	"github.com/corestoreio/errors"
	"github.com/corestoreio/pkg/config"
	"github.com/corestoreio/pkg/config/storage"
	"github.com/corestoreio/pkg/config/storage/cfgfile"
	"github.com/corestoreio/pkg/store/scope"
	"github.com/fortytw2/leaktest"
)

func TestWithLoadYAML(t *testing.T) {

	t.Run("success", func(t *testing.T) {

		cfgSrv, err := config.NewService(
			storage.NewMap(), config.Options{},
			cfgfile.WithLoadYAML(cfgfile.WithFiles([]string{"testdata", "example.yaml"})),
		)
		if err != nil {
			t.Fatalf("%+v", err)
		}

		assert.Exactly(t, `"false"`, cfgSrv.Get(config.MustNewPathWithScope(scope.Website.WithID(2), "google/analytics/active")).String())
		assert.Exactly(t, 2.002, cfgSrv.Get(config.MustNewPathWithScope(scope.Store.WithID(2), "dev/js/merge_files")).UnsafeFloat64())
		assert.Exactly(t, `"http://eshop.dev/"`, cfgSrv.Get(config.MustNewPath("web/unsecure/base_url")).String())
	})

	t.Run("malformed path", func(t *testing.T) {

		cfgSrv, err := config.NewService(
			storage.NewMap(), config.Options{},
			cfgfile.WithLoadYAML(cfgfile.WithFiles([]string{"testdata", "malformed_path.yaml"})),
		)
		assert.Nil(t, cfgSrv)
		assert.True(t, errors.NotValid.Match(err))
		assert.EqualError(t, err, "[config] Invalid Path \"vendorbarenvironment\". Either to short or missing path separator.")
	})

	t.Run("malformed yaml", func(t *testing.T) {

		cfgSrv, err := config.NewService(
			storage.NewMap(), config.Options{},
			cfgfile.WithLoadYAML(cfgfile.WithFiles([]string{"testdata", "malformed_yaml.yaml"})),
		)
		assert.Nil(t, cfgSrv)
		assert.EqualError(t, err, "yaml: unmarshal errors:\n  line 2: cannot unmarshal !!str `192.168...` into map[string]string")
	})
}

func TestWithLoadFieldMetaYAML(t *testing.T) {
	defer leaktest.Check(t)()

	t.Run("success", func(t *testing.T) {

		cfgSrv, err := config.NewService(
			storage.NewMap(), config.Options{},
			cfgfile.WithLoadFieldMetaYAML(cfgfile.WithFiles([]string{"testdata", "example_field_meta.yaml"})),
		)
		if err != nil {
			t.Fatalf("%+v", err)
		}

		scpd13 := cfgSrv.Scoped(1, 3)
		scpd24 := cfgSrv.Scoped(2, 4)
		assert.Exactly(t, `"8080"`, scpd13.Get(scope.Default, "carrier/dpd/port").String())
		assert.Exactly(t, `"60s"`, scpd13.Get(scope.Default, "carrier/dpd/timeout").String())
		assert.Exactly(t, `"50s"`, scpd13.Get(scope.Website, "carrier/dpd/timeout").String())
		assert.Exactly(t, `"40s"`, scpd24.Get(scope.Website, "carrier/dpd/timeout").String())
		assert.Exactly(t, `"prdUser0"`, scpd13.Get(scope.Website, "carrier/dpd/username").String())
		assert.Exactly(t, `"prdUser1"`, scpd13.Get(scope.Store, "carrier/dpd/username").String())
		assert.Exactly(t, `"prdUser2"`, scpd24.Get(scope.Store, "carrier/dpd/username").String())

		err = cfgSrv.Set(config.MustNewPath("carrier/dpd/port").BindWebsite(1), []byte(`return error`))
		assert.True(t, errors.NotAllowed.Match(err), "%+v", err)
		err = cfgSrv.Set(config.MustNewPath("carrier/dpd/timeout").BindStore(1), []byte(`return error`))
		assert.True(t, errors.NotAllowed.Match(err), "%+v", err)
	})

	t.Run("malformed yaml", func(t *testing.T) {
		cfgSrv, err := config.NewService(
			storage.NewMap(), config.Options{},
			cfgfile.WithLoadFieldMetaYAML(cfgfile.WithFiles([]string{"testdata", "example.yaml"})),
		)
		assert.Nil(t, cfgSrv)
		assert.True(t, errors.Fatal.Match(err), "%+v", err)
	})
	t.Run("malformed perm", func(t *testing.T) {
		cfgSrv, err := config.NewService(
			storage.NewMap(), config.Options{},
			cfgfile.WithLoadFieldMetaYAML(cfgfile.WithFiles([]string{"testdata", "malformed_field_meta.yaml"})),
		)
		assert.Nil(t, cfgSrv)
		assert.True(t, errors.NotSupported.Match(err), "%+v", err)
	})
	t.Run("file not found", func(t *testing.T) {
		cfgSrv, err := config.NewService(
			storage.NewMap(), config.Options{},
			cfgfile.WithLoadFieldMetaYAML(cfgfile.WithFiles([]string{"testdata", "malformed_field_meta_XXXZ.yaml"})),
		)
		assert.Nil(t, cfgSrv)
		assert.True(t, errors.NotFound.Match(err), "%+v", err)
	})
}
