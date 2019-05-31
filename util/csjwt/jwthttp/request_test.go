package jwthttp_test

import (
	"fmt"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/pkg/util/assert"
	"github.com/corestoreio/pkg/util/conv"
	"github.com/corestoreio/pkg/util/csjwt"
	"github.com/corestoreio/pkg/util/csjwt/jwtclaim"
	"github.com/corestoreio/pkg/util/csjwt/jwthttp"
)

func makeSample(c jwtclaim.Map) []byte {
	token := csjwt.NewToken(c)
	s, err := token.SignedString(csjwt.NewSigningMethodRS256(), csjwt.WithRSAPrivateKeyFromFile("../test/sample_key"))
	if err != nil {
		panic(err)
	}
	return s
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
	err = jwthttp.NewVerification(sm512).ParseFromRequest(rToken, func(t *csjwt.Token) (csjwt.Key, error) {
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
	vf := jwthttp.NewVerification(rs256)
	vf.CookieName = cookieName

	pubKey := csjwt.WithRSAPublicKeyFromFile("../test/sample_key.pub")

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
		jwthttp.HTTPFormInputName: []string{string(token)},
	}

	rs256 := csjwt.NewSigningMethodRS256()
	vf := jwthttp.NewVerification(rs256)
	vf.CookieName = "unset"
	vf.FormInputName = jwthttp.HTTPFormInputName

	pubKey := csjwt.WithRSAPublicKeyFromFile("../test/sample_key.pub")
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
	vf := jwthttp.NewVerification(rs256)
	vf.CookieName = "unset"
	vf.FormInputName = jwthttp.HTTPFormInputName

	pubKey := csjwt.WithRSAPublicKeyFromFile("../test/sample_key.pub")
	haveToken := csjwt.NewToken(&jwtclaim.Map{})
	haveErr := vf.ParseFromRequest(haveToken, csjwt.NewKeyFunc(rs256, pubKey), r)
	assert.True(t, errors.NotFound.Match(haveErr), "Error: %s", haveErr)
	assert.Empty(t, haveToken.Raw)
	assert.False(t, haveToken.Valid)
}
