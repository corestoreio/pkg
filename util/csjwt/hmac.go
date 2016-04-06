package csjwt

import (
	"crypto"
	"crypto/hmac"
)

// SigningMethodHMAC implements the HMAC-SHA family of signing methods signing methods
type SigningMethodHMAC struct {
	Name string
	Hash crypto.Hash
}

// Specific instances for HS256 and company
var (
	SigningMethodHS256 *SigningMethodHMAC
	SigningMethodHS384 *SigningMethodHMAC
	SigningMethodHS512 *SigningMethodHMAC
)

func init() {
	SigningMethodHS256 = &SigningMethodHMAC{HS256, crypto.SHA256}
	RegisterSigningMethod(SigningMethodHS256)

	SigningMethodHS384 = &SigningMethodHMAC{HS384, crypto.SHA384}
	RegisterSigningMethod(SigningMethodHS384)

	SigningMethodHS512 = &SigningMethodHMAC{HS512, crypto.SHA512}
	RegisterSigningMethod(SigningMethodHS512)
}

func (m *SigningMethodHMAC) Alg() string {
	return m.Name
}

// Verify the signature of HSXXX tokens.  Returns nil if the signature is valid.
// For the key you can use any of the WithPassword*() functions.
func (m *SigningMethodHMAC) Verify(signingString, signature []byte, key Key) error {
	// Verify the key is the right type
	if key.Error != nil {
		return key.Error
	}
	if len(key.hmacPassword) == 0 {
		return ErrInvalidKey
	}

	// Decode signature, for comparison
	sig, err := DecodeSegment(signature)
	if err != nil {
		return err
	}

	// Can we use the specified hashing method?
	if !m.Hash.Available() {
		return ErrHashUnavailable
	}

	// This signing method is symmetric, so we validate the signature
	// by reproducing the signature from the signing string and key, then
	// comparing that against the provided signature.
	hasher := hmac.New(m.Hash.New, key.hmacPassword)
	if _, err := hasher.Write(signingString); err != nil {
		return err
	}

	if !hmac.Equal(sig, hasher.Sum(nil)) {
		return ErrSignatureInvalid
	}

	// No validation errors.  Signature is good.
	return nil
}

// Sign implements the Sign method from SigningMethod interface.
// For the key you can use any of the WithPassword*() functions.
func (m *SigningMethodHMAC) Sign(signingString []byte, key Key) ([]byte, error) {

	if key.Error != nil {
		return nil, key.Error
	}
	if key.hmacPassword == nil {
		return nil, ErrInvalidKey
	}

	if !m.Hash.Available() {
		return nil, ErrHashUnavailable
	}

	hasher := hmac.New(m.Hash.New, key.hmacPassword)
	if _, err := hasher.Write(signingString); err != nil {
		return nil, err
	}

	return EncodeSegment(hasher.Sum(nil)), nil
}
