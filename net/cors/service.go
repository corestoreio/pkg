// Copyright (c) 2014 Olivier Poitrey <rs@dailymotion.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is furnished
// to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

//go:generate go run ../generic/main_copy.go "$GOPACKAGE"

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
		s.optionValidation = func() error {
			s.rwmu.RLock()
			defer s.rwmu.RUnlock()
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
