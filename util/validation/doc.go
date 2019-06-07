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

// Package validation provides validation function for primitive data types.
//
// It does not and will not use any kind of reflection.

// review https://github.com/go-ozzo/ozzo-validation but with an API like uber-go/zap. Ozzo is a slow reflection soup.
// import "github.com/asaskevich/govalidator"
// TODO: review github.com/markbates/validate which has the nicest API ... not really
// https://github.com/gobuffalo/authrecipe/blob/master/models/user.go#53
// TODO: review https://github.com/RussellLuo/validating
//
// https://github.com/go-validator/validator uses struct tags which means hard coded validation and not changeable
// https://github.com/go-ozzo/ozzo-validation uses struct field based validation rules, changeable.
// https://github.com/go-playground/validator uses struct tags

package validation
