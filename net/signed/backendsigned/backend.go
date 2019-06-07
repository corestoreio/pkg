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

package backendsigned

import (
	"github.com/corestoreio/pkg/config/cfgmodel"
	"github.com/corestoreio/pkg/config/element"
	"github.com/corestoreio/pkg/net/signed"
)

// Configuration just exported for the sake of documentation. See fields for more
// information. Please call the New() function for creating a new Backend
// object. Only the New() function will set the paths to the fields.
type Configuration struct {
	*signed.OptionFactories

	// Disabled set to true to disable the middle ware.
	//
	// Path: net/signed/disabled
	Disabled cfgmodel.Bool

	// InTrailer enables to write the signature resp. hash into the HTTP trailer.
	// Note not all clients can read a trailer.
	//
	// Path: net/signed/in_trailer
	InTrailer cfgmodel.Bool

	// AllowedMethods specifies all HTTP methods which will be accepted by the
	// WithRequestSignatureValidation middleware.
	//
	// Path: net/signed/allowed_methods
	AllowedMethods cfgmodel.StringCSV

	// Key defines the symmetric key/password for the hashing algorithms. You
	// must set a type to satisfy the cfgmodel.Encryptor interface or this
	// package panics.
	//
	// Path: net/signed/key
	Key cfgmodel.Obscure

	// Algorithm defines currently supported cryptographical hashing algorithms
	//
	// Path: net/signed/algorithm
	Algorithm cfgmodel.Str

	// HTTPHeaderType sets the type of the HTTP header to either Content-HMAC or
	// Content-Signature.
	//
	// Path: net/signed/http_header_type
	HTTPHeaderType cfgmodel.Str

	// KeyID name or ID of the key which is used in the HMAC algorithm. Only usable when
	// HTTPHeaderType has been set to "signature"
	//
	// Path: net/signed/key_id
	KeyID cfgmodel.Str
}

// New initializes the backend configuration models containing the cfgpath.Route
// variable to the appropriate entries in the storage. The argument Sections
// and opts will be applied to all models.
func New(cfgStruct element.Sections, opts ...cfgmodel.Option) *Configuration {
	be := &Configuration{
		OptionFactories: signed.NewOptionFactories(),
	}

	opts = append(opts, cfgmodel.WithFieldFromSectionSlice(cfgStruct))

	be.Disabled = cfgmodel.NewBool(`net/signed/disabled`, opts...)
	be.InTrailer = cfgmodel.NewBool(`net/signed/in_trailer`, opts...)
	be.AllowedMethods = cfgmodel.NewStringCSV(`net/signed/allowed_methods`, append(opts, cfgmodel.WithSourceByString(
		"POST", "POST",
		"PUT", "PUT",
		"PATCH", "PATCH",
		"DELETE", "DELETE",
	))...)
	be.Key = cfgmodel.NewObscure(`net/signed/key`, opts...)
	be.Algorithm = cfgmodel.NewStr(`net/signed/algorithm`, append(opts, cfgmodel.WithSourceByString(
		"sha256", "SHA 256",
		"sha512", "SHA 512",
		"blake2", "Blake2b",
	))...)
	be.HTTPHeaderType = cfgmodel.NewStr(`net/signed/http_header_type`, append(opts, cfgmodel.WithSourceByString(
		"transparent", "Transparent",
		"hmac", "Content-HMAC",
		"signature", "Content-Signature",
	))...)
	be.KeyID = cfgmodel.NewStr(`net/signed/key_id`, opts...)

	return be
}
