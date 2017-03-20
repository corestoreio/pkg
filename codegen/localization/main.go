// Copyright 2015-2017, Cyrill @ Schumacher.fm and the CoreStore contributors
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

package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"text/template"

	"github.com/corestoreio/csfw/codegen"
	"github.com/corestoreio/csfw/codegen/localization/gen"
	"golang.org/x/text/unicode/cldr"
)

func main() {

	fmt.Println("TODO refactor")
	os.Exit(-1)

	gen.Init()

	// Read the CLDR zip file. Autodownloading if file not found
	r := gen.OpenCLDRCoreZip()
	defer r.Close()

	d := &cldr.Decoder{}
	d.SetDirFilter("main", "supplemental")
	d.SetSectionFilter("localeDisplayNames", "numbers")
	data, err := d.DecodeZip(r)
	codegen.LogFatal(err)

	curW := &bytes.Buffer{}
	for _, loc := range data.Locales() {

		if false == codegen.ConfigLocalization.EnabledLocale.Include(loc) {
			continue
		}

		ldml, err := data.LDML(loc)
		codegen.LogFatal(err)
		fmt.Fprintf(os.Stdout, "Generating: %s\n", loc)

		curB := curBuilder{
			w:      curW,
			locale: loc,
			data:   ldml,
		}
		curB.generate()
	}

	tplData := map[string]interface{}{
		"Package":       codegen.ConfigLocalization.Package,
		"CurrencyDicts": curW.String(),
	}

	formatted, err := codegen.GenerateCode(codegen.ConfigLocalization.Package, tplCode, tplData, nil)
	if err != nil {
		fmt.Printf("\n\n%s\n\n", formatted)
		codegen.LogFatal(err)
	}

	codegen.LogFatal(ioutil.WriteFile(codegen.ConfigLocalization.OutputFile, formatted, 0600))
}

type curBuilder struct {
	w      io.Writer
	locale string
	data   *cldr.LDML
}

func (b *curBuilder) generate() {

	var nameData = bytes.Buffer{}
	var nameIDX []uint16
	for _, cur := range b.data.Numbers.Currencies.Currency {
		var d string
		if len(cur.DisplayName) > 0 {
			d = cur.DisplayName[0].Data()
		}
		nameData.WriteString(d)
		nameIDX = append(nameIDX, uint16(nameData.Len()))
	}
	nameIDX = append(nameIDX, uint16(nameData.Len()))

	var codeData = bytes.Buffer{}
	for _, cur := range b.data.Numbers.Currencies.Currency {
		if len(cur.Type) != 3 {
			panic(fmt.Errorf("Expecting 3 character long currency code: %v\n", cur))
		}
		codeData.WriteString(cur.Type)
	}

	var symbolData = bytes.Buffer{}
	var symbolIDX []uint16
	for _, cur := range b.data.Numbers.Currencies.Currency {
		var d string
		if len(cur.Symbol) > 0 {
			d = cur.Symbol[0].Data()
		}
		symbolData.WriteString(d)
		symbolIDX = append(symbolIDX, uint16(symbolData.Len()))
	}
	symbolIDX = append(symbolIDX, uint16(symbolData.Len()))

	tplData := map[string]interface{}{
		"Locale":      b.locale,
		"Codes3":      codeData.String(),
		"NamesData":   nameData.String(),
		"NamesIDX":    nameIDX,
		"SymbolsData": symbolData.String(),
		"SymbolsIDX":  symbolIDX,
	}

	err := template.Must(template.New("tpl_code_unit").Parse(tplCodeUnit)).Execute(b.w, tplData)
	codegen.LogFatal(err)
}
