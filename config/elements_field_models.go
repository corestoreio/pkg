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
	// ValueLabelSlice type is returned by the SourceModel.Options() interface
	ValueLabelSlice []ValueLabel
	vlSortByLabel   struct {
		ValueLabelSlice
	}
	vlSortByValue struct {
		ValueLabelSlice
	}

	// ValueLabel contains a stringyfied value and a label for print in a browser/ JS api.
	ValueLabel struct {
		Value, Label string
	}

	// ModelConstructor implements different fields/functions which can be differently used
	// by the FieldSourceModeller or FieldBackendModeller types.
	// Nearly all functions will return not nil. The Construct() function takes what it needs.
	ModelConstructor struct {
		// Scope contains a website/store ID or nil (=default scope)
		Scope ScopeIDer
		// ConfigReader returns the configuration reader and never nil
		ConfigReader Reader
		// @todo more fields to be added, depends on the overall requirements of all Magento models.
	}

	// FieldSourceModeller defines how to retrieve all option values. Mostly used for frontend output.
	FieldSourceModeller interface {
		Construct(ModelConstructor)
		Options() ValueLabelSlice
	}

	// FieldBackendModeller defines how to save and load the data @todo rethink AddData
	// In Magento slang: beforeSave() and afterLoad()
	FieldBackendModeller interface {
		Construct(ModelConstructor)
		AddData(interface{})
		Save() error
	}
)

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
