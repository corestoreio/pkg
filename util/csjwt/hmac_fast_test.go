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

package csjwt_test

import (
	"bytes"
	"testing"

	"github.com/corestoreio/csfw/util/csjwt"
)

func TestHMACVerifyFast(t *testing.T) {
	for _, data := range hmacTestData {
		signing, signature, err := csjwt.SplitForVerify(data.tokenString)
		if err != nil {
			t.Fatal(err, "\n", string(data.tokenString))
		}
		alg := csjwt.ToAlgorithm(data.alg) * HMACFastSuffix
		method := csjwt.GetSigningMethod(alg.String())
		err = method.Verify(signing, signature, csjwt.Key{})
		if data.valid && err != nil {
			t.Errorf("[%v] Method %s Error while verifying key: %v", data.name, data.alg, err)
		}
		if !data.valid && err == nil {
			t.Errorf("[%v] Invalid key passed validation", data.name)
		}
	}
}

func TestHMACSignFast(t *testing.T) {
	for _, data := range hmacTestData {
		if data.valid {
			signing, signature, err := csjwt.SplitForVerify(data.tokenString)
			if err != nil {
				t.Fatal(err, "\n", string(data.tokenString))
			}

			alg := csjwt.ToAlgorithm(data.alg) * HMACFastSuffix
			method := csjwt.GetSigningMethod(alg.String())
			sig, err := method.Sign(signing, csjwt.Key{})
			if err != nil {
				t.Errorf("[%v] Error signing token: %v", data.name, err)
			}
			if !bytes.Equal(sig, signature) {
				t.Errorf("[%v] Incorrect signature.\nwas:\n%v\nexpecting:\n%v", data.name, string(sig), string(signature))
			}
		}
	}
}
