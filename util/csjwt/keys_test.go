package csjwt_test

import (
	"github.com/corestoreio/csfw/util/csjwt"
	"testing"
)

func TestKeyParsing(t *testing.T) {

	badKey := []byte("All your base are belong to key")

	key := csjwt.WithRSAPrivateKeyFromFile("test/sample_key")
	// Test parsePrivateKey
	if key.Error != nil {
		t.Errorf("Failed to parse valid private key: %v", key.Error)
	}

	key = csjwt.WithRSAPrivateKeyFromFile("test/sample_key.pub")
	if key.Error == nil {
		t.Errorf("Parsed public key as valid private key: %v", key)
	}

	key = csjwt.WithRSAPrivateKeyFromPEM(badKey)
	if key.Error == nil {
		t.Errorf("Parsed invalid key as valid private key: %v", key)
	}

	// Test parsePublicKey
	key = csjwt.WithRSAPublicKeyFromFile("test/sample_key.pub")
	if key.Error != nil {
		t.Errorf("Failed to parse valid public key: %v", key.Error)
	}

	key = csjwt.WithRSAPublicKeyFromFile("test/sample_key")
	if key.Error == nil {
		t.Errorf("Parsed private key as valid public key: %v", key)
	}

	key = csjwt.WithRSAPublicKeyFromPEM(badKey)
	if key.Error == nil {
		t.Errorf("Parsed invalid key as valid private key: %v", key)
	}
}
