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

package observer

import (
	"fmt"
	"testing"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/pkg/config"
	"github.com/corestoreio/pkg/util/assert"
)

var _ config.Observer = (*ValidateMinMaxInt)(nil)

func TestMinMaxInt_Observe(t *testing.T) {
	t.Parallel()
	var p config.Path
	t.Run("parse failed", func(t *testing.T) {
		mm, err := NewValidateMinMaxInt(1, 2)
		assert.NoError(t, err)
		_, err = mm.Observe(p, []byte("NAN"), false)
		assert.EqualError(t, err, "strconv.ParseInt: parsing \"NAN\": invalid syntax")
	})
	t.Run("null", func(t *testing.T) {
		mm, err := NewValidateMinMaxInt(1, 2)
		assert.NoError(t, err)
		ret, err := mm.Observe(p, nil, false)
		assert.NoError(t, err)
		assert.Nil(t, ret)
	})
	t.Run("not in range1", func(t *testing.T) {
		mm, err := NewValidateMinMaxInt(1, 2)
		assert.NoError(t, err)
		ret, err := mm.Observe(p, []byte(`3`), false)
		assert.True(t, errors.OutOfRange.Match(err), "%+v", err)
		assert.Nil(t, ret)
	})
	t.Run("not in range2", func(t *testing.T) {
		mm, err := NewValidateMinMaxInt(2, 1)
		assert.NoError(t, err)
		ret, err := mm.Observe(p, []byte(`3`), false)
		assert.True(t, errors.OutOfRange.Match(err), "%+v", err)
		assert.Nil(t, ret)
	})
	t.Run("in range1", func(t *testing.T) {
		mm, err := NewValidateMinMaxInt(1, 2)
		assert.NoError(t, err)
		data := []byte(`2`)
		ret, err := mm.Observe(p, data, false)
		assert.NoError(t, err)
		assert.Exactly(t, data, ret)
	})
	t.Run("in range2", func(t *testing.T) {
		mm, err := NewValidateMinMaxInt(1, 2)
		assert.NoError(t, err)
		data := []byte(`2`)
		ret, err := mm.Observe(p, data, false)
		assert.NoError(t, err)
		assert.Exactly(t, data, ret)
	})

	t.Run("partial validation enabled success", func(t *testing.T) {
		mm, err := NewValidateMinMaxInt(1, 2, 5, 6, 7, 8)
		assert.NoError(t, err)
		mm.PartialValidation = true
		data := []byte(`6`)
		ret, err := mm.Observe(p, data, false)
		assert.NoError(t, err)
		assert.Exactly(t, data, ret)
	})

	t.Run("partial validation disabled fails", func(t *testing.T) {
		mm, err := NewValidateMinMaxInt(1, 2, 5, 6, 7, 8)
		assert.NoError(t, err)
		mm.PartialValidation = false
		data := []byte(`6`)
		ret, err := mm.Observe(p, data, false)
		assert.True(t, errors.OutOfRange.Match(err), "%+v", err)
		assert.Nil(t, ret)
	})
}

type observerRegistererFake struct {
	t             *testing.T
	wantEvent     uint8
	wantRoute     string
	wantValidator interface{}
	err           error
}

func (orf observerRegistererFake) RegisterObserver(event uint8, route string, o config.Observer) error {
	if orf.err != nil {
		return orf.err
	}
	if orf.wantEvent != event {
		assert.Exactly(orf.t, orf.wantEvent, event, "Event should be equal: have %d want %d", orf.wantEvent, event)
	}
	if orf.wantRoute != route {
		assert.Exactly(orf.t, orf.wantRoute, route, "Routes")
	}

	// Pointers are different in the final objects hence they get printed and
	// their structure compared, not the address.
	if want, have := fmt.Sprintf("%#v", orf.wantValidator), fmt.Sprintf("%#v", o); want != have {
		assert.Exactly(orf.t, want, have, "Observer internal types should match")
	}
	return nil
}

func (orf observerRegistererFake) DeregisterObserver(event uint8, route string) error {
	if orf.err != nil {
		return orf.err
	}
	if orf.wantEvent != event {
		assert.Exactly(orf.t, orf.wantEvent, event, "Event should be equal: have %d want %d", orf.wantEvent, event)
	}
	if orf.wantRoute != route {
		assert.Exactly(orf.t, orf.wantRoute, route, "Routes")
	}
	return nil
}

func sl(s ...string) []string { return s }

func TestNewValidator(t *testing.T) {
	t.Parallel()

	runner := func(
		validationType []string,
		allowedValues []string,
		csvComma string,
		partialValidation bool,
		data []byte,
		found bool,
		wantNewErr errors.Kind,
		wantObserveErr errors.Kind,
	) func(*testing.T) {
		return func(t *testing.T) {
			s, err := NewValidator(ValidatorArg{
				Funcs:                   validationType,
				PartialValidation:       partialValidation,
				AdditionalAllowedValues: allowedValues,
				CSVComma:                csvComma,
			})
			if !wantNewErr.Empty() {
				assert.Nil(t, s, "validation object s should be nil")
				assert.True(t, wantNewErr.Match(err), "%+v", err)
				return
			}
			assert.NoError(t, err)
			assert.NotNil(t, s)

			haveData, haveErr := s.Observe(config.Path{}, data, found)
			if !wantObserveErr.Empty() {
				assert.Nil(t, haveData, "returned haveData should be nil: %q", string(haveData))
				assert.True(t, wantObserveErr.Match(haveErr), "%+v", haveErr)
				return
			}
			assert.NoError(t, haveErr)
			assert.Exactly(t, data, haveData)
		}
	}

	t.Run("unsupported validationType",
		runner(sl("LifeUniverseEverything"), sl(), "", false, []byte(`42`), true, errors.NotSupported, errors.NoKind),
	)
	t.Run("Custom type is empty",
		runner(sl("Custom"), sl(), "", false, []byte(`42`), true, errors.Empty, errors.NoKind),
	)
	t.Run("Custom type validated",
		runner(sl("Custom"), sl("42", "43"), "", false, []byte(`43`), true, errors.NoKind, errors.NoKind),
	)
	t.Run("Custom type invalid",
		runner(sl("Custom"), sl("42", "43"), "", false, []byte(`44`), true, errors.NoKind, errors.NotValid),
	)
	t.Run("ISO3166Alpha2 validated CSV correct",
		runner(sl("ISO3166Alpha2"), sl(), ",", false, []byte(`DE,CH`), true, errors.NoKind, errors.NoKind),
	)
	t.Run("ISO3166Alpha2 validated  trailing",
		runner(sl("ISO3166Alpha2"), sl(), "@", false, []byte(`DE@CH@`), true, errors.NoKind, errors.NoKind),
	)
	t.Run("ISO3166Alpha2 validated CSV3 heading",
		runner(sl("ISO3166Alpha2"), sl(), ",", false, []byte(`,DE,CH`), true, errors.NoKind, errors.NoKind),
	)
	t.Run("ISO3166Alpha2 not valid CSV1",
		runner(sl("ISO3166Alpha2"), sl(), ",", false, []byte(`,DE,YX`), true, errors.NoKind, errors.NotValid),
	)
	t.Run("ISO3166Alpha2 not valid CSV2 with rune",
		runner(sl("ISO3166Alpha2"), sl(), "", false, []byte(`YX`), true, errors.NoKind, errors.NotValid),
	)
	t.Run("ISO3166Alpha2 not valid CSV3",
		runner(sl("ISO3166Alpha2"), sl(), ",", false, []byte(`YX`), true, errors.NoKind, errors.NotValid),
	)
	t.Run("country_codes2 input not utf8",
		runner(sl("country_codes2"), sl("\xc0\x80"), ",", false, []byte("DE,\xc0\x80"), true, errors.NoKind, errors.NotValid),
	)

	t.Run("ISO3166Alpha3 validated CSV correct",
		runner(sl("ISO3166Alpha3"), sl("XXX"), ";", false, []byte(`DEU;CHE;XXX`), true, errors.NoKind, errors.NoKind),
	)
	t.Run("country_codes3 validated CSV incorrect",
		runner(sl("country_codes3"), sl("FRA"), ";", false, []byte(`FRA;CHE;XXX`), true, errors.NoKind, errors.NotValid),
	)

	t.Run("ISO4217 validated CSV correct",
		runner(sl("ISO4217"), sl("XXY"), ";", false, []byte(`EUR;CHE;XXX;XXY`), true, errors.NoKind, errors.NoKind),
	)
	t.Run("currency3 validated correct",
		runner(sl("currency3"), sl(), "", false, []byte(`CHF`), true, errors.NoKind, errors.NoKind),
	)

	t.Run("Locale validated CSV correct",
		runner(sl("Locale"), sl(), ";", false, []byte(`en-US;de_DE;frfr`), true, errors.NoKind, errors.NoKind),
	)
	t.Run("locale validated correct",
		runner(sl("locale"), sl(), "", false, []byte(`fr-BE`), true, errors.NoKind, errors.NoKind),
	)
	t.Run("locale validated in correct",
		runner(sl("locale"), sl(), "", false, []byte(`fr-DE`), true, errors.NoKind, errors.NotValid),
	)

	t.Run("ISO693Alpha2 validated CSV correct",
		runner(sl("ISO693Alpha2"), sl(), "Ø", false, []byte(`myØzhØce`), true, errors.NoKind, errors.NoKind),
	)
	t.Run("language2 validated correct",
		runner(sl("language2"), sl(), "", false, []byte(`da`), true, errors.NoKind, errors.NoKind),
	)
	t.Run("language2 validated incorrect",
		runner(sl("language2"), sl(), "", false, []byte(`XQ`), true, errors.NoKind, errors.NotValid),
	)

	t.Run("ISO693Alpha3 validated nil",
		runner(sl("ISO693Alpha3"), sl(), "Ø", false, nil, true, errors.NoKind, errors.NoKind),
	)

	t.Run("ISO693Alpha3 validated CSV correct",
		runner(sl("ISO693Alpha3"), sl(), "Ø", false, []byte(`araØaveØdan`), true, errors.NoKind, errors.NoKind),
	)
	t.Run("ISO693Alpha3 validated correct",
		runner(sl("ISO693Alpha3"), sl(), "", false, []byte(`ger`), true, errors.NoKind, errors.NoKind),
	)
	t.Run("language3 validated incorrect",
		runner(sl("language3"), sl(), "", false, []byte(`XxQ`), true, errors.NoKind, errors.NotValid),
	)

	t.Run("uuid validated CSV correct",
		runner(sl("uuid"), sl(), ",", false, []byte(`a987fbc9-4bed-3078-cf07-9141ba07c9f3,a987fbc9-4bed-3078-cf07-9141ba07c9f4`), true, errors.NoKind, errors.NoKind),
	)
	t.Run("uuid validated correct",
		runner(sl("uuid"), sl(), "", false, []byte(`a987fbc9-4bed-3078-cf07-9141ba07c9f3`), true, errors.NoKind, errors.NoKind),
	)
	t.Run("uuid validated incorrect",
		runner(sl("uuid"), sl(), "", false, []byte(`a987fbc9-4bed-3078-cf07-9141ba07c9fZ`), true, errors.NoKind, errors.NotValid),
	)

	t.Run("uuid3 validated CSV correct",
		runner(sl("uuid3"), sl(), ",", false, []byte(`a987fbc9-4bed-3078-cf07-9141ba07c9f3,a987fbc9-4bed-3078-cf07-9141ba07c9f4,`), true, errors.NoKind, errors.NoKind),
	)
	t.Run("uuid3 validated correct",
		runner(sl("uuid3"), sl(), "", false, []byte(`a987fbc9-4bed-3068-cf07-9141ba07c9f3`), true, errors.NoKind, errors.NoKind),
	)
	t.Run("uuid3 validated incorrect",
		runner(sl("uuid3"), sl(), "", false, []byte(`a987fbc9-4bed-4078-8f07-9141ba07c9f3`), true, errors.NoKind, errors.NotValid),
	)

	t.Run("uuid4 validated CSV correct",
		runner(sl("uuid4"), sl(), ",", false, []byte(`57b73598-8764-4ad0-a76a-679bb6641eb1,57b73598-8764-4ad0-a76a-679bb6640eb1,`), true, errors.NoKind, errors.NoKind),
	)
	t.Run("uuid4 validated correct",
		runner(sl("uuid4"), sl(), "", false, []byte(`57b73598-8764-4ad0-a76a-679bb6640eb1`), true, errors.NoKind, errors.NoKind),
	)
	t.Run("uuid4 validated incorrect",
		runner(sl("uuid4"), sl(), "", false, []byte(`a987fbc9-4bed-5078-af07-9141ba07c9f3`), true, errors.NoKind, errors.NotValid),
	)

	t.Run("uuid5 validated CSV correct",
		runner(sl("uuid5"), sl(), ",", false, []byte(`987fbc97-4bed-5078-af07-9141ba07c9f3,987fbc97-4bed-5078-af07-9141ba07c9f3,`), true, errors.NoKind, errors.NoKind),
	)
	t.Run("uuid5 validated correct",
		runner(sl("uuid5"), sl(), "", false, []byte(`987fbc97-4bed-5078-af07-9141ba07c9f3`), true, errors.NoKind, errors.NoKind),
	)
	t.Run("uuid5 validated incorrect",
		runner(sl("uuid5"), sl(), "", false, []byte(`a987fbc9-4bed-3078-cf07-9141ba07c9f3`), true, errors.NoKind, errors.NotValid),
	)

	t.Run("url validated correct",
		runner(sl("url"), sl(), "", false, []byte("http://foobar.中文网/"), true, errors.NoKind, errors.NoKind),
	)
	t.Run("url validated incorrect",
		runner(sl("url"), sl(), "", false, []byte("http://foobar.c_o_m"), true, errors.NoKind, errors.NotValid),
	)

	t.Run("int validated csv correct",
		runner(sl("int"), sl(), ",", false, []byte("1,3,4"), true, errors.NoKind, errors.NoKind),
	)
	t.Run("int validated incorrect",
		runner(sl("int"), sl(), "", false, []byte("h"), true, errors.NoKind, errors.NotValid),
	)

	t.Run("float validated csv correct",
		runner(sl("float"), sl(), ",", false, []byte("1.4,3e4,4.0"), true, errors.NoKind, errors.NoKind),
	)
	t.Run("int validated incorrect",
		runner(sl("float"), sl(), "", false, []byte("h"), true, errors.NoKind, errors.NotValid),
	)

	t.Run("bool validated csv correct",
		runner(sl("bool"), sl(), ",", false, []byte("true,1,True"), true, errors.NoKind, errors.NoKind),
	)
	t.Run("bool validated correct",
		runner(sl("bool"), sl(), "", false, []byte("1"), true, errors.NoKind, errors.NoKind),
	)
	t.Run("bool validated incorrect",
		runner(sl("bool"), sl(), "", false, []byte("h"), true, errors.NoKind, errors.NotValid),
	)

	t.Run("notempty validated correct",
		runner(sl("notempty"), sl(), "", false, []byte("1"), true, errors.NoKind, errors.NoKind),
	)
	t.Run("notempty validated incorrect",
		runner(sl("not_empty"), sl(), "", false, []byte(""), true, errors.NoKind, errors.NotValid),
	)

	t.Run("notemptytrimspace validated correct",
		runner(sl("notemptytrimspace"), sl(), "", false, []byte("\t1\r\t\n"), true, errors.NoKind, errors.NoKind),
	)
	t.Run("notemptytrimspace validated incorrect",
		runner(sl("not_empty_trim_space"), sl(), "", false, []byte("  \t\n\t"), true, errors.NoKind, errors.NotValid),
	)

	t.Run("partialValidation disabled validate correct",
		runner(sl("not_empty", "bool"), sl(), "", false, []byte("true"), true, errors.NoKind, errors.NoKind),
	)
	t.Run("partialValidation disabled validate not correct",
		runner(sl("not_empty", "bool"), sl(), "", false, []byte(""), true, errors.NoKind, errors.NotValid),
	)
	t.Run("partialValidation disabled validate not correct",
		runner(sl("not_empty", "bool"), sl(), "", false, []byte("hello"), true, errors.NoKind, errors.NotValid),
	)

	t.Run("partialValidation enabled validate correct",
		runner(sl("int", "bool"), sl(), "", true, []byte("true"), true, errors.NoKind, errors.NoKind),
	)
	t.Run("partialValidation disabled validate incorrect",
		runner(sl("int", "bool"), sl(), "", false, []byte("true"), true, errors.NoKind, errors.NotValid),
	)

	RegisterValidator("is_euro", isEuro)
	t.Run("RegisterValidator is_euro ok",
		runner(sl("is_euro"), sl(), "", false, []byte("€"), true, errors.NoKind, errors.NoKind),
	)
	t.Run("RegisterValidator is_euro nok",
		runner(sl("is_euro"), sl(), "", false, []byte(""), true, errors.NoKind, errors.NotValid),
	)
}

func isEuro(s string) bool {
	return s == "€"
}

func TestMustNewValidator(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			if err, ok := r.(error); ok {
				assert.True(t, errors.NotSupported.Match(err), "%+v", err)
			} else {
				t.Errorf("Panic should contain an error but got:\n%+v", r)
			}
		} else {
			t.Error("Expecting a panic but got nothing")
		}
	}()
	MustNewValidator(ValidatorArg{
		Funcs: []string{"IsPHP"},
	})
}
