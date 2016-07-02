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

package generic

import (
	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/store/scope"
)

type Service struct {
	service
}

func New(opts ...Option) (*Service, error) {
	return newService(opts...)
}

type scopedConfig struct {
	scopedConfigGeneric
}

func (sc *scopedConfig) isValid() error {
	return nil
}

type Option func(*Service) error

type OptionFactoryFunc func(config.ScopedGetter) []Option

func WithDefaultConfig(scp scope.Scope, id int64) Option {
	return func(s *Service) error {
		return nil
	}
}
