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

package utils

const (
	SortAsc SortDirection = 1 << iota
	SortDesc
)
const (
	_SortDirection_name_0 = "ASC"
	_SortDirection_name_1 = "DESC"
)

type SortDirection uint8

func (i SortDirection) String() string {
	switch {
	case i == 1:
		return _SortDirection_name_0
	case i == 2:
		return _SortDirection_name_1
	default:
		return _SortDirection_name_0
	}
}
