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

package toolsdb

import (
	"log"
	"strings"

	"github.com/juju/errgo"
)

// Camelize transforms from snake case to camelCase e.g. catalog_product_id to CatalogProductID.
func Camelize(s string) string {
	parts := strings.Split(s, "_")
	ret := ""
	for _, p := range parts {
		switch p {
		case "id":
			p = "ID"
			break
		}
		ret = ret + strings.Title(p)
	}
	return ret
}

// LogFatal logs an error as fatal with printed location and exists the program.
func LogFatal(err error) {
	if err == nil {
		return
	}
	s := "Error: " + err.Error()
	if err, ok := err.(errgo.Locationer); ok {
		s += " " + err.Location().String()
	}
	log.Fatalln(s)
}
