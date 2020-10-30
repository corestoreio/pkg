// Copyright 2016 Michal Witkowski. All Rights Reserved.
// See LICENSE for licensing terms.

// +build csall proto

package grpc_auth_test

import (
	"context"
	"crypto/x509"
	"io/ioutil"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/pkg/net/csgrpc"
	grpc_auth "github.com/corestoreio/pkg/net/csgrpc/auth"
	"github.com/corestoreio/pkg/store"
	"github.com/corestoreio/pkg/util/csjwt"
	"github.com/gogo/grpc-example/insecure"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
)

var cc *grpc.ClientConn

func parseToken(token string) (struct{}, error) {
	return struct{}{}, nil
}

func userClaimFromToken(struct{}) string {
	return "foobar"
}

// Simple example of server initialization code.
func Example_serverConfig() {
	exampleAuthFunc := func(ctx context.Context) (context.Context, error) {
		token, err := grpc_auth.AuthFromMD(ctx, "bearer")
		if err != nil {
			return nil, err
		}
		tokenInfo, err := parseToken(token)
		if err != nil {
			return nil, grpc.Errorf(codes.Unauthenticated, "invalid auth token: %v", err)
		}
		grpc_ctxtags.Extract(ctx).Set("auth.sub", userClaimFromToken(tokenInfo))
		newCtx := context.WithValue(ctx, "tokenInfo", tokenInfo)
		return newCtx, nil
	}

	_ = grpc.NewServer(
		grpc.StreamInterceptor(grpc_auth.StreamServerInterceptor(exampleAuthFunc)),
		grpc.UnaryInterceptor(grpc_auth.UnaryServerInterceptor(exampleAuthFunc)),
	)
}

func ExampleNewService() {
	hs256 := csjwt.NewSigningMethodHS256()
	hmacTestKey, err := ioutil.ReadFile("testdata/hmacTestKey")
	if err != nil {
		panic(err)
	}
	key := csjwt.WithPassword(hmacTestKey)

	ss, err := store.NewService()
	if err != nil {
		panic(err)
	}
	grpcStoreService, err := store.NewServiceServer(ss, csgrpc.WithServerAuthFuncOverrider(
		grpc_auth.NewService(
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
					// You can check what ever the certificate contains
					if incoming.SerialNumber.Int64() != 236472634726328 || "ORGANIZATION_NAME" != incoming.Subject.Organization[0] {
						return nil, errors.Unauthorized.Newf("Invalid certificate")
					}
					return ctx, nil
				}),
			grpc_auth.WithTokenAuth(grpc_auth.TokenOptions{Token: "544bbe4d60b4f7593f2be137c7f1190d"}),
		),
	))
	if err != nil {
		panic(err)
	}
	s := grpc.NewServer(
		grpc.Creds(credentials.NewServerTLSFromCert(&insecure.Cert)),
		grpc_middleware.WithUnaryServerChain(
			grpc_auth.UnaryServerInterceptor(nil), // must be nil because grpcStoreService implements the interface
		),
		grpc_middleware.WithStreamServerChain(
			grpc_auth.StreamServerInterceptor(nil), // must be nil because grpcStoreService implements the interface
		),
	)

	store.RegisterStoreServiceServer(s, grpcStoreService)
	// now start the server ...
}
