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
	"strconv"
	"testing"

	"github.com/corestoreio/pkg/util/assert"
	"github.com/corestoreio/pkg/util/codegen"
)

func TestNewProto(t *testing.T) {
	t.Parallel()

	g := codegen.NewProto("config")
	g.AddImport("github.com/gogo/protobuf/gogoproto/gogo.proto", "")
	g.C("These constants", "are used for testing.", "Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an AS IS BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.")
	g.AddOptions(
		"(gogoproto.typedecl_all)", "false",
		"go_package", "dmltestgenerated",
	)

	g.P("message", "CoreConfigData", "{")
	g.In()
	g.P("uint32", "config_id", "= 1 [(gogoproto.customname)=", strconv.Quote("ConfigID"), "];")
	g.P(`google.protobuf.Timestamp`, `version_te`, `= 8 [(gogoproto.customname)=`, strconv.Quote(`VersionTe`), `,(gogoproto.stdtime)=true,(gogoproto.nullable)=false];`)
	g.Out()
	g.P("}")

	var buf bytes.Buffer
	err := g.GenerateFile(&buf)
	assert.NoError(t, err)

	assert.Exactly(t, `// Auto generated source code
syntax = "proto3";
package config;
import "github.com/gogo/protobuf/gogoproto/gogo.proto";
option (gogoproto.typedecl_all) = dmltestgenerated;
option false = go_package;
// These constants are used for testing. Unless required by applicable law or
// agreed to in writing, software distributed under the License is distributed on
// an AS IS BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express
// or implied. See the License for the specific language governing permissions
// and limitations under the License.
message CoreConfigData { 
	uint32 config_id = 1 [(gogoproto.customname)= "ConfigID" ]; 
	google.protobuf.Timestamp version_te = 8 [(gogoproto.customname)= "VersionTe" ,(gogoproto.stdtime)=true,(gogoproto.nullable)=false]; 
} 
`, buf.String())

}
