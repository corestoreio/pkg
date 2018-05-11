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
	"github.com/corestoreio/log"
)

// Options applies configurations to the NewService function. Used mainly by external
// packages for providing different storage engines.
type Options struct {
	// Level1 defines a usually short lived cached like an LRU or TinyLFU. It
	// can be set optionally. Level1 only gets called for Set operation when a
	// value is requested in Level2 via Get.
	Level1       Storager
	Log          log.Logger
	EnablePubSub bool
	// EnvironmentKey loads a string from the environment to use it as a prefix
	// for the path. For example different payment gateways and access
	// credentials for a payment method on STAGING or PRODUCTION systems.
	EnvironmentKey string
}

// LoadDataFn allows other storage backends to pump their data into the
// config.Service during or after initialization.
type LoadDataFn func(*Service) error
