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

package httputil_test

import (
	"testing"

	"github.com/corestoreio/csfw/net/httputil"
	"github.com/stretchr/testify/assert"
)

func TestVersionize(t *testing.T) {
	tests := []struct {
		have, want string
	}{
		{"login", "/V1/login"},
		{"/login", "/V1/login"},
		{"", "/V1/"},
	}
	for _, test := range tests {
		h := httputil.APIRoute.Versionize(test.have)
		assert.Equal(t, test.want, h)
	}
	assert.Equal(t, "/V1/", httputil.APIRoute.String())
}
