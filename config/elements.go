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

const (
	TypeSelect FieldType = iota + 1
	TypeMultiSelect
	TypeText
	TypeObscure
)

type (
	Option struct {
		Value, Label string
	}

	SectionSlice []*Section
	Section      struct {
		ID    string
		Label string
		// Scope: bit value eg: showInDefault="1" showInWebsite="1" showInStore="1"
		Scope     ScopeBits
		SortOrder int
		// Permission some kind of ACL if some is allowed for read or write access
		Permission uint
		Groups     GroupSlice
	}

	GroupSlice []*Group
	Group      struct {
		ID      string
		Label   string
		Comment string
		// Scope: bit value eg: showInDefault="1" showInWebsite="1" showInStore="1"
		Scope     ScopeBits
		SortOrder int
		Fields    FieldSlice
	}

	FieldType uint

	FieldSlice []*Field
	// Element @see magento2/app/code/Magento/Config/etc/system_file.xsd @todo add later more fields IF necessary
	Field struct {
		ID      string
		Type    FieldType
		Label   string
		Comment string
		// Scope: bit value eg: showInDefault="1" showInWebsite="1" showInStore="1"
		Scope     ScopeBits
		SortOrder int
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

// DefaultConfiguration iterates over all slices, creates a path and uses the default value
// to return a map
func (s SectionSlice) DefaultConfiguration() DefaultMap {
	return nil
}
