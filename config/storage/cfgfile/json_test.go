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

	"github.com/corestoreio/errors"
	"github.com/corestoreio/pkg/config"
	"github.com/corestoreio/pkg/config/storage"
	"github.com/corestoreio/pkg/config/storage/cfgfile"
	"github.com/corestoreio/pkg/store/scope"
	"github.com/stretchr/testify/assert"
)

func TestWithLoadJSON(t *testing.T) {
	pUserName := config.MustNewPath("payment/stripe/user_name")

	t.Run("success", func(t *testing.T) {

		cfgSrv, err := config.NewService(
			storage.NewMap(), config.Options{},
			cfgfile.WithLoadJSON(cfgfile.WithFile("testdata", "example.json")),
		)
		if err != nil {
			t.Fatalf("%+v", err)
		}

		assert.Exactly(t, `"AUserName"`, cfgSrv.Get(pUserName).String())
		assert.Exactly(t, `"WS0Username"`, cfgSrv.Get(pUserName.BindWebsite(0)).String())
		assert.Exactly(t, `"WS1Username"`, cfgSrv.Get(pUserName.BindWebsite(1)).String())
		assert.Exactly(t, `"WS2Username"`, cfgSrv.Get(pUserName.BindWebsite(2)).String())

		assert.Exactly(t, `"SO5Username"`, cfgSrv.Get(pUserName.BindStore(5)).String())
		assert.Exactly(t, `"SO11Username"`, cfgSrv.Get(pUserName.BindStore(11)).String())

		assert.Exactly(t, `"1234"`, cfgSrv.Get(config.MustNewPath("payment/stripe/port")).String())
		assert.Exactly(t, `"true"`, cfgSrv.Get(config.MustNewPathWithScope(scope.Website.WithID(0), "payment/stripe/enable")).String())
	})

	runner := func(file string, errKind errors.Kind, errTxt string) func(*testing.T) {
		return func(t *testing.T) {
			cfgSrv, err := config.NewService(
				storage.NewMap(), config.Options{},
				cfgfile.WithLoadJSON(cfgfile.WithFile("testdata", file)),
			)
			assert.True(t, errKind.Match(err), "%+v", err)
			assert.Nil(t, cfgSrv)
			assert.Contains(t, err.Error(), errTxt)
			//t.Logf("%+v", err)
		}
	}
	t.Run("malformed_v1", runner("malformed_v1.json", errors.CorruptData,
		"[cfgfile] WithLoadJSON Unexpected data in \"payment/stripe/port\""))

	t.Run("malformed_v2_dataIF", runner("malformed_v2_dataIF.json", errors.CorruptData,
		"Unable to cast map[string]interface {}{} to []byte"))

	t.Run("malformed_v2_dataIF_scopeID", runner("malformed_v2_dataIF_scopeID.json", errors.CorruptData,
		`failed to parse "-1" to uint: strconv.ParseUint: parsing "-1": invalid syntax`))

	t.Run("malformed_v2t_dataIF", runner("malformed_v2t_dataIF.json", errors.CorruptData,
		`WithLoadJSON unexpected data in []interface {}{}`))

}
