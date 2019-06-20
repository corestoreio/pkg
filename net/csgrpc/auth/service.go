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
	"crypto/subtle"
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"sort"
	"strings"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/log"
	"github.com/corestoreio/pkg/net/csgrpc"
	"github.com/corestoreio/pkg/util/csjwt"
	"github.com/corestoreio/pkg/util/csjwt/jwtclaim"
	"github.com/grpc-ecosystem/go-grpc-middleware/util/metautils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"

	// "github.com/corestoreio/pkg/store/scope"
	"google.golang.org/grpc/status"
)

// TODO add maybe allowed, denied IPs, or IP ranges

// WithLogger adds a logger otherwise logging would be completely disabled.
func WithLogger(l log.Logger) Option {
	return Option{
		sortOrder: 0,
		fn: func(s *service) error {
			s.log = l
			return nil
		},
	}
}

func ctxWithKeyValue(ctx context.Context, key string, value string) context.Context {
	md := metadata.Pairs(key, value)
	nCtx := metautils.NiceMD(md).ToIncoming(ctx)
	return nCtx
}

// BasicOptions sets options to WithBasicAuth.
type BasicOptions struct {
	Username string // required
	Password string // required
	// Scheme sets a custom scheme instead of default: "Basic"
	Scheme string
	// BasicAuthFunc optional custom function to compare username and password.
	// If set, then the fields Username and Password of this struct are ignored.
	BasicAuthFunc func(ctx context.Context, fullMethodName string, userName string, password string) (context.Context, error)
	// KeyInContext sets a custom key to access the username found in basic
	// auth. Defaults to "username".
	KeyInContext string
}

// WithBasicAuth uses basic authentication. Stores the username in the context for later access.
func WithBasicAuth(bo BasicOptions) Option {
	if bo.Scheme == "" {
		bo.Scheme = "Basic"
	}
	if bo.KeyInContext == "" {
		bo.KeyInContext = "username"
	}

	return Option{
		sortOrder: 7,
		fn: func(s *service) error {
			if bo.BasicAuthFunc == nil {
				if bo.Username == "" && bo.Password == "" {
					return errors.Empty.Newf("[grpc_auth] For Basic auth the username and password cannot be empty.")
				}
				bo.BasicAuthFunc = func(ctx context.Context, fullMethodName string, u string, p string) (context.Context, error) {
					if subtle.ConstantTimeCompare([]byte(bo.Username), []byte(u)) == 1 && subtle.ConstantTimeCompare([]byte(bo.Password), []byte(p)) == 1 {
						return ctx, nil
					}
					return nil, errors.Unauthorized.Newf("[auth] Invalid username or password")
				}
			}

			authFn := func(ctx context.Context, fullMethodName string) (context.Context, error) {
				if val := metautils.ExtractIncoming(ctx).Get(HeaderAuthorize); val == "" {
					if s.log != nil && s.log.IsDebug() {
						s.log.Debug("csgrpc.auth.Service.WithBasicAuth.ExtractIncoming",
							log.Bool("value_is_empty", true),
						)
					}
					return nil, unavailableError("basic/basic")
				}
				basicRaw, errCode := authFromMD(ctx, bo.Scheme)
				if errCode > 0 {
					if s.log != nil && s.log.IsDebug() {
						s.log.Debug("csgrpc.auth.Service.WithBasicAuth.authFromMD",
							log.Int("err_code", int(errCode)),
							log.String("scheme", bo.Scheme),
						)
					}
					return nil, unavailableError("basic/basic")
				}

				basicDec, err := base64.StdEncoding.DecodeString(basicRaw)
				if err != nil {
					if s.log != nil && s.log.IsDebug() {
						s.log.Debug("csgrpc.auth.Service.WithBasicAuth.base64.StdEncoding.DecodeString",
							log.Err(err),
						)
					}
					return nil, unavailableError("basic/basic")
				}

				basicStr := string(basicDec)
				colonPos := strings.IndexByte(basicStr, ':')
				if colonPos < 0 {
					return nil, status.Error(codes.Unauthenticated, `invalid basic auth format`)
				}

				user, password := basicStr[:colonPos], basicStr[colonPos+1:]
				ctx, err = bo.BasicAuthFunc(ctx, fullMethodName, user, password)
				if err != nil {
					if s.log != nil && s.log.IsInfo() {
						s.log.Info("csgrpc.auth.Service.WithBasicAuth.BasicAuthFunc", log.Err(err), log.String("username", user))
					}
					return nil, csgrpc.NewStatusBadRequestError(codes.Unauthenticated, "invalid user or password", "error", err.Error())
				}
				// Remove token from headers from here on
				return ctxWithKeyValue(ctx, bo.KeyInContext, user), nil
			}
			s.authFns = append(s.authFns, authFn)
			return nil
		},
	}
}

// WithTLSAuth checks the TLS certificate. Currently only CommonName is supported.
// The common name can be access via key "subject_common_name" in the context.
func WithTLSAuth(authorizeFunc func(ctx context.Context, fullMethodName string, incoming *x509.Certificate) (context.Context, error)) Option {
	return Option{
		sortOrder: 8,
		fn: func(s *service) error {
			authFn := func(ctx context.Context, fullMethodName string) (context.Context, error) {
				p, ok := peer.FromContext(ctx)
				if !ok {
					return nil, unavailableError("tls/peer")
				}

				tlsAuth, ok := p.AuthInfo.(credentials.TLSInfo)
				if !ok {
					return nil, status.Error(codes.Unauthenticated, "unexpected peer transport credentials")
				}

				if len(tlsAuth.State.VerifiedChains) == 0 || len(tlsAuth.State.VerifiedChains[0]) == 0 {
					return nil, status.Error(codes.Unauthenticated, "could not verify peer certificate")
				}

				ctx, err := authorizeFunc(ctx, fullMethodName, tlsAuth.State.VerifiedChains[0][0])
				if err != nil {
					return nil, status.Errorf(codes.Unauthenticated, "authorization failed with error: %s", err)
				}
				return ctx, nil
			}
			s.authFns = append(s.authFns, authFn)
			return nil
		},
	}
}

// TokenOptions to be used in WithTokenAuth
type TokenOptions struct {
	Token string
	// AuthorizeFunc defines an optional function to authorize a request.
	AuthorizeFunc func(ctx context.Context, fullMethodName string, token string) (context.Context, error)
}

// WithTokenAuth checks a simple token carried in the bearer or another optional
// scheme name.
func WithTokenAuth(to TokenOptions) Option {
	return Option{
		sortOrder: 9,
		fn: func(s *service) error {

			const schemeName = "Bearer"
			if to.AuthorizeFunc == nil && to.Token == "" {
				return errors.NotAcceptable.Newf("[grpc_auth] Either Token or AuthorizeFunc func field must be set")
			}
			if to.AuthorizeFunc == nil {
				to.AuthorizeFunc = func(ctx context.Context, _ string, token string) (context.Context, error) {
					if to.Token != token {
						return nil, errors.New("invalid token")
					}
					return ctx, nil
				}
			}

			authFn := func(ctx context.Context, fullMethodName string) (context.Context, error) {
				mdToken, errCode := authFromMD(ctx, schemeName)
				mdTokenDotCount := strings.Count(mdToken, ".")
				if s.log != nil && s.log.IsDebug() {
					s.log.Debug("csgrpc.auth.Service.WithTokenAuth.authFromMD",
						log.Int("err_code", int(errCode)),
						log.String("scheme", schemeName),
						log.Int("md_token_dot_count", mdTokenDotCount),
					)
				}
				if errCode > 0 {
					return nil, unavailableError("token/bearer")
				}
				if mdTokenDotCount == 2 {
					return nil, unavailableError("token/jwt")
				}

				ctx, err := to.AuthorizeFunc(ctx, fullMethodName, mdToken)
				if err != nil {
					return nil, status.Errorf(codes.Unauthenticated, "authorization failed with error: %s", err)
				}
				return ctx, nil
			}
			s.authFns = append(s.authFns, authFn)
			return nil
		},
	}
}

// JWTOptions sets options to WithJWTAuth
type JWTOptions struct {
	// SchemeName optional, e.g. bearer
	SchemeName    string
	TokenFactory  func() *csjwt.Token
	AuthorizeFunc func(ctx context.Context, fullMethodName string, jwtToken *csjwt.Token) (context.Context, error)
}

// WithJWTAuth parses and verifies a token. Puts the parsed token into the
// context for later reuse. To extract the token use: csjwt.FromContextToken
func WithJWTAuth(keyFunc csjwt.Keyfunc, vf *csjwt.Verification, jo JWTOptions) Option {
	return Option{
		sortOrder: 10,
		fn: func(s *service) error {

			if jo.AuthorizeFunc == nil {
				jo.AuthorizeFunc = func(ctx context.Context, fullMethodName string, jwtToken *csjwt.Token) (context.Context, error) {
					return csjwt.WithContextToken(ctx, jwtToken), nil
				}
			}

			if jo.TokenFactory == nil {
				jo.TokenFactory = func() *csjwt.Token {
					return csjwt.NewToken(&jwtclaim.Store{})
				}
			}

			if jo.SchemeName == "" {
				jo.SchemeName = "Bearer"
			}

			authFn := func(ctx context.Context, fullMethodName string) (context.Context, error) {
				tokenRaw, errCode := authFromMD(ctx, jo.SchemeName)
				if errCode > 0 {
					return nil, unavailableError("jwt/bearer")
				}
				if strings.Count(tokenRaw, ".") != 2 {
					return nil, unavailableError("jwt/token")
				}

				t := jo.TokenFactory()
				if err := vf.Parse(t, []byte(tokenRaw), keyFunc); err != nil {
					return nil, status.Error(codes.Unauthenticated, err.Error())
				}
				ctx, err := jo.AuthorizeFunc(ctx, fullMethodName, t)
				if err != nil {
					return nil, status.Errorf(codes.Unauthenticated, "authorization failed with error: %s", err)
				}
				return ctx, nil
			}
			s.authFns = append(s.authFns, authFn)
			return nil
		},
	}
}

// unavailableError gets returned by an auth type to indicate that the next auth
// type should be tried. maybe this can be improved.
type unavailableError string

func (ue unavailableError) ErrorKind() errors.Kind {
	return errors.Unavailable
}

func (ue unavailableError) Error() string {
	return fmt.Sprintf("authentication %q not available, skipping.", ue)
}

type service struct {
	log     log.Logger
	authFns []ServiceAuthFunc
}

// Option applies various settings to NewService
type Option struct {
	sortOrder int
	fn        func(*service) error
}

// NewService creates a new ServiceAuthFuncOverrider containing
// various chained authentication methods. Its function signature matches the
// option function csgrpc.WithServerAuthFuncOverrider.
func NewService(opts ...Option) (ServiceAuthFuncOverrider, error) {
	sort.Slice(opts, func(i, j int) bool {
		return opts[i].sortOrder < opts[j].sortOrder // ascending 0-9 sorting ;-)
	})
	var s service
	for _, opt := range opts {
		if err := opt.fn(&s); err != nil {
			return nil, errors.WithStack(err)
		}
	}
	return s, nil
}

// To implement this method use package csgrpc.WithServiceAuthFuncOverrider
func (s service) AuthFuncOverride(ctx context.Context, fullMethodName string) (context.Context, error) {
	for _, authFn := range s.authFns {
		newCtx, err := authFn(ctx, fullMethodName)
		switch {
		case errors.Unavailable.Match(err):
			continue
		case err != nil:
			return nil, errors.WithStack(err)
		default:
			return purgeHeader(newCtx, "authorization"), nil
		}
	}
	return nil, status.Error(codes.Unauthenticated, "[auth] No successful auth function found")
}

func purgeHeader(ctx context.Context, header string) context.Context {
	md, _ := metadata.FromIncomingContext(ctx)
	mdn := metautils.NiceMD(md)
	mdn = mdn.Del(header)
	return metadata.NewIncomingContext(ctx, metadata.MD(mdn))
}
