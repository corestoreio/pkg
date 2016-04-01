package csjwt

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
// Returns nil for not found.
func GetSigningMethod(alg string) Signer {
	if s, ok := signingMethods[alg]; ok {
		return s
	}
	return nil
}
