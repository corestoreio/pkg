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

import (
	"errors"
	"time"

	"github.com/corestoreio/csfw/config/element"
	"github.com/corestoreio/csfw/config/path"
	"github.com/corestoreio/csfw/storage/csdb"
	"github.com/corestoreio/csfw/storage/dbr"
	"github.com/corestoreio/csfw/store/scope"
	"github.com/corestoreio/csfw/util/cast"
	"github.com/juju/errgo"
)

// LeftDelim and RightDelim are used withing the core_config_data.value field to allow the replacement
// of the placeholder in exchange with the current value.
const (
	LeftDelim  = "{{"
	RightDelim = "}}"
)

type (
	// Getter implements how to receive thread-safe a configuration value from
	// a path and or scope.
	//
	// These functions are also available in the ScopedGetter interface.
	Getter interface {
		NewScoped(websiteID, groupID, storeID int64) ScopedGetter
		String(path.Path) (string, error)
		Bool(path.Path) (bool, error)
		Float64(path.Path) (float64, error)
		Int(path.Path) (int, error)
		DateTime(path.Path) (time.Time, error)
		// maybe add compare and swap function
	}

	// GetterPubSuber implements a configuration Getter and a Subscriber for
	// Publish and Subscribe pattern.
	GetterPubSuber interface {
		Getter
		Subscriber
	}

	// Writer thread safe storing of configuration values under different paths and scopes.
	Writer interface {
		// Write writes a configuration entry and may return an error
		Write(p path.Path, value interface{}) error
	}

	// Service main configuration provider
	Service struct {
		// Storage is the underlying data holding provider. Only access it
		// if you know exactly what you are doing.
		Storage Storager
		*pubSub
	}
)

// TableCollection handles all tables and its columns. init() in generated Go file will set the value.
var TableCollection csdb.Manager

// DefaultService provides a standard NewService via init() func loaded.
var DefaultService *Service

// ErrKeyNotFound will be returned if a key cannot be found or value is nil.
// If you provide your own interface implementation make sure to also return
// ErrKeyNotFound if a key cannot be found.
var ErrKeyNotFound = errors.New("Key not found")

func init() {
	DefaultService = NewService()
}

// ServiceOption applies options to the NewService.
type ServiceOption func(*Service)

// WithDBStorage applies the MySQL storage to a new Service. It
// starts the idle checker of the DBStorage type.
func WithDBStorage(p csdb.Preparer) ServiceOption {
	return func(s *Service) {
		s.Storage = MustNewDBStorage(p).Start()
	}
}

// NewService creates the main new configuration for all scopes: default, website
// and store. Default Storage is a simple map[string]interface{}
func NewService(opts ...ServiceOption) *Service {
	s := &Service{
		pubSub: newPubSub(),
	}
	for _, opt := range opts {
		if opt != nil {
			opt(s)
		}
	}
	if s.Storage == nil {
		s.Storage = newSimpleStorage()
	}
	go s.publish()
	p := path.MustNewByParts(PathCSBaseURL)
	s.Storage.Set(p, CSBaseURL)
	return s
}

// NewScoped creates a new scope base configuration reader
func (s *Service) NewScoped(websiteID, groupID, storeID int64) ScopedGetter {
	return newScopedService(s, websiteID, groupID, storeID)
}

// ApplyDefaults reads slice Sectioner and applies the keys and values to the
// default configuration. Overwrites existing values.
func (s *Service) ApplyDefaults(ss element.Sectioner) (count int, err error) {
	def, err := ss.Defaults()
	if err != nil {
		return
	}
	for k, v := range def {
		if PkgLog.IsDebug() {
			PkgLog.Debug("config.Service.ApplyDefaults", k, v)
		}
		var p path.Path
		p, err = path.NewByParts(k) // default path!
		if err != nil {
			return
		}
		if err = s.Storage.Set(p, v); err != nil {
			return
		}
		count++
	}
	return
}

// ApplyCoreConfigData reads the table core_config_data into the Service and overrides
// existing values. If the column `value` is NULL entry will be ignored. It returns the
// loadedRows which are all rows from the table and the writtenRows which are the applied
// config values where a value is valid.
func (s *Service) ApplyCoreConfigData(dbrSess dbr.SessionRunner) (loadedRows, writtenRows int, err error) {
	var ccd TableCoreConfigDataSlice
	loadedRows, err = csdb.LoadSlice(dbrSess, TableCollection, TableIndexCoreConfigData, &ccd)
	if PkgLog.IsDebug() {
		PkgLog.Debug("config.Service.ApplyCoreConfigData", "rows", loadedRows)
	}
	if err != nil {
		if PkgLog.IsDebug() {
			PkgLog.Debug("config.Service.ApplyCoreConfigData.LoadSlice", "err", err)
		}
		err = errgo.Mask(err)
		return
	}

	for _, cd := range ccd {
		if cd.Value.Valid {
			var p path.Path
			p, err = path.NewByParts(cd.Path)
			if err != nil {
				err = errgo.Mask(err)
				return
			}

			if err = s.Write(p.Bind(scope.FromString(cd.Scope), cd.ScopeID), cd.Value.String); err != nil {
				err = errgo.Mask(err)
				return
			}
			writtenRows++
		}
	}
	return
}

// Write puts a value back into the Service. Example usage:
// 	Default Scope: Write(config.Path("currency", "option", "base"), config.Value("USD"))
// 	Website Scope: Write(config.Path("currency", "option", "base"), config.Value("EUR"), config.ScopeWebsite(w))
// 	Store   Scope: Write(config.Path("currency", "option", "base"), config.ValueReader(resp.Body), config.ScopeStore(s))
func (s *Service) Write(p path.Path, v interface{}) error {
	if PkgLog.IsDebug() {
		PkgLog.Debug("config.Service.Write", "path", p, "val", v)
	}

	s.Storage.Set(p, v)
	s.sendMsg(p)
	return nil
}

// get generic getter ... not sure if this should be public ...
func (s *Service) get(p path.Path) (interface{}, error) {
	if PkgLog.IsDebug() {
		PkgLog.Debug("config.Service.get", "path", p)
	}
	return s.Storage.Get(p)
}

// String returns a string from the Service. Example usage:
//
// Default value: String(config.Path("general/locale/timezone"))
//
// Website value: String(config.Path("general/locale/timezone"), config.ScopeWebsite(w))
//
// Store   value: String(config.Path("general/locale/timezone"), config.ScopeStore(s))
func (s *Service) String(p path.Path) (string, error) {
	vs, err := s.get(p)
	if err != nil {
		return "", err
	}
	return cast.ToStringE(vs)
}

// Bool returns bool from the Service. Example usage see String.
func (s *Service) Bool(p path.Path) (bool, error) {
	vs, err := s.get(p)
	if err != nil {
		return false, err
	}
	return cast.ToBoolE(vs)
}

// Float64 returns a float64 from the Service. Example usage see String.
func (s *Service) Float64(p path.Path) (float64, error) {
	vs, err := s.get(p)
	if err != nil {
		return 0, err
	}
	return cast.ToFloat64E(vs)
}

// Int returns an int from the Service. Example usage see String.
func (s *Service) Int(p path.Path) (int, error) {
	vs, err := s.get(p)
	if err != nil {
		return 0, err
	}
	return cast.ToIntE(vs)
}

// DateTime returns a date and time object from the Service. Example usage see String.
func (s *Service) DateTime(p path.Path) (time.Time, error) {
	vs, err := s.get(p)
	if err != nil {
		return time.Time{}, err
	}
	return cast.ToTimeE(vs)
}

// IsSet checks if a key is in the configuration. Returns false on error.
// Errors will be logged in Debug mode.
func (s *Service) IsSet(p path.Path) bool {
	v, err := s.Storage.Get(p)
	if err != nil {
		if PkgLog.IsDebug() {
			PkgLog.Debug("config.Service.IsSet.Storage.Get", "err", err, "path", p)
		}
		return false
	}
	return v != nil
}

// NotKeyNotFoundError returns true if err is not nil and not of type Key Not Found.
func NotKeyNotFoundError(err error) bool {
	return err != nil && err != ErrKeyNotFound
}
