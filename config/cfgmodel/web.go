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

package cfgmodel

import (
	"net/url"

	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/store/scope"
	"github.com/corestoreio/csfw/util/errors"
)

// Placeholder constants and their values can occur in the table core_config_data.
// These placeholder must be replaced with the current values.
const (
	PlaceholderBaseURL         = config.LeftDelim + "base_url" + config.RightDelim
	PlaceholderBaseURLSecure   = config.LeftDelim + "secure_base_url" + config.RightDelim
	PlaceholderBaseURLUnSecure = config.LeftDelim + "unsecure_base_url" + config.RightDelim
)

// URL represents a path in config.Getter which handles URLs and internal validation
type URL struct{ Str }

// NewURL creates a new URL with validation checks.
func NewURL(path string, opts ...Option) URL {
	return URL{Str: NewStr(path, opts...)}
}

// Get returns an URL. If the underlying value is empty returns nil,nil.
func (p URL) Get(sg config.ScopedGetter) (*url.URL, error) {
	rawurl, err := p.Str.Get(sg)
	if err != nil {
		return nil, errors.Wrap(err, "[cfgmodel] URL.Str.Get")
	}
	if rawurl == "" {
		return nil, nil
	}
	u, err := url.Parse(rawurl)
	if err != nil {
		return nil, errors.NewFatal(err, "[cfgmodel] URL.Parse")
	}
	return u, nil
}

// Write writes a new URL and validates it before saving. If v is nil, an empty value
// will be written.
func (p URL) Write(w config.Writer, v *url.URL, s scope.Scope, scopeID int64) error {
	var val string
	if v != nil {
		val = v.String()
	}
	return p.Str.Write(w, val, s, scopeID)
}

// BaseURL represents a path in config.Getter handles BaseURLs and internal validation
type BaseURL struct{ Str }

// NewBaseURL creates a new BaseURL with validation checks when writing values.
func NewBaseURL(path string, opts ...Option) BaseURL {
	return BaseURL{Str: NewStr(path, opts...)}
}

// Get returns a base URL
func (p BaseURL) Get(sg config.ScopedGetter) (string, error) {
	return p.Str.Get(sg)
}

// Write writes a new base URL and validates it before saving. @TODO
func (p BaseURL) Write(w config.Writer, v string, s scope.Scope, scopeID int64) error {
	// todo URL checks app/code/Magento/Config/Model/Config/Backend/Baseurl.php
	return p.Str.Write(w, v, s, scopeID)
}
