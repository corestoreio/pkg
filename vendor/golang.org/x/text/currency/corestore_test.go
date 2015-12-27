// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package currency

import "testing"

func TestAll(t *testing.T) {
	all := All()
	if c := all[0]; c != "ADP" {
		t.Errorf("first was %c; want ADP", c)
	}
	if c := all[len(all)-1]; c != "ZWR" {
		t.Errorf("last was %c; want ZWR", c)
	}
}

func TestAllRegionsWithUnit(t *testing.T) {
	all := AllRegionsWithUnit()
	for i := 0; i < len(all); i++ {
		if have := all[i].Unit.String(); all[i].Region == "DE" && have != EUR.String() {
			t.Errorf("DE have currency %s; want EUR", have)
		}
	}
}
