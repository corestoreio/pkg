// Copyright 2015, Cyrill @ Schumacher.fm and the CoreStore contributors
//
// Licensed under the Apache License, Version 2.0 (the "License")
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

package crypto

// @see lib/internal/Magento/Framework/Encryption

type Encrypter interface {
	// Encrypt a string
	Encrypt(data []byte) []byte

	// Decrypt a string
	Decrypt(data []byte) []byte

	// Return crypt model, instantiate if it is empty
	// @return \Magento\Framework\Encryption\Crypt
	ValidateKey(key string)
}

type Hasher interface {

	// Generate a [salted] hash.
	// $salt can be:
	// false - salt is not used
	// true - random salt of the default length will be generated
	// integer - random salt of specified length will be generated
	// string - actual salt value to be used
	GenerateHash(password string, salt bool) string

	// Hash a string
	Hash(data string) string

	// Validate hash against hashing method (with or without salt)
	ValidateHash(password, hash string) error
}
