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

// WithFieldMeta sets immutable default values and scope restrictions into the
// service for specific routes. Storage level and sort order are not supported.
// FieldMeta data gets only set once. Reloading is not possible.
func WithFieldMeta(fms ...*FieldMeta) LoadDataOption {
	var once bool
	return LoadDataOption{
		load: func(s *Service) error {
			if once {
				fms = nil
				return nil
			}
			s.mu.Lock()
			defer func() {
				once = true
				s.mu.Unlock()
			}()

			for _, rfm := range fms {
				if rfm.WriteScopePerm > 0 && rfm.ScopeID > scope.DefaultTypeID {
					return errors.NotAcceptable.Newf("[config] WriteScopePerm %q and ScopeID %q cannot be set at once for path %q", rfm.WriteScopePerm.String(), rfm.ScopeID.String(), rfm.Route)
				}
				rfm.valid = true
				if !rfm.DefaultValid && rfm.Default != "" {
					rfm.DefaultValid = true
				}
				s.routeConfig.PutMeta(rfm.Route, rfm)
			}
			return nil
		},
	}
}

// WithApplySections sets the default values and permissionable scopes for
// specific routes. This function option cannot handle a default value for a
// specific website/store scope. Storage level and sort order are not
// supported.
func WithApplySections(sections ...*Section) LoadDataOption {
	secs := Sections(sections)
	var once bool
	return LoadDataOption{
		load: func(s *Service) error {
			if once {
				sections = nil
				secs = nil
				return nil
			}
			if err := secs.Validate(); err != nil {
				return errors.WithStack(err)
			}

			buf := bufferpool.Get()
			s.mu.Lock()
			defer func() {
				s.mu.Unlock()
				bufferpool.Put(buf)
				once = true
			}()

			fm := new(FieldMeta)
			for _, sec := range secs {
				for _, g := range sec.Groups {
					for _, f := range g.Fields {

						parts := [3]string{sec.ID, g.ID, f.ID}
						joinParts(buf, parts[:]...)
						route := buf.String()

						fm.valid = true
						fm.WriteScopePerm = f.Scopes
						fm.Default = f.Default
						fm.DefaultValid = f.Default != ""
						s.routeConfig.PutMeta(route, fm)
						buf.Reset()
					}
				}
			}
			return nil
		},
	}
}
