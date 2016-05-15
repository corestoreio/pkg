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

package mwcors

//import (
//	"testing"
//
//	"github.com/corestoreio/csfw/config/cfgmock"
//
//	"context"
//	"github.com/corestoreio/csfw/store"
//	"github.com/corestoreio/csfw/store/scope"
//	"github.com/corestoreio/csfw/store/storemock"
//	"github.com/stretchr/testify/assert"
//)
//
//func TestCorsCurrent_ShouldCreateANewScopedBasedCors(t *testing.T) {
//
//	be := initBackend(t)
//
//	cfgGet := cfgmock.NewService(
//		cfgmock.WithPV(cfgmock.PathValue{
//			mustToPath(t, be.NetCorsExposedHeaders.FQ, scope.Website, 2):     "X-CoreStore-ID\nContent-Type\n\n",
//			mustToPath(t, be.NetCorsAllowedOrigins.FQ, scope.Website, 2):     "host1.com\nhost2.com\n\n",
//			mustToPath(t, be.NetCorsAllowedMethods.FQ, scope.Website, 2):     "PATCH\nDELETE",
//			mustToPath(t, be.NetCorsAllowedHeaders.FQ, scope.Website, 2):     "Date,X-Header1",
//			mustToPath(t, be.NetCorsAllowCredentials.FQ, scope.Website, 2):   "1",
//			mustToPath(t, be.NetCorsOptionsPassthrough.FQ, scope.Website, 2): "1",
//			mustToPath(t, be.NetCorsMaxAge.FQ, scope.Website, 2):             "2h",
//		}),
//	)
//
//	scpO, err := scope.SetByCode(scope.Website, "oz")
//	if err != nil {
//		t.Fatal(err)
//	}
//	storeSrv := storemock.NewEurozzyService(scpO, store.WithStorageConfig(cfgGet))
//	dftStore, err := storeSrv.Store() // default store for AU
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	if dftStore.Website.Config == nil {
//		t.Fatalf("Website Config unexpected nil: %#v", dftStore.Website)
//	}
//	ctx := store.WithContextProvider(context.Background(), storeSrv, dftStore)
//
//	c := MustNew(WithBackendApplied(be, dftStore.Website.Config)) // OZ website ID = 2 and AU store ID = 5
//
//	csc := newScopeCache(c)
//
//	scopedCors, err := c.current(csc, ctx)
//	if err != nil {
//		t.Fatal(err)
//	}
//	assert.Exactly(t, scope.NewHash(scope.Website, 2), scopedCors.scopedTo)
//
//	// check that we get the same cors back
//	scopedCors2, err := c.current(csc, ctx)
//	if err != nil {
//		t.Fatal(err)
//	}
//	assert.Exactly(t, scope.NewHash(scope.Website, 2), scopedCors.scopedTo)
//	assert.Exactly(t, scopedCors2, scopedCors)
//}
