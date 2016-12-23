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

	"github.com/corestoreio/csfw/config/cfgpath"
	"github.com/corestoreio/csfw/config/element"
	"github.com/corestoreio/errors"
	"github.com/stretchr/testify/assert"
)

var _ element.FieldTyper = (*element.FieldType)(nil)

func TestFieldRouteHash(t *testing.T) {

	tests := []struct {
		preRoutes  []cfgpath.Route
		field      *element.Field
		wantHash   uint64
		wantErrBhf errors.BehaviourFunc
	}{
		{[]cfgpath.Route{cfgpath.NewRoute("aa"), cfgpath.NewRoute("b")}, &element.Field{ID: cfgpath.NewRoute("ca")}, 5676413504385759347, nil},
		{[]cfgpath.Route{cfgpath.NewRoute("aa"), cfgpath.NewRoute("b")}, &element.Field{ID: cfgpath.NewRoute("cb")}, 5676414603897387558, nil},
		{nil, &element.Field{ID: cfgpath.NewRoute("cb")}, 622143294520562096, nil},
		{nil, &element.Field{}, 0, errors.IsEmpty},
		{[]cfgpath.Route{{}, {}}, &element.Field{ID: cfgpath.NewRoute("ca")}, 622146593055446729, nil},
	}
	for i, test := range tests {
		haveHash, haveErr := test.field.RouteHash(test.preRoutes...)
		if test.wantErrBhf != nil {
			assert.Empty(t, haveHash, "Index %d", i)
			assert.True(t, test.wantErrBhf(haveErr), "Index %d => %s", i, haveErr)
			continue
		}
		assert.NoError(t, haveErr, "Index %d", i)
		assert.Exactly(t, test.wantHash, haveHash, "Want: %d Have: %d => Index %d", test.wantHash, haveHash, i)
	}
}

func TestFieldRoute(t *testing.T) {

	tests := []struct {
		preRoutes  []cfgpath.Route
		field      *element.Field
		wantR      string
		wantErrBhf errors.BehaviourFunc
	}{
		{[]cfgpath.Route{cfgpath.NewRoute("aa"), cfgpath.NewRoute("b")}, &element.Field{ID: cfgpath.NewRoute("ca")}, "aa/b/ca", errors.IsNotValid},
		{[]cfgpath.Route{cfgpath.NewRoute("aa"), cfgpath.NewRoute("bb")}, &element.Field{ID: cfgpath.NewRoute("ca")}, "aa/bb/ca", nil},
		{nil, &element.Field{ID: cfgpath.NewRoute("cb")}, "cb", errors.IsNotValid},
		{nil, &element.Field{}, "", errors.IsEmpty},
		{[]cfgpath.Route{{}, {}}, &element.Field{ID: cfgpath.NewRoute("ca")}, "", errors.IsNotValid},
	}
	for i, test := range tests {
		haveR, haveErr := test.field.Route(test.preRoutes...)
		if test.wantErrBhf != nil {
			assert.Exactly(t, cfgpath.Route{}, haveR, "Index %d", i)
			assert.True(t, test.wantErrBhf(haveErr), "Index %d = %s", i, haveErr)
			continue
		}
		assert.NoError(t, haveErr, "Index %d", i)
		assert.Exactly(t, test.wantR, haveR.String(), "Index %d", i)
	}
}
