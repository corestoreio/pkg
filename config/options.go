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

package config

import "github.com/corestoreio/csfw/log"

// ServiceOption applies options to the NewService.
type ServiceOption func(*Service)

// WithLogger sets a logger to the Service and to the internal pubSub
// goroutine. If nil, everything will panic.
// Apply this function before setting other option functions to provide your
// logger to those other option functions.
// Default Logger log.Blackhole with disabled debug and info logging.
func WithLogger(l log.Logger) ServiceOption {
	return func(s *Service) {
		s.Log = l
		if s.pubSub != nil {
			s.pubSub.log = l
		}
	}
}
