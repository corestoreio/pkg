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

package path_test

import (
	"encoding/json"
	"github.com/corestoreio/csfw/config/path"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRouteAppend(t *testing.T) {
	t.Parallel()
	tests := []struct {
		a       path.Route
		b       path.Route
		want    string
		wantErr error
	}{
		{path.Route("aa"), path.Route("bb/cc"), "aa/bb/cc", nil},
		{path.Route("aa"), path.Route("bbcc"), "aa/bbcc", nil},
		{path.Route("aa"), path.Route("bb\x80cc"), "", path.ErrRouteInvalidBytes},
	}
	for i, test := range tests {
		haveErr := test.a.Append(test.b)
		if test.wantErr != nil {
			assert.EqualError(t, haveErr, test.wantErr.Error(), "Index %d", i)
			continue
		}
		assert.NoError(t, haveErr, "Index %d", i)
		assert.Exactly(t, test.want, test.a.String(), "Index %d", i)
	}
}

func TestRouteTextMarshal(t *testing.T) {
	r := path.Route("admin/security/password_lifetime")
	j, err := json.Marshal(r)
	assert.NoError(t, err)
	assert.Exactly(t, "\"admin/security/password_lifetime\"", string(j))
}

func TestRouteUnmarshalTextOk(t *testing.T) {
	var r path.Route
	err := json.Unmarshal([]byte(`"admin/security/password_lifetime"`), &r)
	assert.NoError(t, err)
	assert.Exactly(t, "admin/security/password_lifetime", r.String())
}
