package jwtclaim

import "github.com/corestoreio/errors"

//go:generate ffjson $GOFILE

// Header constants define the main headers used for Set() and Get() functions.
// Those constants are implemented in the HeaderSegments type.
// TODO(cs) add more constants
const (
	HeaderAlg = "alg"
	HeaderTyp = "typ"
)

// ContentTypeJWT defines the content type of a token. At the moment only JWT is
// supported. JWE may be added in the future JSON Web Encryption (JWE).
// https://tools.ietf.org/html/rfc7519
const ContentTypeJWT = `JWT`

// HeadSegments represents a structured version of Header Section, as
// referenced at http://self-issued.info/docs/draft-jones-json-web-token-01.html#anchor5
type HeadSegments struct {
	// Alg (algorithm) header parameter identifies the cryptographic algorithm
	// used to secure the JWT. A list of reserved alg values is in Table 4. The
	// processing of the "alg" (algorithm) header parameter, if present,
	// requires that the value of the "alg" header parameter MUST be one that is
	// both supported and for which there exists a key for use with that
	// algorithm associated with the issuer of the JWT. This header parameter is
	// REQUIRED.
	Algorithm string `json:"alg,omitempty"`
	// Typ (type) header parameter is used to declare that this data structure
	// is a JWT. If a "typ" parameter is present, it is RECOMMENDED that its
	// value be "JWT". This header parameter is OPTIONAL.
	Type string `json:"typ,omitempty"`
	// JKU (JSON Key URL) header parameter is a URL that points to JSON-encoded
	// public key certificates that can be used to validate the signature. The
	// specification for this encoding is TBD. This header parameter is
	// OPTIONAL.
	JKU string `json:"jku,omitempty"`
	// KID (key ID) header parameter is a hint indicating which specific key
	// owned by the signer should be used to validate the signature. This allows
	// signers to explicitly signal a change of key to recipients. Omitting this
	// parameter is equivalent to setting it to an empty string. The
	// interpretation of the contents of the "kid" parameter is unspecified.
	// This header parameter is OPTIONAL.
	KID string `json:"kid,omitempty"`
	// X5U (X.509 URL) header parameter is a URL that points to an X.509 public
	// key certificate that can be used to validate the signature. This
	// certificate MUST conform to RFC 5280 [RFC5280]. This header parameter is
	// OPTIONAL.
	X5U string `json:"x5u,omitempty"`
	// X5T (x.509 certificate thumbprint) header parameter provides a base64url
	// encoded SHA-256 thumbprint (a.k.a. digest) of the DER encoding of an
	// X.509 certificate that can be used to match a certificate. This header
	// parameter is OPTIONAL.
	X5T string `json:"x5t,omitempty"`
}

// NewHeadSegments creates a new header with an optional algorithm. Algorithm
// may be set later depending on the signing method.
func NewHeadSegments(alg ...string) *HeadSegments {
	hs := &HeadSegments{
		Type: ContentTypeJWT,
	}
	if len(alg) == 1 {
		hs.Algorithm = alg[0]
	}
	return hs
}

// Alg returns the underlying algorithm.
func (s *HeadSegments) Alg() string {
	return s.Algorithm
}

// Typ returns the token type.
func (s *HeadSegments) Typ() string {
	return s.Type
}

// Set sets a value. Key must be one of the constants Header*. Error behaviour:
// NotSupported
func (s *HeadSegments) Set(key, value string) (err error) {
	switch key {
	case HeaderAlg:
		s.Algorithm = value
	case HeaderTyp:
		s.Type = value
	default:
		return errors.NewNotSupportedf(errHeaderKeyNotSupported, key)
	}
	return err
}

// Get returns a value or nil or an error. Key must be one of the constants
// Header*. Error behaviour: NotSupported.
func (s *HeadSegments) Get(key string) (value string, err error) {
	switch key {
	case HeaderAlg:
		return s.Algorithm, nil
	case HeaderTyp:
		return s.Type, nil
	}
	return "", errors.NewNotSupportedf(errHeaderKeyNotSupported, key)
}
