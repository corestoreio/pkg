// Copyright 2015-present, Cyrill @ Schumacher.fm and the CoreStore contributors
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

package null

// AllTypes contains value objects of all types in this package. The array might
// change in size in the future. The types are useful for functions like:
// gopkg.in/go-playground/validator.v9/Validate.RegisterCustomTypeFunc or
// gob.Register.
var AllTypes = [...]interface{}{Bool{}, Decimal{}, Float64{}, Int8{}, Int16{}, Int32{}, Int64{}, String{}, Time{}, Uint8{}, Uint16{}, Uint32{}, Uint64{}}
