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
	"bytes"
	"encoding/json"
	"errors"

	"github.com/juju/errgo"
)

const (
	TypeCustom FieldType = iota + 1
	TypeHidden
	TypeObscure
	TypeMultiselect
	TypeSelect
	TypeText
	TypeTime
)

var (
	ErrSectionNotFound = errors.New("Section not found")
	ErrGroupNotFound   = errors.New("Group not found")
	ErrFieldNotFound   = errors.New("Field not found")
)

type (
	// DefaultMap contains the default aka global configuration of a package
	DefaultMap map[string]interface{}

	// Option type is returned by the SourceModel interface
	Option struct {
		Value, Label string
	}

	// Sectioner at the moment only for testing
	Sectioner interface {
		// Defaults generates the default configuration from all fields. Key is the path and value the value.
		Defaults() DefaultMap
	}

	// SectionSlice contains a set of Sections. Some nifty helper functions exists.
	SectionSlice []*Section
	// Section defines the layout for the configuration section which contains groups and fields.
	Section struct {
		// ID unique ID and merged with others. 1st part of the path.
		ID    string
		Label string
		// Scope: bit value eg: showInDefault="1" showInWebsite="1" showInStore="1"
		Scope     ScopePerm
		SortOrder int
		// Permission some kind of ACL if some is allowed for read or write access
		Permission uint
		Groups     GroupSlice
	}
	// GroupSlice contains a set of Groups
	GroupSlice []*Group
	// Group defines the layout of a group containing multiple Fields
	Group struct {
		// ID unique ID and merged with others. 2nd part of the path.
		ID      string
		Label   string
		Comment string
		// Scope: bit value eg: showInDefault="1" showInWebsite="1" showInStore="1"
		Scope     ScopePerm
		SortOrder int
		Fields    FieldSlice
	}

	// FieldType used in constants to define the frontend and input type
	FieldType uint

	// FieldSlice contains a set of Fields
	FieldSlice []*Field
	// Element contains the final path element of a configuration.
	// @see magento2/app/code/Magento/Config/etc/system_file.xsd
	Field struct {
		// ID unique ID and NOT merged with others. 3rd and final part of the path.
		ID   string
		Type interface {
			Type() FieldType
			ToHTML() string // @see \Magento\Framework\Data\Form\Element\AbstractElement
		}
		Label   string
		Comment string
		// Scope: bit value eg: showInDefault="1" showInWebsite="1" showInStore="1"
		Scope     ScopePerm
		SortOrder int
		// Visible used for configuration settings which are not exposed to the user.
		// In Magento2 they do not have an entry in the system.xml
		Visible Visible
		// SourceModel defines how to retrieve all option values
		SourceModel interface {
			Options() []Option
		}
		// BackendModel defines @todo think about AddData
		BackendModel interface {
			AddData(interface{})
			Save() error
		}
		Default interface{}
	}
)

// Type returns the current field type and satisfies the interface of Field.Type
func (t FieldType) Type() FieldType {
	return t
}

// toHTML noop function to satisfies the interface of Field.Type
func (t FieldType) ToHTML() string {
	return ""
}

var _ Sectioner = (*SectionSlice)(nil)

// NewConfiguration creates a new validated SectionSlice with a three level configuration.
// Panics if a path is redundant.
func NewConfiguration(sections ...*Section) SectionSlice {
	ss := SectionSlice(sections)
	if err := ss.Validate(); err != nil {
		logger.WithField("NewConfiguration", "Validate").Warn(err)
		panic(err)
	}
	return ss
}

// Defaults iterates over all slices, creates a path and uses the default value
// to return a map.
func (ss SectionSlice) Defaults() DefaultMap {
	var dm = make(DefaultMap)
	for _, s := range ss {
		for _, g := range s.Groups {
			for _, f := range g.Fields {
				dm[path(s, g, f)] = f.Default
			}
		}
	}
	return dm
}

// Merge merges n SectionSlices into the current slice. Behaviour for duplicates: 1. Warning 2. Last item wins.
func (ss *SectionSlice) MergeAll(sSlices ...SectionSlice) error {
	for _, sl := range sSlices {
		for _, s := range sl {
			if err := (*ss).Merge(s); err != nil {
				return errgo.Mask(err)
			}
		}
	}
	return nil
}

// Merge copies the data from a Section into this slice. Appends if ID is not found
// in this slice otherwise overrides struct fields if not empty.
func (ss *SectionSlice) Merge(s *Section) error {
	cs, err := (*ss).FindByID(s.ID) // cs current section
	if cs == nil || err != nil {
		(*ss).Append(s)
		return nil
	}
	// Maybe that logging is helpful
	logger.WithField("SectionSlice", "merge").Debugf("Label, Scope, SortOrder, Permission of <<%#v>> merged into <<%#v>>", s, cs)
	if s.Label != "" {
		cs.Label = s.Label
	}
	if s.Scope > 0 {
		cs.Scope = s.Scope
	}
	if s.SortOrder != 0 {
		cs.SortOrder = s.SortOrder
	}
	if s.Permission > 0 {
		cs.Permission = s.Permission
	}

	return cs.Groups.Merge(s.Groups...)
}

// FindByID returns a Section pointer or nil if not found. Please check for nil and do not a
func (ss SectionSlice) FindByID(id string) (*Section, error) {
	for _, s := range ss {
		if s != nil && s.ID == id {
			return s, nil
		}
	}
	return nil, ErrSectionNotFound
}

// Append adds 0..n *Section
func (ss *SectionSlice) Append(s ...*Section) *SectionSlice {
	*ss = append(*ss, s...)
	return ss
}

// ToJson transforms the whole slice into JSON
func (ss SectionSlice) ToJson() string {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(ss); err != nil {
		logger.WithField("SectionSlice", "ToJson").Error(err)
		return ""
	}
	return buf.String()
}

// Validate fully validates a configuration for all three hierarchy levels.
// 1. Checks if all slices have at least one entry
// 2. Checks for redundant paths
// 3. ...
func (ss SectionSlice) Validate() error {
	var pc = make(map[string]bool) // pc path checker
	if len(ss) == 0 {
		return errgo.New("SectionSlice is empty")
	}
	for _, s := range ss {
		if len(s.Groups) == 0 {
			return errgo.Newf("%s does not contain groups", s.ID)
		}
		for _, g := range s.Groups {
			if len(g.Fields) == 0 {
				return errgo.Newf("%s/%s does not contain fields", s.ID, g.ID)
			}
			for _, f := range g.Fields {
				p := path(s, g, f)
				if pc[p] {
					return errgo.Newf("Duplicate entry for path %s", p)
				}
				pc[p] = true
			}
		}
	}
	pc = nil
	return nil
}

// FindByID returns a Group pointer or nil if not found
func (gs GroupSlice) FindByID(id string) (*Group, error) {
	for _, g := range gs {
		if g != nil && g.ID == id {
			return g, nil
		}
	}
	return nil, ErrGroupNotFound
}

// Append adds *Group (variadic) to the GroupSlice
func (gs *GroupSlice) Append(g ...*Group) *GroupSlice {
	*gs = append(*gs, g...)
	return gs
}

// Merge copies the data from a groups into this slice. Appends if ID is not found
// in this slice otherwise overrides struct fields if not empty.
func (gs *GroupSlice) Merge(groups ...*Group) error {
	for _, g := range groups {
		if err := (*gs).merge(g); err != nil {
			return errgo.Mask(err)
		}
	}
	return nil
}

func (gs *GroupSlice) merge(g *Group) error {
	cg, err := (*gs).FindByID(g.ID) // cg current group
	if cg == nil || err != nil {
		(*gs).Append(g)
		return nil
	}

	// Maybe that logging is helpful
	logger.WithField("GroupSlice", "merge").Debugf("Label, Comment, Scope, SortOrder of <<%#v>> merged into <<%#v>>", g, cg)

	if g.Label != "" {
		cg.Label = g.Label
	}
	if g.Comment != "" {
		cg.Comment = g.Comment
	}
	if g.Scope > 0 {
		cg.Scope = g.Scope
	}
	if g.SortOrder != 0 {
		cg.SortOrder = g.SortOrder
	}
	cg.Fields.Merge(g.Fields...)
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
	cf, err := (*fs).FindByID(f.ID) // cf current field
	if cf == nil || err != nil {
		(*fs).Append(f)
		return nil
	}

	// Maybe that logging is helpful
	logger.WithField("FieldSlice", "merge").Debugf("Type, Label, Comment, Scope, SortOrder of <<%#v>> merged into <<%#v>>", f, cf)
	if f.Type != nil {
		cf.Type = f.Type
	}
	if f.Label != "" {
		cf.Label = f.Label
	}
	if f.Comment != "" {
		cf.Comment = f.Comment
	}
	if f.Scope > ScopeAbsent {
		cf.Scope = f.Scope
	}
	if f.SortOrder != 0 {
		cf.SortOrder = f.SortOrder
	}
	if f.Visible > VisibleAbsent {
		cf.Visible = f.Visible
	}
	// @todo
	//// SourceModel defines how to retrieve all option values
	//SourceModel interface {
	//Options() []Option
	//}
	//// BackendModel defines @todo think about AddData
	//BackendModel interface {
	//AddData(interface{})
	//Save() error
	//}
	//Default interface{}

	return nil
}

// path creates a valid configuration path with slashes as separators.
func path(s *Section, g *Group, f *Field) string {
	return s.ID + "/" + g.ID + "/" + f.ID
}
