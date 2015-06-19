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

package utils_test

import (
	"testing"

	"github.com/corestoreio/csfw/utils"
	"github.com/stretchr/testify/assert"
)

func TestTranslit(t *testing.T) {
	tests := []struct {
		have string
		want string
	}{
		{"Hello World", "Hello World"},
		{"Weiß, Goldmann, Göbel, Weiss, Göthe, Goethe und Götz", "Weiss, Goldmann, Gobel, Weiss, Gothe, Goethe und Gotz"},
		{"I have 5€ @ home", "I have 5euro at home"},
		{"Weiß, Goldmann, Göbel, Weiss, Göthe, Goethe und Götz", "Weiss, Goldmann, Gobel, Weiss, Gothe, Goethe und Gotz"},
		{"I have 5 € @ home", "I have 5 euro at home"},
	}

	for _, test := range tests {
		assert.Equal(t, test.want, string(utils.Translit([]rune(test.have))), "%#v", test)
	}
}

func TestTranslitURL(t *testing.T) {
	tests := []struct {
		have string
		want string
	}{
		{"Hello World", "hello-world"},
		{"Weiß, Goldmann, Göbel, Weiss, Göthe, Goethe und Götz", "weiss-goldmann-gobel-weiss-gothe-goethe-und-gotz"},
		{"I have 5€ @ home", "i-have-5euro-at-home"},
		{"Weiß, Goldmann, Göbel, Weiss, Göthe, Goethe und Götz", "weiss-goldmann-gobel-weiss-gothe-goethe-und-gotz"},
		{"I have 5 € @ home", "i-have-5-euro-at-home"},
		{"I have 5 € @ home  ∏", "i-have-5-euro-at-home"},
		{"  I have 5 € @ home  ∏ ", "i-have-5-euro-at-home"},
		{"", ""},
	}

	for _, test := range tests {
		assert.Equal(t, test.want, string(utils.TranslitURL([]rune(test.have))), "%#v", test)
	}
}

var benchmarkTranslitURL string

// Benchmark_TranslitURL	  300000	      3905 ns/op	     688 B/op	       3 allocs/op
func Benchmark_TranslitURL(b *testing.B) {
	b.ReportAllocs()
	want := "weiss-goldmann-gobel-weiss-gothe-goethe-und-gotz"
	have := []rune("Weiß$ Goldmann: Göbel; Weiss, Göthe, Goethe und Götz")
	for i := 0; i < b.N; i++ {
		benchmarkTranslitURL = string(utils.TranslitURL(have))
		if benchmarkTranslitURL != want {
			b.Errorf("\nWant: %s\nHave: %s\n", want, benchmarkTranslitURL)
		}
	}
}
