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

// Package csjwt handles JSON web tokens.
//
// See README.md for more info.
// http://self-issued.info/docs/draft-jones-json-web-token.html
//
// Further reading: https://float-middle.com/json-web-tokens-jwt-vs-sessions/
// and http://cryto.net/~joepie91/blog/2016/06/13/stop-using-jwt-for-sessions/
//
// https://news.ycombinator.com/item?id=11929267 => For people using JWT as a
// substitute for stateful sessions, how do you handle renewal (or revocation)?
//
// https://news.ycombinator.com/item?id=14290114 => Things to Use Instead of JSON Web Tokens (inburke.com)
// TL;DR: Refactor the library and strip out RSA/ECDSA/encoding/decoding into its own sub-packages.
//
// A new discussion: https://news.ycombinator.com/item?id=13865459 JSON Web
// Tokens should be avoided (paragonie.com)
//
// TODO: Investigate security bugs: http://blogs.adobe.com/security/2017/03/critical-vulnerability-uncovered-in-json-encryption.html
// Critical Vulnerability Uncovered in JSON Encryption. Executive Summary: If
// you are using go-jose, node-jose, jose2go, Nimbus JOSE+JWT or jose4 with
// ECDH-ES please update to the latest version. RFC 7516 aka JSON Web Encryption
// (JWE) Invalid Curve Attack. This can allow an attacker to recover the secret
// key of a party using JWE with Key Agreement with Elliptic Curve
// Diffie-Hellman Ephemeral Static (ECDH-ES), where the sender could extract
// receiver’s private key.
package csjwt
