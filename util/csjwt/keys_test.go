package csjwt

import (
	"crypto/ecdsa"
	"crypto/rsa"
	"crypto/x509"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

var _ fmt.GoStringer = (*Key)(nil)
var _ fmt.Stringer = (*Key)(nil)

func TestNewKeyFunc(t *testing.T) {
	t.Parallel()
	tests := []struct {
		s       Signer
		key     Key
		token   Token
		wantKey Key
		wantErr error
	}{
		{nil, Key{Error: errors.New("idx1")}, Token{}, Key{}, errors.New("idx1")},
		{
			&SigningMethodHMAC{Name: "Rost"},
			WithPasswordRandom(),
			NewToken(nil),
			Key{},
			ErrTokenUnverifiable,
		},
		{
			NewSigningMethodHS256(),
			WithPassword([]byte(`123456`)),
			Token{
				Header: &Head{Algorithm: HS256},
			},
			WithPassword([]byte(`123456`)),
			nil,
		},
	}
	for i, test := range tests {
		haveKey, haveErr := NewKeyFunc(test.s, test.key)(test.token)
		assert.Exactly(t, test.wantKey, haveKey, "Index %d", i)
		if test.wantErr != nil {
			assert.EqualError(t, haveErr, test.wantErr.Error(), "Index %d", i)
			continue
		}
		assert.NoError(t, haveErr, "Index %d", i)
	}
}

func TestKeyParsing(t *testing.T) {
	t.Parallel()

	badKey := []byte("This is a bad key")
	tests := []struct {
		key     Key
		wantAlg string
		wantErr error
		wantKey interface{}
	}{
		{WithPassword(badKey), HS, nil, []byte{}},
		{WithPassword(nil), "", errKeyEmptyPassword, nil},
		{WithPasswordFromFile("test/hmacTestKey"), HS, nil, []byte{}},
		{WithPasswordFromFile("test/hmacTestKeyNONEXIST"), "", errors.New("open test/hmacTestKeyNONEXIST: no such file or directory"), nil},

		{WithRSAPrivateKey(new(rsa.PrivateKey)), RS, nil, new(rsa.PrivateKey)},
		{WithRSAPrivateKeyFromFile("test/sample_keyOFF"), "", errors.New("open test/sample_keyOFF: no such file or directory"), nil},
		{WithRSAPrivateKeyFromFile("test/sample_key"), RS, nil, new(rsa.PrivateKey)},
		{WithRSAPrivateKeyFromFile("test/test_rsa", []byte("cccamp")), RS, nil, new(rsa.PrivateKey)},
		{WithRSAPrivateKeyFromFile("test/test_rsa", []byte("cCcamp")), "", x509.IncorrectPasswordError, nil},
		{WithRSAPrivateKeyFromFile("test/test_rsa"), "", errKeyMissingPassword, nil},
		{WithRSAPrivateKeyFromFile("test/sample_key.pub"), "", errors.New("asn1: structure error: tags don't match (2 vs {class:0 tag:16 length:13 isCompound:true}) {optional:false explicit:false application:false defaultValue:<nil> tag:<nil> stringType:0 timeType:0 set:false omitEmpty:false} int @2"), nil},
		{WithRSAPrivateKeyFromPEM(badKey), "", errKeyMustBePEMEncoded, nil},

		{WithRSAPublicKeyFromFile("test/sample_key.pubOFF"), "", errors.New("open test/sample_key.pubOFF: no such file or directory with file test/sample_key.pubOFF"), nil},
		{WithRSAPublicKeyFromFile("test/sample_key.pub"), RS, nil, new(rsa.PublicKey)},
		{WithRSAPublicKeyFromFile("test/sample_key"), "", errors.New("asn1: structure error: tags don't match (16 vs {class:0 tag:2 length:1 isCompound:false}) {optional:false explicit:false application:false defaultValue:<nil> tag:<nil> stringType:0 timeType:0 set:false omitEmpty:false} tbsCertificate @2"), nil},
		{WithRSAPublicKeyFromPEM(badKey), "", errKeyMustBePEMEncoded, nil},
		{WithRSAPublicKey(new(rsa.PublicKey)), RS, nil, new(rsa.PublicKey)},

		{WithECPublicKeyFromPEM(badKey), "", errKeyMustBePEMEncoded, nil}, // 17
		{WithECPublicKey(new(ecdsa.PublicKey)), ES, nil, new(ecdsa.PublicKey)},
		{WithECPublicKeyFromFile("test/nothingecdsa"), "", errors.New("open test/nothingecdsa: no such file or directory"), nil},
		{WithECPublicKeyFromFile("test/ec512-public.pem"), ES, nil, new(ecdsa.PublicKey)},

		{WithECPrivateKeyFromPEM(badKey), "", errKeyMustBePEMEncoded, nil},
		{WithECPrivateKey(new(ecdsa.PrivateKey)), ES, nil, new(ecdsa.PrivateKey)},
		{WithECPrivateKeyFromFile("test/nothingecdsa"), "", errors.New("open test/nothingecdsa: no such file or directory"), nil},
		{WithECPrivateKeyFromFile("test/ec512-private.pem"), ES, nil, new(ecdsa.PrivateKey)},
	}
	for i, test := range tests {

		if test.wantErr != nil {
			assert.EqualError(t, test.key.Error, test.wantErr.Error(), "Index %d", i)
		} else {
			assert.NoError(t, test.key.Error, "Index %d", i)
		}

		switch test.wantKey.(type) {
		case *rsa.PrivateKey:
			assert.NotNil(t, test.key.rsaKeyPriv, "Index %d", i)
		case *rsa.PublicKey:
			assert.NotNil(t, test.key.rsaKeyPub, "Index %d", i)
		case []byte:
			assert.NotNil(t, test.key.hmacPassword, "Index %d", i)
		case *ecdsa.PublicKey:
			assert.NotNil(t, test.key.ecdsaKeyPub, "Index %d", i)
		case *ecdsa.PrivateKey:
			assert.NotNil(t, test.key.ecdsaKeyPriv, "Index %d", i)
		case nil:
			assert.True(t, test.key.IsEmpty(), "Index %d", i)
		default:
			t.Fatal("Dude! You missed an entry in this list!")
		}
		assert.Exactly(t, test.wantAlg, test.key.Algorithm(), "Index %d", i)
		assert.Exactly(t, goStringTpl, fmt.Sprintf("%#v", test.key))
		assert.Exactly(t, goStringTpl, fmt.Sprintf("%v", test.key))
		assert.Exactly(t, goStringTpl, fmt.Sprintf("%s", test.key))
	}
}

func TestKeyWithPasswordRandom(t *testing.T) {
	t.Parallel()
	key := WithPasswordRandom()
	assert.Len(t, key.hmacPassword, randomPasswordLenght)
	if len(fmt.Sprintf("%x", key.hmacPassword)) < randomPasswordLenght {
		t.Fatalf("Generated password is too short: %x", key.hmacPassword)
	}
}

func TestKeyWithRSAGenerator(t *testing.T) {
	t.Parallel()
	key := WithRSAGenerated()
	assert.Exactly(t, RS, key.Algorithm())
}
