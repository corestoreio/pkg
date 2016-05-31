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

package mw

import "github.com/corestoreio/csfw/log"

type optionBox struct {
	log                   log.Logger
	genRID                RequestIDGenerator
	methodOverrideFormKey string
}

// Option contains multiple functional options for middlewares.
type Option func(ob *optionBox)

func newOptionBox(opts ...Option) *optionBox {
	ob := &optionBox{
		log:                   log.BlackHole{}, // disabled info and debug logging
		genRID:                &requestIDService{},
		methodOverrideFormKey: MethodOverrideFormKey,
	}
	for _, o := range opts {
		if o != nil {
			o(ob)
		}
	}
	return ob
}

// SetLogger sets a logger to a middleware
func SetLogger(l log.Logger) Option {
	return func(ob *optionBox) {
		ob.log = l
	}
}

// SetRequestIDGenerator sets a custom request ID generator
func SetRequestIDGenerator(g RequestIDGenerator) Option {
	return func(ob *optionBox) {
		ob.genRID = g
	}
}

// SetMethodOverrideFormKey sets a custom form input name
func SetMethodOverrideFormKey(k string) Option {
	return func(ob *optionBox) {
		ob.methodOverrideFormKey = k
	}
}
