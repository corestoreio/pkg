package csjwt

import "fmt"

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

// Algorithm supported Alg[orithms] and hash size by this package.
type Algorithm uint

// All available algorithms which this package supports
const (
	ES256 Algorithm = 1 + iota
	ES384
	ES512
	HS256
	HS384
	HS512
	PS256
	PS384
	PS512
	RS256
	RS384
	RS512
	_maxAlgorithm
	ES
	HS
	RS
)

const _Algorithm_name = "ES256ES384ES512HS256HS384HS512PS256PS384PS512RS256RS384RS512"
const _Alogrithm_template = "Algorithm(%d)"

var _Algorithm_index = [...]uint8{0, 5, 10, 15, 20, 25, 30, 35, 40, 45, 50, 55, 60}

func (i Algorithm) String() string {
	j := i - 1
	if i == 0 || j+1 >= _maxAlgorithm {
		return fmt.Sprintf(_Alogrithm_template, j+1)
	}
	return _Algorithm_name[_Algorithm_index[j]:_Algorithm_index[j+1]]
}

// ToAlgorithm converts a string to an Algorithm type. Returns 0 on error.
func ToAlgorithm(s string) Algorithm {

	for i := range _Algorithm_index {
		if Algorithm(i+1) >= _maxAlgorithm {
			var j Algorithm
			_, _ = fmt.Sscanf(s, _Alogrithm_template, &j)
			return j
		}

		if _Algorithm_name[_Algorithm_index[i]:_Algorithm_index[i+1]] == s {
			return Algorithm(i + 1)
		}
	}
	return 0
}
