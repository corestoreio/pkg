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

	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/config/element"
	"github.com/corestoreio/csfw/config/path"
	"github.com/corestoreio/csfw/config/storage"
	"github.com/corestoreio/csfw/store/scope"
	"github.com/stretchr/testify/assert"
)

var (
	_ config.Getter     = (*config.Service)(nil)
	_ config.Writer     = (*config.Service)(nil)
	_ config.Subscriber = (*config.Service)(nil)
	_ error             = (*config.Service)(nil)
)

func TestService_ApplyDefaults(t *testing.T) {
	t.Parallel()

	pkgCfg := element.MustNewConfiguration(
		&element.Section{
			ID: path.NewRoute("contact"),
			Groups: element.NewGroupSlice(
				&element.Group{
					ID: path.NewRoute("contact"),
					Fields: element.NewFieldSlice(
						&element.Field{
							// Path: `contact/contact/enabled`,
							ID:      path.NewRoute("enabled"),
							Default: true,
						},
					),
				},
				&element.Group{
					ID: path.NewRoute("email"),
					Fields: element.NewFieldSlice(
						&element.Field{
							// Path: `contact/email/recipient_email`,
							ID:      path.NewRoute("recipient_email"),
							Default: `hello@example.com`,
						},
						&element.Field{
							// Path: `contact/email/sender_email_identity`,
							ID:      path.NewRoute("sender_email_identity"),
							Default: 2.7182818284590452353602874713527,
						},
						&element.Field{
							// Path: `contact/email/email_template`,
							ID:      path.NewRoute("email_template"),
							Default: 4711,
						},
					),
				},
			),
		},
	)
	s := config.MustNewService()
	s.ApplyDefaults(pkgCfg)
	cer, err := pkgCfg.FindFieldByID(path.NewRoute("contact", "email", "recipient_email"))
	if err != nil {
		t.Fatal(err)
	}
	email, err := s.String(path.MustNewByParts("contact/email/recipient_email")) // default scope
	assert.NoError(t, err)
	assert.Exactly(t, cer.Default.(string), email)
	assert.NoError(t, s.Close())
}

func TestNewServiceStandard(t *testing.T) {
	t.Parallel()
	srv := config.MustNewService(nil)
	assert.NotNil(t, srv)
	url, err := srv.String(path.MustNewByParts(config.PathCSBaseURL))
	assert.NoError(t, err)
	assert.Exactly(t, config.CSBaseURL, url)
}

func TestWithDBStorage(t *testing.T) {
	t.Skip("todo")
}

func TestNotKeyNotFoundError(t *testing.T) {
	t.Parallel()
	srv := config.MustNewService(nil)
	assert.NotNil(t, srv)

	scopedSrv := srv.NewScoped(1, 1, 1)

	flat, err := scopedSrv.String(path.NewRoute("catalog/product/enable_flat"))
	assert.EqualError(t, err, storage.ErrKeyNotFound.Error())
	assert.Empty(t, flat)
	assert.False(t, config.NotKeyNotFoundError(err))

	val, err := scopedSrv.String(path.NewRoute("catalog"))
	assert.Empty(t, val)
	assert.EqualError(t, err, path.ErrIncorrectPath.Error())
	assert.True(t, config.NotKeyNotFoundError(err))
}

func TestService_NewScoped(t *testing.T) {
	t.Parallel()
	srv := config.MustNewService(nil)
	assert.NotNil(t, srv)

	scopedSrv := srv.NewScoped(1, 1, 1)
	sURL, err := scopedSrv.String(path.NewRoute(config.PathCSBaseURL))
	assert.NoError(t, err)
	assert.Exactly(t, config.CSBaseURL, sURL)

}

func TestService_Write(t *testing.T) {
	t.Parallel()
	srv := config.MustNewService()
	assert.NotNil(t, srv)

	p1 := path.Path{}
	assert.EqualError(t, srv.Write(p1, true), path.ErrRouteEmpty.Error())
}

func TestService_Types(t *testing.T) {
	t.Parallel()
	basePath := path.MustNewByParts("aa/bb/cc")
	tests := []struct {
		p   path.Path
		err error
	}{
		{basePath, nil},
		{path.Path{}, path.ErrRouteEmpty},
		{basePath.Bind(scope.WebsiteID, 10), nil},
		{basePath.Bind(scope.StoreID, 22), nil},
	}

	// vals stores all possible types for which we have functions in config.Service
	values := []interface{}{"Gopher", true, float64(3.14159), int(2016), time.Now()}

	for vi, wantVal := range values {
		for i, test := range tests {
			testServiceTypes(t, test.p, wantVal, wantVal, vi, i, test.err)
			testServiceTypes(t, test.p, struct{}{}, wantVal, vi, i, test.err) // provokes a cast error
		}
	}
}

func testServiceTypes(t *testing.T, p path.Path, writeVal, wantVal interface{}, iFaceIDX, testIDX int, wantErr error) {

	srv := config.MustNewService()

	writeErr := srv.Write(p, writeVal)
	if wantErr != nil {
		assert.EqualError(t, writeErr, wantErr.Error(), "Index Value %d Index Test %d", iFaceIDX, testIDX)
	} else {
		assert.NoError(t, writeErr, "Index Value %d Index Test %d", iFaceIDX, testIDX)
	}

	var haveVal interface{}
	var haveErr error
	switch wantVal.(type) {
	case string:
		haveVal, haveErr = srv.String(p)
	case bool:
		haveVal, haveErr = srv.Bool(p)
	case float64:
		haveVal, haveErr = srv.Float64(p)
	case int:
		haveVal, haveErr = srv.Int(p)
	case time.Time:
		haveVal, haveErr = srv.DateTime(p)
	default:
		t.Fatalf("Unsupported type: %#v in Index Value %d Index Test %d", wantVal, iFaceIDX, testIDX)
	}

	if wantErr != nil {
		// if this fails for time.Time{} then my PR to assert pkg has not yet been merged :-(
		// https://github.com/stretchr/testify/pull/259
		assert.Empty(t, haveVal, "Index %d", testIDX)
		assert.EqualError(t, haveErr, wantErr.Error(), "Index %d", testIDX)
		assert.False(t, srv.IsSet(p))
		return
	}

	if ws, ok := writeVal.(struct{}); ok && ws == (struct{}{}) {
		assert.Contains(t, haveErr.Error(), "Unable to Cast struct {}{} to")
		assert.Empty(t, haveVal)
	} else {
		assert.NoError(t, haveErr, "Index %d", testIDX)
		assert.Exactly(t, wantVal, haveVal, "Index %d", testIDX)
		assert.True(t, srv.IsSet(p))
	}
}
