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

// The MIT License (MIT)
//
// Copyright (c) 2015 LabStack github.com/labstack/echo
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

package response

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"path/filepath"

	"github.com/corestoreio/csfw/net/httputil"
	"github.com/corestoreio/csfw/util/bufferpool"
	"github.com/corestoreio/csfw/util/errors"
)

const indexFile = "index.html"

// NewPrinter creates a non-pointer printer
func NewPrinter(w http.ResponseWriter, r *http.Request) Print {
	return Print{
		Response: w,
		Request:  r,
	}
}

// Print is a helper type for outputting data to a ResponseWriter. Print
// act as a non-pointer type. Print functions uses internally a
// byte buffer pool.
type Print struct {
	// FileSystem stubbed out for testing. Default http.Dir
	FileSystem http.FileSystem
	Response   http.ResponseWriter
	Request    *http.Request
	Renderer   interface {
		// ExecuteTemplate is the interface function to text/template or html/template
		ExecuteTemplate(wr io.Writer, name string, data interface{}) error
	}
}

// Render renders a template with data and sends a text/html response with status
// code. Templates can be registered during `Print` creation.
func (p Print) Render(code int, name string, data interface{}) error {
	if p.Renderer == nil {
		return errors.NewEmptyf("[httputil] Print.Render.Renderer is nil")
	}
	buf := bufferpool.Get()
	defer bufferpool.Put(buf)
	if err := p.Renderer.ExecuteTemplate(buf, name, data); err != nil {
		return errors.NewFatal(err, "[httputil] Print.Render.ExecuteTemplate failed")
	}
	return p.html(code, buf.Bytes())
}

// HTML formats according to a format specifier and sends HTML response with
// status code.
func (p Print) HTML(code int, format string, a ...interface{}) error {
	err := p.html(code, nil)
	if err != nil {
		return errors.Wrap(err, "[httputil] Print.HTML.html")
	}
	_, err = fmt.Fprintf(p.Response, format, a...)
	return errors.NewWriteFailed(err, "[httputil] Print.HTML.Fprintf")
}

func (p Print) html(code int, data []byte) (err error) {
	p.Response.Header().Set(httputil.ContentType, httputil.TextHTMLCharsetUTF8)
	p.Response.WriteHeader(code)
	if data != nil {
		_, err = p.Response.Write(data)
	}
	return errors.NewWriteFailed(err, "[httputil] Print.html.Response.Write")
}

// String formats according to a format specifier and sends text response with
// status code.
func (p Print) String(code int, format string, a ...interface{}) (err error) {
	if err := p.string(code, nil); err != nil {
		return errors.Wrap(err, "[httputil] Print.String.string")
	}
	_, err = fmt.Fprintf(p.Response, format, a...)
	return errors.NewWriteFailed(err, "[httputil] Print.String.Fprintf")
}

// WriteString converts a string into []bytes and outputs it. No formatting
// feature available.
func (p Print) WriteString(code int, s string) (err error) {
	p.Response.Header().Set(httputil.ContentType, httputil.TextPlain)
	p.Response.WriteHeader(code)
	_, err = io.WriteString(p.Response, s)
	return errors.NewWriteFailed(err, "[httputil] Print.WriteString")
}

func (p Print) string(code int, data []byte) (err error) {
	p.Response.Header().Set(httputil.ContentType, httputil.TextPlain)
	p.Response.WriteHeader(code)
	_, err = p.Response.Write(data)
	return errors.NewWriteFailed(err, "[httputil] Print.string")
}

// JSON sends a JSON response with status code.
func (p Print) JSON(code int, i interface{}) (err error) {
	buf := bufferpool.Get()
	defer bufferpool.Put(buf)

	if err := json.NewEncoder(buf).Encode(i); err != nil {
		return errors.NewFatal(err, "[httputil] Print.JSON.NewEncoder.Encode")
	}
	return errors.Wrap(p.json(code, buf.Bytes()), "[httputil] JSON")
}

// JSONIndent sends a JSON response with status code, but it applies prefix and indent to format the output.
func (p Print) JSONIndent(code int, i interface{}, prefix string, indent string) (err error) {
	b, err := json.MarshalIndent(i, prefix, indent)
	if err != nil {
		return err
	}
	return errors.Wrap(p.json(code, b), "[httputil] JSONIndent")
}

func (p Print) json(code int, b []byte) (err error) {
	p.Response.Header().Set(httputil.ContentType, httputil.ApplicationJSONCharsetUTF8)
	p.Response.WriteHeader(code)
	if b != nil {
		_, err = p.Response.Write(b)
	}
	return errors.NewWriteFailed(err, "[httputil] Print.json")
}

// JSONP sends a JSONP response with status code. It uses `callback` to construct
// the JSONP payload.
func (p Print) JSONP(code int, callback string, i interface{}) (err error) {
	buf := bufferpool.Get()
	defer bufferpool.Put(buf)

	buf.WriteString(callback)
	buf.WriteRune('(')

	if err := json.NewEncoder(buf).Encode(i); err != nil {
		return errors.NewFatal(err, "[httputil] Print.JSONP.NewEncoder.Encode")
	}
	buf.WriteString(");")

	p.Response.Header().Set(httputil.ContentType, httputil.ApplicationJavaScriptCharsetUTF8)
	p.Response.WriteHeader(code)

	_, err = p.Response.Write(buf.Bytes())
	return errors.Wrap(err, "[httputil] Print.JSONP.Response.Write")
}

// XML sends an XML response with status code.
func (p Print) XML(code int, i interface{}) (err error) {
	buf := bufferpool.Get()
	defer bufferpool.Put(buf)

	if err := xml.NewEncoder(buf).Encode(i); err != nil {
		return errors.NewFatal(err, "[httputil] Print.XML.NewEncoder.Encode")
	}
	return errors.Wrap(p.xml(code, buf.Bytes()), "[httputil] Print.XML.xml")
}

// XMLIndent sends an XML response with status code, but it applies prefix and indent to format the output.
func (p Print) XMLIndent(code int, i interface{}, prefix string, indent string) (err error) {
	b, err := xml.MarshalIndent(i, prefix, indent)
	if err != nil {
		return errors.NewFatal(err, "[httputil] Print.XMLIndent.MarshalIndent")
	}
	return errors.Wrap(p.xml(code, b), "[httputil] Print.XMLIndent.xml")
}

func (p Print) xml(code int, b []byte) (err error) {
	p.Response.Header().Set(httputil.ContentType, httputil.ApplicationXMLCharsetUTF8)
	p.Response.WriteHeader(code)
	if _, err = p.Response.Write([]byte(xml.Header)); err != nil {
		return errors.Wrap(err, "[httputil] Print.xml")
	}
	if b != nil {
		_, err = p.Response.Write(b)
	}
	return errors.Wrap(err, "[httputil] Print.xml.Response.Write")
}

// File sends a response with the content of the file. If `attachment` is set
// to true, the client is prompted to save the file with provided `name`,
// name can be empty, in that case name of the file is used.
func (p Print) File(path, name string, attachment bool) error {
	dir, file := filepath.Split(path)
	if attachment {
		p.Response.Header().Set(httputil.ContentDisposition, "attachment; filename="+name)
	}
	if err := serveFile(dir, file, p); err != nil {
		p.Response.Header().Del(httputil.ContentDisposition)
		return errors.Wrap(err, "[httputil] Print.File.serveFile")
	}
	return nil
}

// NoContent sends a response with no body and a status code.
func (p Print) NoContent(code int) error {
	p.Response.WriteHeader(code)
	return nil
}

// Redirect redirects the request using http.Redirect with status code.
func (p Print) Redirect(code int, url string) error {
	if code < http.StatusMultipleChoices || code > http.StatusTemporaryRedirect {
		return errors.NewNotValidf("[httputil] Unknown redirect code %d", code)
	}
	http.Redirect(p.Response, p.Request, url, code)
	return nil
}

func serveFile(dir, file string, p Print) error {
	if p.FileSystem == nil {
		p.FileSystem = http.Dir(dir)
	}

	f, err := p.FileSystem.Open(file)
	if err != nil {
		return errors.NewFatalf("[httputil] File not found: %s => %s", dir, file)
	}
	defer f.Close()

	fi, _ := f.Stat()
	if fi.IsDir() {
		file = filepath.Join(file, indexFile)
		f, err = p.FileSystem.Open(file)
		if err != nil {
			return errors.NewFatalf("[httputil] Cannot access file: %s", file) // http.StatusForbidden
		}
		fi, _ = f.Stat()
	}
	http.ServeContent(p.Response, p.Request, fi.Name(), fi.ModTime(), f)
	return nil
}
