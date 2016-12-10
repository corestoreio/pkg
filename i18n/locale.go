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

package i18n

import (
	"github.com/corestoreio/csfw/util/slices"
	"golang.org/x/text/language"
	"golang.org/x/text/language/display"
)

// LocaleSeparator defines the underscore because in Magento land we also have
// the underscore as a separator.
// http://www.unicode.org/reports/tr35/#Language_and_Locale_IDs
const LocaleSeparator = "_"

// Available contains all available locales. One should not modify this slice.
var LocaleAvailable slices.String

// Supported contains all supported locales by this package. One should not modify this slice.
var LocaleSupported slices.String

var (
	// Only import the supported dictionaries here to reduce the amount of
	// data linked into your binary by only using the predefined Dictionary variables
	tags = []language.Tag{
		language.English, // first entry here is the default language
		language.German,
		language.French,
		language.Italian,
		language.Spanish,
		language.Japanese,
		language.Ukrainian,
	}
	dicts = []*display.Dictionary{
		display.English,
		display.German,
		display.French,
		display.Italian,
		display.Spanish,
		display.Japanese,
		display.Ukrainian,
	}
	matcher = language.NewMatcher(tags)
)

func init() {
	// @todo check if this is the proper way to generate the locales
	for _, tag := range display.Values.Tags() {
		b, bc := tag.Base()
		r, rc := tag.Region()
		if bc >= language.Exact && rc >= language.Low && !b.IsPrivateUse() && !r.IsPrivateUse() && r.IsCountry() {
			LocaleAvailable.Append(b.String() + LocaleSeparator + r.String())
		}
	}
	for _, tag := range tags {
		b, bc := tag.Base()
		r, rc := tag.Region()
		if bc >= language.High && rc >= language.Low {
			LocaleSupported.Append(b.String() + LocaleSeparator + r.String())
		}
	}
}

// getDict returns a predefnied dictionary for a given tag. If matching fails
// the fall-back language, which will be the first one passed to NewMatcher,
// will be returned.
func getDict(t language.Tag) *display.Dictionary {
	_, i, _ := matcher.Match(t)
	return dicts[i]
}

// GetLocale creates a new language Tag from a locale
func GetLocaleTag(locale string) (language.Tag, error) {
	return language.Parse(locale)
}

// GetLocale generates a locale from a base and a region and may use an optional script.
// lang
// lang_script
// lang_script_region
// lang_region (aliases to lang_script_region)
func GetLocale(b language.Base, r language.Region, s ...language.Script) string {
	ret := b.String()
	if len(s) == 1 {
		ret += LocaleSeparator + s[0].String()
	}
	return ret + LocaleSeparator + r.String()
}
