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

//go:generate easyjson -build_tags "csall json http proto" $GOFILE

package observer

import (
	"strings"
	"sync"
	"unicode/utf8"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/pkg/config"
	"github.com/corestoreio/pkg/util/byteconv"
	"github.com/corestoreio/pkg/util/validation"
)

// ValidateFn function signature for a validator.
type ValidateFn func(string) bool

type valReg struct {
	sync.RWMutex
	pool map[string]ValidateFn
}

var validatorRegistry = &valReg{
	pool: map[string]ValidateFn{
		"ISO3166Alpha2":        validation.IsISO3166Alpha2,
		"country_codes2":       validation.IsISO3166Alpha2,
		"ISO3166Alpha3":        validation.IsISO3166Alpha3,
		"country_codes3":       validation.IsISO3166Alpha3,
		"ISO4217":              validation.IsISO4217,
		"currency3":            validation.IsISO4217,
		"Locale":               validation.IsLocale,
		"locale":               validation.IsLocale,
		"ISO693Alpha2":         validation.IsISO693Alpha2,
		"language2":            validation.IsISO693Alpha2,
		"ISO693Alpha3":         validation.IsISO693Alpha3b,
		"language3":            validation.IsISO693Alpha3b,
		"uuid":                 validation.IsUUID,
		"uuid3":                validation.IsUUIDv3,
		"uuid4":                validation.IsUUIDv4,
		"uuid5":                validation.IsUUIDv5,
		"url":                  validation.IsURL,
		"int":                  validation.IsInt,
		"float":                validation.IsFloat,
		"bool":                 validation.IsBool,
		"utf8":                 utf8.ValidString,
		"utf8_digit":           validation.IsUTFDigit,
		"utf8_letter":          validation.IsUTFLetter,
		"utf8_letter_numeric":  validation.IsUTFLetterNumeric,
		"notempty":             validation.IsNotEmpty,
		"not_empty":            validation.IsNotEmpty,
		"notemptytrimspace":    validation.IsNotEmptyTrimSpace,
		"not_empty_trim_space": validation.IsNotEmptyTrimSpace,
		"hexadecimal":          validation.IsHexadecimal,
		"hexcolor":             validation.IsHexcolor,
	},
}

// RegisterValidator adds a custom string validation function to the
// global registry for later use with function NewValidator. Adding an
// entry with an already existing `typeName` overwrites the previous validator.
// `typeName` will be handled case-sensitive.
func RegisterValidator(typeName string, vfn ValidateFn) {
	validatorRegistry.Lock()
	validatorRegistry.pool[typeName] = vfn
	validatorRegistry.Unlock()
}

// ValidatorArg gets used as argument to function NewValidator to create a new
// validation observer implementing various checks. Working on CSV data is
// supported to validate each CSV entity. ValidatorArg implements JSON
// encoding/decoding for creating new observer via HTTP or protocol buffers.
//easyjson:json
type ValidatorArg struct {
	// Funcs sets the list of validator functions which can be:
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
	// "custom" for any custom checking if the value is contained in the
	// AdditionalAllowedValues map.
	// Additional all other custom validator functions registered via
	// RegisterValidator are supported.
	Funcs []string `json:"funcs,omitempty"`
	// PartialValidation if true only one of the Configurations must return true /
	// match the string.
	PartialValidation bool `json:"partial_validation,omitempty"`
	// Insecure enables printing in case of errors the values used in a
	// validation function. This might and will leak sensitive information.
	Insecure bool `json:"insecure,omitempty"`
	// CSVComma one character to separate the input data. If empty the
	// validation process does not know to validate CSV.
	CSVComma string `json:"csv_comma,omitempty"`
	// AdditionalAllowedValues can be optionally or solely defined to add more
	// allowed values than Configurations field defines or if Configurations equals
	// "Custom" then AdditionalAllowedValues must have values.
	AdditionalAllowedValues []string `json:"additional_allowed_values,omitempty"`
}

// NewValidator creates a new type specific validator.
func NewValidator(data ValidatorArg) (config.Observer, error) {
	ia := &validators{
		valType:           append([]string{}, data.Funcs...), // copy data
		valFns:            make([]func(string) bool, 0, len(data.Funcs)),
		partialValidation: data.PartialValidation,
		insecure:          data.Insecure,
	}
	validatorRegistry.RLock()
	defer validatorRegistry.RUnlock()

	for _, val := range data.Funcs {
		var valFn ValidateFn
		switch val {
		case "Custom", "custom":
			if len(data.AdditionalAllowedValues) == 0 {
				return nil, errors.Empty.Newf("[config/observer] For type %q the argument allowedValues cannot be empty.", data.Funcs)
			}
		default:
			var ok bool
			valFn, ok = validatorRegistry.pool[val]
			if !ok {
				return nil, errors.NotSupported.Newf("[config/observer] Configurations %q not yet supported.", data.Funcs)
			}
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

// MustNewValidator same as NewValidator but panics on error.
func MustNewValidator(data ValidatorArg) config.Observer {
	o, err := NewValidator(data)
	if err != nil {
		panic(err)
	}
	return o
}

// validators must be used to prevent race conditions during initialization.
// That is the reason we have a separate struct for JSON handling. Having two
// structs allows to refrain from using Locks.
type validators struct {
	csvComma          rune
	partialValidation bool
	valType           []string
	allowedValues     map[string]bool
	valFns            []func(string) bool
	insecure          bool
}

func (v *validators) isValid(val string) error {

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
	if !v.insecure {
		val = "<redacted>"
	}
	return errors.NotValid.Newf("[config/observer] The value %q can't be validated against %q", val, v.valType)
}

// Observe validates the given rawData value. This functions runs in a hot path.
func (v *validators) Observe(_ config.Path, rawData []byte, found bool) (rawData2 []byte, err error) {

	if !utf8.Valid(rawData) {
		return nil, errors.NotValid.Newf("[config/observer] Input data (length:%d) matches no valid UTF8 rune.", len(rawData))
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

// ValidateMinMaxInt validates if a value is between or in range of min and max.
// Provide the field Conditions as a balanced slice where value n defines min
// and n+1 the max value. For JSON handling, see sub-package `json`.
//easyjson:json
type ValidateMinMaxInt struct {
	Conditions []int64 `json:"conditions,omitempty"`
	// PartialValidation if true only one of min/max pairs must be valid.
	PartialValidation bool `json:"partial_validation,omitempty"`
}

// NewValidateMinMaxInt creates a new observer to check if a value is
// contained between min and max values. Argument MinMax must be balanced slice.
func NewValidateMinMaxInt(minMax ...int64) (*ValidateMinMaxInt, error) {
	return &ValidateMinMaxInt{
		Conditions: minMax,
	}, nil
}

// Observe runs the validation process.
func (v ValidateMinMaxInt) Observe(p config.Path, rawData []byte, found bool) (rawData2 []byte, err error) {
	condLen := len(v.Conditions)
	if condLen%2 == 1 || condLen < 1 {
		return nil, errors.NotAcceptable.Newf("[config/observer] ValidateMinMaxInt does not contain a balanced slice. Len: %d", condLen)
	}

	val, ok, err := byteconv.ParseInt(rawData)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if !ok {
		return rawData, nil
	}
	var validations int
	for i := 0; i < condLen; i += 2 {
		if left, right := v.Conditions[i], v.Conditions[i+1]; validation.InRangeInt64(val, left, right) {
			validations++
			if v.PartialValidation {
				return rawData, nil
			}
		}
	}

	if !v.PartialValidation && validations == condLen/2 {
		return rawData, nil
	}
	return nil, errors.OutOfRange.Newf("[config/observer] %q value out of range: %v", val, v.Conditions)
}
