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

package csgrpc_test

import (
	"context"
	"testing"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/log"
	"github.com/corestoreio/pkg/net/csgrpc"
	grpc_auth "github.com/corestoreio/pkg/net/csgrpc/auth"
	"github.com/corestoreio/pkg/util/assert"
)

func TestNewService(t *testing.T) {
	s, err := csgrpc.NewService(
		csgrpc.WithErrorMetrics("randomName"),
		csgrpc.WithLogger(log.BlackHole{}),
		csgrpc.WithServiceAuthFuncOverrider(grpc_auth.ServiceAuthFunc(func(ctx context.Context, fullMethodName string) (context.Context, error) {
			return nil, errors.NotAcceptable.Newf("upssss")
		})),
	)
	assert.NoError(t, err)
	assert.NotNil(t, s.Log)

	ctx, err := s.AuthFuncOverride(context.Background(), "methodName")
	assert.Nil(t, ctx)
	assert.ErrorIsKind(t, errors.NotAcceptable, err)
	s.RecordError(context.Background())
}
