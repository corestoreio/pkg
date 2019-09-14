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

package observer_test

import (
	"bytes"
	"testing"

	"github.com/corestoreio/pkg/config"
	"github.com/corestoreio/pkg/config/observer"
)

// BenchmarkMinMaxInt64/partial-4         	50000000	        30.6 ns/op	       0 B/op	       0 allocs/op
// BenchmarkMinMaxInt64/non-partial-4     	50000000	        33.2 ns/op	       0 B/op	       0 allocs/op
func BenchmarkMinMaxInt64(b *testing.B) {
	var p config.Path

	b.Run("partial", func(b *testing.B) {
		mm, err := observer.NewValidateMinMaxInt(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
		if err != nil {
			b.Fatal(err)
		}

		mm.PartialValidation = true
		data := []byte(`6`)

		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			ret, err := mm.Observe(p, data, false)
			if err != nil {
				b.Fatal(err)
			}
			if !bytes.Equal(ret, data) {
				b.Fatal("Unequal return data")
			}
		}
	})

	b.Run("non-partial", func(b *testing.B) {
		mm, err := observer.NewValidateMinMaxInt(2012, 2016, 2016, 2018) // weird ;-)
		if err != nil {
			b.Fatal(err)
		}

		data := []byte(`2016`)

		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			ret, err := mm.Observe(p, data, false)
			if err != nil {
				b.Fatal(err)
			}
			if !bytes.Equal(ret, data) {
				b.Fatal("Unequal return data")
			}
		}
	})
}

// BenchmarkStrings/CSV_Locale_with_AdditionalAllowedValues_non-partial-4         	 1000000	      2810 ns/op	      64 B/op	       7 allocs/op
func BenchmarkStrings_CSV_Locale(b *testing.B) {
	var p config.Path

	b.Run("with Add.AllowedValues non-partial single", func(b *testing.B) {
		s, err := observer.NewValidator(observer.ValidatorArg{
			Funcs:                   []string{"locale"},
			PartialValidation:       false,
			AdditionalAllowedValues: []string{"tlh"}, // tlh for klingon
			CSVComma:                ";",
		})
		if err != nil {
			b.Fatal(err)
		}

		data := []byte(`en-GB;tlh`)

		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			ret, err := s.Observe(p, data, false)
			if err != nil {
				b.Fatal(err)
			}
			if !bytes.Equal(ret, data) {
				b.Fatal("Unequal return data")
			}
		}
	})

	b.Run("with Add.AllowedValues non-partial multi", func(b *testing.B) {
		s, err := observer.NewValidator(observer.ValidatorArg{
			Funcs:                   []string{"locale"},
			PartialValidation:       false,
			AdditionalAllowedValues: []string{"tlh"}, // tlh for klingon
			CSVComma:                ";",
		})
		if err != nil {
			b.Fatal(err)
		}

		data := []byte(`en-GB;tlh`)

		b.ResetTimer()

		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				ret, err := s.Observe(p, data, false)
				if err != nil {
					b.Fatal(err)
				}
				if !bytes.Equal(ret, data) {
					b.Fatal("Unequal return data")
				}
			}
		})
	})
}

func BenchmarkStrings_Simple(b *testing.B) {
	var p config.Path

	b.Run("notempty,bool", func(b *testing.B) {
		s, err := observer.NewValidator(observer.ValidatorArg{
			Funcs: []string{"notempty", "bool"},
		})
		if err != nil {
			b.Fatal(err)
		}
		data := []byte(`1`)
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			ret, err := s.Observe(p, data, false)
			if err != nil {
				b.Fatal(err)
			}
			if !bytes.Equal(ret, data) {
				b.Fatal("Unequal return data")
			}
		}
	})
}
