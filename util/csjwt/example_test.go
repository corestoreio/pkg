package csjwt_test

import (
	"fmt"
	"time"

	"github.com/corestoreio/csfw/util/csjwt"
)

func ExampleParse(myToken []byte, myLookupKey func(interface{}) (csjwt.Key, error)) {
	token, err := csjwt.Parse(myToken, func(token csjwt.Token) (csjwt.Key, error) {
		return myLookupKey(token.Header["kid"])
	})

	if err == nil && token.Valid {
		fmt.Println("Your token is valid.  I like your style.")
	} else {
		fmt.Println("This token is terrible!  I cannot accept this.")
	}
}

func ExampleNew() {
	// Create the token
	token := csjwt.New(csjwt.SigningMethodRS256)

	// Set some claims
	claims := token.Claims.(csjwt.MapClaims)
	claims["foo"] = "bar"
	claims["exp"] = time.Unix(0, 0).Add(time.Hour * 1).Unix()

	fmt.Printf("<%T> foo:%v exp:%v\n", token.Claims, claims["foo"], claims["exp"])
	//Output: <csjwt.MapClaims> foo:bar exp:3600
}

func ExampleNewWithClaims() {
	// {"alg":"HS256","typ":"JWT"}

	// Create the Claims
	claims := &csjwt.StandardClaims{
		ExpiresAt: 15000,
		Issuer:    "test",
	}

	token := csjwt.NewWithClaims(csjwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(csjwt.WithPassword([]byte("AllYourBase")))
	fmt.Printf("%s %v", string(ss), err)
	//Output: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9Cg.eyJleHAiOjE1MDAwLCJpc3MiOiJ0ZXN0In0K.SY9CGIVRoJ_TBDxtPfv3C8h_qsmOLqd9YrYQcMp-rFM <nil>
}

func ExampleNewWithClaims_customType() {

	type MyCustomClaims struct {
		Foo string `json:"foo"`
		csjwt.StandardClaims
	}

	// Create the Claims
	claims := MyCustomClaims{
		"bar",
		csjwt.StandardClaims{
			ExpiresAt: 15000,
			Issuer:    "test",
		},
	}

	token := csjwt.NewWithClaims(csjwt.SigningMethodHS256, &claims)
	ss, err := token.SignedString(csjwt.WithPassword([]byte("AllYourBase")))
	fmt.Printf("%s %v", string(ss), err)
	//Output: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9Cg.eyJmb28iOiJiYXIiLCJleHAiOjE1MDAwLCJpc3MiOiJ0ZXN0In0K.5T_mKrp6g5d6lWwsRNO67lHorKoEXAOVoUmnFOox6GA <nil>
}

func ExampleParse_errorChecking(myToken []byte, myLookupKey func(interface{}) (csjwt.Key, error)) {
	token, err := csjwt.Parse(myToken, func(token csjwt.Token) (csjwt.Key, error) {
		return myLookupKey(token.Header["kid"])
	})

	if token.Valid {
		fmt.Println("You look nice today")
	} else if ve, ok := err.(*csjwt.ValidationError); ok {
		if ve.Errors&csjwt.ValidationErrorMalformed != 0 {
			fmt.Println("That's not even a token")
		} else if ve.Errors&(csjwt.ValidationErrorExpired|csjwt.ValidationErrorNotValidYet) != 0 {
			// Token is either expired or not active yet
			fmt.Println("Timing is everything")
		} else {
			fmt.Println("Couldn't handle this token:", err)
		}
	} else {
		fmt.Println("Couldn't handle this token:", err)
	}

}
