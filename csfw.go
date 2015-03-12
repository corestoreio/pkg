//go:generate csTableToStruct -p eav -prefixSearch eav -o eav/generated_tables.go -run
//go:generate go install github.com/corestoreio/csfw/tools/csEavToStruct
//go:generate csEavToStruct -p eav -o eav/generated_eav.go -run

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

package csfw
