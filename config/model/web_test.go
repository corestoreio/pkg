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

package model_test

import (
	"testing"

	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/config/model"
	"github.com/corestoreio/csfw/config/path"
	"github.com/corestoreio/csfw/store/scope"
	"github.com/stretchr/testify/assert"
)

func TestBaseURLGet(t *testing.T) {
	t.Parallel()
	const pathWebUnsecUrl = "web/unsecure/base_url"
	wantPath := path.MustNewByParts(pathWebUnsecUrl).Bind(scope.StoreID, 1)
	b := model.NewBaseURL(pathWebUnsecUrl, model.WithConfigStructure(configStructure))

	assert.Empty(t, b.Options())

	sg, err := b.Get(config.NewMockGetter().NewScoped(0, 0, 1))
	if err != nil {
		t.Fatal(err)
	}
	assert.Exactly(t, "{{base_url}}", sg)

	sg, err = b.Get(config.NewMockGetter(
		config.WithMockValues(config.MockPV{
			wantPath.String(): "http://cs.io",
		}),
	).NewScoped(0, 0, 1))
	if err != nil {
		t.Fatal(err)
	}
	assert.Exactly(t, "http://cs.io", sg)

}

func TestBaseURLWrite(t *testing.T) {
	t.Parallel()
	const pathWebUnsecUrl = "web/unsecure/base_url"
	wantPath := path.MustNewByParts(pathWebUnsecUrl).Bind(scope.StoreID, 1)
	b := model.NewBaseURL(pathWebUnsecUrl, model.WithConfigStructure(configStructure))

	mw := &config.MockWrite{}
	assert.NoError(t, b.Write(mw, "dude", scope.StoreID, 1))
	assert.Exactly(t, wantPath.String(), mw.ArgPath)
	assert.Exactly(t, "dude", mw.ArgValue.(string))
}
