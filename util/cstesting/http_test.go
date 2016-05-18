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
	"github.com/corestoreio/csfw/util/cstesting"
	"github.com/corestoreio/csfw/util/errors"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

var _ http.RoundTripper = (*cstesting.HttpTrip)(nil)

func TestNewHttpTrip(t *testing.T) {

	tr := cstesting.NewHttpTrip(333, "Hello Wørld", errors.NewNotValidf("test not valid"))
	req := httptest.NewRequest("GET", "http://noophole.com", nil)
	resp, err := tr.RoundTrip(req)
	assert.True(t, errors.IsNotValid(err), "Error: %s")
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
}
