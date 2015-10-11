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

import (
	"net/http"

	"github.com/corestoreio/csfw/config/valuelabel"
)

// Redirect defines a slice with different redirect codes
var Redirect = valuelabel.NewByInt(valuelabel.Ints{
	{0, "No"},
	{1, "Yes (302 Found)"},                // old from Magento
	{http.StatusFound, "Yes (302 Found)"}, // new correct
	{http.StatusMovedPermanently, "Yes (301 Moved Permanently)"},
})

// Protocol defines a slice with available HTTP protocols
var Protocol = valuelabel.NewByString(
	"", "",
	"http", "HTTP (unsecure)",
	"https", "HTTP (TLS)",
)
