// Copyright 2015 CoreStore Authors
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
	"fmt"

	"github.com/corestoreio/csfw/storage/csdb"
	"github.com/corestoreio/csfw/storage/dbr"
	"github.com/corestoreio/csfw/utils"
	"github.com/juju/errgo"
	"github.com/spf13/viper"
)

const (
	IDScopeAbsent ScopeID = iota // order of the constants is used for comparison
	IDScopeDefault
	IDScopeWebsite
	IDScopeGroup
	IDScopeStore
)

const (
	// StringScopeDefault defines the global scope. Stored in table core_config_data.scope.
	StringScopeDefault = "default"
	// StringScopeWebsites defines the website scope which has default as parent and stores as child.
	//  Stored in table core_config_data.scope.
	StringScopeWebsites = "websites"
	// StringScopeStores defines the store scope which has default and websites as parent.
	//  Stored in table core_config_data.scope.
	StringScopeStores = "stores"

	LeftDelim  = "{{"
	RightDelim = "}}"

	// PathCSBaseURL main CoreStore base URL, used if no configuration on a store level can be found.
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

	// ScopeID used in constants where default is the lowest and store the highest. Func String() attached.
	// Part of ScopePerm.
	ScopeID uint8

	// Retriever implements how to get the website or store ID.
	// Duplicated to avoid import cycles. :-(
	Retriever interface {
		ID() int64
	}

	Reader interface {
		// GetString returns a string from the manager. Example usage:
		// Default value: GetString(config.Path("general/locale/timezone"))
		// Website value: GetString(config.Path("general/locale/timezone"), config.ScopeWebsite(w))
		// Store   value: GetString(config.Path("general/locale/timezone"), config.ScopeStore(s))
		GetString(...OptionFunc) string

		// GetBool returns bool from the manager. Example usage see GetString.
		GetBool(...OptionFunc) bool
	}

	Writer interface {
		// Write puts a value back into the manager. Example usage:
		// Default Scope: Write(config.Path("currency", "option", "base"), config.Value("USD"))
		// Website Scope: Write(config.Path("currency", "option", "base"), config.Value("EUR"), config.ScopeWebsite(w))
		// Store   Scope: Write(config.Path("currency", "option", "base"), config.ValueReader(resp.Body), config.ScopeStore(s))
		Write(...OptionFunc) error
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
	TableCollection csdb.TableStructureSlice
	ErrEmptyKey     = errors.New("Key is empty")
)

// NewManager creates the main new configuration for all scopes: default, website and store
func NewManager() *Manager {
	s := &Manager{
		v: viper.New(),
	}
	s.v.SetDefault(ScopeKey(Path(PathCSBaseURL)), CSBaseURL)
	return s
}

// ApplyDefaults reads the map and applies the keys and values to the default configuration
func (m *Manager) ApplyDefaults(ss Sectioner) *Manager {
	ctxLog := logger.WithField("Scope", "ApplyDefaults")
	for k, v := range ss.Defaults() {
		ctxLog.Debug(k, v)
		m.v.SetDefault(k, v)
	}
	return m
}

func (m *Manager) ApplyCoreConfigData(dbrSess dbr.SessionRunner) error {
	var ccd TableCoreConfigDataSlice
	if _, err := csdb.LoadSlice(dbrSess, TableCollection, TableIndexCoreConfigData, &ccd); err != nil {
		return errgo.Mask(err)
	}
	// @todo
	return nil
}

// Write puts a value back into the manager. Example usage:
// Default Scope: Write(config.Path("currency", "option", "base"), config.Value("USD"))
// Website Scope: Write(config.Path("currency", "option", "base"), config.Value("EUR"), config.ScopeWebsite(w))
// Store   Scope: Write(config.Path("currency", "option", "base"), config.ValueReader(resp.Body), config.ScopeStore(s))
func (m *Manager) Write(o ...OptionFunc) error {
	k, v := ScopeKeyValue(o...)
	if k == "" {
		return ErrEmptyKey
	}
	m.v.Set(k, v)
	return nil
}

// GetString returns a string from the manager. Example usage:
// Default value: GetString(config.Path("general/locale/timezone"))
// Website value: GetString(config.Path("general/locale/timezone"), config.ScopeWebsite(w))
// Store   value: GetString(config.Path("general/locale/timezone"), config.ScopeStore(s))
func (m *Manager) GetString(o ...OptionFunc) string {
	return m.v.GetString(ScopeKey(o...))
}

// @todo use the backend model of a config value. most/all magento string slices are comma lists.
func (m *Manager) GetStringSlice(o ...OptionFunc) []string {
	return m.v.GetStringSlice(ScopeKey(o...))
}

// GetBool returns bool from the manager. Example usage see GetString.
func (m *Manager) GetBool(o ...OptionFunc) bool {
	return m.v.GetBool(ScopeKey(o...))
}

// GetFloat64 returns a float64 from the manager. Example usage see GetString.
func (m *Manager) GetFloat64(o ...OptionFunc) float64 {
	return m.v.GetFloat64(ScopeKey(o...))
}

// @todo consider adding other Get* from the viper package

// AllKeys return all keys regardless where they are set
func (m *Manager) AllKeys() []string { return m.v.AllKeys() }

const _ScopeID_name = "ScopeAbsentScopeDefaultScopeWebsiteScopeGroupScopeStore"

var _ScopeID_index = [...]uint8{0, 11, 23, 35, 45, 55}

// String human readable name of ScopeID. For Marshaling see ScopePerm
func (i ScopeID) String() string {
	if i+1 >= ScopeID(len(_ScopeID_index)) {
		return fmt.Sprintf("ScopeID(%d)", i)
	}
	return _ScopeID_name[_ScopeID_index[i]:_ScopeID_index[i+1]]
}

// ScopeIDNames returns a slice containing all constant names
func ScopeIDNames() (r utils.StringSlice) {
	return r.SplitStringer8(_ScopeID_name, _ScopeID_index[:]...)
}

var _ Reader = (*mockScopeReader)(nil)

// MockScopeReader used for testing
type mockScopeReader struct {
	s func(path string) string
	b func(path string) bool
}

// NewMockScopeReader used for testing
func NewMockScopeReader(
	s func(path string) string,
	b func(path string) bool,
) *mockScopeReader {

	return &mockScopeReader{
		s: s,
		b: b,
	}
}

func (sr mockScopeReader) GetString(opts ...OptionFunc) string {
	if sr.s == nil {
		return ""
	}
	return sr.s(ScopeKey(opts...))
}

func (sr mockScopeReader) GetBool(opts ...OptionFunc) bool {
	if sr.b == nil {
		return false
	}
	return sr.b(ScopeKey(opts...))
}
