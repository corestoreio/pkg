// Copyright 2015-2016, Cyrill @ Schumacher.fm and the CoreStore contributors
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

package geoip

import (
	"context"

	"github.com/corestoreio/csfw/util/errors"
)

// keyctxCountry type is unexported to prevent collisions with context keys defined in
// other packages.
type keyctxCountry struct{}

// ctxCountryWrapper to prevent too much calls to runtime.convT2*
type ctxCountryWrapper struct {
	c   *Country
	err error
}

// WithContextCountry creates a new context with geoip.Country attached.
func WithContextCountry(ctx context.Context, c *Country) context.Context {
	return context.WithValue(ctx, keyctxCountry{}, ctxCountryWrapper{c: c})
}

// withContextError creates a new context with an error attached.
func withContextError(ctx context.Context, err error) context.Context {
	return context.WithValue(ctx, keyctxCountry{}, ctxCountryWrapper{err: err})
}

// FromContextCountry returns the geoip.Country in ctx if it exists or
// and error if that one exists. The error has been previously set
// by WithContextError.
func FromContextCountry(ctx context.Context) (*Country, error) {
	wrp, ok := ctx.Value(keyctxCountry{}).(ctxCountryWrapper)
	if !ok {
		return nil, errors.NewNotFoundf(errContextCountryNotFound)
	}
	if wrp.err != nil {
		return nil, wrp.err
	}
	return wrp.c, nil
}
