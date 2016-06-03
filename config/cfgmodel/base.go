// Copyright 2015-2016, Cyrill @ Schumacher.fm and the CoreStore contributors
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

package cfgmodel

import (
	"fmt"
	"sync"
	"time"

	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/config/cfgpath"
	"github.com/corestoreio/csfw/config/element"
	"github.com/corestoreio/csfw/config/source"
	"github.com/corestoreio/csfw/store/scope"
	"github.com/corestoreio/csfw/util/errors"
)

// PkgBackend used for embedding in the PkgBackend type in each package. The
// mutex protects the init process.
type PkgBackend struct {
	sync.Mutex
}

// optionBox groups different types into one struct to allow multiple option
// functions applied to many different types within this package.
type optionBox struct {
	*baseValue
	*Obscure
	*StringCSV
	*CSV
	*IntCSV
}

// BaseOption as an optional argument for the New*() functions. These options
// will be applied to the underlying unexported baseValue type.
type Option func(*optionBox) error

// WithFieldFromSectionSlice extracts the element.Field from the global
// PackageConfiguration for retrieving the default value of a underlying type
// and for scope permission checking.
func WithFieldFromSectionSlice(cfgStruct element.SectionSlice) Option {
	return func(b *optionBox) error {
		f, _, err := cfgStruct.FindField(b.route)
		if err != nil {
			return errors.Wrap(err, "[cfgmodel] cfgStruct.FindField")
		}
		b.Field = &f
		return nil
	}
}

// WithField adds a Field to the model.
func WithField(f *element.Field) Option {
	return func(b *optionBox) error {
		b.Field = f
		return nil
	}
}

// WithSource sets a source slice for Options() and validation.
func WithSource(vl source.Slice) Option {
	return func(b *optionBox) error {
		b.Source = vl
		return nil
	}
}

// WithSourceByString sets a source slice for Options() and validation.
// Wrapper for source.NewByString
func WithSourceByString(pairs ...string) Option {
	return func(b *optionBox) (err error) {
		b.Source, err = source.NewByString(pairs...)
		return errors.Wrap(err, "[cfgmodel] WithSourceByString")
	}
}

// WithSourceByInt sets a source slice for Options() and validation.
// Wrapper for source.NewByInt
func WithSourceByInt(vli source.Ints) Option {
	return func(b *optionBox) error {
		b.Source = source.NewByInt(vli)
		return nil
	}
}

// baseValue defines the path in the "core_config_data" table like a/b/c. All
// other types in this package inherits from this path type.
type baseValue struct {
	// contains the path like web/cors/exposed_headers but has no scope
	route cfgpath.Route

	// Field is used for scope permission checks and retrieving the default
	// value. A nil field gets ignored. Field will be set through the option
	// functions at creation time of the struct.
	Field *element.Field

	// Source are all available options aka SourceModel in Mage slang. This
	// slice is also used for validation to get and write the correct values.
	// Validation gets triggered only when the slice has been set. The Options()
	// function will be used to access this slice.
	Source source.Slice
	// OptionError might contain an error when an applied function option returns an
	// error. Only used in the function MustNewValue()
	OptionError error
}

// NewValue creates a new baseValue type and the error gets packed into the field
// OptionError which you can check.
func NewValue(path string, opts ...Option) baseValue {
	b := baseValue{
		route: cfgpath.NewRoute(path),
	}
	b.OptionError = (&b).Option(opts...)
	return b
}

// Option sets the options and returns the last set previous option
func (bv *baseValue) Option(opts ...Option) error {
	ob := &optionBox{
		baseValue: bv,
	}
	for _, o := range opts {
		if err := o(ob); err != nil {
			return errors.Wrap(err, "[cfgmodel] baseValue.Option")
		}
	}
	bv = ob.baseValue
	return nil
}

func (bv baseValue) hasField() bool {
	return bv.Field != nil && bv.Field.ID.IsEmpty() == false
}

// Write writes a value v to the config.Writer without checking if the value has
// changed. Checks if the Scope matches as defined in the non-nil
// ConfigStructure.
// Error behaviour: Unauthorized
func (bv baseValue) Write(w config.Writer, v interface{}, s scope.Scope, scopeID int64) error {
	if bv.hasField() {
		if false == bv.Field.Scopes.Has(s) {
			return errors.NewUnauthorizedf("[cfgmodel] Scope permission insufficient: Have '%s'; Want '%s'", s, bv.Field.Scopes)
		}
	}
	pp, err := bv.ToPath(s, scopeID)
	if err != nil {
		return errors.Wrap(err, "[cfgmodel] ToPath")
	}
	return w.Write(pp, v)
}

// String returns the stringyfied route
func (bv baseValue) String() string {
	return bv.route.String()
}

// ToPath creates a new cfgpath.Path bound to a scope. If the argument scope
// does not match the defined scope in the Field, the error behaviour
// Unauthorized gets returned.
func (bv baseValue) ToPath(s scope.Scope, scopeID int64) (cfgpath.Path, error) {

	if err := bv.inScope(s, scopeID); err != nil {
		return cfgpath.Path{}, errors.Wrap(err, "[cfgmodel] ToPath")
	}

	p, err := cfgpath.New(bv.route)
	if err != nil {
		return cfgpath.Path{}, errors.Wrapf(err, "[cfgmodel] cfgpath.New: %q", bv.route)
	}
	return p.Bind(s, scopeID), nil
}

// Route returns a copy of the underlying route.
func (bv baseValue) Route() cfgpath.Route {
	return bv.route.Clone()
}

// InScope checks if a field from a path is allowed for current scope.
// Returns nil on success.
// Error behaviour: Unauthorized
func (bv baseValue) InScope(sg scope.Scoper) (err error) {
	return bv.inScope(sg.Scope())
}

func (bv baseValue) inScope(s scope.Scope, _ int64) (err error) {
	if bv.hasField() && !bv.Field.Scopes.Has(s) {
		err = errors.NewUnauthorizedf("[cfgmodel] Scope permission insufficient: Have '%s'; Want '%s'", s, bv.Field.Scopes)
	}
	return
}

// Options returns a source model for all available options for a configuration
// value.
//
// Usually this function gets customized in a sub-type. Customization can have
// different arguments, etc but must always call this function to set source
// slice.
func (bv baseValue) Options() source.Slice {
	return bv.Source
}

// FQ generates a fully qualified configuration path. Example:
// general/country/allow would transform with StrScope scope.StrStores and
// storeID e.g. 4 into: stores/4/general/country/allow
func (bv baseValue) FQ(s scope.Scope, scopeID int64) (string, error) {
	p, err := bv.ToPath(s, scopeID)
	return p.String(), errors.Wrap(err, "[cfgmodel] ToPath")
}

// MustFQ same as FQ but panics on error. Please use only for testing.
func (bv baseValue) MustFQ(s scope.Scope, scopeID int64) string {
	p, err := bv.ToPath(s, scopeID)
	if err != nil {
		panic(err)
	}
	return p.String()
}

// ValidateString checks if string v is contained in Source source.Slice.
// Error behaviour: NotValid
func (bv baseValue) ValidateString(v string) (err error) {
	if bv.Source != nil && false == bv.Source.ContainsValString(v) {
		jv, jErr := bv.Source.ToJSON()
		if jErr != nil {
			return errors.NewFatal(err, fmt.Sprintf("[cfgmodel] Source: %#v", bv.Source))
		}
		err = errors.NewNotValidf("[cfgmodel] The value '%s' cannot be found within the allowed Options():\n%s", v, jv)
	}
	return
}

// ValidateInt checks if int v is contained in non-nil Source source.Slice.
// Error behaviour: NotValid
func (bv baseValue) ValidateInt(v int) (err error) {
	if bv.Source != nil && false == bv.Source.ContainsValInt(v) {
		jv, jErr := bv.Source.ToJSON()
		if jErr != nil {
			return errors.NewFatal(err, fmt.Sprintf("[cfgmodel] Source: %#v", bv.Source))
		}
		err = errors.NewNotValidf("[cfgmodel] The value '%d' cannot be found within the allowed Options():\n%s", v, jv)
	}
	return
}

// ValidateFloat64 checks if float64 v is contained in non-nil Source source.Slice.
// Error behaviour: NotValid
func (bv baseValue) ValidateFloat64(v float64) (err error) {
	if bv.Source != nil && false == bv.Source.ContainsValFloat64(v) {
		jv, jErr := bv.Source.ToJSON()
		if jErr != nil {
			return errors.NewFatal(err, fmt.Sprintf("[cfgmodel] Source: %#v", bv.Source))
		}
		err = errors.NewNotValidf("[cfgmodel] The value '%.14f' cannot be found within the allowed Options():\n%s", v, jv)
	}
	return
}

// ValidateTime checks if time.Time v is contained in non-nil Source source.Slice.
// Error behaviour: NotValid
func (bv baseValue) ValidateTime(v time.Time) (err error) {
	// todo:
	//if bv.Source != nil && false == bv.Source.ContainsValFloat64(v) {
	//jv, jErr := bv.Source.ToJSON()
	//if jErr != nil {
	//	return errors.NewFatal(err, fmt.Sprintf("[cfgmodel] Source: %#v", bv.Source))
	//}
	//err = errors.NewNotValidf("[cfgmodel] The value '%s' cannot be found within the allowed Options():\n%s", v, jv)
	//}
	return errors.NewNotValidf("[cfgmodel] @todo once someone requires this feature")
}
