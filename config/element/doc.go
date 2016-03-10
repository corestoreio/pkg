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

// Package element represents Magento system.xml configuration elements.
//
// The three elements Section, Group and Field represents front-end configuration fields and more important
// default values and their permissions. They do not define how to handle the source and backend model
// in a Magento sense. Source models are use to load values for displaying in a e.g. HTML select field
// or also known as option values. Backend models know how to save and load a cfgpath.Path
//
// Those three elements represents the PackageConfiguration variable which can be found in any package.
//
// Your app which includes the csfw must merge all "PackageConfiguration"s into a single slice.
// You should submit all default values (interface config.Sectioner) to the config.Service.ApplyDefaults()
// function.
//
// The models included in PackageConfiguration will be later used when handling the values
// for each configuration field.
//
// The JSON enconding of the three elements Section, Group and Field are intended to use
// on the backend REST API and for debugging and testing. Only used in non performance critical parts.
package element
