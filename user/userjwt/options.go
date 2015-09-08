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

package userjwt

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"os"

	"io"

	"io/ioutil"

	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/utils/log"
	"github.com/dgrijalva/jwt-go"
	"github.com/juju/errgo"
)

// PathJWTPassword defines the path where the password has been stored.
const PathJWTPassword = "corestore/userjwt/password"

// @todo add more KeyFrom...()

// OptionFunc can be used as an argument in NewUser to configure a user.
type OptionFunc func(a *AuthManager)

// SetPasswordFromConfig retrieves the password from the configuration with path
// as defined in constant PathJWTPassword
func SetPasswordFromConfig(cr config.Reader) OptionFunc {
	pw := cr.GetString(config.Path(PathJWTPassword))
	return SetPassword([]byte(pw))
}

// SetPassword sets the HMAC 256 bit signing method with a password. Useful to use Magento encryption key.
func SetPassword(key []byte) OptionFunc {
	return func(a *AuthManager) {
		a.lastError = nil
		a.hasKey = true
		a.SigningMethod = jwt.SigningMethodHS256
		a.password = key
	}
}

// SetECDSAFromFile @todo
func SetECDSAFromFile(privateKey string, password ...[]byte) OptionFunc {
	fpk, err := os.Open(privateKey)
	if err != nil {
		return func(a *AuthManager) {
			a.lastError = errgo.Mask(err)
		}
	}
	return SetECDSA(fpk, password...)

}

// SetECDSA @todo
// Default Signing bits 256.
func SetECDSA(privateKey io.Reader, password ...[]byte) OptionFunc {
	if cl, ok := privateKey.(io.Closer); ok {
		defer func() {
			if err := cl.Close(); err != nil { // close file
				log.Error("userjwt.ECDSAKey.ioCloser", "err", err)
			}
		}()
	}

	// @todo implement

	return func(a *AuthManager) {
		a.hasKey = false // set to true if fully implemented
		a.lastError = errgo.New("@todo implement")
		a.SigningMethod = jwt.SigningMethodES256
		a.ecdsapk = nil
	}
}

// SetRSAFromFile reads an RSA private key from a file and applies it as an option
// to the AuthManager. Password as second argument is only required when the
// private key is encrypted. Public key will be derived from the private key.
func SetRSAFromFile(privateKey string, password ...[]byte) OptionFunc {
	fpk, err := os.Open(privateKey)
	if err != nil {
		return func(a *AuthManager) {
			a.lastError = errgo.Mask(err)
		}
	}
	return SetRSA(fpk, password...)
}

// SetRSA reads PEM byte data and decodes it and parses the private key.
// Applies the private and the public key to the AuthManager. Password as second
// argument is only required when the private key is encrypted.
// Checks for io.Close and closes the resource. Public key will be derived from
// the private key. Default Signing bits 256.
func SetRSA(privateKey io.Reader, password ...[]byte) OptionFunc {
	if cl, ok := privateKey.(io.Closer); ok {
		defer func() {
			if err := cl.Close(); err != nil { // close file
				log.Error("userjwt.RSAKey.ioCloser", "err", err)
			}
		}()
	}
	prKeyData, errRA := ioutil.ReadAll(privateKey)
	if errRA != nil {
		return func(a *AuthManager) {
			a.lastError = errgo.Mask(errRA)
		}
	}
	var prKeyPEM *pem.Block
	if prKeyPEM, _ = pem.Decode(prKeyData); prKeyPEM == nil {
		return func(a *AuthManager) {
			a.lastError = errgo.New("Private Key from io.Reader no found")
		}
	}

	var rsaPrivateKey *rsa.PrivateKey
	var err error
	if x509.IsEncryptedPEMBlock(prKeyPEM) {
		if len(password) != 1 || len(password[0]) == 0 {
			return func(a *AuthManager) {
				a.lastError = errgo.New("Private Key is encrypted but password was not set")
			}
		}
		var dd []byte
		var errPEM error
		if dd, errPEM = x509.DecryptPEMBlock(prKeyPEM, password[0]); errPEM != nil {
			return func(a *AuthManager) {
				a.lastError = errgo.Newf("Private Key decryption failed: %s", errPEM.Error())
			}
		}
		rsaPrivateKey, err = x509.ParsePKCS1PrivateKey(dd)
	} else {
		rsaPrivateKey, err = x509.ParsePKCS1PrivateKey(prKeyPEM.Bytes)
	}

	return func(a *AuthManager) {
		a.SigningMethod = jwt.SigningMethodRS256
		a.rsapk = rsaPrivateKey
		a.hasKey = true
		a.lastError = errgo.Mask(err)
	}
}

// SetRSAGenerator creates an in-memory RSA key without persisting it.
// This function may run around ~3secs.
func SetRSAGenerator() OptionFunc {
	pk, err := rsa.GenerateKey(rand.Reader, PrivateKeyBits)
	return func(a *AuthManager) {
		if pk != nil {
			a.rsapk = pk
			a.hasKey = true
			a.SigningMethod = jwt.SigningMethodRS256
		}
		a.lastError = errgo.Mask(err)
	}
}
