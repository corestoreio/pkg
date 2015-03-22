// Copyright 2015 CoreStore Authors
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

package tools

import (
	"bytes"
	"log"
	"math/rand"
	"strings"
	"text/template"
	"time"

	"go/format"

	"github.com/juju/errgo"
)

var (
	logFatalln = log.Fatalln
	letters    = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
)

const Copyright = `// Copyright 2015 CoreStore Authors
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
`

func GenerateCode(pkg, tplCode string, data interface{}) ([]byte, error) {

	var stripPk = func(t string) string {
		l := len(pkg) + 1
		if len(t) <= l {
			return t
		}
		if t[:l] == pkg+TableNameSeparator {
			return t[l:]
		}
		return t
	}

	fm := template.FuncMap{
		"quote":      func(s string) string { return "`" + s + "`" },
		"prepareVar": func(s string) string { return Camelize(stripPk(s)) },
	}
	codeTpl := template.Must(template.New("tpl_code").Funcs(fm).Parse(tplCode))

	var buf = &bytes.Buffer{}
	err := codeTpl.Execute(buf, data)
	if err != nil {
		return nil, errgo.Mask(err)
	}

	fmt, err := format.Source(buf.Bytes())
	if err != nil {
		return buf.Bytes(), err
	}
	return fmt, nil
}

// Camelize transforms from snake case to camelCase e.g. catalog_product_id to CatalogProductID.
func Camelize(s string) string {
	s = strings.ToLower(s)
	parts := strings.Split(s, "_")
	ret := ""
	for _, p := range parts {
		switch p {
		case "id":
			p = "ID"
			break
		case "idx":
			p = "IDX"
			break
		case "eav":
			p = "EAV"
			break
		}
		ret = ret + strings.Title(p)
	}
	return ret
}

// LogFatal logs an error as fatal with printed location and exists the program.
func LogFatal(err error) {
	if err == nil {
		return
	}
	s := "Error: " + err.Error()
	if err, ok := err.(errgo.Locationer); ok {
		s += " " + err.Location().String()
	}
	logFatalln(s)
}

func randSeq(n int) string {
	rand.Seed(time.Now().UTC().UnixNano())
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
