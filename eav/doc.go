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

/*
Package eav contains the logic for the Entity-Attribute-Value model based on the Magento database schema.

To use this library with additional columns in the EAV tables you must run from the
tools folder first `tableToStruct` and then build the program `eavToStruct` and run it.

TODO EAV Models

For what are attribute backend, source, and frontend models for:

- Backend: Provides hooks before and after save, load, and delete operations with an attribute value.
- Source: Provides option values and labels for select and multi-select attributes.
- Frontend: Prepares an attribute value for rendering on the storefront.

Backend models can be an alternative to an observer; for example, when you have to do
something that depends on an attribute value when an entity is saved.
*/
package eav
