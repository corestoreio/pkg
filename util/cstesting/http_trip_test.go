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

package cstesting_test

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"sync"
	"testing"

	"github.com/corestoreio/csfw/util/cstesting"
	"github.com/corestoreio/csfw/util/errors"
	"github.com/stretchr/testify/assert"
)

var _ http.RoundTripper = (*cstesting.HTTPTrip)(nil)

func TestNewHttpTrip_Ok(t *testing.T) {
	t.Parallel()
	tr := cstesting.NewHTTPTrip(333, "Hello Wørld", nil)
	cl := &http.Client{
		Transport: tr,
	}
	const reqURL = `http://corestore.io`
	var wg sync.WaitGroup
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			getReq, err := http.NewRequest("GET", reqURL, nil)
			if err != nil {
				t.Fatal(err)
			}
			resp, err := cl.Do(getReq)
			if err != nil {
				t.Fatal(err)
			}

			defer func() {
				if err := resp.Body.Close(); err != nil {
					t.Fatal(err)
				}
			}()
			data, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Fatal(err)
			}
			assert.Exactly(t, "Hello Wørld", string(data))
			assert.Exactly(t, 333, resp.StatusCode)
		}(&wg)
	}
	wg.Wait()

	tr.RequestsMatchAll(t, func(r *http.Request) bool {
		return r.URL.String() == reqURL
	})
	tr.RequestsCount(t, 10)
}

func TestNewHttpTrip_Error(t *testing.T) {
	t.Parallel()
	tr := cstesting.NewHTTPTrip(501, "Hello Error", errors.NewNotValidf("test not valid"))
	cl := &http.Client{
		Transport: tr,
	}

	var wg sync.WaitGroup
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			getReq, err := http.NewRequest("GET", "http://noophole.com", nil)
			if err != nil {
				t.Fatal("NewRequest", err)
			}
			resp, err := cl.Do(getReq)
			assert.True(t, errors.IsNotValid(err.(*url.Error).Err), "ErrorDo: %#v", err)
			assert.Nil(t, resp)
		}(&wg)
	}
	wg.Wait()
	tr.RequestsCount(t, 10)
}

func TestNewHttpTrip_Error_FromFile(t *testing.T) {
	t.Parallel()
	tr := cstesting.NewHTTPTripFromFile(505, "file_notFOUND.txt")
	cl := &http.Client{
		Transport: tr,
	}

	getReq, err := http.NewRequest("GET", "http://noophole.com", nil)
	if err != nil {
		t.Fatal("NewRequest", err)
	}
	resp, err := cl.Do(getReq)
	assert.True(t, errors.IsNotFound(err.(*url.Error).Err), "ErrorDo: %#v", err)
	assert.Nil(t, resp)

	tr.RequestsCount(t, 1)
}
