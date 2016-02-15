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

package model

import (
	"sync"
	"time"

	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/config/element"
	cfgpath "github.com/corestoreio/csfw/config/path"
	"github.com/corestoreio/csfw/config/source"
	"github.com/corestoreio/csfw/store/scope"
	"github.com/corestoreio/csfw/util/cast"
	"github.com/corestoreio/csfw/util/cserr"
	"github.com/juju/errors"
)

// PkgBackend used for embedding in the PkgBackend type in each package.
// The mutex protects the init process.
type PkgBackend struct {
	sync.Mutex
}

var _ source.Optioner = (*baseValue)(nil)

// Option as an optional argument for the New*() functions.
// To read more about the recursion pattern:
// http://commandcenter.blogspot.com/2014/01/self-referential-functions-and-design.html
type Option func(*baseValue) Option

// WithConfigStructure sets a global PackageConfiguration for retrieving the
// default value of a underlying type and for scope permission checking.
func WithConfigStructure(cfgStruct element.SectionSlice) Option {
	return func(b *baseValue) Option {
		prev := b.ConfigStructure
		b.ConfigStructure = cfgStruct
		return WithConfigStructure(prev)
	}
}

// WithField creates a new SectionSlice and GroupSlice containing this one field.
// The field.ID gets overwritten by the 3rd path parts to match the path.
// Returns nil on error! Errors are stored in the MultiErr field.
func WithField(f *element.Field) Option {
	return func(b *baseValue) Option {
		prev := b.ConfigStructure
		pp, err := b.r.Split()
		if err != nil {
			b.MultiErr = b.AppendErrors(err, errors.Errorf("Route: %s", b.r))
			b.ConfigStructure = nil
			return nil
		}
		f.ID = pp[2]
		b.ConfigStructure = element.MustNewConfiguration(
			&element.Section{
				ID: pp[0],
				Groups: element.NewGroupSlice(
					&element.Group{
						ID:     pp[1],
						Fields: element.NewFieldSlice(f),
					},
				),
			},
		)
		return WithConfigStructure(prev)
	}
}

// WithSource sets a source slice for Options() and validation.
func WithSource(vl source.Slice) Option {
	return func(b *baseValue) Option {
		prev := b.Source
		b.Source = vl
		return WithSource(prev)
	}
}

// WithSourceByString sets a source slice for Options() and validation.
// Wrapper for source.NewByString
func WithSourceByString(pairs ...string) Option {
	return func(b *baseValue) Option {
		prev := b.Source
		b.Source = source.NewByString(pairs...)
		return WithSource(prev)
	}
}

// WithSourceByInt sets a source slice for Options() and validation.
// Wrapper for source.NewByInt
func WithSourceByInt(vli source.Ints) Option {
	return func(b *baseValue) Option {
		prev := b.Source
		b.Source = source.NewByInt(vli)
		return WithSource(prev)
	}
}

// baseValue defines the path in the "core_config_data" table like a/b/c. All other
// types in this package inherits from this path type.
type baseValue struct {
	// MultiErr some errors of the With* option functions gets appended here.
	*cserr.MultiErr

	r cfgpath.Route // contains the path like web/cors/exposed_headers but has no scope

	// ConfigStructure contains the whole package configuration which is used
	// for scope permission checks and retrieving the default value. A nil
	// ConfigStructure gets ignored.
	ConfigStructure element.SectionSlice

	// Source are all available options aka SourceModel in Mage slang.
	// This slice is also used for validation to get and write the correct values.
	// Validation gets triggered only when the slice has been set.
	// The Options() function will be used to access this slice.
	Source source.Slice
}

// NewValue creates a new baseValue type
func NewValue(path string, opts ...Option) baseValue {
	b := baseValue{
		r: cfgpath.NewRoute(path),
	}
	(&b).Option(opts...)
	return b
}

// Option sets the options and returns the last set previous option
func (bv *baseValue) Option(opts ...Option) (previous Option) {
	for _, o := range opts {
		previous = o(bv)
	}
	return
}

// Write writes a value v to the config.Writer without checking if the value
// has changed. Checks if the Scope matches as defined in the non-nil ConfigStructure.
func (bv baseValue) Write(w config.Writer, v interface{}, s scope.Scope, scopeID int64) error {
	if bv.ConfigStructure != nil {
		f, err := bv.ConfigStructure.FindFieldByID(bv.r)
		if err != nil {
			return errors.Mask(err)
		}
		if false == f.Scope.Has(s) {
			return errors.Errorf("Scope permission insufficient: Have '%s'; Want '%s'", s, f.Scope)
		}
	}
	pp, err := bv.ToPath(s, scopeID)
	if err != nil {
		return errors.Mask(err)
	}
	return w.Write(pp, v)
}

// String returns the stringyfied route
func (bv baseValue) String() string {
	return bv.r.String()
}

// ToPath creates a new path.Path bound to a scope.
func (bv baseValue) ToPath(s scope.Scope, scopeID int64) (cfgpath.Path, error) {
	p, err := cfgpath.New(bv.r)
	if err != nil {
		return cfgpath.Path{}, errors.Mask(err)
	}
	return p.Bind(s, scopeID), nil
}

// Route returns a copy of the underlying route.
func (bv baseValue) Route() cfgpath.Route {
	return bv.r.Clone()
}

// InScope checks if a field from a path is allowed for current scope.
// Returns nil on success.
func (bv baseValue) InScope(sg scope.Scoper) (err error) {
	var f *element.Field
	f, err = bv.field(sg)
	if err != nil {
		return
	}
	s, _ := sg.Scope()
	if false == f.Scope.Has(s) {
		err = errors.Errorf("Scope permission insufficient: Have '%s'; Want '%s'", s, f.Scope)
	}
	return
}

// Options returns a source model for all available options for a configuration
// value.
//
// Usually this function gets customized in a sub-type. Customization
// can have different arguments, etc but must always call this function to set
// source slice.
func (bv baseValue) Options() source.Slice {
	return bv.Source
}

// FQ generates a fully qualified configuration path.
// Example: general/country/allow would transform with StrScope scope.StrStores
// and storeID e.g. 4 into: stores/4/general/country/allow
func (bv baseValue) FQ(s scope.Scope, scopeID int64) (string, error) {
	p, err := bv.ToPath(s, scopeID)
	return p.String(), err
}

// field searches for the field in a SectionSlice and checks if the scope in
// ScopedGetter is sufficient.
func (bv baseValue) field(sg scope.Scoper) (f *element.Field, err error) {
	return bv.ConfigStructure.FindFieldByID(bv.r)
}

// lookupString searches the default value in element.SectionSlice and overwrites
// it with a value from ScopedGetter if ScopedGetter is not empty.
// validator can be nil which triggers the default validation method.
func (bv baseValue) lookupString(sg config.ScopedGetter) (string, error) {
	// This code must be kept in sync with other lookup*() functions
	f, err := bv.field(sg)
	if element.NotNotFoundError(err) {
		return "", errors.Maskf(err, "Route %s", bv.r)
	}

	var v string
	if f != nil {
		var err error
		v, err = cast.ToStringE(f.Default)
		if err != nil {
			return "", errors.Mask(err)
		}
	}

	val, err := sg.String(bv.r)
	switch {
	case err == nil: // we found the value in the config service
		v = val
	case config.NotKeyNotFoundError(err):
		err = errors.Maskf(err, "Route %s", bv.r)
	default:
		err = nil // a Err(Section|Group|Field)NotFound error and uninteresting, so reset
	}
	return v, err
}

// ValidateString checks if string v is contained in Source source.Slice.
func (bv baseValue) ValidateString(v string) (err error) {
	if bv.Source != nil && false == bv.Source.ContainsValString(v) {
		jv, jErr := bv.Source.ToJSON()
		err = errors.Errorf("The value '%s' cannot be found within the allowed Options():\n%s\nJSON Error: %s", v, jv, jErr)
	}
	return
}

func (bv baseValue) lookupInt(sg config.ScopedGetter) (int, error) {
	// This code must be kept in sync with other lookup*() functions
	f, err := bv.field(sg)
	if element.NotNotFoundError(err) {
		return 0, errors.Maskf(err, "Route %s", bv.r)
	}

	var v int
	if f != nil {
		var err error
		v, err = cast.ToIntE(f.Default)
		if err != nil {
			return 0, errors.Mask(err)
		}
	}

	val, err := sg.Int(bv.r)
	switch {
	case err == nil: // we found the value in the config service
		v = val
	case config.NotKeyNotFoundError(err):
		err = errors.Maskf(err, "Route %s", bv.r)
	default:
		err = nil // a Err(Section|Group|Field)NotFound error and uninteresting, so reset
	}
	return v, err

}

// ValidateInt checks if int v is contained in non-nil Source source.Slice.
func (bv baseValue) ValidateInt(v int) (err error) {
	if bv.Source != nil && false == bv.Source.ContainsValInt(v) {
		jv, jErr := bv.Source.ToJSON()
		err = errors.Errorf("The value '%d' cannot be found within the allowed Options():\n%s\nJSON Error: %s", v, jv, jErr)
	}
	return
}

func (bv baseValue) lookupFloat64(sg config.ScopedGetter) (float64, error) {
	// This code must be kept in sync with other lookup*() functions
	f, err := bv.field(sg)
	if element.NotNotFoundError(err) {
		return 0, errors.Maskf(err, "Route %s", bv.r)
	}

	var v float64
	if f != nil {
		var err error
		v, err = cast.ToFloat64E(f.Default)
		if err != nil {
			return 0, errors.Mask(err)
		}
	}

	val, err := sg.Float64(bv.r)
	switch {
	case err == nil: // we found the value in the config service
		v = val
	case config.NotKeyNotFoundError(err):
		err = errors.Maskf(err, "Route %s", bv.r)
	default:
		err = nil // a Err(Section|Group|Field)NotFound error and uninteresting, so reset
	}
	return v, err
}

// ValidateFloat64 checks if float64 v is contained in non-nil Source source.Slice.
func (bv baseValue) ValidateFloat64(v float64) (err error) {
	if bv.Source != nil && false == bv.Source.ContainsValFloat64(v) {
		jv, jErr := bv.Source.ToJSON()
		err = errors.Errorf("The value '%.14f' cannot be found within the allowed Options():\n%s\nJSON Error: %s", v, jv, jErr)
	}
	return
}

func (bv baseValue) lookupBool(sg config.ScopedGetter) (bool, error) {
	// This code must be kept in sync with other lookup*() functions
	f, err := bv.field(sg)
	if element.NotNotFoundError(err) {
		return false, errors.Maskf(err, "Route %s", bv.r)
	}

	var v bool
	if f != nil {
		var err error
		v, err = cast.ToBoolE(f.Default)
		if err != nil {
			return false, errors.Mask(err)
		}
	}

	val, err := sg.Bool(bv.r)
	switch {
	case err == nil: // we found the value in the config service
		v = val
	case config.NotKeyNotFoundError(err):
		err = errors.Maskf(err, "Route %s", bv.r)
	default:
		err = nil // a Err(Section|Group|Field)NotFound error and uninteresting, so reset
	}
	return v, err
}

// lookupTime is able to parse available time formats as defined in
// github.com/corestoreio/csfw/util/cast.StringToDate()
func (bv baseValue) lookupTime(sg config.ScopedGetter) (time.Time, error) {
	// This code must be kept in sync with other lookup*() functions
	f, err := bv.field(sg)
	if element.NotNotFoundError(err) {
		return time.Time{}, errors.Maskf(err, "Route %s", bv.r)
	}

	var v time.Time
	if f != nil {
		var err error
		v, err = cast.ToTimeE(f.Default)
		if err != nil {
			return time.Time{}, errors.Mask(err)
		}
	}

	val, err := sg.Time(bv.r)
	switch {
	case err == nil: // we found the value in the config service
		v = val
	case config.NotKeyNotFoundError(err):
		err = errors.Maskf(err, "Route %s", bv.r)
	default:
		err = nil // a Err(Section|Group|Field)NotFound error and uninteresting, so reset
	}
	return v, err
}

// ValidateTime checks if time.Time v is contained in non-nil Source source.Slice.
func (bv baseValue) ValidateTime(v time.Time) (err error) {
	// todo:
	//if bv.Source != nil && false == bv.Source.ContainsValFloat64(v) {
	//	jv, jErr := bv.Source.ToJSON()
	//	err = errors.Errorf("The value '%.14f' cannot be found within the allowed Options():\n%s\nJSON Error: %s", v, jv, jErr)
	//}
	return
}
