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

package proto

//go:generate protoc --gogo_out=Mgoogle/protobuf/timestamp.proto=github.com/gogo/protobuf/types:. --proto_path=../../../../../:../../../../../github.com/gogo/protobuf/protobuf/:. *.proto

import (
	"github.com/corestoreio/pkg/config"
	"github.com/corestoreio/pkg/config/validation"
)

// Validator defines the data retrieved from the outside as JSON to add a new
// validator for a specific route and event.
type Validator struct {
	validation.Validator
}

// Validators a list of Validator types.
type Validators struct {
	Collection []*Validator
}

func RegisterObservers(or config.ObserverRegisterer, protoData []byte) error {

	// vs := make(Validators, 0, 5)
	// if err := vs.Unmarshal(protoData); err != nil {
	// 	return errors.WithStack(err)
	// }
	//
	// for _, v := range vs {
	// 	event, route, o, err := v.NewFromJSON()
	// 	if err != nil {
	// 		return errors.Wrapf(err, "[config/validation] Data: %#v", v)
	// 	}
	// 	if err := or.RegisterObserver(event, route, o); err != nil {
	// 		return errors.WithStack(err)
	// 	}
	// }
	return nil
}
