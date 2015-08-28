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

package utils

import "strings"

func isMagento(tableList []string, proof [4]string) bool {
	lp := len(proof)
	var found int
	for _, table := range tableList {
		for _, check := range proof {
			if strings.Contains(table, check) {
				found++
			}
		}
		if found == lp {
			return true
		}
	}
	return false
}

// IsMagento1 detects by reading the list of tables which Magento version you
// are running. It searches for the tables core_store, core_website,
// core_store_group and api_user.
func IsMagento1(tableList []string) bool {
	return isMagento(tableList, [4]string{"core_store", "core_website", "core_store_group", "api_user"})
}

// IsMagento2 detects by reading the list of tables which Magento version you
// are running. It searches for the tables store", store_website, store_group
// and authorization_role.
func IsMagento2(tableList []string) bool {
	return isMagento(tableList, [4]string{"store", "store_website", "store_group", "authorization_role"})
}
