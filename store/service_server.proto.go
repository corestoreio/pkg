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

package store

import (
	"context"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/log"
	"github.com/corestoreio/pkg/net/csgrpc"
	"github.com/corestoreio/pkg/store/scope"
	"github.com/gogo/protobuf/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// NewServiceServer creates a new object for a gRPC server.
func NewServiceServer(s *Service, opts ...csgrpc.Option) (*ServiceServer, error) {
	gs, err := csgrpc.NewAbstractServer(append([]csgrpc.Option{csgrpc.WithErrorMetrics("store/ServiceServer/errors")}, opts...)...)
	if err != nil {
		return nil, err
	}
	return &ServiceServer{
		service:        s,
		AbstractServer: gs,
	}, nil
}

// ServiceServer a wrapper type for the main Service to be used in a gRPC server.
type ServiceServer struct {
	*csgrpc.AbstractServer
	service *Service
}

func (sp *ServiceServer) IsAllowedStoreID(ctx context.Context, r *ProtoIsAllowedStoreIDRequest) (*ProtoIsAllowedStoreIDResponse, error) {
	isAllowed, storeCode, err := sp.service.IsAllowedStoreID(scope.TypeID(r.RunMode), r.StoreID)
	if sp.Log != nil && sp.Log.IsInfo() {
		sp.Log.Info("store.ServiceServer.IsAllowedStoreID", log.Err(err),
			log.String("request", r.String()), log.Bool("is_allowed", isAllowed), log.String("store_code", storeCode))
	}
	if err != nil {
		sp.RecordError(ctx)
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return &ProtoIsAllowedStoreIDResponse{
		IsAllowed: isAllowed,
		StoreCode: storeCode,
	}, nil
}

func (sp *ServiceServer) DefaultStoreView(ctx context.Context, _ *types.Empty) (*Store, error) {
	store, err := sp.service.DefaultStoreView()
	if sp.Log != nil && sp.Log.IsInfo() {
		sp.Log.Info("store.ServiceServer.DefaultStoreView", log.Err(err),
			log.String("request", ""), log.Stringer("store_code", store))
	}
	if err != nil {
		sp.RecordError(ctx)
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return store, nil
}

func (sp *ServiceServer) DefaultStoreID(ctx context.Context, r *ProtoRunModeRequest) (*ProtoStoreIDWebsiteIDResponse, error) {
	websiteID, storeID, err := sp.service.DefaultStoreID(scope.TypeID(r.RunMode))
	if sp.Log != nil && sp.Log.IsInfo() {
		sp.Log.Info("store.ServiceServer.DefaultStoreID", log.Err(err),
			log.String("request", ""), log.Uint("store_id", uint(storeID)), log.Uint("website_id", uint(websiteID)))
	}
	if err != nil {
		sp.RecordError(ctx)
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return &ProtoStoreIDWebsiteIDResponse{
		StoreID:   storeID,
		WebsiteID: websiteID,
	}, nil
}

func (sp *ServiceServer) StoreIDbyCode(ctx context.Context, r *ProtoStoreIDbyCodeRequest) (*ProtoStoreIDWebsiteIDResponse, error) {
	websiteID, storeID, err := sp.service.StoreIDbyCode(scope.TypeID(r.RunMode), r.StoreCode)
	if sp.Log != nil && sp.Log.IsInfo() {
		sp.Log.Info("store.ServiceServer.StoreIDbyCode", log.Err(err),
			log.String("request", ""), log.Uint("store_id", uint(storeID)), log.Uint("website_id", uint(websiteID)))
	}
	if err != nil {
		sp.RecordError(ctx)
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return &ProtoStoreIDWebsiteIDResponse{
		StoreID:   storeID,
		WebsiteID: websiteID,
	}, nil
}

func (sp *ServiceServer) AllowedStores(ctx context.Context, r *ProtoRunModeRequest) (*Stores, error) {
	stores, err := sp.service.AllowedStores(scope.TypeID(r.RunMode))
	if sp.Log != nil && sp.Log.IsInfo() {
		sp.Log.Info("store.ServiceServer.AllowedStores", log.Err(err),
			log.String("request", ""), log.Int("store_count", stores.Len()))
	}
	if err != nil {
		sp.RecordError(ctx)
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return stores, nil
}

func (sp *ServiceServer) AddWebsite(ctx context.Context, r *StoreWebsite) (*types.Empty, error) {
	err := sp.service.Options(WithWebsites(r))
	if sp.Log != nil && sp.Log.IsInfo() {
		sp.Log.Info("store.ServiceServer.AddWebsite", log.Err(err), log.Stringer("request", r))
	}
	if err != nil {
		sp.RecordError(ctx)
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return &types.Empty{}, nil
}

func (sp *ServiceServer) DeleteWebsite(context.Context, *ProtoIDRequest) (*types.Empty, error) {
	return nil, errors.NotImplemented.Newf("[store] DeleteWebsite not yet implemented")
}

func (sp *ServiceServer) WebsiteByID(ctx context.Context, r *ProtoIDRequest) (*StoreWebsite, error) {
	w, err := sp.service.Website(r.ID)
	if sp.Log != nil && sp.Log.IsInfo() {
		sp.Log.Info("store.ServiceServer.WebsiteByID", log.Err(err), log.Stringer("request", r))
	}
	if err != nil {
		sp.RecordError(ctx)
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return w, nil
}

func (sp *ServiceServer) ListWebsites(ctx context.Context, _ *types.Empty) (*StoreWebsites, error) {
	d := sp.service.Websites()
	return &d, nil
}

func (sp *ServiceServer) AddGroup(ctx context.Context, r *StoreGroup) (*types.Empty, error) {
	err := sp.service.Options(WithGroups(r))
	if sp.Log != nil && sp.Log.IsInfo() {
		sp.Log.Info("store.ServiceServer.AddGroup", log.Err(err), log.Stringer("request", r))
	}
	if err != nil {
		sp.RecordError(ctx)
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return &types.Empty{}, nil
}

func (sp *ServiceServer) DeleteGroup(context.Context, *ProtoIDRequest) (*types.Empty, error) {
	return nil, errors.NotImplemented.Newf("[store] DeleteGroup not yet implemented")
}

func (sp *ServiceServer) GroupByID(ctx context.Context, r *ProtoIDRequest) (*StoreGroup, error) {
	w, err := sp.service.Group(r.ID)
	if sp.Log != nil && sp.Log.IsInfo() {
		sp.Log.Info("store.ServiceServer.GroupByID", log.Err(err), log.Stringer("request", r))
	}
	if err != nil {
		sp.RecordError(ctx)
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return w, nil
}

func (sp *ServiceServer) ListGroups(context.Context, *types.Empty) (*StoreGroups, error) {
	d := sp.service.Groups()
	return &d, nil
}

func (sp *ServiceServer) AddStore(ctx context.Context, r *Store) (*types.Empty, error) {
	err := sp.service.Options(WithStores(r))
	if sp.Log != nil && sp.Log.IsInfo() {
		sp.Log.Info("store.ServiceServer.AddStore", log.Err(err), log.Stringer("request", r))
	}
	if err != nil {
		sp.RecordError(ctx)
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return &types.Empty{}, nil
}

func (sp *ServiceServer) DeleteStore(context.Context, *ProtoIDRequest) (*types.Empty, error) {
	return nil, errors.NotImplemented.Newf("[store] DeleteStore not yet implemented")
}

func (sp *ServiceServer) StoreByID(ctx context.Context, r *ProtoIDRequest) (*Store, error) {
	w, err := sp.service.Store(r.ID)
	if sp.Log != nil && sp.Log.IsInfo() {
		sp.Log.Info("store.ServiceServer.StoreByID", log.Err(err), log.Stringer("request", r))
	}
	if err != nil {
		sp.RecordError(ctx)
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return w, nil
}

func (sp *ServiceServer) ListStores(ctx context.Context, _ *types.Empty) (*Stores, error) {
	d := sp.service.Stores()
	return &d, nil
}
