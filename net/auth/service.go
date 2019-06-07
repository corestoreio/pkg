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

//go:generate go run ../internal/scopedservice/main_copy.go "$GOPACKAGE"

package auth

import "github.com/corestoreio/pkg/config"

// Service implements authentication middleware and scoped based authorization.
type Service struct {
	service
}

// New creates a new authentication service to be used as a middleware or
// standalone.
func New(cfg config.Scoper, opts ...Option) (*Service, error) {
	s, err := newService(cfg, opts...)
	if err != nil {
		return nil, err
	}
	return s, nil
}
