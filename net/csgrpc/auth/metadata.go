// Copyright 2016 Michal Witkowski. All Rights Reserved.
// See LICENSE for licensing terms.

package grpc_auth

import (
	"context"
	"encoding/base64"
	"strings"

	"github.com/corestoreio/errors"
	"github.com/grpc-ecosystem/go-grpc-middleware/util/metautils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const (
	// HeaderAuthorize defines the HTTP header name where to find the token
	HeaderAuthorize = "authorization"
)

// AuthFromMD is a helper function for extracting the :authorization header from
// the gRPC metadata of the request.
//
// It expects the `:authorization` header to be of a certain scheme (e.g.
// `basic`, `bearer`), in a case-insensitive format (see rfc2617, sec 1.2). If
// no such authorization is found, or the token is of wrong scheme, an error
// with gRPC status `Unauthenticated` is returned.
func AuthFromMD(ctx context.Context, expectedScheme string) (string, error) {
	token, errCode := authFromMD(ctx, expectedScheme)
	switch errCode {
	case 0:
		return token, nil
	case 1:
		return "", status.Errorf(codes.Unauthenticated, "Empty request unauthenticated with %q", expectedScheme)
	case 2:
		return "", status.Errorf(codes.Unauthenticated, "Bad authorization string")
	case 3:
		return "", status.Errorf(codes.Unauthenticated, "Request unauthenticated with %q", expectedScheme)
	default:
		return "", errors.NotValid.Newf("[auth] Internal status code %d not valid", errCode)
	}
}

func authFromMD(ctx context.Context, expectedScheme string) (string, uint8) {
	val := metautils.ExtractIncoming(ctx).Get(HeaderAuthorize)
	if val == "" {
		return "", 1
	}
	splits := strings.SplitN(val, " ", 2)
	if len(splits) < 2 {
		return "", 2
	}
	if have, want := strings.ToLower(splits[0]), strings.ToLower(expectedScheme); have != want {
		return "", 3
	}
	return splits[1], 0
}

// AddBasicAuthToOutgoingContext adds a basic authentication header to a new
// outgoing context: "authorization: Basic base64EncodedUserPass"
func AddBasicAuthToOutgoingContext(ctx context.Context, username, password string) context.Context {
	es := base64.StdEncoding.EncodeToString([]byte(username + ":" + password))
	return metadata.NewOutgoingContext(ctx, metadata.Pairs(HeaderAuthorize, "Basic "+es))
}
