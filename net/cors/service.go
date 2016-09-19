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

//go:generate go run ../internal/scopedservice/main_copy.go "$GOPACKAGE"

package cors

import (
	"github.com/corestoreio/csfw/store/scope"
	"github.com/corestoreio/csfw/util/errors"
)

const methodOptions = "OPTIONS"

// Service describes the CrossOriginResourceSharing which is used to create a
// Container Filter that implements CORS. Cross-origin resource sharing (CORS)
// is a mechanism that allows JavaScript on a web page to make XMLHttpRequests
// to another domain, not the domain the JavaScript originated from.
//
// http://en.wikipedia.org/wiki/Cross-origin_resource_sharing
// http://enable-cors.org/server.html
// http://www.html5rocks.com/en/tutorials/cors/#toc-handling-a-not-so-simple-request
type Service struct {
	service
}

// New creates a new Cors handler with the provided options.
func New(opts ...Option) (*Service, error) {
	s, err := newService(opts...)
	if s != nil {
		s.useWebsite = true
		s.optionAfterApply = func() error {
			s.rwmu.Lock()
			defer s.rwmu.Unlock()

			// propagate the logger to all scopes.
			if s.Log != nil {
				for _, sc := range s.scopeCache {
					if sc.log == nil {
						sc.log = s.Log
					}
				}
			}

			// validate that the applied functional options can only be set for
			// scope website. scope store makes no sense.
			for h := range s.scopeCache {
				if scp, _ := h.Unpack(); scp > scope.Website {
					return errors.NewNotSupportedf(errServiceUnsupportedScope, h)
				}
			}
			return nil
		}
	}
	return s, err
}
