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
	"github.com/corestoreio/csfw/utils/log"
	"github.com/spf13/viper"
)

// LeftDelim and RightDelim are used withing the core_config_data.value field to allow the replacement
// of the placeholder in exchange with the current value.
const (
	LeftDelim  = "{{"
	RightDelim = "}}"
)

// PathCSBaseURL main CoreStore base URL, used if no configuration on a store level can be found.
const (
	PathCSBaseURL = "web/corestore/base_url"
	CSBaseURL     = "http://localhost:9500/"
)

// URL* defines the types of available URLs.
const (
	URLTypeAbsent URLType = iota
	// URLTypeWeb defines the URL type to generate the main base URL.
	URLTypeWeb
	// URLTypeStatic defines the URL to the static assets like CSS, JS or theme images
	URLTypeStatic

	// UrlTypeLink hmmm
	// UrlTypeLink

	// URLTypeMedia defines the URL type for generating URLs to product photos
	URLTypeMedia
)

type (
	// URLType defines the type of the URL. Used in constant declaration.
	// @see https://github.com/magento/magento2/blob/0.74.0-beta7/lib/internal/Magento/Framework/UrlInterface.php#L13
	URLType uint8

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
		if log.IsDebug() {
			log.Debug("config.Manager.ApplyDefaults", k, v)
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
	if log.IsDebug() {
		log.Debug("config.Manager.ApplyCoreConfigData", "rows", rows)
	}
	if err != nil {
		return log.Error("config.Manager.ApplyCoreConfigData.LoadSlice", "err", err)
	}

	for _, cd := range ccd {
		if cd.Value.Valid {
			// scope.ID(cd.ScopeID) because cd.ScopeID is a struct field and cannot satisfy interface scope.IDer
			m.Write(Path(cd.Path), Scope(scope.FromString(cd.Scope), cd.ScopeID))
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
		return log.Error("config.Manager.Write.newArg", "err", err)
	}

	if log.IsDebug() {
		log.Debug("config.Manager.Write", "path", a.scopePath(), "val", a.v)
	}

	m.v.Set(a.scopePath(), a.v)
	m.sendMsg(a)
	return nil
}

// get generic getter ... not sure if this should be public ...
func (m *Manager) get(o ...ArgFunc) interface{} {
	a, err := newArg(o...)
	if err != nil {
		return log.Error("config.Manager.get.newArg", "err", err)
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
		log.Error("config.Manager.IsSet.newArg", "err", err)
		return false
	}
	return m.v.IsSet(a.scopePath())
}
