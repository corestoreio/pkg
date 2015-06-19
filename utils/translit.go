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

package utils

import "unicode"

// translitConvertTable courtesy: https://github.com/magento/magento2/blob/master/lib%2Finternal%2FMagento%2FFramework%2FFilter%2FTranslit.php
var translitConvertTable = map[rune][]rune{
	'€': []rune{'e', 'u', 'r', 'o'},
	'&': []rune{'a', 'n', 'd'},
	'@': []rune{'a', 't'},
	'©': []rune{'c'},
	'®': []rune{'r'},
	'À': []rune{'a'},
	'Á': []rune{'a'},
	'Â': []rune{'a'},
	'Ä': []rune{'a'},
	'Å': []rune{'a'},
	'Æ': []rune{'a', 'e'},
	'Ç': []rune{'c'},
	'È': []rune{'e'},
	'É': []rune{'e'},
	'Ë': []rune{'e'},
	'Ì': []rune{'i'},
	'Í': []rune{'i'},
	'Î': []rune{'i'},
	'Ï': []rune{'i'},
	'Ò': []rune{'o'},
	'Ó': []rune{'o'},
	'Ô': []rune{'o'},
	'Õ': []rune{'o'},
	'Ö': []rune{'o'},
	'Ø': []rune{'o'},
	'Ù': []rune{'u'},
	'Ú': []rune{'u'},
	'Û': []rune{'u'},
	'Ü': []rune{'u'},
	'Ý': []rune{'y'},
	'ß': []rune{'s', 's'},
	'à': []rune{'a'},
	'á': []rune{'a'},
	'â': []rune{'a'},
	'ä': []rune{'a'},
	'å': []rune{'a'},
	'æ': []rune{'a', 'e'},
	'ç': []rune{'c'},
	'è': []rune{'e'},
	'é': []rune{'e'},
	'ê': []rune{'e'},
	'ë': []rune{'e'},
	'ì': []rune{'i'},
	'í': []rune{'i'},
	'î': []rune{'i'},
	'ï': []rune{'i'},
	'ò': []rune{'o'},
	'ó': []rune{'o'},
	'ô': []rune{'o'},
	'õ': []rune{'o'},
	'ö': []rune{'o'},
	'ø': []rune{'o'},
	'ù': []rune{'u'},
	'ú': []rune{'u'},
	'û': []rune{'u'},
	'ü': []rune{'u'},
	'ý': []rune{'y'},
	'þ': []rune{'p'},
	'ÿ': []rune{'y'},
	'Ā': []rune{'a'},
	'ā': []rune{'a'},
	'Ă': []rune{'a'},
	'ă': []rune{'a'},
	'Ą': []rune{'a'},
	'ą': []rune{'a'},
	'Ć': []rune{'c'},
	'ć': []rune{'c'},
	'Ĉ': []rune{'c'},
	'ĉ': []rune{'c'},
	'Ċ': []rune{'c'},
	'ċ': []rune{'c'},
	'Č': []rune{'c'},
	'č': []rune{'c'},
	'Ď': []rune{'d'},
	'ď': []rune{'d'},
	'Đ': []rune{'d'},
	'đ': []rune{'d'},
	'Ē': []rune{'e'},
	'ē': []rune{'e'},
	'Ĕ': []rune{'e'},
	'ĕ': []rune{'e'},
	'Ė': []rune{'e'},
	'ė': []rune{'e'},
	'Ę': []rune{'e'},
	'ę': []rune{'e'},
	'Ě': []rune{'e'},
	'ě': []rune{'e'},
	'Ĝ': []rune{'g'},
	'ĝ': []rune{'g'},
	'Ğ': []rune{'g'},
	'ğ': []rune{'g'},
	'Ġ': []rune{'g'},
	'ġ': []rune{'g'},
	'Ģ': []rune{'g'},
	'ģ': []rune{'g'},
	'Ĥ': []rune{'h'},
	'ĥ': []rune{'h'},
	'Ħ': []rune{'h'},
	'ħ': []rune{'h'},
	'Ĩ': []rune{'i'},
	'ĩ': []rune{'i'},
	'Ī': []rune{'i'},
	'ī': []rune{'i'},
	'Ĭ': []rune{'i'},
	'ĭ': []rune{'i'},
	'Į': []rune{'i'},
	'į': []rune{'i'},
	'İ': []rune{'i'},
	'ı': []rune{'i'},
	'Ĳ': []rune{'i', 'j'},
	'ĳ': []rune{'i', 'j'},
	'Ĵ': []rune{'j'},
	'ĵ': []rune{'j'},
	'Ķ': []rune{'k'},
	'ķ': []rune{'k'},
	'ĸ': []rune{'k'},
	'Ĺ': []rune{'l'},
	'ĺ': []rune{'l'},
	'Ļ': []rune{'l'},
	'ļ': []rune{'l'},
	'Ľ': []rune{'l'},
	'ľ': []rune{'l'},
	'Ŀ': []rune{'l'},
	'ŀ': []rune{'l'},
	'Ł': []rune{'l'},
	'ł': []rune{'l'},
	'Ń': []rune{'n'},
	'ń': []rune{'n'},
	'Ņ': []rune{'n'},
	'ņ': []rune{'n'},
	'Ň': []rune{'n'},
	'ň': []rune{'n'},
	'ŉ': []rune{'n'},
	'Ŋ': []rune{'n'},
	'ŋ': []rune{'n'},
	'Ō': []rune{'o'},
	'ō': []rune{'o'},
	'Ŏ': []rune{'o'},
	'ŏ': []rune{'o'},
	'Ő': []rune{'o'},
	'ő': []rune{'o'},
	'Œ': []rune{'o', 'e'},
	'œ': []rune{'o', 'e'},
	'Ŕ': []rune{'r'},
	'ŕ': []rune{'r'},
	'Ŗ': []rune{'r'},
	'ŗ': []rune{'r'},
	'Ř': []rune{'r'},
	'ř': []rune{'r'},
	'Ś': []rune{'s'},
	'ś': []rune{'s'},
	'Ŝ': []rune{'s'},
	'ŝ': []rune{'s'},
	'Ş': []rune{'s'},
	'ş': []rune{'s'},
	'Š': []rune{'s'},
	'š': []rune{'s'},
	'Ţ': []rune{'t'},
	'ţ': []rune{'t'},
	'Ť': []rune{'t'},
	'ť': []rune{'t'},
	'Ŧ': []rune{'t'},
	'ŧ': []rune{'t'},
	'Ũ': []rune{'u'},
	'ũ': []rune{'u'},
	'Ū': []rune{'u'},
	'ū': []rune{'u'},
	'Ŭ': []rune{'u'},
	'ŭ': []rune{'u'},
	'Ů': []rune{'u'},
	'ů': []rune{'u'},
	'Ű': []rune{'u'},
	'ű': []rune{'u'},
	'Ų': []rune{'u'},
	'ų': []rune{'u'},
	'Ŵ': []rune{'w'},
	'ŵ': []rune{'w'},
	'Ŷ': []rune{'y'},
	'ŷ': []rune{'y'},
	'Ÿ': []rune{'y'},
	'Ź': []rune{'z'},
	'ź': []rune{'z'},
	'Ż': []rune{'z'},
	'ż': []rune{'z'},
	'Ž': []rune{'z'},
	'ž': []rune{'z'},
	'ſ': []rune{'z'},
	'Ə': []rune{'e'},
	'ƒ': []rune{'f'},
	'Ơ': []rune{'o'},
	'ơ': []rune{'o'},
	'Ư': []rune{'u'},
	'ư': []rune{'u'},
	'Ǎ': []rune{'a'},
	'ǎ': []rune{'a'},
	'Ǐ': []rune{'i'},
	'ǐ': []rune{'i'},
	'Ǒ': []rune{'o'},
	'ǒ': []rune{'o'},
	'Ǔ': []rune{'u'},
	'ǔ': []rune{'u'},
	'Ǖ': []rune{'u'},
	'ǖ': []rune{'u'},
	'Ǘ': []rune{'u'},
	'ǘ': []rune{'u'},
	'Ǚ': []rune{'u'},
	'ǚ': []rune{'u'},
	'Ǜ': []rune{'u'},
	'ǜ': []rune{'u'},
	'Ǻ': []rune{'a'},
	'ǻ': []rune{'a'},
	'Ǽ': []rune{'a', 'e'},
	'ǽ': []rune{'a', 'e'},
	'Ǿ': []rune{'o'},
	'ǿ': []rune{'o'},
	'ə': []rune{'e'},
	'Ё': []rune{'j', 'o'},
	'Є': []rune{'e'},
	'І': []rune{'i'},
	'Ї': []rune{'i'},
	'А': []rune{'a'},
	'Б': []rune{'b'},
	'В': []rune{'v'},
	'Г': []rune{'g'},
	'Д': []rune{'d'},
	'Е': []rune{'e'},
	'Ж': []rune{'z', 'h'},
	'З': []rune{'z'},
	'И': []rune{'i'},
	'Й': []rune{'j'},
	'К': []rune{'k'},
	'Л': []rune{'l'},
	'М': []rune{'m'},
	'Н': []rune{'n'},
	'О': []rune{'o'},
	'П': []rune{'p'},
	'Р': []rune{'r'},
	'С': []rune{'s'},
	'Т': []rune{'t'},
	'У': []rune{'u'},
	'Ф': []rune{'f'},
	'Х': []rune{'h'},
	'Ц': []rune{'c'},
	'Ч': []rune{'c', 'h'},
	'Ш': []rune{'s', 'h'},
	'Щ': []rune{'s', 'c', 'h'},
	'Ъ': []rune{'-'},
	'Ы': []rune{'y'},
	'Ь': []rune{'-'},
	'Э': []rune{'j', 'e'},
	'Ю': []rune{'j', 'u'},
	'Я': []rune{'j', 'a'},
	'а': []rune{'a'},
	'б': []rune{'b'},
	'в': []rune{'v'},
	'г': []rune{'g'},
	'д': []rune{'d'},
	'е': []rune{'e'},
	'ж': []rune{'z', 'h'},
	'з': []rune{'z'},
	'и': []rune{'i'},
	'й': []rune{'j'},
	'к': []rune{'k'},
	'л': []rune{'l'},
	'м': []rune{'m'},
	'н': []rune{'n'},
	'о': []rune{'o'},
	'п': []rune{'p'},
	'р': []rune{'r'},
	'с': []rune{'s'},
	'т': []rune{'t'},
	'у': []rune{'u'},
	'ф': []rune{'f'},
	'х': []rune{'h'},
	'ц': []rune{'c'},
	'ч': []rune{'c', 'h'},
	'ш': []rune{'s', 'h'},
	'щ': []rune{'s', 'c', 'h'},
	'ъ': []rune{'-'},
	'ы': []rune{'y'},
	'ь': []rune{'-'},
	'э': []rune{'j', 'e'},
	'ю': []rune{'j', 'u'},
	'я': []rune{'j', 'a'},
	'ё': []rune{'j', 'o'},
	'є': []rune{'e'},
	'і': []rune{'i'},
	'ї': []rune{'i'},
	'Ґ': []rune{'g'},
	'ґ': []rune{'g'},
	'א': []rune{'a'},
	'ב': []rune{'b'},
	'ג': []rune{'g'},
	'ד': []rune{'d'},
	'ה': []rune{'h'},
	'ו': []rune{'v'},
	'ז': []rune{'z'},
	'ח': []rune{'h'},
	'ט': []rune{'t'},
	'י': []rune{'i'},
	'ך': []rune{'k'},
	'כ': []rune{'k'},
	'ל': []rune{'l'},
	'ם': []rune{'m'},
	'מ': []rune{'m'},
	'ן': []rune{'n'},
	'נ': []rune{'n'},
	'ס': []rune{'s'},
	'ע': []rune{'e'},
	'ף': []rune{'p'},
	'פ': []rune{'p'},
	'ץ': []rune{'C'},
	'צ': []rune{'c'},
	'ק': []rune{'q'},
	'ר': []rune{'r'},
	'ש': []rune{'w'},
	'ת': []rune{'t'},
	'™': []rune{'t', 'm'},
}

/*
	@todo implmement https://github.com/magento/magento2/blob/master/lib%2Finternal%2FMagento%2FFramework%2FFilter%2FTranslit.php#L345
	needs maybe a re-architecture
*/

// Translit replaces characters in a string using a conversion table,
// e.g. ™ => tm; © => c; @ => at; € => euro
func Translit(str []rune) []rune {
	// @todo thread safe
	i := 0
	for i < len(str) {
		if to, ok := translitConvertTable[str[i]]; ok {
			if len(to) < 1 {
				i++
				continue
			}
			if len(to) == 1 {
				str[i] = to[0]
				i++
				continue
			}
			str = append(str[:i], append(to, str[i+1:]...)...)
			i = 0 // reset i an re-run the str because length changed
		} else {
			i++
		}
	}
	return str
}

// TranslitURL same as Translit() but removes after conversion any non 0-9A-Za-z
// characters and replaces them with a dash.
// This function is responsible for the slug generation for product/category/cms URLs.
func TranslitURL(str []rune) []rune {
	str = Translit(str)
	for i, r := range str {
		str[i] = unicode.ToLower(r)
		switch {
		case '0' <= r && r <= '9':
			continue
		case 'a' <= r && r <= 'z':
			continue
		case 'A' <= r && r <= 'Z':
			continue
		case r == '-':
			continue
		default:
			str[i] = '-'
		}
	}
	i := 0
	// remove multiple dashes
	for i < len(str) {
		j := i + 1
		if str[i] == '-' && j < len(str) && str[j] == '-' { // look ahead
			str = append(str[:i], str[j:]...)
			i = 0 // reset i due to change in length of str
			continue
		}
		i++
	}
	// trim
	if len(str) > 0 {
		if str[0] == '-' {
			str = str[1:]
		}
		if str[len(str)-1] == '-' {
			str = str[:len(str)-1]
		}
	}
	return str
}
