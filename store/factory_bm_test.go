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

package store

import (
	"testing"

	"github.com/corestoreio/csfw/config"
)

func mustNewFactory(cfg config.Getter, opts ...Option) *factory {
	f, err := newFactory(cfg, opts...)
	if err != nil {
		panic(err)
	}
	return f
}

var benchmarkFactoryWebsite Website
var benchmarkFactoryWebsiteDefaultGroup Group

func BenchmarkFactoryWebsiteGetDefaultGroup(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		var err error
		benchmarkFactoryWebsite, err = testFactory.Website(1)
		if err != nil {
			b.Error(err)
		}

		benchmarkFactoryWebsiteDefaultGroup, err = benchmarkFactoryWebsite.DefaultGroup()
		if err != nil {
			b.Error(err)
		}
	}
}

var benchmarkFactoryGroup Group
var benchmarkFactoryGroupDefaultStore Store

func BenchmarkFactoryGroupGetDefaultStore(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		var err error
		benchmarkFactoryGroup, err = testFactory.Group(3)
		if err != nil {
			b.Error(err)
		}

		benchmarkFactoryGroupDefaultStore, err = benchmarkFactoryGroup.DefaultStore()
		if err != nil {
			b.Error(err)
		}
	}
}

var benchmarkFactoryStore Store
var benchmarkFactoryStoreWebsite Website

func BenchmarkFactoryStoreGetWebsite(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		var err error
		benchmarkFactoryStore, err = testFactory.Store(1)
		if err != nil {
			b.Error(err)
		}

		benchmarkFactoryStoreWebsite = benchmarkFactoryStore.Website
		if err := benchmarkFactoryStoreWebsite.Validate(); err != nil {
			b.Errorf("benchmarkFactoryStoreWebsite %+v", err)
		}
	}
}

var benchmarkFactoryStoreID int64

func BenchmarkFactoryDefaultStoreView(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		var err error
		benchmarkFactoryStoreID, err = testFactory.DefaultStoreID()
		if err != nil {
			b.Error(err)
		}
	}
}
