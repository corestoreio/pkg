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

package magento

// Version specific type which defines the current installed version of database
// tables.
type Version int

// Version* defines the return values of Version() function.
const (
	Version1 Version = 2 << (iota + 1)
	Version2
	VersionAll = Version1 | Version2
)

var v1 = [4]string{"core_store", "core_website", "core_store_group", "api_user"}
var v2 = [4]string{"integration", "store_website", "store_group", "authorization_role"}

// DetectVersion detects the running version by reading the list of tables. It
// searches for the tables core_store, core_website, core_store_group and
// api_user for Magento v1. It searches for the tables integration,
// store_website, store_group and authorization_role for Magento v2. Prefix is
// the prefix for each table. Returns zero if the version cannot be detected.
func DetectVersion(prefix string, tableList []string) Version {
	var one, two bool
	lv1 := len(v1)
	f1, f2 := 0, 0
	for _, table := range tableList {
		for i := 0; i < lv1; i++ {
			if table == prefix+v1[i] {
				f1++
			}
			if table == prefix+v2[i] {
				f2++
			}
		}

		if f1 == lv1 {
			one = true
		}
		if f2 == lv1 {
			two = true
		}
	}
	switch {
	case one:
		return Version1
	case two:
		return Version2
	default:
		return 0
	}
}
