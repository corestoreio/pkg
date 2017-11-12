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

package backend_test

import (
	"testing"

	"github.com/corestoreio/pkg/backend"
	"github.com/corestoreio/pkg/config"
	"github.com/corestoreio/pkg/config/cfgmock"
	"github.com/corestoreio/pkg/config/cfgmodel"
	"github.com/corestoreio/pkg/store/scope"
	"github.com/stretchr/testify/assert"
)

func TestConfigRedirectToBase(t *testing.T) {
	t.Parallel()

	r := backend.NewConfigRedirectToBase(
		backend.Backend.WebURLRedirectToBase.String(),
		cfgmodel.WithFieldFromSectionSlice(backend.ConfigStructure),
	)

	redirCode, err := r.Get(cfgmock.NewService().NewScoped(0, 0))
	if err != nil {
		t.Fatal(err)
	}
	assert.Exactly(
		t,
		1, // default value in backend.ConfigStructure
		redirCode,
	)

	redirCode, err = r.Get(cfgmock.NewService().NewScoped(10, 13))
	if err != nil {
		t.Fatal(err)
	}
	// 1 == default value in backend.ConfigStructure
	assert.Exactly(t, 1, redirCode)

	webURLRedirectToBasePath, err := backend.Backend.WebURLRedirectToBase.ToPath(0, 0)
	if err != nil {
		t.Fatal(err)
	}

	cr := cfgmock.NewService(
		cfgmock.WithPV(cfgmock.PathValue{
			webURLRedirectToBasePath.String():                       2,
			webURLRedirectToBasePath.Bind(scope.Store, 33).String(): 34,
		}),
	)

	tests := []struct {
		sg   config.Scoped
		want int
	}{
		{cr.NewScoped(0, 0), 2},
		{cr.NewScoped(1, 2), 2},
		{cr.NewScoped(1, 33), 34},
	}
	for i, test := range tests {
		code, err := r.Get(test.sg)
		if err != nil {
			t.Fatalf("Index %d => %s", i, err)
		}
		assert.Exactly(t, test.want, code, "Index %d", i)
		assert.False(t, r.HasErrors(), "Index %d", i)
	}

	mw := new(cfgmock.Write)
	assert.EqualError(t, r.Write(mw, 200, scope.Default, 0),
		"The value '200' cannot be found within the allowed Options():\n[{\"Value\":0,\"Label\":\"No\"},{\"Value\":1,\"Label\":\"Yes (302 Found)\"},{\"Value\":302,\"Label\":\"Yes (302 Found)\"},{\"Value\":301,\"Label\":\"Yes (301 Moved Permanently)\"}]\n\nJSON Error: %!s(<nil>)",
	) // 200 not allowed
}

func BenchmarkConfigRedirectToBase(b *testing.B) {
	r := backend.NewConfigRedirectToBase(
		backend.Backend.WebURLRedirectToBase.String(),
		cfgmodel.WithFieldFromSectionSlice(backend.ConfigStructure),
	)
	webURLRedirectToBasePath, err := backend.Backend.WebURLRedirectToBase.ToPath(0, 0)
	if err != nil {
		b.Fatal(err)
	}

	sg := cfgmock.NewService(
		cfgmock.WithPV(cfgmock.PathValue{
			webURLRedirectToBasePath.String():                         2,
			webURLRedirectToBasePath.Bind(scope.Website, 33).String(): 34,
		}),
	).NewScoped(33, 1)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		code, err := r.Get(sg)
		if err != nil {
			b.Fatal(err)
		}
		if code != 34 {
			b.Fatalf("Want %d Have %d", 34, code)
		}
	}
}
