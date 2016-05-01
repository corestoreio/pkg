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

package storemock_test

import (
	"testing"

	"github.com/corestoreio/csfw/store"
	"github.com/corestoreio/csfw/store/storemock"
	"github.com/stretchr/testify/assert"
)

var _ store.Provider = (*storemock.NullService)(nil)

func TestNullService(t *testing.T) {

	ns := storemock.NewNullService()
	assert.False(t, ns.IsSingleStoreMode())
	assert.False(t, ns.HasSingleStore())

	ws, err := ns.Website()
	assert.Nil(t, ws)
	assert.EqualError(t, err, store.errWebsiteNotFound.Error())

	wss, err := ns.Websites()
	assert.Nil(t, wss)
	assert.EqualError(t, err, store.errWebsiteNotFound.Error())

	gs, err := ns.Group()
	assert.Nil(t, gs)
	assert.EqualError(t, err, store.errGroupNotFound.Error())

	gss, err := ns.Groups()
	assert.Nil(t, gss)
	assert.EqualError(t, err, store.errGroupNotFound.Error())

	ss, err := ns.Store()
	assert.Nil(t, ss)
	assert.EqualError(t, err, store.errStoreNotFound.Error())

	sss, err := ns.Stores()
	assert.Nil(t, sss)
	assert.EqualError(t, err, store.errStoreNotFound.Error())

	ss, err = ns.DefaultStoreView()
	assert.Nil(t, ss)
	assert.EqualError(t, err, store.errStoreNotFound.Error())

}
