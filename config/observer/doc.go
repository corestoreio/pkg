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

// Package observer provides validators and modifiers for configuration
// values.
//
// Installed validators: ISO3166Alpha2, country_codes2, ISO3166Alpha3,
// country_codes3, ISO4217, currency3, Locale, locale, ISO693Alpha2, language2,
// ISO693Alpha3, language3, uuid, uuid3, uuid4, uuid5, url, int, float, bool,
// utf8, utf8_digit, utf8_letter, utf8_letter_numeric, notempty, not_empty,
// notemptytrimspace, not_empty_trim_space, hexadecimal, hexcolor and any custom
// validator added via function RegisterValidator.
//
// Installed modifiers: upper, lower, trim, title, base64_encode, base64_decode,
// hex_encode, hex_decode, sha256, gzip, gunzip, AES-GCM encrypt/decrypt and any
// custom modifier added via function RegisterModifier.
//
// The list of validators and modifiers will be extended. Please suggest new
// ones.
//
// Other encryption algorithms are getting later added.
//
// Note: When using sha256 the fully qualified path gets prefixed to the value.
//
// To enabled HTTP handler or protobuf you must set build tags on the CLI.
// Supported build tags are:
//	- `json` for JSON encoding and decoding.
//	- `proto` for protocol buffers encoding and registration of modifiers via proto services, includes `json`.
//	- `http` to enable registration of modifiers via HTTP handlers.
//	- `csall` for all features.
package observer
