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

package backendsigned_test

import (
	"github.com/corestoreio/cspkg/config/cfgmodel"
	"github.com/corestoreio/cspkg/net/signed/backendsigned"
)

// backend overall backend models for all tests
var backend *backendsigned.Configuration

var _ cfgmodel.Encrypter = (*noopCrypt)(nil)
var _ cfgmodel.Decrypter = (*noopCrypt)(nil)

type noopCrypt struct{}

func (noopCrypt) Encrypt(s []byte) ([]byte, error) {
	return s, nil
}

func (noopCrypt) Decrypt(s []byte) ([]byte, error) {
	return s, nil
}

// this would belong into the test suit setup
func init() {
	cfgStruct, err := backendsigned.NewConfigStructure()
	if err != nil {
		panic(err)
	}
	backend = backendsigned.New(cfgStruct)

	backend.Key.Encrypter = noopCrypt{}
	backend.Key.Decrypter = noopCrypt{}
}
