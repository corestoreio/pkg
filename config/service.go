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
	"time"

	"github.com/corestoreio/cspkg/config/cfgpath"
	"github.com/corestoreio/cspkg/util/conv"
	"github.com/corestoreio/errors"
	"github.com/corestoreio/log"
)

// LeftDelim and RightDelim are used withing the core_config_data.value field to
// allow the replacement of the placeholder in exchange with the current value.
const (
	LeftDelim  = "{{"
	RightDelim = "}}"
)

// Getter implements how to receive thread-safe a configuration value from an
// underlying backend service. The provided cfgpath.Path as an argument does not
// make any assumptions if the scope of the cfgpath.Path is allowed to retrieve
// the value. The NewScoped() function binds a cfgpath.Route to a scope.Scope
// and gives you the possibility to fallback the hierarchy levels. If a value
// cannot be found it must return an error of behaviour not NotFound.
type Getter interface {
	NewScoped(websiteID, storeID int64) Scoped
	Byte(cfgpath.Path) ([]byte, error)
	String(cfgpath.Path) (string, error)
	Bool(cfgpath.Path) (bool, error)
	Float64(cfgpath.Path) (float64, error)
	Int(cfgpath.Path) (int, error)
	Time(cfgpath.Path) (time.Time, error)
	Duration(cfgpath.Path) (time.Duration, error)
	// maybe add compare and swap function
}

// GetterPubSuber implements a configuration Getter and a Subscriber for Publish
// and Subscribe pattern.
type GetterPubSuber interface {
	Getter
	Subscriber
}

// Writer thread safe storing of configuration values under different paths and
// scopes.
type Writer interface {
	// Write writes a configuration entry and may return an error
	Write(p cfgpath.Path, value interface{}) error
}

// Storager is the underlying data storage for holding the keys and its values.
// Implementations can be spf13/viper or MySQL backed. Default Storager is a
// simple mutex protected map[string]interface{}. The config.Writer function
// calls the config.Storager functions and Storager must make sure of the
// correct type conversions to the supported type of the underlying storage
// engine.
type Storager interface {
	// Set sets a key with a value and returns on success nil or
	// ErrKeyOverwritten, on failure any other error
	Set(key cfgpath.Path, value interface{}) error
	// Get returns the raw value on success or may return a NotFound error
	// behaviour if an entry cannot be found or does not exists. Any other error
	// can also occur.
	Get(key cfgpath.Path) (interface{}, error)
	// AllKeys returns the fully qualified keys
	AllKeys() (cfgpath.PathSlice, error)
}

// Service main configuration provider. Please use the NewService() function
type Service struct {
	// backend is the underlying data holding provider. Only access it if you
	// know exactly what you are doing.
	backend Storager

	// internal service to provide async pub/sub features while reading/writing
	// config values.
	*pubSub

	// Log can be set for debugging purpose. If nil, it panics. Default
	// log.Blackhole with disabled debug and info logging. You should use the
	// option function WithLogger because the logger gets also set to the
	// internal pub/sub service. The exported Log can be used in external
	// package to log within functional option calls. For example in
	// config/storage/ccd.
	Log log.Logger
}

// NewService creates the main new configuration for all scopes: default,
// website and store. Default Storage is a simple map[string]interface{}. A new
// go routine will be startet for the publish and subscribe feature.
func NewService(backend Storager, opts ...Option) (*Service, error) {
	s := &Service{
		backend: backend,
		Log:     log.BlackHole{}, // disabled debug and info logging.
	}

	if err := s.Options(opts...); err != nil {
		if s.pubSub != nil {
			if err2 := s.Close(); err2 != nil {
				// terminate publisher go routine and prevent leaking
				return nil, errors.Wrap(err2, "[config] Service.Option.Close")
			}
		}
		return nil, errors.Wrap(err, "[config] Service.Option")
	}

	p := cfgpath.MustNewByParts(PathCSBaseURL)
	if err := s.backend.Set(p, CSBaseURL); err != nil {
		if err2 := s.Close(); err2 != nil { // terminate publisher go routine and prevent leaking
			return nil, errors.Wrap(err2, "[config] Service.Storage.Close")
		}
		return nil, errors.Wrap(err, "[config] Service.Option")
	}
	return s, nil
}

// MustNewService same as NewService but panics on error. Use only in testing
// or during boot process.
func MustNewService(backend Storager, opts ...Option) *Service {
	s, err := NewService(backend, opts...)
	if err != nil {
		panic(err)
	}
	return s
}

// Options applies service options.
func (s *Service) Options(opts ...Option) error {
	for _, opt := range opts {
		if opt != nil {
			if err := opt(s); err != nil {
				return errors.Wrap(err, "[config] Service.Options")
			}
		}
	}
	return nil
}

// NewScoped creates a new scope base configuration reader
func (s *Service) NewScoped(websiteID, storeID int64) Scoped {
	return NewScoped(s, websiteID, storeID)
}

// Write puts a value back into the Service. Example usage:
//		// Default Scope
//		p, err := cfgpath.NewByParts("currency/option/base") // or use cfgpath.MustNewByParts( ... )
// 		err := Write(p, "USD")
//
//		// Website Scope
//		// 3 for example comes from core_website/store_website database table
//		err := Write(p.Bind(scope.WebsiteID, 3), "EUR")
//
//		// Store Scope
//		// 6 for example comes from core_store/store database table
//		err := Write(p.Bind(scope.StoreID, 6), "CHF")
func (s *Service) Write(p cfgpath.Path, v interface{}) error {
	if s.Log.IsDebug() {
		s.Log.Debug("config.Service.Write", log.Stringer("path", p), log.Object("val", v))
	}

	if err := s.backend.Set(p, v); err != nil {
		return errors.Wrap(err, "[config] sStorage.Set")
	}
	if s.pubSub != nil {
		s.sendMsg(p)
	}
	return nil
}

// get generic getter ... not sure if this should be public ...
func (s *Service) get(p cfgpath.Path) (interface{}, error) {
	if s.Log.IsDebug() {
		s.Log.Debug("config.Service.get", log.Stringer("path", p))
	}
	return s.backend.Get(p)
}

// String returns a string from the Service. Example usage:
//
//		// Default Scope
//		p, err := cfgpath.NewByParts("general/locale/timezone") // or use cfgpath.MustNewByParts( ... )
// 		s, err := String(p)
//
//		// Website Scope
//		// 3 for example comes from core_website/store_website database table
//		s, err := String(p.Bind(scope.WebsiteID, 3))
//
//		// Store Scope
//		// 6 for example comes from core_store/store database table
//		s, err := String(p.Bind(scope.StoreID, 6))
func (s *Service) String(p cfgpath.Path) (string, error) {
	vs, err := s.get(p)
	if err != nil {
		return "", errors.Wrap(err, "[config] Storage.String.get")
	}
	return conv.ToStringE(vs)
}

// Byte returns a byte slice from the Service. Example usage see String.
func (s *Service) Byte(p cfgpath.Path) ([]byte, error) {
	vs, err := s.get(p)
	if err != nil {
		return nil, errors.Wrap(err, "[config] Storage.Byte.get")
	}
	return conv.ToByteE(vs)
}

// Bool returns bool from the Service. Example usage see String.
func (s *Service) Bool(p cfgpath.Path) (bool, error) {
	vs, err := s.get(p)
	if err != nil {
		return false, errors.Wrap(err, "[config] Storage.Bool.get")
	}
	return conv.ToBoolE(vs)
}

// Float64 returns a float64 from the Service. Example usage see String.
func (s *Service) Float64(p cfgpath.Path) (float64, error) {
	vs, err := s.get(p)
	if err != nil {
		return 0, errors.Wrap(err, "[config] Storage.Float64.get")
	}
	return conv.ToFloat64E(vs)
}

// Int returns an int from the Service. Example usage see String.
func (s *Service) Int(p cfgpath.Path) (int, error) {
	vs, err := s.get(p)
	if err != nil {
		return 0, errors.Wrap(err, "[config] Storage.Int.get")
	}
	return conv.ToIntE(vs)
}

// Time returns a date and time object from the Service. Example usage see
// String. Time() is able to parse available time formats as defined in
// github.com/corestoreio/cspkg/util/conv.StringToDate()
func (s *Service) Time(p cfgpath.Path) (time.Time, error) {
	vs, err := s.get(p)
	if err != nil {
		return time.Time{}, errors.Wrap(err, "[config] Storage.Time.get")
	}
	return conv.ToTimeE(vs)
}

// Duration returns a duration from the Service. Example usage see String.
func (s *Service) Duration(p cfgpath.Path) (time.Duration, error) {
	vs, err := s.get(p)
	if err != nil {
		return 0, errors.Wrap(err, "[config] Storage.Duration.get")
	}
	return conv.ToDurationE(vs)
}

// IsSet checks if a key is in the configuration. Returns false on error. Errors
// will be logged in Debug mode. Does not check if the value can be asserted to
// the desired type.
func (s *Service) IsSet(p cfgpath.Path) bool {
	v, err := s.backend.Get(p)
	if err != nil {
		if s.Log.IsDebug() {
			s.Log.Debug("config.Service.IsSet.Storage.Get", log.Err(err), log.Stringer("path", p))
		}
		return false
	}
	return v != nil
}
