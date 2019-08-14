package csjwt

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/rand"
	"math/big"

	"github.com/corestoreio/errors"
)

// SigningMethodECDSA implements the ECDSA family of signing methods signing
// methods.
type SigningMethodECDSA struct {
	Name      string
	Hash      crypto.Hash
	KeySize   int
	CurveBits int
}

func newSigningMethodECDSA(n string, h crypto.Hash, keySize, curveBits int) *SigningMethodECDSA {
	return &SigningMethodECDSA{Name: n, Hash: h, KeySize: keySize, CurveBits: curveBits}
}

// NewSigningMethodES256 creates a new 256bit ECDSA SHA instance and registers
// it.
func NewSigningMethodES256() *SigningMethodECDSA {
	return newSigningMethodECDSA(ES256, crypto.SHA256, 32, 256)
}

// NewSigningMethodES384 creates a new 384bit ECDSA SHA instance and registers
// it.
func NewSigningMethodES384() *SigningMethodECDSA {
	return newSigningMethodECDSA(ES384, crypto.SHA384, 48, 384)
}

// NewSigningMethodES512 creates a new 512bit ECDSA SHA instance and registers
// it.
func NewSigningMethodES512() *SigningMethodECDSA {
	return newSigningMethodECDSA(ES512, crypto.SHA512, 66, 521)
}

// Alg returns the name of the underlying algorithm.
func (m *SigningMethodECDSA) Alg() string {
	return m.Name
}

// Verify implements the Verify method from SigningMethod interface. For the key
// you can use any of the WithEC*Key*() functions Error behaviour: Empty,
// NotImplemented, WriteFailed, NotValid.
func (m *SigningMethodECDSA) Verify(signingString, signature []byte, key Key) error {
	// Get the key
	if key.Error != nil {
		return errors.Wrap(key.Error, "[csjwt] SigningMethodECDSA.Verify.key")
	}
	if key.ecdsaKeyPub == nil {
		return errors.Empty.Newf(errECDSAPublicKeyEmpty)
	}

	// Decode the signature
	sig, err := DecodeSegment(signature)
	if err != nil {
		return errors.Wrap(err, "[csjwt] SigningMethodECDSA.Verify.DecodeSegment")
	}

	if len(sig) != 2*m.KeySize {
		return errors.NotValid.Newf(errECDSAVerification)
	}

	r := big.NewInt(0).SetBytes(sig[:m.KeySize])
	s := big.NewInt(0).SetBytes(sig[m.KeySize:])

	// Create hasher
	if !m.Hash.Available() {
		return errors.NotImplemented.Newf(errECDSAHashUnavailable)
	}

	hasher := m.Hash.New()
	if _, err := hasher.Write(signingString); err != nil {
		return errors.WriteFailed.New(err, "[csjwt] SigningMethodECDSA.Verify.hasher.Write")
	}

	// Verify the signature
	if !ecdsa.Verify(key.ecdsaKeyPub, hasher.Sum(nil), r, s) {
		return errors.NotValid.Newf(errECDSAVerification)
	}
	return nil
}

// Sign implements the Sign method from SigningMethod. For the key you can use
// any of the WithECPrivateKey*() functions. Error behaviour: Empty,
// NotImplemented, WriteFailed, NotValid.
func (m *SigningMethodECDSA) Sign(signingString []byte, key Key) ([]byte, error) {
	if key.Error != nil {
		return nil, errors.Wrap(key.Error, "[csjwt] SigningMethodECDSA.Sign.key")
	}
	if key.ecdsaKeyPriv == nil {
		return nil, errors.Empty.Newf(errECDSAPrivateKeyEmpty)
	}

	// Create the hasher
	if !m.Hash.Available() {
		return nil, errors.NotImplemented.Newf(errECDSAHashUnavailable)
	}

	hasher := m.Hash.New()
	if _, err := hasher.Write(signingString); err != nil {
		return nil, errors.WriteFailed.New(err, "[csjwt] SigningMethodECDSA.Sign.hasher.Write")
	}

	// Sign the string and return r, s
	r, s, err := ecdsa.Sign(rand.Reader, key.ecdsaKeyPriv, hasher.Sum(nil))
	if err != nil {
		return nil, errors.NotValid.New(err, "[csjwt] SigningMethodECDSA.Sign.ecdsa.Sign")
	}

	curveBits := key.ecdsaKeyPriv.Curve.Params().BitSize

	if m.CurveBits != curveBits {
		return nil, errors.NotValid.Newf(errECDSAPrivateInvalidBits)
	}

	keyBytes := curveBits / 8
	if curveBits%8 > 0 {
		keyBytes++
	}

	// We serialize the outpus (r and s) into big-endian byte arrays and pad
	// them with zeros on the left to make sure the sizes work out. Both arrays
	// must be keyBytes long, and the output must be 2*keyBytes long.
	rBytes := r.Bytes()
	rBytesPadded := make([]byte, keyBytes)
	copy(rBytesPadded[keyBytes-len(rBytes):], rBytes)

	sBytes := s.Bytes()
	sBytesPadded := make([]byte, keyBytes)
	copy(sBytesPadded[keyBytes-len(sBytes):], sBytes)

	out := append(rBytesPadded, sBytesPadded...)

	return EncodeSegment(out), nil
}
