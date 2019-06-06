// Copyright 2016 Michal Witkowski. All Rights Reserved.
// See LICENSE for licensing terms.

package grpc_auth_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/corestoreio/pkg/util/assert"
	"github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"github.com/grpc-ecosystem/go-grpc-middleware/testing"
	pb_testproto "github.com/grpc-ecosystem/go-grpc-middleware/testing/testproto"
	"github.com/grpc-ecosystem/go-grpc-middleware/util/metautils"
	"golang.org/x/oauth2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/oauth"
	"google.golang.org/grpc/metadata"
)

var (
	commonAuthToken   = "some_good_token"
	overrideAuthToken = "override_token"

	authedMarker = "some_context_marker"
	goodPing     = &pb_testproto.PingRequest{Value: "something", SleepTimeMs: 9999}
)

// TODO(mwitkow): Add auth from metadata client dialer, which requires TLS.

func buildDummyAuthFunction(expectedScheme string, expectedToken string) func(ctx context.Context) (context.Context, error) {
	return func(ctx context.Context) (context.Context, error) {
		token, err := grpc_auth.AuthFromMD(ctx, expectedScheme)
		if err != nil {
			return nil, err
		}
		if token != expectedToken {
			return nil, grpc.Errorf(codes.PermissionDenied, "buildDummyAuthFunction bad token")
		}
		return context.WithValue(ctx, authedMarker, "marker_exists"), nil
	}
}

func assertAuthMarkerExists(t *testing.T, ctx context.Context) {
	assert.Equal(t, "marker_exists", ctx.Value(authedMarker).(string), "auth marker from buildDummyAuthFunction must be passed around")
}

type assertingPingService struct {
	pb_testproto.TestServiceServer
	T *testing.T
}

func (s *assertingPingService) PingError(ctx context.Context, ping *pb_testproto.PingRequest) (*pb_testproto.Empty, error) {
	assertAuthMarkerExists(s.T, ctx)
	return s.TestServiceServer.PingError(ctx, ping)
}

func (s *assertingPingService) PingList(ping *pb_testproto.PingRequest, stream pb_testproto.TestService_PingListServer) error {
	assertAuthMarkerExists(s.T, stream.Context())
	return s.TestServiceServer.PingList(ping, stream)
}

func ctxWithToken(ctx context.Context, scheme string, token string) context.Context {
	md := metadata.Pairs("authorization", fmt.Sprintf("%s %v", scheme, token))
	nCtx := metautils.NiceMD(md).ToOutgoing(ctx)
	return nCtx
}

func TestAuthTestSuite(t *testing.T) {
	authFunc := buildDummyAuthFunction("bearer", commonAuthToken)
	s := &grpc_testing.InterceptorTestSuite{
		TestService: &assertingPingService{&grpc_testing.TestPingService{T: t}, t},
		ServerOpts: []grpc.ServerOption{
			grpc.StreamInterceptor(grpc_auth.StreamServerInterceptor(authFunc)),
			grpc.UnaryInterceptor(grpc_auth.UnaryServerInterceptor(authFunc)),
		},
	}
	s.SetT(t)
	s.SetupSuite()
	defer s.TearDownSuite()

	t.Run("Unary_NoAuth", func(t *testing.T) {
		_, err := s.Client.Ping(s.SimpleCtx(), goodPing)
		assert.Error(t, err, "there must be an error")
		assert.Equal(t, codes.Unauthenticated, grpc.Code(err), "must error with unauthenticated")
	})
	t.Run("Unary_BadAuth", func(t *testing.T) {
		_, err := s.Client.Ping(ctxWithToken(s.SimpleCtx(), "bearer", "bad_token"), goodPing)
		assert.Error(t, err, "there must be an error")
		assert.Equal(t, codes.PermissionDenied, grpc.Code(err), "must error with permission denied")
	})
	t.Run("Unary_PassesAuth", func(t *testing.T) {
		_, err := s.Client.Ping(ctxWithToken(s.SimpleCtx(), "bearer", commonAuthToken), goodPing)
		assert.NoError(t, err, "no error must occur")
	})
	t.Run("Unary_PassesWithPerRpcCredentials", func(t *testing.T) {
		grpcCreds := oauth.TokenSource{TokenSource: &fakeOAuth2TokenSource{accessToken: commonAuthToken}}
		client := s.NewClient(grpc.WithPerRPCCredentials(grpcCreds))
		_, err := client.Ping(s.SimpleCtx(), goodPing)
		assert.NoError(t, err, "no error must occur")
	})
	t.Run("Stream_NoAuth", func(t *testing.T) {
		stream, err := s.Client.PingList(s.SimpleCtx(), goodPing)
		assert.NoError(t, err, "should not fail on establishing the stream")
		_, err = stream.Recv()
		assert.Error(t, err, "there must be an error")
		assert.Equal(t, codes.Unauthenticated, grpc.Code(err), "must error with unauthenticated")
	})
	t.Run("Stream_BadAuth", func(t *testing.T) {
		stream, err := s.Client.PingList(ctxWithToken(s.SimpleCtx(), "bearer", "bad_token"), goodPing)
		assert.NoError(t, err, "should not fail on establishing the stream")
		_, err = stream.Recv()
		assert.Error(t, err, "there must be an error")
		assert.Equal(t, codes.PermissionDenied, grpc.Code(err), "must error with permission denied")
	})
	t.Run("Stream_PassesAuth", func(t *testing.T) {
		stream, err := s.Client.PingList(ctxWithToken(s.SimpleCtx(), "Bearer", commonAuthToken), goodPing)
		assert.NoError(t, err, "should not fail on establishing the stream")
		pong, err := stream.Recv()
		assert.NoError(t, err, "no error must occur")
		assert.NotNil(t, pong, "pong must not be nil")
	})
	t.Run("Stream_PassesWithPerRpcCredentials", func(t *testing.T) {
		grpcCreds := oauth.TokenSource{TokenSource: &fakeOAuth2TokenSource{accessToken: commonAuthToken}}
		client := s.NewClient(grpc.WithPerRPCCredentials(grpcCreds))
		stream, err := client.PingList(s.SimpleCtx(), goodPing)
		assert.NoError(t, err, "should not fail on establishing the stream")
		pong, err := stream.Recv()
		assert.NoError(t, err, "no error must occur")
		assert.NotNil(t, pong, "pong must not be nil")
	})
}

type authOverrideTestService struct {
	pb_testproto.TestServiceServer
	T *testing.T
}

func (s *authOverrideTestService) AuthFuncOverride(ctx context.Context, fullMethodName string) (context.Context, error) {
	assert.NotEmpty(s.T, fullMethodName, "method name of caller is passed around")
	return buildDummyAuthFunction("bearer", overrideAuthToken)(ctx)
}

func TestAuthOverrideTestSuite(t *testing.T) {
	authFunc := buildDummyAuthFunction("bearer", commonAuthToken)
	s := &grpc_testing.InterceptorTestSuite{
		TestService: &authOverrideTestService{&assertingPingService{&grpc_testing.TestPingService{T: t}, t}, t},
		ServerOpts: []grpc.ServerOption{
			grpc.StreamInterceptor(grpc_auth.StreamServerInterceptor(authFunc)),
			grpc.UnaryInterceptor(grpc_auth.UnaryServerInterceptor(authFunc)),
		},
	}
	s.SetT(t)
	s.SetupSuite()
	defer s.TearDownSuite()

	t.Run("Unary_PassesAuth", func(t *testing.T) {
		_, err := s.Client.Ping(ctxWithToken(s.SimpleCtx(), "bearer", overrideAuthToken), goodPing)
		assert.NoError(t, err, "no error must occur")
	})
	t.Run("Stream_PassesAuth", func(t *testing.T) {
		stream, err := s.Client.PingList(ctxWithToken(s.SimpleCtx(), "Bearer", overrideAuthToken), goodPing)
		assert.NoError(t, err, "should not fail on establishing the stream")
		pong, err := stream.Recv()
		assert.NoError(t, err, "no error must occur")
		assert.NotNil(t, pong, "pong must not be nil")
	})
}

// fakeOAuth2TokenSource implements a fake oauth2.TokenSource for the purpose of credentials test.
type fakeOAuth2TokenSource struct {
	accessToken string
}

func (ts *fakeOAuth2TokenSource) Token() (*oauth2.Token, error) {
	t := &oauth2.Token{
		AccessToken: ts.accessToken,
		Expiry:      time.Now().Add(1 * time.Minute),
		TokenType:   "bearer",
	}
	return t, nil
}
