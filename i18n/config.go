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

package i18n

import (
	"github.com/juju/errgo"
	"golang.org/x/text/cldr"
)

const (
	// LocaleDefault is the overall default locale when no website/store setting is available.
	LocaleDefault = "en_US"
	// CLDRVersionRequired required version to run this package
	CLDRVersionRequired = "27.0.1"
)

func init() {
	if cldr.Version != CLDRVersionRequired {
		panic(errgo.Newf("Incorrect CLDR Version! Expecting %s but got %s. Please check golang.org/x/text/cldr", CLDRVersionRequired, cldr.Version))
	}
}
