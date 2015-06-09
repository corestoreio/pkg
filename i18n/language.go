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
	"fmt"

	"golang.org/x/text/display"
	"golang.org/x/text/language"
)

// GetLanguages returns a list of languages as a key/value slice. Odd index/key = locale,
// even index/value = Humanized readable string. The humanized strings contains the language
// name in its language and language name in requested tag
func GetLanguages(t language.Tag) []string {
	var ret = make([]string, len(tags)*2)
	n := getDict(t)
	i := 0
	for _, t := range tags {
		b, _ := t.Base()
		r, _ := t.Region()
		ret[i] = GetLocale(b, r)
		ret[i+1] = fmt.Sprintf("%-20s (%s)", display.Self.Name(t), n.Languages().Name(t))
		i = i + 2
	}
	return ret
}

func GetAllLanguages() {
	//	for _, bl := range display.Self.Supported.BaseLanguages() {
	//		//		fmt.Print(bl.String(), " ")
	//		fmt.Printf("%s => %s\n ", bl.String(), display.German.Languages().Name(bl))
	//	}

}
