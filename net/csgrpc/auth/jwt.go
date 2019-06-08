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

package grpc_auth

import (
	"context"

	"github.com/corestoreio/pkg/util/csjwt"
	"github.com/corestoreio/pkg/util/csjwt/jwtclaim"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// JWT parses and verifies tokens.
type JWT struct {
	KeyFunc             csjwt.Keyfunc // required
	*csjwt.Verification               // required

	// SchemeName optional, e.g. bearer
	SchemeName string
	// NewClaim optional creates a new custom claim, defaults to jwtclaim.Store.
	NewClaim func() csjwt.Claimer
	NewHead  func() csjwt.Header
}

// NewJWT helper function to create a JWT authentication
func NewJWT(keyFunc csjwt.Keyfunc, availableSigners ...csjwt.Signer) *JWT {
	return &JWT{
		KeyFunc:      keyFunc,
		Verification: csjwt.NewVerification(availableSigners...),
	}
}

func (j *JWT) claim() csjwt.Claimer {
	if j.NewClaim != nil {
		return j.NewClaim()
	}
	return &jwtclaim.Store{}
}

func (j *JWT) header(t *csjwt.Token) {
	if j.NewHead != nil {
		t.Header = j.NewHead()
	}
}

// AuthFuncOverride parses and verifies a token. Puts the parsed token into the
// context for later reuse. To extract the token use: csjwt.FromContextToken
func (j *JWT) AuthFuncOverride(ctx context.Context, fullMethodName string) (context.Context, error) {
	scheme := "bearer"
	if j.SchemeName != "" {
		scheme = j.SchemeName
	}

	tokenRaw, err := AuthFromMD(ctx, scheme)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	t := csjwt.NewToken(j.claim())
	j.header(t)
	if err := j.Verification.Parse(t, []byte(tokenRaw), j.KeyFunc); err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}
	return csjwt.WithContextToken(ctx, t), nil
}
