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

package transcache_test

import (
	"os"
	"sync"
	"testing"

	"fmt"

	"github.com/alicebob/miniredis"
	"github.com/corestoreio/csfw/storage/transcache"
	"github.com/corestoreio/csfw/storage/transcache/tcbigcache"
	"github.com/corestoreio/csfw/storage/transcache/tcboltdb"
	"github.com/corestoreio/csfw/storage/transcache/tcredis"
	"github.com/corestoreio/csfw/util/cstesting"
)

// run this with go test -race .

func TestProcessor_Parallel_GetSet_BigCache(t *testing.T) {
	newTestNewProcessor(t, tcbigcache.With())
}

func TestProcessor_Parallel_GetSet_Bolt(t *testing.T) {
	f := getTempFile(t)
	defer os.Remove(f)
	newTestNewProcessor(t, tcboltdb.WithFile(f, 0600))
}

func TestProcessor_Parallel_GetSet_Redis(t *testing.T) {
	mr := miniredis.NewMiniRedis()
	if err := mr.Start(); err != nil {
		t.Fatal(err)
	}
	defer mr.Close()
	redConURL := fmt.Sprintf("redis://%s/2", mr.Addr())
	newTestNewProcessor(t, tcredis.WithURL(redConURL, nil))
}

func newTestNewProcessor(t *testing.T, opts ...transcache.Option) {
	p, err := transcache.NewProcessor(append(opts, transcache.WithGob())...)
	cstesting.FatalIfError(t, err)

	var wg sync.WaitGroup

	// to detect race conditions run with -race
	wg.Add(1)
	go testCountry(t, &wg, p, []byte("country_one"))

	wg.Add(1)
	go testStoreSlice(t, &wg, p, []byte("stores_one"))

	wg.Add(1)
	go testCountry(t, &wg, p, []byte("country_two"))

	wg.Add(1)
	go testStoreSlice(t, &wg, p, []byte("stores_two"))

	wg.Add(1)
	go testStoreSlice(t, &wg, p, []byte("stores_three"))

	wg.Add(1)
	go testCountry(t, &wg, p, []byte("country_three"))

	wg.Wait()
}
