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

package net

// Method* defines the available methods which this library supports
const (
	MethodHead    = `HEAD`
	MethodGet     = "GET"
	MethodPost    = "POST"
	MethodPut     = "PUT"
	MethodPatch   = "PATCH"
	MethodDelete  = "DELETE"
	MethodTrace   = "TRACE"
	MethodOptions = "OPTIONS"
)

// Courtesy: github.com/labstack/echo

// HTTP methods
const (
	CONNECT = "CONNECT"
	DELETE  = "DELETE"
	GET     = "GET"
	HEAD    = "HEAD"
	OPTIONS = "OPTIONS"
	PATCH   = "PATCH"
	POST    = "POST"
	PUT     = "PUT"
	TRACE   = "TRACE"
)

// Media types
const (
	ApplicationJSON                  = "application/json"
	ApplicationJSONCharsetUTF8       = ApplicationJSON + "; " + CharsetUTF8
	ApplicationJavaScript            = "application/javascript"
	ApplicationJavaScriptCharsetUTF8 = ApplicationJavaScript + "; " + CharsetUTF8
	ApplicationXML                   = "application/xml"
	ApplicationXMLCharsetUTF8        = ApplicationXML + "; " + CharsetUTF8
	ApplicationForm                  = "application/x-www-form-urlencoded"
	ApplicationProtobuf              = "application/protobuf"
	ApplicationMsgpack               = "application/msgpack"
	TextHTML                         = "text/html"
	TextHTMLCharsetUTF8              = TextHTML + "; " + CharsetUTF8
	TextPlain                        = "text/plain"
	TextPlainCharsetUTF8             = TextPlain + "; " + CharsetUTF8
	MultipartForm                    = "multipart/form-data"
	CompressGZIP                     = "gzip"
	CompressDeflate                  = "deflate"
)

// Charset
const (
	CharsetUTF8 = "charset=utf-8"
)

// Headers
const (
	AcceptEncoding     = "Accept-Encoding"
	Authorization      = "Authorization"
	ContentDisposition = "Content-Disposition"
	ContentEncoding    = "Content-Encoding"
	ContentLength      = "Content-Length"
	ContentType        = "Content-Type"

	// Content-Signature: keyId="rsa-key-1",algorithm="rsa-sha256",signature="Base64(RSA-SHA256(signing string))"
	// Content-Signature: keyId="hmac-key-1",algorithm="hmac-sha1",signature="Base64(HMAC-SHA1(signing string))"
	ContentSignature = "Content-Signature"

	Location         = "Location"
	Upgrade          = "Upgrade"
	Vary             = "Vary"
	WWWAuthenticate  = "WWW-Authenticate"
	XForwarded       = "X-Forwarded"
	XForwardedFor    = "X-Forwarded-For"
	XRealIP          = "X-Real-Ip"
	ClientIP         = "Client-Ip"
	Forwarded        = "Forwarded"
	ForwardedFor     = "Forwarded-For"
	XClusterClientIP = "X-Cluster-Client-Ip"
)

// Protocols
const (
	WebSocket = "websocket"
)
