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

package text_test

import (
	"database/sql"
	"database/sql/driver"
	"encoding"
	"encoding/json"
	"testing"

	"errors"
	"github.com/corestoreio/csfw/storage/text"
	"github.com/stretchr/testify/assert"
)

// These checks if a type implements an interface belong into the test package
// and not into its "main" package. Otherwise you would also compile each time
// al the package with their interfaces.
var _ encoding.TextMarshaler = (*text.Long)(nil)
var _ encoding.TextUnmarshaler = (*text.Long)(nil)
var _ sql.Scanner = (*text.Long)(nil)
var _ driver.Valuer = (*text.Long)(nil)

func TestEqual(t *testing.T) {
	tests := []struct {
		a    text.Long
		b    text.Long
		want bool
	}{
		{nil, nil, true},
		{text.Long("a"), text.Long("a"), true},
		{text.Long("a"), text.Long("b"), false},
		{text.Long("a\x80"), text.Long("a"), false},
	}
	for i, test := range tests {
		assert.Exactly(t, test.want, test.a.Equal(test.b), "Index %d", i)
	}
}

func TestLong(t *testing.T) {
	t.Parallel()
	const have string = `Hello fellow Gpher's`
	l := text.Long(have)
	var l1 text.Long
	assert.True(t, l1.IsEmpty())
	assert.False(t, l.IsEmpty())
	assert.Exactly(t, have, l.String())

	l2 := l.Copy()
	assert.Exactly(t, l, l2)
	l2 = nil
	assert.True(t, l2.IsEmpty())
	assert.False(t, l.IsEmpty())
}

func TestTextMarshal(t *testing.T) {
	const have = `admin/security/passwrd_lifetime`
	t.Parallel()
	r := text.Long(have)
	j, err := json.Marshal(r)
	assert.NoError(t, err)
	assert.Exactly(t, `"`+have+`"`, string(j))
}

func TestUnmarshalTextOk(t *testing.T) {
	t.Parallel()
	const have = `admin/security/passwörd_lif‹time`
	var r text.Long
	err := json.Unmarshal([]byte(`"`+have+`"`), &r)
	assert.NoError(t, err)
	assert.Exactly(t, have, r.String())
}

func TestScan(t *testing.T) {
	tests := []struct {
		want    string
		val     interface{}
		wantErr error
	}{
		{"", nil, nil},
		{"hello", "hello", nil},
		{"h€llo", []byte("h€llo"), nil},
		{"", 8, errors.New("Cannot convert value 8 to []byte")},
	}
	for i, test := range tests {
		var l text.Long
		haveErr := l.Scan(test.val)
		if test.wantErr != nil {
			assert.EqualError(t, haveErr, test.wantErr.Error(), "Index %d", i)
			assert.Nil(t, l, "Index %d", i)
			continue
		}
		assert.NoError(t, haveErr, "Index %d", i)
		assert.Exactly(t, test.want, l.String(), "Index %d", i)
	}
}

func TestValue(t *testing.T) {
	t.Parallel()
	l1 := text.Long(`x`)
	v, err := l1.Value()
	assert.NoError(t, err)
	assert.NotNil(t, v)
	assert.Exactly(t, l1, v.(text.Long))

	var l2 text.Long
	v, err = l2.Value()
	assert.NoError(t, err)
	assert.Nil(t, v)

}
