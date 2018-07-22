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

//go:generate easyjson $GOFILE

package validation

import (
	"strings"
	"unicode/utf8"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/pkg/config"
	"github.com/corestoreio/pkg/util/validation"
)

// Strings checks if a value or a CSV value is a valid type of the defined field
// "Type" and/or contained within AdditionalAllowedValues.
//easyjson:json
type Strings struct {
	// Type can be:
	// "ISO3166Alpha2","country_codes2" for two letter country codes,
	// "ISO3166Alpha3","country_codes3" for three letter country codes,
	// "ISO4217" for three letter currency codes,
	// "Locale" for locale codes,
	// "ISO693Alpha2" for two letter language codes,
	// "ISO693Alpha3" for three letter language codes.
	// "uuid" for any UUID.
	// "uuid3" for UUID version 3.
	// "uuid4" for UUID version 4.
	// "uuid5" for UUID version 5.
	// "url" for URLs
	// "int" for integers
	// "float" for floating point numbers
	// "bool" for boolean values
	// "not_empty" to proof values is not empty
	// "not_empty_trim_space" to proof that values with trimmed white spaces are not empty
	// "Custom" for any custom checking if the value is contained in the
	// AdditionalAllowedValues map.
	Validators []string `json:"validators,omitempty"`
	// PartialValidation if true only one of the Validators must return true /
	// match the string.
	PartialValidation bool `json:"partial_validation,omitempty"`
	// CSVComma one character to separate the input data. If empty the
	// validation process does not know to validate CSV.
	CSVComma string `json:"csv_comma,omitempty"`
	// AdditionalAllowedValues can be optionally or solely defined to add more
	// allowed values than Validators field defines or if Validators equals
	// "Custom" then AdditionalAllowedValues must have values.
	AdditionalAllowedValues []string `json:"additional_allowed_values,omitempty"`
}

// NewStrings creates a new type specific validator. Argument validationType can be
// any string listed in the documenation for type Strings.
func NewStrings(data Strings) (config.Observer, error) {
	ia := &observeStrings{
		valType:           append([]string{}, data.Validators...), // copy data
		valFns:            make([]func(string) bool, 0, len(data.Validators)),
		partialValidation: data.PartialValidation,
	}

	for _, val := range data.Validators {
		var valFn func(string) bool
		switch val {
		case "ISO3166Alpha2", "country_codes2":
			valFn = validation.IsISO3166Alpha2
		case "ISO3166Alpha3", "country_codes3":
			valFn = validation.IsISO3166Alpha3
		case "ISO4217", "currency3":
			valFn = validation.IsISO4217
		case "Locale", "locale":
			valFn = validation.IsLocale
		case "ISO693Alpha2", "language2":
			valFn = validation.IsISO693Alpha2
		case "ISO693Alpha3", "language3":
			valFn = validation.IsISO693Alpha3b
		case "uuid":
			valFn = validation.IsUUID
		case "uuid3":
			valFn = validation.IsUUIDv3
		case "uuid4":
			valFn = validation.IsUUIDv4
		case "uuid5":
			valFn = validation.IsUUIDv5
		case "url":
			valFn = validation.IsURL
		case "int":
			valFn = validation.IsInt
		case "float":
			valFn = validation.IsFloat
		case "bool":
			valFn = validation.IsBool
		case "notempty", "not_empty":
			valFn = validation.IsNotEmpty
		case "notemptytrimspace", "not_empty_trim_space":
			valFn = validation.IsNotEmptyTrimSpace

		case "Custom":
			if len(data.AdditionalAllowedValues) == 0 {
				return nil, errors.Empty.Newf("[config/validation] For type %q the argument allowedValues cannot be empty.", data.Validators)
			}
		default:
			return nil, errors.NotSupported.Newf("[config/validation] Validators %q not yet supported.", data.Validators)
		}
		if valFn != nil {
			ia.valFns = append(ia.valFns, valFn)
		}
	}
	if len(data.AdditionalAllowedValues) > 0 {
		ia.allowedValues = make(map[string]bool)
		for _, c := range data.AdditionalAllowedValues {
			if utf8.ValidString(c) {
				ia.allowedValues[c] = true
			}
		}
	}

	if data.CSVComma != "" && utf8.RuneCountInString(data.CSVComma) <= 4 {
		rc := []rune(data.CSVComma)
		ia.csvComma = rc[0]
	}

	return ia, nil
}

// observeStrings must be used to prevent race conditions during initialization.
// That is the reason we have a separate struct for JSON handling. Having two
// structs allows to refrain from using Locks.
type observeStrings struct {
	valType           []string
	csvComma          rune
	allowedValues     map[string]bool
	valFns            []func(string) bool
	partialValidation bool
}

func (v *observeStrings) isValid(val string) error {

	var validations int
	for _, valFn := range v.valFns {
		if valFn(val) {
			validations++
			if v.partialValidation {
				return nil
			}
		}
	}
	if lFns := len(v.valFns); lFns > 0 && !v.partialValidation && validations == lFns {
		return nil
	}
	if v.allowedValues[val] {
		return nil
	}
	return errors.NotValid.Newf("[config/validation] The provided value %q can't be validated against %q", val, v.valType)
}

// Observe validates the given rawData value. This functions runs in a hot path.
func (v *observeStrings) Observe(_ config.Path, rawData []byte, found bool) (rawData2 []byte, err error) {

	if !utf8.Valid(rawData) {
		return nil, errors.NotValid.Newf("[config/validation] Input data (length:%d) matches no valid UTF8 rune.", len(rawData))
	}

	rawString := string(rawData)
	bufLen := len(rawString) - 1

	if v.csvComma == 0 {
		if err := v.isValid(rawString); err != nil {
			return nil, errors.WithStack(err)
		}
		return rawData, nil
	}

	var column strings.Builder
	for pos, r := range rawString {
		switch {
		case r == v.csvComma && pos == 0:
			// do nothing
		case r == v.csvComma && pos > 0:
			if err := v.isValid(column.String()); err != nil {
				return nil, errors.WithStack(err)
			}
			column.Reset()
		case pos == bufLen:
			column.WriteRune(r)

			if err := v.isValid(column.String()); err != nil {
				return nil, errors.WithStack(err)
			}
			column.Reset()
		default:
			column.WriteRune(r)
		}
	}

	return rawData, nil
}
