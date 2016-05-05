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

import (
	"time"

	"github.com/corestoreio/csfw/util/log"
	"github.com/rs/xstats"
)

type optionBox struct {
	log    log.Logger
	xstat  xstats.XStater
	genRID RequestIDGenerator
}

// Option contains multiple functional options for middlewares.
type Option func(ob *optionBox)

func newOptionBox(opts ...Option) *optionBox {
	ob := &optionBox{
		log:    log.BlackHole{}, // disabled info and debug logging
		xstat:  nopS{},
		genRID: &RequestIDService{},
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

// SetXStats sets a stats handler to a middleware
func SetXStats(x xstats.XStater) Option {
	return func(ob *optionBox) {
		ob.xstat = x
	}
}

// SetRequestIDGenerator sets a custom request ID generator
func SetRequestIDGenerator(g RequestIDGenerator) Option {
	return func(ob *optionBox) {
		ob.genRID = g
	}
}

type nopS struct{}

var _ xstats.XStater = (*nopS)(nil)

// AddTag implements XStats interface
func (rc nopS) AddTags(tags ...string) {
}

// Gauge implements XStats interface
func (rc nopS) Gauge(stat string, value float64, tags ...string) {
}

// Count implements XStats interface
func (rc nopS) Count(stat string, count float64, tags ...string) {
}

// Histogram implements XStats interface
func (rc nopS) Histogram(stat string, value float64, tags ...string) {
}

// Timing implements xstats interface
func (rc nopS) Timing(stat string, duration time.Duration, tags ...string) {
}
