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

package directory

import (
	"fmt"

	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/config/model"
	"github.com/corestoreio/csfw/config/scope"
)

type ConfigCurrenciesInstalled struct {
	model.StringCSV
}

func NewConfigCurrenciesInstalled(path string) ConfigCurrenciesInstalled {
	return ConfigCurrenciesInstalled{
		StringCSV: model.NewStringCSV(
			path,
		),
	}
}

// Write writes an int value and checks if the int value is within the allowed Options.
func (p ConfigCurrenciesInstalled) Write(w config.Writer, v []string, s scope.Scope, id int64) error {

	for _, cur := range v {
		if len(cur) != 3 {
			return fmt.Errorf("Incorrect currency %s. Length must the 3 characters.", cur)
		}
	}

	return p.StringCSV.Write(w, v, s, id)
}
