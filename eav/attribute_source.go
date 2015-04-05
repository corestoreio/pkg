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

type (
	// AttributeSourceModeller interface implements the functions needed to retrieve data from a source model
	// @see magento2/site/app/code/Magento/Eav/Model/Entity/Attribute/Source/SourceInterface.php and
	// its abstract class plus other default implementations. Refers to tables eav_attribute_option and
	// eav_attribute_option_value OR other tables. @todo
	AttributeSourceModeller interface {
		// GetAllOptions returns all options in a value/label slice
		GetAllOptions() AttributeSourceOptions
		// GetOptionText returns for a value the appropriate label
		// @see magento2/site/app/code/Magento/Eav/Model/Entity/Attribute/Source/AbstractSource.php
		GetOptionText(value string) string
		// Config to configure the current instance
		Config(...AttributeSourceConfig) AttributeSourceModeller
	}
	// AttributeSource should implement all abstract ideas of
	// @see magento2/site/app/code/Magento/Eav/Model/Entity/Attribute/Source/AbstractSource.php
	// maybe extend also the interface
	AttributeSource struct {
		// a is the reference to the already created attribute during init() call in a package.
		// Do not modify the attribute here
		a *Attribute
		// Source is the internal source []string where i%0 is the value and i%1 is the label.
		Source []string
		// o is the internal option cache
		o AttributeSourceOptions
		// idx references to the generated constant and therefore references to itself. mainly used in
		// backend|source|frontend|etc_model
		idx AttributeIndex
	}
	AttributeSourceConfig func(*AttributeSource)

	// AttributeSourceOptions is a slice of AttributeSourceOption structs
	AttributeSourceOptions []AttributeSourceOption
	// AttributeSourceOption contains a value and label mostly for output in the frontend
	AttributeSourceOption struct {
		// Value can be any value and is now here a temporary string. @todo check if maybe interface is needed
		Value string
		// Label is the name of a value
		Label string
	}
)

var _ AttributeSourceModeller = (*AttributeSource)(nil)

// NewAttributeSource creates a pointer to a new attribute source
func NewAttributeSource(cfgs ...AttributeSourceConfig) *AttributeSource {
	as := &AttributeSource{
		a: nil,
	}
	as.Config(cfgs...)
	return as
}

// AttributeSourceIdx only used in generated code to set the current index in the attribute slice
func AttributeSourceIdx(i AttributeIndex) AttributeSourceConfig {
	return func(as *AttributeSource) {
		as.idx = i
	}
}

// Config runs the configuration functions
func (as *AttributeSource) Config(configs ...AttributeSourceConfig) AttributeSourceModeller {
	for _, cfg := range configs {
		cfg(as)
	}
	return as
}

// GetAllOptions returns an option slice
func (as *AttributeSource) GetAllOptions() AttributeSourceOptions {
	if (len(as.Source) & 1) != 0 { // modulus %2
		as.o = AttributeSourceOptions{
			AttributeSourceOption{
				Label: "Incorrect source. Uneven.",
			},
		}
	}
	if len(as.o) > 0 {
		return as.o
	}

	as.o = make(AttributeSourceOptions, len(as.Source)/2, len(as.Source)/2)
	j := 0
	for i := 0; i < len(as.Source); i = i + 2 {
		as.o[j] = AttributeSourceOption{
			Value: as.Source[i],
			Label: as.Source[i+1],
		}
		j++
	}
	return as.o
}

// GetOptionText returns for a value v the label
func (as *AttributeSource) GetOptionText(v string) string { return as.o.label(v) }

func (os AttributeSourceOptions) label(v string) string {
	for _, o := range os {
		if o.Value == v {
			if o.Label != "" {
				return o.Label
			}
			return o.Value
		}
	}
	return ""
}
