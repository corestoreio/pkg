package csjwt

import "github.com/juju/errors"

var signingMethods = make(map[string]Signer)

// Signer interface to add new methods for signing or verifying tokens.
type Signer interface {
	// Verify returns nil if signature is valid
	Verify(signingString, signature []byte, key Key) error
	// Sign returns encoded signature or error
	Sign(signingString []byte, key Key) ([]byte, error)
	// Alg returns the alg identifier for this method (example: 'HS256')
	Alg() string
}

// RegisterSigningMethod registers the "alg" name and Signer interface
// implementation for a signing method.
func RegisterSigningMethod(s Signer) {
	signingMethods[s.Alg()] = s
}

// GetSigningMethod returns a signing method from an "alg" string.
func GetSigningMethod(alg string) (Signer, error) {
	if s, ok := signingMethods[alg]; ok {
		return s, nil
	}
	return nil, errors.Errorf("SigningMethod %q not registered", alg)
}

// All available algorithms which can be supported
const (
	ES256 = `ES256`
	ES384 = `ES384`
	ES512 = `ES512`
	HS256 = `HS256`
	HS384 = `HS384`
	HS512 = `HS512`
	PS256 = `PS256`
	PS384 = `PS384`
	PS512 = `PS512`
	RS256 = `RS256`
	RS384 = `RS384`
	RS512 = `RS512`
	ES    = `ES`
	HS    = `HS`
	RS    = `RS`
)
