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
	"sort"

	"github.com/corestoreio/csfw/utils"
)

type (
	ValueLabelSlice []ValueLabel
	vlSortByLabel   struct {
		ValueLabelSlice
	}
	vlSortByValue struct {
		ValueLabelSlice
	}

	// Option type is returned by the SourceModel interface
	ValueLabel struct {
		Value, Label string
	}

	// ModelConstructor implements different fields which can be differently used
	// by the FieldSourceModeller or FieldBackendModeller types.
	// Each implementation is on its own responsible to check for the required fields.
	modelConstructor struct {
		ScopeID      ScopeIDer
		ConfigReader Reader
		// more fields to be added ...
	}

	// ModelArgFunc function to set the fields in ModelConstructor
	ModelArgFunc func(*modelConstructor)

	// FieldSourceModeller defines how to retrieve all option values. Mostly used for frontend output.
	FieldSourceModeller interface {
		Construct(...ModelArgFunc)
		Options() ValueLabelSlice
	}

	// FieldBackendModeller defines how to save and load the data @todo think about AddData
	// In Magento slang: beforeSave() and afterLoad()
	FieldBackendModeller interface {
		Construct(...ModelArgFunc)
		AddData(interface{})
		Save() error
	}
)

func ModelScope(s ScopeIDer) ModelArgFunc {
	return func(c *modelConstructor) {
		c.ScopeID = s
	}
}

func ModelConfig(s ScopeIDer) ModelArgFunc { @todo
	return func(c *modelConstructor) {
		c.ScopeID = s
	}
}

// SortByLabel sorts by label in asc or desc direction
func (s ValueLabelSlice) SortByLabel(d utils.SortDirection) ValueLabelSlice {
	fsv := vlSortByLabel{s}
	if d == utils.SortDesc {
		fsv = sort.Reverse(fsv)
	}
	sort.Sort(fsv)
	return s
}

// SortByValue sorts by label in asc or desc direction
func (s ValueLabelSlice) SortByValue(d utils.SortDirection) ValueLabelSlice {
	fsv := vlSortByValue{s}
	if d == utils.SortDesc {
		fsv = sort.Reverse(fsv)
	}
	sort.Sort(fsv)
	return s
}

func (s ValueLabelSlice) Len() int {
	return len(s)
}

func (s ValueLabelSlice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (l vlSortByLabel) Less(i, j int) bool {
	return l.ValueLabelSlice[i].Label < l.ValueLabelSlice[j].Label
}

func (v vlSortByValue) Less(i, j int) bool {
	return v.ValueLabelSlice[i].Value < v.ValueLabelSlice[j].Value
}
