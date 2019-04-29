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

package locale

import (
	"strconv"
	"unicode"
)

// @todo see lib/internal/Magento/Framework/Locale/Format.php

// ExtractNumber returns the first found number from a string.
//
// Examples for input:
//     '  2345.4356,1234' = 23455456.1234
//     '+23,3452.123' = 233452.123
//     ' 12343 ' = 12343
//     '-9456km' = -9456
//     '0' = 0
//     '2 054,10' = 2054.1
//     '2'054.52' = 2054.52
//     '2,46 GB' = 2.46
func ExtractNumber(n string) (float64, error) {

	if n == "" {
		return 0, nil
	}

	var nums = make([]rune, 0, len(n))
	var digitStarted bool

OuterLoop:
	for _, r := range n {
		switch {
		case r == '.', r == ',', unicode.IsDigit(r):
			digitStarted = true
			nums = append(nums, r)
		case r == '-':
			if false == digitStarted {
				nums = append(nums, r)
			} else {
				break OuterLoop // as soon as any other character occurs we break the loop
			}
		case r == '\'', r == '+', unicode.IsSpace(r):
			continue
		default:
			if digitStarted {
				break OuterLoop // as soon as any other character occurs, after we found a number, we break the loop.
			}
		}
	}

	if len(nums) == 0 || false == digitStarted {
		return 0, nil
	}

	var hasComma, hasDot int
	for i, r := range nums { // get positions
		switch r {
		case '.':
			hasDot = i
		case ',':
			hasComma = i
		}
	}

	switch {
	case hasComma > 0 && hasDot == 0: // 1234,56 => 1234.56
		for i, r := range nums {
			if r == ',' {
				nums[i] = '.'
				break
			}
		}
	case hasComma > 0 && hasDot > 0:
		if hasComma > hasDot { // 1.234,56
			nums[hasComma] = '.'                             // replace comma with dot
			nums = append(nums[:hasDot], nums[hasDot+1:]...) // remove dot
		} else { // 1,234.56
			nums = append(nums[:hasComma], nums[hasComma+1:]...) // remove comma
		}
	}

	return strconv.ParseFloat(string(nums), 64)
}

// FormatPrice returns price formatting information for a given locale.
// @todo implement and use github.com/leekchan/accounting
func FormatPrice( /*localCode, currencyCode */ ) /*some kind of struct*/ {
	// if currency code has been provided: load currency from directory package
	// otherwise determine it from the current loaded store

	// if localeCode is empty, the locale code will be determined from the current store

	// move i18n.Symbols including formatting string into this package
}
