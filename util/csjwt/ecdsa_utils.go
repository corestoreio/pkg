package csjwt

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"

	"github.com/corestoreio/csfw/util/errors"
)

// Parse PEM encoded Elliptic Curve Private Key Structure
func parseECPrivateKeyFromPEM(key []byte, password ...[]byte) (*ecdsa.PrivateKey, error) {

	// Parse PEM block
	pemBlock, _ := pem.Decode(key)
	if pemBlock == nil {
		return nil, errors.NewNotSupportedf(errKeyMustBePEMEncoded)
	}

	var blockBytes []byte
	if x509.IsEncryptedPEMBlock(pemBlock) {
		if len(password) != 1 || len(password[0]) == 0 {
			return nil, errors.NewEmptyf(errKeyMissingPassword)
		}
		var errPEM error
		if blockBytes, errPEM = x509.DecryptPEMBlock(pemBlock, password[0]); errPEM != nil {
			return nil, errors.NewNotValid(errPEM, "[csjwt] parseECPrivateKeyFromPEM.x509.DecryptPEMBlock")
		}
	} else {
		blockBytes = pemBlock.Bytes
	}

	// Parse the key
	return x509.ParseECPrivateKey(blockBytes)
}

// Parse PEM encoded PKCS1 or PKCS8 public key
func parseECPublicKeyFromPEM(key []byte) (*ecdsa.PublicKey, error) {

	// Parse PEM block
	block, _ := pem.Decode(key)
	if block == nil {
		return nil, errors.NewNotSupportedf(errKeyMustBePEMEncoded)
	}

	// Parse the key
	parsedKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		cert, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			return nil, errors.Wrap(err, "[csjwt] parseECPublicKeyFromPEM.x509.ParseCertificate")
		}
		parsedKey = cert.PublicKey
	}

	if pkey, ok := parsedKey.(*ecdsa.PublicKey); ok {
		return pkey, nil
	}
	return nil, errors.NewNotValidf(errKeyNonECDSAPublicKey)
}
