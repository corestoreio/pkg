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

package grpc_scope_test

import (
	"context"
	"testing"

	grpc_scope "github.com/corestoreio/pkg/net/csgrpc/scope"
	"github.com/corestoreio/pkg/store/scope"
	"github.com/corestoreio/pkg/util/assert"
	"github.com/grpc-ecosystem/go-grpc-middleware/util/metautils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func TestRunModeNewOutgoingContext(t *testing.T) {
	ctx := grpc_scope.RunModeNewOutgoingContext(context.Background(), scope.Website.WithID(3))
	md := metautils.ExtractOutgoing(ctx)
	tids := md.Get("csgrpc-store-scope")
	assert.Exactly(t, "33554435", tids)
	tid, err := scope.MakeTypeIDString(tids)
	assert.NoError(t, err)
	assert.Exactly(t, scope.Website.WithID(3), tid)
}

type storeFinder struct{}

func (storeFinder) DefaultStoreID(runMode scope.TypeID) (websiteID, storeID uint32, err error) {
	return 21, 11, nil
}

func (storeFinder) StoreIDbyCode(runMode scope.TypeID, storeCode string) (websiteID, storeID uint32, err error) {
	panic("implement me")
}

func TestUnaryServerInterceptor(t *testing.T) {
	t.Run("runmode not found", func(t *testing.T) {
		handler := func(ctx context.Context, req interface{}) (interface{}, error) {
			return nil, nil
		}
		retIF, err := grpc_scope.UnaryServerInterceptor(storeFinder{})(context.Background(), nil, nil, handler)
		assert.EqualError(t, err, "rpc error: code = InvalidArgument desc = [grpc_scope] Run mode in context not found. Client needs to add the run mode in the context with function grpc_scope.RunModeNewOutgoingContext")
		assert.Nil(t, retIF)
	})
	t.Run("runmode not parsable", func(t *testing.T) {
		handler := func(ctx context.Context, req interface{}) (interface{}, error) {
			return nil, nil
		}
		ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs("csgrpc-store-scope", "sdasddfas"))
		retIF, err := grpc_scope.UnaryServerInterceptor(storeFinder{})(ctx, nil, nil, handler)
		assert.EqualError(t, err, `rpc error: code = InvalidArgument desc = [grpc_scope] Cannot parse run mode: [scope] MakeTypeIDString with text "sdasddfas": strconv.ParseUint: parsing "sdasddfas": invalid syntax`)
		assert.Nil(t, retIF)
	})
	t.Run("runmode success", func(t *testing.T) {
		called := false
		handler := func(ctx context.Context, req interface{}) (interface{}, error) {
			websiteID, storeID, ok := scope.FromContext(ctx)
			assert.True(t, ok)
			assert.Exactly(t, uint32(21), websiteID)
			assert.Exactly(t, uint32(11), storeID)
			called = true
			return nil, nil
		}
		ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs("csgrpc-store-scope", scope.Website.WithID(4).ToIntString()))
		_, err := grpc_scope.UnaryServerInterceptor(storeFinder{})(ctx, nil, nil, handler)
		assert.NoError(t, err)
		assert.True(t, called)
	})
}

type mockSrvStream struct {
	ctx context.Context
}

func (mockSrvStream) SetHeader(metadata.MD) error {
	panic("implement me")
}

func (mockSrvStream) SendHeader(metadata.MD) error {
	panic("implement me")
}

func (mockSrvStream) SetTrailer(metadata.MD) {
	panic("implement me")
}

func (ms mockSrvStream) Context() context.Context {
	return ms.ctx
}

func (mockSrvStream) SendMsg(m interface{}) error {
	panic("implement me")
}

func (mockSrvStream) RecvMsg(m interface{}) error {
	panic("implement me")
}

func TestStreamServerInterceptor(t *testing.T) {
	t.Run("runmode not found", func(t *testing.T) {
		handler := func(srv interface{}, stream grpc.ServerStream) error {
			return nil
		}
		err := grpc_scope.StreamServerInterceptor(storeFinder{})(nil, mockSrvStream{ctx: context.Background()}, nil, handler)
		assert.EqualError(t, err, "rpc error: code = InvalidArgument desc = [grpc_scope] Run mode in context not found. Client needs to add the run mode in the context with function grpc_scope.RunModeNewOutgoingContext")

	})
	t.Run("runmode not parsable", func(t *testing.T) {
		handler := func(srv interface{}, stream grpc.ServerStream) error {
			return nil
		}
		ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs("csgrpc-store-scope", "sdasddfas"))
		err := grpc_scope.StreamServerInterceptor(storeFinder{})(nil, mockSrvStream{ctx: ctx}, nil, handler)
		assert.EqualError(t, err, `rpc error: code = InvalidArgument desc = [grpc_scope] Cannot parse run mode: [scope] MakeTypeIDString with text "sdasddfas": strconv.ParseUint: parsing "sdasddfas": invalid syntax`)

	})
	t.Run("runmode success", func(t *testing.T) {
		called := false
		handler := func(srv interface{}, stream grpc.ServerStream) error {
			websiteID, storeID, ok := scope.FromContext(stream.Context())
			assert.True(t, ok)
			assert.Exactly(t, uint32(21), websiteID)
			assert.Exactly(t, uint32(11), storeID)
			called = true
			return nil
		}
		stream := mockStream{
			ctx: metadata.NewIncomingContext(context.Background(), metadata.Pairs("csgrpc-store-scope", scope.Website.WithID(4).ToIntString())),
		}
		err := grpc_scope.StreamServerInterceptor(storeFinder{})(nil, stream, nil, handler)
		assert.NoError(t, err)
		assert.True(t, called)
	})
}

type mockStream struct{ ctx context.Context }

func (mockStream) SetHeader(metadata.MD) error {
	panic("implement me")
}

func (mockStream) SendHeader(metadata.MD) error {
	panic("implement me")
}

func (mockStream) SetTrailer(metadata.MD) {
	panic("implement me")
}

func (m mockStream) Context() context.Context {
	return m.ctx
}

func (mockStream) SendMsg(m interface{}) error {
	panic("implement me")
}

func (mockStream) RecvMsg(m interface{}) error {
	panic("implement me")
}
