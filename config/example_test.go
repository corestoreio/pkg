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

package config_test

import (
	"fmt"

	"github.com/corestoreio/pkg/config"
	"github.com/corestoreio/pkg/config/storage"
	"github.com/corestoreio/pkg/store/scope"
)

func panicIfErr(errs ...error) {
	for _, err := range errs {
		if err != nil {
			panic(fmt.Sprintf("%+v", err))
		}
	}
}

func ExamplePath() {
	p := config.MustMakePath("system/smtp/host")

	fmt.Println(p.String())
	fmt.Println(p.BindWebsite(1).String())
	// alternative way
	fmt.Println(p.BindWebsite(1).String())

	fmt.Println(p.BindStore(3).String())
	// alternative way
	fmt.Println(p.BindStore(3).String())
	// Group is not supported and falls back to default
	fmt.Println(p.Bind(scope.Group.WithID(4)).String())

	p, err := config.MakePath("system/smtp/host")
	if err != nil {
		fmt.Printf("%s\n", err)
	}
	fmt.Println(p.String())

	routes := config.MustMakePath("dev/css/merge_css_files")
	rs, err := routes.Split()
	if err != nil {
		fmt.Printf("%s\n", err)
	}
	fmt.Println("dev/css/merge_css_files => ", rs[0], rs[1], rs[2])

	// Output:
	// default/0/system/smtp/host
	// websites/1/system/smtp/host
	// websites/1/system/smtp/host
	// stores/3/system/smtp/host
	// stores/3/system/smtp/host
	// default/0/system/smtp/host
	// default/0/system/smtp/host
	// dev/css/merge_css_files =>  dev css merge_css_files
}

// ExampleValue shows how to use the Value type to convert the raw byte slice
// value to the appropriate final type.
func ExampleValue() {
	// Default storage engine with build-in in-memory map.
	// the NewService gets only instantiated once during app start up.
	configSrv := config.MustNewService(storage.NewMap(), config.Options{})

	const (
		routeCountries      = "carriers/freeshipping/specificcountry"
		routeListingCount   = "catalog/frontend/list_per_page_values"
		routeCookieLifetime = "web/cookie/cookie_lifetime"
	)

	routesVals := []struct {
		route string
		data  string
	}{
		{routeCountries, `CH,LI,DE`},
		{routeListingCount, `5,10,15,20,25`},
		{routeCookieLifetime, `7200s`},
	}
	var p config.Path
	for _, pv := range routesVals {
		panicIfErr(
			p.Parse(pv.route),
			configSrv.Set(p, []byte(pv.data)),
		)
	}

	// This instance of makeScoped gets created somewhere deep in the program.
	// The numbers 1 and 2 are chosen here randomly.
	scpd := configSrv.Scoped(1, 2)
	cntryVal := scpd.Get(scope.Default, routeCountries)

	countries, err := cntryVal.Strs()
	panicIfErr(err)
	fmt.Printf("%s: %#v\n", routeCountries, countries)
	fmt.Printf("%s: %q\n", routeCountries, cntryVal.UnsafeStr())

	listingCount, err := scpd.Get(scope.Default, routeListingCount).Ints()
	panicIfErr(err)
	fmt.Printf("%s: %#v\n", routeListingCount, listingCount)

	keksLifetime, ok, err := scpd.Get(scope.Default, routeCookieLifetime).Duration()
	panicIfErr(err)
	fmt.Printf("%s: found: %t %#v\n", routeCookieLifetime, ok, keksLifetime)

	// Output:
	// carriers/freeshipping/specificcountry: []string{"CH", "LI", "DE"}
	// carriers/freeshipping/specificcountry: "CH,LI,DE"
	// catalog/frontend/list_per_page_values: []int{5, 10, 15, 20, 25}
	// web/cookie/cookie_lifetime: found: true 7200000000000
}
