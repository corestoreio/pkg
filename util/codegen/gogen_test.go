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

package codegen_test

import (
	"bytes"
	"testing"

	"github.com/corestoreio/pkg/util/assert"
	"github.com/corestoreio/pkg/util/codegen"
)

func TestEncloseBT(t *testing.T) {
	assert.Exactly(t, "`a`", codegen.EncloseBT(`a`))
}

func TestNewGo(t *testing.T) {
	t.Parallel()

	g := codegen.NewGo("config")
	g.BuildTags = []string{"ignoring"}
	g.AddImport("fmt", "")
	g.AddImport("github.com/corestoreio/pkg/storage/null", "null")
	g.C("These constants", "are used for testing.", "Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an AS IS BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.")
	g.AddConstString("TableA", "table_a")

	g.Pln("type", "CatalogProductEntity", "struct {")
	g.In()
	g.Pln("EntityID", "int64")
	g.Pln("StoreID", "uint32", `// store_id smallint(5) unsigned NOT NULL PRI   "Store ID"`)
	g.Pln("Value", "null.Decimal", `// value decimal(12,4) NOT NULL PRI   "Value"`)
	g.Out()
	g.P(`// Hello World`)
	g.Pln("\n}")

	var buf bytes.Buffer
	err := g.GenerateFile(&buf)
	assert.NoError(t, err)

	assert.Exactly(t, `// +build ignoring

package config

// Auto generated source code
import (
	"fmt"
	null "github.com/corestoreio/pkg/storage/null"
)

const (
	TableA = "table_a"
)

// These constants are used for testing. Unless required by applicable law or
// agreed to in writing, software distributed under the License is distributed on
// an AS IS BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express
// or implied. See the License for the specific language governing permissions
// and limitations under the License.
type CatalogProductEntity struct {
	EntityID int64
	StoreID  uint32       // store_id smallint(5) unsigned NOT NULL PRI   "Store ID"
	Value    null.Decimal // value decimal(12,4) NOT NULL PRI   "Value"
	// Hello World
}
`, buf.String())

}
