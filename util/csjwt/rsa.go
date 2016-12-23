package csjwt

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"

	"github.com/corestoreio/errors"
)

// SigningMethodRSA implements the RSA family of signing methods signing methods
type SigningMethodRSA struct {
	Name string
	Hash crypto.Hash
}

func newSigningMethodRSA(n string, h crypto.Hash) *SigningMethodRSA {
	return &SigningMethodRSA{Name: n, Hash: h}
}

// NewSigningMethodRS256 creates a new 256bit RSA SHA instance and registers it.
func NewSigningMethodRS256() *SigningMethodRSA {
	return newSigningMethodRSA(RS256, crypto.SHA256)
}

// NewSigningMethodRS384 creates a new 384bit RSA SHA instance and registers it.
func NewSigningMethodRS384() *SigningMethodRSA {
	return newSigningMethodRSA(RS384, crypto.SHA384)
}

// NewSigningMethodRS512 creates a new 512bit RSA SHA instance and registers it.
func NewSigningMethodRS512() *SigningMethodRSA {
	return newSigningMethodRSA(RS512, crypto.SHA512)
}

func (m *SigningMethodRSA) Alg() string {
	return m.Name
}

// Verify implements the Verify method from SigningMethod interface. For the key
// you can use any of the WithRSA*Key*() functions. Error behaviour: Empty,
// NotImplemented, WriteFailed, NotValid
func (m *SigningMethodRSA) Verify(signingString, signature []byte, key Key) error {
	if key.Error != nil {
		return errors.Wrap(key.Error, "[csjwt] SigningMethodRSA.Verify.key")
	}
	if key.rsaKeyPub == nil {
		return errors.NewEmptyf(errRSAPublicKeyEmpty)
	}

	// Decode the signature
	sig, err := DecodeSegment(signature)
	if err != nil {
		return errors.Wrap(err, "[csjwt] SigningMethodRSA.Verify.DecodeSegment")
	}

	// Create hasher
	if !m.Hash.Available() {
		return errors.NewNotImplementedf(errRSAHashUnavailable)
	}
	hasher := m.Hash.New()
	if _, err := hasher.Write(signingString); err != nil {
		return errors.NewWriteFailed(err, "[csjwt] SigningMethodRSA.Verify.hasher.Write")
	}

	// Verify the signature
	return errors.NewNotValid(rsa.VerifyPKCS1v15(key.rsaKeyPub, m.Hash, hasher.Sum(nil), sig), "[csjwt] SigningMethodRSA.Verify.VerifyPKCS1v15")
}

// Sign implements the Sign method from SigningMethod interface. For the key you
// can use any of the WithRSAPrivateKey*() functions. Error behaviour: Empty,
// NotImplemented, WriteFailed, NotValid.
func (m *SigningMethodRSA) Sign(signingString []byte, key Key) ([]byte, error) {
	if key.Error != nil {
		return nil, errors.Wrap(key.Error, "[csjwt] SigningMethodRSA.Sign.key")
	}
	if key.rsaKeyPriv == nil {
		return nil, errors.NewEmptyf(errRSAPrivateKeyEmpty)
	}

	// Create the hasher
	if !m.Hash.Available() {
		return nil, errors.NewNotImplementedf(errRSAHashUnavailable)
	}

	hasher := m.Hash.New()
	if _, err := hasher.Write(signingString); err != nil {
		return nil, errors.NewWriteFailed(err, "[csjwt] SigningMethodRSA.Sign.hasher.Write")
	}

	// Sign the string and return the encoded bytes
	sigBytes, err := rsa.SignPKCS1v15(rand.Reader, key.rsaKeyPriv, m.Hash, hasher.Sum(nil))
	if err != nil {
		return nil, errors.NewNotValid(err, "[csjwt] SigningMethodRSA.Sign.SignPKCS1v15")
	}
	return EncodeSegment(sigBytes), nil
}
