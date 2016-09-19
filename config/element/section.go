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
	"encoding/json"
	"sort"

	"github.com/corestoreio/csfw/config/cfgpath"
	"github.com/corestoreio/csfw/storage/text"
	"github.com/corestoreio/csfw/store/scope"
	"github.com/corestoreio/csfw/util/bufferpool"
	"github.com/corestoreio/csfw/util/errors"
)

// SectionSlice contains a set of Sections. Some nifty helper functions exists.
// Thread safe for reading. A section slice can be used in many goroutines. It
// must remain lock-free.
type SectionSlice []Section

// Section defines the layout for the configuration section which contains
// groups and fields. Thread safe for reading but not for modifying.
type Section struct {
	// ID unique ID and merged with others. 1st part of the path.
	ID    cfgpath.Route
	Label text.Chars `json:",omitempty"`
	// Scopes: bit value eg: showInDefault="1" showInWebsite="1" showInStore="1"
	// Scopes can contain multiple Scope but no more than Default, Website and Store.
	Scopes    scope.Perm `json:",omitempty"`
	SortOrder int        `json:",omitempty"`
	// Resource some kind of ACL if someone is allowed for no,read or write access @todo
	Resource uint `json:",omitempty"`
	Groups   GroupSlice
}

// NewSectionSlice wrapper function, for now.
func NewSectionSlice(s ...Section) SectionSlice {
	return SectionSlice(s)
}

// NewConfiguration creates a new validated SectionSlice with a three level
// configuration.
func NewConfiguration(sections ...Section) (SectionSlice, error) {
	ss := NewSectionSlice(sections...)
	if err := ss.Validate(); err != nil {
		return nil, errors.Wrap(err, "[element] NewSectionSlice.Validate")
	}
	return ss, nil
}

// MustNewConfiguration same as NewConfiguration but panics on error.
func MustNewConfiguration(sections ...Section) SectionSlice {
	s, err := NewConfiguration(sections...)
	if err != nil {
		panic(err)
	}
	return s
}

// NewConfigurationMerge creates a new validated SectionSlice with a three level
// configuration. Before validation, slices are all merged together. Panics if a
// path is redundant. Only use this function if your package elementuration
// really has duplicated entries.
func NewConfigurationMerge(sections ...Section) (SectionSlice, error) {
	var ss SectionSlice
	if err := ss.Merge(sections...); err != nil {
		return nil, errors.Wrap(err, "[element] SectionSlice.Merge")
	}
	if err := ss.Validate(); err != nil {
		return nil, errors.Wrap(err, "[element] SectionSlice.Validate")
	}
	return ss, nil
}

// MustNewConfigurationMerge same as NewConfigurationMerge but panics on error.
func MustNewConfigurationMerge(sections ...Section) SectionSlice {
	s, err := NewConfigurationMerge(sections...)
	if err != nil {
		panic(err)
	}
	return s
}

// Defaults iterates over all slices, creates a path and uses the default value
// to return a map.
func (ss SectionSlice) Defaults() (DefaultMap, error) {
	var dm = make(DefaultMap)
	for _, s := range ss {
		for _, g := range s.Groups {
			for _, f := range g.Fields {
				r, err := f.Route(s.ID, g.ID)
				if err != nil {
					return nil, errors.Wrapf(err, "[element] SectionSlice.Defaults.Field.Route. Section %q Group %q", s.ID, g.ID)
				}
				dm[r.String()] = f.Default
			}
		}
	}
	return dm, nil
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

// MergeMultiple merges n SectionSlices into the current slice. Behaviour for
// duplicates: Last item wins. Not thread safe.
func (ss *SectionSlice) MergeMultiple(sSlices ...SectionSlice) error {
	for _, sl := range sSlices {
		if err := ss.Merge(sl...); err != nil {
			return errors.Wrap(err, "[element] SectionSlice.MergeMultiple")
		}
	}
	return nil
}

// Merge merges n Sections into the current slice. Behaviour for duplicates:
// Last item wins. Not thread safe.
func (ss *SectionSlice) Merge(sections ...Section) error {
	for _, s := range sections {
		if err := ss.merge(s); err != nil {
			return errors.Wrap(err, "[element] SectionSlice.merge")
		}
	}
	return nil
}

// Merge copies the data from a Section into this slice. Appends if ID is not
// found in this slice otherwise overrides struct fields if not empty. Not
// thread safe.
func (ss *SectionSlice) merge(s Section) error {
	cs, idx, err := (*ss).Find(s.ID) // cs current section
	if err != nil {
		(*ss) = append(*ss, s)
		idx = len(*ss) - 1
	}

	cs.ID = s.ID.Clone()
	if s.Label.IsEmpty() == false {
		cs.Label = s.Label.Clone()
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
	if err := cs.Groups.Merge(s.Groups...); err != nil {
		return errors.Wrap(err, "[element] SectionSlice.merge.Groups.Merge")
	}

	(*ss)[idx] = cs
	return nil
}

// Find returns a Section pointer or ErrSectionNotFound. Route must be a single
// part. E.g. if you have path "a/b/c" route would be in this case "a". For
// comparison the field Sum32 of a route will be used. 2nd return parameter
// contains the position of the Section within the SectionSlice. Error
// behaviour: NotFound
func (ss SectionSlice) Find(id cfgpath.Route) (Section, int, error) {
	for i, s := range ss {
		if s.ID.Sum32 > 0 && s.ID.Sum32 == id.Sum32 {
			return s, i, nil
		}
	}
	return Section{}, 0, errors.NewNotFoundf("[element] Section %q", id)
}

// FindGroup searches for a group using the first two path segments. Route must
// have the format a/b/c. 2nd return parameter contains the position of the
// Group within the GgroupSlice of a Section. Error behaviour: NotFound
func (ss SectionSlice) FindGroup(r cfgpath.Route) (Group, int, error) {

	spl, err := r.Split()
	if err != nil {
		return Group{}, 0, errors.NewNotFoundf("[element] Route %q", r)
	}
	cs, _, err := ss.Find(spl[0])
	if err != nil {
		return Group{}, 0, errors.Wrap(err, "[element] SectionSlice.FindGroup")
	}
	return cs.Groups.Find(spl[1]) // annotation missing !?
}

// FindField searches for a field using all three path segments. Route must have
// the format a/b/c. Error behaviour: NotFound, NotValid
func (ss SectionSlice) FindField(r cfgpath.Route) (Field, int, error) {
	spl, err := r.Split()
	if err != nil {
		return Field{}, 0, errors.Wrapf(err, "[element] Route %q", r)
	}
	sec, _, err := ss.Find(spl[0])
	if err != nil {
		return Field{}, 0, errors.Wrapf(err, "[element] Route %q", r)
	}
	cg, _, err := sec.Groups.Find(spl[1])
	if err != nil {
		return Field{}, 0, errors.Wrapf(err, "[element] Route %q", r)
	}
	return cg.Fields.Find(spl[2]) // annotation missing !?
}

// UpdateField searches for a field using all three path segments and updates
// the found field with the new field data. Not thread safe! Error behaviour:
// NotFound, NotValid
func (ss SectionSlice) UpdateField(r cfgpath.Route, new Field) error {
	spl, err := r.Split()
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

	ss[sIDX].Groups[gIDX].Fields[fIDX] = cf.Update(new)

	return nil
}

// Append adds 0..n Section. Not thread safe.
func (ss *SectionSlice) Append(s ...Section) *SectionSlice {
	*ss = append(*ss, s...)
	return ss
}

// AppendFields adds 0..n *Fields. Path must have at least two path parts like
// a/b more path parts gets ignored. Not thread safe. Error behaviour: NotFound,
// NotValid
func (ss *SectionSlice) AppendFields(r cfgpath.Route, fs ...Field) error {
	spl, err := r.Split()
	if err != nil {
		return errors.NewNotFoundf("[element] Route %q", r)
	}
	cs, sIDX, err := ss.Find(spl[0])
	if err != nil {
		return errors.Wrapf(err, "[element] Route %q", r)
	}
	cg, gIDX, err := cs.Groups.Find(spl[1])
	if err != nil {
		return errors.Wrapf(err, "[element] Route %q", r)
	}
	cg.Fields.Append(fs...)
	(*ss)[sIDX].Groups[gIDX] = cg
	return nil
}

// ToJSON transforms the whole slice into JSON
func (ss SectionSlice) ToJSON() string {
	buf := bufferpool.Get()
	defer bufferpool.Put(buf)
	if err := json.NewEncoder(buf).Encode(ss); err != nil {
		return "[element] Error: " + err.Error()
	}
	return buf.String()
}

// Validate checks for duplicated configuration paths in all three hierarchy
// levels. Error behaviour: NotValid
func (ss SectionSlice) Validate() error {
	if len(ss) == 0 {
		return errors.NewNotValidf("[element] SectionSlice length is zero")
	}

	var hashes = make([]uint64, ss.TotalFields(), ss.TotalFields()) // pc path checker

	i := 0
	for _, s := range ss {
		for _, g := range s.Groups {
			for _, f := range g.Fields {

				fnv1a, err := f.RouteHash(s.ID, g.ID)
				if err != nil {
					return errors.Wrapf(err, "[element] Route Section %q Group %q", s.ID, g.ID)
				}

				for _, h := range hashes {
					if h == fnv1a {
						p, err := f.Route(s.ID, g.ID)
						if err != nil {
							return errors.Wrapf(err, "[element] Route Section %q Group %q", s.ID, g.ID)
						}
						return errors.NewNotValidf("[element] Duplicate entry for path %q :: %s", p.String(), ss.ToJSON())
					}
				}
				hashes[i] = fnv1a
				i++
			}
		}
	}
	return nil
}

// SortAll recursively sorts all slices. Not thread safe.
func (ss SectionSlice) SortAll() SectionSlice {
	for _, s := range ss {
		for _, g := range s.Groups {
			g.Fields.Sort()
		}
		s.Groups.Sort()
	}
	return ss.Sort()
}

// Sort convenience helper. Not thread safe.
func (ss SectionSlice) Sort() SectionSlice {
	sort.Sort(ss)
	return ss
}

func (ss SectionSlice) Len() int {
	return len(ss)
}

func (ss SectionSlice) Swap(i, j int) {
	ss[i], ss[j] = ss[j], ss[i]
}

func (ss SectionSlice) Less(i, j int) bool {
	return ss[i].SortOrder < ss[j].SortOrder
}

// ConfigurationWriter thread safe storing of configuration values under
// different paths and scopes. This interface has been copied from config.Writer
// to avoid bloated imports. Due to testing both interfaces will be kept in
// sync.
type ConfigurationWriter interface {
	// Write writes a configuration entry and may return an error
	Write(p cfgpath.Path, value interface{}) error
}

// ApplyDefaults reads slice Sectioner and applies the keys and values to the
// default configuration writer. Overwrites existing values.
// TODO maybe use a flag to prevent overwriting
func (ss SectionSlice) ApplyDefaults(s ConfigurationWriter) (count int, err error) {
	def, err := ss.Defaults()
	if err != nil {
		return 0, errors.Wrap(err, "[element] SectionSlice.ApplyDefaults.Defaults")
	}
	for k, v := range def {
		var p cfgpath.Path
		p, err = cfgpath.NewByParts(k) // default path!
		if err != nil {
			err = errors.Wrap(err, "[element] SectionSlice.ApplyDefaults.cfgpath.NewByParts")
			return
		}
		if err = s.Write(p, v); err != nil {
			err = errors.Wrap(err, "[element] SectionSlice.ApplyDefaults.Storage.Set")
			return
		}
		count++
	}
	return
}
