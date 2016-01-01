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

package element

import (
	"bytes"
	"encoding/json"
	"errors"
	"sort"

	"github.com/corestoreio/csfw/store/scope"
	"github.com/corestoreio/csfw/util"
	"github.com/juju/errgo"
)

// ErrSectionNotFound error when a section cannot be found.
var ErrSectionNotFound = errors.New("Section not found")

type (
	// SectionSlice contains a set of Sections. Some nifty helper functions
	// exists. Thread safe for reading. A section slice can be used in many
	// goroutines. It must remain lock-free.
	SectionSlice []*Section
	// Section defines the layout for the configuration section which contains groups and fields.
	Section struct {
		// ID unique ID and merged with others. 1st part of the path.
		ID    string
		Label string `json:",omitempty"`
		// Scope: bit value eg: showInDefault="1" showInWebsite="1" showInStore="1"
		Scope     scope.Perm `json:",omitempty"`
		SortOrder int        `json:",omitempty"`
		// Resource some kind of ACL if someone is allowed for no,read or write access @todo
		Resource uint `json:",omitempty"`
		Groups   GroupSlice
	}
)

var _ Sectioner = (*SectionSlice)(nil)

// NewConfiguration creates a new validated SectionSlice with a three level configuration.
// Panics if a path is redundant.
func NewConfiguration(sections ...*Section) (SectionSlice, error) {
	ss := SectionSlice(sections)
	if err := ss.Validate(); err != nil {
		if PkgLog.IsDebug() {
			PkgLog.Debug("config.NewConfiguration.Validate", "err", err)
		}
		return nil, errgo.Mask(err)
	}
	return ss, nil
}

// MustNewConfiguration same as NewConfiguration but panics on error.
func MustNewConfiguration(sections ...*Section) SectionSlice {
	s, err := NewConfiguration(sections...)
	if err != nil {
		panic(err)
	}
	return s
}

// NewConfigurationMerge creates a new validated SectionSlice with a three level configuration.
// Before validation, slices are all merged together. Panics if a path is redundant.
// Only use this function if your package elementuration really has duplicated entries.
func NewConfigurationMerge(sections ...*Section) (SectionSlice, error) {
	var ss SectionSlice
	if err := ss.Merge(sections...); err != nil {
		if PkgLog.IsDebug() {
			PkgLog.Debug("config.NewConfigurationMerge.Merge", "err", err, "sections", sections)
		}
		return nil, errgo.Mask(err)
	}
	if err := ss.Validate(); err != nil {
		if PkgLog.IsDebug() {
			PkgLog.Debug("config.NewConfigurationMerge.Validate", "err", err)
		}
		return nil, errgo.Mask(err)
	}
	return ss, nil
}

// MustNewConfigurationMerge same as NewConfigurationMerge but panics on error.
func MustNewConfigurationMerge(sections ...*Section) SectionSlice {
	s, err := NewConfigurationMerge(sections...)
	if err != nil {
		panic(err)
	}
	return s
}

// Defaults iterates over all slices, creates a path and uses the default value
// to return a map.
func (ss SectionSlice) Defaults() DefaultMap {
	var dm = make(DefaultMap)
	for _, s := range ss {
		for _, g := range s.Groups {
			for _, f := range g.Fields {
				dm[f.FQPathDefault(s.ID, g.ID)] = f.Default
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
			for range g.Fields {
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
	if s.Resource > 0 {
		cs.Resource = s.Resource
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
		paths = scope.PathSplit(paths[0])
	}
	if len(paths) < 2 {
		return nil, ErrGroupNotFound
	}
	cs, err := ss.FindByID(paths[0])
	if err != nil {
		return nil, errgo.Mask(err)
	}
	return cs.Groups.FindByID(paths[1])
}

// FindFieldByPath searches for a field using all three path segments.
// If one argument is given then considered as the full path e.g. a/b/c
// If three arguments are given then each argument will be treated as a path part.
func (ss SectionSlice) FindFieldByPath(paths ...string) (*Field, error) {
	if len(paths) == 1 {
		paths = scope.PathSplit(paths[0])
	}
	if len(paths) < 3 {
		return nil, ErrFieldNotFound
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

// AppendFields adds 0..n *Fields. Path must have at least two path parts like a/b
// more path parts gets ignored.
func (ss SectionSlice) AppendFields(path string, fs ...*Field) error {
	g, err := ss.FindGroupByPath(path)
	if err != nil {
		return errgo.Mask(err)
	}
	g.Fields.Append(fs...)
	return nil
}

// ToJSON transforms the whole slice into JSON
func (ss SectionSlice) ToJSON() string {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(ss); err != nil {
		PkgLog.Debug("config.SectionSlice.ToJSON.Encode", "err", err)
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
	var pc = make(util.StringSlice, ss.TotalFields()) // pc path checker
	i := 0
	for _, s := range ss {
		for _, g := range s.Groups {
			for _, f := range g.Fields {
				p := f.FQPathDefault(s.ID, g.ID)
				if pc.Include(p) {
					return errgo.Newf("Duplicate entry for path %s :: %s", p, ss.ToJSON())
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
