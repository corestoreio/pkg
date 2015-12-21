// Copyright 2015, Cyrill @ Schumacher.fm and the CoreStore contributors
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

package model_test

import (
	"testing"

	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/config/model"
	"github.com/corestoreio/csfw/config/scope"
	"github.com/stretchr/testify/assert"
)

func TestBaseURL(t *testing.T) {

	wantPath := scope.StrStores.FQPathInt64(1, "web/unsecure/base_url")
	b := model.NewBaseURL("web/unsecure/base_url")

	assert.Empty(t, b.Options())

	assert.Exactly(t, "{{base_url}}", b.Get(packageConfiguration, config.NewMockGetter().NewScoped(0, 0, 1)))

	assert.Exactly(t, "http://cs.io", b.Get(packageConfiguration, config.NewMockGetter(
		config.WithMockValues(config.MockPV{
			wantPath: "http://cs.io",
		}),
	).NewScoped(0, 0, 1)))

	mw := &config.MockWrite{}
	assert.NoError(t, b.Write(mw, "dude", scope.StoreID, 1))
	assert.Exactly(t, wantPath, mw.ArgPath)
	assert.Exactly(t, "dude", mw.ArgValue.(string))

}
