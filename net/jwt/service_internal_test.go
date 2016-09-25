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

package jwt

import (
	"testing"

	"github.com/corestoreio/csfw/config/cfgmock"
	"github.com/corestoreio/csfw/store/scope"
	"github.com/corestoreio/csfw/util/csjwt"
	"github.com/corestoreio/csfw/util/errors"
	"github.com/stretchr/testify/assert"
)

func TestServiceWithBackend_NoBackend(t *testing.T) {

	jwts := MustNew()
	// a hack for testing to remove the default setting or make it invalid
	jwts.scopeCache[scope.DefaultTypeID] = &ScopedConfig{}

	cr := cfgmock.NewService()
	sc, err := jwts.ConfigByScopedGetter(cr.NewScoped(0, 0))
	assert.True(t, errors.IsNotValid(err), "Error: %+v", err)
	assert.Exactly(t, ScopedConfig{}, sc)
}

func TestServiceWithBackend_DefaultConfig(t *testing.T) {

	jwts := MustNew()

	cr := cfgmock.NewService()
	sc, err := jwts.ConfigByScopedGetter(cr.NewScoped(0, 0))
	assert.NoError(t, err, "%+v", err)
	dsc := newScopedConfig()
	if err := dsc.isValid(); err != nil {
		t.Fatal(err)
	}
	assert.Exactly(t, csjwt.HS256, sc.SigningMethod.Alg())
	assert.Exactly(t, dsc.Key.Algorithm(), sc.Key.Algorithm())

	assert.NotNil(t, dsc.ErrorHandler)
	assert.NotNil(t, sc.ErrorHandler)
	assert.NotNil(t, jwts.scopeCache[scope.DefaultTypeID].ErrorHandler)
	assert.Exactly(t, DefaultExpire, dsc.Expire)
	assert.False(t, dsc.Key.IsEmpty())
	assert.False(t, sc.Key.IsEmpty())
}
