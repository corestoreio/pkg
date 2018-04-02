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

package config

import (
	"github.com/corestoreio/errors"
	"github.com/corestoreio/log"
)

// Option applies options to the NewService function. Used mainly by external
// packages for providing different storage engines.
type Option func(*Service) error

// WithLogger sets a logger to the Service and to the internal pubSub goroutine.
// If nil, everything will panic. Apply this function before setting other
// option functions to provide your logger to those other option functions.
// Default Logger log.Blackhole with disabled debug and info logging.
func WithLogger(l log.Logger) Option {
	return func(s *Service) error {
		s.Log = l
		if s.pubSub != nil {
			s.pubSub.log = l
		}
		return nil
	}
}

// WithPubSub starts the internal publish and subscribe service as a goroutine.
func WithPubSub() Option {
	return func(s *Service) error {
		if s.pubSub != nil && !s.pubSub.closed {
			return errors.AlreadyExists.Newf("[config] PubSub Service already exists and is running.")
		}

		s.pubSub = newPubSub(s.Log)

		go s.publish() // yes we know how to quit this goroutine, just call Service.Close()

		return nil
	}
}
