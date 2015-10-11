// Copyright 2015, Cyrill @ Schumacher.fm and the CoreStore contributors
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

package scope

var _ WebsiteIDer = MockID(0)
var _ GroupIDer = MockID(0)
var _ StoreIDer = MockID(0)

// AdminScope is always 0 in each of the three scopes.
var AdminScope MockID

// MockID is convenience helper to satisfy the interface WebsiteIDer, GroupIDer and StoreIDer.
type MockID int64

// WebsiteID is convenience helper to satisfy the interface WebsiteIDer
func (i MockID) WebsiteID() int64 { return int64(i) }

// GroupID is convenience helper to satisfy the interface GroupIDer
func (i MockID) GroupID() int64 { return int64(i) }

// StoreID is convenience helper to satisfy the interface StoreIDer
func (i MockID) StoreID() int64 { return int64(i) }

var _ StoreCoder = MockCode("")
var _ WebsiteCoder = MockCode("")
var _ WebsiteIDer = MockCode("")
var _ GroupIDer = MockCode("")
var _ StoreIDer = MockCode("")

// MockCode is convenience helper to satisfy the interface WebsiteCoder, StoreCoder,
// WebsiteIDer, GroupIDer and StoreIDer. Reason: In package store all functions have
// as argument an *IDer interface but once they detect an *Coder interface, they
// will once the *Coder return value.
type MockCode string

// WebsiteID is convenience helper to satisfy the interface WebsiteIDer. Returns -1.
func (c MockCode) WebsiteID() int64 { return -1 }

// WebsiteCode mock helper to return a website code
func (c MockCode) WebsiteCode() string { return string(c) }

// StoreID is convenience helper to satisfy the interface StoreIDer. Returns -1.
func (c MockCode) StoreID() int64 { return -1 }

// StoreCode mock helper to return a store code
func (c MockCode) StoreCode() string { return string(c) }

// GroupID is convenience helper to satisfy the interface GroupIDer. Returns -1.
func (c MockCode) GroupID() int64 { return -1 }
