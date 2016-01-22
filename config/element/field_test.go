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

package element_test

import (
	"testing"

	"github.com/corestoreio/csfw/config/element"
	"github.com/corestoreio/csfw/config/path"
	"github.com/stretchr/testify/assert"
)

var _ error = (*element.FieldError)(nil)
var _ element.FieldTyper = (*element.FieldType)(nil)

func TestFieldRouteHash(t *testing.T) {
	t.Parallel()
	tests := []struct {
		preRoutes []path.Route
		field     *element.Field
		wantHash  uint64
		wantErr   error
	}{
		{[]path.Route{path.NewRoute("aa"), path.NewRoute("b")}, &element.Field{ID: path.NewRoute("ca")}, 5676413504385759347, nil},
		{[]path.Route{path.NewRoute("aa"), path.NewRoute("b")}, &element.Field{ID: path.NewRoute("cb")}, 5676414603897387558, nil},
		{nil, &element.Field{ID: path.NewRoute("cb")}, 622143294520562096, nil},
		{nil, &element.Field{}, 0, path.ErrRouteEmpty},
		{[]path.Route{{}, {}}, &element.Field{ID: path.NewRoute("ca")}, 622146593055446729, nil},
	}
	for i, test := range tests {
		haveHash, haveErr := test.field.RouteHash(test.preRoutes...)
		if test.wantErr != nil {
			assert.Empty(t, haveHash, "Index %d", i)
			assert.EqualError(t, haveErr, test.wantErr.Error(), "Index %d", i)
			continue
		}
		assert.NoError(t, haveErr, "Index %d", i)
		assert.Exactly(t, test.wantHash, haveHash, "Want: %d Have: %d => Index %d", test.wantHash, haveHash, i)
	}
}

func TestFieldFQPathDefault(t *testing.T) {
	t.Parallel()
	tests := []struct {
		preRoutes []path.Route
		field     *element.Field
		wantFQ    string
		wantErr   error
	}{
		{[]path.Route{path.NewRoute("aa"), path.NewRoute("b")}, &element.Field{ID: path.NewRoute("ca")}, "aa/b/ca", path.ErrIncorrectPath},
		{[]path.Route{path.NewRoute("aa"), path.NewRoute("bb")}, &element.Field{ID: path.NewRoute("ca")}, "default/0/aa/bb/ca", nil},
		{nil, &element.Field{ID: path.NewRoute("cb")}, "cb", path.ErrIncorrectPath},
		{nil, &element.Field{}, "", path.ErrRouteEmpty},
		{[]path.Route{{}, {}}, &element.Field{ID: path.NewRoute("ca")}, "", path.ErrIncorrectPath},
	}
	for i, test := range tests {
		haveFQ, haveErr := test.field.FQPathDefault(test.preRoutes...)
		if test.wantErr != nil {
			assert.Empty(t, haveFQ, "Index %d", i)
			assert.EqualError(t, haveErr, test.wantErr.Error(), "Index %d", i)
			continue
		}
		assert.NoError(t, haveErr, "Index %d", i)
		assert.Exactly(t, test.wantFQ, haveFQ, "Index %d", i)
	}
}
