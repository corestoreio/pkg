// Copyright 2015-present, Cyrill @ Schumacher.fm and the CoreStore contributors
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
	"sort"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/pkg/store/scope"
)

// Fields contains a set of Fields. Has several method receivers attached.
// Thread safe for reading but not for modifying.
type Fields []*Field

// Field contains the final path element of a configuration. Includes several
// options. Thread safe for reading but not for modifying. @see
// magento2/app/code/Magento/Config/etc/system_file.xsd
type Field struct {
	// ID unique ID and NOT merged with others. 3rd and final part of the path.
	ID string
	// ConfigPath if provided defines the storage path and overwrites the path from
	// section.id + group.id + field.id. ConfigPath can be nil.
	ConfigPath string `json:",omitempty"` // omitempty does not yet work on non-pointer structs that is reason for the interface
	// Type is used for the front end on how to display a Field
	Type FieldTyper `json:",omitempty"`
	// Label a short label of the field
	Label string `json:",omitempty"`
	// Comment can contain HTML
	Comment string `json:",omitempty"`
	// Tooltip used for frontend and can contain HTML
	Tooltip string `json:",omitempty"`
	// Scopes: bit value eg: showInDefault="1" showInWebsite="1" showInStore="1"
	// Scopes can contain multiple Scope but no more than Default, Website and
	// Store.
	Scopes scope.Perm `json:",omitempty"`
	// SortOrder in ascending order
	SortOrder int `json:",omitempty"`
	// Visible used for configuration settings which are not exposed to the user.
	Visible bool `json:",omitempty"`

	// CanBeEmpty only used in HTML forms for multiselect fields
	// Use case: lib/internal/Magento/Framework/Data/Form/Element/Multiselect.php::getElementHtml()
	CanBeEmpty bool `json:",omitempty"`
	// Default can contain any default config value: float64, int64, string, bool
	Default interface{} `json:",omitempty"`
}

// MakeFields wrapper to create a new Fields
func MakeFields(fs ...*Field) Fields {
	return Fields(fs)
}

// Find returns a Field pointer or ErrFieldNotFound. Route must be a single
// part. E.g. if you have path "a/b/c" route would be in this case "c". 2nd
// argument int contains the slice index of the field. Error behaviour: NotFound
func (fs Fields) Find(id string) (*Field, int, error) {
	for i, f := range fs {
		if f.ID != "" && f.ID == id {
			return f, i, nil
		}
	}
	return nil, 0, errors.NotFound.Newf("[element] Field %q not found", id)
}

// Append adds *Field (variadic) to the Fields. Not thread safe.
func (fs Fields) Append(f ...*Field) Fields {
	return append(fs, f...)
}

// Merge copies the data from a Field into this slice and returns the new slice.
// Appends if ID is not found in this slice otherwise overrides struct fields if
// not empty. Not thread safe.
func (fs Fields) Merge(fields ...*Field) Fields {
	for _, f := range fields {
		fs = fs.merge(f)
	}
	return fs
}

// merge merges field f into the slice. Appends the field if the Id is new.
func (fs Fields) merge(f *Field) Fields {

	cf, idx, err := fs.Find(f.ID) // cf current field
	if err != nil {
		cf = f
		fs = append(fs, cf)
		idx = len(fs) - 1
	}

	fs[idx] = cf.Update(f)
	return fs
}

// Sort convenience helper. Not thread safe.
func (fs Fields) Sort() Fields {
	sort.Sort(fs)
	return fs
}

func (fs Fields) Len() int {
	return len(fs)
}

func (fs Fields) Swap(i, j int) {
	fs[i], fs[j] = fs[j], fs[i]
}

func (fs Fields) Less(i, j int) bool {
	return fs[i].SortOrder < fs[j].SortOrder
}

// Update applies the data from the new Field to the old field and returns the
// updated Field. Only non-empty values will be copied and byte slices gets
// cloned. The returned Field allows modifications.
func (f *Field) Update(new *Field) *Field {
	if new.Type != nil {
		f.Type = new.Type
	}
	if "" != new.Label {
		f.Label = new.Label
	}
	if "" != new.Comment {
		f.Comment = new.Comment
	}
	if "" != new.Tooltip {
		f.Tooltip = new.Tooltip
	}
	if new.Scopes > 0 {
		f.Scopes = new.Scopes
	}
	if new.SortOrder != 0 {
		f.SortOrder = new.SortOrder
	}
	if new.Visible > VisibleAbsent {
		f.Visible = new.Visible
	}
	f.CanBeEmpty = new.CanBeEmpty
	if new.Default != nil {
		f.Default = new.Default
	}
	return f
}

// Route returns the merged route of either Section.ID + Group.ID + Field.ID OR
// Field.ConfgPath if set. The returned route has not been validated.
func (f Field) Route(preRoutes ...string) string {
	if "" != f.ConfigPath {
		return f.ConfigPath
	}
	return JoinRoutes(preRoutes...)
}
