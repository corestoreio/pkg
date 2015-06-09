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

package i18n_test

import (
	"errors"
	"testing"

	"github.com/corestoreio/csfw/i18n"
	"github.com/stretchr/testify/assert"
)

func TestGetLanguages(t *testing.T) {

	tests := []struct {
		locale  string
		wantErr error
		want    []string
	}{
		{"de_DE", nil, []string{"en_US", "English              (Englisch)", "de_DE", "Deutsch              (Deutsch)", "fr_FR", "français             (Französisch)", "it_IT", "italiano             (Italienisch)", "es_ES", "español              (Spanisch)", "ja_JP", "日本語                  (Japanisch)", "uk_UA", "українська           (Ukrainisch)"}},
		{"en", nil, []string{"en_US", "English              (English)", "de_DE", "Deutsch              (German)", "fr_FR", "français             (French)", "it_IT", "italiano             (Italian)", "es_ES", "español              (Spanish)", "ja_JP", "日本語                  (Japanese)", "uk_UA", "українська           (Ukrainian)"}},
		{"es_ES", nil, []string{"en_US", "English              (inglés)", "de_DE", "Deutsch              (alemán)", "fr_FR", "français             (francés)", "it_IT", "italiano             (italiano)", "es_ES", "español              (español)", "ja_JP", "日本語                  (japonés)", "uk_UA", "українська           (ucraniano)"}},
		{"unknown", errors.New("language: tag is not well-formed"), nil},
	}

	for _, test := range tests {
		haveTag, err := i18n.GetLocaleTag(test.locale)
		if test.wantErr != nil {
			assert.Error(t, err)
			assert.EqualError(t, err, test.wantErr.Error())
		} else {
			assert.NoError(t, err)
			have := i18n.GetLanguages(haveTag)
			assert.EqualValues(t, test.want, have)
		}
	}
}
