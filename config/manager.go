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
	"time"

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

const (
	URLTypeAbsent URLType = iota
	// UrlTypeWeb defines the ULR type to generate the main base URL.
	URLTypeWeb
	// UrlTypeStatic defines the url to the static assets like css, js or theme images
	URLTypeStatic
	// UrlTypeLink hmmm
	// UrlTypeLink
	// UrlTypeMedia defines the ULR type for generating URLs to product photos
	URLTypeMedia
)

type (
	// UrlType defines the type of the URL. Used in const declaration.
	// @see https://github.com/magento/magento2/blob/0.74.0-beta7/lib/internal/Magento/Framework/UrlInterface.php#L13
	URLType uint8

	// Reader implements how to receive thread safe a configuration value from a path and or scope.
	Reader interface {
		GetString(...ArgFunc) string
		GetBool(...ArgFunc) bool
		GetFloat64(o ...ArgFunc) float64
		GetInt(o ...ArgFunc) int
		GetDateTime(o ...ArgFunc) time.Time
	}

	// Writer thread safe storing of configuration values under different paths and scopes.
	Writer interface {
		Write(...ArgFunc) error
	}

	// Manager main configuration struct
	Manager struct {
		// why is Viper private? Because it can maybe replaced by something else ...
		v *viper.Viper
	}
)

var (
	_ Reader = (*Manager)(nil)
	_ Writer = (*Manager)(nil)

	// TableCollection handles all tables and its columns. init() in generated Go file will set the value.
	TableCollection csdb.Manager

	// DefaultManager provides a default manager
	DefaultManager *Manager
)

func init() {
	DefaultManager = NewManager()
}

// NewManager creates the main new configuration for all scopes: default, website and store
func NewManager() *Manager {
	s := &Manager{
		v: viper.New(),
	}
	s.v.Set(newArg(Path(PathCSBaseURL)).scopePath(), CSBaseURL)
	return s
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
			// ScopeID(cd.ScopeID) because cd.ScopeID is a struct field and cannot satisfy interface ScopeIDer
			m.Write(Path(cd.Path), Scope(GetScopeGroup(cd.Scope), ScopeID(cd.ScopeID)))
		}
	}
	return nil
}

// Write puts a value back into the manager. Example usage:
// Default Scope: Write(config.Path("currency", "option", "base"), config.Value("USD"))
// Website Scope: Write(config.Path("currency", "option", "base"), config.Value("EUR"), config.ScopeWebsite(w))
// Store   Scope: Write(config.Path("currency", "option", "base"), config.ValueReader(resp.Body), config.ScopeStore(s))
func (m *Manager) Write(o ...ArgFunc) error {
	a := newArg(o...)
	if a.isBubbling() {
		if log.IsDebug() {
			log.Debug("config.Manager.Write.isBubbling", "path", a.scopePathDefault(), "bubble", a.isBubbling(), "val", a.v)
		}
		m.v.Set(a.scopePathDefault(), a.v)
	}

	if log.IsDebug() {
		log.Debug("config.Manager.Write", "path", a.scopePath(), "val", a.v)
	}
	m.v.Set(a.scopePath(), a.v)

	return nil
}

// get generic getter ... not sure if this should be public ...
func (m *Manager) get(o ...ArgFunc) interface{} {
	a := newArg(o...)
	vs := m.v.Get(a.scopePath()) // vs = value scope
	if vs == nil && a.isBubbling() {
		vs = m.v.Get(a.scopePathDefault())
	}
	return vs
}

// GetString returns a string from the manager. Example usage:
// Default value: GetString(config.Path("general/locale/timezone"))
// Website value: GetString(config.Path("general/locale/timezone"), config.ScopeWebsite(w))
// Store   value: GetString(config.Path("general/locale/timezone"), config.ScopeStore(s))
func (m *Manager) GetString(o ...ArgFunc) string {
	vs := m.get(o...)
	if vs == nil {
		return ""
	}
	return cast.ToString(vs)
}

// @todo use the backend model of a config value. most/all magento string slices are comma lists.
func (m *Manager) GetStringSlice(o ...ArgFunc) []string {
	return nil
	//	return m.v.GetStringSlice(newArg(o...))
}

// GetBool returns bool from the manager. Example usage see GetString.
func (m *Manager) GetBool(o ...ArgFunc) bool {
	vs := m.get(o...)
	if vs == nil {
		return false
	}
	return cast.ToBool(vs)
}

// GetFloat64 returns a float64 from the manager. Example usage see GetString.
func (m *Manager) GetFloat64(o ...ArgFunc) float64 {
	vs := m.get(o...)
	if vs == nil {
		return 0.0
	}
	return cast.ToFloat64(vs)
}

// GetInt returns an int from the manager. Example usage see GetString.
func (m *Manager) GetInt(o ...ArgFunc) int {
	vs := m.get(o...)
	if vs == nil {
		return 0
	}
	return cast.ToInt(vs)
}

// GetDateTime returns a date and time object from the manager. Example usage see GetString.
func (m *Manager) GetDateTime(o ...ArgFunc) time.Time {
	vs := m.get(o...)
	t, err := cast.ToTimeE(vs)
	if err != nil {
		log.Error("config.Manager.GetDateTime.ToTimeE", "err", err, "val", vs)
	}
	return t
}

// @todo consider adding other Get* from the viper package

// AllKeys return all keys regardless where they are set
func (m *Manager) AllKeys() []string { return m.v.AllKeys() }

// IsSet checks if a key is in the config. Does not bubble.
func (m *Manager) IsSet(o ...ArgFunc) bool {
	return m.v.IsSet(newArg(o...).scopePath())
}
