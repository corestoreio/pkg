// Copyright 2015-2016, Cyrill @ Schumacher.fm and the CoreStore contributors
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

package config

import (
	"net/url"

	"github.com/corestoreio/csfw/util/errors"
)

// PathCSBaseURL main CoreStore base URL, used if no configuration on a store level can be found.
const (
	PathCSBaseURL = "web/corestore/base_url"
	CSBaseURL     = "http://localhost:9500/"
)

// URL* defines the types of available URLs.
const (
	URLTypeAbsent URLType = iota
	// URLTypeWeb defines the URL type to generate the main base URL.
	URLTypeWeb
	// URLTypeStatic defines the URL to the static assets like CSS, JS or theme images
	URLTypeStatic

	// UrlTypeLink hmmm
	// UrlTypeLink

	// URLTypeMedia defines the URL type for generating URLs to product photos
	URLTypeMedia
	maxURLTypes
)

// URLType defines the type of the URL. Used in constant declaration.
// @see https://github.com/magento/magento2/blob/0.74.0-beta7/lib/internal/Magento/Framework/UrlInterface.php#L13
type URLType uint8

// URLCache represents a type for embedding into other structs to act as an
// internal cache for parsed URLs.
type URLCache struct {
	urls []*url.URL
}

// NewURLCache creates a new cache
func NewURLCache() *URLCache {
	return &URLCache{
		urls: make([]*url.URL, maxURLTypes, maxURLTypes),
	}
}

// Get returns a parsed URL by its Type
func (uc *URLCache) Get(t URLType) *url.URL {
	if t < maxURLTypes {
		return uc.urls[t]
	}
	return nil
}

// Set parses a rawURL and adds it to the cache by its Type. Multiple calls
// with the same type will overwrite existing values.
// Error behaviour: NotFound, Empty and NotValid.
func (uc *URLCache) Set(t URLType, rawURL string) (*url.URL, error) {
	if t >= maxURLTypes {
		return nil, errors.NewNotFoundf("[config] Unknown Index %d", t)
	}
	if rawURL == "" {
		return nil, errors.NewEmptyf("[config] rawURL")
	}
	u, err := url.Parse(rawURL)
	if err != nil {
		return nil, errors.NewNotValid(err, "[config] url.Parse")
	}
	uc.urls[t] = u
	return u, nil
}

// Clear clears the internal cache bucket. Returns nil on success.
func (uc *URLCache) Clear() error {
	uc.urls = make([]*url.URL, maxURLTypes, maxURLTypes)
	return nil
}
