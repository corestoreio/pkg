// Copyright 2015, Cyrill @ Schumacher.fm and the CoreStore contributors
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

	"github.com/corestoreio/csfw/config/scope"
	"github.com/corestoreio/csfw/storage/csdb"
	"github.com/corestoreio/csfw/storage/dbr"
	"github.com/corestoreio/csfw/utils/cast"
	"github.com/juju/errgo"
	"github.com/spf13/viper"
)

// LeftDelim and RightDelim are used withing the core_config_data.value field to allow the replacement
// of the placeholder in exchange with the current value.
const (
	LeftDelim  = "{{"
	RightDelim = "}}"
)

type (
	// Reader implements how to receive thread-safe a configuration value from
	// a path and or scope.
	//
	// These functions are also available in the ScopedReader interface.
	Reader interface {
		NewScoped(websiteID, groupID, storeID int64) ScopedReader
		GetString(...ArgFunc) (string, error)
		GetBool(...ArgFunc) (bool, error)
		GetFloat64(...ArgFunc) (float64, error)
		GetInt(...ArgFunc) (int, error)
		GetDateTime(...ArgFunc) (time.Time, error)
	}

	// ReaderPubSuber implements a configuration Reader and a Subscriber for
	// Publish and Subscribe pattern.
	ReaderPubSuber interface {
		Reader
		Subscriber
	}

	// Writer thread safe storing of configuration values under different paths and scopes.
	Writer interface {
		// Write writes a configuration entry and may return an error
		Write(...ArgFunc) error
	}

	// Manager main configuration provider
	Manager struct {
		// why is Viper private? Because it can maybe replaced by something else ...
		v *viper.Viper
		*pubSub
	}
)

var (
	_ Reader     = (*Manager)(nil)
	_ Writer     = (*Manager)(nil)
	_ Subscriber = (*Manager)(nil)
)

// TableCollection handles all tables and its columns. init() in generated Go file will set the value.
var TableCollection csdb.Manager

// DefaultManager provides a default manager
var DefaultManager *Manager = NewManager()

// ErrKeyNotFound will be returned if a key cannot be found or value is nil.
var ErrKeyNotFound = errors.New("Key not found")

func init() {
	DefaultManager = NewManager()
}

// NewManager creates the main new configuration for all scopes: default, website and store
func NewManager() *Manager {
	m := &Manager{
		v:      viper.New(),
		pubSub: newPubSub(),
	}
	m.v.Set(mustNewArg(Path(PathCSBaseURL)).scopePath(), CSBaseURL)
	go m.publish()
	return m
}

// NewScoped creates a new scope base configuration reader
func (m *Manager) NewScoped(websiteID, groupID, storeID int64) ScopedReader {
	return newScopedManager(m, websiteID, groupID, storeID)
}

// ApplyDefaults reads the map and applies the keys and values to the default configuration
func (m *Manager) ApplyDefaults(ss Sectioner) *Manager {
	for k, v := range ss.Defaults() {
		if PkgLog.IsDebug() {
			PkgLog.Debug("config.Manager.ApplyDefaults", k, v)
		}
		m.v.Set(k, v)
	}
	return m
}

// ApplyCoreConfigData reads the table core_config_data into the Manager and overrides
// existing values. If the column value is NULL entry will be ignored.
func (m *Manager) ApplyCoreConfigData(dbrSess dbr.SessionRunner) error {
	var ccd TableCoreConfigDataSlice
	rows, err := csdb.LoadSlice(dbrSess, TableCollection, TableIndexCoreConfigData, &ccd)
	if PkgLog.IsDebug() {
		PkgLog.Debug("config.Manager.ApplyCoreConfigData", "rows", rows)
	}
	if err != nil {
		if PkgLog.IsDebug() {
			PkgLog.Debug("config.Manager.ApplyCoreConfigData.LoadSlice", "err", err)
		}
		return errgo.Mask(err)
	}

	for _, cd := range ccd {
		if cd.Value.Valid {
			// scope.ID(cd.ScopeID) because cd.ScopeID is a struct field and cannot satisfy interface scope.IDer
			if err := m.Write(Path(cd.Path), Scope(scope.FromString(cd.Scope), cd.ScopeID)); err != nil {
				return errgo.Mask(err)
			}
		}
	}
	return nil
}

// Write puts a value back into the manager. Example usage:
// Default Scope: Write(config.Path("currency", "option", "base"), config.Value("USD"))
// Website Scope: Write(config.Path("currency", "option", "base"), config.Value("EUR"), config.ScopeWebsite(w))
// Store   Scope: Write(config.Path("currency", "option", "base"), config.ValueReader(resp.Body), config.ScopeStore(s))
func (m *Manager) Write(o ...ArgFunc) error {
	a, err := newArg(o...)
	if err != nil {
		if PkgLog.IsDebug() {
			PkgLog.Debug("config.Manager.Write.newArg", "err", err)
		}
		return errgo.Mask(err)
	}

	if PkgLog.IsDebug() {
		PkgLog.Debug("config.Manager.Write", "path", a.scopePath(), "val", a.v)
	}

	m.v.Set(a.scopePath(), a.v)
	m.sendMsg(a)
	return nil
}

// get generic getter ... not sure if this should be public ...
func (m *Manager) get(o ...ArgFunc) interface{} {
	a, err := newArg(o...)
	if err != nil {
		if PkgLog.IsDebug() {
			PkgLog.Debug("config.Manager.get.newArg", "err", err)
		}
		return errgo.Mask(err)
	}

	return m.v.Get(a.scopePath())
}

// GetString returns a string from the manager. Example usage:
// Default value: GetString(config.Path("general/locale/timezone"))
// Website value: GetString(config.Path("general/locale/timezone"), config.ScopeWebsite(w))
// Store   value: GetString(config.Path("general/locale/timezone"), config.ScopeStore(s))
func (m *Manager) GetString(o ...ArgFunc) (string, error) {
	vs := m.get(o...)
	if vs == nil {
		return "", ErrKeyNotFound
	}
	return cast.ToStringE(vs)
}

// GetBool returns bool from the manager. Example usage see GetString.
func (m *Manager) GetBool(o ...ArgFunc) (bool, error) {
	vs := m.get(o...)
	if vs == nil {
		return false, ErrKeyNotFound
	}
	return cast.ToBoolE(vs)
}

// GetFloat64 returns a float64 from the manager. Example usage see GetString.
func (m *Manager) GetFloat64(o ...ArgFunc) (float64, error) {
	vs := m.get(o...)
	if vs == nil {
		return 0.0, ErrKeyNotFound
	}
	return cast.ToFloat64E(vs)
}

// GetInt returns an int from the manager. Example usage see GetString.
func (m *Manager) GetInt(o ...ArgFunc) (int, error) {
	vs := m.get(o...)
	if vs == nil {
		return 0, ErrKeyNotFound
	}
	return cast.ToIntE(vs)
}

// GetDateTime returns a date and time object from the manager. Example usage see GetString.
func (m *Manager) GetDateTime(o ...ArgFunc) (time.Time, error) {
	vs := m.get(o...)
	if vs == nil {
		return time.Time{}, ErrKeyNotFound
	}
	return cast.ToTimeE(vs)
}

// @todo consider adding other Get* from the viper package

// GetStringSlice returns a slice of strings with config values.
// @todo use the backend model of a config value. most/all magento string slices are comma lists.
func (m *Manager) GetStringSlice(o ...ArgFunc) ([]string, error) {
	return nil, ErrKeyNotFound
	//	return m.v.GetStringSlice(newArg(o...))
}

// AllKeys return all keys regardless where they are set
func (m *Manager) AllKeys() []string { return m.v.AllKeys() }

// IsSet checks if a key is in the configuration. Returns false on error.
func (m *Manager) IsSet(o ...ArgFunc) bool {
	a, err := newArg(o...)
	if err != nil {
		if PkgLog.IsDebug() {
			PkgLog.Debug("config.Manager.IsSet.newArg", "err", err)
		}
		return false
	}
	return m.v.IsSet(a.scopePath())
}

// NotKeyNotFoundError returns true if err is not nil and not of type Key Not Found.
func NotKeyNotFoundError(err error) bool {
	return err != nil && err != ErrKeyNotFound
}
