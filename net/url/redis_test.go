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

package url_test

import (
	"testing"

	"github.com/corestoreio/csfw/net/url"
	"github.com/corestoreio/csfw/util/errors"
)

func TestRedisParseURL(t *testing.T) {
	tests := []struct {
		raw          string
		wantAddress  string
		wantPassword string
		wantDB       int64
		wantErrBhf   errors.BehaviourFunc
	}{
		{
			"localhost",
			"",
			"",
			0,
			errors.IsNotValid, // "invalid redis URL scheme",
		},
		// The error message for invalid hosts is diffferent in different
		// versions of Go, so just check that there is an error message.
		{
			"redis://weird url",
			"",
			"",
			0,
			errors.IsFatal,
		},
		{
			"redis://foo:bar:baz",
			"",
			"",
			0,
			errors.IsNotValid,
		},
		{
			"http://www.google.com",
			"",
			"",
			0,
			errors.IsNotValid, // "invalid redis URL scheme: http",
		},
		{
			"redis://localhost:6379/abc123",
			"",
			"",
			0,
			errors.IsNotValid, // "invalid database: abc123",
		},
		{
			"redis://localhost:6379/123",
			"localhost:6379",
			"",
			123,
			nil,
		},
		{
			"redis://:6379/123",
			"localhost:6379",
			"",
			123,
			nil,
		},
		{
			"redis://",
			"localhost:6379",
			"",
			0,
			nil,
		},
		{
			"redis://192.168.0.234/123",
			"192.168.0.234:6379",
			"",
			123,
			nil,
		},
		{
			"redis://192.168.0.234/",
			"",
			"",
			0,
			errors.IsNotValid,
		},
		{
			"redis://empty:SuperSecurePa55w0rd@192.168.0.234/3",
			"192.168.0.234:6379",
			"SuperSecurePa55w0rd",
			3,
			nil,
		},
	}
	for i, test := range tests {

		haveAddress, havePW, haveDB, haveErr := url.RedisParseURL(test.raw)

		if have, want := haveAddress, test.wantAddress; have != want {
			t.Errorf("(%d) Address: Have: %v Want: %v", i, have, want)
		}
		if have, want := havePW, test.wantPassword; have != want {
			t.Errorf("(%d) Password: Have: %v Want: %v", i, have, want)
		}
		if have, want := haveDB, test.wantDB; have != want {
			t.Errorf("(%d) DB: Have: %v Want: %v", i, have, want)
		}
		if test.wantErrBhf != nil {
			if have, want := test.wantErrBhf(haveErr), true; have != want {
				t.Errorf("(%d) Error: Have: %v Want: %v", i, have, want)
			}
		} else {
			if haveErr != nil {
				t.Errorf("(%d) Did not expect an Error: %+v", i, haveErr)
			}
		}
	}
}
