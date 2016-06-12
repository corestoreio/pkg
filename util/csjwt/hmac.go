package csjwt

import (
	"crypto"
	"crypto/hmac"

	"github.com/corestoreio/csfw/util/errors"
)

// SigningMethodHMAC implements the HMAC-SHA family of signing methods signing
// methods.
type SigningMethodHMAC struct {
	Name string
	Hash crypto.Hash
}

func newSigningMethodHMAC(n string, h crypto.Hash) *SigningMethodHMAC {
	return &SigningMethodHMAC{Name: n, Hash: h}
}

// NewSigningMethodHS256 creates a new 256bit HMAC SHA instance and registers
// it.
func NewSigningMethodHS256() *SigningMethodHMAC {
	return newSigningMethodHMAC(HS256, crypto.SHA256)
}

// NewSigningMethodHS384 creates a new 384bit HMAC SHA instance and registers
// it.
func NewSigningMethodHS384() *SigningMethodHMAC {
	return newSigningMethodHMAC(HS384, crypto.SHA384)
}

// NewSigningMethodHS512 creates a new 512bit HMAC SHA instance and registers
// it.
func NewSigningMethodHS512() *SigningMethodHMAC {
	return newSigningMethodHMAC(HS512, crypto.SHA512)
}

// Alg returns the used algorithm.
func (m *SigningMethodHMAC) Alg() string {
	return m.Name
}

// Verify the signature of HSXXX tokens.  Returns nil if the signature is valid.
// For the key you can use any of the WithPassword*() functions. Error
// behaviour: Empty, NotImplemented, WriteFailed, NotValid.
func (m *SigningMethodHMAC) Verify(signingString, signature []byte, key Key) error {
	// Verify the key is the right type
	if key.Error != nil {
		return errors.Wrap(key.Error, "[csjwt] SigningMethodHMAC.Verify.key")
	}
	if len(key.hmacPassword) == 0 {
		return errors.NewEmptyf(errHmacPasswordEmpty)
	}

	// Decode signature, for comparison
	sig, err := DecodeSegment(signature)
	if err != nil {
		return errors.Wrap(err, "[csjwt] SigningMethodHMAC.Verify.DecodeSegment")
	}

	// Can we use the specified hashing method?
	if !m.Hash.Available() {
		return errors.NewNotImplementedf(errHmacHashUnavailable)
	}

	// This signing method is symmetric, so we validate the signature by
	// reproducing the signature from the signing string and key, then comparing
	// that against the provided signature.
	hasher := hmac.New(m.Hash.New, key.hmacPassword)
	if _, err := hasher.Write(signingString); err != nil {
		return errors.NewWriteFailed(err, "[csjwt] SigningMethodHMAC.Verify.hasher.Write")
	}

	if !hmac.Equal(sig, hasher.Sum(nil)) {
		return errors.NewNotValidf(errHmacSignatureInvalid)
	}

	// No validation errors.  Signature is good.
	return nil
}

// Sign implements the Sign method from SigningMethod interface. For the key you
// can use any of the WithPassword*() functions. Error behaviour: Empty,
// NotImplemented, WriteFailed
func (m *SigningMethodHMAC) Sign(signingString []byte, key Key) ([]byte, error) {

	if key.Error != nil {
		return nil, errors.Wrap(key.Error, "[csjwt] SigningMethodHMAC.Sign.key")
	}
	if key.hmacPassword == nil {
		return nil, errors.NewEmptyf(errHmacPasswordEmpty)
	}

	if !m.Hash.Available() {
		return nil, errors.NewNotImplementedf(errHmacHashUnavailable)
	}

	hasher := hmac.New(m.Hash.New, key.hmacPassword)
	if _, err := hasher.Write(signingString); err != nil {
		return nil, errors.NewWriteFailed(err, "[csjwt] SigningMethodHMAC.Sign.hasher.Write")
	}

	return EncodeSegment(hasher.Sum(nil)), nil
}
