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

package element

// Visible* defines yes/no/null values if a configuration field should be visible.
// If null then the field is a candidate for merging.
const (
	VisibleAbsent Visible = iota // must start from 0
	VisibleYes
	VisibleNo
)

// Visible because GoLang bool cannot be nil 8-) and also in love to
// https://github.com/magento/magento2/blob/0.74.0-beta9/app/code/Magento/Catalog/Model/Product/Attribute/Source/Status.php#L14
// Main reason is to detect a change when merging section, group and field slices
type Visible uint8

// MarshalJSON implements marshaling into a human readable string. @todo UnMarshal
func (v Visible) MarshalJSON() ([]byte, error) {
	switch v {
	case VisibleAbsent:
		return []byte("null"), nil
	case VisibleYes:
		return []byte("true"), nil
	}
	return []byte(`false`), nil
}
