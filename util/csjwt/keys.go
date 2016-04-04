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

package csjwt

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/rsa"
	"io/ioutil"

	"github.com/juju/errors"
)

// PrivateKeyBits used when auto generating a private key
const PrivateKeyBits = 2048

// ErrHMACEmptyPassword whenever the password length is 0.
var ErrHMACEmptyPassword = errors.New("Empty passwords are forbidden")

// Keyfunc used by Parse methods, this callback function supplies
// the key for verification.  The function receives the parsed,
// but unverified Token.  This allows you to use propries in the
// Header of the token (such as `kid`) to identify which key to use.
type Keyfunc func(Token) (Key, error)

// Key defines a container for the HMAC password, RSA and ECDSA public and
// private keys. The Error fields gets filled out when loading/parsing the keys.
type Key struct {
	hmacPassword []byte
	ecdsaKeyPub  *ecdsa.PublicKey
	ecdsaKeyPriv *ecdsa.PrivateKey
	rsaKeyPub    *rsa.PublicKey
	rsaKeyPriv   *rsa.PrivateKey
	Error        error
}

const goStringTpl = `csjwt.Key{}`

// GoString protects keys and enforces privacy.
func (k Key) GoString() string {
	return goStringTpl
}

// String protects keys and enforces privacy.
func (k Key) String() string {
	return goStringTpl
}

// IsEmpty returns true when no field has been used in the Key struct.
// Error is excluded from the check
func (k Key) IsEmpty() bool {
	return k.hmacPassword == nil && k.ecdsaKeyPub == nil && k.ecdsaKeyPriv == nil && k.rsaKeyPub == nil && k.rsaKeyPriv == nil
}

// Algorithm returns the supported algorithm but not the bit size.
// Returns 0 on error, or one of the constants: ES, HS or RS.
func (k Key) Algorithm() (a Algorithm) {
	switch {
	case len(k.hmacPassword) > 0:
		a = HS
	case k.rsaKeyPriv != nil:
		a = RS // also matches RSA-PSS
	case k.ecdsaKeyPriv != nil:
		a = ES
	case k.rsaKeyPub != nil:
		a = RS // also matches RSA-PSS
	case k.ecdsaKeyPub != nil:
		a = ES
	}
	return a
}

// WithPassword uses the byte slice as the password for the HMAC-SHA signing method.
func WithPassword(password []byte) Key {
	var err error
	if len(password) == 0 {
		err = ErrHMACEmptyPassword
	}
	return Key{
		hmacPassword: password,
		Error:        err,
	}
}

const randomPasswordLenght = 32

// WithPasswordRandom creates cryptographically secure random password which you
// cannot obtain. Whenever you restart your app with a random password, all
// HMAC-SHA tokens get invalided.
func WithPasswordRandom() Key {
	pw := make([]byte, randomPasswordLenght)
	_, err := rand.Read(pw)
	return Key{
		hmacPassword: pw,
		Error:        err,
	}
}

// WithPasswordFromFile loads the content of a file and uses that content as
// the password for the HMAC-SHA signing method.
func WithPasswordFromFile(pathToFile string) Key {
	var k Key
	k.hmacPassword, k.Error = ioutil.ReadFile(pathToFile)
	if k.Error != nil {
		k.Error = k.Error
	}
	return k
}

// WithRSAPublicKeyFromPEM parses PEM encoded PKCS1 or PKCS8 public key
func WithRSAPublicKeyFromPEM(publicKey []byte) (k Key) {
	k.rsaKeyPub, k.Error = parseRSAPublicKeyFromPEM(publicKey)
	return
}

// WithRSAPublicKeyFromFile parses PEM encoded PKCS1 or PKCS8 public key found
// in a file.
func WithRSAPublicKeyFromFile(pathToFile string) (k Key) {
	pk, err := ioutil.ReadFile(pathToFile)
	if err != nil {
		k.Error = errors.Errorf("%s with file %s", err, pathToFile)
		return k
	}
	return WithRSAPublicKeyFromPEM(pk)
}

// WithRSAPublicKey sets the public key
func WithRSAPublicKey(publicKey *rsa.PublicKey) (k Key) {
	k.rsaKeyPub = publicKey
	return
}

// WithRSAPrivateKeyFromPEM parses PEM encoded PKCS1 or PKCS8 private key.
// Password as second argument is only required when the
// private key is encrypted. Public key will be derived from the private key.
func WithRSAPrivateKeyFromPEM(privateKey []byte, password ...[]byte) (k Key) {
	k.rsaKeyPriv, k.Error = parseRSAPrivateKeyFromPEM(privateKey, password...)
	if k.rsaKeyPriv != nil {
		k.rsaKeyPub = &k.rsaKeyPriv.PublicKey
	}
	return
}

// WithRSAPrivateKeyFromFile parses PEM encoded PKCS1 or PKCS8 private key
// found in a file.
// Password as second argument is only required when the
// private key is encrypted. Public key will be derived from the private key.
func WithRSAPrivateKeyFromFile(pathToFile string, password ...[]byte) (k Key) {
	pk, err := ioutil.ReadFile(pathToFile)
	if err != nil {
		k.Error = err
		return k
	}
	return WithRSAPrivateKeyFromPEM(pk, password...)
}

// WithRSAPrivateKey sets the private key.
// Public key will be derived from the private key.
func WithRSAPrivateKey(privateKey *rsa.PrivateKey) (k Key) {
	k.rsaKeyPriv = privateKey
	k.rsaKeyPub = &privateKey.PublicKey
	return
}

// WithRSAGenerated creates an in-memory private key to be used for signing
// and verifying. Bit size see constant: PrivateKeyBits
// Public key will be derived from the private key.
func WithRSAGenerated() (k Key) {
	pk, err := rsa.GenerateKey(rand.Reader, PrivateKeyBits)
	if err != nil {
		k.Error = err
		return
	}
	k.rsaKeyPriv = pk
	k.rsaKeyPub = &pk.PublicKey
	return
}

// WithECPublicKeyFromPEM parses PEM encoded Elliptic Curve Public Key Structure
func WithECPublicKeyFromPEM(publicKey []byte) (k Key) {
	k.ecdsaKeyPub, k.Error = parseECPublicKeyFromPEM(publicKey)
	return
}

// WithECPublicKeyFromFile parses a file PEM encoded Elliptic Curve Public Key Structure
func WithECPublicKeyFromFile(pathToFile string) (k Key) {
	pk, err := ioutil.ReadFile(pathToFile)
	if err != nil {
		k.Error = err
		return k
	}
	k.ecdsaKeyPub, k.Error = parseECPublicKeyFromPEM(pk)
	return
}

// WithECPublicKey sets the ECDSA public key
func WithECPublicKey(publicKey *ecdsa.PublicKey) (k Key) {
	k.ecdsaKeyPub = publicKey
	return
}

// WithECPrivateKeyFromPEM parses PEM encoded Elliptic Curve Private Key Structure.
// Public key will be derived from the private key.
func WithECPrivateKeyFromPEM(privateKey []byte) (k Key) {
	k.ecdsaKeyPriv, k.Error = parseECPrivateKeyFromPEM(privateKey)
	if k.ecdsaKeyPriv != nil {
		k.ecdsaKeyPub = &k.ecdsaKeyPriv.PublicKey
	}
	return
}

// WithECPrivateKeyFromFile parses file PEM encoded Elliptic Curve Private Key Structure.
// Public key will be derived from the private key.
func WithECPrivateKeyFromFile(pathToFile string) (k Key) {
	pk, err := ioutil.ReadFile(pathToFile)
	if err != nil {
		k.Error = err
		return k
	}
	return WithECPrivateKeyFromPEM(pk)
}

// WithECPrivateKey sets the ECDSA private key.
// Public key will be derived from the private key.
func WithECPrivateKey(privateKey *ecdsa.PrivateKey) (k Key) {
	k.ecdsaKeyPriv = privateKey
	k.ecdsaKeyPub = &privateKey.PublicKey
	return
}
