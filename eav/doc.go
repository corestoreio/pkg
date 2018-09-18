// Copyright 2015-2017, Cyrill @ Schumacher.fm and the CoreStore contributors
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

// Package eav contains the logic for the Entity-Attribute-Value pattern (WIP).
//
// To use this library with additional columns in the EAV tables you must run
// from the tools folder first `tableToStruct` and then build the program
// `eavToStruct` and run it.
//
// Definition of attribute backend, source, and frontend models:
//
// - Backend: Provides hooks before and after save, load, and delete operations
// with an attribute value.
// - Source: Provides option values and labels for select and multi-select attributes.
// - Frontend: Prepares an attribute value for rendering on the storefront and admin backend.
//
// Backend models can be an alternative to an observer; for example, when you
// have to do something that depends on an attribute value when an entity is
// saved.
//
// TODO(CSC): idea to import data quickly: see Entity-Attribute-Value_(EAV)_The_Antipattern_Too_Great_to_Give_Up_-__Andy_Novick_2016-03-19.pdf
// Break it down to single partition operations
// • SQLCLR proc breaks the file by attribute_id
// • SEND attribute_id’s data to a Service Broker QUEUE
// • Each task is working on ONE attribute_id
// – That’s one HOBT / Partition
// • Run 1-2 tasks per core
// CSC: in our case run a pool goroutines to work on the attribute IDs or for
// each attribute_id a dedicated goroutine (maybe code generated)
package eav
