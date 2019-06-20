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
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
)

// Option applies various options to the AbstractServer.
type Option func(*AbstractServer) error

// WithErrorMetrics enables OpenCensus error counting.
func WithErrorMetrics(name string) Option {
	return func(s *AbstractServer) error {
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

// WithLogger adds a logger otherwise logging would be completely disabled.
func WithLogger(l log.Logger) Option {
	return func(s *AbstractServer) error {
		s.Log = l
		return nil
	}
}

// ServerAuthFuncOverrider allows a given gRPC service implementation to
// override the global `AuthFunc`.
//
// If a service implements the AuthFuncOverride method, it takes precedence over
// the `AuthFunc` method, and will be called instead of AuthFunc for all method
// invocations within that service.
type ServerAuthFuncOverrider interface {
	AuthFuncOverride(ctx context.Context, fullMethodName string) (context.Context, error)
}

// WithServerAuthFuncOverrider adds a custom authentication type. If not set,
// there will be no authentication. In the subpackage `auth` a couple of
// authentication methods are provided like JWT or basic auth.
func WithServerAuthFuncOverrider(a ServerAuthFuncOverrider, err error) Option {
	return func(s *AbstractServer) error {
		s.auth = a
		return errors.WithStack(err)
	}
}

// NewAbstractServer creates a new abstract service for embedding into other services.
func NewAbstractServer(opts ...Option) (AbstractServer, error) {
	var s AbstractServer
	for _, o := range opts {
		if err := o(&s); err != nil {
			return AbstractServer{}, errors.WithStack(err)
		}
	}
	return s, nil
}

// AbstractServer provides general functions and options for a concrete service. A
// AbstractServer type must be embedded into other services. For example see the Store
// gRPC service.
type AbstractServer struct {
	Log         log.Logger
	auth        ServerAuthFuncOverrider
	statsErrors *stats.Int64Measure
}

// RecordError records an error if WithErrorMetrics has been applied.
func (s AbstractServer) RecordError(ctx context.Context) {
	if s.statsErrors != nil {
		stats.Record(ctx, s.statsErrors.M(1))
	}
}

// AuthFuncOverride calls the custom authentication function provided by
// ServerRPCOptions. AuthFuncOverride gets called by the middleware of package
// "github.com/corestoreio/pkg/net/csgrpc/auth". When implementing, make
// sure that `grpc_auth.UnaryServerInterceptor(nil)` has the nil argument.
func (s AbstractServer) AuthFuncOverride(ctx context.Context, fullMethodName string) (context.Context, error) {
	if s.auth == nil {
		return ctx, nil
	}
	return s.auth.AuthFuncOverride(ctx, fullMethodName)
}
