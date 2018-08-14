// Copyright 2015-present, Cyrill @ Schumacher.fm and the CoreStore contributors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package modification

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
	"os"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/pkg/config"
)

// AESGCMOptions sets the Key and Nonce from this struct or from environment
// variables.
//easyjson:json
type AESGCMOptions struct {
	// The key argument should be the AES key, either 16, 24, or 32 bytes to
	// select AES-128, AES-192, or AES-256. If empty, a random key will be
	// generated and discarded once the service gets shut down.
	Key string
	// KeyEnvironmentVariableName defines the name of the environment variable
	// which contains the key used for encryption and decryption.
	KeyEnvironmentVariableName   string
	Nonce                        []byte // max 12 bytes
	NonceEnvironmentVariableName string
}

const nonceLength = 12

type aesGCM struct {
	o         AESGCMOptions
	aead      cipher.AEAD
	eventType uint8
}

// NewAESGCM creates a new observer which can encrypt or decrypt a value with
// the AES-GCM mode. Only two events are supported: config.EventOnBeforeSet for
// encryption and config.EventOnAfterGet for decryption.
func NewAESGCM(eventType uint8, eo *AESGCMOptions) (config.Observer, error) {

	if eventType != config.EventOnBeforeSet && eventType != config.EventOnAfterGet {
		return nil, errors.NotValid.Newf("[config/modification] Event type can only be: EventOnBeforeSet (encryption) or EventOnAfterGet (decryption)")
	}

	key := []byte(eo.Key)
	if k, ok := os.LookupEnv(eo.KeyEnvironmentVariableName); ok && eo.KeyEnvironmentVariableName != "" && k != "" {
		key = []byte(k)
	}

	if len(key) == 0 {
		key = make([]byte, 32) // AES-256
		if _, err := io.ReadFull(rand.Reader, key); err != nil {
			return nil, errors.ReadFailed.New(err, "[config/modification] ReadFull failed")
		}
		eo.Key = string(key)
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, errors.NotValid.New(err, "[config/modification] The encryption key has a wrong format.")
	}

	// Never use more than 2^32 = [12]byte random nonces with a given key
	// because of the risk of a repeat.
	nonce := eo.Nonce
	if k, ok := os.LookupEnv(eo.NonceEnvironmentVariableName); ok && eo.NonceEnvironmentVariableName != "" && k != "" {
		nonce = []byte(k)
	}

	if len(eo.Nonce) != nonceLength {
		nonce = make([]byte, nonceLength)
		if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
			return nil, errors.ReadFailed.New(err, "[config/modification] ReadFull failed")
		}
		eo.Nonce = append(eo.Nonce[:0], nonce...)
	}

	aead, err := cipher.NewGCM(block)
	if err != nil {
		return nil, errors.Fatal.New(err, "[config/modification] cipher GCM failed")
	}

	enc := &aesGCM{
		o:         *eo,
		aead:      aead,
		eventType: eventType,
	}

	return enc, nil
}

func (v *aesGCM) Observe(p config.Path, rawData []byte, found bool) ([]byte, error) {

	switch v.eventType {
	case config.EventOnBeforeSet:
		return v.aead.Seal(nil, v.o.Nonce, rawData, nil), nil
	case config.EventOnAfterGet:
		if !found {
			return nil, nil
		}
		plaintext, err := v.aead.Open(nil, v.o.Nonce, rawData, nil)
		if err != nil {
			return nil, errors.Wrapf(err, "[config/modification] For Path %q", p.String())
		}
		return plaintext, nil
	}
	return nil, errors.Fatal.Newf("[config/modification] A programmer made an error")
}
