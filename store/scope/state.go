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

package scope

import (
	"sync"
	"sync/atomic"

	"github.com/corestoreio/csfw/util/errors"
)

// HashState defines a struct which can atomically and in concurrent cases
// detect if a process runs, idles or has been finished. If another Goroutine
// checks the running state and returns true this Goroutine will be suspended
// and must wait until the main Goroutine calls Done().
//
// Use case for the Hashstate is mainly in middleware Service types for
// net/http.Requests to load configuration values atomically and make other
// requests wait until the configuration has been fully loaded and applied.
// After that skip the whole access to HashState and use the configuration
// values cached in the middleware service type.
type HashState struct {
	mu     *sync.RWMutex
	states map[Hash]state
}

// NewHashState creates a new idle HashState.
func NewHashState() HashState {
	return HashState{
		mu:     &sync.RWMutex{},
		states: make(map[Hash]state),
	}
}

// Initialized returns true if the type HashState has been initialized with the
// call to NewHashState().
func (shs HashState) Initialized() bool {
	return shs.mu != nil && shs.states != nil
}

// Reset clears the internal list of Hash(es). Will panic if called on an
// uninitialized HashState.
func (shs *HashState) Reset() {
	shs.mu.Lock()
	defer shs.mu.Unlock()
	shs.states = make(map[Hash]state)
}

// Len returns the number of processed Hashes.
func (shs HashState) Len() int {
	shs.mu.Lock()
	defer shs.mu.Unlock()
	return len(shs.states)
}

// ShouldStart reports true atomically for a specific Hash, if a process can
// start. Safe for concurrent use. You should check ShouldStart() first in your
// switch statement:
//		switch {
//			case hs.ShouldStart(scopeHash):
//				// start here your process
//				err := hs.Done(scopeHash)
//			case hs.ShouldWait(scopeHash):
//				// do nothing and wait ...
//		}
//		// proceed here with your program
func (shs HashState) ShouldStart(h Hash) bool {
	if !shs.Initialized() {
		return false
	}
	shs.mu.Lock()
	defer shs.mu.Unlock()
	_, ok := shs.states[h]
	if !ok {
		shs.states[h] = newRunningState()
	}
	return !ok // we've created a new entry in the map and now we can run
}

// ShouldWait checks atomically if the HashState has been set to run and if so
// the calling Goroutine waits until Done() has been called. You should use
// ShouldWait() as second case in your switch statement:
//		switch {
//			case hs.ShouldStart(scopeHash):
//				// start here your process
//				err := hs.Done(scopeHash)
//			case hs.ShouldWait(scopeHash):
//				// do nothing and wait ...
//		}
//		// proceed here with your program
func (shs HashState) ShouldWait(h Hash) bool {
	if !shs.Initialized() {
		return false
	}
	shs.mu.RLock()
	st, ok := shs.states[h]
	shs.mu.RUnlock()
	if !ok {
		return false
	}

	if atomic.LoadUint32(st.status) != stateRunning {
		return false
	}

	st.Lock()
	for atomic.LoadUint32(st.status) != stateDone {
		st.Wait()
	}
	st.Unlock()

	return true
}

// Done releases all the waiting Goroutines caught with the function IsRunning()
// and sets the internal state to done. Any subsequent calls to ShouldWait() and
// ShouldStart() will fail. You must call Reset() once all is done.
func (shs HashState) Done(h Hash) error {
	if !shs.Initialized() {
		return errors.NewFatalf("[scope] HashState not initialized")
	}
	shs.mu.Lock()
	defer shs.mu.Unlock()
	st, ok := shs.states[h]
	if !ok {
		return errors.NewFatalf("[scope] Done(%s) called without calling CanRun(%s)", h, h)
	}
	atomic.StoreUint32(st.status, stateDone)
	st.Broadcast()
	shs.states[h] = state{
		// not needed anymore, so set sync.* pointers to nil.
		status: st.status,
	}
	return nil
}

const (
	stateRunning uint32 = 1 << iota
	stateDone
)

type state struct {
	// if manual applied options those three values are nil
	*sync.Mutex
	*sync.Cond
	status *uint32
}

func newRunningState() state {
	var sr = stateRunning
	mu := &sync.Mutex{}
	return state{
		Mutex:  mu,
		Cond:   sync.NewCond(mu),
		status: &sr,
	}
}
