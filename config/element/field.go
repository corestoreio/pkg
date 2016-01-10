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

import (
	"errors"
	"sort"

	"github.com/corestoreio/csfw/config/path"
	"github.com/corestoreio/csfw/store/scope"
	"github.com/juju/errgo"
)

// ErrFieldNotFound error when a field cannot be found.
var ErrFieldNotFound = errors.New("Field not found")

// FieldSlice contains a set of Fields. Has several method receivers attached.
//  Thread safe for reading but not for modifying.
type FieldSlice []*Field

// Field contains the final path element of a configuration. Includes several options.
//  Thread safe for reading but not for modifying.
// @see magento2/app/code/Magento/Config/etc/system_file.xsd
type Field struct {
	// ID unique ID and NOT merged with others. 3rd and final part of the path.
	ID string
	// ConfigPath if provided defines the storage path and overwrites the path from
	// section.id + group.id + field.id
	ConfigPath string `json:",omitempty"`
	// Type is used for the front end on how to display a Field
	Type FieldTyper `json:",omitempty"`
	// Label a short label of the field
	Label string `json:",omitempty"`
	// Comment can contain HTML
	Comment LongText `json:",omitempty"`
	// Tooltip used for frontend and can contain HTML
	Tooltip LongText `json:",omitempty"`
	// Scope: bit value eg: showInDefault="1" showInWebsite="1" showInStore="1"
	Scope scope.Perm `json:",omitempty"`
	// SortOrder in ascending order
	SortOrder int `json:",omitempty"`
	// Visible used for configuration settings which are not exposed to the user.
	// In Magento2 they do not have an entry in the system.xml
	Visible Visible `json:",omitempty"`

	// CanBeEmpty only used in HTML forms for multiselect fields
	// Use case: lib/internal/Magento/Framework/Data/Form/Element/Multiselect.php::getElementHtml()
	CanBeEmpty bool `json:",omitempty"`
	// Default can contain any default config value: float64, int64, string, bool
	Default interface{} `json:",omitempty"`
}

// NewFieldSlice wrapper to create a new FieldSlice
func NewFieldSlice(fs ...*Field) FieldSlice {
	return FieldSlice(fs)
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

// Append adds *Field (variadic) to the FieldSlice. Not thread safe.
func (fs *FieldSlice) Append(f ...*Field) *FieldSlice {
	*fs = append(*fs, f...)
	return fs
}

// Merge copies the data from a Field into this slice. Appends if ID is not found
// in this slice otherwise overrides struct fields if not empty. Not thread safe.
func (fs *FieldSlice) Merge(fields ...*Field) error {
	for _, f := range fields {
		if err := (*fs).merge(f); err != nil {
			return errgo.Mask(err)
		}
	}
	return nil
}

// merge merges field f into the slice. Appends the field if the Id is new.
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
	if !f.Comment.IsEmpty() {
		cf.Comment = f.Comment.Copy()
	}
	if !f.Tooltip.IsEmpty() {
		cf.Tooltip = f.Tooltip.Copy()
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
	cf.CanBeEmpty = f.CanBeEmpty
	if f.Default != nil {
		cf.Default = f.Default
	}

	return nil
}

// Sort convenience helper. Not thread safe.
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

// FQPathDefault returns the default fully qualified path of either
// Section.ID + Group.ID + Field.ID OR Field.ConfgPath if set.
// Function may panic.
func (f *Field) FQPathDefault(prePaths ...string) string {
	p := path.Path{
		Scope: scope.DefaultID,
	}
	if f.ConfigPath != "" {
		p.Parts = []string{f.ConfigPath}
	} else {
		p.Parts = append(prePaths, f.ID)
	}
	return p.String()

}
