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

package backendratelimit

import (
	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/config/cfgmodel"
	"github.com/corestoreio/csfw/net/ratelimit"
	"github.com/corestoreio/csfw/util/errors"
	"gopkg.in/throttled/throttled.v2"
)

// ConfigDuration handles the allowed duration values like s,i,h and d.
type ConfigRate struct {
	cfgmodel.Str
}

// NewConfigDuration creates a new duration model with a predefined source slice
// of all allowed values.
func NewConfigDuration(path string, opts ...cfgmodel.Option) ConfigRate {
	return ConfigRate{
		Str: cfgmodel.NewStr(
			path,
			append(opts, cfgmodel.WithSourceByString(
				"s", "Second",
				"i", "Minute",
				"h", "Hour",
				"d", "Day",
			))...,
		),
	}
}

// Get returns a new rate. The requests argument declares the number of requests
// allowed per time period. Invalid duration setting falls back to hourly calculation.
func (md ConfigRate) Get(sg config.ScopedGetter, requests int) (throttled.Rate, error) {
	val, err := md.Str.Get(sg)
	if err != nil {
		return throttled.Rate{}, errors.Wrap(err, "[ratelimit] ConfigDuration.Str.Get")
	}
	return ratelimit.CalculateRate(val, requests), nil
}
