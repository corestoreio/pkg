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

package scopedservice

import (
	"io"

	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/log/logw"
	"github.com/corestoreio/csfw/store/scope"
)

// Service DO NOT USE
type Service struct {
	service
}

// New DO NOT USE
func New(opts ...Option) (*Service, error) {
	return newService(opts...)
}

type scopedConfig struct {
	scopedConfigGeneric
	value string
}

func (sc *scopedConfig) isValid() error {
	if sc.lastErr != nil {
		return sc.lastErr
	}
	return nil
}

func newScopedConfig() *scopedConfig {
	return &scopedConfig{}
}

// Option DO NOT USE
type Option func(*Service) error

// OptionFactoryFunc DO NOT USE
type OptionFactoryFunc func(config.ScopedGetter) []Option

// WithDefaultConfig DO NOT USE
func WithDefaultConfig(scp scope.Scope, id int64) Option {
	return func(s *Service) error {
		return withValue(scp, id, "Hello Default Gophers")(s)
	}
}

func withValue(scp scope.Scope, id int64, val string) Option {
	return func(s *Service) error {
		h := scope.NewHash(scp, id)
		s.rwmu.Lock()
		defer s.rwmu.Unlock()

		sc := s.scopeCache[h]
		if sc == nil {
			sc = optionInheritDefault(s)
		}
		sc.value = val
		sc.scopeHash = h
		s.scopeCache[h] = sc
		return nil
	}
}

// withDebugLogger w must be thread safe
func withDebugLogger(w io.Writer) Option {
	return func(s *Service) error {
		s.rwmu.Lock()
		defer s.rwmu.Unlock()
		s.Log = logw.NewLog(logw.WithWriter(w), logw.WithLevel(logw.LevelDebug))
		return nil
	}
}
