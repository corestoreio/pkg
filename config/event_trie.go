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

// The MIT License (MIT)
//
// Copyright (c) 2014 Dalton Hubble
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package config

import (
	"strings"

	"github.com/corestoreio/errors"
)

// Event constants defines where and when a specific blocking event gets
// dispatched. The concrete value of a constant may change without notice.
const (
	EventOnBeforeSet uint8 = iota // must start with zero
	EventOnAfterSet
	EventOnBeforeGet
	EventOnAfterGet
	eventMaxCount
)

// EventObserver gets called when an event gets dispatched. Not all events
// support modifying the raw data argument.
// For example the EventOnAfterGet allows to decrypt encrypted data.
type EventObserver interface {
	// Observe must always return the rawData argument or an error.
	// Observer can transform and modify the raw data in any case.
	Observe(p Path, found bool, rawData []byte) (rawData2 []byte, err error)
}

type eventObservers []EventObserver

func (fns eventObservers) dispatch(p *Path, found bool, v []byte) (_ []byte, err error) {
	if len(fns) == 0 {
		return v, nil
	}
	p2 := *p
	for idx, fn := range fns {
		if v, err = fn.Observe(p2, found, v); err != nil {
			return nil, errors.Wrapf(err, "[config] At index %d", idx)
		}
	}
	return v, nil
}

// walkFn defines some action to take on the given key and value during
// a Trie Walk. Returning a non-nil error will terminate the Walk.
type walkFn func(key string, value eventObservers) error

// segmentPath segments string key paths by slash separators. For example,
// "/a/b/c" -> ("/a", 2), ("/b", 4), ("/c", -1) in successive calls. It does
// not allocate any heap memory.
func segmentPath(path string, start int) (segment string, next int) {
	if len(path) == 0 || start < 0 || start > len(path)-1 {
		return "", -1
	}
	end := strings.IndexRune(path[start+1:], '/') // next '/' after 0th rune
	if end == -1 {
		return path[start:], -1
	}
	return path[start : start+end+1], start + end + 1
}

// triePath is a trie of paths with string keys and interface{} values.

// triePath is a trie of string keys and interface{} values. Internal nodes
// have nil values so stored nil values cannot be distinguished and are
// excluded from walks. By default, triePath will segment keys by forward
// slashes with segmentPath (e.g. "/a/b/c" -> "/a", "/b", "/c"). A custom
// segmentStringFn may be used to customize how strings are segmented into
// nodes. A classic trie might segment keys by rune (i.e. unicode points).
type triePath struct {
	value    eventObservers
	children map[string]*triePath
}

// newTriePath allocates and returns a new *triePath.
func newTriePath() *triePath {
	return &triePath{
		children: make(map[string]*triePath),
	}
}

// Get returns the value stored at the given key. Returns nil for internal
// nodes or for nodes with a value of nil.
func (trie *triePath) Get(key string) eventObservers {
	node := trie
	for part, i := segmentPath(key, 0); ; part, i = segmentPath(key, i) {
		node = node.children[part]
		if node == nil {
			return nil
		}
		if i == -1 {
			break
		}
	}
	return node.value
}

func (trie *triePath) dispatch(p *Path, found bool, v []byte) (_ []byte, err error) {
	if trie == nil {
		return v, nil
	}
	node := trie
	key := p.separatorPrefixRoute()
	for part, i := segmentPath(key, 0); ; part, i = segmentPath(key, i) {
		node = node.children[part]
		if node == nil {
			return v, nil
		}

		if v, err = node.value.dispatch(p, found, v); err != nil {
			return nil, errors.WithStack(err)
		}

		if i == -1 {
			break
		}
	}
	return v, nil
}

// Put inserts the value into the trie at the given key, replacing any
// existing items. It returns true if the put adds a new value, false
// if it replaces an existing value.
// Note that internal nodes have nil values so a stored nil value will not
// be distinguishable and will not be included in Walks.
func (trie *triePath) Put(key string, value EventObserver) bool {
	if value == nil {
		return true
	}
	node := trie
	for part, i := segmentPath(key, 0); ; part, i = segmentPath(key, i) {
		child := node.children[part]
		if child == nil {
			child = newTriePath()
			node.children[part] = child
		}
		node = child
		if i == -1 {
			break
		}
	}
	// does node have an existing value?
	isNewVal := node.value == nil
	node.value = append(node.value, value)
	return isNewVal
}

// Delete removes the value associated with the given key. Returns true if a
// node was found for the given key. If the node or any of its ancestors
// becomes childless as a result, it is removed from the trie.
func (trie *triePath) Delete(key string) bool {
	var path []nodeStr // record ancestors to check later
	node := trie
	for part, i := segmentPath(key, 0); ; part, i = segmentPath(key, i) {
		path = append(path, nodeStr{part: part, node: node})
		node = node.children[part]
		if node == nil {
			// node does not exist
			return false
		}
		if i == -1 {
			break
		}
	}
	// delete the node value
	node.value = nil
	// if leaf, remove it from its parent's children map. Repeat for ancestor path.
	if node.isLeaf() {
		// iterate backwards over path
		for i := len(path) - 1; i >= 0; i-- {
			parent := path[i].node
			part := path[i].part
			delete(parent.children, part)
			if parent.value != nil || !parent.isLeaf() {
				// parent has a value or has other children, stop
				break
			}
		}
	}
	return true // node (internal or not) existed and its value was nil'd
}

// Walk iterates over each key/value stored in the trie and calls the given
// walker function with the key and value. If the walker function returns
// an error, the walk is aborted.
// The traversal is depth first with no guaranteed order.
func (trie *triePath) Walk(walker walkFn) error {
	return trie.walk("", walker)
}

// triePath node and the part string key of the child the path descends into.
type nodeStr struct {
	node *triePath
	part string
}

func (trie *triePath) walk(key string, walker walkFn) error {
	if trie.value != nil {
		walker(key, trie.value)
	}
	for part, child := range trie.children {
		err := child.walk(key+part, walker)
		if err != nil {
			return err
		}
	}
	return nil
}

func (trie *triePath) isLeaf() bool {
	return len(trie.children) == 0
}
