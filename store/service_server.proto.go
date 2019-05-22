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
	"github.com/corestoreio/pkg/store/scope"
	"github.com/gogo/protobuf/types"
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/trace"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// todo think about an instrumented service with opentracing, metrics, etc
// todo authentication with different roles and permissions

type ServiceRPCOptions struct {
	// TODO use config package, except for logger
	Trace   bool
	Metrics bool
	Log     log.Logger
}

func NewServiceRPC(serice *Service, o ServiceRPCOptions) (*ServiceRPC, error) {
	var mErrors *stats.Int64Measure
	if o.Metrics {
		mErrors = stats.Int64("store/ServiceRPC/errors", "The number of errors encountered", stats.UnitDimensionless)
		if err := view.Register(&view.View{
			Name:        "store/ServiceRPC/errors",
			Measure:     mErrors,
			Description: "The number of errors encountered",
			Aggregation: view.Count(),
		}); err != nil {
			return nil, errors.WithStack(err)
		}
	}
	return &ServiceRPC{
		service:     serice,
		opt:         o,
		statsErrors: mErrors,
	}, nil
}

type ServiceRPC struct {
	service     *Service
	opt         ServiceRPCOptions
	statsErrors *stats.Int64Measure
}

func (sp *ServiceRPC) recordError(ctx context.Context) {
	if sp.statsErrors != nil {
		stats.Record(ctx, sp.statsErrors.M(1))
	}
}

func (sp *ServiceRPC) IsAllowedStoreID(ctx context.Context, r *ProtoIsAllowedStoreIDRequest) (*ProtoIsAllowedStoreIDResponse, error) {
	if sp.opt.Trace {
		_, span := trace.StartSpan(ctx, "store/ServiceRPC.IsAllowedStoreID")
		defer span.End()
	}
	isAllowed, storeCode, err := sp.service.IsAllowedStoreID(scope.TypeID(r.RunMode), r.StoreID)
	if sp.opt.Log != nil && sp.opt.Log.IsInfo() {
		sp.opt.Log.Info("store.ServiceRPC.IsAllowedStoreID", log.Err(err),
			log.String("request", r.String()), log.Bool("is_allowed", isAllowed), log.String("store_code", storeCode))
	}
	if err != nil {
		sp.recordError(ctx)
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return &ProtoIsAllowedStoreIDResponse{
		IsAllowed: isAllowed,
		StoreCode: storeCode,
	}, nil
}

func (sp *ServiceRPC) DefaultStoreView(ctx context.Context, _ *types.Empty) (*Store, error) {
	if sp.opt.Trace {
		_, span := trace.StartSpan(ctx, "store/ServiceRPC.DefaultStoreView")
		defer span.End()
	}
	store, err := sp.service.DefaultStoreView()
	if sp.opt.Log != nil && sp.opt.Log.IsInfo() {
		sp.opt.Log.Info("store.ServiceRPC.DefaultStoreView", log.Err(err),
			log.String("request", ""), log.Stringer("store_code", store))
	}
	if err != nil {
		sp.recordError(ctx)
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return store, nil
}

func (sp *ServiceRPC) DefaultStoreID(ctx context.Context, r *ProtoRunModeRequest) (*ProtoStoreIDWebsiteIDResponse, error) {
	if sp.opt.Trace {
		_, span := trace.StartSpan(ctx, "store/ServiceRPC.DefaultStoreID")
		defer span.End()
	}
	storeID, websiteID, err := sp.service.DefaultStoreID(scope.TypeID(r.RunMode))
	if sp.opt.Log != nil && sp.opt.Log.IsInfo() {
		sp.opt.Log.Info("store.ServiceRPC.DefaultStoreID", log.Err(err),
			log.String("request", ""), log.Uint("store_id", uint(storeID)), log.Uint("website_id", uint(websiteID)))
	}
	if err != nil {
		sp.recordError(ctx)
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return &ProtoStoreIDWebsiteIDResponse{
		StoreID:   storeID,
		WebsiteID: websiteID,
	}, nil
}

func (sp *ServiceRPC) StoreIDbyCode(ctx context.Context, r *ProtoStoreIDbyCodeRequest) (*ProtoStoreIDWebsiteIDResponse, error) {
	if sp.opt.Trace {
		_, span := trace.StartSpan(ctx, "store/ServiceRPC.StoreIDbyCode")
		defer span.End()
	}
	storeID, websiteID, err := sp.service.StoreIDbyCode(scope.TypeID(r.RunMode), r.StoreCode)
	if sp.opt.Log != nil && sp.opt.Log.IsInfo() {
		sp.opt.Log.Info("store.ServiceRPC.StoreIDbyCode", log.Err(err),
			log.String("request", ""), log.Uint("store_id", uint(storeID)), log.Uint("website_id", uint(websiteID)))
	}
	if err != nil {
		sp.recordError(ctx)
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return &ProtoStoreIDWebsiteIDResponse{
		StoreID:   storeID,
		WebsiteID: websiteID,
	}, nil
}

func (sp *ServiceRPC) AllowedStores(ctx context.Context, r *ProtoRunModeRequest) (*Stores, error) {
	if sp.opt.Trace {
		_, span := trace.StartSpan(ctx, "store/ServiceRPC.AllowedStores")
		defer span.End()
	}
	stores, err := sp.service.AllowedStores(scope.TypeID(r.RunMode))
	if sp.opt.Log != nil && sp.opt.Log.IsInfo() {
		sp.opt.Log.Info("store.ServiceRPC.AllowedStores", log.Err(err),
			log.String("request", ""), log.Int("store_count", stores.Len()))
	}
	if err != nil {
		sp.recordError(ctx)
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return stores, nil
}

func (sp *ServiceRPC) AddWebsite(ctx context.Context, r *StoreWebsite) (*types.Empty, error) {
	if sp.opt.Trace {
		_, span := trace.StartSpan(ctx, "store/ServiceRPC.AddWebsite")
		defer span.End()
	}
	err := sp.service.Options(WithWebsites(r))
	if sp.opt.Log != nil && sp.opt.Log.IsInfo() {
		sp.opt.Log.Info("store.ServiceRPC.AddWebsite", log.Err(err), log.Stringer("request", r))
	}
	if err != nil {
		sp.recordError(ctx)
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return &types.Empty{}, nil
}

func (sp *ServiceRPC) DeleteWebsite(context.Context, *ProtoIDRequest) (*types.Empty, error) {
	return nil, errors.NotImplemented.Newf("[store] DeleteWebsite not yet implemented")
}

func (sp *ServiceRPC) WebsiteByID(ctx context.Context, r *ProtoIDRequest) (*StoreWebsite, error) {
	if sp.opt.Trace {
		_, span := trace.StartSpan(ctx, "store/ServiceRPC.WebsiteByID")
		defer span.End()
	}
	w, err := sp.service.Website(r.ID)
	if sp.opt.Log != nil && sp.opt.Log.IsInfo() {
		sp.opt.Log.Info("store.ServiceRPC.WebsiteByID", log.Err(err), log.Stringer("request", r))
	}
	if err != nil {
		sp.recordError(ctx)
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return w, nil
}

func (sp *ServiceRPC) ListWebsites(ctx context.Context, _ *types.Empty) (*StoreWebsites, error) {
	if sp.opt.Trace {
		_, span := trace.StartSpan(context.Background(), "store/ServiceRPC.ListWebsites")
		defer span.End()
	}
	d := sp.service.Websites()
	return &d, nil
}

func (sp *ServiceRPC) AddGroup(ctx context.Context, r *StoreGroup) (*types.Empty, error) {
	if sp.opt.Trace {
		_, span := trace.StartSpan(ctx, "store/ServiceRPC.AddGroup")
		defer span.End()
	}
	err := sp.service.Options(WithGroups(r))
	if sp.opt.Log != nil && sp.opt.Log.IsInfo() {
		sp.opt.Log.Info("store.ServiceRPC.AddGroup", log.Err(err), log.Stringer("request", r))
	}
	if err != nil {
		sp.recordError(ctx)
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return &types.Empty{}, nil
}

func (sp *ServiceRPC) DeleteGroup(context.Context, *ProtoIDRequest) (*types.Empty, error) {
	return nil, errors.NotImplemented.Newf("[store] DeleteGroup not yet implemented")
}

func (sp *ServiceRPC) GroupByID(ctx context.Context, r *ProtoIDRequest) (*StoreGroup, error) {
	if sp.opt.Trace {
		_, span := trace.StartSpan(ctx, "store/ServiceRPC.GroupByID")
		defer span.End()
	}
	w, err := sp.service.Group(r.ID)
	if sp.opt.Log != nil && sp.opt.Log.IsInfo() {
		sp.opt.Log.Info("store.ServiceRPC.GroupByID", log.Err(err), log.Stringer("request", r))
	}
	if err != nil {
		sp.recordError(ctx)
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return w, nil
}

func (sp *ServiceRPC) ListGroups(context.Context, *types.Empty) (*StoreGroups, error) {
	if sp.opt.Trace {
		_, span := trace.StartSpan(context.Background(), "store/ServiceRPC.ListGroups")
		defer span.End()
	}
	d := sp.service.Groups()
	return &d, nil
}

func (sp *ServiceRPC) AddStore(ctx context.Context, r *Store) (*types.Empty, error) {
	if sp.opt.Trace {
		_, span := trace.StartSpan(ctx, "store/ServiceRPC.AddStore")
		defer span.End()
	}
	err := sp.service.Options(WithStores(r))
	if sp.opt.Log != nil && sp.opt.Log.IsInfo() {
		sp.opt.Log.Info("store.ServiceRPC.AddStore", log.Err(err), log.Stringer("request", r))
	}
	if err != nil {
		sp.recordError(ctx)
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return &types.Empty{}, nil
}

func (sp *ServiceRPC) DeleteStore(context.Context, *ProtoIDRequest) (*types.Empty, error) {
	return nil, errors.NotImplemented.Newf("[store] DeleteStore not yet implemented")
}

func (sp *ServiceRPC) StoreByID(ctx context.Context, r *ProtoIDRequest) (*Store, error) {
	if sp.opt.Trace {
		_, span := trace.StartSpan(ctx, "store/ServiceRPC.StoreByID")
		defer span.End()
	}
	w, err := sp.service.Store(r.ID)
	if sp.opt.Log != nil && sp.opt.Log.IsInfo() {
		sp.opt.Log.Info("store.ServiceRPC.StoreByID", log.Err(err), log.Stringer("request", r))
	}
	if err != nil {
		sp.recordError(ctx)
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return w, nil
}

func (sp *ServiceRPC) ListStores(ctx context.Context, _ *types.Empty) (*Stores, error) {
	if sp.opt.Trace {
		_, span := trace.StartSpan(context.Background(), "store/ServiceRPC.ListStores")
		defer span.End()
	}
	d := sp.service.Stores()
	return &d, nil
}
