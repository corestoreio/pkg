package csjwt_test

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"runtime"
	"runtime/debug"
	"testing"
	"time"

	"github.com/corestoreio/csfw/storage/text"
	"github.com/corestoreio/csfw/util/cserr"
	"github.com/corestoreio/csfw/util/csjwt"
	"github.com/corestoreio/csfw/util/csjwt/jwtclaim"
	"github.com/stretchr/testify/assert"
)

var (
	defaultKeyFunc csjwt.Keyfunc = func(t csjwt.Token) (csjwt.Key, error) {
		return csjwt.WithRSAPublicKeyFromFile("test/sample_key.pub"), nil
	}
	emptyKeyFunc csjwt.Keyfunc = func(t csjwt.Token) (csjwt.Key, error) { return csjwt.Key{}, nil }
	errorKeyFunc csjwt.Keyfunc = func(t csjwt.Token) (csjwt.Key, error) { return csjwt.Key{}, fmt.Errorf("error loading key") }
	nilKeyFunc   csjwt.Keyfunc = nil
)

var jwtTestData = []struct {
	name        string
	tokenString []byte
	keyfunc     csjwt.Keyfunc
	claims      jwtclaim.Map
	valid       bool
	wantErr     error
	parser      *csjwt.Verification
}{
	{
		"basic",
		[]byte("eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiJ9.eyJmb28iOiJiYXIifQ.FhkiHkoESI_cG3NPigFrxEk9Z60_oXrOT2vGm9Pn6RDgYNovYORQmmA0zs1AoAOf09ly2Nx2YAg6ABqAYga1AcMFkJljwxTT5fYphTuqpWdy4BELeSYJx5Ty2gmr8e7RonuUztrdD5WfPqLKMm1Ozp_T6zALpRmwTIW0QPnaBXaQD90FplAg46Iy1UlDKr-Eupy0i5SLch5Q-p2ZpaL_5fnTIUDlxC3pWhJTyx_71qDI-mAA_5lE_VdroOeflG56sSmDxopPEG3bFlSu1eowyBfxtu0_CuVd-M42RU75Zc4Gsj6uV77MBtbMrf4_7M_NUTSgoIF3fRqxrj0NzihIBg"),
		defaultKeyFunc,
		jwtclaim.Map{"foo": "bar"},
		true,
		nil,
		csjwt.NewVerification(csjwt.NewSigningMethodRS256()),
	},
	{
		"basic expired",
		nil, // autogen
		defaultKeyFunc,
		jwtclaim.Map{"foo": "bar", "exp": float64(time.Now().Unix() - 100)},
		false,
		jwtclaim.ErrValidationExpired,
		csjwt.NewVerification(),
	},
	{
		"basic nbf",
		nil, // autogen
		defaultKeyFunc,
		jwtclaim.Map{"foo": "bar", "nbf": float64(time.Now().Unix() + 100)},
		false,
		jwtclaim.ErrValidationNotValidYet,
		csjwt.NewVerification(),
	},
	{
		"expired and nbf",
		nil, // autogen
		defaultKeyFunc,
		jwtclaim.Map{"foo": "bar", "nbf": float64(time.Now().Unix() + 100), "exp": float64(time.Now().Unix() - 100)},
		false,
		cserr.NewMultiErr(jwtclaim.ErrValidationNotValidYet, jwtclaim.ErrValidationExpired),
		csjwt.NewVerification(),
	},
	{
		"basic invalid",
		[]byte("eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiJ9.eyJmb28iOiJiYXIifQ.EhkiHkoESI_cG3NPigFrxEk9Z60_oXrOT2vGm9Pn6RDgYNovYORQmmA0zs1AoAOf09ly2Nx2YAg6ABqAYga1AcMFkJljwxTT5fYphTuqpWdy4BELeSYJx5Ty2gmr8e7RonuUztrdD5WfPqLKMm1Ozp_T6zALpRmwTIW0QPnaBXaQD90FplAg46Iy1UlDKr-Eupy0i5SLch5Q-p2ZpaL_5fnTIUDlxC3pWhJTyx_71qDI-mAA_5lE_VdroOeflG56sSmDxopPEG3bFlSu1eowyBfxtu0_CuVd-M42RU75Zc4Gsj6uV77MBtbMrf4_7M_NUTSgoIF3fRqxrj0NzihIBg"),
		defaultKeyFunc,
		jwtclaim.Map{"foo": "bar"},
		false,
		csjwt.ErrSignatureInvalid,
		csjwt.NewVerification(csjwt.NewSigningMethodRS256()),
	},
	{
		"basic nokeyfunc",
		[]byte("eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiJ9.eyJmb28iOiJiYXIifQ.FhkiHkoESI_cG3NPigFrxEk9Z60_oXrOT2vGm9Pn6RDgYNovYORQmmA0zs1AoAOf09ly2Nx2YAg6ABqAYga1AcMFkJljwxTT5fYphTuqpWdy4BELeSYJx5Ty2gmr8e7RonuUztrdD5WfPqLKMm1Ozp_T6zALpRmwTIW0QPnaBXaQD90FplAg46Iy1UlDKr-Eupy0i5SLch5Q-p2ZpaL_5fnTIUDlxC3pWhJTyx_71qDI-mAA_5lE_VdroOeflG56sSmDxopPEG3bFlSu1eowyBfxtu0_CuVd-M42RU75Zc4Gsj6uV77MBtbMrf4_7M_NUTSgoIF3fRqxrj0NzihIBg"),
		nilKeyFunc,
		jwtclaim.Map{"foo": "bar"},
		false,
		errors.New("[csjwt] Missing KeyFunc"),
		csjwt.NewVerification(csjwt.NewSigningMethodRS256()),
	},
	{
		"basic nokey",
		[]byte("eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiJ9.eyJmb28iOiJiYXIifQ.FhkiHkoESI_cG3NPigFrxEk9Z60_oXrOT2vGm9Pn6RDgYNovYORQmmA0zs1AoAOf09ly2Nx2YAg6ABqAYga1AcMFkJljwxTT5fYphTuqpWdy4BELeSYJx5Ty2gmr8e7RonuUztrdD5WfPqLKMm1Ozp_T6zALpRmwTIW0QPnaBXaQD90FplAg46Iy1UlDKr-Eupy0i5SLch5Q-p2ZpaL_5fnTIUDlxC3pWhJTyx_71qDI-mAA_5lE_VdroOeflG56sSmDxopPEG3bFlSu1eowyBfxtu0_CuVd-M42RU75Zc4Gsj6uV77MBtbMrf4_7M_NUTSgoIF3fRqxrj0NzihIBg"),
		emptyKeyFunc,
		jwtclaim.Map{"foo": "bar"},
		false,
		csjwt.ErrSignatureInvalid,
		csjwt.NewVerification(csjwt.NewSigningMethodRS256()),
	},
	{
		"basic errorkey",
		[]byte("eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiJ9.eyJmb28iOiJiYXIifQ.FhkiHkoESI_cG3NPigFrxEk9Z60_oXrOT2vGm9Pn6RDgYNovYORQmmA0zs1AoAOf09ly2Nx2YAg6ABqAYga1AcMFkJljwxTT5fYphTuqpWdy4BELeSYJx5Ty2gmr8e7RonuUztrdD5WfPqLKMm1Ozp_T6zALpRmwTIW0QPnaBXaQD90FplAg46Iy1UlDKr-Eupy0i5SLch5Q-p2ZpaL_5fnTIUDlxC3pWhJTyx_71qDI-mAA_5lE_VdroOeflG56sSmDxopPEG3bFlSu1eowyBfxtu0_CuVd-M42RU75Zc4Gsj6uV77MBtbMrf4_7M_NUTSgoIF3fRqxrj0NzihIBg"),
		errorKeyFunc,
		jwtclaim.Map{"foo": "bar"},
		false,
		csjwt.ErrTokenUnverifiable,
		csjwt.NewVerification(csjwt.NewSigningMethodRS256()),
	},
	{
		"invalid signing method",
		nil, // token gets generated with RS256 method
		defaultKeyFunc,
		map[string]interface{}{"foo": "bar"},
		false,
		errors.New("[csjwt] Algorithm \"RS256\" not found in method list \"HS256\""),
		csjwt.NewVerification(), // default signing method is HS256
	},
	{
		"valid signing method",
		nil,
		defaultKeyFunc,
		map[string]interface{}{"foo": "bar"},
		true,
		nil,
		csjwt.NewVerification(csjwt.NewSigningMethodES256(), csjwt.NewSigningMethodRS256()),
	},
}

func makeSample(c jwtclaim.Map) []byte {
	token := csjwt.NewToken(c)
	s, err := token.SignedString(csjwt.NewSigningMethodRS256(), csjwt.WithRSAPrivateKeyFromFile("test/sample_key"))
	if err != nil {
		panic(err)
	}
	return s
}

func TestParseWithMap(t *testing.T) {
	t.Parallel()
	for _, data := range jwtTestData {
		if len(data.tokenString) == 0 {
			data.tokenString = makeSample(data.claims)
		}

		token, err := data.parser.Parse(csjwt.NewToken(&jwtclaim.Map{}), data.tokenString, data.keyfunc)

		if !reflect.DeepEqual(&data.claims, token.Claims) {
			t.Errorf("[%v] Claims mismatch. Expecting: %v  Got: %v", data.name, data.claims, token.Claims)
		}
		if data.valid && err != nil {
			t.Errorf("[%v] Error while verifying token: %T:%v", data.name, err, err)
		}
		if !data.valid && err == nil {
			t.Errorf("[%v] Invalid token passed validation", data.name)
		}
		if data.wantErr != nil {
			if err == nil {
				t.Errorf("[%v] Expecting error.  Didn't get one.", data.name)
			} else {

				if !cserr.Contains(err, data.wantErr) {
					t.Errorf("[%v] Errors don't match expectation:\n@%#v@ != |%#v|\n", data.name, err, data.wantErr)
				}
			}
		}
		if data.valid && len(token.Signature) == 0 {
			t.Errorf("[%v] Signature is left unpopulated after parsing", data.name)
		}
	}
}

func TestParseFromRequest(t *testing.T) {
	t.Parallel()
	// Bearer token request
	for _, data := range jwtTestData {

		if len(data.tokenString) == 0 {
			data.tokenString = makeSample(data.claims)
		}

		r, _ := http.NewRequest("GET", "/", nil)
		r.Header.Set("Authorization", fmt.Sprintf("Bearer %s", data.tokenString))
		token, err := data.parser.ParseFromRequest(csjwt.NewToken(&jwtclaim.Map{}), data.keyfunc, r)

		if token.Raw == nil {
			t.Errorf("[%v] Token was not found: %v", data.name, err)
			continue
		}
		if !reflect.DeepEqual(&data.claims, token.Claims) {
			t.Errorf("[%v] Claims mismatch. Expecting: %v  Got: %v", data.name, data.claims, token.Claims)
		}
		if data.valid && err != nil {
			t.Errorf("[%v] Error while verifying token: %v", data.name, err)
		}
		if !data.valid && err == nil {
			t.Errorf("[%v] Invalid token passed validation", data.name)
		}
	}
}

func TestParseFromRequestComplex(t *testing.T) {
	t.Parallel()

	key := csjwt.WithPassword([]byte(`csjwt.SigningMethodHS512!`))
	clm := jwtclaim.Map{
		"foo":               "bar",
		"user_id":           "hello_gophers",
		"cart_items":        "234234,12;34234,34;234234,1;123123,12;234234,12;34234,34;234234,1;123123,12;234234,12;34234,34;234234,1;123123,12;",
		"last_viewed_items": "234234,12;34234,34;234234,1;123123,12;234234,12;34234,34;234234,1;123123,12;234234,12;34234,34;234234,1;123123,12;",
		"requested_price":   float64(3.141592 * 2.718281 / 3),
		"checkout_step":     float64(3),
		"payment_valid":     true,
		"exp":               float64(time.Now().Add(time.Hour * 72).Unix()),
	}
	sm512 := csjwt.NewSigningMethodHS512()
	token := csjwt.NewToken(clm)
	tokenString, err := token.SignedString(sm512, key)
	if err != nil {
		t.Fatal(err)
	}

	r, _ := http.NewRequest("GET", "/", nil)
	r.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tokenString))

	var newClaim = make(jwtclaim.Map)
	rToken, err := csjwt.NewVerification(sm512).ParseFromRequest(csjwt.NewToken(&newClaim), func(t csjwt.Token) (csjwt.Key, error) {
		if have, want := t.Alg(), sm512.Alg(); have != want {
			return csjwt.Key{}, fmt.Errorf("Have: %s Want: %s", have, want)
		}
		return key, nil
	}, r)
	if err != nil {
		t.Fatal(err)
	}
	assert.Exactly(t, rToken.Claims, &clm)
}

func TestParseWithClaimsBearerInHeader(t *testing.T) {
	t.Parallel()
	token := text.Chars(`BEaRER `)
	token = append(token, jwtTestData[0].tokenString...)

	haveToken, haveErr := csjwt.NewVerification().Parse(csjwt.NewToken(&jwtclaim.Map{}), token, nil)
	assert.NotNil(t, haveToken)
	assert.NotNil(t, haveToken.Claims)
	assert.Exactly(t, haveToken.Raw, token)
	assert.EqualError(t, haveErr, `[csjwt] tokenstring should not contain 'bearer '`)
	assert.True(t, bytes.Contains(haveToken.Raw, token))
}

func TestParseWithClaimInvalidSegments(t *testing.T) {
	t.Parallel()
	token, err := csjwt.NewVerification().Parse(csjwt.NewToken(nil), []byte(`hello.gopher`), nil)
	assert.False(t, token.Valid)
	assert.EqualError(t, err, `[csjwt] token contains an invalid number of segments`)
}

func TestSplitForVerify(t *testing.T) {
	t.Parallel()
	tests := []struct {
		rawToken      []byte
		signingString []byte
		signature     []byte
		wantErr       error
	}{
		{
			[]byte(`Hello.World.Gophers`),
			[]byte(`Hello.World`),
			[]byte(`Gophers`),
			nil,
		},
		{
			[]byte(`Hello.WorldGophers`),
			nil,
			nil,
			errors.New("[csjwt] token contains an invalid number of segments"),
		},
		{
			[]byte(`Hello.World.Gop.hers`),
			nil,
			nil,
			errors.New("[csjwt] token contains an invalid number of segments"),
		},
	}
	for _, test := range tests {
		haveSS, haveSig, haveErr := csjwt.SplitForVerify(test.rawToken)
		if test.wantErr != nil {
			assert.EqualError(t, haveErr, test.wantErr.Error())
		} else {
			assert.NoError(t, haveErr)
		}
		assert.Exactly(t, test.signingString, haveSS)
		assert.Exactly(t, test.signature, haveSig)
	}
}

func benchmarkSigning(b *testing.B, method csjwt.Signer, key csjwt.Key) {
	sc := &jwtclaim.Standard{}
	b.ResetTimer()
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			if _, err := csjwt.NewToken(sc).SignedString(method, key); err != nil {
				b.Fatal(err)
			}
		}
	})
}

func benchmarkMethodVerify(b *testing.B, method csjwt.Signer, signingString []byte, signature []byte, key csjwt.Key) {
	b.ResetTimer()
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			if err := method.Verify(signingString, signature, key); err != nil {
				b.Fatal(err)
			}
		}
	})
}

func BenchmarkParseFromRequest_HS256(b *testing.B) {
	key := csjwt.WithPassword([]byte(`csjwt.SigningMethodHS256!`))
	hmacFast, err := csjwt.NewHMACFast256(key)
	if err != nil {
		b.Fatal(err)
	}
	benchmarkParseFromRequest(
		b,
		hmacFast, // csjwt.SigningMethodHS256,
		key,
		func(t csjwt.Token) (csjwt.Key, error) {
			if have, want := t.Alg(), hmacFast.Alg(); have != want {
				return csjwt.Key{}, fmt.Errorf("Have: %s Want: %s", have, want)
			}
			return key, nil
		},
	)
}
func BenchmarkParseFromRequest_HS384(b *testing.B) {
	key := csjwt.WithPassword([]byte(`csjwt.SigningMethodHS384!`))
	hmacFast, err := csjwt.NewHMACFast384(key)
	if err != nil {
		b.Fatal(err)
	}
	benchmarkParseFromRequest(
		b,
		hmacFast, // csjwt.SigningMethodHS384,
		key,
		func(t csjwt.Token) (csjwt.Key, error) {
			if have, want := t.Alg(), hmacFast.Alg(); have != want {
				return csjwt.Key{}, fmt.Errorf("Have: %s Want: %s", have, want)
			}
			return key, nil
		},
	)
}
func BenchmarkParseFromRequest_HS512(b *testing.B) {
	key := csjwt.WithPassword([]byte(`csjwt.SigningMethodHS512!`))
	hmacFast, err := csjwt.NewHMACFast512(key)
	if err != nil {
		b.Fatal(err)
	}
	benchmarkParseFromRequest(
		b,
		hmacFast, // csjwt.SigningMethodHS512,
		key,
		func(t csjwt.Token) (csjwt.Key, error) {
			if have, want := t.Alg(), hmacFast.Alg(); have != want {
				return csjwt.Key{}, fmt.Errorf("Have: %s Want: %s", have, want)
			}
			return key, nil
		},
	)
}

type ShoppingCartClaim struct {
	*jwtclaim.Standard
	CartPID       []int
	LastViewedPID []int
	RequestPrice  float64
	CheckoutStep  uint8
	PaymentValid  bool
}

func benchmarkParseFromRequest(b *testing.B, sm csjwt.Signer, key csjwt.Key, keyFunc csjwt.Keyfunc) {

	clm := &ShoppingCartClaim{
		Standard: &jwtclaim.Standard{
			ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		CartPID:       []int{12345, 6789, 2345, 3456, 7564, 45678, 5678578, 345234, 2345234},
		LastViewedPID: []int{6456, 3453, 45345, 234235, 345345, 645646, 567567, 345635, 85689, 5678, 5674567, 345635, 4356, 245645},
		RequestPrice:  2.718281 * 3.141592,
		CheckoutStep:  3,
		PaymentValid:  true,
	}

	token := csjwt.NewToken(clm)
	tokenString, err := token.SignedString(sm, key)
	if err != nil {
		b.Fatal(err)
	}

	r, _ := http.NewRequest("GET", "/", nil)
	r.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tokenString))

	veri := csjwt.NewVerification(sm)

	mc := &ShoppingCartClaim{Standard: new(jwtclaim.Standard)}

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {

		rToken, err := veri.ParseFromRequest(csjwt.NewToken(mc), keyFunc, r)
		if err != nil {
			b.Fatal(err)
		}
		if !rToken.Valid {
			b.Fatalf("Token not valid: %#v", rToken)
		}
	}

	if have, want := mc.ExpiresAt, clm.ExpiresAt; have != want {
		b.Fatalf("Mismatch of claims: Have %d Want %d", have, want)
	}
	//b.Log("GC Pause:", gcPause())
}

func gcPause() time.Duration {
	runtime.GC()
	var stats debug.GCStats
	debug.ReadGCStats(&stats)
	return stats.Pause[0]
}
