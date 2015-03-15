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
	"go/format"
	"text/template"

	"github.com/juju/errgo"
)

func GenerateCode(tplCode string, data interface{}) ([]byte, error) {

	fm := template.FuncMap{
		"quote":    func(s string) string { return "`" + s + "`" },
		"camelize": Camelize,
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
