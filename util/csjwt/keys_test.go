package csjwt

import (
	"crypto/ecdsa"
	"crypto/rsa"
	"fmt"
	"testing"

	"github.com/corestoreio/csfw/util/errors"
	"github.com/stretchr/testify/assert"
)

var _ fmt.GoStringer = (*Key)(nil)
var _ fmt.Stringer = (*Key)(nil)

func TestNewKeyFunc(t *testing.T) {

	tests := []struct {
		s          Signer
		key        Key
		token      Token
		wantKey    Key
		wantErrBhf errors.BehaviourFunc
	}{
		{nil, Key{Error: errors.NewAlreadyClosedf("idx1")}, Token{}, Key{}, errors.IsAlreadyClosed},
		{
			&SigningMethodHMAC{Name: "Rost"},
			WithPasswordRandom(),
			NewToken(nil),
			Key{},
			errors.IsNotValid,
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
		haveKey, haveErr := NewKeyFunc(test.s, test.key)(&test.token)
		assert.Exactly(t, test.wantKey, haveKey, "Index %d", i)
		if test.wantErrBhf != nil {
			assert.True(t, test.wantErrBhf(haveErr), "Index %d => %s", i, haveErr)
			continue
		}
		assert.NoError(t, haveErr, "Index %d", i)
	}
}

func TestKeyParsing(t *testing.T) {

	badKey := []byte("This is a bad key")
	tests := []struct {
		key        Key
		wantAlg    string
		wantErrBhf errors.BehaviourFunc
		wantKey    interface{}
	}{
		{WithPassword(badKey), HS, nil, []byte{}},
		{WithPassword(nil), "", errors.IsEmpty, nil},
		{WithPasswordFromFile("test/hmacTestKey"), HS, nil, []byte{}},
		{WithPasswordFromFile("test/hmacTestKeyNONEXIST"), "", errors.IsNotValid, nil},

		{WithRSAPrivateKey(new(rsa.PrivateKey)), RS, nil, new(rsa.PrivateKey)},
		{WithRSAPrivateKeyFromFile("test/sample_keyOFF"), "", errors.IsNotValid, nil},
		{WithRSAPrivateKeyFromFile("test/sample_key"), RS, nil, new(rsa.PrivateKey)},
		{WithRSAPrivateKeyFromFile("test/test_rsa", []byte("cccamp")), RS, nil, new(rsa.PrivateKey)},
		{WithRSAPrivateKeyFromFile("test/test_rsa", []byte("cCcamp")), "", errors.IsNotValid, nil},
		{WithRSAPrivateKeyFromFile("test/test_rsa"), "", errors.IsEmpty, nil},
		{WithRSAPrivateKeyFromFile("test/sample_key.pub"), "", errors.IsNotValid, nil},
		{WithRSAPrivateKeyFromPEM(badKey), "", errors.IsNotSupported, nil},

		{WithRSAPublicKeyFromFile("test/sample_key.pubOFF"), "", errors.IsNotValid, nil},
		{WithRSAPublicKeyFromFile("test/sample_key.pub"), RS, nil, new(rsa.PublicKey)},
		{WithRSAPublicKeyFromFile("test/sample_key"), "", errors.IsNotValid, nil},
		{WithRSAPublicKeyFromPEM(badKey), "", errors.IsNotSupported, nil},
		{WithRSAPublicKey(new(rsa.PublicKey)), RS, nil, new(rsa.PublicKey)},

		{WithECPublicKeyFromPEM(badKey), "", errors.IsNotSupported, nil}, // 17
		{WithECPublicKey(new(ecdsa.PublicKey)), ES, nil, new(ecdsa.PublicKey)},
		{WithECPublicKeyFromFile("test/nothingecdsa"), "", errors.IsNotValid, nil},
		{WithECPublicKeyFromFile("test/ec512-public.pem"), ES, nil, new(ecdsa.PublicKey)},

		{WithECPrivateKeyFromPEM(badKey), "", errors.IsNotSupported, nil},
		{WithECPrivateKey(new(ecdsa.PrivateKey)), ES, nil, new(ecdsa.PrivateKey)},
		{WithECPrivateKeyFromFile("test/nothingecdsa"), "", errors.IsNotValid, nil},
		{WithECPrivateKeyFromFile("test/ec512-private.pem"), ES, nil, new(ecdsa.PrivateKey)},
	}
	for i, test := range tests {

		if test.wantErrBhf != nil {
			assert.True(t, test.wantErrBhf(test.key.Error), "Index %d => %s\n", i, errors.PrintLoc(test.key.Error))
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

	key := WithPasswordRandom()
	assert.Len(t, key.hmacPassword, randomPasswordLenght)
	if len(fmt.Sprintf("%x", key.hmacPassword)) < randomPasswordLenght {
		t.Fatalf("Generated password is too short: %x", key.hmacPassword)
	}
}

func TestKeyWithRSAGenerator(t *testing.T) {

	key := WithRSAGenerated()
	assert.Exactly(t, RS, key.Algorithm())
}
