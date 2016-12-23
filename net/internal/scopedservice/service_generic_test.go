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
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/config/cfgmock"
	"github.com/corestoreio/csfw/store/scope"
	"github.com/corestoreio/csfw/util/cstesting"
	"github.com/corestoreio/errors"
	"github.com/corestoreio/log"
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
		withString("Default=0", scope.Default.Pack(0)),
		withString("Website=1", scope.Website.Pack(1)),
		WithDebugLog(logBuf),
	)

	if err := s.Options(withString("Store=1", scope.Store.Pack(2))); err != nil {
		t.Errorf("%+v", err)
	}

	hpu := cstesting.NewHTTPParallelUsers(10, 10, 100, time.Millisecond)
	r := httptest.NewRequest("GET", "http://corestore.io", nil)
	hpu.ServeHTTP(r, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		tests := []struct {
			cfg  config.Scoped
			want string
		}{
			0: {cfgmock.NewService().NewScoped(0, 0), "Default=0"},
			1: {cfgmock.NewService().NewScoped(1, 999), "Website=1"},   // store 999 not found, fall back to website
			2: {cfgmock.NewService().NewScoped(888, 777), "Default=0"}, // store 777 + website 888 not found, fall back to Default
			3: {cfgmock.NewService().NewScoped(1, 0), "Website=1"},
			4: {cfgmock.NewService().NewScoped(1, 2), "Store=1"},
			5: {cfgmock.NewService().NewScoped(334, 2), "Store=1"},
		}
		for i, test := range tests {

			cfg, err := s.ConfigByScopedGetter(test.cfg)
			if err != nil {
				t.Fatalf("%+v", err)
			}

			if have, want := cfg.string, test.want; have != want {
				t.Errorf("(%d) Have: %q Want: %q (%s)", i, have, want, cfg.ScopeID)
			}
		}
	}))

	var comparePointers = func(h1, h2 scope.TypeID) {
		if have, want := reflect.ValueOf(s.scopeCache[h1]).Pointer(), reflect.ValueOf(s.scopeCache[h2]).Pointer(); have != want {
			t.Errorf("H1 Pointer: %d H2 Pointer: %d | %s => %s", have, want, h1, h2)
		}
	}
	// the second argument must have the pointer of the first argument to avoid
	// configuration duplication.
	comparePointers(scope.DefaultTypeID, scope.MakeTypeID(scope.Store, 777))
	comparePointers(scope.MakeTypeID(scope.Website, 1), scope.MakeTypeID(scope.Store, 999))

	buf := &bytes.Buffer{}
	if err := s.DebugCache(buf); err != nil {
		t.Fatalf("%+v", err)
	}
	assert.Contains(t, buf.String(), `Type(Default) ID(0) => `)
	assert.Contains(t, buf.String(), `Type(Website) ID(1) => `)
	assert.Contains(t, buf.String(), `Type(Store) ID(2) => `)
	assert.Contains(t, buf.String(), `Type(Store) ID(777) => `)
	assert.Contains(t, buf.String(), `Type(Store) ID(999) => `)
}

func TestService_ClearCache(t *testing.T) {
	srv := MustNew(withString("Gopher", scope.Website.Pack(33)), WithRootConfig(cfgmock.NewService()))
	cfg, err := srv.ConfigByScope(33, 44)
	assert.NoError(t, err, "%+v", err)
	assert.Exactly(t, cfg.string, "Gopher")

	assert.NoError(t, srv.ClearCache())

	cfg, err = srv.ConfigByScopeID(scope.Website.Pack(33), 0)
	assert.True(t, errors.IsNotFound(err), "%+v", err)
	assert.Exactly(t, cfg.string, "")
}

func TestService_MultiScope_Fallback(t *testing.T) {
	// see for default values: newScopedConfig()
	s := MustNew(
		withString("Website=1", scope.Website.Pack(1)),
		withInt(130, scope.Website.Pack(1)),

		withString("Website=2", scope.Website.Pack(2)), // int must be 42

		withString("Store=1", scope.Store.Pack(1)),
		withInt(132, scope.Store.Pack(1)),

		withString("Store=2", scope.Store.Pack(2), scope.Website.Pack(1)), // int must be 130
		withString("Store=3", scope.Store.Pack(3)),                        // int must be 42
	)

	tests := []struct {
		cfg  config.Scoped
		want string
	}{
		// Default values
		0: {cfgmock.NewService().NewScoped(0, 0), configDefaultString + " => 42 => Type(Default) ID(0) => Type(Default) ID(0)"},
		// Store 99 does not exists so we get the pointer from Website 1
		1: {cfgmock.NewService().NewScoped(1, 99), "Website=1 => 130 => Type(Website) ID(1) => Type(Website) ID(1)"},
		// Store 0 does not exists so we get the pointer from Website 1
		2: {cfgmock.NewService().NewScoped(1, 0), "Website=1 => 130 => Type(Website) ID(1) => Type(Website) ID(1)"},
		// programmer made an error. Store 99 cannot have multiple parents (1
		// and 2) and Store 99 already checked above and assigned to Website 1.
		3: {cfgmock.NewService().NewScoped(2, 99), "Website=1 => 130 => Type(Website) ID(1) => Type(Website) ID(1)"},
		// Store 98 does not exists and gets pointer to Website 2 assigned
		4: {cfgmock.NewService().NewScoped(2, 98), "Website=2 => 42 => Type(Website) ID(2) => Type(Website) ID(2)"},
		// store 777 + website 888 not found, fall back to Default
		5: {cfgmock.NewService().NewScoped(888, 777), configDefaultString + " => 42 => Type(Default) ID(0) => Type(Default) ID(0)"},
		// 130 value from Website 1
		6: {cfgmock.NewService().NewScoped(1, 2), "Store=2 => 130 => Type(Store) ID(2) => Type(Website) ID(1)"},
		7: {cfgmock.NewService().NewScoped(1, 1), "Store=1 => 132 => Type(Store) ID(1) => Type(Default) ID(0)"},
		8: {cfgmock.NewService().NewScoped(1, 3), "Store=3 => 42 => Type(Store) ID(3) => Type(Default) ID(0)"},
		//{cfgmock.NewService().NewScoped(334, 2), "Store=1"},
	}
	for j, test := range tests {

		// food for the race detector
		const iterations = 10
		var wg sync.WaitGroup
		wg.Add(iterations)
		for i := 0; i < iterations; i++ {
			go func(wg *sync.WaitGroup) {
				defer wg.Done()
				cfg, err := s.ConfigByScopedGetter(test.cfg)
				if err != nil {
					t.Fatalf("%+v", err)
				}

				if have, want := fmt.Sprintf("%s => %d => %s => %s", cfg.string, cfg.int, cfg.ScopeID, cfg.ParentID), test.want; have != want {
					t.Errorf("Index %d\nHave: %q\nWant: %q\n ScopeID: %s", j, have, want, cfg.ScopeID)
				}
			}(&wg)
		}
		wg.Wait()
	}
}

func TestService_DefaultOverwritesAScope(t *testing.T) {
	scopeStore1 := scope.Store.Pack(1)
	s := MustNew(
		withString("Store=1", scopeStore1),
		WithDefaultConfig(scopeStore1),
	)
	// WithDefaultConfig overwrites the previously set
	assert.Exactly(t, configDefaultString, s.scopeCache[scopeStore1].string)

}

func TestService_ConfigByScopeID_Error(t *testing.T) {
	ws44 := scope.Website.Pack(44)
	s := MustNew(
		WithDefaultConfig(ws44),
	)
	s.scopeCache[scope.DefaultTypeID].lastErr = errors.NewNotImplementedf("Brain")
	s.scopeCache[ws44].lastErr = errors.NewAlreadyClosedf("Brain")

	t.Run("Invalid Parent", func(t *testing.T) {
		scpCfg, err := s.ConfigByScopeID(2, 1)
		assert.True(t, errors.IsNotValid(err), "%+v", err)
		assert.Exactly(t, ScopedConfig{}, scpCfg)
	})

	t.Run("Invalid Default", func(t *testing.T) {
		scpCfg, err := s.ConfigByScopeID(scope.Store.Pack(1), scope.Website.Pack(2))
		assert.True(t, errors.IsNotImplemented(err), "%+v", err)
		assert.Exactly(t, ScopedConfig{}, scpCfg)
	})

	t.Run("Invalid Default", func(t *testing.T) {
		scpCfg, err := s.ConfigByScopeID(scope.Store.Pack(55), ws44)
		assert.True(t, errors.IsAlreadyClosed(err), "%+v", err)
		assert.Exactly(t, ScopedConfig{}, scpCfg)
	})
}

func TestService_ScopedConfig_NotFound(t *testing.T) {
	s := MustNew(
		WithRootConfig(cfgmock.NewService()),
		withString("Store=1", scope.Group.Pack(98765)),
		WithErrorHandler(func(err error) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				panic(fmt.Sprintf("Should not get called. Scoped Error Handler\n\n%+v", err))
			})
		},
			scope.DefaultTypeID),
	)

	if err := s.ClearCache(); err != nil {
		t.Fatalf("%+v", err)
	}

	ctx := scope.WithContext(context.Background(), 1, 1)
	scpCfg, err := s.configByContext(ctx)
	assert.True(t, errors.IsNotFound(err), "%+v", err)
	assert.Exactly(t, ScopedConfig{}, scpCfg)
}

func TestService_ConfigByContext(t *testing.T) {
	s := MustNew(
		WithRootConfig(cfgmock.NewService()),
	)

	t.Run("Error", func(t *testing.T) {
		scpCfg, err := s.configByContext(context.Background())
		assert.True(t, errors.IsNotFound(err), "%+v", err)
		assert.Exactly(t, ScopedConfig{}, scpCfg)
	})

	t.Run("Success", func(t *testing.T) {
		scpCfg, err := s.configByContext(scope.WithContext(context.Background(), 3, 4))
		assert.NoError(t, err, "%+v", err)
		assert.Exactly(t, configDefaultString, scpCfg.string)
		assert.Exactly(t, configDefaultInt, scpCfg.int)
	})
}

func TestService_ConfigByScope(t *testing.T) {
	s := MustNew(
		WithRootConfig(cfgmock.NewService()),
		withInt(33, scope.Website.Pack(3)),
		withInt(44, scope.Store.Pack(4), scope.Website.Pack(3)),
		withInt(55, scope.Website.Pack(5)),
		withInt(66, scope.Store.Pack(6), scope.Website.Pack(5)),
	)
	s.useWebsite = false
	t.Run("Scope Store", func(t *testing.T) {
		scpCfg, err := s.ConfigByScope(3, 4)
		assert.NoError(t, err, "%+v", err)
		assert.Exactly(t, scope.Store.Pack(4), scpCfg.ScopeID)
		assert.Exactly(t, scope.Website.Pack(3), scpCfg.ParentID)
	})

	s.useWebsite = true
	t.Run("Scope Website", func(t *testing.T) {
		scpCfg, err := s.ConfigByScope(5, 6)
		assert.NoError(t, err, "%+v", err)
		assert.Exactly(t, scope.Website.Pack(5), scpCfg.ScopeID)
		assert.Exactly(t, scope.DefaultTypeID, scpCfg.ParentID)
	})
}

func TestService_OptionAfterApply_Error(t *testing.T) {
	s := MustNew()
	s.optionAfterApply = func() error {
		return errors.NewWriteFailedf("Do'h!")
	}
	err := s.Options(WithRootConfig(cfgmock.NewService()))
	assert.True(t, errors.IsWriteFailed(err), "%+v", err)
}

func TestWithOptionFactory(t *testing.T) {
	var off OptionFactoryFunc = func(config.Scoped) []Option {
		return []Option{
			withString("a value for the website 2 scope", scope.Website.Pack(2)),
			withString("a value for the store 1 scope", scope.Store.Pack(1), scope.Website.Pack(2)),
		}
	}
	s := MustNew(
		WithRootConfig(cfgmock.NewService()),
		WithOptionFactory(off),
	)

	sg := cfgmock.NewService().NewScoped(2, 1)
	scpCfg, err := s.ConfigByScopedGetter(sg)
	assert.NoError(t, err)
	assert.Exactly(t, scope.Store.Pack(1), scpCfg.ScopeID)
	assert.Exactly(t, scope.Website.Pack(2), scpCfg.ParentID)
	assert.Exactly(t, `a value for the store 1 scope`, scpCfg.string)
}

func TestWithOptionFactory_Error(t *testing.T) {
	var off OptionFactoryFunc = func(config.Scoped) []Option {
		return []Option{
			withError(),
		}
	}
	s := MustNew(
		WithRootConfig(cfgmock.NewService()),
		WithOptionFactory(off),
	)

	sg := cfgmock.NewService().NewScoped(2, 1)
	scpCfg, err := s.ConfigByScopedGetter(sg)
	assert.True(t, errors.IsNotValid(err), "%+v", err)
	assert.Exactly(t, scope.TypeID(0), scpCfg.ScopeID)
	assert.Exactly(t, scope.TypeID(0), scpCfg.ParentID)
	assert.Exactly(t, ``, scpCfg.string)
}
