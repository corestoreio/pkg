package csjwt

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"

	"github.com/corestoreio/csfw/util/errors"
)

// Parse PEM encoded PKCS1 or PKCS8 private key. Provide optionally a password if
// key is encrypted.
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
			return nil, errors.NewNotValid(errPEM, "[csjwt] parseRSAPrivateKeyFromPEM.DecryptPEMBlock")
		}
	} else {
		block = prKeyPEM.Bytes
	}

	var err error
	var parsedKey interface{}
	if parsedKey, err = x509.ParsePKCS1PrivateKey(block); err != nil {
		if parsedKey, err = x509.ParsePKCS8PrivateKey(block); err != nil {
			return nil, errors.NewNotValid(err, "[csjwt] parseRSAPrivateKeyFromPEM.ParsePKCS8PrivateKey")
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
			return nil, errors.NewNotValid(err, "[csjwt] parseRSAPublicKeyFromPEM.ParseCertificate")
		}
	}

	var pkey *rsa.PublicKey
	var ok bool
	if pkey, ok = parsedKey.(*rsa.PublicKey); !ok {
		return nil, errors.NewNotValidf(errKeyNonRSAPrivateKey)
	}

	return pkey, nil
}
