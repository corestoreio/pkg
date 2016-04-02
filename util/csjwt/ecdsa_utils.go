package csjwt

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
)

var (
	ErrNotECPublicKey = errors.New("Key is not a valid ECDSA public key")
)

// Parse PEM encoded Elliptic Curve Private Key Structure
func parseECPrivateKeyFromPEM(key []byte) (*ecdsa.PrivateKey, error) {

	// Parse PEM block
	block, _ := pem.Decode(key)
	if block == nil {
		return nil, ErrKeyMustBePEMEncoded
	}

	// Parse the key
	return x509.ParseECPrivateKey(block.Bytes)
}

// Parse PEM encoded PKCS1 or PKCS8 public key
func parseECPublicKeyFromPEM(key []byte) (*ecdsa.PublicKey, error) {

	// Parse PEM block
	block, _ := pem.Decode(key)
	if block == nil {
		return nil, ErrKeyMustBePEMEncoded
	}

	// Parse the key
	parsedKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		cert, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			return nil, err
		}
		parsedKey = cert.PublicKey
	}

	if pkey, ok := parsedKey.(*ecdsa.PublicKey); ok {
		return pkey, nil
	}
	return nil, ErrNotECPublicKey
}
