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

package store

// ErrStoreChangeNotAllowed if a given store within a website would like to
// switch to another store in a different website.
const errStoreChangeNotAllowed = "[store] Store change not allowed"

const (
	errStoreNotFound         = "[store] Store not found"
	errStoreDefaultNotFound  = "[store] Default Store ID not found"
	errStoreNotActive        = "[store] Store not active"
	errStoreIncorrectGroup   = "[store] Incorrect group"
	errStoreIncorrectWebsite = "[store] Incorrect website"
	errStoreCodeInvalid      = "[store] The store code may contain only letters (a-z), numbers (0-9) or underscore(_). The first character must be a letter. Have: %q"
)

const (
	errGroupDefaultStoreNotFound   = "[store] Group default store %d not found"
	errGroupWebsiteNotFound        = "[store] Group Website not found or nil or ID do not match"
	errGroupWebsiteIntegrityFailed = "[store] Groups WebsiteID does not match the Websites ID"
)

// ErrWebsite* are general errors when handling with the Website type.
// They are self explanatory.
const (
	errWebsiteDefaultGroupNotFound = "[store] Website Default Group not found"
)
