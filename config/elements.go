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

package config

import (
	"sort"
	"strings"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/pkg/store/scope"
)

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

// FieldType used in constants to define the frontend and input type
type FieldType uint8

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

// Sections contains a set of Sections. Some nifty helper functions exists.
// Thread safe for reading. A section slice can be used in many goroutines. It
// must remain lock-free.
type Sections []*Section

// Section defines the layout for the configuration section which contains
// groups and fields. Thread safe for reading but not for modifying.
type Section struct {
	// ID unique ID and merged with others. 1st part of the path.
	ID    string
	Label string `json:",omitempty"`
	// Scopes: bit value eg: showInDefault="1" showInWebsite="1" showInStore="1"
	// Scopes can contain multiple Scope but no more than Default, Website and
	// Store.
	Scopes    scope.Perm `json:",omitempty"`
	SortOrder int        `json:",omitempty"`
	// Resource some kind of ACL if someone is allowed for no,read or write access @todo
	Resource uint   `json:",omitempty"`
	Groups   Groups `json:",omitempty"`
}

// MakeSections wrapper function, for now.
func MakeSections(s ...*Section) Sections {
	return Sections(s)
}

// MakeSectionsValidated creates a new validated Sections with a three level
// configuration.
func MakeSectionsValidated(sections ...*Section) (Sections, error) {
	ss := MakeSections(sections...)
	if err := ss.Validate(); err != nil {
		return nil, errors.Wrap(err, "[element] MakeSections.Validate")
	}
	return ss, nil
}

// MustMakeSectionsValidate same as MakeSectionsValidated but panics on error.
func MustMakeSectionsValidate(sections ...*Section) Sections {
	s, err := MakeSectionsValidated(sections...)
	if err != nil {
		panic(err)
	}
	return s
}

// MakeSectionsMerged creates a new validated Sections with a three level
// configuration. Before validation, slices are all merged together. Panics if a
// path is redundant. Only use this function if your package elements really has
// duplicated entries.
func MakeSectionsMerged(sections ...*Section) (Sections, error) {
	var ss Sections
	ss = ss.Merge(sections...)
	if err := ss.Validate(); err != nil {
		return nil, errors.Wrap(err, "[element] Sections.Validate")
	}
	return ss, nil
}

// MustMakeSectionsMerged same as MakeSectionsMerged but panics on error.
func MustMakeSectionsMerged(sections ...*Section) Sections {
	s, err := MakeSectionsMerged(sections...)
	if err != nil {
		panic(err)
	}
	return s
}

// TotalFields calculates the total amount of all fields
func (ss Sections) TotalFields() (fs int) {
	for _, s := range ss {
		for _, g := range s.Groups {
			for range g.Fields {
				fs++
			}
		}
	}
	return fs
}

// MergeMultiple merges n SectionSlices into the current slice. Behaviour for
// duplicates: Last item wins. Not thread safe.
func (ss Sections) MergeMultiple(sSlices ...Sections) Sections {
	for _, sl := range sSlices {
		ss = ss.Merge(sl...)
	}
	return ss
}

// Merge merges n Sections into the current slice. Behaviour for duplicates:
// Last item wins. Not thread safe.
func (ss Sections) Merge(sections ...*Section) Sections {
	for _, s := range sections {
		ss = ss.merge(s)
	}
	return ss
}

// Merge copies the data from a Section into this slice. Appends if ID is not
// found in this slice otherwise overrides struct fields if not empty. Not
// thread safe.
func (ss Sections) merge(s *Section) Sections {
	cs, idx, err := ss.Find(s.ID) // cs = current section
	if err != nil {
		ss = append(ss, s)
		idx = len(ss) - 1
	}

	cs.ID = s.ID
	if s.Label != "" {
		cs.Label = s.Label
	}
	if s.Scopes > 0 {
		cs.Scopes = s.Scopes
	}
	if s.SortOrder != 0 {
		cs.SortOrder = s.SortOrder
	}
	if s.Resource > 0 {
		cs.Resource = s.Resource
	}
	cs.Groups = cs.Groups.Merge(s.Groups...)
	ss[idx] = cs
	return ss
}

// Find returns a Section pointer or ErrSectionNotFound. Route must be a single
// part. E.g. if you have path "a/b/c" route would be in this case "a". For
// comparison the field Sum32 of a route will be used. 2nd return parameter
// contains the position of the Section within the Sections. Error
// behaviour: NotFound
func (ss Sections) Find(id string) (*Section, int, error) {
	for i, s := range ss {
		if s.ID == id {
			return s, i, nil
		}
	}
	return nil, 0, errors.NotFound.Newf("[element] Section %q", id)
}

// FindGroup searches for a group using the first two path segments. Route must
// have the format a/b/c. 2nd return parameter contains the position of the
// Group within the GgroupSlice of a Section. Error behaviour: NotFound
func (ss Sections) FindGroup(r string) (*Group, int, error) {
	p := &Path{
		route: r,
	}
	spl, err := p.Split()
	if err != nil {
		return nil, 0, errors.NotFound.Newf("[element] Route %q", r)
	}
	cs, _, err := ss.Find(spl[0])
	if err != nil {
		return nil, 0, errors.Wrap(err, "[element] Sections.FindGroup")
	}
	return cs.Groups.Find(spl[1]) // annotation missing !?
}

// FindField searches for a field using all three path segments. Route must have
// the format a/b/c. Error behaviour: NotFound, NotValid
func (ss Sections) FindField(r string) (*Field, int, error) {
	p := &Path{
		route: r,
	}
	spl, err := p.Split()
	if err != nil {
		return nil, 0, errors.Wrapf(err, "[element] Route %q", r)
	}
	sec, _, err := ss.Find(spl[0])
	if err != nil {
		return nil, 0, errors.Wrapf(err, "[element] Route %q", r)
	}
	cg, _, err := sec.Groups.Find(spl[1])
	if err != nil {
		return nil, 0, errors.Wrapf(err, "[element] Route %q", r)
	}
	return cg.Fields.Find(spl[2]) // annotation missing !?
}

// UpdateField searches for a field using all three path segments and updates
// the found field with the new field data. Not thread safe! Error behaviour:
// NotFound, NotValid
func (ss Sections) UpdateField(r string, nf *Field) error {
	p := &Path{
		route: r,
	}
	spl, err := p.Split()
	if err != nil {
		return errors.Wrapf(err, "[element] Route %q", r)
	}
	sec, sIDX, err := ss.Find(spl[0])
	if err != nil {
		return errors.Wrapf(err, "[element] Route %q", r)
	}
	cg, gIDX, err := sec.Groups.Find(spl[1])
	if err != nil {
		return errors.Wrapf(err, "[element] Route %q", r)
	}
	cf, fIDX, err := cg.Fields.Find(spl[2])
	if err != nil {
		return errors.Wrapf(err, "[element] Route %q", r)
	}

	ss[sIDX].Groups[gIDX].Fields[fIDX] = cf.Update(nf)

	return nil
}

// Append adds 0..n Section. Not thread safe.
func (ss Sections) Append(s ...*Section) Sections {
	return append(ss, s...)
}

// AppendFields adds 0..n *Fields. Path must have at least two path parts like
// a/b more path parts gets ignored. Not thread safe. Error behaviour: NotFound,
// NotValid
func (ss Sections) AppendFields(r string, fs ...*Field) (Sections, error) {
	p := &Path{
		route: r,
	}
	spl, err := p.Split()
	if err != nil {
		return nil, errors.NotFound.Newf("[element] Route %q", r)
	}
	cs, sIDX, err := ss.Find(spl[0])
	if err != nil {
		return nil, errors.Wrapf(err, "[element] Route %q", r)
	}
	cg, gIDX, err := cs.Groups.Find(spl[1])
	if err != nil {
		return nil, errors.Wrapf(err, "[element] Route %q", r)
	}
	cg.Fields = cg.Fields.Append(fs...)
	ss[sIDX].Groups[gIDX] = cg
	return ss, nil
}

// Validate checks for duplicated configuration paths in all three hierarchy
// levels. Error behaviour: NotValid
func (ss Sections) Validate() error {
	if len(ss) == 0 {
		return errors.NotValid.Newf("[element] Sections length is zero")
	}

	dups := make(map[string]bool) // pc path checker
	var buf strings.Builder
	for _, s := range ss {
		for _, g := range s.Groups {
			for _, f := range g.Fields {
				buf.WriteString(s.ID)
				buf.WriteByte(PathSeparator)
				buf.WriteString(g.ID)
				buf.WriteByte(PathSeparator)
				buf.WriteString(f.ID)
				if !dups[buf.String()] {
					dups[buf.String()] = true
				} else {
					return errors.Duplicated.Newf("[config] Within sections the path %q appears two times.", buf.String())
				}
				buf.Reset()
			}
		}
	}
	return nil
}

// SortAll recursively sorts all slices. Not thread safe.
func (ss Sections) SortAll() Sections {
	for _, s := range ss {
		for _, g := range s.Groups {
			g.Fields.Sort()
		}
		s.Groups.Sort()
	}
	return ss.Sort()
}

// Sort convenience helper. Not thread safe.
func (ss Sections) Sort() Sections {
	sort.Sort(ss)
	return ss
}

func (ss Sections) Len() int {
	return len(ss)
}

func (ss Sections) Swap(i, j int) {
	ss[i], ss[j] = ss[j], ss[i]
}

func (ss Sections) Less(i, j int) bool {
	return ss[i].SortOrder < ss[j].SortOrder
}

///////////////////////////////////////////////////////////////////////////////
// GROUPS
///////////////////////////////////////////////////////////////////////////////

// Groups contains a set of Groups.
//  Thread safe for reading but not for modifying.
type Groups []*Group

// Group defines the layout of a group containing multiple Fields
//  Thread safe for reading but not for modifying.
type Group struct {
	// ID unique ID and merged with others. 2nd part of the path.
	ID      string
	Label   string `json:",omitempty"`
	Comment string `json:",omitempty"`
	// Scopes: bit value eg: showInDefault="1" showInWebsite="1" showInStore="1"
	Scopes    scope.Perm `json:",omitempty"`
	SortOrder int        `json:",omitempty"`

	HelpURL  string `json:",omitempty"`
	MoreURL  string `json:",omitempty"`
	DemoLink string `json:",omitempty"`

	Fields Fields `json:",omitempty"`
	// Groups     Groups @todo see recursive options <xs:element name="group"> in app/code/Magento/Config/etc/system_file.xsd
}

// MakeGroups wrapper function, for now.
func MakeGroups(gs ...*Group) Groups {
	return Groups(gs)
}

// Find returns a Group pointer or ErrGroupNotFound. Route must be a single
// part. E.g. if you have path "a/b/c" route would be in this case "b". For
// comparison the field Sum32 of a route will be used. Error behaviour: NotFound
func (gs Groups) Find(id string) (*Group, int, error) {
	for i, g := range gs {
		if g.ID != "" && g.ID == id {
			return g, i, nil
		}
	}
	return nil, 0, errors.NotFound.Newf("[element] Group %q not found", id)
}

// Merge copies the data from a groups into this slice and returns the new
// slice. Appends if ID is not found in this slice otherwise overrides struct
// fields if not empty. Not thread safe.
func (gs Groups) Merge(groups ...*Group) Groups {
	for _, g := range groups {
		gs = gs.merge(g)
	}
	return gs
}

func (gs Groups) merge(g *Group) Groups {
	cg, idx, err := gs.Find(g.ID) // cg current group
	if err != nil {
		cg = g
		gs = append(gs, cg)
		idx = len(gs) - 1
	}

	if "" != g.Label {
		cg.Label = g.Label
	}
	if "" != g.Comment {
		cg.Comment = g.Comment
	}
	if g.Scopes > 0 {
		cg.Scopes = g.Scopes
	}
	if g.SortOrder != 0 {
		cg.SortOrder = g.SortOrder
	}
	cg.Fields = cg.Fields.Merge(g.Fields...)

	gs[idx] = cg
	return gs
}

// Sort convenience helper. Not thread safe.
func (gs Groups) Sort() Groups {
	sort.Sort(gs)
	return gs
}

func (gs Groups) Len() int {
	return len(gs)
}

func (gs Groups) Swap(i, j int) {
	gs[i], gs[j] = gs[j], gs[i]
}

func (gs Groups) Less(i, j int) bool {
	return gs[i].SortOrder < gs[j].SortOrder
}

///////////////////////////////////////////////////////////////////////////////
// FIELDS
///////////////////////////////////////////////////////////////////////////////

// FieldMeta sets meta data for a field into the config.Service object. The meta
// data defines scope access restrictions and default values for different
// scopes.
// Intermediate type for function WithFieldMeta
type FieldMeta struct {
	valid  bool
	Events [eventMaxCount]eventObservers
	// Route defines the route or storage key, e.g.: customer/address/prefix_options
	Route string
	// WriteScopePerm sets the permission to allow setting values to this route.
	// For example WriteScopePerm equals scope.PermStore can be set from
	// default, website and store. If you restrict WriteScopePerm to
	// scope.PermDefault, the route and its value can only be set from default
	// but not from websites or stores. It is only allowed to set WriteScopePerm
	// when ScopeID is zero or scope.DefaultTypeID.
	WriteScopePerm scope.Perm
	// ScopeID defines the scope ID for which a default value is valid. Scope
	// type can only contain three scopes (default,websites or stores). ID
	// relates to the corresponding website or store ID.
	ScopeID      scope.TypeID
	DefaultValid bool
	// Default sets the default value which gets later parsed into the desired
	// final Go type. An empty string means not set or null.
	Default string
}

// Fields contains a set of Fields. Has several method receivers attached.
// Thread safe for reading but not for modifying.
type Fields []*Field

// Field contains the final path element of a configuration. Includes several
// options. Thread safe for reading but not for modifying. @see
// magento2/app/code/Magento/Config/etc/system_file.xsd
type Field struct {
	// ID unique ID and NOT merged with others. 3rd and final part of the path.
	ID string
	// ConfigRoute if provided defines the storage path and overwrites the path from
	// section.id + group.id + field.id. ConfigRoute can be nil.
	ConfigRoute string `json:",omitempty"`
	// Type is used for the front end on how to display a Field
	Type FieldType `json:",omitempty"`
	// Label a short label of the field
	Label string `json:",omitempty"`
	// Comment can contain HTML
	Comment string `json:",omitempty"`
	// Tooltip used for frontend and can contain HTML
	Tooltip string `json:",omitempty"`
	// SortOrder in ascending order
	SortOrder int `json:",omitempty"`
	// Visible used for configuration settings which are not exposed to the user.
	Visible bool `json:",omitempty"`
	// CanBeEmpty only used in HTML forms for multiselect fields
	// Use case: lib/internal/Magento/Framework/Data/Form/Element/Multiselect.php::getElementHtml()
	CanBeEmpty bool `json:",omitempty"`
	// Scopes defines the max allowed scope. Some paths or values can only act
	// on default, website or store scope. So perm checks if the provided
	// path has a scope equal or lower than defined in perm.
	Scopes scope.Perm `json:",omitempty"`
	// Default can contain any default config value: float64, int64, string,
	// bool. An empty string is equal to NULL. A default gets requests if the
	// value for a path cannot be retrieved from Level1 or Level2 storage.
	Default string `json:",omitempty"`
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
	if new.Type > 0 {
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

	f.Visible = new.Visible
	f.CanBeEmpty = new.CanBeEmpty

	if new.Default != "" {
		f.Default = new.Default
	}
	return f
}
