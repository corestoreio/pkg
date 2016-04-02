package csjwt

import (
	"crypto/ecdsa"
	"crypto/rsa"
	"io/ioutil"

	"github.com/juju/errors"
)

// ErrHMACEmptyPassword whenever the password length is 0.
var ErrHMACEmptyPassword = errors.New("Empty passwords are forbidden")

// Keyfunc used by Parse methods, this callback function supplies
// the key for verification.  The function receives the parsed,
// but unverified Token.  This allows you to use propries in the
// Header of the token (such as `kid`) to identify which key to use.
type Keyfunc func(Token) (Key, error)

// Key defines a container for the HMAC password, RSA and ECDSA public and
// private keys. The Error fields gets filled out when loading/parsing the keys.
type Key struct {
	hmacPassword []byte
	ecdsaKeyPub  *ecdsa.PublicKey
	ecdsaKeyPriv *ecdsa.PrivateKey
	rsaKeyPub    *rsa.PublicKey
	rsaKeyPriv   *rsa.PrivateKey
	Error        error
}

const goStringTpl = `csjwt.Key{}`

// GoString protects keys and enforces privacy.
func (k Key) GoString() string {
	return goStringTpl
}

// String protects keys and enforces privacy.
func (k Key) String() string {
	return goStringTpl
}

// IsEmpty returns true when no field has been used in the Key struct.
// Error is excluded from the check
func (k Key) IsEmpty() bool {
	return k.hmacPassword == nil && k.ecdsaKeyPub == nil && k.ecdsaKeyPriv == nil && k.rsaKeyPub == nil && k.rsaKeyPriv == nil
}

// WithPassword uses the byte slice as the password for the HMAC-SHA signing method.
func WithPassword(password []byte) Key {
	var err error
	if len(password) == 0 {
		err = ErrHMACEmptyPassword
	}
	return Key{
		hmacPassword: password,
		Error:        err,
	}
}

// WithPasswordFromFile loads the content of a file and uses that content as
// the password for the HMAC-SHA signing method.
func WithPasswordFromFile(pathToFile string) Key {
	var k Key
	k.hmacPassword, k.Error = ioutil.ReadFile(pathToFile)
	if k.Error != nil {
		k.Error = k.Error
	}
	return k
}

// WithRSAPublicKeyFromPEM parses PEM encoded PKCS1 or PKCS8 public key
func WithRSAPublicKeyFromPEM(publicKey []byte) (k Key) {
	k.rsaKeyPub, k.Error = parseRSAPublicKeyFromPEM(publicKey)
	return
}

// WithRSAPublicKeyFromFile parses PEM encoded PKCS1 or PKCS8 public key found
// in a file.
func WithRSAPublicKeyFromFile(pathToFile string) (k Key) {
	pk, err := ioutil.ReadFile(pathToFile)
	if err != nil {
		k.Error = errors.Errorf("%s with file %s", err, pathToFile)
		return k
	}
	return WithRSAPublicKeyFromPEM(pk)
}

// WithRSAPublicKey sets the public key
func WithRSAPublicKey(publicKey *rsa.PublicKey) (k Key) {
	k.rsaKeyPub = publicKey
	return
}

// WithRSAPrivateKeyFromPEM parses PEM encoded PKCS1 or PKCS8 private key
func WithRSAPrivateKeyFromPEM(privateKey []byte) (k Key) {
	k.rsaKeyPriv, k.Error = parseRSAPrivateKeyFromPEM(privateKey)
	return
}

// WithRSAPrivateKeyFromFile parses PEM encoded PKCS1 or PKCS8 private key
// found in a file.
func WithRSAPrivateKeyFromFile(pathToFile string) (k Key) {
	pk, err := ioutil.ReadFile(pathToFile)
	if err != nil {
		k.Error = err
		return k
	}
	return WithRSAPrivateKeyFromPEM(pk)
}

// WithRSAPrivateKey sets the private key
func WithRSAPrivateKey(privateKey *rsa.PrivateKey) (k Key) {
	k.rsaKeyPriv = privateKey
	return
}

// WithECPublicKeyFromPEM parses PEM encoded Elliptic Curve Public Key Structure
func WithECPublicKeyFromPEM(publicKey []byte) (k Key) {
	k.ecdsaKeyPub, k.Error = parseECPublicKeyFromPEM(publicKey)
	return
}

// WithECPublicKeyFromFile parses a file PEM encoded Elliptic Curve Public Key Structure
func WithECPublicKeyFromFile(pathToFile string) (k Key) {
	pk, err := ioutil.ReadFile(pathToFile)
	if err != nil {
		k.Error = err
		return k
	}
	k.ecdsaKeyPub, k.Error = parseECPublicKeyFromPEM(pk)
	return
}

// WithECPublicKey sets the ECDSA public key
func WithECPublicKey(publicKey *ecdsa.PublicKey) (k Key) {
	k.ecdsaKeyPub = publicKey
	return
}

// WithECPrivateKeyFromPEM parses PEM encoded Elliptic Curve Private Key Structure
func WithECPrivateKeyFromPEM(privateKey []byte) (k Key) {
	k.ecdsaKeyPriv, k.Error = parseECPrivateKeyFromPEM(privateKey)
	return
}

// WithECPrivateKeyFromFile parses file PEM encoded Elliptic Curve Private Key Structure
func WithECPrivateKeyFromFile(pathToFile string) (k Key) {
	pk, err := ioutil.ReadFile(pathToFile)
	if err != nil {
		k.Error = err
		return k
	}
	k.ecdsaKeyPriv, k.Error = parseECPrivateKeyFromPEM(pk)
	return
}

// WithECPrivateKey sets the ECDSA private key
func WithECPrivateKey(privateKey *ecdsa.PrivateKey) (k Key) {
	k.ecdsaKeyPriv = privateKey
	return
}
