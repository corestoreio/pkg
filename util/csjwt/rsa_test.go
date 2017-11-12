package csjwt_test

import (
	"bytes"
	"testing"

	"github.com/corestoreio/pkg/util/csjwt"
)

func init() {
	_ = csjwt.NewSigningMethodRS256()
	_ = csjwt.NewSigningMethodRS384()
	_ = csjwt.NewSigningMethodRS512()
}

var rsaTestData = []struct {
	name        string
	tokenString []byte
	method      csjwt.Signer
	claims      map[string]interface{}
	valid       bool
}{
	{
		"Basic RS256",
		[]byte("eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiJ9.eyJmb28iOiJiYXIifQ.FhkiHkoESI_cG3NPigFrxEk9Z60_oXrOT2vGm9Pn6RDgYNovYORQmmA0zs1AoAOf09ly2Nx2YAg6ABqAYga1AcMFkJljwxTT5fYphTuqpWdy4BELeSYJx5Ty2gmr8e7RonuUztrdD5WfPqLKMm1Ozp_T6zALpRmwTIW0QPnaBXaQD90FplAg46Iy1UlDKr-Eupy0i5SLch5Q-p2ZpaL_5fnTIUDlxC3pWhJTyx_71qDI-mAA_5lE_VdroOeflG56sSmDxopPEG3bFlSu1eowyBfxtu0_CuVd-M42RU75Zc4Gsj6uV77MBtbMrf4_7M_NUTSgoIF3fRqxrj0NzihIBg"),
		csjwt.NewSigningMethodRS256(),
		map[string]interface{}{"foo": "bar"},
		true,
	},
	{
		"Basic RS384",
		[]byte("eyJhbGciOiJSUzM4NCIsInR5cCI6IkpXVCJ9.eyJmb28iOiJiYXIifQ.W-jEzRfBigtCWsinvVVuldiuilzVdU5ty0MvpLaSaqK9PlAWWlDQ1VIQ_qSKzwL5IXaZkvZFJXT3yL3n7OUVu7zCNJzdwznbC8Z-b0z2lYvcklJYi2VOFRcGbJtXUqgjk2oGsiqUMUMOLP70TTefkpsgqDxbRh9CDUfpOJgW-dU7cmgaoswe3wjUAUi6B6G2YEaiuXC0XScQYSYVKIzgKXJV8Zw-7AN_DBUI4GkTpsvQ9fVVjZM9csQiEXhYekyrKu1nu_POpQonGd8yqkIyXPECNmmqH5jH4sFiF67XhD7_JpkvLziBpI-uh86evBUadmHhb9Otqw3uV3NTaXLzJw"),
		csjwt.NewSigningMethodRS384(),
		map[string]interface{}{"foo": "bar"},
		true,
	},
	{
		"Basic RS512",
		[]byte("eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJmb28iOiJiYXIifQ.zBlLlmRrUxx4SJPUbV37Q1joRcI9EW13grnKduK3wtYKmDXbgDpF1cZ6B-2Jsm5RB8REmMiLpGms-EjXhgnyh2TSHE-9W2gA_jvshegLWtwRVDX40ODSkTb7OVuaWgiy9y7llvcknFBTIg-FnVPVpXMmeV_pvwQyhaz1SSwSPrDyxEmksz1hq7YONXhXPpGaNbMMeDTNP_1oj8DZaqTIL9TwV8_1wb2Odt_Fy58Ke2RVFijsOLdnyEAjt2n9Mxihu9i3PhNBkkxa2GbnXBfq3kzvZ_xxGGopLdHhJjcGWXO-NiwI9_tiu14NRv4L2xC0ItD9Yz68v2ZIZEp_DuzwRQ"),
		csjwt.NewSigningMethodRS512(),
		map[string]interface{}{"foo": "bar"},
		true,
	},
	{
		"basic invalid: foo => bar",
		[]byte("eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiJ9.eyJmb28iOiJiYXIifQ.EhkiHkoESI_cG3NPigFrxEk9Z60_oXrOT2vGm9Pn6RDgYNovYORQmmA0zs1AoAOf09ly2Nx2YAg6ABqAYga1AcMFkJljwxTT5fYphTuqpWdy4BELeSYJx5Ty2gmr8e7RonuUztrdD5WfPqLKMm1Ozp_T6zALpRmwTIW0QPnaBXaQD90FplAg46Iy1UlDKr-Eupy0i5SLch5Q-p2ZpaL_5fnTIUDlxC3pWhJTyx_71qDI-mAA_5lE_VdroOeflG56sSmDxopPEG3bFlSu1eowyBfxtu0_CuVd-M42RU75Zc4Gsj6uV77MBtbMrf4_7M_NUTSgoIF3fRqxrj0NzihIBg"),
		csjwt.NewSigningMethodRS256(),
		map[string]interface{}{"foo": "bar"},
		false,
	},
}

func TestRSAVerify(t *testing.T) {

	key := csjwt.WithRSAPublicKeyFromFile("test/sample_key.pub")
	for _, data := range rsaTestData {
		signing, signature, err := csjwt.SplitForVerify(data.tokenString)
		if err != nil {
			t.Fatal(err, "\n", string(data.tokenString))
		}

		err = data.method.Verify(signing, signature, key)
		if data.valid && err != nil {
			t.Errorf("[%v] Error while verifying key: %v", data.name, err)
		}
		if !data.valid && err == nil {
			t.Errorf("[%v] Invalid key passed validation", data.name)
		}
	}
}

func TestRSASign(t *testing.T) {

	key := csjwt.WithRSAPrivateKeyFromFile("test/sample_key")

	for _, data := range rsaTestData {
		if data.valid {
			signing, signature, err := csjwt.SplitForVerify(data.tokenString)
			if err != nil {
				t.Fatal(err, "\n", string(data.tokenString))
			}

			sig, err := data.method.Sign(signing, key)
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

	key := csjwt.WithRSAPublicKeyFromFile("test/sample_key.pub")

	testData := rsaTestData[0]

	signing, signature, err := csjwt.SplitForVerify(testData.tokenString)
	if err != nil {
		t.Fatal(err, "\n", string(testData.tokenString))
	}

	sm256 := csjwt.NewSigningMethodRS256()
	err = sm256.Verify(signing, signature, key)
	if err != nil {
		t.Errorf("[%v] Error while verifying key: %v", testData.name, err)
	}
}

func TestRSAWithPreParsedPrivateKey(t *testing.T) {

	key := csjwt.WithRSAPrivateKeyFromFile("test/sample_key")

	testData := rsaTestData[0]

	signing, signature, err := csjwt.SplitForVerify(testData.tokenString)
	if err != nil {
		t.Fatal(err, "\n", string(testData.tokenString))
	}

	sm256 := csjwt.NewSigningMethodRS256()
	sig, err := sm256.Sign(signing, key)
	if err != nil {
		t.Errorf("[%v] Error signing token: %v", testData.name, err)
	}
	if !bytes.Equal(sig, signature) {
		t.Errorf("[%v] Incorrect signature.\nwas:\n%v\nexpecting:\n%v", testData.name, string(sig), string(signature))
	}
}

func BenchmarkRS256Signing(b *testing.B) {
	key := csjwt.WithRSAPrivateKeyFromFile("test/sample_key")
	benchmarkSigning(b, csjwt.NewSigningMethodRS256(), key)
}

func BenchmarkRS384Signing(b *testing.B) {
	key := csjwt.WithRSAPrivateKeyFromFile("test/sample_key")
	benchmarkSigning(b, csjwt.NewSigningMethodRS384(), key)
}

func BenchmarkRS512Signing(b *testing.B) {
	key := csjwt.WithRSAPrivateKeyFromFile("test/sample_key")
	benchmarkSigning(b, csjwt.NewSigningMethodRS512(), key)
}

func BenchmarkRS256Verify(b *testing.B) {
	signing, signature, err := csjwt.SplitForVerify(rsaTestData[0].tokenString)
	if err != nil {
		b.Fatal(err)
	}
	key := csjwt.WithRSAPublicKeyFromFile("test/sample_key.pub")
	benchmarkMethodVerify(b, csjwt.NewSigningMethodRS256(), signing, signature, key)
}

func BenchmarkRS384Verify(b *testing.B) {
	signing, signature, err := csjwt.SplitForVerify(rsaTestData[1].tokenString)
	if err != nil {
		b.Fatal(err)
	}
	key := csjwt.WithRSAPublicKeyFromFile("test/sample_key.pub")
	benchmarkMethodVerify(b, csjwt.NewSigningMethodRS384(), signing, signature, key)
}

func BenchmarkRS512Verify(b *testing.B) {
	signing, signature, err := csjwt.SplitForVerify(rsaTestData[2].tokenString)
	if err != nil {
		b.Fatal(err)
	}
	key := csjwt.WithRSAPublicKeyFromFile("test/sample_key.pub")
	benchmarkMethodVerify(b, csjwt.NewSigningMethodRS512(), signing, signature, key)
}
