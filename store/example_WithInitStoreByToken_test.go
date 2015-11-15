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
	"bytes"
	"fmt"
	"io/ioutil"
	std "log"
	"net/http"
	"net/http/httptest"

	"github.com/corestoreio/csfw/config/scope"
	"github.com/corestoreio/csfw/net/ctxhttp"
	"github.com/corestoreio/csfw/net/ctxjwt"
	"github.com/corestoreio/csfw/storage/dbr"
	"github.com/corestoreio/csfw/store"
	"github.com/corestoreio/csfw/utils/log"
	"golang.org/x/net/context"
)

var testDebugLogBuf bytes.Buffer
var testStoreService store.Reader

func initStore() {
	store.PkgLog = log.NewStdLogger(
		log.SetStdDebug(&testDebugLogBuf, "testErr: ", std.Lshortfile),
	)
	store.PkgLog.SetLevel(log.StdLevelDebug)

	testStoreService = store.MustNewService(
		scope.Option{
			Website: scope.MockID(1), // bound to website ID 1 = Europe
		},
		store.NewStorage(
			// Storage gets usually loaded from the database tables containing
			// website, group and store. For the sake of this example the storage
			// is hard coded.
			store.SetStorageWebsites(
				&store.TableWebsite{WebsiteID: 0, Code: dbr.NewNullString("admin"), Name: dbr.NewNullString("Admin"), SortOrder: 0, DefaultGroupID: 0, IsDefault: dbr.NewNullBool(false)},
				&store.TableWebsite{WebsiteID: 1, Code: dbr.NewNullString("euro"), Name: dbr.NewNullString("Europe"), SortOrder: 0, DefaultGroupID: 1, IsDefault: dbr.NewNullBool(true)},
				&store.TableWebsite{WebsiteID: 2, Code: dbr.NewNullString("oz"), Name: dbr.NewNullString("OZ"), SortOrder: 20, DefaultGroupID: 3, IsDefault: dbr.NewNullBool(false)},
			),
			store.SetStorageGroups(
				&store.TableGroup{GroupID: 3, WebsiteID: 2, Name: "Australia", RootCategoryID: 2, DefaultStoreID: 5},
				&store.TableGroup{GroupID: 1, WebsiteID: 1, Name: "DACH Group", RootCategoryID: 2, DefaultStoreID: 2},
				&store.TableGroup{GroupID: 0, WebsiteID: 0, Name: "Default", RootCategoryID: 0, DefaultStoreID: 0},
				&store.TableGroup{GroupID: 2, WebsiteID: 1, Name: "UK Group", RootCategoryID: 2, DefaultStoreID: 4},
			),
			store.SetStorageStores(
				&store.TableStore{StoreID: 0, Code: dbr.NewNullString("admin"), WebsiteID: 0, GroupID: 0, Name: "Admin", SortOrder: 0, IsActive: true},
				&store.TableStore{StoreID: 5, Code: dbr.NewNullString("au"), WebsiteID: 2, GroupID: 3, Name: "Australia", SortOrder: 10, IsActive: true},
				&store.TableStore{StoreID: 1, Code: dbr.NewNullString("de"), WebsiteID: 1, GroupID: 1, Name: "Germany", SortOrder: 10, IsActive: true},
				&store.TableStore{StoreID: 4, Code: dbr.NewNullString("uk"), WebsiteID: 1, GroupID: 2, Name: "UK", SortOrder: 10, IsActive: true},
				&store.TableStore{StoreID: 2, Code: dbr.NewNullString("at"), WebsiteID: 1, GroupID: 1, Name: "Österreich", SortOrder: 20, IsActive: true},
				&store.TableStore{StoreID: 6, Code: dbr.NewNullString("nz"), WebsiteID: 2, GroupID: 3, Name: "Kiwi", SortOrder: 30, IsActive: true},
				&store.TableStore{IsActive: false, StoreID: 3, Code: dbr.NewNullString("ch"), WebsiteID: 1, GroupID: 1, Name: "Schweiz", SortOrder: 30},
			),
		),
	)
}

func ExampleWithInitStoreByToken() {
	initStore()
	ctx := store.NewContextReader(context.Background(), testStoreService)

	jwtService, err := ctxjwt.NewService(ctxjwt.WithPassword([]byte(`GÒph3r`)))

	finalHandler := ctxhttp.Chain(
		ctxhttp.HandlerFunc(func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			_, haveReqStore, err := store.FromContextReader(ctx)
			if err != nil {
				return err
			}
			// now we know that the current request depends on the store view DE.
			fmt.Fprintf(w, "StoreCode: %s\n", haveReqStore.StoreCode())
			return nil
		}), // 								  executed 3rd
		store.WithInitStoreByToken(),      // executed 2nd
		jwtService.WithParseAndValidate(), // executed 1st
	)

	ts := httptest.NewServer(ctxhttp.NewAdapter(ctx, finalHandler))
	defer ts.Close()

	// Setup GET request
	token, _, err := jwtService.GenerateToken(
		map[string]interface{}{
			// Despite default store for Website ID 1 is AT we are currently
			// in the store context of DE.
			store.ParamName: "de",
		},
	)
	if err != nil {
		log.Fatal("jwtService.GenerateToken", "err", err)
	}

	req, err := http.NewRequest("GET", ts.URL, nil)
	if err != nil {
		log.Fatal("http.Get", "err", err)
	}
	ctxjwt.SetHeaderAuthorization(req, token)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal("http.DefaultClient.Do", "err", err)
	}

	response, err := ioutil.ReadAll(res.Body)
	if errC := res.Body.Close(); errC != nil {
		log.Fatal("res.Body.Close", "err", errC)
	}

	if err != nil {
		log.Fatal("ioutil.ReadAll", "err", err)
	}

	fmt.Printf("Response: %s\n", response)
	fmt.Printf("Log: %s\n", testDebugLogBuf.String())
	// Output:
	// Response: StoreCode: de
	//
	// Log:
}
