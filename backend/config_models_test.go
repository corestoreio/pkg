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

package backend_test

import (
	"testing"

	"github.com/corestoreio/csfw/backend"
	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/config/model"
	"github.com/corestoreio/csfw/store/scope"
	"github.com/stretchr/testify/assert"
)

func TestConfigRedirectToBase(t *testing.T) {
	defer debugLogBuf.Reset()
	t.Parallel()

	r := backend.NewConfigRedirectToBase(
		backend.Backend.WebURLRedirectToBase.String(),
		model.WithConfigStructure(backend.ConfigStructure),
	)

	cr := config.NewMockGetter(
		config.WithMockValues(config.MockPV{
			backend.Backend.WebURLRedirectToBase.String(): 2,
		}),
	)

	code := r.Get(cr.NewScoped(0, 0, 0))
	assert.Exactly(t, 2, code)
	code = r.Get(cr.NewScoped(1, 1, 2))
	assert.Exactly(t, 0, code)

	// that is crap we should return an error
	assert.Contains(t, debugLogBuf.String(), "Scope permission insufficient: Have 'Store'; Want 'Default'")

	mw := new(config.MockWrite)
	assert.EqualError(t, r.Write(mw, 200, scope.DefaultID, 0),
		"Cannot find 200 in list: [{\"Value\":0,\"Label\":\"No\"},{\"Value\":1,\"Label\":\"Yes (302 Found)\"},{\"Value\":302,\"Label\":\"Yes (302 Found)\"},{\"Value\":301,\"Label\":\"Yes (301 Moved Permanently)\"}]\n",
	) // 200 not allowed

}
