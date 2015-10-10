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

package configsource

import "github.com/corestoreio/csfw/config/valuelabel"

// YesNo defines a slice with yes and no options.
var YesNo = valuelabel.NewByBool(valuelabel.Bools{
	{false, "No"},
	{true, "Yes"},
})

// EnableDisable defines a slice with enable and disable options.
var EnableDisable = valuelabel.NewByBool(valuelabel.Bools{
	{false, "Disable"},
	{true, "Enable"},
})
