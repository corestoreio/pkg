package csjwt

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
)

// SigningMethodRSAPSS implements the RSAPSS family of signing methods signing methods
type SigningMethodRSAPSS struct {
	*SigningMethodRSA
	Options *rsa.PSSOptions
}

// Specific instances for RS/PS and company
var (
	SigningMethodPS256 *SigningMethodRSAPSS
	SigningMethodPS384 *SigningMethodRSAPSS
	SigningMethodPS512 *SigningMethodRSAPSS
)

func init() {
	// PS256
	SigningMethodPS256 = &SigningMethodRSAPSS{
		&SigningMethodRSA{
			Name: "PS256",
			Hash: crypto.SHA256,
		},
		&rsa.PSSOptions{
			SaltLength: rsa.PSSSaltLengthAuto,
			Hash:       crypto.SHA256,
		},
	}
	RegisterSigningMethod(SigningMethodPS256)

	// PS384
	SigningMethodPS384 = &SigningMethodRSAPSS{
		&SigningMethodRSA{
			Name: "PS384",
			Hash: crypto.SHA384,
		},
		&rsa.PSSOptions{
			SaltLength: rsa.PSSSaltLengthAuto,
			Hash:       crypto.SHA384,
		},
	}
	RegisterSigningMethod(SigningMethodPS384)

	// PS512
	SigningMethodPS512 = &SigningMethodRSAPSS{
		&SigningMethodRSA{
			Name: "PS512",
			Hash: crypto.SHA512,
		},
		&rsa.PSSOptions{
			SaltLength: rsa.PSSSaltLengthAuto,
			Hash:       crypto.SHA512,
		},
	}
	RegisterSigningMethod(SigningMethodPS512)
}

// Verify implements the Verify method from SigningMethod interface.
// For the key you can use any of the WithRSAPublicKey*() functions.
func (m *SigningMethodRSAPSS) Verify(signingString, signature []byte, key Key) error {
	if key.Error != nil {
		return key.Error
	}
	if key.rsaKeyPub == nil {
		return ErrInvalidKey
	}

	// Decode the signature
	sig, err := DecodeSegment(signature)
	if err != nil {
		return err
	}

	// Create hasher
	if !m.Hash.Available() {
		return ErrHashUnavailable
	}
	hasher := m.Hash.New()
	if _, err := hasher.Write(signingString); err != nil {
		return err
	}

	return rsa.VerifyPSS(key.rsaKeyPub, m.Hash, hasher.Sum(nil), sig, m.Options)
}

// Sign implements the Sign method from SigningMethod interface.
// For the key you can use any of the WithRSAPrivateKey*() functions.
func (m *SigningMethodRSAPSS) Sign(signingString []byte, key Key) ([]byte, error) {
	if key.Error != nil {
		return nil, key.Error
	}
	if key.rsaKeyPriv == nil {
		return nil, ErrInvalidKey
	}

	// Create the hasher
	if !m.Hash.Available() {
		return nil, ErrHashUnavailable
	}

	hasher := m.Hash.New()
	if _, err := hasher.Write(signingString); err != nil {
		return nil, err
	}

	// Sign the string and return the encoded bytes
	sigBytes, err := rsa.SignPSS(rand.Reader, key.rsaKeyPriv, m.Hash, hasher.Sum(nil), m.Options)
	if err != nil {
		return nil, err
	}
	return EncodeSegment(sigBytes), nil
}
