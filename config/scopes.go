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

// TableCollection handles all tables and its columns. init() in generated Go file will set the value.
var TableCollection csdb.TableStructureSlice

type (
	// UrlType defines the type of the URL. Used in const declaration.
	// @see https://github.com/magento/magento2/blob/0.74.0-beta7/lib/internal/Magento/Framework/UrlInterface.php#L13
	URLType uint8

	// ScopeID used in constants where default is the lowest and store the highest. Func String() attached.
	// Part of ScopePerm.
	ScopeID uint8

	// Retriever implements how to get the ID. If Retriever implements CodeRetriever
	// then CodeRetriever has precedence. ID can be any of the website, group or store IDs.
	// Duplicated to avoid import cycles.
	Retriever interface {
		ID() int64
	}

	Reader interface {
		// GetString retrieves a config string value
		GetString(...OptionFunc) string

		// GetBool retrieves a config flag by path, ScopeID and/or ID
		GetBool(...OptionFunc) bool
	}

	Writer interface {
		// SetString sets config value in the corresponding config scope
		Write(...OptionFunc)
		//Write(path, value interface{}, scope ScopeID, r ...Retriever)
	}

	// Scope main configuration struct which includes Viper, unhappy with the name Scope
	Manager struct {
		v *viper.Viper
	}
)

var _ Reader = (*Manager)(nil)
var _ Writer = (*Manager)(nil)

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

	return nil
}

func (m *Manager) Write(o ...OptionFunc) {
	k, v := ScopeKeyValue(o...)
	m.v.Set(k, v)
}

func (m *Manager) GetString(o ...OptionFunc) string {
	return m.v.GetString(ScopeKey(o...))
}

func (m *Manager) GetBool(o ...OptionFunc) bool {
	return m.v.GetBool(ScopeKey(o...))
}

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
