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

// Generates code for all attribute types
package main

const tplEav = `// Copyright 2015 CoreStore Authors
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

// Package {{ .Package }} file is auto generated via eavToStruct
package {{ .Package }}
import (

)

`

var defaultMapping = []byte(`{
  "eav/entity_attribute_backend_datetime": {
    "import_path": "github.com/corestoreio/csfw/eav",
    "backend_model": "eav.Attribute().Backend().Datetime()",
  }
  "catalog/product_attribute_backend_price": {
    "import_path": "github.com/corestoreio/csfw/catalog",
    "backend_model": "catalog.Attribute().Backend().Price()",
  }
  SELECT COUNT(*) AS Rows, backend_model FROM eav_attribute GROUP BY backend_model ORDER BY Rows desc
  rethink that ...


    "frontend_model": "customer.Attribute()",
    "source_model": "customer.Attribute()",
}`)
