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
	"os"
	"strings"
	"unicode"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/log"
	"github.com/corestoreio/pkg/store/scope"
	"github.com/corestoreio/pkg/util/bufferpool"
)

// EnvNamePlaceHolder replaces in a file name or pattern argument applied to a
// WithFiles, WithFile or WithGlob function with the current environment name of
// *config.Service. You can load environment dependent configuration files.
const EnvNamePlaceHolder = `{CS_ENV}`

// DefaultOSEnvVariableName default name of the OS environment variable.
const DefaultOSEnvVariableName = `CS_ENV`

// Options applies configurations to the NewService function. Used mainly by external
// packages for providing different storage engines.
type Options struct {
	// Level1 defines a usually short lived cached like an LRU or TinyLFU or a
	// storage with un-resetble data (locked data). It can be set optionally.
	// Level1 only gets called for Set operation when a value is requested in
	// Level2 via Get.
	Level1       Storager
	Log          log.Logger
	EnablePubSub bool
	// OSEnvVariableName loads a string from an applied environment variable to
	// use it as a prefix for the Path type and when loading configuration files
	// as part of their filename or path (see cfgfile.EnvNamePlaceHolder). For
	// example different payment gateways and access credentials for a payment
	// method on STAGING or PRODUCTION systems.
	OSEnvVariableName string
	// EnvName defines the current name of the environment directly instead of
	// loading it from an operating system environment variable. It gets used as
	// a prefix for the Path type or as part of the directory or filename when
	// loading configuration files (see cfgfile.EnvNamePlaceHolder).
	EnvName string

	// EnableHotReload if the Service receives an OS signal, it triggers a hot
	// reload of the cached functions of type LoadDataOption. Errors during hot
	// reloading do not trigger an exit of the config.Service.
	EnableHotReload bool
	// HotReloadSignals specifies custom signals to listen to. Defaults to
	// syscall.SIGUSR2
	HotReloadSignals []os.Signal
}

// LoadDataOption allows other storage backends to pump their data into the
// config.Service during or after initialization via an OS signal and hot
// reloading.
type LoadDataOption struct {
	level     int // either 1 or 2, if other value, falls back to 2.
	sortOrder int
	load      func(*Service) error
}

// MakeLoadDataOption a wrapper helper function.
func MakeLoadDataOption(fn func(*Service) error) LoadDataOption {
	return LoadDataOption{
		load: fn,
	}
}

// WithSortOrder sets the order/position while loading the data in
// config.Service.
func (o LoadDataOption) WithSortOrder(sortOrderPosition int) LoadDataOption {
	o.sortOrder = sortOrderPosition
	return o
}

// WithUseStorageLevel executes the current Load function in storage level one.
// Argument level can either be 1 or 2. Any other integer value falls back to 2,
// for now.
func (o LoadDataOption) WithUseStorageLevel(level int) LoadDataOption {
	o.level = level
	return o
}

type loadDataOptions []LoadDataOption

func (o loadDataOptions) Len() int           { return len(o) }
func (o loadDataOptions) Less(i, j int) bool { return o[i].sortOrder < o[j].sortOrder }
func (o loadDataOptions) Swap(i, j int)      { o[i], o[j] = o[j], o[i] }

func isLetter(str string) bool {
	for _, r := range str {
		if !unicode.IsLetter(r) {
			return false
		}
	}
	return true
}

func isDigitOnly(str string) bool {
	for _, r := range str {
		if !unicode.IsDigit(r) {
			return false
		}
	}
	return true
}

func hasScopePrefix(path string) bool {
	firstSepPos := strings.IndexByte(path, PathSeparator)
	if firstSepPos < 1 { // case: /aa/bb/cc OR no PathSeparator at all
		return false
	}
	return scope.Valid(path[:firstSepPos])
}

// WithApplyDefaultValues sets immutable default values and scope restrictions
// into the service for specific routes. `routeValueScope` is a triple balanced
// slice, e.g.:
//		WithApplyDefaults(
// 			"currency/options/allow","d","CHF,EUR",
// 			"tax/classes/shipping_tax_class","w","2",
// 			"carriers/freeshipping/name","s","Free Shipping",
// 		)
// First part the route, then default value and last part the scope permission
// where "d" or "" defines the default scope, "w" defines websites and "s"
// defines stores.
func WithApplyDefaultValues(routeValueScope ...string) LoadDataOption {
	return LoadDataOption{
		load: func(s *Service) error {
			if len(routeValueScope)%3 != 0 {
				return errors.NotValid.Newf("[config] WithApplyDefaultValues: routeValueScope is not a triple balanced slice: %v", routeValueScope)
			}

			s.mu.Lock()
			defer s.mu.Unlock()

			for i := 0; i < len(routeValueScope); i = i + 3 {
				r := routeValueScope[i]
				if hasScopePrefix(r) {
					return errors.NotAllowed.Newf("[config] WithApplyDefaults path %q should be qualified with a scope. Expecting just e.g.: aa/bb/cc", r)
				}

				perm, err := scope.MakePerm(routeValueScope[i+1])
				if err != nil {
					return errors.Wrapf(err, "[config] WithApplyDefaultValues for route %q", r)
				}

				s.routeConfig.PutMeta(r, perm, routeValueScope[i+2])
			}
			return nil
		},
	}
}

// WithApplySections sets the default values and permissionable scopes for
// specific routes.
func WithApplySections(sections ...*Section) LoadDataOption {
	secs := Sections(sections)
	return LoadDataOption{
		load: func(s *Service) error {
			if err := secs.Validate(); err != nil {
				return errors.WithStack(err)
			}

			buf := bufferpool.Get()
			s.mu.Lock()
			defer func() {
				s.mu.Unlock()
				bufferpool.Put(buf)
			}()

			for _, sec := range secs {
				for _, g := range sec.Groups {
					for _, f := range g.Fields {

						parts := [3]string{sec.ID, g.ID, f.ID}
						joinParts(buf, parts[:]...)
						route := buf.String()

						s.routeConfig.PutMeta(route, f.ScopePerm, f.Default)

						buf.Reset()
					}
				}
			}

			return nil
		},
	}
}
