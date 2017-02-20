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

package problem_test

import (
	"encoding/json"
	"net/http"
	"net/mail"
	"net/textproto"
	"net/url"
	"sort"
	"testing"

	"github.com/corestoreio/csfw/net/problem"
	"github.com/corestoreio/errors"
	"github.com/mailru/easyjson"
	"github.com/stretchr/testify/assert"
)

var _ json.Marshaler = (*problem.Detail)(nil)
var _ json.Unmarshaler = (*problem.Detail)(nil)
var _ easyjson.Marshaler = (*problem.Detail)(nil)
var _ easyjson.Unmarshaler = (*problem.Detail)(nil)

func TestDetail(t *testing.T) {
	t.Parallel()

	t.Run("MustNewDetail OK", func(t *testing.T) {
		d := problem.MustNewDetail("Insufficient funds")
		assert.NotNil(t, d)
	})
	t.Run("MustNew Panics", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				err := r.(error)
				assert.True(t, errors.IsEmpty(err), "Error should have behaviour empty but was: %+v", err)
			}
		}()
		d := problem.MustNewDetail("")
		assert.Nil(t, d)
	})
	t.Run("Options fails", func(t *testing.T) {
		d, err := problem.NewDetail("You are the problem ;-)", problem.WithExtensionString("key only"))
		assert.Nil(t, d)
		assert.True(t, errors.IsNotValid(err), "Error should have behaviour not valid: %+v", err)
	})
	t.Run("WithExtensionString", func(t *testing.T) {
		d, err := problem.NewDetail("You are the problem ;-)", problem.WithExtensionString("key", "val"))
		assert.NoError(t, err)
		assert.Exactly(t, []string{"key", "val"}, d.Extension)
	})
	t.Run("WithExtensionInt", func(t *testing.T) {
		d, err := problem.NewDetail("You are the problem ;-)", problem.WithExtensionInt("key", -4711))
		assert.NoError(t, err)
		assert.Exactly(t, []string{"key", "-4711"}, d.Extension)
	})
	t.Run("WithExtensionUint", func(t *testing.T) {
		d, err := problem.NewDetail("You are the problem ;-)", problem.WithExtensionUint("key", 4711))
		assert.NoError(t, err)
		assert.Exactly(t, []string{"key", "4711"}, d.Extension)
	})
	t.Run("WithCause", func(t *testing.T) {
		d, err := problem.NewDetail("You are the problem ;-)",
			problem.WithCause("No you are the problem", problem.WithExtensionUint("key", 4242)),
		)
		assert.NoError(t, err)
		assert.Exactly(t, []string{"key", "4242"}, d.Cause.Extension)
	})
	t.Run("WithExtensionMapString", func(t *testing.T) {
		d, err := problem.NewDetail("test",
			problem.WithExtensionMapString(map[string]string{"key3": "val3"}),
		)
		assert.NoError(t, err)
		sort.Strings(d.Extension)
		assert.Exactly(t, []string{"key3", "val3"}, d.Extension)
	})
	t.Run("WithExtensionMapStringSlice", func(t *testing.T) {
		d, err := problem.NewDetail("test",
			problem.WithExtensionMapStringSlice(url.Values{"key4": []string{"val4a", "val4b"}}),
			problem.WithExtensionMapStringSlice(http.Header{"key5": []string{"val5a", "val5b"}}),
			problem.WithExtensionMapStringSlice(mail.Header{"key6": []string{"val6a", "val6b"}}),
			problem.WithExtensionMapStringSlice(textproto.MIMEHeader{"key7": []string{"val7a", "val7b"}}),
		)
		assert.NoError(t, err)
		sort.Strings(d.Extension)
		assert.Exactly(t, []string{"key4", "key4", "key5", "key5", "key6", "key6", "key7", "key7", "val4a", "val4b", "val5a", "val5b", "val6a", "val6b", "val7a", "val7b"}, d.Extension)
	})
}

func TestDetail_Validate(t *testing.T) {
	runner := func(v interface {
		Validate() error
	}, wantErrBhf errors.BehaviourFunc) func(*testing.T) {
		return func(t *testing.T) {
			have := v.Validate()
			if wantErrBhf != nil {
				assert.True(t, wantErrBhf(have), "%+v", have)
			} else {
				assert.NoError(t, have)
			}
		}
	}
	t.Run("Title empty", runner(&problem.Detail{}, errors.IsEmpty))
	t.Run("Wrong URI: empty", runner(&problem.Detail{
		Title: "x",
		Type:  ``,
	}, errors.IsNotValid))
	t.Run("Wrong URI: dot", runner(&problem.Detail{
		Title: "x",
		Type:  `.`,
	}, errors.IsNotValid))
	t.Run("Wrong URI: %", runner(&problem.Detail{
		Title: "x",
		Type:  `http://192.168.0.%31:8080/`,
	}, errors.IsNotValid))
	t.Run("Imbalanced extension", runner(&problem.Detail{
		Title:     "x",
		Type:      `http://192.168.0.31:8080/`,
		Extension: []string{"key"},
	}, errors.IsNotValid))
}

func TestDetail_JSON(t *testing.T) {
	t.Parallel()

	d := problem.MustNewDetail("A title", problem.WithExtensionString("ke\x00y1", "val\"ue1", "ky\n2", "valu€2"))
	d.Status = 505
	d.Detail = "I could freak out!"
	d.Instance = "https://news.ycombinator.com/item?id=13679499"

	const wantJSON = `{"type":"about:blank","title":"A title","status":505,"detail":"I could freak out!","instance":"https://news.ycombinator.com/item?id=13679499","extension":{"ke\u0000y1":"val\"ue1","ky\n2":"valu€2"}}`

	t.Run("Marshal", func(t *testing.T) {
		j, err := d.MarshalJSON()
		assert.NoError(t, err)
		if s := string(j); wantJSON != s {
			t.Errorf("\nWant: %s\nHave: %s\n", wantJSON, s)
		}
	})

	t.Run("Unmarshal", func(t *testing.T) {
		data := []byte(wantJSON)
		d2 := new(problem.Detail)
		err := d2.UnmarshalJSON(data)
		assert.NoError(t, err)
		assert.Exactly(t, d, d2)
	})
}
