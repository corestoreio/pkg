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

package grpc_auth_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/corestoreio/pkg/net/csgrpc/auth"
	"github.com/corestoreio/pkg/util/assert"
	"github.com/corestoreio/pkg/util/csjwt"
	"github.com/corestoreio/pkg/util/csjwt/jwtclaim"
	"google.golang.org/grpc/metadata"
)

var _ grpc_auth.ServiceAuthFuncOverrider = (*grpc_auth.JWT)(nil)

func TestNewJWT(t *testing.T) {
	key := csjwt.WithPasswordRandom()
	hs256 := csjwt.MustSigningMethodFactory(csjwt.HS256)

	ajt := grpc_auth.NewJWT(csjwt.NewKeyFunc(
		hs256,
		key,
	), hs256)

	t.Run("valid token", func(t *testing.T) {
		clientToken := csjwt.NewToken(&jwtclaim.Store{
			Store:  "ch-de",
			UserID: "peter.steam@hotwater.ai",
		})
		clientTokenRaw, err := clientToken.SignedString(hs256, key)
		assert.NoError(t, err)

		ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs("authorization", fmt.Sprintf("bearer %s", clientTokenRaw)))
		ctx, err = ajt.AuthFuncOverride(ctx, "fullMetalJacket")
		assert.NoError(t, err)
		tk, ok := csjwt.FromContextToken(ctx)
		assert.True(t, ok)

		assert.Exactly(t, "HS256", tk.Header.Alg())
		assert.Exactly(t, "peter.steam@hotwater.ai", tk.Claims.(*jwtclaim.Store).UserID)
	})

	t.Run("token not in header", func(t *testing.T) {
		ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs("authorization", fmt.Sprintf("beaLeL clientTokenRaw")))
		_, err := ajt.AuthFuncOverride(ctx, "fullMetalJacket")
		assert.EqualError(t, err, "rpc error: code = Unauthenticated desc = rpc error: code = Unauthenticated desc = Request unauthenticated with \"bearer\"")
	})

	t.Run("invalid token", func(t *testing.T) {
		ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs("authorization", fmt.Sprintf("bearer clientTokenRaw.clientTokenRaw.clientTokenRaw")))
		ctx, err := ajt.AuthFuncOverride(ctx, "fullMetalJacket")
		assert.EqualError(t, err, "rpc error: code = Unauthenticated desc = [csjwt] token is malformed: parse error: syntax error near offset 0 of 'rX\x9e\x9e\xd4\xe8\x91\xe9\xd1k'")
	})

}
