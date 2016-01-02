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

package i18n_test

import "testing"

func TestAllLanguages(t *testing.T) {

	//	d := cldr.Decoder{}
	//	c, err := d.DecodePath("/Users/cys/Downloads/cldr/27.0.1/core/common/")
	//	assert.NoError(t, err)
	//
	//	//	t.Logf("\n%#v\n", c.Locales())
	//
	//	ldml, err := c.LDML("de_DE")
	//	assert.NoError(t, err)
	//
	//	for _, common := range ldml.LocaleDisplayNames.Languages.Language {
	//		t.Logf("%s => ", common.Type)
	//		t.Logf("%s\n", common.CharData)
	//	}

	//	t.Logf("\n%#v\n", ldml.LocaleDisplayNames.Languages)

}

func TestAllCurrencies(t *testing.T) {

	//	for _, bl := range display.Self.Supported.BaseLanguages() {
	//		//		fmt.Print(bl.String(), " ")
	//		fmt.Printf("%s => %s\n ", bl.String(), display.German.Languages().Name(bl))
	//	}

	//	//	fmt.Println("Tags:")
	//	for _, tag := range display.Self.Supported.Tags() {
	//
	//		fmt.Printf("%#v\n", display.German.Languages().Name(tag))
	//
	//		//		r, _ := tag.Region()
	//		//		fmt.Printf("%s => %s => %s\n", tag.String(), r.String())
	//	}

}
