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

package validation_test

import (
	"testing"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/pkg/config"
	"github.com/corestoreio/pkg/config/validation"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func sl(s ...string) []string { return s }

func TestNewStrings(t *testing.T) {
	t.Parallel()

	runner := func(
		validationType string,
		allowedValues []string,
		csvComma string,
		data []byte,
		found bool,
		wantNewErr errors.Kind,
		wantObserveErr errors.Kind,
	) func(*testing.T) {
		return func(t *testing.T) {

			s, err := validation.NewStrings(validation.Strings{
				Type: validationType,
				AdditionalAllowedValues: allowedValues,
				CSVComma:                csvComma,
			})
			if wantNewErr > 0 {
				assert.Nil(t, s, "validation object s should be nil")
				assert.True(t, wantNewErr.Match(err), "%+v", err)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, s)

			haveData, haveErr := s.Observe(config.Path{}, data, found)
			if wantObserveErr > 0 {
				assert.Nil(t, haveData, "returned haveData should be nil: %q", string(haveData))
				assert.True(t, wantObserveErr.Match(haveErr), "%+v", haveErr)
				return
			}
			require.NoError(t, haveErr)
			assert.Exactly(t, data, haveData)
		}
	}

	t.Run("unsupported validationType",
		runner("LifeUniverseEverything", sl(), "", []byte(`42`), true, errors.NotSupported, errors.NoKind),
	)
	t.Run("Custom type is empty",
		runner("Custom", sl(), "", []byte(`42`), true, errors.Empty, errors.NoKind),
	)
	t.Run("Custom type validated",
		runner("Custom", sl("42", "43"), "", []byte(`43`), true, errors.NoKind, errors.NoKind),
	)
	t.Run("Custom type invalid",
		runner("Custom", sl("42", "43"), "", []byte(`44`), true, errors.NoKind, errors.NotValid),
	)
	t.Run("ISO3166Alpha2 validated CSV correct",
		runner("ISO3166Alpha2", sl(), ",", []byte(`DE,CH`), true, errors.NoKind, errors.NoKind),
	)
	t.Run("ISO3166Alpha2 validated  trailing",
		runner("ISO3166Alpha2", sl(), "@", []byte(`DE@CH@`), true, errors.NoKind, errors.NoKind),
	)
	t.Run("ISO3166Alpha2 validated CSV3 heading",
		runner("ISO3166Alpha2", sl(), ",", []byte(`,DE,CH`), true, errors.NoKind, errors.NoKind),
	)
	t.Run("ISO3166Alpha2 not valid CSV1",
		runner("ISO3166Alpha2", sl(), ",", []byte(`,DE,YX`), true, errors.NoKind, errors.NotValid),
	)
	t.Run("ISO3166Alpha2 not valid CSV2 with rune",
		runner("ISO3166Alpha2", sl(), "", []byte(`YX`), true, errors.NoKind, errors.NotValid),
	)
	t.Run("ISO3166Alpha2 not valid CSV3",
		runner("ISO3166Alpha2", sl(), ",", []byte(`YX`), true, errors.NoKind, errors.NotValid),
	)
	t.Run("country_codes2 input not utf8",
		runner("country_codes2", sl("\xc0\x80"), ",", []byte("DE,\xc0\x80"), true, errors.NoKind, errors.NotValid),
	)

	t.Run("ISO3166Alpha3 validated CSV correct",
		runner("ISO3166Alpha3", sl("XXX"), ";", []byte(`DEU;CHE;XXX`), true, errors.NoKind, errors.NoKind),
	)
	t.Run("country_codes3 validated CSV incorrect",
		runner("country_codes3", sl("FRA"), ";", []byte(`FRA;CHE;XXX`), true, errors.NoKind, errors.NotValid),
	)

	t.Run("ISO4217 validated CSV correct",
		runner("ISO4217", sl("XXY"), ";", []byte(`EUR;CHE;XXX;XXY`), true, errors.NoKind, errors.NoKind),
	)
	t.Run("currency3 validated correct",
		runner("currency3", sl(), "", []byte(`CHF`), true, errors.NoKind, errors.NoKind),
	)

	t.Run("Locale validated CSV correct",
		runner("Locale", sl(), ";", []byte(`en-US;de_DE;frfr`), true, errors.NoKind, errors.NoKind),
	)
	t.Run("locale validated correct",
		runner("locale", sl(), "", []byte(`fr-BE`), true, errors.NoKind, errors.NoKind),
	)
	t.Run("locale validated in correct",
		runner("locale", sl(), "", []byte(`fr-DE`), true, errors.NoKind, errors.NotValid),
	)

	t.Run("ISO693Alpha2 validated CSV correct",
		runner("ISO693Alpha2", sl(), "Ø", []byte(`myØzhØce`), true, errors.NoKind, errors.NoKind),
	)
	t.Run("language2 validated correct",
		runner("language2", sl(), "", []byte(`da`), true, errors.NoKind, errors.NoKind),
	)
	t.Run("language2 validated incorrect",
		runner("language2", sl(), "", []byte(`XQ`), true, errors.NoKind, errors.NotValid),
	)

	t.Run("ISO693Alpha3 validated nil",
		runner("ISO693Alpha3", sl(), "Ø", nil, true, errors.NoKind, errors.NoKind),
	)

	t.Run("ISO693Alpha3 validated CSV correct",
		runner("ISO693Alpha3", sl(), "Ø", []byte(`araØaveØdan`), true, errors.NoKind, errors.NoKind),
	)
	t.Run("ISO693Alpha3 validated correct",
		runner("ISO693Alpha3", sl(), "", []byte(`ger`), true, errors.NoKind, errors.NoKind),
	)
	t.Run("language3 validated incorrect",
		runner("language3", sl(), "", []byte(`XxQ`), true, errors.NoKind, errors.NotValid),
	)

	t.Run("uuid validated CSV correct",
		runner("uuid", sl(), ",", []byte(`a987fbc9-4bed-3078-cf07-9141ba07c9f3,a987fbc9-4bed-3078-cf07-9141ba07c9f4`), true, errors.NoKind, errors.NoKind),
	)
	t.Run("uuid validated correct",
		runner("uuid", sl(), "", []byte(`a987fbc9-4bed-3078-cf07-9141ba07c9f3`), true, errors.NoKind, errors.NoKind),
	)
	t.Run("uuid validated incorrect",
		runner("uuid", sl(), "", []byte(`a987fbc9-4bed-3078-cf07-9141ba07c9fZ`), true, errors.NoKind, errors.NotValid),
	)

	t.Run("uuid3 validated CSV correct",
		runner("uuid3", sl(), ",", []byte(`a987fbc9-4bed-3078-cf07-9141ba07c9f3,a987fbc9-4bed-3078-cf07-9141ba07c9f4,`), true, errors.NoKind, errors.NoKind),
	)
	t.Run("uuid3 validated correct",
		runner("uuid3", sl(), "", []byte(`a987fbc9-4bed-3068-cf07-9141ba07c9f3`), true, errors.NoKind, errors.NoKind),
	)
	t.Run("uuid3 validated incorrect",
		runner("uuid3", sl(), "", []byte(`a987fbc9-4bed-4078-8f07-9141ba07c9f3`), true, errors.NoKind, errors.NotValid),
	)

	t.Run("uuid4 validated CSV correct",
		runner("uuid4", sl(), ",", []byte(`57b73598-8764-4ad0-a76a-679bb6641eb1,57b73598-8764-4ad0-a76a-679bb6640eb1,`), true, errors.NoKind, errors.NoKind),
	)
	t.Run("uuid4 validated correct",
		runner("uuid4", sl(), "", []byte(`57b73598-8764-4ad0-a76a-679bb6640eb1`), true, errors.NoKind, errors.NoKind),
	)
	t.Run("uuid4 validated incorrect",
		runner("uuid4", sl(), "", []byte(`a987fbc9-4bed-5078-af07-9141ba07c9f3`), true, errors.NoKind, errors.NotValid),
	)

	t.Run("uuid5 validated CSV correct",
		runner("uuid5", sl(), ",", []byte(`987fbc97-4bed-5078-af07-9141ba07c9f3,987fbc97-4bed-5078-af07-9141ba07c9f3,`), true, errors.NoKind, errors.NoKind),
	)
	t.Run("uuid5 validated correct",
		runner("uuid5", sl(), "", []byte(`987fbc97-4bed-5078-af07-9141ba07c9f3`), true, errors.NoKind, errors.NoKind),
	)
	t.Run("uuid5 validated incorrect",
		runner("uuid5", sl(), "", []byte(`a987fbc9-4bed-3078-cf07-9141ba07c9f3`), true, errors.NoKind, errors.NotValid),
	)

	t.Run("url validated correct",
		runner("url", sl(), "", []byte("http://foobar.中文网/"), true, errors.NoKind, errors.NoKind),
	)
	t.Run("url validated incorrect",
		runner("url", sl(), "", []byte("http://foobar.c_o_m"), true, errors.NoKind, errors.NotValid),
	)

	t.Run("int validated csv correct",
		runner("int", sl(), ",", []byte("1,3,4"), true, errors.NoKind, errors.NoKind),
	)
	t.Run("int validated incorrect",
		runner("int", sl(), "", []byte("h"), true, errors.NoKind, errors.NotValid),
	)

	t.Run("float validated csv correct",
		runner("float", sl(), ",", []byte("1.4,3e4,4.0"), true, errors.NoKind, errors.NoKind),
	)
	t.Run("int validated incorrect",
		runner("float", sl(), "", []byte("h"), true, errors.NoKind, errors.NotValid),
	)

	t.Run("bool validated csv correct",
		runner("bool", sl(), ",", []byte("true,1,True"), true, errors.NoKind, errors.NoKind),
	)
	t.Run("bool validated correct",
		runner("bool", sl(), "", []byte("1"), true, errors.NoKind, errors.NoKind),
	)
	t.Run("bool validated incorrect",
		runner("bool", sl(), "", []byte("h"), true, errors.NoKind, errors.NotValid),
	)

}

func TestStrings_UnmarshalJSON(t *testing.T) {
	t.Parallel()
	rawJSON := []byte(`{"type":"Locale","csv_comma":"|","additional_allowed_values":["Klingon","Vulcan"]}`)

	var data validation.Strings
	require.NoError(t, data.UnmarshalJSON(rawJSON))

	sv, err := validation.NewStrings(data)
	require.NoError(t, err)
	value := []byte(`de_DE|Vulcan`)
	nValue, err := sv.Observe(config.Path{}, value, true)
	require.NoError(t, err)
	assert.Exactly(t, value, nValue)

}
