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
	"sort"
	"strings"

	"github.com/corestoreio/csfw/utils"
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
		Label string `json:",omitempty"`
		// Scope: bit value eg: showInDefault="1" showInWebsite="1" showInStore="1"
		Scope     ScopePerm `json:",omitempty"`
		SortOrder int       `json:",omitempty"`
		// Permission some kind of ACL if someone is allowed for no,read or write access @todo
		Permission uint `json:",omitempty"`
		Groups     GroupSlice
	}
	// GroupSlice contains a set of Groups
	GroupSlice []*Group
	// Group defines the layout of a group containing multiple Fields
	Group struct {
		// ID unique ID and merged with others. 2nd part of the path.
		ID      string
		Label   string `json:",omitempty"`
		Comment string `json:",omitempty"`
		// Scope: bit value eg: showInDefault="1" showInWebsite="1" showInStore="1"
		Scope     ScopePerm `json:",omitempty"`
		SortOrder int       `json:",omitempty"`
		Fields    FieldSlice
	}

	// FieldType used in constants to define the frontend and input type
	FieldType uint8

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
		Scope     ScopePerm `json:",omitempty"`
		SortOrder int       `json:",omitempty"`
		// Visible used for configuration settings which are not exposed to the user.
		// In Magento2 they do not have an entry in the system.xml
		Visible Visible `json:",omitempty"`
		// SourceModel defines how to retrieve all option values
		SourceModel FieldSourceModeller `json:",omitempty"`
		// BackendModel defines how to save and load? the data
		BackendModel FieldBackendModeller `json:",omitempty"`
		// Default can contain any default config value: float64, int64, string, bool
		Default interface{} `json:",omitempty"`
	}

	// FieldTyper defines which front end type a configuration value is and generates the HTML for it
	FieldTyper interface {
		Type() FieldType
		ToHTML() []byte // @see \Magento\Framework\Data\Form\Element\AbstractElement
	}

	// FieldSourceModeller defines how to retrieve all option values
	FieldSourceModeller interface {
		Options() []Option
	}

	// FieldBackendModeller defines how to save and load? the data @todo think about AddData
	FieldBackendModeller interface {
		AddData(interface{})
		Save() error
	}
)

var _ Sectioner = (*SectionSlice)(nil)
var _ FieldTyper = (*FieldType)(nil)

// Type returns the current field type and satisfies the interface of Field.Type
func (t FieldType) Type() FieldType {
	return t
}

// ToHTML noop function to satisfies the interface of Field.Type
func (t FieldType) ToHTML() []byte {
	return nil
}

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

// NewConfigurationMerge creates a new validated SectionSlice with a three level configuration.
// Before validation, slices are all merged together. Panics if a path is redundant.
// Only use this function if your package configuration really has duplicated entries.
func NewConfigurationMerge(sections ...*Section) SectionSlice {
	var ss SectionSlice
	ss.Merge(sections...)
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
				dm[ScopeKey(Path(s.ID, g.ID, f.ID))] = f.Default
			}
		}
	}
	return dm
}

// TotalFields calculates the total amount of all fields
func (ss SectionSlice) TotalFields() int {
	fs := 0
	for _, s := range ss {
		for _, g := range s.Groups {
			for _ = range g.Fields {
				fs++
			}
		}
	}
	return -^+^-fs
}

// MergeMultiple merges n SectionSlices into the current slice. Behaviour for duplicates: Last item wins.
func (ss *SectionSlice) MergeMultiple(sSlices ...SectionSlice) error {
	for _, sl := range sSlices {
		if err := (*ss).Merge(sl...); err != nil {
			return err
		}
	}
	return nil
}

// Merge merges n Sections into the current slice. Behaviour for duplicates: Last item wins.
func (ss *SectionSlice) Merge(sections ...*Section) error {
	for _, s := range sections {
		if s != nil {
			if err := (*ss).merge(s); err != nil {
				return errgo.Mask(err)
			}
		}
	}
	return nil
}

// Merge copies the data from a Section into this slice. Appends if ID is not found
// in this slice otherwise overrides struct fields if not empty.
func (ss *SectionSlice) merge(s *Section) error {
	if s == nil {
		return nil
	}
	cs, err := (*ss).FindByID(s.ID) // cs current section
	if cs == nil || err != nil {
		cs = &Section{ID: s.ID}
		(*ss).Append(cs)
	}

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

// FindGroupByPath searches for a group using the first two path segments.
// If one argument is given then considered as the full path e.g. a/b/c
// If two or more arguments are given then each argument will be treated as a path part.
func (ss SectionSlice) FindGroupByPath(paths ...string) (*Group, error) {
	if len(paths) == 1 {
		paths = strings.Split(paths[0], "/")
	}
	if len(paths) < 2 {
		return nil, errgo.Mask(ErrGroupNotFound)
	}
	cs, err := ss.FindByID(paths[0])
	if err != nil {
		return nil, errgo.Mask(err)
	}
	return cs.Groups.FindByID(paths[1])
}

// FindGroupByPath searches for a field using the all three path segments.
// If one argument is given then considered as the full path e.g. a/b/c
// If three arguments are given then each argument will be treated as a path part.
func (ss SectionSlice) FindFieldByPath(paths ...string) (*Field, error) {
	if len(paths) == 1 {
		paths = strings.Split(paths[0], "/")
	}
	if len(paths) < 3 {
		return nil, errgo.Mask(ErrFieldNotFound)
	}
	cg, err := ss.FindGroupByPath(paths...)
	if err != nil {
		return nil, errgo.Mask(err)
	}
	return cg.Fields.FindByID(paths[2])
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

// Validate checks for duplicated configuration paths in all three hierarchy levels.
func (ss SectionSlice) Validate() error {
	if len(ss) == 0 {
		return errgo.New("SectionSlice is empty")
	}
	// @todo try to pick the right strategy between maps and slice depending on the overall size of a full SectionSlice
	var pc = make(utils.StringSlice, ss.TotalFields()) // pc path checker
	i := 0
	for _, s := range ss {
		for _, g := range s.Groups {
			for _, f := range g.Fields {
				p := ScopeKey(Path(s.ID, g.ID, f.ID))
				if pc.Include(p) {
					return errgo.Newf("Duplicate entry for path %s :: %s", p, ss.ToJson())
				}
				pc[i] = p
				i++
			}
		}
	}
	return nil
}

// SortAll recursively sorts all slices
func (ss *SectionSlice) SortAll() *SectionSlice {
	for _, s := range *ss {
		for _, g := range s.Groups {
			g.Fields.Sort()
		}
		s.Groups.Sort()
	}
	return ss.Sort()
}

// Sort convenience helper
func (ss *SectionSlice) Sort() *SectionSlice {
	sort.Sort(ss)
	return ss
}

func (ss *SectionSlice) Len() int {
	return len(*ss)
}

func (ss *SectionSlice) Swap(i, j int) {
	(*ss)[i], (*ss)[j] = (*ss)[j], (*ss)[i]
}

func (ss *SectionSlice) Less(i, j int) bool {
	return (*ss)[i].SortOrder < (*ss)[j].SortOrder
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
	if g == nil {
		return nil
	}
	cg, err := (*gs).FindByID(g.ID) // cg current group
	if cg == nil || err != nil {
		cg = g
		(*gs).Append(cg)
	}

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

// ToJson transforms the whole slice into JSON
func (gs GroupSlice) ToJson() string {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(gs); err != nil {
		logger.WithField("GroupSlice", "ToJson").Error(err)
		return ""
	}
	return buf.String()
}

// Sort convenience helper
func (gs *GroupSlice) Sort() *GroupSlice {
	sort.Sort(gs)
	return gs
}

func (gs *GroupSlice) Len() int {
	return len(*gs)
}

func (gs *GroupSlice) Swap(i, j int) {
	(*gs)[i], (*gs)[j] = (*gs)[j], (*gs)[i]
}

func (gs *GroupSlice) Less(i, j int) bool {
	return (*gs)[i].SortOrder < (*gs)[j].SortOrder
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

const _FieldType_name = "TypeButtonTypeCustomTypeLabelTypeHiddenTypeImageTypeObscureTypeMultiselectTypeSelectTypeTextTypeTextareaTypeTime"

var _FieldType_index = [...]uint8{10, 20, 29, 39, 48, 59, 74, 84, 92, 104, 112}

func (i FieldType) String() string {
	i -= 1
	if i >= FieldType(len(_FieldType_index)) {
		return "FieldType(?)"
	}
	hi := _FieldType_index[i]
	lo := uint8(0)
	if i > 0 {
		lo = _FieldType_index[i-1]
	}
	return _FieldType_name[lo:hi]
}

// MarshalJSON implements marshaling into a human readable string. @todo UnMarshal
func (i FieldType) MarshalJSON() ([]byte, error) {
	return []byte(`"` + strings.ToLower(i.String()[4:]) + `"`), nil
}
