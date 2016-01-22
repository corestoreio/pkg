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

import "strings"

// Type* defines the type of the front end user input/display form
const (
	TypeButton FieldType = iota + 1 // must be + 1 because 0 is not set
	TypeCustom
	TypeLabel
	TypeHidden
	TypeImage
	TypeObscure
	TypeMultiselect
	TypeSelect
	TypeText
	TypeTextarea
	TypeTime
	TypeDuration
	TypeZMaximum
)

type (

	// FieldType used in constants to define the frontend and input type
	FieldType uint8

	// FieldTyper defines which front end type a configuration value is and generates the HTML for it
	FieldTyper interface {
		Type() FieldType
		ToHTML() []byte // @see \Magento\Framework\Data\Form\Element\AbstractElement
	}
)

// Type returns the current field type and satisfies the interface of Field.Type
func (i FieldType) Type() FieldType {
	return i
}

// ToHTML noop function to satisfies the interface of Field.Type
func (i FieldType) ToHTML() []byte {
	return nil
}

const fieldTypeName = "TypeButtonTypeCustomTypeLabelTypeHiddenTypeImageTypeObscureTypeMultiselectTypeSelectTypeTextTypeTextareaTypeTime"

var fieldTypeIndex = [...]uint8{10, 20, 29, 39, 48, 59, 74, 84, 92, 104, 112}

func (i FieldType) String() string {
	i--
	if i >= FieldType(len(fieldTypeIndex)) {
		return "FieldType(?)"
	}
	hi := fieldTypeIndex[i]
	lo := uint8(0)
	if i > 0 {
		lo = fieldTypeIndex[i-1]
	}
	return fieldTypeName[lo:hi]
}

// MarshalJSON implements marshaling into a human readable string. @todo UnMarshal
func (i FieldType) MarshalJSON() ([]byte, error) {
	return []byte(`"` + strings.ToLower(i.String()[4:]) + `"`), nil
}
