// Copyright 2015-2017, Cyrill @ Schumacher.fm and the CoreStore contributors
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

package eav

import "net/http"

const (
	OutputFormatJSON uint8 = iota + 1
	OutputFormatText
	OutputFormatHTML
	OutputFormatPDF
	OutputFormatOneline
	OutputFormatArray // not sure about that one
)

type (
	//AttributeDataModeller implements methods from magento2/site/app/code/Magento/Eav/Model/Attribute/Data/AbstractData.php
	// All functions: @todo implementation, parameters, returns and if function is really needed.
	AttributeDataModeller interface {

		// ExtractValue Extract data from request and return value
		ExtractValue(req *http.Request)

		// ValidateValue Validate data
		ValidateValue(value []string)

		// CompactValue Export attribute value to entity model
		CompactValue(value []string)

		// RestoreValue Restore attribute value from SESSION to entity model
		RestoreValue(value []string)

		//OutputValue return formatted attribute value from entity model
		OutputValue(format uint8)

		// Config to configure the current instance
		Config(...AttributeDataConfig) AttributeDataModeller
	}

	AttributeData struct {
		a *Attribute
		// idx references to the generated constant and therefore references to itself. mainly used in
		// backend|source|frontend|etc_model
		idx AttributeIndex
	}
	AttributeDataConfig func(*AttributeData)
)

var _ AttributeDataModeller = (*AttributeData)(nil)

// NewAttributeData creates a pointer to a new attribute source
func NewAttributeData(cfgs ...AttributeDataConfig) *AttributeData {
	ad := &AttributeData{
		a: nil,
	}
	ad.Config(cfgs...)
	return ad
}

// AttributeDataIdx only used in generated code to set the current index in the attribute slice
func AttributeDataIdx(i AttributeIndex) AttributeDataConfig {
	return func(as *AttributeData) {
		as.idx = i
	}
}

func (as *AttributeData) Config(configs ...AttributeDataConfig) AttributeDataModeller {
	for _, cfg := range configs {
		cfg(as)
	}
	return as
}

func (AttributeData) ExtractValue(req *http.Request) {}
func (AttributeData) ValidateValue(value []string)   {}
func (AttributeData) CompactValue(value []string)    {}
func (AttributeData) RestoreValue(value []string)    {}
func (AttributeData) OutputValue(format uint8)       {}
