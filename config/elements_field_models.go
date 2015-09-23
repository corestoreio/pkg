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

package config

import (
	"bytes"
	"encoding/json"
	"sort"

	"github.com/corestoreio/csfw/config/scope"
	"github.com/corestoreio/csfw/utils"
	"github.com/corestoreio/csfw/utils/log"
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

	// ValueLabel contains a stringyfied value and a label for printing in a browser / JS api.
	ValueLabel struct {
		Value, Label string
	}

	// ModelConstructor implements different fields/functions which can be differently used
	// by the FieldSourceModeller or FieldBackendModeller types.
	// Nearly all functions will return not nil. The Construct() function takes what it needs.
	ModelConstructor struct {
		// WebsiteID contains a website/store ID or nil (=default scope) both can be nil or just one
		WebsiteID scope.WebsiteIDer
		StoreID   scope.StoreIDer
		// ConfigReader returns the configuration reader and never nil
		ConfigReader Reader
		// @todo more fields to be added, depends on the overall requirements of all Magento models.
	}

	// FieldSourceModeller defines how to retrieve all option values. Mostly used for frontend output.
	// The Construct() must be used because NOT all fields of ModelConstructor are available during
	// init process and can of course change during the running app. Also to prevent circular dependencies.
	FieldSourceModeller interface {
		Construct(ModelConstructor) error
		Options() ValueLabelSlice
	}

	// FieldBackendModeller defines how to save and load the data @todo rethink AddData
	// In Magento slang: beforeSave() and afterLoad().
	// The Construct() must be used because NOT all fields of ModelConstructor are available during
	// init process and can of course change during the running app. Also to prevent circular dependencies.
	FieldBackendModeller interface {
		Construct(ModelConstructor) error
		AddData(interface{})
		Save() error
	}
)

// SortByLabel sorts by label in asc or desc direction
func (s ValueLabelSlice) SortByLabel(d utils.SortDirection) ValueLabelSlice {
	var si sort.Interface
	si = vlSortByLabel{s}
	if d == utils.SortDesc {
		si = sort.Reverse(si)
	}
	sort.Sort(si)
	return s
}

// SortByValue sorts by label in asc or desc direction
func (s ValueLabelSlice) SortByValue(d utils.SortDirection) ValueLabelSlice {
	var si sort.Interface
	si = vlSortByValue{s}
	if d == utils.SortDesc {
		si = sort.Reverse(si)
	}
	sort.Sort(si)
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

// ToJSON returns a JSON string
func (s ValueLabelSlice) ToJSON() string {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(s); err != nil {
		log.Error("config.ValueLabelSlice.ToJSON.Encode", "err", err)
		return ""
	}
	return buf.String()
}
