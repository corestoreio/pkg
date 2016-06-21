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
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"

	"github.com/corestoreio/csfw/util/errors"
)

// Parse PEM encoded PKCS1 or PKCS8 private key. Provide optionally a password
// if key is encrypted.
func parseRSAPrivateKeyFromPEM(privateKey []byte, password ...[]byte) (*rsa.PrivateKey, error) {
	var prKeyPEM *pem.Block
	if prKeyPEM, _ = pem.Decode(privateKey); prKeyPEM == nil {
		return nil, errors.NewNotSupportedf(errKeyMustBePEMEncoded)
	}

	var block []byte
	if x509.IsEncryptedPEMBlock(prKeyPEM) {
		if len(password) != 1 || len(password[0]) == 0 {
			return nil, errors.NewEmptyf(errKeyMissingPassword)
		}
		var errPEM error
		if block, errPEM = x509.DecryptPEMBlock(prKeyPEM, password[0]); errPEM != nil {
			return nil, errors.NewNotValidf(errKeyDecryptPEMBlockFailed, errPEM)
		}
	} else {
		block = prKeyPEM.Bytes
	}

	var err error
	var parsedKey interface{}
	if parsedKey, err = x509.ParsePKCS1PrivateKey(block); err != nil {
		if parsedKey, err = x509.ParsePKCS8PrivateKey(block); err != nil {
			return nil, errors.NewNotValidf(errKeyParsePKCS8PrivateKeyFailed, err)
		}
	}

	var pkey *rsa.PrivateKey
	var ok bool
	if pkey, ok = parsedKey.(*rsa.PrivateKey); !ok {
		return nil, errors.NewNotValidf(errKeyNonRSAPrivateKey)
	}

	return pkey, nil
}

// Parse PEM encoded PKCS1 or PKCS8 public key
func parseRSAPublicKeyFromPEM(key []byte) (*rsa.PublicKey, error) {
	var err error

	// Parse PEM block
	var block *pem.Block
	if block, _ = pem.Decode(key); block == nil {
		return nil, errors.NewNotSupportedf(errKeyMustBePEMEncoded)
	}

	// Parse the key
	var parsedKey interface{}
	if parsedKey, err = x509.ParsePKIXPublicKey(block.Bytes); err != nil {
		if cert, err := x509.ParseCertificate(block.Bytes); err == nil {
			parsedKey = cert.PublicKey
		} else {
			return nil, errors.NewNotValidf(errKeyParseCertificateFailed, err)
		}
	}

	var pkey *rsa.PublicKey
	var ok bool
	if pkey, ok = parsedKey.(*rsa.PublicKey); !ok {
		return nil, errors.NewNotValidf(errKeyNonRSAPrivateKey)
	}

	return pkey, nil
}
