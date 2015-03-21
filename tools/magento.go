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

package tools

type (
	// mageModelMap key refers to a *_model column and value is a JSON default mapping
	mageModelMap map[int][]byte
)

const (
	// EavAttributeBackendModel relates to table column eav_attribute.backend_model
	EavAttributeBackendModel int = iota + 1
	// EavAttributeFrontendModel relates to table column eav_attribute.frontend_model
	EavAttributeFrontendModel
	// EavAttributeSourceModel relates to table column eav_attribute.source_model
	EavAttributeSourceModel
)
