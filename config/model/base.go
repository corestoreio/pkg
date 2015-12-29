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

package model

import (
	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/config/valuelabel"
	"github.com/corestoreio/csfw/store/scope"
	"github.com/corestoreio/csfw/util/cast"
	"github.com/juju/errgo"
)

// SourceModeller defines how to retrieve all option values. Mostly used for frontend output.
type SourceModeller interface {
	Options() valuelabel.Slice
}

var _ SourceModeller = (*basePath)(nil)

// basePath defines the path in the "core_config_data" table like a/b/c. All other
// types in this package inherits from this path type.
type basePath struct {
	string // contains the path

	// ValueLabel are all available options aka SourceModel in Magento slang.
	// This slice is also used for validation to get and write the correct values.
	// Validation gets triggered only when the slice has been set.
	// The Options() function will be used to access this slice.
	ValueLabel valuelabel.Slice
}

// NewPath creates a new basePath type optional value label Pair
func NewPath(path string, vlPairs ...valuelabel.Pair) basePath {
	return basePath{
		string:     path,
		ValueLabel: valuelabel.Slice(vlPairs),
	}
}

// Write writes a value v to the config.Writer without checking if the value
// has changed.
func (p basePath) Write(w config.Writer, v interface{}, s scope.Scope, id int64) error {
	return w.Write(config.Path(p.string), config.Value(v), config.Scope(s, id))
}

// String returns the path
func (p basePath) String() string {
	return p.string
}

// InScope checks if a field from a path is allowed for current ScopedGetter.
func (p basePath) InScope(f *config.Field, sg config.ScopedGetter) (err error) {
	if s, _ := sg.Scope(); false == f.Scope.Has(s) {
		err = errgo.Newf("Scope permission insufficient: Have '%s'; Want '%s'", s, f.Scope)
	}
	return
}

// Options returns a source model for all available options for a configuration
// value.
//
// Usually this function gets customized in a sub-type. Customization
// can have different arguments, etc but must always call this function to set
// valuelabel slice.
func (p basePath) Options() valuelabel.Slice {
	return p.ValueLabel
}

// FQPathInt64 generates a fully qualified configuration path.
// Example: general/country/allow would transform with StrScope scope.StrStores
// and storeID e.g. 4 into: stores/4/general/country/allow
func (p basePath) FQPathInt64(strScope scope.StrScope, scopeID int64) string {
	return strScope.FQPathInt64(scopeID, p.string)
}

// field searches for the field in a SectionSlice and checks if the scope in
// ScopedGetter is sufficient.
func (p basePath) field(pkgCfg config.SectionSlice, sg config.ScopedGetter) (field *config.Field, err error) {
	if field, err = pkgCfg.FindFieldByPath(p.string); err != nil {
		return
	}
	err = p.InScope(field, sg)
	return
}

// lookupString searches the default value in config.SectionSlice and overwrites
// it with a value from ScopedGetter if ScopedGetter is not empty.
// validator can be nil which triggers the default validation method.
func (p basePath) lookupString(pkgCfg config.SectionSlice, sg config.ScopedGetter) (v string, err error) {

	var f *config.Field
	if f, err = p.field(pkgCfg, sg); err != nil {
		return
	}

	v, err = cast.ToStringE(f.Default)
	if err != nil {
		err = errgo.Mask(err)
		return
	}

	if val, errSG := sg.String(p.string); errSG == nil {
		v = val
	} else {
		// errSG is usually a key not found error, but that one is uninteresting
		if PkgLog.IsDebug() {
			PkgLog.Debug("model.basePath.lookupString.ScopedGetter.String", "err", errSG, "path", p.string, "previousErr", err)
		}
	}
	return
}

func (p basePath) validateString(v string) (err error) {
	if len(p.ValueLabel) > 0 && false == p.ValueLabel.ContainsValString(v) {
		jv, jErr := p.ValueLabel.ToJSON()
		err = errgo.Newf("The value '%s' cannot be found within the allowed Options():\n%s\nJSON Error: %s", v, jv, jErr)
	}
	return
}

func (p basePath) lookupInt(pkgCfg config.SectionSlice, sg config.ScopedGetter) (v int, err error) {

	var f *config.Field
	if f, err = p.field(pkgCfg, sg); err != nil {
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
	if len(p.ValueLabel) > 0 && false == p.ValueLabel.ContainsValInt(v) {
		jv, jErr := p.ValueLabel.ToJSON()
		err = errgo.Newf("The value '%d' cannot be found within the allowed Options():\n%s\nJSON Error: %s", v, jv, jErr)
	}
	return
}

func (p basePath) lookupFloat64(pkgCfg config.SectionSlice, sg config.ScopedGetter) (v float64, err error) {

	var f *config.Field
	if f, err = p.field(pkgCfg, sg); err != nil {
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
	if len(p.ValueLabel) > 0 && false == p.ValueLabel.ContainsValFloat64(v) {
		jv, jErr := p.ValueLabel.ToJSON()
		err = errgo.Newf("The value '%.14f' cannot be found within the allowed Options():\n%s\nJSON Error: %s", v, jv, jErr)
	}
	return
}

func (p basePath) lookupBool(pkgCfg config.SectionSlice, sg config.ScopedGetter) (v bool, err error) {

	var f *config.Field
	if f, err = p.field(pkgCfg, sg); err != nil {
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
