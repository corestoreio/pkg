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

package config_test

import (
	"testing"
	"time"

	"github.com/corestoreio/cspkg/config"
	"github.com/corestoreio/cspkg/config/cfgpath"
	"github.com/corestoreio/cspkg/config/element"
	"github.com/corestoreio/errors"
	"github.com/stretchr/testify/assert"
)

var (
	_ config.Getter               = (*config.Service)(nil)
	_ config.Writer               = (*config.Service)(nil)
	_ config.Subscriber           = (*config.Service)(nil)
	_ element.ConfigurationWriter = (*config.Service)(nil)
)

func TestNewServiceStandard(t *testing.T) {

	srv := config.MustNewService(config.NewInMemoryStore())
	assert.NotNil(t, srv)
	url, err := srv.String(cfgpath.MustNewByParts(config.PathCSBaseURL))
	assert.NoError(t, err)
	assert.Exactly(t, config.CSBaseURL, url)
}

func TestWithDBStorage(t *testing.T) {
	t.Skip("todo")
}

func TestNotKeyNotFoundError(t *testing.T) {

	srv := config.MustNewService(config.NewInMemoryStore())
	assert.NotNil(t, srv)

	scopedSrv := srv.NewScoped(1, 1)

	flat, err := scopedSrv.String(cfgpath.NewRoute("catalog/product/enable_flat"))
	assert.True(t, errors.IsNotFound(err), "Error: %s", err)
	assert.Empty(t, flat)
	//assert.Exactly(t, scope.DefaultTypeID.String(), h.String())

	val, err := scopedSrv.String(cfgpath.NewRoute("catalog"))
	assert.Empty(t, val)
	assert.True(t, errors.IsNotValid(err), "Error: %s", err)
	assert.False(t, errors.IsNotFound(err), "Error: %s", err)
	//assert.Exactly(t, scope.TypeID(0).String(), h.String())
}

func TestService_NewScoped(t *testing.T) {

	srv := config.MustNewService(config.NewInMemoryStore())
	assert.NotNil(t, srv)

	scopedSrv := srv.NewScoped(1, 1)
	sURL, err := scopedSrv.String(cfgpath.NewRoute(config.PathCSBaseURL))
	assert.NoError(t, err)
	assert.Exactly(t, config.CSBaseURL, sURL)
	//assert.Exactly(t, scope.DefaultTypeID.String(), h.String())

}

func TestService_Write(t *testing.T) {

	srv := config.MustNewService(config.NewInMemoryStore())
	assert.NotNil(t, srv)

	p1 := cfgpath.Path{}
	err := srv.Write(p1, true)
	assert.True(t, errors.IsNotValid(err), "Error: %s", err)
}

func TestService_Types(t *testing.T) {

	basePath := cfgpath.MustNewByParts("aa/bb/cc")
	tests := []struct {
		p          cfgpath.Path
		wantErrBhf errors.BehaviourFunc
	}{
		{basePath, nil},
		{cfgpath.Path{}, errors.IsNotValid},
		{basePath.BindWebsite(10), nil},
		{basePath.BindStore(22), nil},
	}

	// vals stores all possible types for which we have functions in config.Service
	values := []interface{}{"Gopher", true, float64(3.14159), int(2016), time.Now(), []byte(`Hello Goph€rs`)}

	for vi, wantVal := range values {
		for i, test := range tests {
			testServiceTypes(t, test.p, wantVal, wantVal, vi, i, test.wantErrBhf)
			testServiceTypes(t, test.p, struct{}{}, wantVal, vi, i, test.wantErrBhf) // provokes a cast error
		}
	}
}

func testServiceTypes(t *testing.T, p cfgpath.Path, writeVal, wantVal interface{}, iFaceIDX, testIDX int, wantErrBhf errors.BehaviourFunc) {

	srv := config.MustNewService(config.NewInMemoryStore())

	if writeErr := srv.Write(p, writeVal); wantErrBhf != nil {
		assert.True(t, wantErrBhf(writeErr), "Index Value %d Index Test %d => %s", iFaceIDX, testIDX, writeErr)
	} else {
		assert.NoError(t, writeErr, "Index Value %d Index Test %d", iFaceIDX, testIDX)
	}

	var haveVal interface{}
	var haveErr error
	switch wantVal.(type) {
	case []byte:
		haveVal, haveErr = srv.Byte(p)
	case string:
		haveVal, haveErr = srv.String(p)
	case bool:
		haveVal, haveErr = srv.Bool(p)
	case float64:
		haveVal, haveErr = srv.Float64(p)
	case int:
		haveVal, haveErr = srv.Int(p)
	case time.Time:
		haveVal, haveErr = srv.Time(p)
	default:
		t.Fatalf("Unsupported type: %#v in Index Value %d Index Test %d", wantVal, iFaceIDX, testIDX)
	}

	if wantErrBhf != nil {
		assert.Empty(t, haveVal, "Index %d", testIDX)
		assert.True(t, wantErrBhf(haveErr), "Index Value %d Index Test %d => %s", iFaceIDX, testIDX, haveErr)
		assert.False(t, srv.IsSet(p))
		return
	}

	if ws, ok := writeVal.(struct{}); ok && ws == (struct{}{}) {
		assert.True(t, errors.IsNotValid(haveErr), "Error: %s", haveErr)
		assert.Empty(t, haveVal)
	} else {
		assert.NoError(t, haveErr, "Index %d", testIDX)
		assert.Exactly(t, wantVal, haveVal, "Index %d", testIDX)
		assert.True(t, srv.IsSet(p))
	}
}
