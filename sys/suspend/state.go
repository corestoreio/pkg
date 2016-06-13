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

// Package suspend provides a complicated duplicate function call suppression
// mechanism.
//
// It picks from n-Goroutines the first one to do a job and suspends
// the following.
//
// The following Goroutines may continue once the first one calls a broadcast
// signal to release the suspended.
package suspend

import (
	"hash"
	"sync"
	"sync/atomic"

	"github.com/corestoreio/csfw/util/errors"
)

// State defines a struct which can atomically and in concurrent cases
// detect if a process runs, idles or has been finished. If another Goroutine
// checks the running state and returns true this Goroutine will be suspended
// and must wait until the main Goroutine calls Done().
//
// Use case for the Hashstate is mainly in middleware Service types for
// net/http.Requests to load configuration values atomically and make other
// requests wait until the configuration has been fully loaded and applied.
// After that skip the whole access to State and use the configuration
// values cached in the middleware service type.
type State struct {
	mu     *sync.RWMutex
	states map[uint64]state
	// Hash64 optional feature to be used when calling the functions which
	// accepts a byte slice as input argument. The byte slice gets hashed into
	// an uint64. If Hash64 remains empty, a panic will pop up. Depending on the
	// hashing algorithm, collisions can occur.
	// http://programmers.stackexchange.com/questions/49550/which-hashing-algorithm-is-best-for-uniqueness-and-speed
	hash.Hash64
}

// NewState creates a new idle State.
func NewState() State {
	return State{
		mu:     &sync.RWMutex{},
		states: make(map[uint64]state),
	}
}

// NewStateWithHash convenience helper function to create a new State with a hash algorithm.
func NewStateWithHash(h hash.Hash64) State {
	s := NewState()
	s.Hash64 = h
	return s
}

// Initialized returns true if the type State has been initialized with the
// call to NewHashState().
func (s State) Initialized() bool {
	return s.mu != nil && s.states != nil
}

// Reset clears the internal list of uint64(es). Will panic if called on an
// uninitialized State.
func (s *State) Reset() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.states = make(map[uint64]state)
}

// Len returns the number of processed Hashes.
func (s State) Len() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.states)
}

func (s State) byteToUint64(key []byte) uint64 {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Hash64.Reset()
	_, _ = s.Hash64.Write(key)
	return s.Hash64.Sum64()
}

// ShouldStartBytes same as ShouldStart.
func (s State) ShouldStartBytes(key []byte) bool {
	return s.ShouldStart(s.byteToUint64(key))
}

// ShouldStart reports true atomically for a specific uint64, if a process can
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
func (s State) ShouldStart(key uint64) bool {
	if !s.Initialized() {
		return false
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	_, ok := s.states[key]
	if !ok {
		s.states[key] = newRunningState()
	}
	return !ok // we've created a new entry in the map and now we can run
}

// ShouldWaitBytes same as ShouldWait.
func (s State) ShouldWaitBytes(key []byte) bool {
	return s.ShouldWait(s.byteToUint64(key))
}

// ShouldWait checks atomically if the State has been set to run and if so
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
func (s State) ShouldWait(key uint64) bool {
	if !s.Initialized() {
		return false
	}
	s.mu.RLock()
	st, ok := s.states[key]
	s.mu.RUnlock()
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

// DoneBytes same as Done.
func (s State) DoneBytes(key []byte) error {
	return errors.Wrap(s.Done(s.byteToUint64(key)), "[suspend] DoneBytes")
}

// Done releases all the waiting Goroutines caught with the function IsRunning()
// and sets the internal state to done. Any subsequent calls to ShouldWait() and
// ShouldStart() will fail. You must call Reset() once all is done.
func (s State) Done(key uint64) error {
	if !s.Initialized() {
		return errors.NewFatalf("[suspend] State not initialized")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	st, ok := s.states[key]
	if !ok {
		return errors.NewFatalf("[suspend] Done(%s) called without calling CanRun(%s)", key, key)
	}
	atomic.StoreUint32(st.status, stateDone)
	st.Broadcast()
	s.states[key] = state{
		// not needed anymore, so set sync.* pointers to nil and copy status
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
