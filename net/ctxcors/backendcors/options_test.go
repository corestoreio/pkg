package backendcors_test

import (
	"testing"

	"github.com/corestoreio/csfw/store/scope"
)

func mustToPath(t *testing.T, fqf func(scope.Scope, int64) (string, error), s scope.Scope, id int64) string {
	fq, err := fqf(s, id)
	if err != nil {
		t.Fatal(err, s, id)
	}
	return fq
}

//func initBackend(t *testing.T) *Backend {
//	cfgStruct, err := NewConfigStructure()
//	if err != nil {
//		t.Fatal(err)
//	}
//	return New(cfgStruct)
//}

//func TestWithBackend(t *testing.T) {
//
//	c := MustNew()
//	assert.Nil(t, c.Backend)
//
//	be := initBackend(t)
//	_, err := c.Options(WithBackend(be))
//	assert.NoError(t, err)
//	assert.Exactly(t, be, c.Backend)
//}

//func TestWithBackendApplied(t *testing.T) {
//
//	be := initBackend(t)
//
//	cfgGet := cfgmock.NewService(
//		cfgmock.WithPV(cfgmock.PathValue{
//			mustToPath(t, be.NetCtxcorsExposedHeaders.FQ, scope.Website, 2):     "X-CoreStore-ID\nContent-Type\n\n",
//			mustToPath(t, be.NetCtxcorsAllowedOrigins.FQ, scope.Website, 2):     "host1.com\nhost2.com\n\n",
//			mustToPath(t, be.NetCtxcorsAllowedMethods.FQ, scope.Default, 0):     "PATCH\nDELETE",
//			mustToPath(t, be.NetCtxcorsAllowedHeaders.FQ, scope.Default, 0):     "Date,X-Header1",
//			mustToPath(t, be.NetCtxcorsAllowCredentials.FQ, scope.Website, 2):   "1",
//			mustToPath(t, be.NetCtxcorsOptionsPassthrough.FQ, scope.Website, 2): "1",
//			mustToPath(t, be.NetCtxcorsMaxAge.FQ, scope.Website, 2):             "2h",
//		}),
//	)
//
//	c := MustNew(WithBackendApplied(be, cfgGet.NewScoped(2, 4)))
//
//	assert.Exactly(t, []string{"X-Corestore-Id", "Content-Type"}, c.exposedHeaders)
//	assert.Exactly(t, []string{"host1.com", "host2.com"}, c.allowedOrigins)
//	assert.Exactly(t, []string{"PATCH", "DELETE"}, c.allowedMethods)
//	assert.Exactly(t, []string{"Date,X-Header1", "Origin"}, c.allowedHeaders)
//	assert.Exactly(t, true, c.AllowCredentials)
//	assert.Exactly(t, true, c.OptionsPassthrough)
//	assert.Exactly(t, "7200", c.maxAge)
//}

//func TestWithBackendAppliedErrors(t *testing.T) {
//
//	be := initBackend(t)
//
//	cfgErr := errors.New("Test Error")
//	cfgGet := cfgmock.NewService(
//		cfgmock.WithBool(func(_ string) (bool, error) {
//			return false, cfgErr
//		}),
//		cfgmock.WithString(func(_ string) (string, error) {
//			return "", cfgErr
//		}),
//	)
//
//	c, err := New(WithBackendApplied(be, cfgGet.NewScoped(223, 43213)))
//	assert.Nil(t, c)
//	assert.EqualError(t, err, "Route net/ctxcors/exposed_headers: Test Error\nRoute net/ctxcors/allowed_origins: Test Error\nRoute net/ctxcors/allowed_methods: Test Error\nRoute net/ctxcors/allowed_headers: Test Error\nRoute net/ctxcors/allow_credentials: Test Error\nRoute net/ctxcors/allow_credentials: Test Error\nRoute net/ctxcors/max_age: Test Error\nMaxAge: Invalid Duration seconds: 0")
//}
