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
	"os/signal"
	"strings"
	"syscall"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/log"
	"github.com/corestoreio/pkg/store/scope"
)

// Getter implements how to receive thread-safe a configuration value from an
// underlying backend service. The provided route as an argument does not
// make any assumptions if the scope of the Path is allowed to retrieve
// the value. The NewScoped() function binds a route to a scope.Scope
// and gives you the possibility to fallback the hierarchy levels. If a value
// cannot be found, it must return false as the 2nd return argument.
type Getter interface {
	NewScoped(websiteID, storeID int64) Scoped
	// Value returns a guaranteed non-nil Value.
	Get(p *Path) *Value
}

// GetterPubSuber implements a configuration Getter and a Subscriber for Publish
// and Subscribe pattern.
type GetterPubSuber interface {
	Getter
	Subscriber
}

// Setter thread safe storing of configuration values under different paths and
// scopes.
type Setter interface {
	// Write writes a configuration entry and may return an error
	Set(p *Path, value []byte) error
}

// Storager is the underlying data storage for holding the keys and its values.
// Implementations can be spf13/viper or MySQL backed. Default Storager is a
// simple mutex protected in memory map[string]string.
// Storager must be safe for concurrent use.
type Storager interface {
	// TODO maybe add context as first param and observe in storage implementations if the ctx got cancelled.
	// Set sets a key with a value and returns on success nil or
	// ErrKeyOverwritten, on failure any other error
	Setter
	// Get returns the raw value `v` on success and sets `found` to true
	// because the value has been found. If the value cannot be found for the
	// desired path, return value `found` must be false. A nil value `v`
	// indicates also a value and hence `found` is true, if found.
	Get(p *Path) (v []byte, found bool, err error)
}

// Service main configuration provider. Please use the NewService() function.
// Safe for concurrent use.
// TODO build in to support different environments for the same keys. E.g. different API credentials for external APIs (staging, production, etc).
type Service struct {
	level2  Storager
	envName string
	config  Options
	Log     log.Logger
	// internal service to provide async pub/sub features while reading/writing
	// config values.
	pubSub          *pubSub
	hotReloadSignal chan os.Signal
	loadDataFns     []LoadDataFn
	envReplacer     *strings.Replacer
	// TODO auto detect unmarshaller and get one from a pool pass the correct unmarshaler to the Value.
}

// NewService creates the main new configuration for all scopes: default,
// website and store. Default Storage is a simple map[string]interface{}. A new
// go routine will be startet for the publish and subscribe feature.
// Level2 is the required underlying data holding provider.
func NewService(level2 Storager, o Options, fns ...LoadDataFn) (s *Service, err error) {

	s = &Service{
		level2: level2,
		config: o,
		Log:    o.Log,
	}

	if err := s.setupEnv(); err != nil {
		return nil, errors.WithStack(err)
	}

	if o.EnablePubSub {
		var l log.Logger
		if o.Log != nil {
			l = o.Log.With(log.Bool("pubSub", true))
		}
		s.pubSub = newPubSub(l)
		go s.pubSub.publish() // yes we know how to quit this goroutine, just call Service.Close()
	}

	s.loadDataFns = append(s.loadDataFns, fns...) // make a copy of fns slice
	if err := s.loadData(); err != nil {
		if err2 := s.Close(); err2 != nil {
			// terminate publisher go routine and prevent leaking
			return nil, errors.WithStack(err2)
		}
		return nil, errors.WithStack(err)
	}

	if err := s.shouldEnableHotReload(); err != nil {
		return nil, errors.WithStack(err)
	}

	return s, nil
}

// MustNewService same as NewService but panics on error. Use only in testing
// or during boot process.
func MustNewService(level2 Storager, o Options, fns ...LoadDataFn) *Service {
	s, err := NewService(level2, o, fns...)
	if err != nil {
		panic(err)
	}
	return s
}

func (s *Service) setupEnv() error {
	osEnvVar := DefaultOSEnvVariableName
	if s.config.OSEnvVariableName != "" {
		osEnvVar = s.config.OSEnvVariableName
	}
	if envVal := os.Getenv(osEnvVar); envVal != "" {
		if !isLetter(envVal) {
			return errors.NotValid.Newf("[config] Environment key %q contains invalid non-letter characters: %q", osEnvVar, envVal)
		}
		s.envName = envVal
	}

	if s.envName == "" {
		s.envName = s.config.EnvName
	}
	if s.envName != "" {
		s.envReplacer = strings.NewReplacer(EnvNamePlaceHolder, s.EnvName())
	}
	return nil
}

func (s *Service) shouldEnableHotReload() error {
	if !s.config.EnableHotReload {
		return nil
	}

	s.hotReloadSignal = make(chan os.Signal, 1)
	signals := []os.Signal{syscall.SIGUSR2}
	if len(s.config.HotReloadSignals) > 0 {
		signals = append(signals[:0], s.config.HotReloadSignals...)
	}

	signal.Notify(s.hotReloadSignal, signals...)

	go func() {
		for sgnl := range s.hotReloadSignal {
			if s.config.Log != nil && s.config.Log.IsDebug() {
				s.config.Log.Debug("config.Service.HotReload.Signal", log.String("signal", sgnl.String()))
			}
			err := s.loadData()
			if s.config.Log != nil && s.config.Log.IsInfo() && err != nil {
				s.config.Log.Debug("config.Service.HotReload.LoadingError", log.String("signal", sgnl.String()), log.Err(err))
			}
		}
	}()

	return nil
}

// loadData used for hot reloading and runs also within another goroutine but
// reads only from *Service.
func (s *Service) loadData() error {
	for _, fn := range s.loadDataFns {
		if err := fn(s); err != nil {
			return errors.WithStack(err)
		}
	}
	return nil
}

// EnvName returns the environment name to which this service is bound to.
func (s *Service) EnvName() string {
	return s.envName
}

// ReplaceEnvName replaces the occurrence of the value of constant
// EnvNamePlaceHolder with the current configured environment name. If no
// environment name has been defined it returns the argument unchanged.
func (s *Service) ReplaceEnvName(str string) string {
	if s.envReplacer == nil {
		return str
	}
	return s.envReplacer.Replace(str)
}

// Close closes and terminates the internal goroutines and connections.
func (s *Service) Close() error {

	if s.config.EnableHotReload {
		signal.Stop(s.hotReloadSignal)
		close(s.hotReloadSignal)
	}

	if s.config.EnablePubSub {
		if err := s.pubSub.Close(); err != nil {
			return errors.WithStack(err)
		}
	}
	return nil
}

// Flush flushes the internal caches. Write operation also flushes the entry for
// the given key. If a Storage service implements
//		type flusher interface {
//			Flush() error
//		}
// then it can flush its internal caches.
func (s *Service) Flush() error {
	type flusher interface {
		Flush() error
	}
	if s.config.Level1 != nil {
		if f, ok := s.config.Level1.(flusher); ok {
			if err := f.Flush(); err != nil {
				return errors.WithStack(err)
			}
		}
	}
	if f, ok := s.level2.(flusher); ok {
		if err := f.Flush(); err != nil {
			return errors.WithStack(err)
		}
	}
	return nil
}

// NewScoped creates a new scope base configuration reader
func (s *Service) NewScoped(websiteID, storeID int64) Scoped {
	return NewScoped(s, websiteID, storeID)
}

// Set puts a value back into the Service. Safe for concurrent use. Example
// usage:
//		// Default Scope
//		p, err := config.MakeByString("currency/option/base") // or use config.MustMakeByString( ... )
// 		err := Write(p, "USD")
//
//		// Website Scope
//		// 3 for example comes from core_website/store_website database table
//		err := Write(p.Bind(scope.WebsiteID, 3), "EUR")
//
//		// Store Scope
//		// 6 for example comes from core_store/store database table
//		err := Write(p.Bind(scope.StoreID, 6), "CHF")
func (s *Service) Set(p *Path, v []byte) error {
	if s.config.Log != nil && s.config.Log.IsDebug() {
		log.WhenDone(s.config.Log).Debug("config.Service.Write", log.Stringer("path", p), log.Int("data_length", len(v)))
	}
	if err := p.IsValid(); err != nil {
		return errors.WithStack(err)
	}
	if err := s.level2.Set(p, v); err != nil {
		return errors.Wrap(err, "[config] Service.level2.Set")
	}
	if s.pubSub != nil {
		s.pubSub.sendMsg(*p)
	}
	return nil
}

// Get returns a configuration value from the Service, ignoring the scopes
// using a direct match. Safe for concurrent use. Example usage:
//
//		// Default Scope
//		dp := config.MustNewPath("general/locale/timezone")
//
//		// Website Scope
//		// 3 for example comes from store_website database table
//		ws := dp.Bind(scope.WebsiteID, 3)
//
//		// Store Scope
//		// 6 for example comes from store database table
//		ss := p.Bind(scope.StoreID, 6)
//
// Returns a guaranteed non-nil value.
func (s *Service) Get(p *Path) (v *Value) {
	if s.config.Log != nil && s.config.Log.IsDebug() {
		wdl := log.WhenDone(s.config.Log)
		defer func() {
			// this func captures the scope of v or the structured log entries would be Go's default values.
			wdl.Debug("config.Service.Get", log.Stringer("path", p), log.String("found", valFoundStringer(v.found)), log.Err(v.lastErr))
		}()
	}

	var ok bool
	v = &Value{
		Path: *p,
	}

	if s.config.Level1 != nil {
		v.data, ok, v.lastErr = s.config.Level1.Get(p)
		if v.lastErr != nil {
			return
		}
		if ok {
			v.found = valFoundL1
			return
		}
	}

	v.data, ok, v.lastErr = s.level2.Get(p)
	if ok {
		v.found = valFoundL2
	}
	if v.lastErr != nil {
		v.lastErr = errors.Wrapf(v.lastErr, "[config] Service.Value with path %q", p)
	}
	if ok && s.config.Level1 != nil && v.lastErr == nil {
		if v.lastErr = s.config.Level1.Set(p, v.data); v.lastErr != nil {
			v.lastErr = errors.Wrapf(v.lastErr, "[config] Service.Level1.Set with path %q", p)
			return
		}
	}
	return v
}

// Subscribe adds a Subscriber to be called when a write event happens. See
// interface Subscriber for a detailed description. Route can be any kind of
// level and can contain StrScope and Scope ID. Valid routes can be for example:
//		- StrScope/ID/currency/options/base
//		- StrScope/ID/currency/options
//		- StrScope/ID/currency
//		- currency/options/base
//		- currency/options
//		- currency
func (s *Service) Subscribe(path string, mr MessageReceiver) (subscriptionID int, err error) {
	if s.pubSub == nil {
		return 0, errors.NotImplemented.Newf("[config] PubSub not enabled")
	}
	return s.pubSub.Subscribe(path, mr)
}

// Unsubscribe removes a subscriber with a specific ID.
func (s *Service) Unsubscribe(subscriptionID int) error {
	if s.pubSub == nil {
		return nil
	}
	return s.pubSub.Unsubscribe(subscriptionID)
}

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
// mandatory. Scoped must act as non-pointer value.
type Scoped struct {
	// Root holds the main functions for retrieving values by paths from the
	// storage.
	Root      Getter
	WebsiteID int64
	StoreID   int64
}

// TODO: Scoped should support websites/0/ and stores/0/ to provide a top level
// websites or stores specific configuration.

// NewScoped instantiates a ScopedGetter implementation.  Getter
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
		return scope.Website.WithID(ss.WebsiteID)
	}
	return scope.DefaultTypeID
}

// ScopeID tells you the current underlying scope and its ID to which this
// configuration has been bound to.
func (ss Scoped) ScopeID() scope.TypeID {
	if ss.StoreID > 0 {
		return scope.Store.WithID(ss.StoreID)
	}
	if ss.WebsiteID > 0 {
		return scope.Website.WithID(ss.WebsiteID)
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

func (ss Scoped) isAllowedStore(restrictUpTo scope.Type) bool {
	scp := ss.ScopeID().Type()
	if restrictUpTo > scope.Absent {
		scp = restrictUpTo
	}
	return ss.StoreID > 0 && scope.PermStoreReverse.Has(scp)
}

func (ss Scoped) isAllowedWebsite(restrictUpTo scope.Type) bool {
	scp := ss.ScopeID().Type()
	if restrictUpTo > scope.Absent {
		scp = restrictUpTo
	}
	return ss.WebsiteID > 0 && scope.PermWebsiteReverse.Has(scp)
}

// Get traverses through the scopes store->website->default to find a matching
// byte slice value. The argument `restrictUpTo` scope.Type restricts the
// bubbling. For example a path gets stored in all three scopes but argument
// `restrictUpTo` specifies only website scope, then the store scope will be
// ignored for querying. If argument `restrictUpTo` has been set to zero aka.
// scope.Absent, then all three scopes are considered for querying.
// Returns a guaranteed non-nil Value.
func (ss Scoped) Get(restrictUpTo scope.Type, route string) (v *Value) {
	// fallback to next parent scope if value does not exists
	p := Path{
		route: route,
	}
	if ss.isAllowedStore(restrictUpTo) {
		p.ScopeID = scope.Store.WithID(ss.StoreID)
		v := ss.Root.Get(&p)
		if v.found > valFoundNo || v.lastErr != nil {
			// value found or err is not a NotFound error
			if v.lastErr != nil {
				v.lastErr = errors.WithStack(v.lastErr) // hmm, maybe can be removed if no one gets confused
			}
			return v
		}
	}
	if ss.isAllowedWebsite(restrictUpTo) {
		p.ScopeID = scope.Website.WithID(ss.WebsiteID)
		v := ss.Root.Get(&p)
		if v.found > valFoundNo || v.lastErr != nil {
			if v.lastErr != nil {
				v.lastErr = errors.WithStack(v.lastErr) // hmm, maybe can be removed if no one gets confused
			}
			return v
		}
	}
	p.ScopeID = scope.DefaultTypeID
	return ss.Root.Get(&p)
}
