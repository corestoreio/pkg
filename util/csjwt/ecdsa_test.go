package csjwt_test

import (
	"io/ioutil"
	"testing"

	"bytes"

	"github.com/corestoreio/csfw/util/csjwt"
)

var ecdsaTestData = []struct {
	name        string
	keys        map[string]string
	tokenString []byte
	alg         string
	claims      map[string]interface{}
	valid       bool
}{
	{
		"Basic ES256",
		map[string]string{"private": "test/ec256-private.pem", "public": "test/ec256-public.pem"},
		[]byte("eyJ0eXAiOiJKV1QiLCJhbGciOiJFUzI1NiJ9.eyJmb28iOiJiYXIifQ.feG39E-bn8HXAKhzDZq7yEAPWYDhZlwTn3sePJnU9VrGMmwdXAIEyoOnrjreYlVM_Z4N13eK9-TmMTWyfKJtHQ"),
		"ES256",
		map[string]interface{}{"foo": "bar"},
		true,
	},
	{
		"Basic ES384",
		map[string]string{"private": "test/ec384-private.pem", "public": "test/ec384-public.pem"},
		[]byte("eyJ0eXAiOiJKV1QiLCJhbGciOiJFUzM4NCJ9.eyJmb28iOiJiYXIifQ.ngAfKMbJUh0WWubSIYe5GMsA-aHNKwFbJk_wq3lq23aPp8H2anb1rRILIzVR0gUf4a8WzDtrzmiikuPWyCS6CN4-PwdgTk-5nehC7JXqlaBZU05p3toM3nWCwm_LXcld"),
		"ES384",
		map[string]interface{}{"foo": "bar"},
		true,
	},
	{
		"Basic ES512",
		map[string]string{"private": "test/ec512-private.pem", "public": "test/ec512-public.pem"},
		[]byte("eyJ0eXAiOiJKV1QiLCJhbGciOiJFUzUxMiJ9.eyJmb28iOiJiYXIifQ.AAU0TvGQOcdg2OvrwY73NHKgfk26UDekh9Prz-L_iWuTBIBqOFCWwwLsRiHB1JOddfKAls5do1W0jR_F30JpVd-6AJeTjGKA4C1A1H6gIKwRY0o_tFDIydZCl_lMBMeG5VNFAjO86-WCSKwc3hqaGkq1MugPRq_qrF9AVbuEB4JPLyL5"),
		"ES512",
		map[string]interface{}{"foo": "bar"},
		true,
	},
	{
		"basic ES256 invalid: foo => bar",
		map[string]string{"private": "test/ec256-private.pem", "public": "test/ec256-public.pem"},
		[]byte("eyJhbGciOiJFUzI1NiIsInR5cCI6IkpXVCJ9.eyJmb28iOiJiYXIifQ.MEQCIHoSJnmGlPaVQDqacx_2XlXEhhqtWceVopjomc2PJLtdAiAUTeGPoNYxZw0z8mgOnnIcjoxRuNDVZvybRZF3wR1l8W"),
		"ES256",
		map[string]interface{}{"foo": "bar"},
		false,
	},
}

func TestECDSAVerify(t *testing.T) {
	for _, data := range ecdsaTestData {

		key, err := ioutil.ReadFile(data.keys["public"])
		if err != nil {
			t.Fatal(err)
		}

		signing, signature, err := csjwt.SplitForVerify(data.tokenString)
		if err != nil {
			t.Fatal(err, "\n", string(data.tokenString))
		}

		method := csjwt.GetSigningMethod(data.alg)
		err = method.Verify(signing, signature, csjwt.WithECPublicKeyFromPEM(key))
		if data.valid && err != nil {
			t.Errorf("[%v] Error while verifying key: %v", data.name, err)
		}
		if !data.valid && err == nil {
			t.Errorf("[%v] Invalid key passed validation", data.name)
		}
	}
}

func TestECDSASign(t *testing.T) {
	for _, data := range ecdsaTestData {

		key, err := ioutil.ReadFile(data.keys["private"])
		if err != nil {
			t.Fatal(err)
		}

		if data.valid {

			signing, signature, err := csjwt.SplitForVerify(data.tokenString)
			if err != nil {
				t.Fatal(err, "\n", string(data.tokenString))
			}

			method := csjwt.GetSigningMethod(data.alg)
			sig, err := method.Sign(signing, csjwt.WithECPrivateKeyFromPEM(key))
			if err != nil {
				t.Errorf("[%v] Error signing token: %v", data.name, err)
			}
			if bytes.Equal(sig, signature) {
				t.Errorf("[%v] Identical signatures\nbefore:\n%v\nafter:\n%v", data.name, string(signature), string(sig))
			}
		}
	}
}

func BenchmarkES256Signing(b *testing.B) {
	key := csjwt.WithECPrivateKeyFromFile("test/ec256-private.pem")
	benchmarkSigning(b, csjwt.SigningMethodES256, key)
}

func BenchmarkES256Verify(b *testing.B) {
	signing, signature, err := csjwt.SplitForVerify(ecdsaTestData[0].tokenString)
	if err != nil {
		b.Fatal(err)
	}
	key := csjwt.WithECPublicKeyFromFile("test/ec256-public.pem")
	benchmarkMethodVerify(b, csjwt.SigningMethodES256, signing, signature, key)
}

func BenchmarkES384Signing(b *testing.B) {
	key := csjwt.WithECPrivateKeyFromFile("test/ec384-private.pem")
	benchmarkSigning(b, csjwt.SigningMethodES384, key)
}

func BenchmarkES384Verify(b *testing.B) {
	signing, signature, err := csjwt.SplitForVerify(ecdsaTestData[1].tokenString)
	if err != nil {
		b.Fatal(err)
	}
	key := csjwt.WithECPublicKeyFromFile("test/ec384-public.pem")
	benchmarkMethodVerify(b, csjwt.SigningMethodES384, signing, signature, key)
}

func BenchmarkES512Signing(b *testing.B) {
	key := csjwt.WithECPrivateKeyFromFile("test/ec512-private.pem")
	benchmarkSigning(b, csjwt.SigningMethodES512, key)
}

func BenchmarkES512Verify(b *testing.B) {
	signing, signature, err := csjwt.SplitForVerify(ecdsaTestData[2].tokenString)
	if err != nil {
		b.Fatal(err)
	}
	key := csjwt.WithECPublicKeyFromFile("test/ec512-public.pem")
	benchmarkMethodVerify(b, csjwt.SigningMethodES512, signing, signature, key)
}
