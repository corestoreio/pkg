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

func ExampleNew(mySigningKey []byte) ([]byte, error) {
	// Create the token
	token := csjwt.New(csjwt.SigningMethodHS256)
	// Set some claims
	token.Claims["foo"] = "bar"
	token.Claims["exp"] = time.Now().Add(time.Hour * 72).Unix()
	// Sign and get the complete encoded token as a string
	return token.SignedString(csjwt.WithPassword(mySigningKey))
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
