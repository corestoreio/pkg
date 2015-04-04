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

package eav

import "net/http"

const (
	OutputFormatJson uint8 = iota + 1
	OutputFormatText
	OutputFormatHtml
	OutputFormatPdf
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
	}

	AttributeData struct {
		*Attribute
	}
)

func (AttributeData) ExtractValue(req *http.Request) {}
func (AttributeData) ValidateValue(value []string)   {}
func (AttributeData) CompactValue(value []string)    {}
func (AttributeData) RestoreValue(value []string)    {}
func (AttributeData) OutputValue(format uint8)       {}
