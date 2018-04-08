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
	"github.com/corestoreio/errors"
	"github.com/corestoreio/log"
	"github.com/corestoreio/pkg/store/scope"
)

// LeftDelim and RightDelim are used withing the core_config_data.value field to
// allow the replacement of the placeholder in exchange with the current value.
// deprecated gets not refactored
const (
	LeftDelim  = "{{"
	RightDelim = "}}"
)

// Getter implements how to receive thread-safe a configuration value from an
// underlying backend service. The provided route as an argument does not
// make any assumptions if the scope of the Path is allowed to retrieve
// the value. The NewScoped() function binds a route to a scope.Scope
// and gives you the possibility to fallback the hierarchy levels. If a value
// cannot be found, it must return false as the 2nd return argument.
type Getter interface {
	NewScoped(websiteID, storeID int64) Scoped
	Value(Path) (v Value, ok bool, err error)
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
	Write(p Path, value []byte) error
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
	Set(scp scope.TypeID, path string, value []byte) error
	// Get returns the raw value `v` on success and sets `ok` to true because
	// the value has been found. If the value cannot be found for the desired
	// path, return value `ok` must be false. A nil value `v` indicates also a
	// value and hence `ok` is true, if found.
	Value(scp scope.TypeID, path string) (v []byte, ok bool, err error)
	// AllKeys returns the fully qualified keys
	AllKeys() (scps scope.TypeIDs, paths []string, err error)
}

// Service main configuration provider. Please use the NewService() function.
// TODO build in to support different environments for the same keys. E.g. different API credentials for external APIs (staging, production, etc).
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

	// TODO auto detect unmarhaller and get one from a pool pass the correct unmarshaler to the Value.
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
//		p, err := cfgpath.MakeByString("currency/option/base") // or use cfgpath.MustMakeByString( ... )
// 		err := Write(p, "USD")
//
//		// Website Scope
//		// 3 for example comes from core_website/store_website database table
//		err := Write(p.Bind(scope.WebsiteID, 3), "EUR")
//
//		// Store Scope
//		// 6 for example comes from core_store/store database table
//		err := Write(p.Bind(scope.StoreID, 6), "CHF")
func (s *Service) Write(p Path, v []byte) error {
	if s.Log.IsDebug() {
		log.WhenDone(s.Log).Debug("config.Service.Write", log.Stringer("path", p), log.Int("data_length", len(v)))
	}
	if err := p.IsValid(); err != nil {
		return errors.WithStack(err)
	}

	if err := s.backend.Set(p.ScopeID, p.route, v); err != nil {
		return errors.Wrap(err, "[config] Service.backend.Set")
	}
	if s.pubSub != nil {
		s.sendMsg(p)
	}
	return nil
}

// String returns a string from the Service. Example usage:
//
//		// Default Scope
//		p, err := cfgpath.MakeByString("general/locale/timezone") // or use cfgpath.MustMakeByString( ... )
// 		s, err := String(p)
//
//		// Website Scope
//		// 3 for example comes from core_website/store_website database table
//		s, err := String(p.Bind(scope.WebsiteID, 3))
//
//		// Store Scope
//		// 6 for example comes from core_store/store database table
//		s, err := String(p.Bind(scope.StoreID, 6))
func (s *Service) Value(p Path) (v Value, ok bool, err error) {
	if s.Log.IsDebug() {
		defer log.WhenDone(s.Log).Debug("config.Service.get", log.Stringer("path", p), log.Bool("found", ok), log.Err(err))
	}

	v.data, ok, err = s.backend.Value(p.ScopeID, p.route)
	return
}

// IsSet checks if a key is in the configuration. Returns false on error. Errors
// will be logged in Debug mode. Does not check if the value can be asserted to
// the desired type.
func (s *Service) IsSet(p Path) bool {
	_, ok, err := s.Value(p)
	return err == nil && ok
}

// think about that segregation
//type ScopedStringer interface {
//  Parent() (scope.Scope, int64)
//	scope.Scoper
//	Bind(scope.Scope) ScopedGetter
//	String(p config.Path, s scope.Scope) (string, error)
//}
// and so on ...

// Scoped is equal to Getter but not an interface and the underlying
// implementation takes care of providing the correct scope: default, website or
// store and bubbling up the scope chain from store -> website -> default if a
// value won't get found in the desired scope. The Path for each primitive type
// represents always a path like "section/group/element" without the scope
// string and scope ID.
//
// To restrict bubbling up you can provide a second argument scope.Scope. You
// can restrict a configuration path to be only used with the default, website
// or store scope. See the examples. This second argument will mainly be used by
// the cfgmodel package to use a defined scope in a config.Structure. If you
// access the ScopedGetter from a store.Store, store.Website type the second
// argument must already be internally pre-filled.
//
// WebsiteID and StoreID must be in a relation like enforced in the database
// tables via foreign keys. Empty storeID triggers the website scope. Empty
// websiteID and empty storeID are triggering the default scope.
//
// You can use the function NewScoped() to create a new object but not
// mandatory. Returned error has mostly the behaviour of NotFound. Debug logging
// can be implemented in the config.Getter.
type Scoped struct {
	// Root holds the main functions for retrieving values by paths from the
	// storage.
	Root      Getter
	WebsiteID int64
	StoreID   int64
}

// NewScopedService instantiates a ScopedGetter implementation.  Getter
// specifies the root Getter which does not know about any scope.
func NewScoped(r Getter, websiteID, storeID int64) Scoped {
	return Scoped{
		Root:      r,
		WebsiteID: websiteID,
		StoreID:   storeID,
	}
}

// IsValid checks if the object has been set up correctly.
func (ss Scoped) IsValid() bool {
	return ss.Root != nil && ((ss.WebsiteID == 0 && ss.StoreID == 0) ||
		(ss.WebsiteID > 0 && ss.StoreID == 0) ||
		(ss.WebsiteID > 0 && ss.StoreID > 0))
}

// ParentID tells you the parent underlying scope and its ID. Store falls back
// to website and website falls back to default.
func (ss Scoped) ParentID() scope.TypeID {
	if ss.StoreID > 0 {
		return scope.Website.Pack(ss.WebsiteID)
	}
	return scope.DefaultTypeID
}

// ScopeID tells you the current underlying scope and its ID to which this
// configuration has been bound to.
func (ss Scoped) ScopeID() scope.TypeID {
	if ss.StoreID > 0 {
		return scope.Store.Pack(ss.StoreID)
	}
	if ss.WebsiteID > 0 {
		return scope.Website.Pack(ss.WebsiteID)
	}
	return scope.DefaultTypeID
}

// ScopeIDs returns the hierarchical order of the scopes containing ScopeID() on
// position zero and ParentID() on position one. This function guarantees that
// the returned slice contains at least two entries.
func (ss Scoped) ScopeIDs() scope.TypeIDs {
	var ids = [2]scope.TypeID{
		ss.ScopeID(),
		ss.ParentID(),
	}
	return ids[:]
}

func (ss Scoped) isAllowedStore(s scope.Type) bool {
	scp := ss.ScopeID().Type()
	if s > scope.Absent {
		scp = s
	}
	return ss.StoreID > 0 && scope.PermStoreReverse.Has(scp)
}

func (ss Scoped) isAllowedWebsite(s scope.Type) bool {
	scp := ss.ScopeID().Type()
	if s > scope.Absent {
		scp = s
	}
	return ss.WebsiteID > 0 && scope.PermWebsiteReverse.Has(scp)
}

// Byte traverses through the scopes store->website->default to find
// a matching byte slice value.
func (ss Scoped) Value(s scope.Type, route string) (v Value, ok bool, err error) {
	// fallback to next parent scope if value does not exists
	p := Path{
		route:   route,
		ScopeID: scope.DefaultTypeID,
	}
	if ss.isAllowedStore(s) {
		v, ok, err := ss.Root.Value(p.BindStore(ss.StoreID))
		if ok || err != nil {
			// value found or err is not a NotFound error
			return v, ok, errors.WithStack(err)
		}
	}
	if ss.isAllowedWebsite(s) {
		v, ok, err := ss.Root.Value(p.BindWebsite(ss.WebsiteID))
		if ok || err != nil {
			return v, ok, errors.WithStack(err)
		}
	}
	return ss.Root.Value(p)
}
