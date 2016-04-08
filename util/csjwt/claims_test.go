package csjwt_test

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"errors"
	"math"

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

func TestClaimsValid(t *testing.T) {
	tests := []struct {
		sc        csjwt.Claimer
		wantValid error
	}{
		{&csjwt.StandardClaims{}, csjwt.ErrValidationClaimsInvalid},
		{&csjwt.StandardClaims{ExpiresAt: time.Now().Add(time.Second).Unix()}, nil},
		{&csjwt.StandardClaims{ExpiresAt: time.Now().Add(-time.Second).Unix()}, csjwt.ErrValidationExpired},
		{&csjwt.StandardClaims{IssuedAt: time.Now().Add(-time.Second).Unix()}, nil},
		{&csjwt.StandardClaims{IssuedAt: time.Now().Add(time.Second * 5).Unix()}, csjwt.ErrValidationUsedBeforeIssued},
		{&csjwt.StandardClaims{NotBefore: time.Now().Add(-time.Second).Unix()}, nil},
		{&csjwt.StandardClaims{NotBefore: time.Now().Add(time.Second * 5).Unix()}, csjwt.ErrValidationNotValidYet},
		{
			&csjwt.StandardClaims{
				ExpiresAt: time.Now().Add(-time.Second).Unix(),
				IssuedAt:  time.Now().Add(time.Second * 5).Unix(),
				NotBefore: time.Now().Add(time.Second * 5).Unix(),
			},
			fmt.Errorf("%s\n%s\n%s", csjwt.ErrValidationExpired, csjwt.ErrValidationUsedBeforeIssued, csjwt.ErrValidationNotValidYet),
		},

		{csjwt.MapClaims{}, csjwt.ErrValidationClaimsInvalid},                                         // 7
		{csjwt.MapClaims{"exp": time.Now().Add(time.Second).Unix()}, nil},                             // 8
		{csjwt.MapClaims{"exp": time.Now().Add(-time.Second * 2).Unix()}, csjwt.ErrValidationExpired}, // 9
		{csjwt.MapClaims{"iat": time.Now().Add(-time.Second).Unix()}, nil},                            // 10
		{csjwt.MapClaims{"iat": time.Now().Add(time.Second * 5).Unix()}, csjwt.ErrValidationUsedBeforeIssued},
		{csjwt.MapClaims{"nbf": time.Now().Add(-time.Second).Unix()}, nil},
		{csjwt.MapClaims{"nbf": time.Now().Add(time.Second * 5).Unix()}, csjwt.ErrValidationNotValidYet},
		{
			csjwt.MapClaims{
				"exp": time.Now().Add(-time.Second).Unix(),
				"iat": time.Now().Add(time.Second * 5).Unix(),
				"nbf": time.Now().Add(time.Second * 5).Unix(),
			},
			fmt.Errorf("%s\n%s\n%s", csjwt.ErrValidationExpired, csjwt.ErrValidationUsedBeforeIssued, csjwt.ErrValidationNotValidYet),
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

func TestClaimsGetSet(t *testing.T) {
	tests := []struct {
		sc         csjwt.Claimer
		key        string
		val        interface{}
		wantSetErr error
		wantGetErr error
	}{
		{&csjwt.StandardClaims{}, csjwt.ClaimAudience, '', errors.New("Cannot convert key \"aud\" with value 63743 to string or int64"), nil},
		{&csjwt.StandardClaims{}, csjwt.ClaimAudience, "Go", nil, nil},
		{&csjwt.StandardClaims{}, csjwt.ClaimExpiresAt, time.Now().Unix(), nil, nil},
		{&csjwt.StandardClaims{}, "Not Supported", time.Now().Unix(), errors.New("Claim \"Not Supported\" not supported. Please see constants Claim*."), errors.New("Claim \"Not Supported\" not supported. Please see constants Claim*.")},

		{csjwt.MapClaims{}, csjwt.ClaimAudience, "Go", nil, nil},
		{csjwt.MapClaims{}, csjwt.ClaimExpiresAt, time.Now().Unix(), nil, nil},
		{csjwt.MapClaims{}, "Not Supported", math.Pi, nil, nil},
	}
	for i, test := range tests {

		haveSetErr := test.sc.Set(test.key, test.val)
		if test.wantSetErr != nil {
			assert.EqualError(t, haveSetErr, test.wantSetErr.Error(), "Index %d", i)
		} else {
			assert.NoError(t, haveSetErr, "Index %d", i)
		}

		haveVal, haveGetErr := test.sc.Get(test.key)
		if test.wantGetErr != nil {
			assert.EqualError(t, haveGetErr, test.wantGetErr.Error(), "Index %d", i)
			continue
		} else {
			assert.NoError(t, haveGetErr, "Index %d", i)
		}

		if test.wantSetErr == nil {
			assert.Exactly(t, test.val, haveVal, "Index %d", i)
		}
	}
}

func TestClaimsExpires(t *testing.T) {
	tm := time.Now()
	tests := []struct {
		sc      csjwt.Claimer
		wantExp time.Duration
	}{
		{&csjwt.StandardClaims{ExpiresAt: tm.Add(time.Second * 2).Unix()}, time.Second * 1},
		{&csjwt.StandardClaims{ExpiresAt: tm.Add(time.Second * 5).Unix()}, time.Second * 4},
		{&csjwt.StandardClaims{ExpiresAt: -123123}, time.Duration(0)},
		{&csjwt.StandardClaims{}, time.Duration(0)},

		{csjwt.MapClaims{"exp": tm.Add(time.Second * 2).Unix()}, time.Second * 1},
		{csjwt.MapClaims{"exp": tm.Add(time.Second * 22).Unix()}, time.Second * 21},
		{csjwt.MapClaims{"exp": -123123}, time.Duration(0)},
		{csjwt.MapClaims{"eXp": 23}, time.Duration(0)},
		{csjwt.MapClaims{"exp": fmt.Sprintf("%d", tm.Unix()+10)}, time.Second * 9},
	}
	for i, test := range tests {
		assert.Exactly(t, int64(test.wantExp.Seconds()), int64(test.sc.Expires().Seconds()), "Index %d", i)
	}
}
