package csjwt

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"

	"github.com/corestoreio/errors"
)

// SigningMethodRSAPSS implements the RSAPSS family of signing methods signing
// methods.
type SigningMethodRSAPSS struct {
	SigningMethodRSA
	Options rsa.PSSOptions
}

func newSigningMethodRSAPSS(n string, h crypto.Hash) *SigningMethodRSAPSS {
	return &SigningMethodRSAPSS{
		SigningMethodRSA{
			Name: n,
			Hash: h,
		},
		rsa.PSSOptions{
			SaltLength: rsa.PSSSaltLengthAuto,
			Hash:       h,
		},
	}
}

// NewSigningMethodPS256 creates a new 256bit RSAPSS SHA instance and registers
// it.
func NewSigningMethodPS256() *SigningMethodRSAPSS {
	return newSigningMethodRSAPSS(PS256, crypto.SHA256)
}

// NewSigningMethodPS384 creates a new 384bit RSAPSS SHA instance and registers
// it.
func NewSigningMethodPS384() *SigningMethodRSAPSS {
	return newSigningMethodRSAPSS(PS384, crypto.SHA384)
}

// NewSigningMethodPS512 creates a new 512bit RSAPSS SHA instance and registers
// it.
func NewSigningMethodPS512() *SigningMethodRSAPSS {
	return newSigningMethodRSAPSS(PS512, crypto.SHA512)
}

// Verify implements the Verify method from SigningMethod interface. For the key
// you can use any of the WithRSA*Key*() functions. Error behaviour: Empty,
// NotImplemented, WriteFailed, NotValid
func (m *SigningMethodRSAPSS) Verify(signingString, signature []byte, key Key) error {
	if key.Error != nil {
		return errors.Wrap(key.Error, "[csjwt] SigningMethodRSAPSS.Verify.key")
	}
	if key.rsaKeyPub == nil {
		return errors.Empty.Newf(errRSAPublicKeyEmpty)
	}

	// Decode the signature
	sig, err := DecodeSegment(signature)
	if err != nil {
		return errors.Wrap(err, "[csjwt] SigningMethodRSAPSS.Verify.DecodeSegment")
	}

	// Create hasher
	if !m.Hash.Available() {
		return errors.NotImplemented.Newf(errRSAHashUnavailable)
	}
	hasher := m.Hash.New()
	if _, err := hasher.Write(signingString); err != nil {
		return errors.WriteFailed.New(err, "[csjwt] SigningMethodRSA.Verify.hasher.Write")
	}

	return errors.NotValid.New(rsa.VerifyPSS(key.rsaKeyPub, m.Hash, hasher.Sum(nil), sig, &m.Options), "[csjwt] SigningMethodRSAPSS.Verify.VerifyPSS")
}

// Sign implements the Sign method from SigningMethod interface. For the key you
// can use any of the WithRSAPrivateKey*() functions. Error behaviour: Empty,
// NotImplemented, WriteFailed, NotValid.
func (m *SigningMethodRSAPSS) Sign(signingString []byte, key Key) ([]byte, error) {
	if key.Error != nil {
		return nil, errors.Wrap(key.Error, "[csjwt] SigningMethodRSAPSS.Sign.key")
	}
	if key.rsaKeyPriv == nil {
		return nil, errors.Empty.Newf(errRSAPrivateKeyEmpty)
	}

	// Create the hasher
	if !m.Hash.Available() {
		return nil, errors.NotImplemented.Newf(errRSAHashUnavailable)
	}

	hasher := m.Hash.New()
	if _, err := hasher.Write(signingString); err != nil {
		return nil, errors.WriteFailed.New(err, "[csjwt] SigningMethodRSA.Sign.hasher.Write")
	}

	// Sign the string and return the encoded bytes
	sigBytes, err := rsa.SignPSS(rand.Reader, key.rsaKeyPriv, m.Hash, hasher.Sum(nil), &m.Options)
	if err != nil {
		return nil, errors.NotValid.New(err, "[csjwt] SigningMethodRSAPSS.Sign.SignPSS")
	}
	return EncodeSegment(sigBytes), nil
}
