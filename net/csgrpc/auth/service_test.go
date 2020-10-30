// Copyright 2015-present, Cyrill @ Schumacher.fm and the CoreStore contributors
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

package grpc_auth_test

import (
	"context"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"io/ioutil"
	"math/big"
	"testing"
	"time"

	"github.com/corestoreio/errors"
	grpc_auth "github.com/corestoreio/pkg/net/csgrpc/auth"
	"github.com/corestoreio/pkg/util/assert"
	"github.com/corestoreio/pkg/util/csjwt"
	"github.com/corestoreio/pkg/util/csjwt/jwtclaim"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

func addBasicAuthToIncomingContext(ctx context.Context, username, password string) context.Context {
	es := base64.StdEncoding.EncodeToString([]byte(username + ":" + password))
	return metadata.NewIncomingContext(ctx, metadata.Pairs(grpc_auth.HeaderAuthorize, "Basic "+es))
}

func contextWithTLSCert(t *testing.T) context.Context {
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		t.Fatal(err)
	}

	ctx := peer.NewContext(context.Background(), &peer.Peer{
		AuthInfo: credentials.TLSInfo{
			State: tls.ConnectionState{
				VerifiedChains: [][]*x509.Certificate{
					{
						&x509.Certificate{
							SerialNumber: serialNumber,
							Subject: pkix.Name{
								Organization:  []string{"ORGANIZATION_NAME"},
								Country:       []string{"COUNTRY_CODE"},
								Province:      []string{"PROVINCE"},
								Locality:      []string{"CITY"},
								StreetAddress: []string{"ADDRESS"},
								PostalCode:    []string{"POSTAL_CODE"},
							},
							NotBefore:             time.Now(),
							NotAfter:              time.Now().AddDate(10, 0, 0),
							IsCA:                  true,
							ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
							KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign | x509.KeyUsageKeyEncipherment,
							BasicConstraintsValid: true,
						},
					},
				},
			},
		},
	})
	return ctx
}

func TestWithBasicAuth(t *testing.T) {
	srv, err := grpc_auth.NewService(
		grpc_auth.WithBasicAuth(grpc_auth.BasicOptions{
			Username:     "test_user",
			Password:     "123456",
			KeyInContext: "",
		}),
	)
	assert.NoError(t, err)

	t.Run("username password correct", func(t *testing.T) {
		ctx := addBasicAuthToIncomingContext(context.Background(), "test_user", "123456")
		ctx, err = srv.AuthFuncOverride(ctx, "/name.Service/FullMethodName")
		assert.NoError(t, err)
		md, ok := metadata.FromIncomingContext(ctx)
		assert.True(t, ok)
		assert.Exactly(t, []string{"test_user"}, md.Get("username"))
	})
	t.Run("username password incorrect", func(t *testing.T) {
		ctx := addBasicAuthToIncomingContext(context.Background(), "test_user", "123457")
		ctx, err = srv.AuthFuncOverride(ctx, "/name.Service/FullMethodName")

		st, ok := status.FromError(errors.Cause(err))
		assert.True(t, ok)
		assert.Exactly(t, codes.Unauthenticated, st.Code())

		assert.EqualError(t, err, "rpc error: code = Unauthenticated desc = invalid user or password")
		assert.Nil(t, ctx)
	})
	t.Run("invalid encoded basic string", func(t *testing.T) {
		ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs(grpc_auth.HeaderAuthorize, "Basic NotBase_64=Encoded"))
		ctx, err = srv.AuthFuncOverride(ctx, "/name.Service/FullMethodName")
		assert.Nil(t, ctx)
		assert.EqualError(t, err, "rpc error: code = Unauthenticated desc = [auth] No successful auth function found")
	})
	t.Run("missing metadata", func(t *testing.T) {
		ctx, err := srv.AuthFuncOverride(context.Background(), "/name.Service/FullMethodName")
		assert.Nil(t, ctx)
		assert.EqualError(t, err, "rpc error: code = Unauthenticated desc = [auth] No successful auth function found")
	})
}

func TestWithTLSAuth(t *testing.T) {
	srv, err := grpc_auth.NewService(
		grpc_auth.WithTLSAuth(
			func(ctx context.Context, fullMethodName string, incoming *x509.Certificate) (context.Context, error) {
				assert.Exactly(t, "ORGANIZATION_NAME", incoming.Subject.Organization[0])
				return ctx, nil
			}),
	)
	assert.NoError(t, err)

	t.Run("unavailable", func(t *testing.T) {
		ctx, err := srv.AuthFuncOverride(context.Background(), "/name.Service/FullMethodName")
		assert.Nil(t, ctx)
		assert.EqualError(t, err, "rpc error: code = Unauthenticated desc = [auth] No successful auth function found")
	})

	t.Run("no peer info", func(t *testing.T) {
		ctx := peer.NewContext(context.Background(), &peer.Peer{})
		ctx, err := srv.AuthFuncOverride(ctx, "/name.Service/FullMethodName")
		assert.Nil(t, ctx)
		assert.EqualError(t, err, "rpc error: code = Unauthenticated desc = unexpected peer transport credentials")
	})

	t.Run("no peer cert", func(t *testing.T) {
		ctx := peer.NewContext(context.Background(), &peer.Peer{
			AuthInfo: credentials.TLSInfo{},
		})
		ctx, err := srv.AuthFuncOverride(ctx, "/name.Service/FullMethodName")
		assert.Nil(t, ctx)
		assert.EqualError(t, err, "rpc error: code = Unauthenticated desc = could not verify peer certificate")
	})

	t.Run("auth success", func(t *testing.T) {
		ctx, err := srv.AuthFuncOverride(contextWithTLSCert(t), "/name.Service/FullMethodName")
		assert.NotNil(t, ctx)
		assert.NoError(t, err)

		p, ok := peer.FromContext(ctx)
		assert.True(t, ok)
		assert.Exactly(t, "ORGANIZATION_NAME", p.AuthInfo.(credentials.TLSInfo).State.VerifiedChains[0][0].Subject.Organization[0])
	})
}

func addTokenToIncomingContext(ctx context.Context, token string) context.Context {
	return metadata.NewIncomingContext(ctx, metadata.Pairs(grpc_auth.HeaderAuthorize, "Bearer "+token))
}

func TestWithTokenAuth(t *testing.T) {
	srv, err := grpc_auth.NewService(
		grpc_auth.WithTokenAuth(grpc_auth.TokenOptions{Token: "aTokenn"}),
	)
	assert.NoError(t, err)
	t.Run("correct", func(t *testing.T) {
		ctx := addTokenToIncomingContext(context.Background(), "aTokenn")
		ctx, err = srv.AuthFuncOverride(ctx, "/name.Service/FullMethodName")
		assert.NoError(t, err)
		assert.NotNil(t, ctx)
	})
	t.Run("incorrect", func(t *testing.T) {
		ctx := addTokenToIncomingContext(context.Background(), "aToken")
		ctx, err = srv.AuthFuncOverride(ctx, "/name.Service/FullMethodName")
		assert.EqualError(t, err, `rpc error: code = Unauthenticated desc = authorization failed with error: invalid token`)
		assert.Nil(t, ctx)
	})
	t.Run("no token at all", func(t *testing.T) {
		ctx := context.Background()
		ctx, err = srv.AuthFuncOverride(ctx, "/name.Service/FullMethodName")
		assert.EqualError(t, err, "rpc error: code = Unauthenticated desc = [auth] No successful auth function found")
		assert.Nil(t, ctx)
	})
	t.Run("is JWT and skip", func(t *testing.T) {
		ctx := addTokenToIncomingContext(context.Background(), "eyJ0eXAiOiJKV1QiLA0KICJhbGciOiJIUzI1NiJ9.eyJpc3MiOiJqb2UiLA0KICJleHAiOjEzMDA4MTkzODAsDQogImh0dHA6Ly9leGFtcGxlLmNvbS9pc19yb290Ijp0cnVlfQ.dBjftJeZ4CVP-mB92K27uhbUJU1p1r_wW1gFWFOEjXk")
		ctx, err = srv.AuthFuncOverride(ctx, "/name.Service/FullMethodName")
		assert.EqualError(t, err, "rpc error: code = Unauthenticated desc = [auth] No successful auth function found")
		assert.Nil(t, ctx)
	})
	t.Run("wrong TokenOptions", func(t *testing.T) {
		_, err := grpc_auth.NewService(
			grpc_auth.WithTokenAuth(grpc_auth.TokenOptions{}),
		)
		assert.ErrorIsKind(t, errors.NotAcceptable, err)
	})
}

func TestWithJWTAuth(t *testing.T) {
	hs256 := csjwt.NewSigningMethodHS256()
	hmacTestKey, err := ioutil.ReadFile("testdata/hmacTestKey")
	assert.NoError(t, err)
	key := csjwt.WithPassword(hmacTestKey)

	srv, err := grpc_auth.NewService(
		grpc_auth.WithJWTAuth(
			csjwt.NewKeyFunc(hs256, key),
			csjwt.NewVerification(hs256),
			grpc_auth.JWTOptions{},
		),
	)
	assert.NoError(t, err)

	t.Run("token expired", func(t *testing.T) {
		ctx := addTokenToIncomingContext(context.Background(), "eyJ0eXAiOiJKV1QiLA0KICJhbGciOiJIUzI1NiJ9.eyJpc3MiOiJqb2UiLA0KICJleHAiOjEzMDA4MTkzODAsDQogImh0dHA6Ly9leGFtcGxlLmNvbS9pc19yb290Ijp0cnVlfQ.dBjftJeZ4CVP-mB92K27uhbUJU1p1r_wW1gFWFOEjXk")
		ctx, err = srv.AuthFuncOverride(ctx, "/name.Service/FullMethodName")
		assert.Nil(t, ctx)
		assert.Contains(t, err.Error(), "rpc error: code = Unauthenticated desc = [csjwt] token claims validation failed: [jwtclaim] token is expired")
	})

	t.Run("token not JWT", func(t *testing.T) {
		ctx := addTokenToIncomingContext(context.Background(), "544bbe4d60b4f7593f2be137c7f1190d")
		ctx, err = srv.AuthFuncOverride(ctx, "/name.Service/FullMethodName")
		assert.Nil(t, ctx)
		assert.EqualError(t, err, "rpc error: code = Unauthenticated desc = [auth] No successful auth function found")
	})

	t.Run("token valid", func(t *testing.T) {
		tk := csjwt.NewToken(&jwtclaim.Store{Store: "de_CH"})
		tks, err := tk.SignedString(hs256, key)
		assert.NoError(t, err)

		ctx := addTokenToIncomingContext(context.Background(), string(tks))
		ctx, err = srv.AuthFuncOverride(ctx, "/name.Service/FullMethodName")
		assert.NotNil(t, ctx)
		assert.NoError(t, err)
		ctxTk, ok := csjwt.FromContextToken(ctx)
		assert.True(t, ok)
		s, _ := ctxTk.Claims.Get("store")
		assert.Exactly(t, "de_CH", s)
	})
}

func TestNewServiceFullChain(t *testing.T) {
	hs256 := csjwt.NewSigningMethodHS256()
	hmacTestKey, err := ioutil.ReadFile("testdata/hmacTestKey")
	assert.NoError(t, err)
	key := csjwt.WithPassword(hmacTestKey)

	srv, err := grpc_auth.NewService(
		grpc_auth.WithJWTAuth(
			csjwt.NewKeyFunc(hs256, key),
			csjwt.NewVerification(hs256),
			grpc_auth.JWTOptions{},
		),
		grpc_auth.WithBasicAuth(grpc_auth.BasicOptions{
			Username:     "test_user",
			Password:     "123456",
			KeyInContext: "",
		}),
		grpc_auth.WithTLSAuth(
			func(ctx context.Context, fullMethodName string, incoming *x509.Certificate) (context.Context, error) {
				assert.Exactly(t, "ORGANIZATION_NAME", incoming.Subject.Organization[0])
				return ctx, nil
			}),
		grpc_auth.WithTokenAuth(grpc_auth.TokenOptions{Token: "544bbe4d60b4f7593f2be137c7f1190d"}),
	)
	assert.NoError(t, err)

	t.Run("no auth header", func(t *testing.T) {
		ctx, err := srv.AuthFuncOverride(context.Background(), "/name.Service/FullMethodName")
		assert.Nil(t, ctx)
		assert.EqualError(t, err, "rpc error: code = Unauthenticated desc = [auth] No successful auth function found")
	})
	t.Run("token valid", func(t *testing.T) {
		ctx := addTokenToIncomingContext(context.Background(), "544bbe4d60b4f7593f2be137c7f1190d")
		ctx, err = srv.AuthFuncOverride(ctx, "/name.Service/FullMethodName")
		assert.NotNil(t, ctx)
		assert.NoError(t, err)
	})

	t.Run("jwt valid", func(t *testing.T) {
		tk := csjwt.NewToken(&jwtclaim.Store{Store: "de_CH"})
		tks, err := tk.SignedString(hs256, key)
		assert.NoError(t, err)

		ctx := addTokenToIncomingContext(context.Background(), string(tks))
		ctx, err = srv.AuthFuncOverride(ctx, "/name.Service/FullMethodName")
		assert.NotNil(t, ctx)
		assert.NoError(t, err)
		ctxTk, ok := csjwt.FromContextToken(ctx)
		assert.True(t, ok)
		s, _ := ctxTk.Claims.Get("store")
		assert.Exactly(t, "de_CH", s)
	})
	t.Run("username password correct", func(t *testing.T) {
		ctx := addBasicAuthToIncomingContext(context.Background(), "test_user", "123456")
		ctx, err = srv.AuthFuncOverride(ctx, "/name.Service/FullMethodName")
		assert.NoError(t, err)
		md, ok := metadata.FromIncomingContext(ctx)
		assert.True(t, ok)
		assert.Exactly(t, []string{"test_user"}, md.Get("username"))
	})
	t.Run("TLS cert valid", func(t *testing.T) {
		ctx, err := srv.AuthFuncOverride(contextWithTLSCert(t), "/name.Service/FullMethodName")
		assert.NotNil(t, ctx)
		assert.NoError(t, err)

		p, ok := peer.FromContext(ctx)
		assert.True(t, ok)
		assert.Exactly(t, "ORGANIZATION_NAME", p.AuthInfo.(credentials.TLSInfo).State.VerifiedChains[0][0].Subject.Organization[0])
	})
}
