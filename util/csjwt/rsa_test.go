package csjwt_test

import (
	"io/ioutil"
	"testing"

	"bytes"

	"github.com/corestoreio/csfw/util/csjwt"
)

var rsaTestData = []struct {
	name        string
	tokenString []byte
	alg         string
	claims      map[string]interface{}
	valid       bool
}{
	{
		"Basic RS256",
		[]byte("eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiJ9.eyJmb28iOiJiYXIifQ.FhkiHkoESI_cG3NPigFrxEk9Z60_oXrOT2vGm9Pn6RDgYNovYORQmmA0zs1AoAOf09ly2Nx2YAg6ABqAYga1AcMFkJljwxTT5fYphTuqpWdy4BELeSYJx5Ty2gmr8e7RonuUztrdD5WfPqLKMm1Ozp_T6zALpRmwTIW0QPnaBXaQD90FplAg46Iy1UlDKr-Eupy0i5SLch5Q-p2ZpaL_5fnTIUDlxC3pWhJTyx_71qDI-mAA_5lE_VdroOeflG56sSmDxopPEG3bFlSu1eowyBfxtu0_CuVd-M42RU75Zc4Gsj6uV77MBtbMrf4_7M_NUTSgoIF3fRqxrj0NzihIBg"),
		"RS256",
		map[string]interface{}{"foo": "bar"},
		true,
	},
	{
		"Basic RS384",
		[]byte("eyJhbGciOiJSUzM4NCIsInR5cCI6IkpXVCJ9.eyJmb28iOiJiYXIifQ.W-jEzRfBigtCWsinvVVuldiuilzVdU5ty0MvpLaSaqK9PlAWWlDQ1VIQ_qSKzwL5IXaZkvZFJXT3yL3n7OUVu7zCNJzdwznbC8Z-b0z2lYvcklJYi2VOFRcGbJtXUqgjk2oGsiqUMUMOLP70TTefkpsgqDxbRh9CDUfpOJgW-dU7cmgaoswe3wjUAUi6B6G2YEaiuXC0XScQYSYVKIzgKXJV8Zw-7AN_DBUI4GkTpsvQ9fVVjZM9csQiEXhYekyrKu1nu_POpQonGd8yqkIyXPECNmmqH5jH4sFiF67XhD7_JpkvLziBpI-uh86evBUadmHhb9Otqw3uV3NTaXLzJw"),
		"RS384",
		map[string]interface{}{"foo": "bar"},
		true,
	},
	{
		"Basic RS512",
		[]byte("eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJmb28iOiJiYXIifQ.zBlLlmRrUxx4SJPUbV37Q1joRcI9EW13grnKduK3wtYKmDXbgDpF1cZ6B-2Jsm5RB8REmMiLpGms-EjXhgnyh2TSHE-9W2gA_jvshegLWtwRVDX40ODSkTb7OVuaWgiy9y7llvcknFBTIg-FnVPVpXMmeV_pvwQyhaz1SSwSPrDyxEmksz1hq7YONXhXPpGaNbMMeDTNP_1oj8DZaqTIL9TwV8_1wb2Odt_Fy58Ke2RVFijsOLdnyEAjt2n9Mxihu9i3PhNBkkxa2GbnXBfq3kzvZ_xxGGopLdHhJjcGWXO-NiwI9_tiu14NRv4L2xC0ItD9Yz68v2ZIZEp_DuzwRQ"),
		"RS512",
		map[string]interface{}{"foo": "bar"},
		true,
	},
	{
		"basic invalid: foo => bar",
		[]byte("eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiJ9.eyJmb28iOiJiYXIifQ.EhkiHkoESI_cG3NPigFrxEk9Z60_oXrOT2vGm9Pn6RDgYNovYORQmmA0zs1AoAOf09ly2Nx2YAg6ABqAYga1AcMFkJljwxTT5fYphTuqpWdy4BELeSYJx5Ty2gmr8e7RonuUztrdD5WfPqLKMm1Ozp_T6zALpRmwTIW0QPnaBXaQD90FplAg46Iy1UlDKr-Eupy0i5SLch5Q-p2ZpaL_5fnTIUDlxC3pWhJTyx_71qDI-mAA_5lE_VdroOeflG56sSmDxopPEG3bFlSu1eowyBfxtu0_CuVd-M42RU75Zc4Gsj6uV77MBtbMrf4_7M_NUTSgoIF3fRqxrj0NzihIBg"),
		"RS256",
		map[string]interface{}{"foo": "bar"},
		false,
	},
}

func TestRSAVerify(t *testing.T) {
	key, _ := ioutil.ReadFile("test/sample_key.pub")

	for _, data := range rsaTestData {
		signing, signature, err := csjwt.SplitForVerify(data.tokenString)
		if err != nil {
			t.Fatal(err, "\n", string(data.tokenString))
		}

		method := csjwt.GetSigningMethod(data.alg)
		err = method.Verify(signing, signature, key)
		if data.valid && err != nil {
			t.Errorf("[%v] Error while verifying key: %v", data.name, err)
		}
		if !data.valid && err == nil {
			t.Errorf("[%v] Invalid key passed validation", data.name)
		}
	}
}

func TestRSASign(t *testing.T) {
	key, _ := ioutil.ReadFile("test/sample_key")

	for _, data := range rsaTestData {
		if data.valid {
			signing, signature, err := csjwt.SplitForVerify(data.tokenString)
			if err != nil {
				t.Fatal(err, "\n", string(data.tokenString))
			}

			method := csjwt.GetSigningMethod(data.alg)
			sig, err := method.Sign(signing, key)
			if err != nil {
				t.Errorf("[%v] Error signing token: %v", data.name, err)
			}
			if !bytes.Equal(sig, signature) {
				t.Errorf("[%v] Incorrect signature.\nwas:\n%v\nexpecting:\n%v", data.name, string(sig), string(signature))
			}
		}
	}
}

func TestRSAVerifyWithPreParsedPrivateKey(t *testing.T) {
	key, _ := ioutil.ReadFile("test/sample_key.pub")
	parsedKey, err := csjwt.ParseRSAPublicKeyFromPEM(key)
	if err != nil {
		t.Fatal(err)
	}
	testData := rsaTestData[0]

	signing, signature, err := csjwt.SplitForVerify(testData.tokenString)
	if err != nil {
		t.Fatal(err, "\n", string(testData.tokenString))
	}

	err = csjwt.SigningMethodRS256.Verify(signing, signature, parsedKey)
	if err != nil {
		t.Errorf("[%v] Error while verifying key: %v", testData.name, err)
	}
}

func TestRSAWithPreParsedPrivateKey(t *testing.T) {
	key, _ := ioutil.ReadFile("test/sample_key")
	parsedKey, err := csjwt.ParseRSAPrivateKeyFromPEM(key)
	if err != nil {
		t.Fatal(err)
	}
	testData := rsaTestData[0]

	signing, signature, err := csjwt.SplitForVerify(testData.tokenString)
	if err != nil {
		t.Fatal(err, "\n", string(testData.tokenString))
	}

	sig, err := csjwt.SigningMethodRS256.Sign(signing, parsedKey)
	if err != nil {
		t.Errorf("[%v] Error signing token: %v", testData.name, err)
	}
	if !bytes.Equal(sig, signature) {
		t.Errorf("[%v] Incorrect signature.\nwas:\n%v\nexpecting:\n%v", testData.name, string(sig), string(signature))
	}
}

func TestRSAKeyParsing(t *testing.T) {
	key, _ := ioutil.ReadFile("test/sample_key")
	pubKey, _ := ioutil.ReadFile("test/sample_key.pub")
	badKey := []byte("All your base are belong to key")

	// Test parsePrivateKey
	if _, e := csjwt.ParseRSAPrivateKeyFromPEM(key); e != nil {
		t.Errorf("Failed to parse valid private key: %v", e)
	}

	if k, e := csjwt.ParseRSAPrivateKeyFromPEM(pubKey); e == nil {
		t.Errorf("Parsed public key as valid private key: %v", k)
	}

	if k, e := csjwt.ParseRSAPrivateKeyFromPEM(badKey); e == nil {
		t.Errorf("Parsed invalid key as valid private key: %v", k)
	}

	// Test parsePublicKey
	if _, e := csjwt.ParseRSAPublicKeyFromPEM(pubKey); e != nil {
		t.Errorf("Failed to parse valid public key: %v", e)
	}

	if k, e := csjwt.ParseRSAPublicKeyFromPEM(key); e == nil {
		t.Errorf("Parsed private key as valid public key: %v", k)
	}

	if k, e := csjwt.ParseRSAPublicKeyFromPEM(badKey); e == nil {
		t.Errorf("Parsed invalid key as valid private key: %v", k)
	}

}

func BenchmarkRS256Signing(b *testing.B) {
	key, _ := ioutil.ReadFile("test/sample_key")
	parsedKey, err := csjwt.ParseRSAPrivateKeyFromPEM(key)
	if err != nil {
		b.Fatal(err)
	}

	benchmarkSigning(b, csjwt.SigningMethodRS256, parsedKey)
}

func BenchmarkRS384Signing(b *testing.B) {
	key, _ := ioutil.ReadFile("test/sample_key")
	parsedKey, err := csjwt.ParseRSAPrivateKeyFromPEM(key)
	if err != nil {
		b.Fatal(err)
	}

	benchmarkSigning(b, csjwt.SigningMethodRS384, parsedKey)
}

func BenchmarkRS512Signing(b *testing.B) {
	key, _ := ioutil.ReadFile("test/sample_key")
	parsedKey, err := csjwt.ParseRSAPrivateKeyFromPEM(key)
	if err != nil {
		b.Fatal(err)
	}

	benchmarkSigning(b, csjwt.SigningMethodRS512, parsedKey)
}
