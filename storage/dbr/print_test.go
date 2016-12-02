// Copyright 2015-2016, Cyrill @ Schumacher.fm and the CoreStore contributors
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

package dbr_test

import (
	"fmt"

	"github.com/corestoreio/csfw/storage/dbr"
)

// check if the types implement the interfaces

var _ fmt.Stringer = (*dbr.Delete)(nil)
var _ fmt.Stringer = (*dbr.Insert)(nil)
var _ fmt.Stringer = (*dbr.Update)(nil)
var _ fmt.Stringer = (*dbr.Select)(nil)

var _ dbr.QueryBuilder = (*dbr.Select)(nil)
var _ dbr.QueryBuilder = (*dbr.Delete)(nil)
var _ dbr.QueryBuilder = (*dbr.Update)(nil)
var _ dbr.QueryBuilder = (*dbr.Insert)(nil)
