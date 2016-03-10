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

	"github.com/corestoreio/csfw/config/cfgpath"
	"github.com/corestoreio/csfw/config/element"
	"github.com/corestoreio/csfw/config/storage"
	"github.com/corestoreio/csfw/util/conv"
	"github.com/corestoreio/csfw/util/cserr"
)

// LeftDelim and RightDelim are used withing the core_config_data.value field to allow the replacement
// of the placeholder in exchange with the current value.
const (
	LeftDelim  = "{{"
	RightDelim = "}}"
)

type (
	// Getter implements how to receive thread-safe a configuration value from
	// a cfgpath.Path.
	// Providing a cfgpath.Path argument does not make any assumptions if the
	// scope of the cfgpath.Path is allowed to retrieve the value.
	// The NewScoped() function binds a cfgpath.Route to a scope.Scope and
	// gives you the possibility to fallback the hierarchy levels.
	// These functions are also available in the ScopedGetter interface.
	Getter interface {
		NewScoped(websiteID, storeID int64) ScopedGetter
		String(cfgpath.Path) (string, error)
		Bool(cfgpath.Path) (bool, error)
		Float64(cfgpath.Path) (float64, error)
		Int(cfgpath.Path) (int, error)
		Time(cfgpath.Path) (time.Time, error)
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
		Write(p cfgpath.Path, value interface{}) error
	}

	// Service main configuration provider
	Service struct {
		// Storage is the underlying data holding provider. Only access it
		// if you know exactly what you are doing.
		Storage storage.Storager
		// MultiErr which ServiceOption function arguments are generating
		// Usually empty (= nil) ;-)
		*cserr.MultiErr
		*pubSub
	}
)

// DefaultService provides a standard NewService via init() func loaded.
var DefaultService *Service

func init() {
	// todo: think about removing it or explain properly why a global code smell is needed.
	DefaultService = MustNewService()
}

// ServiceOption applies options to the NewService.
type ServiceOption func(*Service)

// NewService creates the main new configuration for all scopes: default, website
// and store. Default Storage is a simple map[string]interface{}. A new go routine
// will be startet for the publish and subscribe feature.
func NewService(opts ...ServiceOption) (*Service, error) {
	s := &Service{
		pubSub:  newPubSub(),
		Storage: storage.NewKV(),
	}

	go s.publish()

	_ = s.Options(opts...)

	if s.HasErrors() {
		if err := s.Close(); err != nil { // terminate publisher go routine and prevent leaking
			s.MultiErr = s.AppendErrors(err)
		}
		return nil, s
	}

	p := cfgpath.MustNewByParts(PathCSBaseURL)
	if err := s.Storage.Set(p, CSBaseURL); err != nil {
		if err := s.Close(); err != nil { // terminate publisher go routine and prevent leaking
			s.MultiErr = s.AppendErrors(err)
		}
		return nil, err
	}
	return s, nil
}

// MustNewService same as NewService but panics on error. Use only in testing
// or during boot process.
func MustNewService(opts ...ServiceOption) *Service {
	s, err := NewService(opts...)
	if err != nil {
		panic(err)
	}
	return s
}

// Options applies service options.
func (s *Service) Options(opts ...ServiceOption) error {
	for _, opt := range opts {
		if opt != nil {
			opt(s)
		}
	}
	if s.HasErrors() {
		return s
	}
	return nil
}

// NewScoped creates a new scope base configuration reader
func (s *Service) NewScoped(websiteID, storeID int64) ScopedGetter {
	return NewScopedService(s, websiteID, storeID)
}

// ApplyDefaults reads slice Sectioner and applies the keys and values to the
// default configuration. Overwrites existing values. TODO maybe use a flag to
// prevent overwriting
func (s *Service) ApplyDefaults(ss element.Sectioner) (count int, err error) {
	def, err := ss.Defaults()
	if err != nil {
		return
	}
	for k, v := range def {
		if PkgLog.IsDebug() {
			PkgLog.Debug("config.Service.ApplyDefaults", k, v)
		}
		var p cfgpath.Path
		p, err = cfgpath.NewByParts(k) // default path!
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
	if PkgLog.IsDebug() {
		PkgLog.Debug("config.Service.Write", "path", p, "val", v)
	}

	if err := s.Storage.Set(p, v); err != nil {
		return err
	}
	s.sendMsg(p)
	return nil
}

// get generic getter ... not sure if this should be public ...
func (s *Service) get(p cfgpath.Path) (interface{}, error) {
	if PkgLog.IsDebug() {
		PkgLog.Debug("config.Service.get", "path", p)
	}
	return s.Storage.Get(p)
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
		return "", err
	}
	return conv.ToStringE(vs)
}

// Bool returns bool from the Service. Example usage see String.
func (s *Service) Bool(p cfgpath.Path) (bool, error) {
	vs, err := s.get(p)
	if err != nil {
		return false, err
	}
	return conv.ToBoolE(vs)
}

// Float64 returns a float64 from the Service. Example usage see String.
func (s *Service) Float64(p cfgpath.Path) (float64, error) {
	vs, err := s.get(p)
	if err != nil {
		return 0, err
	}
	return conv.ToFloat64E(vs)
}

// Int returns an int from the Service. Example usage see String.
func (s *Service) Int(p cfgpath.Path) (int, error) {
	vs, err := s.get(p)
	if err != nil {
		return 0, err
	}
	return conv.ToIntE(vs)
}

// Time returns a date and time object from the Service. Example usage see String.
// Time() is able to parse available time formats as defined in
// github.com/corestoreio/csfw/util/conv.StringToDate()
func (s *Service) Time(p cfgpath.Path) (time.Time, error) {
	vs, err := s.get(p)
	if err != nil {
		return time.Time{}, err
	}
	return conv.ToTimeE(vs)
}

// IsSet checks if a key is in the configuration. Returns false on error.
// Errors will be logged in Debug mode. Does not check if the value can be asserted
// to the desired type.
func (s *Service) IsSet(p cfgpath.Path) bool {
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
	err = cserr.UnwrapMasked(err)
	return err != nil && err != storage.ErrKeyNotFound
}
