// Copyright 2015 CoreGroup Authors
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

//func TestGroup(t *testing.T) {
//	s1 := storeManager.Group().Collection()
//	assert.True(t, len(s1) > 2, "1. There should be at least three groups in the slice")
//	assert.True(t, s1.Len() > 2, "2. There should be at least three groups in the slice")
//
//	for _, group := range storeManager.Group().Collection() {
//		assert.NotNil(t, group, "Expecting first index to be nil")
//		assert.True(t, len(group.Name) > 1, "group.Name should be longer than 1 char: %#v", group)
//	}
//	assert.Equal(t, utils.Int64Slice{4, 0, 1}, s1.IDs())
//}
//
//func TestGetGroupByID(t *testing.T) {
//	g, err := storeManager.Group().Get(1)
//	if err != nil {
//		t.Error(err)
//		assert.Nil(t, g)
//	} else {
//		assert.NoError(t, err)
//		assert.Equal(t, "Madison Island", g.Name)
//	}
//	gInvalid, err := storeManager.Group().Get(10000)
//	assert.EqualError(t, err, store.ErrGroupNotFound.Error())
//	assert.Nil(t, gInvalid)
//}
//
//func TestGroupGetStores(t *testing.T) {
//	sInvalid, err := storeManager.Group().Stores(321)
//	assert.EqualError(t, err, store.ErrGroupStoresNotFound.Error())
//	assert.Nil(t, sInvalid)
//
//	stores, err := storeManager.Group().Stores(1)
//	assert.NoError(t, err)
//	assert.Equal(t, "default,french,german", stores.Codes().Join(","))
//}
//
//var benchGroupGetStoreCodes utils.StringSlice
//
//// BenchmarkGroupGetStoreCodes	10000000	       177 ns/op	      48 B/op	       1 allocs/op
//func BenchmarkGroupGetStoreCodes(b *testing.B) {
//	b.ReportAllocs()
//	for i := 0; i < b.N; i++ {
//		stores, err := storeManager.Group().Stores(1)
//		if err != nil {
//			b.Error(err)
//		}
//		benchGroupGetStoreCodes = stores.Codes()
//	}
//}
//
//func TestGroupGetWebsite(t *testing.T) {
//
//	wsInvalid, err := storeManager.Group().Website(321)
//	assert.EqualError(t, store.ErrGroupWebsiteNotFound, err.Error())
//	assert.Nil(t, wsInvalid)
//
//	website, err := storeManager.Group().Website(1)
//	assert.NoError(t, err)
//	assert.Equal(t, "Main Website", website.Name.String)
//}
//
//func TestGroupDefaultStore(t *testing.T) {
//
//	wsInvalid, err := storeManager.Group().DefaultStore(321)
//	assert.EqualError(t, store.ErrGroupNotFound, err.Error())
//	assert.Nil(t, wsInvalid)
//
//	store, err := storeManager.Group().DefaultStore(1)
//	assert.NoError(t, err)
//	assert.Equal(t, "default", store.Code.String)
//}
//
//var benchGroupDefaultStore *store.TableStore
//
//// BenchmarkGroupDefaultStore	50000000	        30.2 ns/op	       0 B/op	       0 allocs/op
//func BenchmarkGroupDefaultStore(b *testing.B) {
//	b.ReportAllocs()
//	for i := 0; i < b.N; i++ {
//		store, err := storeManager.Group().DefaultStore(1)
//		if err != nil {
//			b.Error(err)
//		}
//		benchGroupDefaultStore = store
//	}
//}
