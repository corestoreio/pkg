package csjwt_test

//
//import (
//	"fmt"
//	"time"
//
//	"github.com/corestoreio/pkg/util/cserr"
//	"github.com/corestoreio/pkg/util/csjwt"
//	"github.com/corestoreio/pkg/util/csjwt/jwtclaim"
//)
//
//func ExampleVerification(myToken []byte, myLookupKey func(interface{}) (csjwt.Key, error)) {
//	base := csjwt.NewToken(&jwtclaim.Map{})
//	token, err := csjwt.NewVerification().Parse(base, myToken, func(token csjwt.Token) (csjwt.Key, error) {
//		kid, err := token.Header.Get("kid")
//		if err != nil {
//			return csjwt.Key{}, err
//		}
//		return myLookupKey(kid)
//	})
//
//	if err == nil && token.Valid {
//		fmt.Println("Your token is valid.  I like your style.")
//	} else {
//		fmt.Println("This token is terrible!  I cannot accept this.")
//	}
//}
//
//func ExampleToken() {
//	// Create the token
//	token := csjwt.NewToken(&jwtclaim.Map{})
//
//	// Set some claims
//	if err := token.Claims.Set("foo", "bar"); err != nil {
//		panic(err) // only use panic while testing
//	}
//	if err := token.Claims.Set(jwtclaim.KeyExpiresAt, time.Now().Add(time.Hour*1).Unix()); err != nil {
//		panic(err) // only use panic while testing
//	}
//
//	foo, err := token.Claims.Get("foo")
//	exp := token.Claims.Expires()
//	fmt.Printf("<%T> foo:%v foo:err:%s exp:%dm\n", token.Claims, foo, err, int(exp.Minutes()))
//	//Output: <*jwtclaim.Map> foo:bar foo:err:%!s(<nil>) exp:59m
//}
//
//func ExampleToken_WithClaims() {
//	// {"alg":"HS256","typ":"JWT"}
//
//	// Create the Claims
//	claims := &jwtclaim.Standard{
//		ExpiresAt: 15000,
//		Issuer:    "test",
//	}
//
//	token := csjwt.NewToken(claims)
//	ss, err := token.SignedString(csjwt.NewSigningMethodHS256(), csjwt.WithPassword([]byte("AllYourBase")))
//	fmt.Printf("%s %v", string(ss), err)
//	//Output: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9Cg.eyJleHAiOjE1MDAwLCJpc3MiOiJ0ZXN0In0K.SY9CGIVRoJ_TBDxtPfv3C8h_qsmOLqd9YrYQcMp-rFM <nil>
//}
//
//func ExampleToken_customType() {
//
//	type MyCustomClaims struct {
//		Foo string `json:"foo"`
//		*jwtclaim.Standard
//	}
//
//	// Create the Claims
//	claims := MyCustomClaims{
//		"bar",
//		&jwtclaim.Standard{
//			ExpiresAt: 15000,
//			Issuer:    "test",
//		},
//	}
//
//	token := csjwt.NewToken(&claims)
//	ss, err := token.SignedString(csjwt.NewSigningMethodHS256(), csjwt.WithPassword([]byte("AllYourBase")))
//	fmt.Printf("%s %v\n", string(ss), err)
//
//	ss, err = token.SignedString(csjwt.NewSigningMethodHS512(), csjwt.WithPassword([]byte("AllYourZombies")))
//	fmt.Printf("%s %v", string(ss), err)
//	//Output:
//	// eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9Cg.eyJmb28iOiJiYXIiLCJleHAiOjE1MDAwLCJpc3MiOiJ0ZXN0In0K.5T_mKrp6g5d6lWwsRNO67lHorKoEXAOVoUmnFOox6GA <nil>
//	// eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9Cg.eyJmb28iOiJiYXIiLCJleHAiOjE1MDAwLCJpc3MiOiJ0ZXN0In0K.zy8pbCUubZKYssrZaHN5riSWvvM7utQO4aeflq9Y5PN_ibL75lEJ3JKOmzgRNGbEUMpGv1MMIFnP9JQ6zUzTMw <nil>
//}
//
//func ExampleVerification_errorChecking(myToken []byte, myLookupKey func(interface{}) (csjwt.Key, error)) {
//	baseToken := csjwt.Token{
//		Header: jwtclaim.NewHeadSegments(),
//		Claims: jwtclaim.Map{},
//	}
//	token, err := csjwt.NewVerification().Parse(baseToken, myToken, func(token csjwt.Token) (csjwt.Key, error) {
//		kid, err := token.Header.Get("kid")
//		if err != nil {
//			return csjwt.Key{}, err
//		}
//		return myLookupKey(kid)
//	})
//
//	if err != nil {
//		switch {
//		case cserr.Contains(err, csjwt.ErrTokenMalformed):
//			fmt.Println("That's not even a token")
//		case cserr.Contains(err, jwtclaim.ErrValidationExpired):
//			fmt.Println("Token Expired")
//		case cserr.Contains(err, jwtclaim.ErrValidationNotValidYet):
//			fmt.Println("Token not yet valid!")
//		default:
//			fmt.Println(err.Error())
//		}
//		return
//	}
//
//	if token.Valid {
//		fmt.Println("You look nice today")
//	} else {
//		fmt.Println("Couldn't handle this token:", err)
//	}
//}
