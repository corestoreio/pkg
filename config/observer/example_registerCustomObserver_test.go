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

// +build csall json

package observer

// import (
// 	"github.com/corestoreio/pkg/config"
// 	"github.com/corestoreio/pkg/config/validation"
// )

// type ValidFileExtension struct {
// 	Extensions []string
// }
//
// func (fe *ValidFileExtension) UnmarshalJSON(data []byte) error {
// 	return nil
// }
//
// func (fe ValidFileExtension) Observe(p config.Path, rawData []byte, found bool) (newRawData []byte, err error) {
// 	return rawData, nil
// }

func ExampleRegisterCustomObserver() {

	// validation.RegisterCustomObserver("image_file", &ValidFileExtension{})
}
