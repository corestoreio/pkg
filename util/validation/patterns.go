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

import "regexp"

// TODO remove REGEX as many as possible, already removed 5 regexes

// Basic regular expressions for validating strings
const (
	rxStrEmail          = "^(((([a-zA-Z]|\\d|[!#\\$%&'\\*\\+\\-\\/=\\?\\^_`{\\|}~]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])+(\\.([a-zA-Z]|\\d|[!#\\$%&'\\*\\+\\-\\/=\\?\\^_`{\\|}~]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])+)*)|((\\x22)((((\\x20|\\x09)*(\\x0d\\x0a))?(\\x20|\\x09)+)?(([\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x7f]|\\x21|[\\x23-\\x5b]|[\\x5d-\\x7e]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(\\([\\x01-\\x09\\x0b\\x0c\\x0d-\\x7f]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}]))))*(((\\x20|\\x09)*(\\x0d\\x0a))?(\\x20|\\x09)+)?(\\x22)))@((([a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(([a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])([a-zA-Z]|\\d|-|\\.|_|~|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])*([a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])))\\.)+(([a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(([a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])([a-zA-Z]|\\d|-|_|~|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])*([a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])))\\.?$"
	rxStrCreditCard     = "^(?:4[0-9]{12}(?:[0-9]{3})?|5[1-5][0-9]{14}|6(?:011|5[0-9][0-9])[0-9]{12}|3[47][0-9]{13}|3(?:0[0-5]|[68][0-9])[0-9]{11}|(?:2131|1800|35\\d{3})\\d{11})$"
	rxStrISBN10         = "^(?:[0-9]{9}X|[0-9]{10})$"
	rxStrISBN13         = "^(?:[0-9]{13})$"
	rxStrInt            = "^(?:[-+]?(?:0|[1-9][0-9]*))$"
	rxStrFloat          = "^(?:[-+]?(?:[0-9]+))?(?:\\.[0-9]*)?(?:[eE][\\+\\-]?(?:[0-9]+))?$"
	rxStrRGBcolor       = "^rgb\\(\\s*(0|[1-9]\\d?|1\\d\\d?|2[0-4]\\d|25[0-5])\\s*,\\s*(0|[1-9]\\d?|1\\d\\d?|2[0-4]\\d|25[0-5])\\s*,\\s*(0|[1-9]\\d?|1\\d\\d?|2[0-4]\\d|25[0-5])\\s*\\)$"
	rxStrASCII          = "^[\x00-\x7F]+$"
	rxStrMultibyte      = "[^\x00-\x7F]"
	rxStrFullWidth      = "[^\u0020-\u007E\uFF61-\uFF9F\uFFA0-\uFFDC\uFFE8-\uFFEE0-9a-zA-Z]"
	rxStrHalfWidth      = "[\u0020-\u007E\uFF61-\uFF9F\uFFA0-\uFFDC\uFFE8-\uFFEE0-9a-zA-Z]"
	rxStrPrintableASCII = "^[\x20-\x7E]+$"
	rxStrDataURI        = "^data:.+\\/(.+);base64$"
	rxStrLatitude       = "^[-+]?([1-8]?\\d(\\.\\d+)?|90(\\.0+)?)$"
	rxStrLongitude      = "^[-+]?(180(\\.0+)?|((1[0-7]\\d)|([1-9]?\\d))(\\.\\d+)?)$"
	rxStrDNSName        = `^([a-zA-Z0-9_]{1}[a-zA-Z0-9_-]{0,62}){1}(\.[a-zA-Z0-9_]{1}[a-zA-Z0-9_-]{0,62})*[\._]?$`
	rxStrIP             = `(([0-9a-fA-F]{1,4}:){7,7}[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,7}:|([0-9a-fA-F]{1,4}:){1,6}:[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,5}(:[0-9a-fA-F]{1,4}){1,2}|([0-9a-fA-F]{1,4}:){1,4}(:[0-9a-fA-F]{1,4}){1,3}|([0-9a-fA-F]{1,4}:){1,3}(:[0-9a-fA-F]{1,4}){1,4}|([0-9a-fA-F]{1,4}:){1,2}(:[0-9a-fA-F]{1,4}){1,5}|[0-9a-fA-F]{1,4}:((:[0-9a-fA-F]{1,4}){1,6})|:((:[0-9a-fA-F]{1,4}){1,7}|:)|fe80:(:[0-9a-fA-F]{0,4}){0,4}%[0-9a-zA-Z]{1,}|::(ffff(:0{1,4}){0,1}:){0,1}((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])|([0-9a-fA-F]{1,4}:){1,4}:((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9]))`
	rxStrURLSchema      = `((ftp|tcp|udp|wss?|https?):\/\/)`
	rxStrURLUsername    = `(\S+(:\S*)?@)`
	rxStrURLPath        = `((\/|\?|#)[^\s]*)`
	rxStrURLPort        = `(:(\d{1,5}))`
	rxStrURLIP          = `([1-9]\d?|1\d\d|2[01]\d|22[0-3])(\.(1?\d{1,2}|2[0-4]\d|25[0-5])){2}(?:\.([0-9]\d?|1\d\d|2[0-4]\d|25[0-4]))`
	rxStrURLSubdomain   = `((www\.)|([a-zA-Z0-9]([-\.][-\._a-zA-Z0-9]+)*))`
	rxStrURL            = `^` + rxStrURLSchema + `?` + rxStrURLUsername + `?` + `((` + rxStrURLIP + `|(\[` + rxStrIP + `\])|(([a-zA-Z0-9]([a-zA-Z0-9-_]+)?[a-zA-Z0-9]([-\.][a-zA-Z0-9]+)*)|(` + rxStrURLSubdomain + `?))?(([a-zA-Z\x{00a1}-\x{ffff}0-9]+-?-?)*[a-zA-Z\x{00a1}-\x{ffff}0-9]+)(?:\.([a-zA-Z\x{00a1}-\x{ffff}]{1,}))?))\.?` + rxStrURLPort + `?` + rxStrURLPath + `?$`
	rxStrSSN            = `^\d{3}[- ]?\d{2}[- ]?\d{4}$`
	rxStrWinPath        = `^[a-zA-Z]:\\(?:[^\\/:*?"<>|\r\n]+\\)*[^\\/:*?"<>|\r\n]*$`
	rxStrUnixPath       = `^(/[^/\x00]*)+/?$`
	rxStrSemver         = "^v?(?:0|[1-9]\\d*)\\.(?:0|[1-9]\\d*)\\.(?:0|[1-9]\\d*)(-(0|[1-9]\\d*|\\d*[a-zA-Z-][0-9a-zA-Z-]*)(\\.(0|[1-9]\\d*|\\d*[a-zA-Z-][0-9a-zA-Z-]*))*)?(\\+[0-9a-zA-Z-]+(\\.[0-9a-zA-Z-]+)*)?$"
	hasLowerCase        = ".*[[:lower:]]"
	hasUpperCase        = ".*[[:upper:]]"
)

// Used by IsFilePath func
const (
	// Unknown is unresolved OS type
	Unknown = iota
	// Win is Windows type
	Win
	// Unix is *nix OS types
	Unix
)

// regexp uses a global mutex which will become the bottle in certain operations.
var (
	userRegexp       = regexp.MustCompile("^[a-zA-Z0-9!#$%&'*+/=?^_`{|}~.-]+$")
	hostRegexp       = regexp.MustCompile("^[^\\s]+\\.[^\\s]+$")
	userDotRegexp    = regexp.MustCompile("(^[.]{1})|([.]{1}$)|([.]{2,})")
	rxEmail          = regexp.MustCompile(rxStrEmail)
	rxCreditCard     = regexp.MustCompile(rxStrCreditCard)
	rxISBN10         = regexp.MustCompile(rxStrISBN10)
	rxISBN13         = regexp.MustCompile(rxStrISBN13)
	rxInt            = regexp.MustCompile(rxStrInt)
	rxFloat          = regexp.MustCompile(rxStrFloat)
	rxRGBcolor       = regexp.MustCompile(rxStrRGBcolor)
	rxASCII          = regexp.MustCompile(rxStrASCII)
	rxPrintableASCII = regexp.MustCompile(rxStrPrintableASCII)
	rxMultibyte      = regexp.MustCompile(rxStrMultibyte)
	rxFullWidth      = regexp.MustCompile(rxStrFullWidth)
	rxHalfWidth      = regexp.MustCompile(rxStrHalfWidth)
	rxDataURI        = regexp.MustCompile(rxStrDataURI)
	rxLatitude       = regexp.MustCompile(rxStrLatitude)
	rxLongitude      = regexp.MustCompile(rxStrLongitude)
	rxDNSName        = regexp.MustCompile(rxStrDNSName)
	rxURL            = regexp.MustCompile(rxStrURL)
	rxSSN            = regexp.MustCompile(rxStrSSN)
	rxWinPath        = regexp.MustCompile(rxStrWinPath)
	rxUnixPath       = regexp.MustCompile(rxStrUnixPath)
	rxSemver         = regexp.MustCompile(rxStrSemver)
	rxHasLowerCase   = regexp.MustCompile(hasLowerCase)
	rxHasUpperCase   = regexp.MustCompile(hasUpperCase)
)

// Matches check if string matches the pattern (pattern is regular expression)
// In case of error return false
func Matches(str, pattern string) bool {
	match, err := regexp.MatchString(pattern, str)
	return match && err == nil
}
