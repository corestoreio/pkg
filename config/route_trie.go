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
	"github.com/corestoreio/pkg/store/scope"
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
// support modifying the raw data argument. For example the EventOnAfterGet
// allows to decrypt encrypted data. Or check if some one is allowed to read or
// write a special path with its value. Or validate data for correctness.
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
type walkFn func(key string, value fieldMeta) error

// segmentRoute segments string key paths by slash separators. For example,
// "/a/b/c" -> ("/a", 2), ("/b", 4), ("/c", -1) in successive calls. It does
// not allocate any heap memory.
func segmentRoute(path string, start int) (segment string, next int) {
	if len(path) == 0 || start < 0 || start > len(path)-1 {
		return "", -1
	}
	end := strings.IndexRune(path[start+1:], '/') // next '/' after 0th rune
	if end == -1 {
		return path[start:], -1
	}
	return path[start : start+end+1], start + end + 1
}

// trieRoute is a trie of paths with string keys and interface{} values.

// trieRoute is a trie of string keys and interface{} values. Internal nodes
// have nil values so stored nil values cannot be distinguished and are
// excluded from walks. By default, trieRoute will segment keys by forward
// slashes with segmentRoute (e.g. "/a/b/c" -> "/a", "/b", "/c"). A custom
// segmentStringFn may be used to customize how strings are segmented into
// nodes. A classic trie might segment keys by rune (i.e. unicode points).
type trieRoute struct {
	fm       fieldMeta
	children map[string]*trieRoute
}

// newTrieRoute allocates and returns a new *trieRoute.
func newTrieRoute() *trieRoute {
	return &trieRoute{
		children: make(map[string]*trieRoute),
	}
}

// Get returns the fm stored at the given key. Returns nil for internal
// nodes or for nodes with a fm of nil.
func (trie *trieRoute) Get(key string) fieldMeta {
	if !strings.HasPrefix(key, sPathSeparator) {
		key = sPathSeparator + key
	}
	node := trie
	for part, i := segmentRoute(key, 0); ; part, i = segmentRoute(key, i) {
		node = node.children[part]
		if node == nil {
			return fieldMeta{}
		}
		if i == -1 {
			break
		}
	}
	return node.fm
}

// process runs on each tree level and dispatches the events and checks for
// scope permission and default value.
func (trie *trieRoute) process(key string, event uint8, p *Path, found bool, v []byte) (_ []byte, err error) {
	if trie == nil {
		return v, nil
	}
	if !strings.HasPrefix(key, sPathSeparator) {
		key = sPathSeparator + key
	}
	node := trie
	for part, i := segmentRoute(key, 0); ; part, i = segmentRoute(key, i) {
		node = node.children[part]

		if node == nil {
			return v, nil
		}
		if node.fm.valid && node.fm.ScopePerm > 0 && p.ScopeID > 0 && !node.fm.ScopePerm.Has(p.ScopeID.Type()) {
			return nil, errors.NotAllowed.Newf("[config] The path %q is not allowed to access this scope %s", p.String(), node.fm.ScopePerm.String())
		}

		if v, err = node.fm.Events[event].dispatch(p, found, v); err != nil {
			return nil, errors.WithStack(err)
		}
		if node.fm.valid && event == EventOnAfterGet && !found && v == nil && node.fm.Default != "" { // dispatch should also return `found`
			v = []byte(node.fm.Default)
		}

		if i == -1 {
			break
		}
	}
	return v, nil
}

// PutEvent inserts the fm into the trie at the given key, replacing any
// existing items. It returns true if the put adds a new fm, false
// if it replaces an existing fm.
// Note that internal nodes have nil values so a stored nil fm will not
// be distinguishable and will not be included in Walks.
func (trie *trieRoute) PutEvent(event uint8, key string, eo EventObserver) bool {
	if eo == nil {
		return true
	}

	if !strings.HasPrefix(key, sPathSeparator) {
		key = sPathSeparator + key
	}

	node := trie
	for part, i := segmentRoute(key, 0); ; part, i = segmentRoute(key, i) {
		child := node.children[part]
		if child == nil {
			child = newTrieRoute()
			node.children[part] = child
		}
		node = child
		if i == -1 {
			break
		}
	}
	// does node have an existing eo?
	isNewVal := node.fm.valid == false
	node.fm.Events[event] = append(node.fm.Events[event], eo)
	node.fm.valid = true
	return isNewVal
}

func (trie *trieRoute) PutMeta(key string, perm scope.Perm, defaultVal string) bool {
	if !strings.HasPrefix(key, sPathSeparator) {
		key = sPathSeparator + key
	}

	node := trie
	for part, i := segmentRoute(key, 0); ; part, i = segmentRoute(key, i) {
		child := node.children[part]
		if child == nil {
			child = newTrieRoute()
			node.children[part] = child
		}
		node = child
		if i == -1 {
			break
		}
	}
	// does node have an existing eo?
	isNewVal := node.fm.valid == false
	node.fm.ScopePerm = perm
	node.fm.Default = defaultVal
	node.fm.valid = true
	return isNewVal
}

// DeleteEvent removes the fm associated with the given key. Returns true if a
// node was found for the given key. If the node or any of its ancestors becomes
// childless as a result, it is removed from the trie.
func (trie *trieRoute) DeleteEvent(event uint8, key string) bool {
	if !strings.HasPrefix(key, sPathSeparator) {
		key = sPathSeparator + key
	}
	var path []nodeStr // record ancestors to check later
	node := trie
	for part, i := segmentRoute(key, 0); ; part, i = segmentRoute(key, i) {
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
	// delete the node fm
	node.fm.Events[event] = nil
	// if leaf, remove it from its parent's children map. Repeat for ancestor path.

	return true // node (internal or not) existed and its fm was nil'd
}

// DeleteEvent removes the fm associated with the given key. Returns true if a
// node was found for the given key. If the node or any of its ancestors becomes
// childless as a result, it is removed from the trie.
func (trie *trieRoute) Delete(key string) bool {
	if !strings.HasPrefix(key, sPathSeparator) {
		key = sPathSeparator + key
	}
	var path []nodeStr // record ancestors to check later
	node := trie
	for part, i := segmentRoute(key, 0); ; part, i = segmentRoute(key, i) {
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
	// delete the node fm
	node.fm = fieldMeta{}
	// if leaf, remove it from its parent's children map. Repeat for ancestor path.
	if node.isLeaf() {
		// iterate backwards over path
		for i := len(path) - 1; i >= 0; i-- {
			parent := path[i].node
			part := path[i].part
			delete(parent.children, part)
			if parent.fm.valid || !parent.isLeaf() {
				// parent has a fm or has other children, stop
				break
			}
		}
	}
	return true // node (internal or not) existed and its fm was nil'd
}

// Walk iterates over each key/value stored in the trie and calls the given
// walker function with the key and fm. If the walker function returns
// an error, the walk is aborted.
// The traversal is depth first with no guaranteed order.
func (trie *trieRoute) Walk(walker walkFn) error {
	return trie.walk("", walker)
}

// trieRoute node and the part string key of the child the path descends into.
type nodeStr struct {
	node *trieRoute
	part string
}

func (trie *trieRoute) walk(key string, walker walkFn) error {
	if trie.fm.valid {
		walker(key, trie.fm)
	}
	for part, child := range trie.children {
		err := child.walk(key+part, walker)
		if err != nil {
			return err
		}
	}
	return nil
}

func (trie *trieRoute) isLeaf() bool {
	return len(trie.children) == 0
}
