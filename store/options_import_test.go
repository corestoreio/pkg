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

// +build csall yaml cue

package store_test

import (
	"sort"
	"testing"

	"github.com/corestoreio/pkg/store"
	"github.com/corestoreio/pkg/util/assert"
)

func TestWithLoadFromYAML(t *testing.T) {

	compareCodes := func(t *testing.T, have func(...string) []string, want ...string) {
		c := have()
		sort.Strings(c)
		assert.Exactly(t, c, want)
	}

	t.Run("structure1", func(t *testing.T) {
		s, err := store.NewService(store.WithLoadFromYAML("testdata/store-structure1.yaml", store.YAMLOptions{}))
		assert.NoError(t, err)
		assert.Len(t, s.Stores().Data, 43)

		ss := s.Stores()
		compareCodes(t, ss.Codes, "pos_ae_en", "pos_al_en", "pos_am_en",
			"pos_az_en", "pos_by_en", "pos_ch_de", "pos_ch_en", "pos_ch_fr",
			"pos_cn_en", "pos_de_en", "pos_dz_en", "pos_eg_en", "pos_es_en",
			"pos_et_en", "pos_fr_en", "pos_ga_en", "pos_gb_en", "pos_ge_en",
			"pos_gh_en", "pos_id_en", "pos_ie_en", "pos_iq_en", "pos_it_en",
			"pos_jo_en", "pos_ke_en", "pos_kz_en", "pos_lb_en", "pos_ly_en",
			"pos_ma_en", "pos_mg_en", "pos_ng_en", "pos_ph_en", "pos_rs_en",
			"pos_ru_en", "pos_th_en", "pos_tn_en", "pos_tr_en", "pos_tz_en",
			"pos_ua_en", "pos_ug_en", "pos_vn_en", "pos_za_en", "pos_zw_en")

		gs := s.Groups()
		compareCodes(t, gs.Codes, "pos_ae", "pos_al", "pos_am", "pos_az",
			"pos_by", "pos_ch", "pos_cn", "pos_de", "pos_dz", "pos_eg", "pos_es",
			"pos_et", "pos_fr", "pos_ga", "pos_gb", "pos_ge", "pos_gh", "pos_id",
			"pos_ie", "pos_iq", "pos_it", "pos_jo", "pos_ke", "pos_kz", "pos_lb",
			"pos_ly", "pos_ma", "pos_mg", "pos_ng", "pos_ph", "pos_rs", "pos_ru",
			"pos_th", "pos_tn", "pos_tr", "pos_tz", "pos_ua", "pos_ug", "pos_vn",
			"pos_za", "pos_zw")

		ws := s.Websites()
		compareCodes(t, ws.Codes, "ae", "al", "am", "az", "by", "ch", "cn",
			"de", "dz", "eg", "es", "et", "fr", "ga", "gb", "ge", "gh", "id", "ie",
			"iq", "it", "jo", "ke", "kz", "lb", "ly", "ma", "mg", "ng", "ph", "rs",
			"ru", "th", "tn", "tr", "tz", "ua", "ug", "vn", "za", "zw")

		haveData := toJSON(s, false)
		assert.LenBetween(t, haveData, 29600, 30000)

	})

	t.Run("structure2", func(t *testing.T) {
		s, err := store.NewService(store.WithLoadFromYAML("testdata/store-structure2.yaml", store.YAMLOptions{}))
		assert.NoError(t, err)
		assert.Len(t, s.Stores().Data, 17)

		ss := s.Stores()
		compareCodes(t, ss.Codes, "pos_de_at", "pos_de_ch", "pos_de_de",
			"pos_de_de_b2b", "pos_en_at", "pos_en_ch", "pos_en_de", "pos_en_gb",
			"pos_en_gd", "pos_en_it", "pos_en_us", "pos_es_ca", "pos_es_es",
			"pos_es_us", "pos_fr_ch", "pos_fr_fr", "pos_it_it")

		gs := s.Groups()
		compareCodes(t, gs.Codes, "pos_at", "pos_at", "pos_ch", "pos_ch",
			"pos_ch", "pos_de", "pos_de", "pos_es", "pos_fr", "pos_gb", "pos_it",
			"pos_it", "pos_us", "pos_us")

		ws := s.Websites()
		compareCodes(t, ws.Codes, "de", "en", "es", "fr", "it")

		haveData := toJSON(s, false)
		assert.LenBetween(t, haveData, 9000, 10000)
	})

}
