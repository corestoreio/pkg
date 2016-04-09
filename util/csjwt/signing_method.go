package csjwt

import "bytes"

// Signer interface to add new methods for signing or verifying tokens.
type Signer interface {
	// Verify returns nil if signature is valid
	Verify(signingString, signature []byte, key Key) error
	// Sign returns encoded signature or error
	Sign(signingString []byte, key Key) ([]byte, error)
	// Alg returns the alg identifier for this method (example: 'HS256')
	Alg() string
}

// All available algorithms which are supported by this package.
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
	PS    = `PS`
	RS    = `RS`
)

type methods []Signer

func (ms methods) String() string {
	var buf bytes.Buffer
	for i, m := range ms {
		_, _ = buf.WriteString(m.Alg())
		if i < len(ms)-1 {
			_, _ = buf.WriteString(`, `)
		}
	}
	return buf.String()
}
