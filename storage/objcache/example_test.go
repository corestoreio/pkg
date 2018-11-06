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

// +build bigcache csall

package objcache_test

import (
	"context"
	"encoding/gob"
	"fmt"
	"log"

	"github.com/allegro/bigcache"
	"github.com/corestoreio/pkg/storage/objcache"
)

type P struct {
	X, Y, Z int
	Name    string
}

type Q struct {
	X, Y *int32
	Name string
}

type R struct {
	Name string
	Rune rune
}

// This example shows the basic usage of the package: Create the objcache
// processor, set some values, get some, re-prime gob, set values get some.
func ExampleWithPooledEncoder() {

	gob.Register(P{})
	gob.Register(Q{})
	gob.Register(R{})

	// Use the gob encoder and prime it with the types.
	tc, err := objcache.NewService(
		nil,
		objcache.NewBigCacheClient(bigcache.Config{}),
		// Playing around? Try removing P{}, Q{}, R{} from the next line and see what happens.
		&objcache.ServiceOptions{
			Codec:        gobCodec{},
			PrimeObjects: []interface{}{P{}, Q{}, R{}},
		},
	)
	if err != nil {
		log.Fatalf("NewService error: %+v", err)
	}

	pythagorasKey := `Pythagoras`
	if err := tc.Set(context.TODO(), pythagorasKey, P{3, 4, 5, "Pythagoras"}, 0); err != nil {
		log.Fatalf("Set error 1: %+v", err)
	}
	treeHouseKey := `TreeHouse`
	if err := tc.Set(context.TODO(), treeHouseKey, P{1782, 1841, 1922, "Treehouse"}, 0); err != nil {
		log.Fatalf("Set error 2: %+v", err)
	}

	// Get from cache and print the values. Get operations are called more frequently
	// than Set operations so we're simulating that with 5 repetitions.
	for i := 0; i < 5; i++ {
		var q Q
		if err := tc.Get(context.TODO(), pythagorasKey, &q); err != nil {
			log.Fatalf("Get error 1: %+v", err)
		}
		fmt.Printf("%q: {%d, %d}\n", q.Name, *q.X, *q.Y)

		if err := tc.Get(context.TODO(), treeHouseKey, &q); err != nil {
			log.Fatalf("Get error: %+v", err)
		}
		fmt.Printf("%q: {%d, %d}\n", q.Name, *q.X, *q.Y)
	}

	// We overwrite the previously set values
	if err := tc.Set(context.TODO(), pythagorasKey, R{"Pythagoras2", 'P'}, 0); err != nil {
		log.Fatalf("Set error 1: %+v", err)
	}
	if err := tc.Set(context.TODO(), treeHouseKey, R{"Treehouse2", 'T'}, 0); err != nil {
		log.Fatalf("Set error 2: %+v", err)
	}

	// Get from cache and print the values. Get operations are called more frequently
	// than Set operations so we're simulating that with 5 repetitions.
	for i := 0; i < 5; i++ {
		var r R
		if err := tc.Get(context.TODO(), pythagorasKey, &r); err != nil {
			log.Fatalf("Get error 3: %+v", err)
		}
		fmt.Printf("%q: {%d}\n", r.Name, r.Rune)

		if err := tc.Get(context.TODO(), treeHouseKey, &r); err != nil {
			log.Fatalf("Get error: %+v", err)
		}
		fmt.Printf("%q: {%d}\n", r.Name, r.Rune)
	}
	// Output:
	//"Pythagoras": {3, 4}
	//"Treehouse": {1782, 1841}
	//"Pythagoras": {3, 4}
	//"Treehouse": {1782, 1841}
	//"Pythagoras": {3, 4}
	//"Treehouse": {1782, 1841}
	//"Pythagoras": {3, 4}
	//"Treehouse": {1782, 1841}
	//"Pythagoras": {3, 4}
	//"Treehouse": {1782, 1841}
	//"Pythagoras2": {80}
	//"Treehouse2": {84}
	//"Pythagoras2": {80}
	//"Treehouse2": {84}
	//"Pythagoras2": {80}
	//"Treehouse2": {84}
	//"Pythagoras2": {80}
	//"Treehouse2": {84}
	//"Pythagoras2": {80}
	//"Treehouse2": {84}
}
