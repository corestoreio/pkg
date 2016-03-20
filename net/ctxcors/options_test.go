// Copyright (c) 2014 Olivier Poitrey <rs@dailymotion.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is furnished
// to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package ctxcors

import (
	"errors"
	"testing"
	"time"

	"github.com/corestoreio/csfw/config/cfgmock"
	"github.com/corestoreio/csfw/store/scope"
	"github.com/corestoreio/csfw/util/log"
	"github.com/stretchr/testify/assert"
)

func mustToPath(t *testing.T, fqf func(scope.Scope, int64) (string, error), s scope.Scope, id int64) string {
	fq, err := fqf(s, id)
	if err != nil {
		t.Fatal(err, s, id)
	}
	return fq
}

func TestWithOptionsPassthrough(t *testing.T) {
	t.Parallel()

	c := MustNew()
	assert.False(t, c.OptionsPassthrough)
	if _, err := c.Options(WithOptionsPassthrough()); err != nil {
		t.Fatal(err)
	}
	assert.True(t, c.OptionsPassthrough)
}

func TestWithAllowCredentials(t *testing.T) {
	t.Parallel()

	c := MustNew()
	assert.False(t, c.AllowCredentials)
	if _, err := c.Options(WithAllowCredentials()); err != nil {
		t.Fatal(err)
	}
	assert.True(t, c.AllowCredentials)
}

func TestWithMaxAge(t *testing.T) {
	t.Parallel()

	c := MustNew()
	_, err := c.Options(WithMaxAge(-1 * time.Second))
	assert.EqualError(t, err, "MaxAge: Invalid Duration seconds: -1")

	c = MustNew()
	_, err = c.Options(WithMaxAge(2 * time.Second))
	assert.NoError(t, err)
	assert.Exactly(t, "2", c.maxAge)
}

func TestWithLogger(t *testing.T) {
	t.Parallel()

	c := MustNew()
	assert.Exactly(t, log.BlackHole{}, c.Log)

	logga := log.NewBlackHole()
	_, err := c.Options(WithLogger(&logga))
	assert.NoError(t, err)
	assert.Exactly(t, &logga, c.Log)
}

func initBackend(t *testing.T) *PkgBackend {
	cfgStruct, err := NewConfigStructure()
	if err != nil {
		t.Fatal(err)
	}
	return NewBackend(cfgStruct)
}

func TestWithBackend(t *testing.T) {
	t.Parallel()

	c := MustNew()
	assert.Nil(t, c.Backend)

	be := initBackend(t)
	_, err := c.Options(WithBackend(be))
	assert.NoError(t, err)
	assert.Exactly(t, be, c.Backend)
}

func TestWithBackendApplied(t *testing.T) {
	t.Parallel()
	be := initBackend(t)

	cfgGet := cfgmock.NewService(
		cfgmock.WithPV(cfgmock.PathValue{
			mustToPath(t, be.NetCtxcorsExposedHeaders.FQ, scope.WebsiteID, 2):     "X-CoreStore-ID\nContent-Type\n\n",
			mustToPath(t, be.NetCtxcorsAllowedOrigins.FQ, scope.WebsiteID, 2):     "host1.com\nhost2.com\n\n",
			mustToPath(t, be.NetCtxcorsAllowedMethods.FQ, scope.DefaultID, 0):     "PATCH\nDELETE",
			mustToPath(t, be.NetCtxcorsAllowedHeaders.FQ, scope.DefaultID, 0):     "Date,X-Header1",
			mustToPath(t, be.NetCtxcorsAllowCredentials.FQ, scope.WebsiteID, 2):   "1",
			mustToPath(t, be.NetCtxcorsOptionsPassthrough.FQ, scope.WebsiteID, 2): "1",
			mustToPath(t, be.NetCtxcorsMaxAge.FQ, scope.WebsiteID, 2):             "2h",
		}),
	)

	c := MustNew(WithBackendApplied(be, cfgGet.NewScoped(2, 4)))

	assert.Exactly(t, []string{"X-Corestore-Id", "Content-Type"}, c.exposedHeaders)
	assert.Exactly(t, []string{"host1.com", "host2.com"}, c.allowedOrigins)
	assert.Exactly(t, []string{"PATCH", "DELETE"}, c.allowedMethods)
	assert.Exactly(t, []string{"Date,X-Header1", "Origin"}, c.allowedHeaders)
	assert.Exactly(t, true, c.AllowCredentials)
	assert.Exactly(t, true, c.OptionsPassthrough)
	assert.Exactly(t, "7200", c.maxAge)
}

func TestWithBackendAppliedErrors(t *testing.T) {
	t.Parallel()
	be := initBackend(t)

	cfgErr := errors.New("Test Error")
	cfgGet := cfgmock.NewService(
		cfgmock.WithBool(func(_ string) (bool, error) {
			return false, cfgErr
		}),
		cfgmock.WithString(func(_ string) (string, error) {
			return "", cfgErr
		}),
	)

	c, err := New(WithBackendApplied(be, cfgGet.NewScoped(223, 43213)))
	assert.Nil(t, c)
	assert.EqualError(t, err, "Route net/ctxcors/exposed_headers: Test Error\nRoute net/ctxcors/allowed_origins: Test Error\nRoute net/ctxcors/allowed_methods: Test Error\nRoute net/ctxcors/allowed_headers: Test Error\nRoute net/ctxcors/allow_credentials: Test Error\nRoute net/ctxcors/allow_credentials: Test Error\nRoute net/ctxcors/max_age: Test Error\nMaxAge: Invalid Duration seconds: 0")
}
