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

package httputils_test

import (
	"errors"
	"math"
	"net/http/httptest"
	"testing"
	"text/template"

	"github.com/corestoreio/csfw/net/httputils"
	"github.com/stretchr/testify/assert"
)

var nonMarshallableChannel chan bool

type errorWriter struct {
	*httptest.ResponseRecorder
}

func (errorWriter) Write(p []byte) (n int, err error) {
	return 0, errors.New("Not in the mood to write today")
}

func TestPrintRender(t *testing.T) {
	w := httptest.NewRecorder()
	p := httputils.NewPrinter(w, nil)
	tpl, err := template.New("foo").Parse(`{{define "T"}}Hello, {{.}}!{{end}}`)
	assert.NoError(t, err)
	p.Renderer = tpl
	assert.NoError(t, p.Render(3141, "T", "<script>alert('you have been pwned')</script>"))
	assert.Exactly(t, `Hello, <script>alert('you have been pwned')</script>!`, w.Body.String())
	assert.Exactly(t, 3141, w.Code)
	assert.Equal(t, httputils.TextHTMLCharsetUTF8, w.Header().Get(httputils.ContentType))
}

func TestPrintRenderErrors(t *testing.T) {
	assert.EqualError(t, httputils.NewPrinter(nil, nil).Render(0, "", nil), httputils.ErrRendererNotRegistered.Error())

	w := httptest.NewRecorder()
	p := httputils.NewPrinter(w, nil)
	tpl, err := template.New("foo").Parse(`{{define "T"}}Hello, {{.}}!{{end}}`)
	assert.NoError(t, err)
	p.Renderer = tpl
	assert.EqualError(t, p.Render(3141, "X", nil), "template: no template \"X\" associated with template \"foo\"")
	assert.Exactly(t, ``, w.Body.String())

}

func TestPrintHTML(t *testing.T) {
	w := httptest.NewRecorder()
	p := httputils.NewPrinter(w, nil)

	assert.NoError(t, p.HTML(3141, "Hello %s. Wanna have some %.5f?", "Gophers", math.Pi))
	assert.Exactly(t, `Hello Gophers. Wanna have some 3.14159?`, w.Body.String())
	assert.Exactly(t, 3141, w.Code)
	assert.Equal(t, httputils.TextHTMLCharsetUTF8, w.Header().Get(httputils.ContentType))
}

func TestPrintHTMLError(t *testing.T) {
	w := new(errorWriter)
	w.ResponseRecorder = httptest.NewRecorder()
	p := httputils.NewPrinter(w, nil)

	assert.EqualError(t, p.HTML(31415, "Hello %s", "Gophers"), "Not in the mood to write today")
	assert.Exactly(t, ``, w.Body.String())
	assert.Exactly(t, 31415, w.Code)
	assert.Equal(t, httputils.TextHTMLCharsetUTF8, w.Header().Get(httputils.ContentType))
}

func TestPrintString(t *testing.T) {
	w := httptest.NewRecorder()
	p := httputils.NewPrinter(w, nil)

	assert.NoError(t, p.String(3141, "Hello %s. Wanna have some %.5f?", "Gophers", math.Pi))
	assert.Exactly(t, `Hello Gophers. Wanna have some 3.14159?`, w.Body.String())
	assert.Exactly(t, 3141, w.Code)
	assert.Equal(t, httputils.TextPlain, w.Header().Get(httputils.ContentType))
}

func TestPrintStringByte(t *testing.T) {
	w := httptest.NewRecorder()
	p := httputils.NewPrinter(w, nil)

	assert.NoError(t, p.StringByte(3141, "Hello %s. Wanna have some %.5f?"))
	assert.Exactly(t, `Hello %s. Wanna have some %.5f?`, w.Body.String())
	assert.Exactly(t, 3141, w.Code)
	assert.Equal(t, httputils.TextPlain, w.Header().Get(httputils.ContentType))
}

func TestPrintStringError(t *testing.T) {
	w := new(errorWriter)
	w.ResponseRecorder = httptest.NewRecorder()
	p := httputils.NewPrinter(w, nil)

	assert.EqualError(t, p.String(31415, "Hello %s", "Gophers"), "Not in the mood to write today")
	assert.Exactly(t, ``, w.Body.String())
	assert.Exactly(t, 31415, w.Code)
	assert.Equal(t, httputils.TextPlain, w.Header().Get(httputils.ContentType))
}

var jsonData = []struct {
	Title string
	SKU   string
	Price float64
}{
	{"Camera", "323423423", 45.12},
	{"LCD TV", "8785344", 145.99},
}

func TestPrintJSON(t *testing.T) {
	w := httptest.NewRecorder()
	p := httputils.NewPrinter(w, nil)

	assert.NoError(t, p.JSON(3141, jsonData))
	assert.Exactly(t, "[{\"Title\":\"Camera\",\"SKU\":\"323423423\",\"Price\":45.12},{\"Title\":\"LCD TV\",\"SKU\":\"8785344\",\"Price\":145.99}]\n", w.Body.String())
	assert.Exactly(t, 3141, w.Code)
	assert.Equal(t, httputils.ApplicationJSONCharsetUTF8, w.Header().Get(httputils.ContentType))
}

func TestPrintJSONError(t *testing.T) {
	w := httptest.NewRecorder()
	p := httputils.NewPrinter(w, nil)

	assert.EqualError(t, p.JSON(3141, nonMarshallableChannel), "json: unsupported type: chan bool")
	assert.Exactly(t, "", w.Body.String())
	assert.Exactly(t, 200, w.Code)
	assert.Equal(t, "", w.Header().Get(httputils.ContentType))
}

func TestPrintJSONIndent(t *testing.T) {
	w := httptest.NewRecorder()
	p := httputils.NewPrinter(w, nil)

	assert.NoError(t, p.JSONIndent(3141, jsonData, "  ", "\t"))
	assert.Exactly(t, "[\n  \t{\n  \t\t\"Title\": \"Camera\",\n  \t\t\"SKU\": \"323423423\",\n  \t\t\"Price\": 45.12\n  \t},\n  \t{\n  \t\t\"Title\": \"LCD TV\",\n  \t\t\"SKU\": \"8785344\",\n  \t\t\"Price\": 145.99\n  \t}\n  ]", w.Body.String())
	assert.Exactly(t, 3141, w.Code)
	assert.Equal(t, httputils.ApplicationJSONCharsetUTF8, w.Header().Get(httputils.ContentType))
}

func TestPrintJSONIndentError(t *testing.T) {
	w := httptest.NewRecorder()
	p := httputils.NewPrinter(w, nil)

	assert.EqualError(t, p.JSONIndent(3141, nonMarshallableChannel, "  ", "\t"), "json: unsupported type: chan bool")
	assert.Exactly(t, "", w.Body.String())
	assert.Exactly(t, 200, w.Code)
	assert.Equal(t, "", w.Header().Get(httputils.ContentType))
}

func TestPrintJSONP(t *testing.T) {
	w := httptest.NewRecorder()
	p := httputils.NewPrinter(w, nil)

	assert.NoError(t, p.JSONP(3141, "awesomeReact", jsonData))
	assert.Exactly(t, "awesomeReact([{\"Title\":\"Camera\",\"SKU\":\"323423423\",\"Price\":45.12},{\"Title\":\"LCD TV\",\"SKU\":\"8785344\",\"Price\":145.99}]\n);", w.Body.String())
	assert.Exactly(t, 3141, w.Code)
	assert.Equal(t, httputils.ApplicationJavaScriptCharsetUTF8, w.Header().Get(httputils.ContentType))
}

func TestPrintJSONPError(t *testing.T) {
	w := httptest.NewRecorder()
	p := httputils.NewPrinter(w, nil)

	assert.EqualError(t, p.JSONP(3141, "awesomeReact", nonMarshallableChannel), "json: unsupported type: chan bool")
	assert.Exactly(t, "", w.Body.String())
	assert.Exactly(t, 200, w.Code)
	assert.Equal(t, "", w.Header().Get(httputils.ContentType))
}
