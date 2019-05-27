package csjwt_test

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"testing"
	"time"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/pkg/util/assert"
	"github.com/corestoreio/pkg/util/conv"
	"github.com/corestoreio/pkg/util/csjwt"
	"github.com/corestoreio/pkg/util/csjwt/jwtclaim"
)

var (
	defaultKeyFunc csjwt.Keyfunc = func(_ *csjwt.Token) (csjwt.Key, error) {
		return csjwt.WithRSAPublicKeyFromFile("test/sample_key.pub"), nil
	}
	emptyKeyFunc csjwt.Keyfunc = func(_ *csjwt.Token) (csjwt.Key, error) { return csjwt.Key{}, nil }
	errorKeyFunc csjwt.Keyfunc = func(_ *csjwt.Token) (csjwt.Key, error) { return csjwt.Key{}, fmt.Errorf("error loading key") }
	nilKeyFunc   csjwt.Keyfunc = nil
)

var jwtTestData = []struct {
	name        string
	tokenString []byte
	keyfunc     csjwt.Keyfunc
	claims      jwtclaim.Map
	valid       bool
	wantErrKind errors.Kind
	parser      *csjwt.Verification
}{
	{
		"basic",
		[]byte("eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiJ9.eyJmb28iOiJiYXIifQ.FhkiHkoESI_cG3NPigFrxEk9Z60_oXrOT2vGm9Pn6RDgYNovYORQmmA0zs1AoAOf09ly2Nx2YAg6ABqAYga1AcMFkJljwxTT5fYphTuqpWdy4BELeSYJx5Ty2gmr8e7RonuUztrdD5WfPqLKMm1Ozp_T6zALpRmwTIW0QPnaBXaQD90FplAg46Iy1UlDKr-Eupy0i5SLch5Q-p2ZpaL_5fnTIUDlxC3pWhJTyx_71qDI-mAA_5lE_VdroOeflG56sSmDxopPEG3bFlSu1eowyBfxtu0_CuVd-M42RU75Zc4Gsj6uV77MBtbMrf4_7M_NUTSgoIF3fRqxrj0NzihIBg"),
		defaultKeyFunc,
		jwtclaim.Map{"foo": "bar"},
		true,
		errors.NoKind,
		csjwt.NewVerification(csjwt.NewSigningMethodRS256()),
	},
	{
		"basic expired",
		nil, // autogen
		defaultKeyFunc,
		jwtclaim.Map{"foo": "bar", "exp": float64(time.Now().Unix() - 100)},
		false,
		errors.NotValid,
		csjwt.NewVerification(),
	},
	{
		"basic nbf",
		nil, // autogen
		defaultKeyFunc,
		jwtclaim.Map{"foo": "bar", "nbf": float64(time.Now().Unix() + 100)},
		false,
		errors.NotValid,
		csjwt.NewVerification(),
	},
	{
		"expired and nbf",
		nil, // autogen
		defaultKeyFunc,
		jwtclaim.Map{"foo": "bar", "nbf": float64(time.Now().Unix() + 100), "exp": float64(time.Now().Unix() - 100)},
		false,
		errors.NotValid,
		csjwt.NewVerification(),
	},
	{
		"basic invalid",
		[]byte("eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiJ9.eyJmb28iOiJiYXIifQ.EhkiHkoESI_cG3NPigFrxEk9Z60_oXrOT2vGm9Pn6RDgYNovYORQmmA0zs1AoAOf09ly2Nx2YAg6ABqAYga1AcMFkJljwxTT5fYphTuqpWdy4BELeSYJx5Ty2gmr8e7RonuUztrdD5WfPqLKMm1Ozp_T6zALpRmwTIW0QPnaBXaQD90FplAg46Iy1UlDKr-Eupy0i5SLch5Q-p2ZpaL_5fnTIUDlxC3pWhJTyx_71qDI-mAA_5lE_VdroOeflG56sSmDxopPEG3bFlSu1eowyBfxtu0_CuVd-M42RU75Zc4Gsj6uV77MBtbMrf4_7M_NUTSgoIF3fRqxrj0NzihIBg"),
		defaultKeyFunc,
		jwtclaim.Map{"foo": "bar"},
		false,
		errors.NotValid,
		csjwt.NewVerification(csjwt.NewSigningMethodRS256()),
	},
	{
		"basic nokeyfunc",
		[]byte("eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiJ9.eyJmb28iOiJiYXIifQ.FhkiHkoESI_cG3NPigFrxEk9Z60_oXrOT2vGm9Pn6RDgYNovYORQmmA0zs1AoAOf09ly2Nx2YAg6ABqAYga1AcMFkJljwxTT5fYphTuqpWdy4BELeSYJx5Ty2gmr8e7RonuUztrdD5WfPqLKMm1Ozp_T6zALpRmwTIW0QPnaBXaQD90FplAg46Iy1UlDKr-Eupy0i5SLch5Q-p2ZpaL_5fnTIUDlxC3pWhJTyx_71qDI-mAA_5lE_VdroOeflG56sSmDxopPEG3bFlSu1eowyBfxtu0_CuVd-M42RU75Zc4Gsj6uV77MBtbMrf4_7M_NUTSgoIF3fRqxrj0NzihIBg"),
		nilKeyFunc,
		jwtclaim.Map{"foo": "bar"},
		false,
		errors.Empty,
		csjwt.NewVerification(csjwt.NewSigningMethodRS256()),
	},
	{
		"basic nokey",
		[]byte("eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiJ9.eyJmb28iOiJiYXIifQ.FhkiHkoESI_cG3NPigFrxEk9Z60_oXrOT2vGm9Pn6RDgYNovYORQmmA0zs1AoAOf09ly2Nx2YAg6ABqAYga1AcMFkJljwxTT5fYphTuqpWdy4BELeSYJx5Ty2gmr8e7RonuUztrdD5WfPqLKMm1Ozp_T6zALpRmwTIW0QPnaBXaQD90FplAg46Iy1UlDKr-Eupy0i5SLch5Q-p2ZpaL_5fnTIUDlxC3pWhJTyx_71qDI-mAA_5lE_VdroOeflG56sSmDxopPEG3bFlSu1eowyBfxtu0_CuVd-M42RU75Zc4Gsj6uV77MBtbMrf4_7M_NUTSgoIF3fRqxrj0NzihIBg"),
		emptyKeyFunc,
		jwtclaim.Map{"foo": "bar"},
		false,
		errors.NotValid,
		csjwt.NewVerification(csjwt.NewSigningMethodRS256()),
	},
	{
		"basic errorkey",
		[]byte("eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiJ9.eyJmb28iOiJiYXIifQ.FhkiHkoESI_cG3NPigFrxEk9Z60_oXrOT2vGm9Pn6RDgYNovYORQmmA0zs1AoAOf09ly2Nx2YAg6ABqAYga1AcMFkJljwxTT5fYphTuqpWdy4BELeSYJx5Ty2gmr8e7RonuUztrdD5WfPqLKMm1Ozp_T6zALpRmwTIW0QPnaBXaQD90FplAg46Iy1UlDKr-Eupy0i5SLch5Q-p2ZpaL_5fnTIUDlxC3pWhJTyx_71qDI-mAA_5lE_VdroOeflG56sSmDxopPEG3bFlSu1eowyBfxtu0_CuVd-M42RU75Zc4Gsj6uV77MBtbMrf4_7M_NUTSgoIF3fRqxrj0NzihIBg"),
		errorKeyFunc,
		jwtclaim.Map{"foo": "bar"},
		false,
		errors.NotValid,
		csjwt.NewVerification(csjwt.NewSigningMethodRS256()),
	},
	{
		"invalid signing method",
		nil, // token gets generated with RS256 method
		defaultKeyFunc,
		map[string]interface{}{"foo": "bar"},
		false,
		errors.NotFound,
		csjwt.NewVerification(csjwt.NewSigningMethodHS256()),
	},
	{
		"valid signing method",
		nil,
		defaultKeyFunc,
		map[string]interface{}{"foo": "bar"},
		true,
		errors.NoKind,
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

func TestVerification_ParseFromRequest_WithMap(t *testing.T) {

	for _, data := range jwtTestData {
		if len(data.tokenString) == 0 {
			data.tokenString = makeSample(data.claims)
		}
		token := csjwt.NewToken(&jwtclaim.Map{})
		err := data.parser.Parse(token, data.tokenString, data.keyfunc)

		if !reflect.DeepEqual(&data.claims, token.Claims) {
			t.Errorf("[%v] Claims mismatch. Expecting: %v  Got: %v", data.name, data.claims, token.Claims)
		}
		if data.valid && err != nil {
			t.Errorf("[%v] Error while verifying token: %T:%v", data.name, err, err)
		}
		if !data.valid && err == nil {
			t.Errorf("[%v] Invalid token passed validation", data.name)
		}
		if !data.wantErrKind.Empty() {
			if err == nil {
				t.Errorf("[%v] Expecting error.  Didn't get one.", data.name)
			} else {

				if !data.wantErrKind.Match(err) {
					t.Errorf("[%v] Errors don't match expectation: %s\n", data.name, err)
				}
			}
		}
		if data.valid && len(token.Signature) == 0 {
			t.Errorf("[%v] Signature is left unpopulated after parsing", data.name)
		}
	}
}

func TestVerification_ParseFromRequest_LoopTestData(t *testing.T) {

	// Bearer token request
	for _, data := range jwtTestData {

		if len(data.tokenString) == 0 {
			data.tokenString = makeSample(data.claims)
		}

		r, _ := http.NewRequest("GET", "/", nil)
		r.Header.Set("Authorization", fmt.Sprintf("Bearer %s", data.tokenString))
		token := csjwt.NewToken(&jwtclaim.Map{})
		err := data.parser.ParseFromRequest(token, data.keyfunc, r)

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

func TestVerification_ParseFromRequest_Complex(t *testing.T) {

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
	rToken := csjwt.NewToken(&newClaim)
	err = csjwt.NewVerification(sm512).ParseFromRequest(rToken, func(t *csjwt.Token) (csjwt.Key, error) {
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

func TestVerification_Parse_BearerInHeader(t *testing.T) {

	token := []byte(`BEaRER `)
	token = append(token, jwtTestData[0].tokenString...)

	haveToken := csjwt.NewToken(&jwtclaim.Map{})
	haveErr := csjwt.NewVerification().Parse(haveToken, token, nil)
	assert.NotNil(t, haveToken)
	assert.NotNil(t, haveToken.Claims)
	assert.Exactly(t, haveToken.Raw, token)
	assert.True(t, errors.NotValid.Match(haveErr), "Error: %s", haveErr)
	assert.True(t, bytes.Contains(haveToken.Raw, token))
}

func TestVerification_Parse_InvalidSegments(t *testing.T) {
	token := csjwt.NewToken(nil)
	haveErr := csjwt.NewVerification().Parse(token, []byte(`hello.gopher`), nil)
	assert.False(t, token.Valid)
	assert.True(t, errors.NotValid.Match(haveErr), "Error: %s", haveErr)
}

func TestVerification_ParseFromRequest_Cookie(t *testing.T) {

	const cookieName = "store_bearer"
	token := makeSample(jwtclaim.Map{
		"where": "in the cookie dude!",
	})

	r, _ := http.NewRequest("GET", "/", nil)
	r.AddCookie(&http.Cookie{
		Name:  cookieName,
		Value: string(token),
	})

	rs256 := csjwt.NewSigningMethodRS256()
	vf := csjwt.NewVerification(rs256)
	vf.CookieName = cookieName

	pubKey := csjwt.WithRSAPublicKeyFromFile("test/sample_key.pub")

	haveToken := csjwt.NewToken(&jwtclaim.Map{})
	haveErr := vf.ParseFromRequest(haveToken, csjwt.NewKeyFunc(rs256, pubKey), r)
	if haveErr != nil {
		t.Fatalf("%+v", haveErr)
	}

	where, err := haveToken.Claims.Get(`where`)
	if err != nil {
		t.Fatalf("%+v", err)
	}
	assert.Exactly(t, `in the cookie dude!`, conv.ToString(where))
}

func TestVerification_ParseFromRequest_Form(t *testing.T) {

	token := makeSample(jwtclaim.Map{
		"where": "in the form dude!",
	})

	r, _ := http.NewRequest("GET", "/", nil)
	r.AddCookie(&http.Cookie{
		Name:  "not important",
		Value: "x1",
	})
	r.Form = url.Values{
		csjwt.HTTPFormInputName: []string{string(token)},
	}

	rs256 := csjwt.NewSigningMethodRS256()
	vf := csjwt.NewVerification(rs256)
	vf.CookieName = "unset"
	vf.FormInputName = csjwt.HTTPFormInputName

	pubKey := csjwt.WithRSAPublicKeyFromFile("test/sample_key.pub")
	haveToken := csjwt.NewToken(&jwtclaim.Map{})
	haveErr := vf.ParseFromRequest(haveToken, csjwt.NewKeyFunc(rs256, pubKey), r)
	if haveErr != nil {
		t.Fatalf("%+v", haveErr)
	}

	where, err := haveToken.Claims.Get(`where`)
	if err != nil {
		t.Fatalf("%+v", err)
	}
	assert.Exactly(t, `in the form dude!`, conv.ToString(where))
}

func TestVerification_ParseFromRequest_NoTokenInRequest(t *testing.T) {

	r, _ := http.NewRequest("GET", "/", nil)

	rs256 := csjwt.NewSigningMethodRS256()
	vf := csjwt.NewVerification(rs256)
	vf.CookieName = "unset"
	vf.FormInputName = csjwt.HTTPFormInputName

	pubKey := csjwt.WithRSAPublicKeyFromFile("test/sample_key.pub")
	haveToken := csjwt.NewToken(&jwtclaim.Map{})
	haveErr := vf.ParseFromRequest(haveToken, csjwt.NewKeyFunc(rs256, pubKey), r)
	assert.True(t, errors.NotFound.Match(haveErr), "Error: %s", haveErr)
	assert.Empty(t, haveToken.Raw)
	assert.False(t, haveToken.Valid)
}

func TestSplitForVerify(t *testing.T) {

	tests := []struct {
		rawToken      []byte
		signingString []byte
		signature     []byte
		wantErrKind   errors.Kind
	}{
		{
			[]byte(`Hello.World.Gophers`),
			[]byte(`Hello.World`),
			[]byte(`Gophers`),
			errors.NoKind,
		},
		{
			[]byte(`Hello.WorldGophers`),
			nil,
			nil,
			errors.NotValid,
		},
		{
			[]byte(`Hello.World.Gop.hers`),
			nil,
			nil,
			errors.NotValid,
		},
	}
	for i, test := range tests {
		haveSS, haveSig, haveErr := csjwt.SplitForVerify(test.rawToken)
		if !test.wantErrKind.Empty() {
			assert.True(t, test.wantErrKind.Match(haveErr), "Index %d => %s", i, haveErr)
		} else {
			assert.NoError(t, haveErr, "Index %d", i)
		}
		assert.Exactly(t, test.signingString, haveSS, "Index %d", i)
		assert.Exactly(t, test.signature, haveSig, "Index %d", i)
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
	hmacFast, err := csjwt.NewSigningMethodHS256Fast(key)
	if err != nil {
		b.Fatal(err)
	}
	benchmarkParseFromRequest(
		b,
		hmacFast, // csjwt.SigningMethodHS256,
		key,
		func(t *csjwt.Token) (csjwt.Key, error) {
			if have, want := t.Alg(), hmacFast.Alg(); have != want {
				return csjwt.Key{}, fmt.Errorf("Have: %s Want: %s", have, want)
			}
			return key, nil
		},
	)
}
func BenchmarkParseFromRequest_HS384(b *testing.B) {
	key := csjwt.WithPassword([]byte(`csjwt.SigningMethodHS384!`))
	hmacFast, err := csjwt.NewSigningMethodHS384Fast(key)
	if err != nil {
		b.Fatal(err)
	}
	benchmarkParseFromRequest(
		b,
		hmacFast, // csjwt.SigningMethodHS384,
		key,
		func(t *csjwt.Token) (csjwt.Key, error) {
			if have, want := t.Alg(), hmacFast.Alg(); have != want {
				return csjwt.Key{}, fmt.Errorf("Have: %s Want: %s", have, want)
			}
			return key, nil
		},
	)
}
func BenchmarkParseFromRequest_HS512(b *testing.B) {
	key := csjwt.WithPassword([]byte(`csjwt.SigningMethodHS512!`))
	hmacFast, err := csjwt.NewSigningMethodHS512Fast(key)
	if err != nil {
		b.Fatal(err)
	}
	benchmarkParseFromRequest(
		b,
		hmacFast, // csjwt.SigningMethodHS512,
		key,
		func(t *csjwt.Token) (csjwt.Key, error) {
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
		b.Fatalf("%+v", err)
	}

	r, _ := http.NewRequest("GET", "/", nil)
	r.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tokenString))

	veri := csjwt.NewVerification(sm)

	mc := &ShoppingCartClaim{Standard: new(jwtclaim.Standard)}

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		rToken := csjwt.NewToken(mc)
		err := veri.ParseFromRequest(rToken, keyFunc, r)
		if err != nil {
			b.Fatalf("%+v", err)
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
