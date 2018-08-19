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

// Package observer provides validators and modificators for configuration
// values.
//
// The list of validators and modificators will be extended. Please suggest new
// ones.
//
// Modification can be trim, lower case, upper case, title, base64
// encode/decode, sha256, gzip, gunzip and AES-GCM encrypt/decrypt.
//
// Other encryption algorithms are getting later added.
//
// Note: When using sha256 the fully qualified path gets prefixed to the value.
//
// To enabled HTTP handler or protobuf you must set build tags on the CLI.
// Supported build tags are:
// - `json` for JSON encoding and decoding.
// - `proto` for protocol buffers encoding and registration of modificators via
//    proto services, includes `json`.
// - `http` to enable registration of modificators via HTTP handlers.
// - `csall` for all features.
package observer
