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

	"github.com/corestoreio/csfw/utils"
)

const (
	VisibleAbsent Visible = iota // must start from 0
	VisibleYes
	VisibleNo
)

type (
	// Visible because GoLang bool cannot be nil 8-) and also in love to
	// https://github.com/magento/magento2/blob/0.74.0-beta9/app/code/Magento/Catalog/Model/Product/Attribute/Source/Status.php#L14
	// Main reason is to detect a change when merging section, group and field slices
	Visible uint8
)

const _Visible_name = "VisibleAbsentVisibleYesVisibleNo"

var _Visible_index = [...]uint8{0, 13, 23, 32}

func (i Visible) String() string {
	if i+1 >= Visible(len(_Visible_index)) {
		return fmt.Sprintf("Visible(%d)", i)
	}
	return _Visible_name[_Visible_index[i]:_Visible_index[i+1]]
}

// VisibleNames returns a slice containing all constant names
func VisibleNames() (r utils.StringSlice) {
	return r.SplitStringer8(_Visible_name, _Visible_index[:]...)
}

// MarshalJSON implements marshaling into a human readable string. @todo UnMarshal
func (v Visible) MarshalJSON() ([]byte, error) {
	return []byte(`"` + v.String() + `"`), nil
}
