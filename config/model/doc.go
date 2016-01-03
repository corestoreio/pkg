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

// Package model provides types for getting and setting values of configuration
// fields/values with validation and their default value handling.
//
// In Mage world this would be called BackendModel.
//
// The default value gets returned if the Get call to the store configuration
// value fails or value is not set.
//
// The signature of a getter function states in most cases:
//		Get(sg config.ScopedGetter) (v <T>)
// The global PackageConfiguration variable which is present in each
// package gets set to the Path* variables during init process and then
// shall not change. PackageConfiguration knows the default value of a
// configuration path.
// sg config.ScopedGetter is the current config.Getter but bounded to a
// scope. If sg finds not a value then the default value gets returned.
//
// The Get() function signature may vary between the packages.
//
// The signature of the setter function states in most cases:
// 		Write(w config.Writer, v interface{}, s scope.Scope, id int64) error
// The interface v gets in the parent type replaced by the correct type and
// this type gets converted most times to a string or int or float.
// Sometimes the Write() function signature can differ in packages.
//
// This package stays pointer free because these types will be more often
// used as global variables, cough cough, through different packages.
// With non-pointers we reduce the pressure on the GC.
package model
