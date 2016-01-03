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

	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/config/element"
	"github.com/corestoreio/csfw/config/source"
	"github.com/corestoreio/csfw/store/scope"
	"github.com/corestoreio/csfw/util/cast"
	"github.com/juju/errgo"
)

// PkgPath used for embedding in the PkgPath type in each package.
// The mutex protects the init process.
type PkgPath struct {
	sync.Mutex
}

var _ source.Optioner = (*basePath)(nil)

// Option as an optional argument for the New*() functions.
// To read more about the recursion pattern:
// http://commandcenter.blogspot.com/2014/01/self-referential-functions-and-design.html
type Option func(*basePath) Option

// WithConfigStructure sets a global PackageConfiguration for retrieving the
// default value of a underlying type and for scope permission checking.
func WithConfigStructure(cfgStruct element.SectionSlice) Option {
	return func(b *basePath) Option {
		prev := b.ConfigStructure
		b.ConfigStructure = cfgStruct
		return WithConfigStructure(prev)
	}
}

// WithField creates a new SectionSlice and GroupSlice containing this one field.
// The field.ID gets overwritten by the 3rd path parts to match the path.
func WithField(f *element.Field) Option {
	return func(b *basePath) Option {
		prev := b.ConfigStructure
		pp := scope.PathSplit(b.string)
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
	return func(b *basePath) Option {
		prev := b.Source
		b.Source = vl
		return WithSource(prev)
	}
}

// WithSourceByString sets a source slice for Options() and validation.
// Wrapper for source.NewByString
func WithSourceByString(pairs ...string) Option {
	return func(b *basePath) Option {
		prev := b.Source
		b.Source = source.NewByString(pairs...)
		return WithSource(prev)
	}
}

// WithSourceByInt sets a source slice for Options() and validation.
// Wrapper for source.NewByInt
func WithSourceByInt(vli source.Ints) Option {
	return func(b *basePath) Option {
		prev := b.Source
		b.Source = source.NewByInt(vli)
		return WithSource(prev)
	}
}

// basePath defines the path in the "core_config_data" table like a/b/c. All other
// types in this package inherits from this path type.
type basePath struct {
	string // contains the path

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

// NewPath creates a new basePath type
func NewPath(path string, opts ...Option) basePath {
	b := basePath{
		string: path,
	}
	(&b).Option(opts...)
	return b
}

// Option sets the options and returns the last set previous option
func (p *basePath) Option(opts ...Option) (previous Option) {
	for _, o := range opts {
		previous = o(p)
	}
	return
}

// Write writes a value v to the config.Writer without checking if the value
// has changed. Checks if the Scope matches as defined in the non-nil ConfigStructure.
func (p basePath) Write(w config.Writer, v interface{}, s scope.Scope, id int64) error {
	if p.ConfigStructure != nil {
		f, err := p.ConfigStructure.FindFieldByPath(p.string)
		if err != nil {
			return errgo.Mask(err)
		}
		if false == f.Scope.Has(s) {
			return errgo.Newf("Scope permission insufficient: Have '%s'; Want '%s'", s, f.Scope)
		}
	}
	return w.Write(config.Path(p.string), config.Value(v), config.Scope(s, id))
}

// String returns the path
func (p basePath) String() string {
	return p.string
}

// InScope checks if a field from a path is allowed for current scope.
// Returns nil on success.
func (p basePath) InScope(sg scope.Scoper) (err error) {
	_, err = p.field(sg)
	return
}

// Options returns a source model for all available options for a configuration
// value.
//
// Usually this function gets customized in a sub-type. Customization
// can have different arguments, etc but must always call this function to set
// source slice.
func (p basePath) Options() source.Slice {
	return p.Source
}

// FQPathInt64 generates a fully qualified configuration path.
// Example: general/country/allow would transform with StrScope scope.StrStores
// and storeID e.g. 4 into: stores/4/general/country/allow
func (p basePath) FQPathInt64(strScope scope.StrScope, scopeID int64) string {
	return strScope.FQPathInt64(scopeID, p.string)
}

// field searches for the field in a SectionSlice and checks if the scope in
// ScopedGetter is sufficient.
func (p basePath) field(sg scope.Scoper) (f *element.Field, err error) {
	f, err = p.ConfigStructure.FindFieldByPath(p.string)
	if err != nil {
		return nil, errgo.Mask(err)
	}
	s, _ := sg.Scope()
	if false == f.Scope.Has(s) {
		return nil, errgo.Newf("Scope permission insufficient: Have '%s'; Want '%s'", s, f.Scope)
	}
	return
}

// lookupString searches the default value in element.SectionSlice and overwrites
// it with a value from ScopedGetter if ScopedGetter is not empty.
// validator can be nil which triggers the default validation method.
func (p basePath) lookupString(sg config.ScopedGetter) (v string, err error) {

	if f, errF := p.field(sg); errF == nil {
		v, err = cast.ToStringE(f.Default)
	} else if PkgLog.IsDebug() {
		PkgLog.Debug("model.basePath.lookupString.field", "err", errF, "path", p.string)
	}

	if val, errSG := sg.String(p.string); errSG == nil {
		v = val
	} else if PkgLog.IsDebug() {
		// errSG is usually a key not found error, but that one is uninteresting
		PkgLog.Debug("model.basePath.lookupString.ScopedGetter.String", "err", errSG, "path", p.string, "previousErr", err)
	}
	return
}

func (p basePath) validateString(v string) (err error) {
	if p.Source != nil && false == p.Source.ContainsValString(v) {
		jv, jErr := p.Source.ToJSON()
		err = errgo.Newf("The value '%s' cannot be found within the allowed Options():\n%s\nJSON Error: %s", v, jv, jErr)
	}
	return
}

func (p basePath) lookupInt(sg config.ScopedGetter) (v int, err error) {

	var f *element.Field
	if f, err = p.field(sg); err != nil {
		return
	}

	v, err = cast.ToIntE(f.Default)
	if err != nil {
		err = errgo.Mask(err)
		return
	}

	if val, errSG := sg.Int(p.string); errSG == nil {
		v = val
	} else {
		// errSG is usually a key not found error, but that one is uninteresting
		if PkgLog.IsDebug() {
			PkgLog.Debug("model.path.lookupString.ScopedGetter.Int", "err", errSG, "path", p.string, "previousErr", err)
		}
	}
	return
}

func (p basePath) validateInt(v int) (err error) {
	if p.Source != nil && false == p.Source.ContainsValInt(v) {
		jv, jErr := p.Source.ToJSON()
		err = errgo.Newf("The value '%d' cannot be found within the allowed Options():\n%s\nJSON Error: %s", v, jv, jErr)
	}
	return
}

func (p basePath) lookupFloat64(sg config.ScopedGetter) (v float64, err error) {

	var f *element.Field
	if f, err = p.field(sg); err != nil {
		return
	}

	v, err = cast.ToFloat64E(f.Default)
	if err != nil {
		err = errgo.Mask(err)
		return
	}

	if val, errSG := sg.Float64(p.string); errSG == nil {
		v = val
	} else {
		// errSG is usually a key not found error, but that one is uninteresting
		if PkgLog.IsDebug() {
			PkgLog.Debug("model.path.lookupString.ScopedGetter.Float64", "err", errSG, "path", p.string, "previousErr", err)
		}
	}
	return
}

func (p basePath) validateFloat64(v float64) (err error) {
	if p.Source != nil && false == p.Source.ContainsValFloat64(v) {
		jv, jErr := p.Source.ToJSON()
		err = errgo.Newf("The value '%.14f' cannot be found within the allowed Options():\n%s\nJSON Error: %s", v, jv, jErr)
	}
	return
}

func (p basePath) lookupBool(sg config.ScopedGetter) (v bool, err error) {

	var f *element.Field
	if f, err = p.field(sg); err != nil {
		return
	}

	v, err = cast.ToBoolE(f.Default)
	if err != nil {
		err = errgo.Mask(err)
		return
	}

	if val, errSG := sg.Bool(p.string); errSG == nil {
		v = val
	} else {
		// errSG is usually a key not found error, but that one is uninteresting
		if PkgLog.IsDebug() {
			PkgLog.Debug("model.path.lookupString.ScopedGetter.Bool", "err", errSG, "path", p.string, "previousErr", err)
		}
	}

	return
}
