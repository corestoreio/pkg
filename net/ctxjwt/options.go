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

package ctxjwt

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"io"
	"io/ioutil"
	"os"

	"github.com/corestoreio/csfw/config"
	"github.com/dgrijalva/jwt-go"
	"github.com/pborman/uuid"
)

// PathJWTPassword defines the path where the password has been stored.
const PathJWTPassword = "corestore/jwt/password"

// ErrPrivateKeyNotFound will be returned when the PK cannot be read from the Reader
var ErrPrivateKeyNotFound = errors.New("Private Key from io.Reader no found")

// ErrPrivateKeyNoPassword will be returned when the PK is encrypted but you
// forgot to provide a password.
var ErrPrivateKeyNoPassword = errors.New("Private Key is encrypted but password was not set")

// PrivateKeyBits used when auto generating a private key
const PrivateKeyBits = 4096

// @todo add more WithKeyFrom...()

// Option can be used as an argument in NewService to configure a token service.
type Option func(a *Service)

// WithBlacklist sets a new black list service.
func WithBlacklist(blacklist Blacklister) Option {
	return func(s *Service) {
		s.Blacklist = blacklist
	}
}

// WithPasswordFromConfig retrieves the password from the configuration with path
// as defined in constant PathJWTPassword
func WithPasswordFromConfig(cr config.Getter) Option {
	pw, err := cr.String(config.Path(PathJWTPassword))
	if config.NotKeyNotFoundError(err) {
		pw = string(uuid.NewRandom())
	}
	return WithPassword([]byte(pw))
}

// WithPassword sets the HMAC 256 bit signing method with a password. Useful to use Magento encryption key.
func WithPassword(key []byte) Option {
	return func(s *Service) {
		s.lastError = nil
		s.hasKey = true
		s.SigningMethod = jwt.SigningMethodHS256
		s.password = key
	}
}

// WithECDSAFromFile loads the ECDSA key from a file @todo
func WithECDSAFromFile(fileName string, password ...[]byte) Option {
	fpk, err := os.Open(fileName)
	if err != nil {
		return func(s *Service) {
			if PkgLog.IsDebug() {
				PkgLog.Debug("ctxjwt.WithECDSAFromFile.os.Open", "err", err, "file", fileName)
			}
			s.lastError = err
		}
	}
	return WithECDSA(fpk, password...)

}

// WithECDSA @todo
// Default Signing bits 256.
func WithECDSA(privateKey io.Reader, password ...[]byte) Option {
	if cl, ok := privateKey.(io.Closer); ok {
		defer func() {
			if err := cl.Close(); err != nil { // close file
				PkgLog.Debug("ctxjwt.ECDSAKey.ioCloser", "err", err)
			}
		}()
	}

	return func(s *Service) {
		s.hasKey = false // set to true if fully implemented
		s.lastError = errors.New("@todo implement")
		s.SigningMethod = jwt.SigningMethodES256
		s.ecdsapk = nil
	}
}

// WithRSAFromFile reads an RSA private key from a file and applies it as an option
// to the AuthManager. Password as second argument is only required when the
// private key is encrypted. Public key will be derived from the private key.
func WithRSAFromFile(fileName string, password ...[]byte) Option {
	fpk, err := os.Open(fileName)
	if err != nil {
		return func(s *Service) {
			if PkgLog.IsDebug() {
				PkgLog.Debug("ctxjwt.WithRSAFromFile.os.Open", "err", err, "file", fileName)
			}
			s.lastError = err
		}
	}
	return WithRSA(fpk, password...)
}

// WithRSA reads PEM byte data and decodes it and parses the private key.
// Applies the private and the public key to the AuthManager. Password as second
// argument is only required when the private key is encrypted.
// Checks for io.Close and closes the resource. Public key will be derived from
// the private key. Default Signing bits 256.
func WithRSA(privateKey io.Reader, password ...[]byte) Option {
	if cl, ok := privateKey.(io.Closer); ok {
		defer func() {
			if err := cl.Close(); err != nil { // close file
				PkgLog.Debug("ctxjwt.RSAKey.ioCloser", "err", err)
			}
		}()
	}
	prKeyData, errRA := ioutil.ReadAll(privateKey)
	if errRA != nil {
		return func(a *Service) {
			if PkgLog.IsDebug() {
				PkgLog.Debug("ctxjwt.WithRSA.ioutil.ReadAll", "err", errRA, "privateKey", privateKey)
			}
			a.lastError = errRA
		}
	}
	var prKeyPEM *pem.Block
	if prKeyPEM, _ = pem.Decode(prKeyData); prKeyPEM == nil {
		return func(s *Service) {
			if PkgLog.IsDebug() {
				PkgLog.Debug("ctxjwt.WithRSA.pem.Decode", "err", ErrPrivateKeyNotFound, "prKeyData", prKeyData)
			}
			s.lastError = ErrPrivateKeyNotFound
		}
	}

	var rsaPrivateKey *rsa.PrivateKey
	var err error
	if x509.IsEncryptedPEMBlock(prKeyPEM) {
		if len(password) != 1 || len(password[0]) == 0 {
			return func(s *Service) {
				if PkgLog.IsDebug() {
					PkgLog.Debug("ctxjwt.WithRSA.IsEncryptedPEMBlock", "err", ErrPrivateKeyNoPassword)
				}
				s.lastError = ErrPrivateKeyNoPassword
			}
		}
		var dd []byte
		var errPEM error
		if dd, errPEM = x509.DecryptPEMBlock(prKeyPEM, password[0]); errPEM != nil {
			return func(s *Service) {
				if PkgLog.IsDebug() {
					PkgLog.Debug("ctxjwt.WithRSA.DecryptPEMBlock", "err", errPEM)
				}
				s.lastError = errPEM
			}
		}
		rsaPrivateKey, err = x509.ParsePKCS1PrivateKey(dd)
	} else {
		rsaPrivateKey, err = x509.ParsePKCS1PrivateKey(prKeyPEM.Bytes)
	}

	return func(s *Service) {
		s.SigningMethod = jwt.SigningMethodRS256
		s.rsapk = rsaPrivateKey
		s.hasKey = true
		if err != nil {
			if PkgLog.IsDebug() {
				PkgLog.Debug("ctxjwt.WithRSA.ParsePKCS1PrivateKey", "err", err)
			}
			s.lastError = err
		}
	}
}

// WithRSAGenerator creates an in-memory RSA key without persisting it.
// This function may run around ~3secs.
func WithRSAGenerator() Option {
	pk, err := rsa.GenerateKey(rand.Reader, PrivateKeyBits)
	return func(s *Service) {
		if pk != nil {
			s.rsapk = pk
			s.hasKey = true
			s.SigningMethod = jwt.SigningMethodRS256
		}
		if err != nil {
			if PkgLog.IsDebug() {
				PkgLog.Debug("ctxjwt.WithRSAGenerator.GenerateKey", "err", err)
			}
			s.lastError = err
		}
	}
}
