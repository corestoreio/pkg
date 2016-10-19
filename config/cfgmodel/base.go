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
	"time"

	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/config/cfgpath"
	"github.com/corestoreio/csfw/config/element"
	"github.com/corestoreio/csfw/config/source"
	"github.com/corestoreio/csfw/store/scope"
	"github.com/corestoreio/csfw/util/errors"
)

// optionBox groups different types into one struct to allow multiple option
// functions applied to many different types within this package.
type optionBox struct {
	*baseValue
	*Obscure
	*StringCSV
	*CSV
	*IntCSV
	*Encode
}

// Option as an optional argument for all New*() functions. These options will
// be applied to the underlying unexported baseValue type.
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

// WithField adds a Field to the model. Convenient helper function.
func WithField(f *element.Field) Option {
	return func(b *optionBox) error {
		b.Field = f
		return nil
	}
}

// WithSource sets a source slice for Options() and validation. Convenient
// helper function.
func WithSource(vl source.Slice) Option {
	return func(b *optionBox) error {
		b.Source = vl
		return nil
	}
}

// WithSourceByString sets a source slice for Options() and validation. Wrapper
// for source.NewByString.  Convenient helper function.
func WithSourceByString(pairs ...string) Option {
	return func(b *optionBox) (err error) {
		b.Source, err = source.NewByString(pairs...)
		return errors.Wrap(err, "[cfgmodel] WithSourceByString")
	}
}

// WithSourceByInt sets a source slice for Options() and validation. Wrapper for
// source.NewByInt. Convenient helper function.
func WithSourceByInt(vli source.Ints) Option {
	return func(b *optionBox) error {
		b.Source = source.NewByInt(vli)
		return nil
	}
}

// WithScopeStore sets the initial scope to Store. Not needed when using Fields.
// Convenient helper function.
func WithScopeStore() Option {
	return func(b *optionBox) error {
		b.Scopes = scope.PermStore
		return nil
	}
}

// WithScopeWebsite sets the initial scope to Website. Not needed when using
// Fields. Convenient helper function.
func WithScopeWebsite() Option {
	return func(b *optionBox) error {
		b.Scopes = scope.PermWebsite
		return nil
	}
}

// BaseValue represents a Value Object and defines the path in the
// "core_config_data" table like a/b/c. It gets embedded into other types within
// this package.
type baseValue struct {
	// route contains the path like web/cors/exposed_headers but has no scope.
	// Immutable.
	route cfgpath.Route

	// Scopes defaults to scope.Default and is used as an initial value for
	// triggering the hierarchical fallback in the Get() functions. This value
	// gets overwritten when the field *Field below gets set.
	Scopes scope.Perm

	// Field is used for scope permission checks and retrieving the default
	// value. A nil field gets ignored. Field will be set through the option
	// functions at creation time of the struct.
	Field *element.Field

	// Source are all available options aka SourceModel in Mage slang. This
	// slice is also used for validation to get and write the correct values.
	// Validation gets triggered only when the slice has been set. The Options()
	// function will be used to access this slice.
	Source source.Slice
	// LastError might contain an error when an applied functional option
	// returns an error in any New*() constructor. Exported for testing reasons.
	// Every Get() function in a primitive type checks for this error.
	LastError error
}

// newBaseValue creates a new BaseValue type and applies different options.
// Those options can also be set via the structs direct field.
func newBaseValue(path string, opts ...Option) baseValue {
	b := baseValue{
		route: cfgpath.NewRoute(path),
	}
	b.LastError = (&b).Option(opts...)
	return b
}

// Option sets the options and resets the LastError field to nil.
func (bv *baseValue) Option(opts ...Option) error {
	bv.LastError = nil
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

// HasField returns true if the Field has been set and the Fields ID is not
// empty.
func (bv baseValue) HasField() bool {
	return bv.Field != nil && bv.Field.ID.IsEmpty() == false
}

func (bv baseValue) initScope() (p scope.Perm) {
	p = scope.PermDefault
	if bv.Scopes > 0 {
		p = bv.Scopes
	}
	return
}

// Write writes a value v to the config.Writer without checking if the value has
// changed. Checks if the Scope matches as defined in the non-nil
// ConfigStructure. Error behaviour: Unauthorized
func (bv baseValue) Write(w config.Writer, v interface{}, h scope.TypeID) error {
	pp, err := bv.ToPath(h)
	if err != nil {
		return errors.Wrap(err, "[cfgmodel] baseValue.ToPath")
	}
	return w.Write(pp, v)
}

// String returns the stringyfied route
func (bv baseValue) String() string {
	return bv.route.String()
}

// ToPath creates a new cfgpath.Path bound to a scope. If the argument scope
// does not match the defined scope in the Field, the error behaviour
// Unauthorized gets returned. If no argument has been provided it falls back to
// the default scope with ID 0.
//
// If you need a string returned, consider calling FQ() or
// MustFQ*(). FQ = fully qualified path. The returned route in the
// path is owned by the callee.
func (bv baseValue) ToPath(h ...scope.TypeID) (cfgpath.Path, error) {
	t := scope.DefaultTypeID
	if len(h) == 1 {
		t = h[0]
	}
	if err := bv.inScope(t); err != nil {
		return cfgpath.Path{}, errors.Wrap(err, "[cfgmodel] ToPath")
	}

	p, err := cfgpath.New(bv.route)
	if err != nil {
		return cfgpath.Path{}, errors.Wrapf(err, "[cfgmodel] cfgpath.New: %q %s", bv.route, t)
	}
	p.ScopeID = t
	return p, nil
}

// Route returns a copy of the underlying route.
func (bv baseValue) Route() cfgpath.Route {
	return bv.route.Clone()
}

// InScope checks if a field from a path is allowed for current scope. Returns
// nil on success. Error behaviour: Unauthorized
func (bv baseValue) InScope(h scope.TypeID) error {
	return bv.inScope(h)
}

func (bv baseValue) inScope(h scope.TypeID) (err error) {
	s, _ := h.Unpack()
	if bv.HasField() {
		if !bv.Field.Scopes.Has(s) {
			return errors.NewUnauthorizedf(errScopePermissionInsufficient, h, bv.Field.Scopes, bv)
		}
		return nil
	}
	if perms := bv.initScope(); !perms.Has(s) {
		err = errors.NewUnauthorizedf(errScopePermissionInsufficient, h, perms, bv)
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
// storeID e.g. 4 into: stores/4/general/country/allow If no argument has been
// provided it falls back to the default scope with ID 0.
func (bv baseValue) FQ(h ...scope.TypeID) (string, error) {
	p, err := bv.ToPath(h...)
	return p.String(), errors.Wrap(err, "[cfgmodel] ToPath")
}

// MustFQ same as FQ but panics on error. Please use only for testing. If no
// argument has been provided it falls back to the default scope with ID 0.
func (bv baseValue) MustFQ(h ...scope.TypeID) string {
	p, err := bv.ToPath(h...)
	if err != nil {
		panic(err)
	}
	return p.String()
}

// MustFQ same as FQ but for scope website and panics on error. Please use only
// for testing.
func (bv baseValue) MustFQWebsite(id int64) string {
	p, err := bv.ToPath(scope.Website.Pack(id))
	if err != nil {
		panic(err)
	}
	return p.String()
}

// MustFQStore same as FQ but for scope store and panics on error. Please use
// only for testing.
func (bv baseValue) MustFQStore(id int64) string {
	p, err := bv.ToPath(scope.Store.Pack(id))
	if err != nil {
		panic(err)
	}
	return p.String()
}

// IsSet checks if the route has been set.
func (bv baseValue) IsSet() bool {
	return bv.route.IsEmpty() == false
}

// ValidateString checks if string v is contained in Source source.Slice.
// Error behaviour: NotValid
func (bv baseValue) ValidateString(v string) (err error) {
	if bv.Source != nil && false == bv.Source.ContainsValString(v) {
		jv, jErr := bv.Source.ToJSON()
		if jErr != nil {
			return errors.NewFatal(err, fmt.Sprintf("[cfgmodel] Source: %#v", bv.Source))
		}
		err = errors.NewNotValidf(errValueNotFoundInOptions, v, jv)
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
