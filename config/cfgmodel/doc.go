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

// Package cfgmodel provides types to get/set values of a configuration.
//
// Package cfgmodel handles the scope permission checking, validation based
// on source models and default value handling based on element.Field type.
//
// In Mage world this would be called BackendModel.
//
// The signature of a getter function states in most cases:
//		Get(sg config.ScopedGetter) (v <T>)
// The Get() function signature may vary between the packages.
//
// The signature of the setter function states in most cases:
// 		Write(w config.Writer, v interface{}, s scope.Scope, id int64) error
// The Write() function signature differs within the types to mainly force the
// type safety. In other packages the Write() signature can be totally different.
//
// The responsibility of config.Writer adheres to the correct type conversion
// to the supported type of the underlying storage engine. E.g. for package
// config/storage/ccd it config.Writer converts all types to a byte slice.
//
//
// The global PackageConfiguration variable (type element.SectionSlice), which
// is present in each package, gets set to the cfgmodel.New* variables during init
// process. The element.Field will be extracted to allow scope checks and access
// to default values.
//
// The default value gets returned if the Get call to the store configuration
// value fails or a value is not set.
package cfgmodel
