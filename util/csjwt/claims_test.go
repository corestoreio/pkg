package csjwt_test

import (
	"encoding/json"
	"testing"
	"time"

	"errors"
	"github.com/corestoreio/csfw/util/csjwt"
	"github.com/stretchr/testify/assert"
)

var _ csjwt.Claimer = (*csjwt.StandardClaims)(nil)
var _ csjwt.Claimer = (*csjwt.MapClaims)(nil)

func TestStandardClaimsParseJSON(t *testing.T) {

	sc := csjwt.StandardClaims{
		Audience:  `LOTR`,
		ExpiresAt: time.Now().Add(time.Hour).Unix(),

		IssuedAt:  time.Now().Unix(),
		Issuer:    `Gandalf`,
		NotBefore: time.Now().Unix(),
		Subject:   `Test Subject`,
	}
	rawJSON, err := json.Marshal(sc)
	if err != nil {
		t.Fatal(err)
	}
	assert.Len(t, rawJSON, 102, "JSON: %s", rawJSON)

	scNew := csjwt.StandardClaims{}
	if err := json.Unmarshal(rawJSON, &scNew); err != nil {
		t.Fatal(err)
	}
	assert.Exactly(t, sc, scNew)
	assert.NoError(t, scNew.Valid())
}

func TestStandardClaimsValid(t *testing.T) {
	tests := []struct {
		sc        *csjwt.StandardClaims
		wantValid error
	}{
		{
			&csjwt.StandardClaims{
				ExpiresAt: time.Now().Add(time.Second).Unix(),
			},
			nil,
		},
		{
			&csjwt.StandardClaims{
				ExpiresAt: time.Now().Add(-time.Second).Unix(),
			},
			errors.New("Token is expired"),
		},
	}
	for i, test := range tests {
		if test.wantValid != nil {
			assert.EqualError(t, test.sc.Valid(), test.wantValid.Error(), "Index %d", i)
		} else {
			assert.NoError(t, test.sc.Valid(), "Index %d", i)
		}
	}
}
