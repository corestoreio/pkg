// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:generate go run gen.go gen_common.go -output tables.go

// Package currency contains currency-related functionality.
//
// NOTE: the formatting functionality is currently under development and may
// change without notice.
package currency // import "golang.org/x/text/currency"

// All returns all available currencies on this planet.
func All() []string {
	all := make([]string, numCurrencies)
	for i := 0; i < numCurrencies; i++ {
		all[i] = currency.Elem(i + 1)[:3]
	}
	return all
}

type regionUnit struct {
	Region string
	Unit
}

// AllRegionsWithUnit returns all regions with its assigned Unit
func AllRegionsWithUnit() []regionUnit {
	all := make([]regionUnit, len(regionToCurrency))
	for i := 0; i < len(regionToCurrency); i++ {
		r := regionToCurrency[i].region
		all[i] = regionUnit{
			Region: string(r>>8) + string(r^((r>>8)<<8)),
			Unit:   Unit{regionToCurrency[i].code},
		}
	}
	return all
}
