// A useful example app.  You can use this to debug your tokens on the command line.
// This is also a great place to look at how you might use this library.
//
// Example usage:
// The following will create and sign a token, then verify it and output the original claims.
//     echo {\"foo\":\"bar\"} | bin/jwt -key test/sample_key -alg RS256 -sign - | bin/jwt -key test/sample_key.pub -verify -
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	"github.com/corestoreio/csfw/util/csjwt"
	"github.com/corestoreio/csfw/util/csjwt/jwtclaim"
)

var (
	// Options
	flagAlg     = flag.String("alg", "", "signing algorithm identifier")
	flagKey     = flag.String("key", "", "path to key file or '-' to read from stdin")
	flagCompact = flag.Bool("compact", false, "output compact JSON")
	flagDebug   = flag.Bool("debug", false, "print out all kinds of debug data")

	// Modes - exactly one of these is required
	flagSign   = flag.String("sign", "", "path to claims object to sign or '-' to read from stdin")
	flagVerify = flag.String("verify", "", "path to JWT token to verify or '-' to read from stdin")
)

func main() {
	// Usage message if you ask for -help or if you mess up inputs.
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  One of the following flags is required: sign, verify\n")
		flag.PrintDefaults()
	}

	// Parse command line options
	flag.Parse()

	// Do the thing.  If something goes wrong, print error to stderr
	// and exit with a non-zero status code
	if err := start(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

// Figure out which thing to do and then do that
func start() error {
	if *flagSign != "" {
		return signToken()
	} else if *flagVerify != "" {
		return verifyToken()
	} else {
		flag.Usage()
		return fmt.Errorf("None of the required flags are present.  What do you want me to do?")
	}
}

// Helper func:  Read input from specified file or stdin
func loadData(p string) ([]byte, error) {
	if p == "" {
		return nil, fmt.Errorf("No path specified")
	}

	var rdr io.Reader
	if p == "-" {
		rdr = os.Stdin
	} else {
		if f, err := os.Open(p); err == nil {
			rdr = f
			defer f.Close()
		} else {
			return nil, err
		}
	}
	return ioutil.ReadAll(rdr)
}

// Print a json object in accordance with the prophecy (or the command line options)
func printJSON(j interface{}) error {
	var out []byte
	var err error

	if *flagCompact == false {
		out, err = json.MarshalIndent(j, "", "    ")
	} else {
		out, err = json.Marshal(j)
	}

	if err == nil {
		fmt.Println(string(out))
	}

	return err
}

// Verify a token and output the claims.  This is a great example
// of how to verify and view a token.
func verifyToken() error {
	// get the token
	tokData, err := loadData(*flagVerify)
	if err != nil {
		return fmt.Errorf("Couldn't read token: %v", err)
	}

	// trim possible whitespace from token
	tokData = regexp.MustCompile(`\s*$`).ReplaceAll(tokData, []byte{})
	if *flagDebug {
		fmt.Fprintf(os.Stderr, "Token len: %v bytes\n", len(tokData))
	}

	v := csjwt.NewVerification(csjwt.NewSigningMethodES256(), csjwt.NewSigningMethodHS256())

	mc := &jwtclaim.Map{}
	// Parse the token.  Load the key from command line option
	token, err := v.ParseWithClaim(tokData, func(t csjwt.Token) (csjwt.Key, error) {
		data, err := loadData(*flagKey)
		if err != nil {
			return csjwt.Key{}, err
		}
		if isEs() {
			return csjwt.WithECPublicKeyFromPEM(data), nil
		}
		return csjwt.WithPassword(data), nil
	}, mc)

	// Print some debug data
	if *flagDebug && token.Raw != nil {
		fmt.Fprintf(os.Stderr, "Header:\n%v\n", token.Header)
		fmt.Fprintf(os.Stderr, "Claims:\n%v\n", token.Claims)
	}

	// Print an error if we can't parse for some reason
	if err != nil {
		return fmt.Errorf("Couldn't parse token: %v", err)
	}

	// Is token invalid?
	if !token.Valid {
		return fmt.Errorf("Token is invalid")
	}

	// Print the token details
	if err := printJSON(token.Claims); err != nil {
		return fmt.Errorf("Failed to output claims: %v", err)
	}

	return nil
}

// Create, sign, and output a token.  This is a great, simple example of
// how to use this library to create and sign a token.
func signToken() error {
	// get the token data from command line arguments
	tokData, err := loadData(*flagSign)
	if err != nil {
		return fmt.Errorf("Couldn't read token: %v", err)
	} else if *flagDebug {
		fmt.Fprintf(os.Stderr, "Token: %v bytes", len(tokData))
	}

	// parse the JSON of the claims
	var claims jwtclaim.Map
	if err := json.Unmarshal(tokData, &claims); err != nil {
		return fmt.Errorf("Couldn't parse claims JSON: %v", err)
	}

	// get the key

	rawKey, err := loadData(*flagKey)
	if err != nil {
		return fmt.Errorf("Couldn't read key: %v", err)
	}

	// get the signing alg
	sMethod, err := csjwt.NewSigningMethodByAlg(*flagAlg)
	if sMethod == nil || err != nil {
		return fmt.Errorf("Couldn't find signing method: %q - Error: %s", *flagAlg, err)
	}

	// create a new token
	token := csjwt.NewToken(claims)

	var key csjwt.Key
	if isEs() {
		key = csjwt.WithECPrivateKeyFromPEM(rawKey)
	} else {
		key = csjwt.WithPassword(rawKey)
	}

	if out, err := token.SignedString(sMethod, key); err == nil {
		fmt.Println(out)
	} else {
		return fmt.Errorf("Error signing token: %v", err)
	}

	return nil
}

func isEs() bool {
	return strings.HasPrefix(*flagAlg, "ES")
}
