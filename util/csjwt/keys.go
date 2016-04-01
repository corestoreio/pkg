package csjwt

import (
	"crypto/ecdsa"
	"crypto/rsa"
	"github.com/juju/errors"
	"io/ioutil"
)

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

// GoString protects keys and enforces privacy.
func (k Key) GoString() string {
	return `csjwt.Key{}`
}

// IsEmpty returns true when no field has been used in the Key struct.
func (k Key) IsEmpty() bool {
	return k.hmacPassword == nil && k.ecdsaKeyPub == nil && k.ecdsaKeyPriv == nil && k.rsaKeyPub == nil && k.rsaKeyPriv == nil && k.Error == nil
}

// WithPassword uses the byte slice as the password for the HMAC-SHA signing method.
func WithPassword(password []byte) Key {
	return Key{
		hmacPassword: password,
	}
}

// WithPasswordFromFile loads the content of a file and uses that content as
// the password for the HMAC-SHA signing method.
func WithPasswordFromFile(pathToFile string) Key {
	var k Key
	k.hmacPassword, k.Error = ioutil.ReadFile(pathToFile)
	if k.Error != nil {
		k.Error = errors.Errorf("%s with file %s", k.Error, pathToFile)
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
		k.Error = errors.Errorf("%s with file %s", err, pathToFile)
		return k
	}
	return WithRSAPrivateKeyFromPEM(pk)
}

// WithRSAPrivateKey sets the private key
func WithRSAPrivateKey(privateKey *rsa.PrivateKey) (k Key) {
	k.rsaKeyPriv = privateKey
	return
}

func WithECPublicKeyFromPEM(publicKey []byte) (k Key) {
	k.ecdsaKeyPub, k.Error = ParseECPublicKeyFromPEM(publicKey)
	return
}

func WithECPublicKey(publicKey *ecdsa.PublicKey) (k Key) {
	k.ecdsaKeyPub = publicKey
	return
}

func WithECPrivateKeyFromPEM(privateKey []byte) (k Key) {
	k.ecdsaKeyPriv, k.Error = ParseECPrivateKeyFromPEM(privateKey)
	return
}

func WithECPrivateKey(privateKey *ecdsa.PrivateKey) (k Key) {
	k.ecdsaKeyPriv = privateKey
	return
}
