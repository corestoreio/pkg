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

package store_test

import (
	"testing"

	"github.com/corestoreio/csfw/config/scope"
	"github.com/corestoreio/csfw/storage/dbr"
	"github.com/corestoreio/csfw/store"
	storemock "github.com/corestoreio/csfw/store/mock"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

func TestContextReaderError(t *testing.T) {
	haveMr, s, err := store.FromContextReader(context.Background())
	assert.Nil(t, haveMr)
	assert.Nil(t, s)
	assert.EqualError(t, err, store.ErrContextServiceNotFound.Error())

	ctx := store.NewContextReader(context.Background(), nil)
	assert.NotNil(t, ctx)
	haveMr, s, err = store.FromContextReader(ctx)
	assert.Nil(t, haveMr)
	assert.Nil(t, s)
	assert.EqualError(t, err, store.ErrContextServiceNotFound.Error())

	mr := storemock.NewNullService()
	ctx = store.NewContextReader(context.Background(), mr)
	assert.NotNil(t, ctx)
	haveMr, s, err = store.FromContextReader(ctx)
	assert.EqualError(t, err, store.ErrStoreNotFound.Error())
	assert.Nil(t, haveMr)
	assert.Nil(t, s)

}

func TestContextReaderSuccess(t *testing.T) {
	ctx := storemock.NewContextService(scope.Option{},
		func(ms *storemock.Storage) {
			ms.MockStore = func() (*store.Store, error) {
				return store.NewStore(
					&store.TableStore{StoreID: 6, Code: dbr.NewNullString("de"), WebsiteID: 1, GroupID: 1, Name: "Germany", SortOrder: 10, IsActive: true},
					&store.TableWebsite{WebsiteID: 1, Code: dbr.NewNullString("euro"), Name: dbr.NewNullString("Europe"), SortOrder: 0, DefaultGroupID: 1, IsDefault: dbr.NewNullBool(true)},
					&store.TableGroup{GroupID: 1, WebsiteID: 1, Name: "DACH Group", RootCategoryID: 2, DefaultStoreID: 6},
				)
			}
		},
	)

	haveMr, s, err := store.FromContextReader(ctx)
	assert.NoError(t, err)
	assert.Exactly(t, int64(6), s.StoreID())

	s2, err2 := haveMr.Store()
	assert.NoError(t, err2)
	assert.Exactly(t, int64(6), s2.StoreID())

}
