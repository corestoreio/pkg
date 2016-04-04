package csjwt

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
)

var (
	ErrKeyMustBePEMEncoded       = errors.New("Invalid Key: Key must be PEM encoded PKCS1 or PKCS8 private key")
	ErrNotRSAPrivateKey          = errors.New("Key is not a valid RSA private key")
	ErrPrivateKeyMissingPassword = errors.New("Missing password to decrypt private key")
)

// Parse PEM encoded PKCS1 or PKCS8 private key. Provide optionally a password if
// key is encrypted.
func parseRSAPrivateKeyFromPEM(privateKey []byte, password ...[]byte) (*rsa.PrivateKey, error) {
	var prKeyPEM *pem.Block
	if prKeyPEM, _ = pem.Decode(privateKey); prKeyPEM == nil {
		return nil, ErrKeyMustBePEMEncoded
	}

	var block []byte
	if x509.IsEncryptedPEMBlock(prKeyPEM) {
		if len(password) != 1 || len(password[0]) == 0 {
			return nil, ErrPrivateKeyMissingPassword
		}
		var errPEM error
		if block, errPEM = x509.DecryptPEMBlock(prKeyPEM, password[0]); errPEM != nil {
			return nil, errPEM
		}
	} else {
		block = prKeyPEM.Bytes
	}

	var err error
	var parsedKey interface{}
	if parsedKey, err = x509.ParsePKCS1PrivateKey(block); err != nil {
		if parsedKey, err = x509.ParsePKCS8PrivateKey(block); err != nil {
			return nil, err
		}
	}

	var pkey *rsa.PrivateKey
	var ok bool
	if pkey, ok = parsedKey.(*rsa.PrivateKey); !ok {
		return nil, ErrNotRSAPrivateKey
	}

	return pkey, nil
}

// Parse PEM encoded PKCS1 or PKCS8 public key
func parseRSAPublicKeyFromPEM(key []byte) (*rsa.PublicKey, error) {
	var err error

	// Parse PEM block
	var block *pem.Block
	if block, _ = pem.Decode(key); block == nil {
		return nil, ErrKeyMustBePEMEncoded
	}

	// Parse the key
	var parsedKey interface{}
	if parsedKey, err = x509.ParsePKIXPublicKey(block.Bytes); err != nil {
		if cert, err := x509.ParseCertificate(block.Bytes); err == nil {
			parsedKey = cert.PublicKey
		} else {
			return nil, err
		}
	}

	var pkey *rsa.PublicKey
	var ok bool
	if pkey, ok = parsedKey.(*rsa.PublicKey); !ok {
		return nil, ErrNotRSAPrivateKey
	}

	return pkey, nil
}
