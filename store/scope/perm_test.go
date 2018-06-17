// Copyright 2015-present, Cyrill @ Schumacher.fm and the CoreStore contributors
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

package scope_test

import (
	"encoding/json"
	"testing"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/pkg/store/scope"
	"github.com/stretchr/testify/assert"
)

var (
	_ json.Unmarshaler = (*scope.Perm)(nil)
	_ json.Marshaler   = (*scope.Perm)(nil)
)

func TestPermAll(t *testing.T) {
	var p scope.Perm
	pa := p.All()
	assert.True(t, pa.Has(scope.Default))
	assert.True(t, pa.Has(scope.Website))
	assert.True(t, pa.Has(scope.Store))
}

func TestPermTop(t *testing.T) {
	assert.Exactly(t, scope.Website, scope.PermWebsite.Top())
	assert.Exactly(t, scope.Store, scope.PermStore.Top())
	assert.Exactly(t, scope.Default, scope.PermDefault.Top())
	assert.Exactly(t, scope.Website, scope.Perm(44).Top())
	assert.Exactly(t, scope.Store, scope.PermWebsiteReverse.Top())
	assert.Exactly(t, scope.Store, scope.PermStoreReverse.Top())
}

func TestPerm_JSON(t *testing.T) {

	t.Run("Marshal", func(t *testing.T) {
		var p scope.Perm
		pa := p.All()
		jd, err := pa.MarshalJSON()
		assert.NoError(t, err)
		assert.Exactly(t, "\"stores\"", string(jd))
	})
	t.Run("null Marshal", func(t *testing.T) {
		var p scope.Perm
		jd, err := p.MarshalJSON()
		assert.NoError(t, err)
		assert.Exactly(t, "null", string(jd))
	})
	t.Run("unmarshal", func(t *testing.T) {
		data := []byte(`"websites"`)
		var p scope.Perm
		assert.NoError(t, p.UnmarshalJSON(data))
		assert.Exactly(t, scope.PermWebsite, p)
	})
	t.Run("error unmarshal", func(t *testing.T) {
		data := []byte(`Websites`)
		var p scope.Perm
		err := p.UnmarshalJSON(data)
		assert.True(t, errors.NotSupported.Match(err), "%+v", err)
		assert.Exactly(t, scope.Perm(0), p)
	})
}

func TestMakePerm(t *testing.T) {
	tests := []struct {
		name     string
		wantPerm scope.Perm
		wantErr  errors.Kind
	}{
		{"default", scope.PermDefault, errors.NoKind},
		{"d", scope.PermDefault, errors.NoKind},
		{"", scope.PermDefault, errors.NoKind},

		{"websites", scope.PermWebsite, errors.NoKind},
		{"w", scope.PermWebsite, errors.NoKind},

		{"stores", scope.PermStore, errors.NoKind},
		{"s", scope.PermStore, errors.NoKind},

		{"g", 0, errors.NotSupported},
	}
	for i, test := range tests {
		haveP, haveErr := scope.MakePerm(test.name)
		if test.wantErr > 0 {
			assert.True(t, test.wantErr.Match(haveErr), "Index %d %+v", i, haveErr)
		} else {
			assert.NoError(t, haveErr)
			assert.Exactly(t, test.wantPerm, haveP, "Index %d", i)
		}
	}
}

func TestPerm(t *testing.T) {

	const (
		scope1 scope.Type = iota + 1
		scope2
		scope3
		scope4
		scope5
	)

	tests := []struct {
		have    []scope.Type
		want    scope.Type
		notWant scope.Type
		human   []string
		string
	}{
		{[]scope.Type{scope1}, scope1, scope3, []string{"Default"}, "default"},
		{[]scope.Type{scope1, scope2}, scope2, scope3, []string{"Default", "Website"}, "websites"},
		{[]scope.Type{scope3, scope4}, scope3, scope2, []string{"Group", "Store"}, "stores"},
		{[]scope.Type{scope4, scope5}, scope4, scope2, []string{"Store"}, "stores"},
	}

	for _, test := range tests {
		var b = scope.Perm(0).Set(test.have...)
		if !b.Has(test.want) {
			t.Errorf("%d should contain %d", b, test.want)
		}
		if b.Has(test.notWant) {
			t.Errorf("%d should not contain %d", b, test.notWant)
		}
		assert.EqualValues(t, test.human, b.Human())
		assert.EqualValues(t, test.string, b.String())
	}
}
