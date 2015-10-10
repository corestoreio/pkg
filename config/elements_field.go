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
	"errors"
	"sort"
	"strings"

	"github.com/corestoreio/csfw/config/scope"
	"github.com/juju/errgo"
)

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
)

// ErrFieldNotFound error when a field cannot be found.
var ErrFieldNotFound = errors.New("Field not found")

type (

	// FieldType used in constants to define the frontend and input type
	FieldType uint8

	// FieldTyper defines which front end type a configuration value is and generates the HTML for it
	FieldTyper interface {
		Type() FieldType
		ToHTML() []byte // @see \Magento\Framework\Data\Form\Element\AbstractElement
	}

	// FieldSlice contains a set of Fields
	FieldSlice []*Field
	// Element contains the final path element of a configuration.
	// @see magento2/app/code/Magento/Config/etc/system_file.xsd
	Field struct {
		// ID unique ID and NOT merged with others. 3rd and final part of the path.
		ID string
		// Type is used for the front end on how to display a Field
		Type    FieldTyper `json:",omitempty"`
		Label   string     `json:",omitempty"`
		Comment string     `json:",omitempty"`
		// Scope: bit value eg: showInDefault="1" showInWebsite="1" showInStore="1"
		Scope     scope.Perm `json:",omitempty"`
		SortOrder int        `json:",omitempty"`
		// Visible used for configuration settings which are not exposed to the user.
		// In Magento2 they do not have an entry in the system.xml
		Visible Visible `json:",omitempty"`
		// SourceModel defines how to retrieve all option values
		SourceModel SourceModeller `json:",omitempty"`
		// BackendModel defines how to save and load? the data
		BackendModel BackendModeller `json:",omitempty"`
		// Default can contain any default config value: float64, int64, string, bool
		Default interface{} `json:",omitempty"`

		// ConfigPath    string  @todo see typeConfigPath in app/code/Magento/Config/etc/system_file.xsd

	}
)

var _ FieldTyper = (*FieldType)(nil)

// Type returns the current field type and satisfies the interface of Field.Type
func (i FieldType) Type() FieldType {
	return i
}

// ToHTML noop function to satisfies the interface of Field.Type
func (i FieldType) ToHTML() []byte {
	return nil
}

// FindByID returns a Field pointer or nil if not found
func (fs FieldSlice) FindByID(id string) (*Field, error) {
	for _, f := range fs {
		if f != nil && f.ID == id {
			return f, nil
		}
	}
	return nil, ErrFieldNotFound
}

// Append adds *Field (variadic) to the FieldSlice
func (fs *FieldSlice) Append(f ...*Field) *FieldSlice {
	*fs = append(*fs, f...)
	return fs
}

// Merge copies the data from a Section into this slice. Appends if ID is not found
// in this slice otherwise overrides struct fields if not empty.
func (fs *FieldSlice) Merge(fields ...*Field) error {
	for _, f := range fields {
		if err := (*fs).merge(f); err != nil {
			return errgo.Mask(err)
		}
	}
	return nil
}

func (fs *FieldSlice) merge(f *Field) error {
	if f == nil {
		return nil
	}
	cf, err := (*fs).FindByID(f.ID) // cf current field
	if cf == nil || err != nil {
		cf = f
		(*fs).Append(cf)
	}

	if f.Type != nil {
		cf.Type = f.Type
	}
	if f.Label != "" {
		cf.Label = f.Label
	}
	if f.Comment != "" {
		cf.Comment = f.Comment
	}
	if f.Scope > 0 {
		cf.Scope = f.Scope
	}
	if f.SortOrder != 0 {
		cf.SortOrder = f.SortOrder
	}
	if f.Visible > VisibleAbsent {
		cf.Visible = f.Visible
	}
	if f.SourceModel != nil {
		cf.SourceModel = f.SourceModel
	}
	if f.BackendModel != nil {
		cf.BackendModel = f.BackendModel
	}
	if f.Default != nil {
		cf.Default = f.Default
	}

	return nil
}

// Sort convenience helper
func (fs *FieldSlice) Sort() *FieldSlice {
	sort.Sort(fs)
	return fs
}

func (fs *FieldSlice) Len() int {
	return len(*fs)
}

func (fs *FieldSlice) Swap(i, j int) {
	(*fs)[i], (*fs)[j] = (*fs)[j], (*fs)[i]
}

func (fs *FieldSlice) Less(i, j int) bool {
	return (*fs)[i].SortOrder < (*fs)[j].SortOrder
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
