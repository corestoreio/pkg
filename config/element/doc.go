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
// Thw following diagram shows the tree structure:
//    +---------+
//    | Section |
//    +---------+
//    |
//    |   +-------+
//    +-->+ Group |
//    |   +-------+
//    |   |
//    |   |    +--------+
//    |   +--->+ Field  |
//    |   |    +--------+
//    |   |
//    |   |    +--------+
//    |   +--->+ Field  |
//    |   |    +--------+
//    |   |
//    |   |    +--------+
//    |   +--->+ Field  |
//    |        +--------+
//    |
//    |   +-------+
//    +-->+ Group |
//    |   +-------+
//    |   |
//    |   |    +--------+
//    |   +--->+ Field  |
//    |   |    +--------+
//    |   |
//    |   |    +--------+
//    |   +--->+ Field  |
//    |   |    +--------+
//    |   |
//    |   |    +--------+
//    |   +--->+ Field  |
//    |        +--------+
//    |
//    http://asciiflow.com/
// The three elements Section, Group and Field represents front-end configuration
// fields and more important default values and their permissions. A permission
// is of type scope.Perm and defines which of elements in which scope can be shown.
//
// Those three elements represents the tree in function NewConfigStructure() which
// can be found in any package.
//
// Unclear: Your app which includes the csfw must merge all "PackageConfiguration"s into a single slice.
// You should submit all default values (interface config.Sectioner) to the config.Service.ApplyDefaults()
// function.
//
// The JSON encoding of the three elements Section, Group and Field are intended to use
// on the backend REST API and for debugging and testing. Only used in non performance critical parts.
package element
