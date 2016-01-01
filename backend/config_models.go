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

package backend

import (
	"fmt"
	"net/http"

	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/config/model"
	"github.com/corestoreio/csfw/config/valuelabel"
	"github.com/corestoreio/csfw/store/scope"
)

// ConfigRedirectToBase enables if a redirect to the base URL should
// happen and with which status code.
type ConfigRedirectToBase struct {
	model.Int
}

// NewConfigRedirectToBase creates a new type.
func NewConfigRedirectToBase(path string, opts ...model.Option) ConfigRedirectToBase {
	return ConfigRedirectToBase{
		Int: model.NewInt(
			path,
			model.WithValueLabelByInt(valuelabel.Ints{
				{0, "No"},
				{1, "Yes (302 Found)"},                // old from Magento
				{http.StatusFound, "Yes (302 Found)"}, // new correct
				{http.StatusMovedPermanently, "Yes (301 Moved Permanently)"},
			}),
			opts...,
		),
	}
}

// Write writes an int value and checks if the int value is within the allowed Options.
func (p ConfigRedirectToBase) Write(w config.Writer, v int, s scope.Scope, id int64) error {

	if false == p.ValueLabel.ContainsValInt(v) {
		return fmt.Errorf("Cannot find %d in list %#v", v, p.ValueLabel)
	}

	return p.Int.Write(w, v, s, id)
}
