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
	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/config/cfgmock"
	"github.com/corestoreio/csfw/store"
	"github.com/corestoreio/csfw/store/scope"
	"github.com/corestoreio/csfw/store/storemock"
	"github.com/corestoreio/csfw/util/cstesting"
	"github.com/corestoreio/csfw/util/errors"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
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

func TestExposedHeader_MultiScope(t *testing.T) {
	s := MustNew(
		withValue(scope.Default, 0, "Default=0"),
		withValue(scope.Website, 1, "Website=1"),
	)

	eur := storemock.NewEurozzyService(scope.Option{Website: scope.MockID(1)}, store.WithStorageConfig(cfgmock.NewService()))
	atStore, atErr := eur.Store(scope.MockID(2)) // ID = 2 store Austria
	if atErr != nil {
		t.Fatalf("%+v", atErr)
	}

	if err := s.Options(withValue(scope.Store, 1, "Store=1")); err != nil {
		t.Errorf("%+v", err)
	}

	hpu := cstesting.NewHTTPParallelUsers(1, 1, 100, time.Millisecond)
	r := httptest.NewRequest("GET", "http://corestore.io", nil)
	hpu.ServeHTTP(r, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		tests := []struct {
			cfg  config.ScopedGetter
			want string
		}{
			{atStore.Config, "Store=1"},
			{atStore.Website.Config, "Website=1"},
		}
		for _, test := range tests {

			cfg := s.configByScopedGetter(test.cfg)

			if have, want := cfg.value, test.want; have != want {
				t.Errorf("Have: %q Want: %q (%s)", have, want, cfg.scopeHash)
			}
		}

	}))

}
