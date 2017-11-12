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
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/corestoreio/pkg/config"
	"github.com/corestoreio/pkg/config/cfgmock"
	"github.com/corestoreio/pkg/store/scope"
	"github.com/corestoreio/pkg/util/cstesting"
	"github.com/corestoreio/errors"
	"github.com/corestoreio/log"
	"github.com/stretchr/testify/assert"
)

func TestWithConfigGetter_Panic(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			assert.NotNil(t, r)
		} else {
			t.Fatal("Expecting a panic")
		}
	}()
	_ = WithRootConfig(nil)
}

func TestWithConfigGetter(t *testing.T) {
	cfg := cfgmock.NewService()

	src, err := newService(WithRootConfig(cfg))
	assert.NoError(t, err)
	assert.NotNil(t, src.RootConfig)
}

func TestWithErrorHandler(t *testing.T) {
	var eh = func(error) http.Handler { return nil }
	s, err := newService(WithErrorHandler(eh, scope.Store.Pack(44)))
	assert.NoError(t, err, "%+v", err)
	cfg, err := s.ConfigByScopeID(scope.MakeTypeID(scope.Store, 44), 0)
	assert.NoError(t, err, "%+v", err)
	assert.NotNil(t, cfg.ErrorHandler)
	cstesting.EqualPointers(t, eh, cfg.ErrorHandler)
	cstesting.EqualPointers(t, s.ErrorHandler, defaultErrorHandler)
}

func TestWithServiceErrorHandler(t *testing.T) {
	var eh = func(error) http.Handler { return nil }
	s, err := newService(WithServiceErrorHandler(eh))
	assert.NoError(t, err)
	cstesting.EqualPointers(t, s.ErrorHandler, eh)
	assert.Nil(t, s.ErrorHandler(errors.New("Error handler returns nil")))
}

func TestOptionsError(t *testing.T) {
	opts := OptionsError(errors.NewAlreadyClosedf("Something has already been closed."))
	s, err := New(opts...)
	assert.Nil(t, s)
	assert.True(t, errors.IsAlreadyClosed(err), "%+v", err)
}

func TestOptionFactories(t *testing.T) {

	var off OptionFactoryFunc = func(config.Scoped) []Option {
		return []Option{
			withString("a value for the store 1 scope", scope.Store.Pack(1)),
			withString("a value for the website 2 scope", scope.Website.Pack(2)),
		}
	}

	of := NewOptionFactories()
	of.Register("key", off)
	assert.Exactly(t, []string{"key"}, of.Names())

	off2, err := of.Lookup("key")
	assert.NoError(t, err)
	assert.Exactly(t, fmt.Sprintf("%#v", off), fmt.Sprintf("%#v", off2)) // yes weird but it does the job

	off3, err := of.Lookup("not found")
	assert.Nil(t, off3)
	assert.True(t, errors.IsNotFound(err), "%+v", err)
}

func TestNewScopedConfigGeneric(t *testing.T) {

	scg := newScopedConfigGeneric(0, 0)
	assert.Exactly(t, scope.TypeID(0), scg.ParentID)
	assert.Exactly(t, scope.TypeID(0), scg.ScopeID)
	assert.Nil(t, scg.lastErr)
	assert.NotNil(t, scg.ErrorHandler)

	rec := httptest.NewRecorder()
	scg.ErrorHandler(errors.New("A programmer made a mistake")).ServeHTTP(rec, nil)
	assert.Exactly(t, http.StatusServiceUnavailable, rec.Code)
	assert.Contains(t, rec.Body.String(), "A programmer made a mistake")
}

func TestWithDebugLog(t *testing.T) {
	logBuf := new(log.MutexBuffer)
	srv, err := newService(WithDebugLog(logBuf))
	assert.NoError(t, err, "%+v", err)

	_, err = srv.ConfigByScopedGetter(cfgmock.NewService().NewScoped(0, 0))
	assert.NoError(t, err, "%+v", err)
	assert.Contains(t, logBuf.String(), `scopedservice.Service.ConfigByScopedGetter.IsValid requested_scope: "Type(Default) ID(0)" requested_parent_scope: "Type(Absent) ID(0)" responded_scope: "Type(Default) ID(0)"`)
}

func TestWithLogger(t *testing.T) {
	nl := log.BlackHole{}
	srv := MustNew(WithLogger(nl))
	assert.Exactly(t, nl, srv.Log)
}

func TestWithDisable(t *testing.T) {
	srv := MustNew(
		WithRootConfig(cfgmock.NewService()),
		WithDisable(true, scope.Website.Pack(2)),
		WithDisable(true, scope.Store.Pack(3)),
	)
	scpCfg, err := srv.ConfigByScope(2, 0)
	assert.NoError(t, err, "%+v", err)
	assert.True(t, scpCfg.Disabled)

	scpCfg, err = srv.ConfigByScope(22, 3)
	assert.NoError(t, err, "%+v", err)
	assert.True(t, scpCfg.Disabled)
}

func TestWithTriggerOptionFactories(t *testing.T) {
	srv := MustNew(
		WithRootConfig(cfgmock.NewService()),
		WithMarkPartiallyApplied(true, scope.Store.Pack(4)),
	)
	_, err := srv.ConfigByScope(22, 4)
	assert.True(t, errors.IsTemporary(err), "%+v", err)

	assert.NoError(t, srv.Options(WithMarkPartiallyApplied(false, scope.Store.Pack(4))))
	_, err = srv.ConfigByScope(22, 4)
	assert.NoError(t, err, "%+v")
}
