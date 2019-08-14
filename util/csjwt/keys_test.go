package csjwt

import (
	"crypto/ecdsa"
	"crypto/rsa"
	"fmt"
	"testing"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/pkg/util/assert"
)

var _ fmt.GoStringer = (*Key)(nil)
var _ fmt.Stringer = (*Key)(nil)

func TestNewKeyFunc(t *testing.T) {
	tests := []struct {
		s           Signer
		key         Key
		token       *Token
		wantKey     Key
		wantErrKind errors.Kind
	}{
		{nil, Key{Error: errors.AlreadyClosed.Newf("idx1")}, &Token{}, Key{}, errors.AlreadyClosed},
		{
			&SigningMethodHMAC{Name: "Rost"},
			WithPasswordRandom(),
			NewToken(nil),
			Key{},
			errors.NotValid,
		},
		{
			NewSigningMethodHS256(),
			WithPassword([]byte(`123456`)),
			&Token{
				Header: &Head{Algorithm: HS256},
			},
			WithPassword([]byte(`123456`)),
			errors.NoKind,
		},
	}
	for i, test := range tests {
		haveKey, haveErr := NewKeyFunc(test.s, test.key)(test.token)
		assert.Exactly(t, test.wantKey, haveKey, "Index %d", i)
		if !test.wantErrKind.Empty() {
			assert.True(t, test.wantErrKind.Match(haveErr), "Index %d => %s", i, haveErr)
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
		wantErrBhf errors.Kind
		wantKey    interface{}
	}{
		{WithPassword(badKey), HS, errors.NoKind, []byte{}},
		{WithPassword(nil), "", errors.Empty, nil},
		{WithPasswordFromFile("test/hmacTestKey"), HS, errors.NoKind, []byte{}},
		{WithPasswordFromFile("test/hmacTestKeyNONEXIST"), "", errors.NotValid, nil},

		{WithRSAPrivateKey(new(rsa.PrivateKey)), RS, errors.NoKind, new(rsa.PrivateKey)},
		{WithRSAPrivateKeyFromFile("test/sample_keyOFF"), "", errors.NotValid, nil},
		{WithRSAPrivateKeyFromFile("test/sample_key"), RS, errors.NoKind, new(rsa.PrivateKey)},
		{WithRSAPrivateKeyFromFile("test/test_rsa", []byte("cccamp")), RS, errors.NoKind, new(rsa.PrivateKey)},
		{WithRSAPrivateKeyFromFile("test/test_rsa", []byte("cCcamp")), "", errors.NotValid, nil},
		{WithRSAPrivateKeyFromFile("test/test_rsa"), "", errors.Empty, nil},
		{WithRSAPrivateKeyFromFile("test/sample_key.pub"), "", errors.NotValid, nil},
		{WithRSAPrivateKeyFromPEM(badKey), "", errors.NotSupported, nil},

		{WithRSAPublicKeyFromFile("test/sample_key.pubOFF"), "", errors.NotValid, nil},
		{WithRSAPublicKeyFromFile("test/sample_key.pub"), RS, errors.NoKind, new(rsa.PublicKey)},
		{WithRSAPublicKeyFromFile("test/sample_key"), "", errors.NotValid, nil},
		{WithRSAPublicKeyFromPEM(badKey), "", errors.NotSupported, nil},
		{WithRSAPublicKey(new(rsa.PublicKey)), RS, errors.NoKind, new(rsa.PublicKey)},

		{WithECPublicKeyFromPEM(badKey), "", errors.NotSupported, nil}, // 17
		{WithECPublicKey(new(ecdsa.PublicKey)), ES, errors.NoKind, new(ecdsa.PublicKey)},
		{WithECPublicKeyFromFile("test/nothingecdsa"), "", errors.NotValid, nil},
		{WithECPublicKeyFromFile("test/ec512-public.pem"), ES, errors.NoKind, new(ecdsa.PublicKey)},

		{WithECPrivateKeyFromPEM(badKey), "", errors.NotSupported, nil},
		{WithECPrivateKey(new(ecdsa.PrivateKey)), ES, errors.NoKind, new(ecdsa.PrivateKey)},
		{WithECPrivateKeyFromFile("test/nothingecdsa"), "", errors.NotValid, nil},
		{WithECPrivateKeyFromFile("test/ec512-private.pem"), ES, errors.NoKind, new(ecdsa.PrivateKey)},
	}
	for i, test := range tests {

		if !test.wantErrBhf.Empty() {
			assert.True(t, test.wantErrBhf.Match(test.key.Error), "Index %d => %+v\n", i, test.key.Error)
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
	assert.Len(t, key.hmacPassword, randomPasswordLength)
	if len(fmt.Sprintf("%x", key.hmacPassword)) < randomPasswordLength {
		t.Fatalf("Generated password is too short: %x", key.hmacPassword)
	}
}

func TestKeyWithRSAGenerator(t *testing.T) {
	key := WithRSAGenerated()
	assert.Exactly(t, RS, key.Algorithm())
}
