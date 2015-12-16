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

package model

import (
	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/config/scope"
)

// BaseURL represents a path in config.Getter handles BaseURLs and internal validation
type BaseURL struct{ Path }

// NewBaseURL creates a new BaseURL with validation checks when writing values.
func NewBaseURL(path string) BaseURL {
	return BaseURL{Path{path}}
}

// Get returns a base URL
func (p BaseURL) Get(pkgCfg config.SectionSlice, sg config.ScopedGetter) string {
	return p.Path.LookupString(pkgCfg, sg)
}

// Set writes a new base URL and validates it before saving.
func (p BaseURL) Set(w config.Writer, v string, s scope.Scope, id int64) error {
	// todo URL checks app/code/Magento/Config/Model/Config/Backend/Baseurl.php
	return p.Path.Set(w, v, s, id)
}
