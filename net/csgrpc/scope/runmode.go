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

package grpc_scope

import (
	"context"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/pkg/store"
	"github.com/corestoreio/pkg/store/scope"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/go-grpc-middleware/util/metautils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const scopeKey = "csgrpc-store-scope"

// both are the runmode

// RunModeNewOutgoingContext use on client side to transmit a store/scope to the
// server.
func RunModeNewOutgoingContext(ctx context.Context, runMode scope.TypeID) context.Context {
	return metadata.NewOutgoingContext(ctx, metadata.Pairs(scopeKey, runMode.ToIntString()))
}

// runModeFromIncomingContext extracts the store/scope on the server side transmitted by a
// client. Returns 0,nil in case no scope is present.
func runModeFromIncomingContext(ctx context.Context) (scope.TypeID, error) {
	val := metautils.ExtractIncoming(ctx).Get(scopeKey)
	if val == "" {
		return 0, nil
	}
	tID, err := scope.MakeTypeIDString(val)
	if err != nil {
		return 0, errors.WithStack(err)
	}
	return tID, nil
}

// UnaryServerInterceptor returns a new unary server interceptors that performs
// per-request the store ID and website ID detection by the run mode. Store and
// website ID gets packed into the context with function scope.WithContext.
// Which can be later extracted via scope.FromContext. A client must add the run
// mode to the context via function RunModeNewOutgoingContext.
func UnaryServerInterceptor(sf store.Finder) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

		runMode, err := runModeFromIncomingContext(ctx)
		if runMode == 0 && err == nil {
			return nil, status.Error(codes.InvalidArgument, "[grpc_scope] Run mode in context not found. Client needs to add the run mode in the context with function grpc_scope.RunModeNewOutgoingContext")
		}
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "[grpc_scope] Cannot parse run mode: %s", err)
		}

		websiteID, storeID, err := sf.DefaultStoreID(runMode)
		if err != nil {
			return nil, status.Errorf(codes.FailedPrecondition, "[grpc_scope] Failed to query DefaultStoreID with run mode %s. Error: %s", runMode, err)
		}

		return handler(scope.WithContext(ctx, websiteID, storeID), req)
	}
}

// StreamServerInterceptor returns a new stream server interceptor that performs
// per-request the store ID and website ID detection by the run mode. Store and
// website ID gets packed into the context with function scope.WithContext.
// Which can be later extracted via scope.FromContext. A client must add the run
// mode to the context via function RunModeNewOutgoingContext.
func StreamServerInterceptor(sf store.Finder) grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, _ *grpc.StreamServerInfo, handler grpc.StreamHandler) error {

		runMode, err := runModeFromIncomingContext(stream.Context())
		if runMode == 0 && err == nil {
			return status.Error(codes.InvalidArgument, "[grpc_scope] Run mode in context not found. Client needs to add the run mode in the context with function grpc_scope.RunModeNewOutgoingContext")
		}
		if err != nil {
			return status.Errorf(codes.InvalidArgument, "[grpc_scope] Cannot parse run mode: %s", err)
		}

		websiteID, storeID, err := sf.DefaultStoreID(runMode)
		if err != nil {
			return status.Errorf(codes.FailedPrecondition, "[grpc_scope] Failed to query DefaultStoreID with run mode %s. Error: %s", runMode, err)
		}

		wrapped := grpc_middleware.WrapServerStream(stream)
		wrapped.WrappedContext = scope.WithContext(stream.Context(), websiteID, storeID)
		return handler(srv, wrapped)
	}
}
