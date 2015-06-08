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

package i18n

import "sort"

var currencyDict CurrencyDictSlice

type (
	CurrencyDictSlice []currencyDict
	currencyDict      struct {
		l  string   // locale
		t  tagIndex // indexes currency codes
		cn header   // currency names
		cs header   // currency symbol
	}
)

func SetCurrencyDict(cds ...currencyDict) {
	currencyDict = CurrencyDictSlice(cds)
}

func NewCurrencyDict(locale string, ti tagIndex, currencyNames, currencySymbols header) currencyDict {
	return currencyDict{
		l:  locale,
		t:  ti,
		cn: currencyNames,
		cs: currencySymbols,
	}
}

func NewHeader(d string, i []uint16) header {
	return header{d, i}
}

func NewTagIndex(ti [3]string) tagIndex {
	return tagIndex(ti)
}

// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// header contains the data and indexes for a single namer.
// data contains a series of strings concatenated into one. index contains the
// offsets for a string in data. For example, consider a header that defines
// strings for the languages de, el, en, fi, and nl:
//
// 		header{
// 			data: "GermanGreekEnglishDutch",
//  		index: []uint16{ 0, 6, 11, 18, 18, 23 },
// 		}
//
// For a language with index i, the string is defined by
// data[index[i]:index[i+1]]. So the number of elements in index is always one
// greater than the number of languages for which header defines a value.
// A string for a language may be empty, which means the name is undefined. In
// the above example, the name for fi (Finnish) is undefined.
type header struct {
	data  string
	index []uint16
}

// name looks up the name for a tag in the dictionary, given its index.
func (h *header) name(i int) string {
	if i < len(h.index)-1 {
		return h.data[h.index[i]:h.index[i+1]]
	}
	return ""
}

// tagIndex holds a concatenated lists of subtags of length 2 to 4, one string
// for each length, which can be used in combination with binary search to get
// the index associated with a tag.
// For example, a tagIndex{
//   "arenesfrruzh",  // 6 2-byte tags.
//   "barwae",        // 2 3-byte tags.
//   "",
// }
// would mean that the 2-byte tag "fr" had an index of 3, and the 3-byte tag
// "wae" had an index of 7.
type tagIndex [3]string

func (t *tagIndex) index(s string) int {
	sz := len(s)
	if sz < 2 || 4 < sz {
		return -1
	}
	a := t[sz-2]
	index := sort.Search(len(a)/sz, func(i int) bool {
		p := i * sz
		return a[p:p+sz] >= s
	})
	p := index * sz
	if end := p + sz; end > len(a) || a[p:end] != s {
		return -1
	}
	// Add the number of tags for smaller sizes.
	for i := 0; i < sz-2; i++ {
		index += len(t[i]) / (i + 2)
	}
	return index
}

// len returns the number of tags that are contained in the tagIndex.
func (t *tagIndex) len() (n int) {
	for i, s := range t {
		n += len(s) / (i + 2)
	}
	return n
}

// keys calls f for each tag.
func (t *tagIndex) keys(f func(key string)) {
	for i, s := range *t {
		for ; s != ""; s = s[i+2:] {
			f(s[:i+2])
		}
	}
}
