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
	ApplicationForm                  = "application/x-www-form-urlencoded"
	ApplicationGob                   = "application/gob"
	ApplicationJSON                  = "application/json"
	ApplicationJSONCharsetUTF8       = ApplicationJSON + "; " + CharsetUTF8
	ApplicationJavaScript            = "application/javascript"
	ApplicationJavaScriptCharsetUTF8 = ApplicationJavaScript + "; " + CharsetUTF8
	ApplicationMsgpack               = "application/msgpack"
	ApplicationProtobuf              = "application/protobuf"
	ApplicationXML                   = "application/xml"
	ApplicationXMLCharsetUTF8        = ApplicationXML + "; " + CharsetUTF8
	CompressDeflate                  = "deflate"
	CompressGZIP                     = "gzip"
	MultipartForm                    = "multipart/form-data"
	TextHTML                         = "text/html"
	TextHTMLCharsetUTF8              = TextHTML + "; " + CharsetUTF8
	TextPlain                        = "text/plain"
	TextPlainCharsetUTF8             = TextPlain + "; " + CharsetUTF8
)

// Charset
const (
	CharsetUTF8 = "charset=utf-8"
)

// Headers
const (
	AcceptEncoding     = "Accept-Encoding"
	Authorization      = "Authorization"
	ClientIP           = "Client-Ip"
	ContentDisposition = "Content-Disposition"
	ContentEncoding    = "Content-Encoding"
	ContentLength      = "Content-Length"
	ContentSignature   = "Content-Signature"
	ContentType        = "Content-Type"
	Forwarded          = "Forwarded"
	ForwardedFor       = "Forwarded-For"
	Location           = "Location"
	Trailer            = "Trailer"
	Upgrade            = "Upgrade"
	Vary               = "Vary"
	WWWAuthenticate    = "WWW-Authenticate"
	XClusterClientIP   = "X-Cluster-Client-Ip"
	XForwarded         = "X-Forwarded"
	XForwardedFor      = "X-Forwarded-For"
	XRealIP            = "X-Real-Ip"
)

// Protocols
const (
	WebSocket = "websocket"
)
