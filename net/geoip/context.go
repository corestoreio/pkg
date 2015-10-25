// Copyright 2015, Cyrill @ Schumacher.fm and the CoreStore contributors
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

import "golang.org/x/net/context"

// ctxKey type is unexported to prevent collisions with context keys defined in
// other packages.
type ctxKey uint

// key* defines the keys to access a value in a context.Context
const (
	keyctxCountry ctxKey = iota
	keyctxErr
)

// NewContextCountry creates a new context with geoip.Country attached.
func NewContextCountry(ctx context.Context, c *IPCountry) context.Context {
	return context.WithValue(ctx, keyctxCountry, c)
}

// NewContextWithError creates a new context with an error attached.
func NewContextWithError(ctx context.Context, err error) context.Context {
	return context.WithValue(ctx, keyctxErr, err)
}

// FromContextCountry returns the geoip.Country in ctx if it exists or
// and error if that one exists. The error has been previously set
// by NewContextWithError.
func FromContextCountry(ctx context.Context) (*IPCountry, error, bool) {
	err, ok := ctx.Value(keyctxErr).(error)
	if ok {
		return nil, err, ok
	}
	c, ok := ctx.Value(keyctxCountry).(*IPCountry)
	return c, nil, ok
}
