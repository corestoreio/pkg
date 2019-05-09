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

package mock

import (
	"github.com/corestoreio/pkg/storage/null"
	"github.com/corestoreio/pkg/store"
)

// NewServiceEuroOZ creates a fully initialized store.Service with 3 websites,
// 4 groups and 7 stores used for testing. Panics on error. Website 1 contains
// Europe and website 2 contains Australia/New Zealand.
func NewServiceEuroOZ(opts ...store.Option) *store.Service {
	defaultOpts := []store.Option{
		store.WithWebsites(
			&store.StoreWebsite{WebsiteID: 0, Code: `admin`, Name: null.MakeString(`Admin`), SortOrder: 0, DefaultGroupID: 0, IsDefault: false},
			&store.StoreWebsite{WebsiteID: 1, Code: `euro`, Name: null.MakeString(`Europe`), SortOrder: 0, DefaultGroupID: 1, IsDefault: true},
			&store.StoreWebsite{WebsiteID: 2, Code: `oz`, Name: null.MakeString(`OZ`), SortOrder: 20, DefaultGroupID: 3, IsDefault: false},
		),
		store.WithGroups(
			&store.StoreGroup{GroupID: 3, WebsiteID: 2, Name: `Australia`, RootCategoryID: 2, DefaultStoreID: 5},
			&store.StoreGroup{GroupID: 1, WebsiteID: 1, Name: `DACH Group`, RootCategoryID: 2, DefaultStoreID: 2},
			&store.StoreGroup{GroupID: 0, WebsiteID: 0, Name: `Default`, RootCategoryID: 0, DefaultStoreID: 0},
			&store.StoreGroup{GroupID: 2, WebsiteID: 1, Name: `UK Group`, RootCategoryID: 2, DefaultStoreID: 4},
		),
		store.WithStores(
			&store.Store{StoreID: 0, Code: `admin`, WebsiteID: 0, GroupID: 0, Name: `Admin`, SortOrder: 0, IsActive: true},
			&store.Store{StoreID: 5, Code: `au`, WebsiteID: 2, GroupID: 3, Name: `Australia`, SortOrder: 10, IsActive: true},
			&store.Store{StoreID: 1, Code: `de`, WebsiteID: 1, GroupID: 1, Name: `Germany`, SortOrder: 10, IsActive: true},
			&store.Store{StoreID: 4, Code: `uk`, WebsiteID: 1, GroupID: 2, Name: `UK`, SortOrder: 10, IsActive: true},
			&store.Store{StoreID: 2, Code: `at`, WebsiteID: 1, GroupID: 1, Name: `Österreich`, SortOrder: 20, IsActive: true},
			&store.Store{StoreID: 6, Code: `nz`, WebsiteID: 2, GroupID: 3, Name: `Kiwi`, SortOrder: 30, IsActive: true},
			&store.Store{IsActive: false, StoreID: 3, Code: `ch`, WebsiteID: 1, GroupID: 1, Name: `Schweiz`, SortOrder: 30},
		),
	}
	return store.MustNewService(append(defaultOpts, opts...)...)
}

func NewServiceEuroW11G11S19(opts ...store.Option) *store.Service {
	defaultOpts := []store.Option{
		store.WithWebsites(
			&store.StoreWebsite{Code: `admin`, Name: null.MakeString(`Admin`)},
			&store.StoreWebsite{WebsiteID: 3, Code: `ch`, Name: null.MakeString(`Schweiz`), SortOrder: 2, DefaultGroupID: 3},
			&store.StoreWebsite{WebsiteID: 8, Code: `at`, Name: null.MakeString(`Österreich`), SortOrder: 7, DefaultGroupID: 8},
			&store.StoreWebsite{WebsiteID: 5, Code: `fr`, Name: null.MakeString(`Frankreich`), SortOrder: 4, DefaultGroupID: 5},
			&store.StoreWebsite{WebsiteID: 2, Code: `de`, Name: null.MakeString(`Deutschland`), SortOrder: 1, DefaultGroupID: 2, IsDefault: true},
			&store.StoreWebsite{WebsiteID: 4, Code: `it`, Name: null.MakeString(`Italien`), SortOrder: 3, DefaultGroupID: 4},
			&store.StoreWebsite{WebsiteID: 6, Code: `be`, Name: null.MakeString(`Belgien`), SortOrder: 5, DefaultGroupID: 6},
			&store.StoreWebsite{WebsiteID: 9, Code: `int`, Name: null.MakeString(`International`), SortOrder: 8, DefaultGroupID: 9},
			&store.StoreWebsite{WebsiteID: 10, Code: `nl`, Name: null.MakeString(`Netherlands`), SortOrder: 9, DefaultGroupID: 10},
			&store.StoreWebsite{WebsiteID: 11, Code: `uk`, Name: null.MakeString(`United Kingdom`), SortOrder: 10, DefaultGroupID: 11},
			&store.StoreWebsite{WebsiteID: 7, Code: `lu`, Name: null.MakeString(`Luxemburg`), SortOrder: 6, DefaultGroupID: 7},
		),
		store.WithGroups(
			&store.StoreGroup{Name: `Default`, Code: null.String{}},
			&store.StoreGroup{GroupID: 2, WebsiteID: 2, Name: `b2c`, RootCategoryID: 2, DefaultStoreID: 2},
			&store.StoreGroup{GroupID: 7, WebsiteID: 7, Name: `b2c`, RootCategoryID: 2, DefaultStoreID: 15},
			&store.StoreGroup{GroupID: 8, WebsiteID: 8, Name: `b2c`, RootCategoryID: 2, DefaultStoreID: 17},
			&store.StoreGroup{GroupID: 11, WebsiteID: 11, Name: `b2c`, RootCategoryID: 2, DefaultStoreID: 20},
			&store.StoreGroup{GroupID: 9, WebsiteID: 9, Name: `b2c`, RootCategoryID: 2, DefaultStoreID: 18},
			&store.StoreGroup{GroupID: 6, WebsiteID: 6, Name: `b2c`, RootCategoryID: 2, DefaultStoreID: 13},
			&store.StoreGroup{GroupID: 4, WebsiteID: 4, Name: `b2c`, RootCategoryID: 2, DefaultStoreID: 10},
			&store.StoreGroup{GroupID: 5, WebsiteID: 5, Name: `b2c`, RootCategoryID: 2, DefaultStoreID: 12},
			&store.StoreGroup{GroupID: 3, WebsiteID: 3, Name: `b2c`, RootCategoryID: 2, DefaultStoreID: 6},
			&store.StoreGroup{GroupID: 10, WebsiteID: 10, Name: `b2c`, RootCategoryID: 2, DefaultStoreID: 19},
		),
		store.WithStores(
			&store.Store{Code: `admin`, Name: `Admin`, IsActive: true},
			&store.Store{StoreID: 16, Code: `lude`, WebsiteID: 7, GroupID: 7, Name: `de`, SortOrder: 2, IsActive: true},
			&store.Store{StoreID: 3, Code: `detr`, WebsiteID: 2, GroupID: 2, Name: `tr`, SortOrder: 4, IsActive: true},
			&store.Store{StoreID: 18, Code: `inten`, WebsiteID: 9, GroupID: 9, Name: `en`, SortOrder: 1, IsActive: true},
			&store.Store{StoreID: 12, Code: `frfr`, WebsiteID: 5, GroupID: 5, Name: `fr`, SortOrder: 1, IsActive: true},
			&store.Store{StoreID: 2, Code: `dede`, WebsiteID: 2, GroupID: 2, Name: `de`, SortOrder: 1, IsActive: true},
			&store.Store{StoreID: 10, Code: `itit`, WebsiteID: 4, GroupID: 4, Name: `it`, SortOrder: 1, IsActive: true},
			&store.Store{StoreID: 20, Code: `uken`, WebsiteID: 11, GroupID: 11, Name: `en`, SortOrder: 1, IsActive: true},
			&store.Store{StoreID: 15, Code: `lufr`, WebsiteID: 7, GroupID: 7, Name: `fr`, SortOrder: 1, IsActive: true},
			&store.Store{StoreID: 8, Code: `chit`, WebsiteID: 3, GroupID: 3, Name: `it`, SortOrder: 3, IsActive: true},
			&store.Store{StoreID: 19, Code: `nlen`, WebsiteID: 10, GroupID: 10, Name: `en`, SortOrder: 1, IsActive: true},
			&store.Store{StoreID: 6, Code: `chde`, WebsiteID: 3, GroupID: 3, Name: `de`, SortOrder: 1, IsActive: true},
			&store.Store{StoreID: 7, Code: `chfr`, WebsiteID: 3, GroupID: 3, Name: `fr`, SortOrder: 2, IsActive: true},
			&store.Store{StoreID: 17, Code: `atde`, WebsiteID: 8, GroupID: 8, Name: `de`, SortOrder: 1, IsActive: true},
			&store.Store{StoreID: 5, Code: `deen`, WebsiteID: 2, GroupID: 2, Name: `en`, SortOrder: 4, IsActive: true},
			&store.Store{StoreID: 9, Code: `chen`, WebsiteID: 3, GroupID: 3, Name: `en`, SortOrder: 4, IsActive: true},
			&store.Store{StoreID: 14, Code: `been`, WebsiteID: 6, GroupID: 6, Name: `en`, SortOrder: 2, IsActive: true},
			&store.Store{StoreID: 13, Code: `befr`, WebsiteID: 6, GroupID: 6, Name: `fr`, SortOrder: 1, IsActive: true},
			&store.Store{StoreID: 11, Code: `itde`, WebsiteID: 4, GroupID: 4, Name: `de`, SortOrder: 2, IsActive: true},
		),
	}
	return store.MustNewService(append(defaultOpts, opts...)...)
}
