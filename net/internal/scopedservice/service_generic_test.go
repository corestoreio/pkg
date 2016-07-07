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

package scopedservice

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/config/cfgmock"
	"github.com/corestoreio/csfw/log"
	"github.com/corestoreio/csfw/store/scope"
	"github.com/corestoreio/csfw/util/cstesting"
	"github.com/corestoreio/csfw/util/errors"
	"github.com/stretchr/testify/assert"
)

func withError() Option {
	return func(s *Service) error {
		return errors.NewNotValidf("Paaaaaaaniic!")
	}
}

func TestMustNew_Default(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			err := r.(error)
			assert.True(t, errors.IsNotValid(err), "Error: %s", err)
		} else {
			t.Fatal("Expecting a Panic")
		}
	}()
	_ = MustNew(withError())
}

func TestService_MultiScope_NoFallback(t *testing.T) {
	logBuf := new(log.MutexBuffer)

	s := MustNew(
		withValue(scope.Default, 0, "Default=0"),
		withValue(scope.Website, 1, "Website=1"),
		withDebugLogger(logBuf),
	)

	if err := s.Options(withValue(scope.Store, 2, "Store=1")); err != nil {
		t.Errorf("%+v", err)
	}

	hpu := cstesting.NewHTTPParallelUsers(10, 10, 100, time.Millisecond)
	r := httptest.NewRequest("GET", "http://corestore.io", nil)
	hpu.ServeHTTP(r, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		tests := []struct {
			cfg  config.ScopedGetter
			want string
		}{
			{cfgmock.NewService().NewScoped(0, 0), "Default=0"},
			{cfgmock.NewService().NewScoped(1, 999), "Website=1"},   // store 999 not found, fall back to website
			{cfgmock.NewService().NewScoped(888, 777), "Default=0"}, // store 777 + website 888 not found, fall back to Default
			{cfgmock.NewService().NewScoped(1, 0), "Website=1"},
			{cfgmock.NewService().NewScoped(1, 2), "Store=1"},
			{cfgmock.NewService().NewScoped(334, 2), "Store=1"},
		}
		for i, test := range tests {

			cfg := s.configByScopedGetter(test.cfg)

			if have, want := cfg.value, test.want; have != want {
				t.Errorf("(%d) Have: %q Want: %q (%s)", i, have, want, cfg.scopeHash)
			}
		}
	}))

	var comparePointers = func(h1, h2 scope.Hash) {
		if have, want := reflect.ValueOf(s.scopeCache[h1]).Pointer(), reflect.ValueOf(s.scopeCache[h2]).Pointer(); have != want {
			t.Errorf("H1 Pointer: %d H2 Pointer: %d | %s => %s", have, want, h1, h2)
		}
	}
	// the second argument must have the pointer of the first argument to avoid
	// configuration duplication.
	comparePointers(scope.DefaultHash, scope.NewHash(scope.Store, 777))
	comparePointers(scope.NewHash(scope.Website, 1), scope.NewHash(scope.Store, 999))

	buf := &bytes.Buffer{}
	if err := s.DebugCache(buf); err != nil {
		t.Fatalf("%+v", err)
	}
	assert.Regexp(t, `Scope\(Default\) ID\(0\) => \[[0-9a-zA-Z]+\]=scopedConfigGeneric{lastErr: %!q\(<nil>\), scopeHash: scope.NewHash\(scope.Default, 0\)}
Scope\(Website\) ID\(1\) => \[[0-9a-zA-Z]+\]=scopedConfigGeneric{lastErr: %!q\(<nil>\), scopeHash: scope.NewHash\(scope.Website, 1\)}
Scope\(Store\) ID\(2\) => \[[0-9a-zA-Z]+\]=scopedConfigGeneric{lastErr: %!q\(<nil>\), scopeHash: scope.NewHash\(scope.Store, 2\)}
Scope\(Store\) ID\(777\) => \[[0-9a-zA-Z]+\]=scopedConfigGeneric{lastErr: %!q\(<nil>\), scopeHash: scope.NewHash\(scope.Default, 0\)}
Scope\(Store\) ID\(999\) => \[[0-9a-zA-Z]+\]=scopedConfigGeneric{lastErr: %!q\(<nil>\), scopeHash: scope.NewHash\(scope.Website, 1\)}
`, buf.String())

}
