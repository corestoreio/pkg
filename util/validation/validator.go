// Copyright 2015-present, Cyrill @ Schumacher.fm and the CoreStore contributors
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

// The MIT License (MIT)
//
// Copyright (c) 2014 Alex Saskevich
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package validation

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"net"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/corestoreio/pkg/util/conv"
)

var (
	notNumberRegexp     = regexp.MustCompile("[^0-9]+")
	whiteSpacesAndMinus = regexp.MustCompile("[\\s-]+")
)

const maxURLRuneCount = 2083
const minURLRuneCount = 3
const RF3339WithoutZone = "2006-01-02T15:04:05"

// IsEmailRegexp runs a complex regex check if the string is an email.
func IsEmailRegexp(str string) bool {
	// TODO uppercase letters are not supported
	return rxEmail.MatchString(str)
}

// IsEmailSimple checks if the str contains one @.
func IsEmailSimple(str string) bool {
	idx := strings.IndexByte(str, '@')
	if idx < 0 {
		return false
	}
	user := str[:idx]
	host := str[idx+1:]
	return strings.Count(str, "@") == 1 && user != "" && host != ""
}

// IsExistingEmail check if the string is an email of existing domain. Looks up
// the MX record.
func IsExistingEmail(email string) bool {

	if len(email) < 6 || len(email) > 254 {
		return false
	}
	at := strings.LastIndex(email, "@")
	if at <= 0 || at > len(email)-3 {
		return false
	}
	user := email[:at]
	host := email[at+1:]
	if len(user) > 64 {
		return false
	}
	if userDotRegexp.MatchString(user) || !userRegexp.MatchString(user) || !hostRegexp.MatchString(host) {
		return false
	}
	switch host {
	case "localhost", "example.com":
		return true
	}
	if _, err := net.LookupMX(host); err != nil {
		if _, err := net.LookupIP(host); err != nil {
			return false
		}
	}

	return true
}

// IsURL check if the string is an URL.
func IsURL(str string) bool {
	if str == "" || utf8.RuneCountInString(str) >= maxURLRuneCount || len(str) <= minURLRuneCount || strings.HasPrefix(str, ".") {
		return false
	}
	strTemp := str
	if strings.Index(str, ":") >= 0 && strings.Index(str, "://") == -1 {
		// support no indicated urlscheme but with colon for port number
		// http:// is appended so url.Parse will succeed, strTemp used so it does not impact rxURL.MatchString
		strTemp = "http://" + str
	}
	u, err := url.Parse(strTemp)
	if err != nil {
		return false
	}
	if strings.HasPrefix(u.Host, ".") {
		return false
	}
	if u.Host == "" && (u.Path != "" && !strings.Contains(u.Path, ".")) {
		return false
	}
	return rxURL.MatchString(str)
}

// IsRequestURL check if the string rawurl, assuming
// it was received in an HTTP request, is a valid
// URL confirm to RFC 3986
func IsRequestURL(rawurl string) bool {
	url, err := url.ParseRequestURI(rawurl)
	if err != nil {
		return false //Couldn't even parse the rawurl
	}
	if len(url.Scheme) == 0 {
		return false //No Scheme found
	}
	return true
}

// IsRequestURI check if the string rawurl, assuming
// it was received in an HTTP request, is an
// absolute URI or an absolute path.
func IsRequestURI(rawurl string) bool {
	_, err := url.ParseRequestURI(rawurl)
	return err == nil
}

// IsAlpha check if the string contains only letters (a-zA-Z). Empty string is valid.
func IsAlpha(str string) bool {
	// "^[a-zA-Z]+$"
	if IsNull(str) {
		return true
	}
	l := utf8.RuneCountInString(str)
	if l != len(str) {
		return false
	}

	for _, r := range str {
		switch {
		case 'a' <= r && r <= 'z', 'A' <= r && r <= 'Z':
			// ok
		default:
			return false
		}
	}
	return true
}

//IsUTFLetter check if the string contains only unicode letter characters.
//Similar to IsAlpha but for all languages. Empty string is valid.
func IsUTFLetter(str string) bool {
	if IsNull(str) {
		return true
	}

	for _, c := range str {
		if !unicode.IsLetter(c) {
			return false
		}
	}
	return true

}

// IsAlphanumeric check if the string contains only letters and numbers. Empty string is valid.
func IsAlphanumeric(str string) bool {
	// "^[a-zA-Z0-9]+$"
	if IsNull(str) {
		return true
	}
	l := utf8.RuneCountInString(str)
	if l != len(str) {
		return false
	}

	for _, r := range str {
		switch {
		case 'a' <= r && r <= 'z', 'A' <= r && r <= 'Z', '0' <= r && r <= '9':
			// ok
		default:
			return false
		}
	}
	return true
}

// IsUTFLetterNumeric check if the string contains only unicode letters and numbers. Empty string is valid.
func IsUTFLetterNumeric(str string) bool {
	if IsNull(str) {
		return true
	}
	for _, c := range str {
		if !unicode.IsLetter(c) && !unicode.IsNumber(c) { //letters && numbers are ok
			return false
		}
	}
	return true

}

// IsNumeric check if the string contains only numbers. Empty string is valid.
func IsNumeric(str string) bool {
	if IsNull(str) {
		return true
	}
	l := utf8.RuneCountInString(str)
	if l != len(str) {
		return false
	}

	for _, r := range str {
		switch {
		case '0' <= r && r <= '9':
			// ok
		default:
			return false
		}
	}
	return true
}

// IsUTFNumeric check if the string contains only unicode numbers of any kind.
// Numbers can be 0-9 but also Fractions ¾,Roman Ⅸ and Hangzhou 〩. Empty string is valid.
func IsUTFNumeric(str string) bool {
	if IsNull(str) {
		return true
	}
	if strings.IndexAny(str, "+-") > 0 {
		return false
	}
	if len(str) > 1 {
		str = strings.TrimPrefix(str, "-")
		str = strings.TrimPrefix(str, "+")
	}
	for _, c := range str {
		if !unicode.IsNumber(c) { //numbers && minus sign are ok
			return false
		}
	}
	return true

}

// IsUTFDigit check if the string contains only unicode radix-10 decimal digits. Empty string is valid.
func IsUTFDigit(str string) bool {
	if IsNull(str) {
		return true
	}
	if strings.IndexAny(str, "+-") > 0 {
		return false
	}
	if len(str) > 1 {
		str = strings.TrimPrefix(str, "-")
		str = strings.TrimPrefix(str, "+")
	}
	for _, c := range str {
		if !unicode.IsDigit(c) { //digits && minus sign are ok
			return false
		}
	}
	return true

}

// IsHexadecimal check if the string is a hexadecimal number.
func IsHexadecimal(str string) bool {
	// "^[0-9a-fA-F]+$"
	l := utf8.RuneCountInString(str)
	if l == 0 || l != len(str) {
		return false
	}

	for _, r := range str {
		switch {
		case 'a' <= r && r <= 'f', 'A' <= r && r <= 'F', '0' <= r && r <= '9':
			// ok
		default:
			return false
		}
	}
	return true
}

// IsHexcolor check if the string is a hexadecimal color.
func IsHexcolor(str string) bool {
	// "^#?([0-9a-fA-F]{3}|[0-9a-fA-F]{6})$"
	switch utf8.RuneCountInString(str) {
	case 0:
		return false
	case 4, 7: // #fff,#ffffff
		if str[:1] != "#" {
			return false
		}
		str = str[1:]
		fallthrough
	case 3, 6: // fff,FFFFFF
		for _, r := range str {
			switch {
			case 'a' <= r && r <= 'f', 'A' <= r && r <= 'F', '0' <= r && r <= '9':
				// ok
			default:
				return false
			}
		}
	}
	return true
}

// IsRGBcolor check if the string is a valid RGB color in form rgb(RRR, GGG, BBB).
func IsRGBcolor(str string) bool {
	return rxRGBcolor.MatchString(str)
}

// IsLowerCase check if the string is lowercase. Empty string is valid.
func IsLowerCase(str string) bool {
	if IsNull(str) {
		return true
	}
	return str == strings.ToLower(str)
}

// IsUpperCase check if the string is uppercase. Empty string is valid.
func IsUpperCase(str string) bool {
	if IsNull(str) {
		return true
	}
	return str == strings.ToUpper(str)
}

// HasLowerCase check if the string contains at least 1 lowercase. Empty string is valid.
func HasLowerCase(str string) bool {
	if IsNull(str) {
		return true
	}
	return rxHasLowerCase.MatchString(str)
}

// HasUpperCase check if the string contians as least 1 uppercase. Empty string is valid.
func HasUpperCase(str string) bool {
	if IsNull(str) {
		return true
	}
	return rxHasUpperCase.MatchString(str)
}

// IsInt check if the string is an integer. Empty string is valid.
func IsInt(str string) bool {
	// using no REGEX here speeds up this function by factor 25.
	if IsNull(str) {
		return true
	}
	if strings.HasPrefix(str, "-") || strings.HasPrefix(str, "+") {
		str = str[2:]
	}

	if str == "0" {
		return true
	}

	for i, r := range str {
		switch {
		case i == 0 && '1' <= r && r <= '9': // must start not with zero
			// ok
		case i > 0 && '0' <= r && r <= '9':
			// ok
		default:
			return false
		}
	}
	return true
}

// IsBool check if the string is a bool. Empty string is valid.
func IsBool(str string) bool {
	if IsNull(str) {
		return true
	}
	_, err := strconv.ParseBool(str)
	return err == nil
}

// IsFloat check if the string is a float.
func IsFloat(str string) bool {
	// "^(?:[-+]?(?:[0-9]+))?(?:\\.[0-9]*)?(?:[eE][\\+\\-]?(?:[0-9]+))?$"
	// return str != "" && rxFloat.MatchString(str)
	_, err := strconv.ParseFloat(str, 64)
	return err == nil
}

// IsDivisibleBy check if the string is a number that's divisible by another. If
// second argument is not valid integer or zero, it's return false. Otherwise,
// if first argument is not valid integer or zero, it's return true (Invalid
// string converts to zero).
func IsDivisibleBy(str, num string) bool {
	f := conv.ToFloat64(str)
	p := int64(f)
	q := conv.ToInt64(num)
	if q == 0 {
		return false
	}
	return (p == 0) || (p%q == 0)
}

// IsNull check if the string is null.
func IsNull(str string) bool {
	return len(str) == 0
}

// IsByteLength check if the string's length (in bytes) falls in a range.
func IsByteLength(str string, min, max int) bool {
	return len(str) >= min && len(str) <= max
}

func isUUIDvX(str, version string) bool {
	// "^[0-9a-f]{8}-[0-9a-f]{4}-3[0-9a-f]{3}-[0-9a-f]{4}-[0-9a-f]{12}$"
	// a987fbc9-4bed-3078-cf07-9141ba07c9f3
	l := utf8.RuneCountInString(str)
	if l == 0 || l != len(str) || l != 36 {
		return false
	}
	if version != "" && str[14:15] != version {
		return false
	}
	if str[8:9] != "-" || str[13:14] != "-" || str[18:19] != "-" || str[23:24] != "-" {
		return false
	}
	for _, r := range str {
		switch {
		case r == '-', 'a' <= r && r <= 'f', 'A' <= r && r <= 'F', '0' <= r && r <= '9':
			// ok
		default:
			return false
		}
	}
	return true
}

// IsUUIDv3 check if the string is a UUID version 3. Does not use regex. Highly
// optimized.
func IsUUIDv3(str string) bool {
	// "^[0-9a-f]{8}-[0-9a-f]{4}-3[0-9a-f]{3}-[0-9a-f]{4}-[0-9a-f]{12}$"
	// a987fbc9-4bed-3078-cf07-9141ba07c9f3
	return isUUIDvX(str, "3")
}

// IsUUIDv4 check if the string is a UUID version 4. Does not use regex. Highly
// optimized.
func IsUUIDv4(str string) bool {
	// "^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$"
	// 625e63f3-58f5-40b7-83a1-a72ad31acffb
	if !isUUIDvX(str, "4") {
		return false
	}
	switch str[19:20] {
	case "8", "9", "a", "b":
		return true
	}
	return false
}

// IsUUIDv5 check if the string is a UUID version 5. Does not use regex. Highly
// optimized.
func IsUUIDv5(str string) bool {
	// "^[0-9a-f]{8}-[0-9a-f]{4}-5[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$"
	// 987fbc97-4bed-5078-9f07-9141ba07c9f3
	if !isUUIDvX(str, "5") {
		return false
	}
	switch str[19:20] {
	case "8", "9", "a", "b":
		return true
	}
	return false
}

// IsUUID check if the string is a UUID (version 3, 4 or 5). Does not use regex.
// Highly optimized.
func IsUUID(str string) bool {
	// "^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"
	// a987fbc9-4bed-3078-cf07-9141ba07c9f3
	return isUUIDvX(str, "")
}

// IsCreditCard check if the string is a credit card.
func IsCreditCard(str string) bool {
	sanitized := notNumberRegexp.ReplaceAllString(str, "")
	if !rxCreditCard.MatchString(sanitized) {
		return false
	}
	var sum int64
	var digit string
	var tmpNum int64
	var shouldDouble bool
	for i := len(sanitized) - 1; i >= 0; i-- {
		digit = sanitized[i:(i + 1)]
		tmpNum = conv.ToInt64(digit)
		if shouldDouble {
			tmpNum *= 2
			if tmpNum >= 10 {
				sum += (tmpNum % 10) + 1
			} else {
				sum += tmpNum
			}
		} else {
			sum += tmpNum
		}
		shouldDouble = !shouldDouble
	}

	return sum%10 == 0
}

// IsISBN10 check if the string is an ISBN version 10.
func IsISBN10(str string) bool {
	return IsISBN(str, 10)
}

// IsISBN13 check if the string is an ISBN version 13.
func IsISBN13(str string) bool {
	return IsISBN(str, 13)
}

// IsISBN check if the string is an ISBN (version 10 or 13).
// If version value is not equal to 10 or 13, it will be check both variants.
func IsISBN(str string, version int) bool {
	sanitized := whiteSpacesAndMinus.ReplaceAllString(str, "")
	var checksum int32
	var i int32
	if version == 10 {
		if !rxISBN10.MatchString(sanitized) {
			return false
		}
		for i = 0; i < 9; i++ {
			checksum += (i + 1) * int32(sanitized[i]-'0')
		}
		if sanitized[9] == 'X' {
			checksum += 10 * 10
		} else {
			checksum += 10 * int32(sanitized[9]-'0')
		}
		return checksum%11 == 0
	} else if version == 13 {
		if !rxISBN13.MatchString(sanitized) {
			return false
		}
		factor := []int32{1, 3}
		for i = 0; i < 12; i++ {
			checksum += factor[i%2] * int32(sanitized[i]-'0')
		}
		return (int32(sanitized[12]-'0'))-((10-(checksum%10))%10) == 0
	}
	return IsISBN(str, 10) || IsISBN(str, 13)
}

// IsMultibyte check if the string contains one or more multibyte chars. Empty string is valid.
func IsMultibyte(str string) bool {
	if IsNull(str) {
		return true
	}
	return rxMultibyte.MatchString(str)
}

// IsASCII check if the string contains ASCII chars only. Empty string is valid.
func IsASCII(str string) bool {
	if IsNull(str) {
		return true
	}
	return rxASCII.MatchString(str)
}

// IsPrintableASCII check if the string contains printable ASCII chars only. Empty string is valid.
func IsPrintableASCII(str string) bool {
	if IsNull(str) {
		return true
	}
	return rxPrintableASCII.MatchString(str)
}

// IsFullWidth check if the string contains any full-width chars. Empty string is valid.
func IsFullWidth(str string) bool {
	if IsNull(str) {
		return true
	}
	return rxFullWidth.MatchString(str)
}

// IsHalfWidth check if the string contains any half-width chars. Empty string is valid.
func IsHalfWidth(str string) bool {
	if IsNull(str) {
		return true
	}
	return rxHalfWidth.MatchString(str)
}

// IsVariableWidth check if the string contains a mixture of full and half-width chars. Empty string is valid.
func IsVariableWidth(str string) bool {
	if IsNull(str) {
		return true
	}
	return rxHalfWidth.MatchString(str) && rxFullWidth.MatchString(str)
}

// IsBase64 check if a string is base64 encoded. Optimized version without
// regex.
func IsBase64(str string) bool {
	// "^(?:[A-Za-z0-9+\\/]{4})*(?:[A-Za-z0-9+\\/]{2}==|[A-Za-z0-9+\\/]{3}=|[A-Za-z0-9+\\/]{4})$"
	// Jon Skeet on SO: Well you can: Check that the length is a multiple of 4
	// characters Check that every character is in the set A-Z, a-z, 0-9, +, /
	// except for padding at the end which is 0, 1 or 2 '=' characters If you're
	// expecting that it will be base64, then you can probably just use whatever
	// library is available on your platform to try to decode it to a byte
	// array, throwing an exception if it's not valid base 64. That depends on
	// your platform, of course.
	l := utf8.RuneCountInString(str)
	if l == 0 || l != len(str) || l%4 != 0 {
		return false
	}

	var equals int
	for _, r := range str {
		switch {
		case r == '=':
			equals++
		case r == '+', r == '/':
		case 'a' <= r && r <= 'z', 'A' <= r && r <= 'Z', '0' <= r && r <= '9':
			// ok
		default:
			return false
		}
	}
	return equals <= 3

}

// IsFilePath check is a string is Win or Unix file path and returns it's type.
func IsFilePath(str string) (bool, int) {
	if rxWinPath.MatchString(str) {
		//check windows path limit see:
		//  http://msdn.microsoft.com/en-us/library/aa365247(VS.85).aspx#maxpath
		if len(str[3:]) > 32767 {
			return false, Win
		}
		return true, Win
	} else if rxUnixPath.MatchString(str) {
		return true, Unix
	}
	return false, Unknown
}

// IsDataURI checks if a string is base64 encoded data URI such as an image
func IsDataURI(str string) bool {
	dataURI := strings.Split(str, ",")
	if !rxDataURI.MatchString(dataURI[0]) {
		return false
	}
	return IsBase64(dataURI[1])
}

// IsISO3166Alpha2 checks if a string is valid two-letter country code
func IsISO3166Alpha2(str string) bool {
	for _, entry := range ISO3166List {
		if str == entry.Alpha2Code {
			return true
		}
	}
	return false
}

// IsISO3166Alpha3 checks if a string is valid three-letter country code
func IsISO3166Alpha3(str string) bool {
	for _, entry := range ISO3166List {
		if str == entry.Alpha3Code {
			return true
		}
	}
	return false
}

// IsISO693Alpha2 checks if a string is valid two-letter language code
func IsISO693Alpha2(str string) bool {
	for _, entry := range ISO693List {
		if str == entry.Alpha2Code {
			return true
		}
	}
	return false
}

// IsISO693Alpha3b checks if a string is valid three-letter language code.
func IsISO693Alpha3b(str string) bool {
	for _, entry := range ISO693List {
		if str == entry.Alpha3bCode {
			return true
		}
	}
	return false
}

// IsLocale checks case-insensitive if the provided string is a locale.
func IsLocale(lcid string) bool {
	lcid = strings.ToLower(lcid)
	var buf strings.Builder
	for _, r := range lcid {
		if 'a' <= r && r <= 'z' || '0' <= r && r <= '9' {
			buf.WriteRune(r)
		}
	}
	_, ok := LocaleList[buf.String()]
	return ok
}

// IsDNSName will validate the given string as a DNS name
func IsDNSName(str string) bool {
	if str == "" || len(strings.Replace(str, ".", "", -1)) > 255 {
		// constraints already violated
		return false
	}
	return !IsIP(str) && rxDNSName.MatchString(str)
}

// IsHash checks if a string is a hash of type algorithm. Algorithm is one of
// ['md4', 'md5', 'sha1', 'sha256', 'sha384', 'sha512', 'ripemd128',
// 'ripemd160', 'tiger128', 'tiger160', 'tiger192', 'crc32', 'crc32b']
func IsHash(str string, algorithm string) bool {
	var aLen int
	algo := strings.ToLower(algorithm)

	switch algo {
	case "crc32", "crc32b":
		aLen = 8
	case "md5", "md4", "ripemd128", "tiger128":
		aLen = 32
	case "sha1", "ripemd160", "tiger160":
		aLen = 40
	case "tiger192":
		aLen = 48
	case "sha256":
		aLen = 64
	case "sha384":
		aLen = 96
	case "sha512":
		aLen = 128
	default:
		return false
	}
	sLen := len(str)
	if sLen != aLen {
		return false
	}
	var rLen int
	for _, r := range str {
		switch {
		case r >= 'a' && r <= 'f':
			rLen++
		case r >= 'A' && r <= 'F':
			rLen++
		case r >= '0' && r <= '9':
			rLen++
		}
	}

	return sLen == rLen
}

// IsDialString validates the given string for usage with the various Dial() functions
func IsDialString(str string) bool {

	if h, p, err := net.SplitHostPort(str); err == nil && h != "" && p != "" && (IsDNSName(h) || IsIP(h)) && IsPort(p) {
		return true
	}

	return false
}

// IsIP checks if a string is either IP version 4 or 6.
func IsIP(str string) bool {
	return net.ParseIP(str) != nil
}

// IsPort checks if a string represents a valid port
func IsPort(str string) bool {
	if i, err := strconv.Atoi(str); err == nil && i > 0 && i < 65536 {
		return true
	}
	return false
}

// IsIPv4 check if the string is an IP version 4.
func IsIPv4(str string) bool {
	ip := net.ParseIP(str)
	return ip != nil && strings.Contains(str, ".")
}

// IsIPv6 check if the string is an IP version 6.
func IsIPv6(str string) bool {
	ip := net.ParseIP(str)
	return ip != nil && strings.Contains(str, ":")
}

// IsCIDR check if the string is an valid CIDR notiation (IPV4 & IPV6)
func IsCIDR(str string) bool {
	_, _, err := net.ParseCIDR(str)
	return err == nil
}

// IsMAC check if a string is valid MAC address.
// Possible MAC formats:
// 01:23:45:67:89:ab
// 01:23:45:67:89:ab:cd:ef
// 01-23-45-67-89-ab
// 01-23-45-67-89-ab-cd-ef
// 0123.4567.89ab
// 0123.4567.89ab.cdef
func IsMAC(str string) bool {
	_, err := net.ParseMAC(str)
	return err == nil
}

// IsHost checks if the string is a valid IP (both v4 and v6) or a valid DNS name
func IsHost(str string) bool {
	return IsIP(str) || IsDNSName(str)
}

// IsMongoID check if the string is a valid hex-encoded representation of a MongoDB ObjectId.
func IsMongoID(str string) bool {
	return IsHexadecimal(str) && (len(str) == 24)
}

// IsLatitude check if a string is valid latitude.
func IsLatitude(str string) bool {
	return rxLatitude.MatchString(str)
}

// IsLongitude check if a string is valid longitude.
func IsLongitude(str string) bool {
	return rxLongitude.MatchString(str)
}

// IsRsaPublicKey check if a string is valid public key with provided length
func IsRsaPublicKey(str string, keylen int) bool {
	block, _ := pem.Decode([]byte(str))
	if block != nil && block.Type != "PUBLIC KEY" {
		return false
	}
	var der []byte
	if block != nil {
		der = block.Bytes
	} else {
		var err error
		if der, err = base64.StdEncoding.DecodeString(str); err != nil {
			return false
		}
	}

	key, err := x509.ParsePKIXPublicKey(der)
	if err != nil {
		return false
	}
	pubkey, ok := key.(*rsa.PublicKey)
	if !ok {
		return false
	}
	bitlen := len(pubkey.N.Bytes()) * 8
	return bitlen == keylen
}

// IsSSN will validate the given string as a U.S. Social Security Number
func IsSSN(str string) bool {
	if str == "" || len(str) != 11 {
		return false
	}
	return rxSSN.MatchString(str)
}

// IsSemver check if string is valid semantic version
func IsSemver(str string) bool {
	return rxSemver.MatchString(str)
}

// IsTime check if string is valid according to given format
func IsTime(str string, format string) bool {
	_, err := time.Parse(format, str)
	return err == nil
}

// IsRFC3339 check if string is valid timestamp value according to RFC3339
func IsRFC3339(str string) bool {
	return IsTime(str, time.RFC3339)
}

// IsRFC3339WithoutZone check if string is valid timestamp value according to RFC3339 which excludes the timezone.
func IsRFC3339WithoutZone(str string) bool {
	return IsTime(str, RF3339WithoutZone)
}

// IsISO4217 check if string is valid ISO currency code. Code must be upper case.
func IsISO4217(str string) bool {
	return iso4217List[str]
}

// ByteLength check string's length
func ByteLength(str []byte, min, max int) bool {
	return len(str) >= min && len(str) <= max
}

// StringLength check string's length (including multi byte strings)
func StringLength(str string, min, max int) bool {
	strLength := utf8.RuneCountInString(str)
	return strLength >= min && strLength <= max
}

// IsIn check if string str is a member of the set of strings params
func IsIn(str string, params ...string) bool {
	for _, param := range params {
		if str == param {
			return true
		}
	}
	return false
}
