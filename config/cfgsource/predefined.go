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

package cfgsource

// YesNo defines an immutable slice with yes and no options.
var YesNo = NewByBool(Bools{
	{false, "No"},
	{true, "Yes"},
})

// EnableDisable defines an immutable slice with enable and disable options.
var EnableDisable = NewByBool(Bools{
	{false, "Disable"},
	{true, "Enable"},
})

// Protocol defines an immutable slice with available HTTP protocols
var Protocol = MustNewByString(
	"", "",
	"http", "HTTP (unsecure)",
	"https", "HTTP (TLS)",
)
