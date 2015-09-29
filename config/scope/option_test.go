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

package scope_test

import (
	"testing"

	"github.com/corestoreio/csfw/config/scope"
	"github.com/stretchr/testify/assert"
)

func TestApplyCode(t *testing.T) {
	tests := []struct {
		wantStoreCode   string
		wantWebsiteCode string
		haveCode        string
		s               scope.Scope
		err             error
	}{
		{"", "de1", "de1", scope.WebsiteID, nil},
		{"de2", "", "de2", scope.StoreID, nil},
		{"", "", "de3", scope.GroupID, scope.ErrUnsupportedScope},
		{"", "", "de4", scope.AbsentID, scope.ErrUnsupportedScope},
	}

	for _, test := range tests {
		so, err := scope.SetByCode(test.haveCode, test.s)
		assert.NotNil(t, so)
		if test.err != nil {
			assert.EqualError(t, err, test.err.Error())
		} else {
			assert.NoError(t, err)
			assert.Equal(t, test.s, so.Scope())
			assert.Equal(t, test.wantStoreCode, so.StoreCode())
			assert.Equal(t, test.wantWebsiteCode, so.WebsiteCode())
		}
	}
}

func TestApplyID(t *testing.T) {
	tests := []struct {
		wantWebsiteID scope.WebsiteIDer
		wantGroupID   scope.GroupIDer
		wantStoreID   scope.StoreIDer

		haveID int64
		s      scope.Scope
		err    error
	}{
		{scope.MockID(1), nil, nil, 1, scope.WebsiteID, nil},
		{nil, scope.MockID(3), nil, 3, scope.GroupID, nil},
		{nil, nil, scope.MockID(2), 2, scope.StoreID, nil},
		{nil, nil, nil, 4, scope.AbsentID, scope.ErrUnsupportedScope},
	}

	for _, test := range tests {
		so, err := scope.SetByID(test.haveID, test.s)
		assert.NotNil(t, so)
		if test.err != nil {
			assert.EqualError(t, err, test.err.Error())
			assert.Nil(t, so.Website)
			assert.Nil(t, so.Group)
			assert.Nil(t, so.Store)
		} else {
			assert.NoError(t, err)
			assert.Equal(t, test.s, so.Scope())
			assert.Equal(t, "", so.StoreCode())
			assert.Equal(t, "", so.WebsiteCode())

			if test.wantWebsiteID != nil {
				assert.Equal(t, test.wantWebsiteID.WebsiteID(), so.Website.WebsiteID())
			} else {
				assert.Nil(t, test.wantWebsiteID)
			}

			if test.wantGroupID != nil {
				assert.Equal(t, test.wantGroupID.GroupID(), so.Group.GroupID())
			} else {
				assert.Nil(t, test.wantGroupID)
			}

			if test.wantStoreID != nil {
				assert.Equal(t, test.wantStoreID.StoreID(), so.Store.StoreID())
			} else {
				assert.Nil(t, test.wantStoreID)
			}

		}
	}
}

func TestApplyWebsite(t *testing.T) {

	so := scope.Option{Website: scope.MockID(3)}
	assert.NotNil(t, so)
	assert.Equal(t, int64(3), so.Website.WebsiteID())
	assert.Nil(t, so.Group)
	assert.Nil(t, so.Store)
}

func TestApplyGroup(t *testing.T) {

	so := scope.Option{Group: scope.MockID(3)}
	assert.NotNil(t, so)

	assert.Equal(t, int64(3), so.Group.GroupID())
	assert.Nil(t, so.Website)
	assert.Nil(t, so.Store)
}

func TestApplyStore(t *testing.T) {

	so := scope.Option{Store: scope.MockID(3)}
	assert.NotNil(t, so)

	assert.Equal(t, int64(3), so.Store.StoreID())
	assert.Nil(t, so.Website)
	assert.Nil(t, so.Group)

}
