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

package ctxhttp_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/net/ctxhttp"
	"github.com/corestoreio/csfw/net/httputils"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

func TestWithValidateBaseUrlRedirect(t *testing.T) {

	cr := config.NewMockReader(
		config.MockInt(func(path string) (int, error) {
			t.Logf("GetInt %s", path)
			return 0, nil
		}),
		config.MockBool(func(path string) (bool, error) {
			t.Logf("GetBool %s", path)
			return true, nil
		}),
	)

	mw := ctxhttp.WithValidateBaseUrl(cr)
	w := httptest.NewRecorder()
	req, err := http.NewRequest(httputils.MethodGet, "http://corestore.io/catalog/product/view", nil)
	assert.NoError(t, err)
	ctx := context.Background()

	err = mw(nil).ServeHTTPContext(ctx, w, req)
	assert.EqualError(t, err, "Cannot extract config.Reader from context")

	t.Log("@todo proper testing")
}
