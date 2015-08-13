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

	"github.com/corestoreio/csfw/utils/log"
	"github.com/juju/errgo"
)

// @todo add more KeyFrom...()

// PrivateKeyBits used when auto generating a private key
const PrivateKeyBits = 4096

// OptionFunc applies options to the AuthManager
type OptionFunc func(a *AuthManager)

// KeyFromFile reads an RSA private key from a file and applies it as an option
// to the AuthManager. Password as second argument is only required when the
// private key is encrypted. Public key will be derived from the private key.
func KeyFromFile(privateKey string, password ...[]byte) OptionFunc {
	fpk, err := os.Open(privateKey)
	if err != nil {
		return func(a *AuthManager) {
			a.lastError = errgo.Mask(err)
		}
	}
	return Key(fpk, password...)
}

// Key reads PEM byte data and decodes it and parses the private key.
// Applies the private and the public key to the AuthManager. Password as second
// argument is only required when the private key is encrypted.
// Checks for io.Close and closes the resource. Public key will be derived from
// the private key.
func Key(privateKey io.Reader, password ...[]byte) OptionFunc {
	if cl, ok := privateKey.(io.Closer); ok {
		defer func() {
			if err := cl.Close(); err != nil { // close file
				log.Error("userjwt.Key.ioCloser", "err", err)
			}
		}()
	}
	prKeyData, errRA := ioutil.ReadAll(privateKey)
	if errRA != nil {
		return func(a *AuthManager) {
			a.lastError = errgo.Mask(errRA)
		}
	}

	prKeyPEM, _ := pem.Decode(prKeyData)
	if prKeyPEM == nil {
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
		if dd, err := x509.DecryptPEMBlock(prKeyPEM, password[0]); err != nil {
			return func(a *AuthManager) {
				a.lastError = errgo.Newf("Private Key decryption failed: %s", err.Error())
			}
		}
		rsaPrivateKey, err = x509.ParsePKCS1PrivateKey(dd)
	} else {
		rsaPrivateKey, err = x509.ParsePKCS1PrivateKey(prKeyPEM.Bytes)
	}

	return func(a *AuthManager) {
		a.privateKey = rsaPrivateKey
		if rsaPrivateKey != nil {
			a.publicKey = &(rsaPrivateKey.PublicKey)
		}
		a.lastError = errgo.Mask(err)
	}
}

func generatePrivateKey(a *AuthManager) {
	pk, err := rsa.GenerateKey(rand.Reader, PrivateKeyBits)

	if pk != nil {
		a.privateKey = pk
		a.publicKey = &(pk.PublicKey)
	}
	a.lastError = errgo.Mask(err)

	//	privateKeyDer := x509.MarshalPKCS1PrivateKey(privateKey)
	//	privateKeyBlock := pem.Block{
	//		Type:    "RSA PRIVATE KEY",
	//		Headers: nil,
	//		Bytes:   privateKeyDer,
	//	}
	//	privateKeyPem := string(pem.EncodeToMemory(&privateKeyBlock))
	//
	//	publicKey := privateKey.PublicKey
	//	publicKeyDer, err := x509.MarshalPKIXPublicKey(&publicKey)
	//	if err != nil {
	//		return nil, err
	//	}
	//
	//	publicKeyBlock := pem.Block{
	//		Type:    "PUBLIC KEY",
	//		Headers: nil,
	//		Bytes:   publicKeyDer,
	//	}
	//	publicKeyPem := string(pem.EncodeToMemory(&publicKeyBlock))
	//
	//	fmt.Println(privateKeyPem)
	//	fmt.Println(publicKeyPem)
}
