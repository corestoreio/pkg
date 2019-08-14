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

package search

// TODO JSON and protobuf

//'eq'            => "{{fieldName}} = ?",
//'neq'           => "{{fieldName}} != ?",
//'like'          => "{{fieldName}} LIKE ?",
//'nlike'         => "{{fieldName}} NOT LIKE ?",
//'in'            => "{{fieldName}} IN(?)",
//'nin'           => "{{fieldName}} NOT IN(?)",
//'is'            => "{{fieldName}} IS ?",
//'notnull'       => "{{fieldName}} IS NOT NULL",
//'null'          => "{{fieldName}} IS NULL",
//'gt'            => "{{fieldName}} > ?",
//'lt'            => "{{fieldName}} < ?",
//'gteq'          => "{{fieldName}} >= ?",
//'lteq'          => "{{fieldName}} <= ?",
//'finset'        => "FIND_IN_SET(?, {{fieldName}})",
//'regexp'        => "{{fieldName}} REGEXP ?",
//'from'          => "{{fieldName}} >= ?",
//'to'            => "{{fieldName}} <= ?",
//'seq'           => null,
//'sneq'          => null,
//'ntoa'          => "INET_NTOA({{fieldName}}) LIKE ?",

type Filter struct {
	Field     string
	Value     string
	Condition string // default 'eq'
}

type Filters struct {
	Filters  []Filter
	JoinType string // and/or
}

type SortOrder struct {
	Field  string
	IsDesc bool // Default sort order ascending
}

type Criteria struct {
	FilterGroups []Filters
	SortOrders   []SortOrder
	PageSize     uint
	CurrentPage  uint
}

func (c *Criteria) AddFilters(fs ...Filter) *Criteria {
	c.FilterGroups = append(c.FilterGroups, Filters{
		Filters:  fs,
		JoinType: "???",
	})
	return c
}
