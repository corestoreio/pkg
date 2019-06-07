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

package csgrpc

import (
	"context"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/log"
	grpc_auth "github.com/corestoreio/pkg/net/csgrpc/auth"
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
)

// ServiceOptionFn applies various options to the Service.
type ServiceOptionFn func(*Service) error

// WithErrorMetrics enables OpenCensus error counting.
func WithErrorMetrics(name string) ServiceOptionFn {
	return func(s *Service) error {
		s.statsErrors = stats.Int64(name, "The number of errors encountered", stats.UnitDimensionless)
		err := view.Register(&view.View{
			Name:        name,
			Measure:     s.statsErrors,
			Description: "The number of errors encountered",
			Aggregation: view.Count(),
		})
		return errors.WithStack(err)
	}
}

// WithLogging adds a logger otherwise logging would be completely disabled.
func WithLogging(l log.Logger) ServiceOptionFn {
	return func(s *Service) error {
		s.Log = l
		return nil
	}
}

// WithServiceAuthFuncOverrider adds a custom authentication type. If not set,
// there will be no authentication. In the subpackage `auth` a couple of
// authentication methods are provided like JWT or basic auth.
func WithServiceAuthFuncOverrider(a grpc_auth.ServiceAuthFuncOverrider) ServiceOptionFn {
	return func(s *Service) error {
		s.auth = a
		return nil
	}
}

// NewService creates a new abstract service for embedding into other services.
func NewService(sos ...ServiceOptionFn) (*Service, error) {
	s := &Service{}
	if err := s.Options(sos...); err != nil {
		return nil, errors.WithStack(err)
	}
	return s, nil
}

// Service provides general functions and options for a concrete service. A
// Service type must be embedded into other services. For example see the Store
// gRPC service.
type Service struct {
	Log         log.Logger
	auth        grpc_auth.ServiceAuthFuncOverrider
	statsErrors *stats.Int64Measure
}

// Options applies the options.
func (s *Service) Options(opts ...ServiceOptionFn) error {
	for _, o := range opts {
		if err := o(s); err != nil {
			return errors.WithStack(err)
		}
	}
	return nil
}

// RecordError records an error if WithErrorMetrics has been applied.
func (s *Service) RecordError(ctx context.Context) {
	if s.statsErrors != nil {
		stats.Record(ctx, s.statsErrors.M(1))
	}
}

// AuthFuncOverride calls the custom authentication function provided by
// ServiceRPCOptions. AuthFuncOverride gets called by the middleware of package
// "github.com/corestoreio/pkg/net/csgrpc/auth". When implementing, make
// sure that `grpc_auth.UnaryServerInterceptor(nil)` has the nil argument.
func (s *Service) AuthFuncOverride(ctx context.Context, fullMethodName string) (context.Context, error) {
	if s.auth == nil {
		return ctx, nil
	}
	return s.auth.AuthFuncOverride(ctx, fullMethodName)
}
