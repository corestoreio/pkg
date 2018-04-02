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
	"encoding/json"
	"sort"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/pkg/config"
	"github.com/corestoreio/pkg/config/cfgpath"
	"github.com/corestoreio/pkg/store/scope"
	"github.com/corestoreio/pkg/util/bufferpool"
)

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
	// Scopes can contain multiple Scope but no more than Default, Website and Store.
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

// Defaults iterates over all slices, creates a path and uses the default value
// to return a map.
func (ss Sections) Defaults() (DefaultMap, error) {
	var dm = make(DefaultMap)
	for _, s := range ss {
		for _, g := range s.Groups {
			for _, f := range g.Fields {
				dm[f.Route(s.ID, g.ID)] = f.Default
			}
		}
	}
	return dm, nil
}

// TotalFields calculates the total amount of all fields
func (ss Sections) TotalFields() int {
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
func (ss Sections) Find(id cfgpath.Route) (*Section, int, error) {
	for i, s := range ss {
		if s.ID.Equal(id) {
			return s, i, nil
		}
	}
	return nil, 0, errors.NotFound.Newf("[element] Section %q", id)
}

// FindGroup searches for a group using the first two path segments. Route must
// have the format a/b/c. 2nd return parameter contains the position of the
// Group within the GgroupSlice of a Section. Error behaviour: NotFound
func (ss Sections) FindGroup(r cfgpath.Route) (*Group, int, error) {

	spl, err := r.Split()
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
func (ss Sections) FindField(r cfgpath.Route) (*Field, int, error) {
	spl, err := r.Split()
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
func (ss Sections) UpdateField(r cfgpath.Route, nf *Field) error {
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
func (ss Sections) AppendFields(r cfgpath.Route, fs ...*Field) (Sections, error) {
	spl, err := r.Split()
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

// ToJSON transforms the whole slice into JSON
func (ss Sections) ToJSON() string {
	buf := bufferpool.Get()
	defer bufferpool.Put(buf)
	if err := json.NewEncoder(buf).Encode(ss); err != nil {
		return "[element] Error: " + err.Error()
	}
	return buf.String()
}

// Validate checks for duplicated configuration paths in all three hierarchy
// levels. Error behaviour: NotValid
func (ss Sections) Validate() error {
	if len(ss) == 0 {
		return errors.NotValid.Newf("[element] Sections length is zero")
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
						return errors.NotValid.Newf("[element] Duplicate entry for path %q :: %s", p.String(), ss.ToJSON())
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

// ConfigurationWriter thread safe storing of configuration values under
// different paths and scopes. This interface has been copied from config.Writer
// to avoid bloated imports. Due to testing both interfaces will be kept in
// sync.
// deprecated this package might get merged into config.
type ConfigurationWriter interface {
	// Write writes a configuration entry and may return an error
	Write(p config.Path, value interface{}) error
}

// ApplyDefaults reads slice Sectioner and applies the keys and values to the
// default configuration writer. Overwrites existing values.
// TODO maybe use a flag to prevent overwriting
func (ss Sections) ApplyDefaults(s ConfigurationWriter) (count int, err error) {
	def, err := ss.Defaults()
	if err != nil {
		return 0, errors.Wrap(err, "[element] Sections.ApplyDefaults.Defaults")
	}
	for k, v := range def {
		var p cfgpath.Path
		p, err = cfgpath.MakeByString(k) // default path!
		if err != nil {
			err = errors.Wrap(err, "[element] Sections.ApplyDefaults.cfgpath.MakeByString")
			return
		}
		if err = s.Write(p, v); err != nil {
			err = errors.Wrap(err, "[element] Sections.ApplyDefaults.Storage.Set")
			return
		}
		count++
	}
	return
}
