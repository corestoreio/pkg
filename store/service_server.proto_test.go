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

// +build csall proto

package store_test

import (
	"context"
	"fmt"
	"net"
	"testing"

	"github.com/corestoreio/pkg/net/csgrpc"
	grpc_auth "github.com/corestoreio/pkg/net/csgrpc/auth"
	"github.com/corestoreio/pkg/storage/null"
	"github.com/corestoreio/pkg/store"
	"github.com/corestoreio/pkg/store/mock"
	"github.com/corestoreio/pkg/store/scope"
	"github.com/corestoreio/pkg/util/assert"
	"github.com/corestoreio/pkg/util/cstesting"
	"github.com/gogo/grpc-example/insecure"
	"github.com/gogo/protobuf/types"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_opentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	"github.com/grpc-ecosystem/go-grpc-middleware/util/metautils"
	"github.com/opentracing/opentracing-go/mocktracer"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
)

const validToken = `im_a_valid_good_token'`
const headerScheme = `bearer`

func ctxWithToken(ctx context.Context, scheme string, token string) context.Context {
	md := metadata.Pairs("authorization", fmt.Sprintf("%s %v", scheme, token))
	nCtx := metautils.NiceMD(md).ToOutgoing(ctx)
	return nCtx
}

func initServer(t *testing.T) (port int, _ func()) {
	mockTracerServer := mocktracer.New()
	srv := mock.NewServiceEuroOZ()
	srvRPC, err := store.NewServiceServer(srv,
		csgrpc.WithErrorMetrics("store"), // overwrite default settings, just for testing
		csgrpc.WithServerAuthFuncOverrider(
			grpc_auth.NewService(
				grpc_auth.WithTokenAuth(grpc_auth.TokenOptions{
					Token: validToken,
				}),
				grpc_auth.WithBasicAuth(
					grpc_auth.BasicOptions{
						Username: "t3st",
						Password: "fa!l3d",
					},
				),
			),
		),
	)
	assert.NoError(t, err)

	s := grpc.NewServer(
		grpc.Creds(credentials.NewServerTLSFromCert(&insecure.Cert)),
		grpc_middleware.WithUnaryServerChain(
			grpc_opentracing.UnaryServerInterceptor(grpc_opentracing.WithTracer(mockTracerServer)),
			grpc_auth.UnaryServerInterceptor(nil), // must be nil because srvRPC implements the interface
		),
	)

	store.RegisterStoreServiceServer(s, srvRPC)

	port = cstesting.MustFreePort()
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	go func() {
		if err := s.Serve(lis); err != nil {
			panic(err)
		}
	}()
	return port, func() {
		t.Run("FinishedSpans Server", func(t *testing.T) {
			assert.Len(t, mockTracerServer.FinishedSpans(), 26)
		})
		s.Stop()
		lis.Close() // do not check for error "close tcp 127.0.0.1:61497: use of closed network connection"; no idea
	}
}

func initClient(t *testing.T, port int) (store.StoreServiceClient, func()) {
	mockTracerClient := mocktracer.New()
	// See https://github.com/grpc/grpc/blob/master/doc/naming.md
	// for gRPC naming standard information.
	dialAddr := fmt.Sprintf("passthrough://localhost/localhost:%d", port)
	conn, err := grpc.DialContext(
		context.Background(),
		dialAddr,
		grpc.WithTransportCredentials(credentials.NewClientTLSFromCert(insecure.CertPool, "")),
		grpc.WithBlock(),
		grpc.WithUnaryInterceptor(grpc_opentracing.UnaryClientInterceptor(grpc_opentracing.WithTracer(mockTracerClient))),
	)
	assert.NoError(t, err)

	client := store.NewStoreServiceClient(conn)
	return client, func() {
		t.Run("FinishedSpans Client", func(t *testing.T) {
			assert.Len(t, mockTracerClient.FinishedSpans(), 26)
		})
		cstesting.Close(t, conn)
	}
}

func TestNewServiceRPC(t *testing.T) {

	port, srvClose := initServer(t)
	defer srvClose()
	client, clientClose := initClient(t, port)
	defer clientClose()

	ctxToken := ctxWithToken(context.Background(), headerScheme, validToken)

	t.Run("Missing Token", func(t *testing.T) {
		rpcResp, err := client.IsAllowedStoreID(context.Background(), &store.ProtoIsAllowedStoreIDRequest{
			RunMode: uint32(scope.Website.WithID(2)),
			StoreID: 6,
		})
		assert.Nil(t, rpcResp)
		assert.EqualError(t, err, "rpc error: code = Unauthenticated desc = Request unauthenticated with bearer")
	})
	t.Run("IsAllowedStoreID_OK", func(t *testing.T) {
		rpcResp, err := client.IsAllowedStoreID(ctxToken, &store.ProtoIsAllowedStoreIDRequest{
			RunMode: uint32(scope.Website.WithID(2)),
			StoreID: 6,
		})
		assert.NoError(t, err)
		assert.Exactly(t, "nz", rpcResp.StoreCode)
		assert.True(t, rpcResp.IsAllowed)
	})
	t.Run("IsAllowedStoreID_Err", func(t *testing.T) {
		rpcResp, err := client.IsAllowedStoreID(ctxToken, &store.ProtoIsAllowedStoreIDRequest{
			RunMode: uint32(scope.Store.WithID(0)),
			StoreID: 666,
		})
		assert.Nil(t, rpcResp)
		assert.EqualError(t, err, "rpc error: code = InvalidArgument desc = [store] Cannot find Store ID 666")
	})

	t.Run("DefaultStoreID_OK", func(t *testing.T) {
		rpcResp, err := client.DefaultStoreID(ctxToken, &store.ProtoRunModeRequest{
			RunMode: uint32(scope.Website.WithID(2)),
		})
		assert.NoError(t, err)
		assert.Exactly(t, uint32(5), rpcResp.StoreID)
		assert.Exactly(t, uint32(2), rpcResp.WebsiteID)
	})
	t.Run("DefaultStoreID_Err", func(t *testing.T) {
		rpcResp, err := client.DefaultStoreID(ctxToken, &store.ProtoRunModeRequest{
			RunMode: uint32(scope.Group.WithID(110)),
		})
		assert.Nil(t, rpcResp)
		assert.EqualError(t, err, "rpc error: code = InvalidArgument desc = [store] DefaultStoreID Scope Group ID 110: [store] Cannot find Group ID 110")
	})
	t.Run("DefaultStoreView_OK", func(t *testing.T) {
		rpcResp, err := client.DefaultStoreView(ctxToken, &types.Empty{})
		assert.NoError(t, err)
		assert.Exactly(t, uint32(2), rpcResp.StoreID)
		assert.Exactly(t, "at", rpcResp.Code)
	})

	t.Run("StoreIDbyCode_OK", func(t *testing.T) {
		rpcResp, err := client.StoreIDbyCode(ctxToken, &store.ProtoStoreIDbyCodeRequest{
			RunMode:   uint32(scope.Website.WithID(1)),
			StoreCode: "uk",
		})
		assert.NoError(t, err)
		assert.Exactly(t, uint32(4), rpcResp.StoreID)
		assert.Exactly(t, uint32(1), rpcResp.WebsiteID)
	})
	t.Run("StoreIDbyCode_Err", func(t *testing.T) {
		rpcResp, err := client.StoreIDbyCode(ctxToken, &store.ProtoStoreIDbyCodeRequest{
			RunMode:   uint32(scope.Group.WithID(3)),
			StoreCode: "nsw",
		})
		assert.Nil(t, rpcResp)
		assert.EqualError(t, err, `rpc error: code = InvalidArgument desc = [store] Code "nsw" not found for runMode Type(Group) ID(3)`)
	})

	t.Run("AllowedStores_OK", func(t *testing.T) {
		rpcResp, err := client.AllowedStores(ctxToken, &store.ProtoRunModeRequest{
			RunMode: uint32(scope.Website.WithID(1)),
		})
		assert.NoError(t, err)
		assert.Exactly(t, uint32(1), rpcResp.Data[0].StoreID)
		assert.Exactly(t, uint32(2), rpcResp.Data[1].StoreID)
		assert.Exactly(t, uint32(4), rpcResp.Data[2].StoreID)
	})
	t.Run("AllowedStores_Empty", func(t *testing.T) {
		rpcResp, err := client.AllowedStores(ctxToken, &store.ProtoRunModeRequest{
			RunMode: uint32(scope.Group.WithID(333)),
		})
		assert.NoError(t, err)
		assert.Exactly(t, &store.Stores{}, rpcResp)
	})
	t.Run("AllowedStores_Err", func(t *testing.T) {
		rpcResp, err := client.AllowedStores(ctxToken, &store.ProtoRunModeRequest{
			RunMode: uint32(999999),
		})
		assert.EqualError(t, err, "rpc error: code = InvalidArgument desc = [store] Scope Absent not yet implemented.")
		assert.Nil(t, rpcResp)
	})

	t.Run("AddWebsite_OK", func(t *testing.T) {
		_, err := client.AddWebsite(ctxToken,
			&store.StoreWebsite{WebsiteID: 3, Code: `africa`, Name: null.MakeString(`Africa Continent`), SortOrder: 30, DefaultGroupID: 3, IsDefault: false},
		)
		assert.NoError(t, err)
	})
	t.Run("AddWebsite_Empty", func(t *testing.T) {
		_, err := client.AddWebsite(ctxToken, &store.StoreWebsite{})
		assert.EqualError(t, err, `rpc error: code = InvalidArgument desc = [store] Website 0: Empty code`)
	})
	t.Run("WebsiteByID_OK", func(t *testing.T) {
		protoW, err := client.WebsiteByID(ctxToken, &store.ProtoIDRequest{ID: 3})
		assert.NoError(t, err)
		assert.Exactly(t, "africa", protoW.Code)
	})
	t.Run("WebsiteByID_Err", func(t *testing.T) {
		protoW, err := client.WebsiteByID(ctxToken, &store.ProtoIDRequest{ID: 3333})
		assert.Nil(t, protoW)
		assert.EqualError(t, err, "rpc error: code = InvalidArgument desc = [store] Cannot find Website ID 3333")
	})
	t.Run("ListWebsites_OK", func(t *testing.T) {
		protoWs, err := client.ListWebsites(ctxToken, &types.Empty{})
		assert.NoError(t, err)
		assert.Exactly(t, "admin", protoWs.Data[0].Code)
		assert.Exactly(t, "euro", protoWs.Data[1].Code)
		assert.Exactly(t, "oz", protoWs.Data[2].Code)
		assert.Exactly(t, "africa", protoWs.Data[3].Code)
	})

	t.Run("AddGroup_OK", func(t *testing.T) {
		_, err := client.AddGroup(ctxToken,
			&store.StoreGroup{GroupID: 4, WebsiteID: 3, Name: `Northern States`, Code: `afno`, RootCategoryID: 2, DefaultStoreID: 0},
		)
		assert.NoError(t, err)
	})
	t.Run("AddGroup_Empty", func(t *testing.T) {
		_, err := client.AddGroup(ctxToken, &store.StoreGroup{})
		assert.EqualError(t, err, `rpc error: code = InvalidArgument desc = [store] Group 0: Empty code`)
	})
	t.Run("GroupByID_OK", func(t *testing.T) {
		protoW, err := client.GroupByID(ctxToken, &store.ProtoIDRequest{ID: 4})
		assert.NoError(t, err)
		assert.Exactly(t, "afno", protoW.Code)
	})
	t.Run("GroupByID_Err", func(t *testing.T) {
		protoW, err := client.GroupByID(ctxToken, &store.ProtoIDRequest{ID: 3333})
		assert.Nil(t, protoW)
		assert.EqualError(t, err, "rpc error: code = InvalidArgument desc = [store] Cannot find Group ID 3333")
	})
	t.Run("ListGroups_OK", func(t *testing.T) {
		protoWs, err := client.ListGroups(ctxToken, &types.Empty{})
		assert.NoError(t, err)
		assert.Exactly(t, "admin", protoWs.Data[0].Code)
		assert.Exactly(t, "dach", protoWs.Data[1].Code)
		assert.Exactly(t, "uk", protoWs.Data[2].Code)
		assert.Exactly(t, "au", protoWs.Data[3].Code)
		assert.Exactly(t, "afno", protoWs.Data[4].Code)
	})

	t.Run("AddStore_OK", func(t *testing.T) {
		_, err := client.AddStore(ctxToken,
			&store.Store{StoreID: 7, Code: `mo`, WebsiteID: 3, GroupID: 4, Name: `Morocco`, SortOrder: 40, IsActive: true},
		)
		assert.NoError(t, err)
	})
	t.Run("AddStore_Empty", func(t *testing.T) {
		_, err := client.AddStore(ctxToken, &store.Store{})
		assert.EqualError(t, err, `rpc error: code = InvalidArgument desc = [store] Store 0: Empty code`)
	})
	t.Run("StoreByID_OK", func(t *testing.T) {
		protoW, err := client.StoreByID(ctxToken, &store.ProtoIDRequest{ID: 7})
		assert.NoError(t, err)
		assert.Exactly(t, "mo", protoW.Code)
	})
	t.Run("StoreByID_Err", func(t *testing.T) {
		protoW, err := client.StoreByID(ctxToken, &store.ProtoIDRequest{ID: 3333})
		assert.Nil(t, protoW)
		assert.EqualError(t, err, "rpc error: code = InvalidArgument desc = [store] Cannot find Store ID 3333")
	})
	t.Run("ListStores_OK", func(t *testing.T) {
		protoWs, err := client.ListStores(ctxToken, &types.Empty{})
		assert.NoError(t, err)
		assert.Exactly(t, "admin", protoWs.Data[0].Code)
		assert.Exactly(t, "de", protoWs.Data[1].Code)
		assert.Exactly(t, "at", protoWs.Data[2].Code)
		assert.Exactly(t, "ch", protoWs.Data[3].Code)
		assert.Exactly(t, "uk", protoWs.Data[4].Code)
		assert.Exactly(t, "au", protoWs.Data[5].Code)
		assert.Exactly(t, "nz", protoWs.Data[6].Code)
		assert.Exactly(t, "mo", protoWs.Data[7].Code)
	})
}
