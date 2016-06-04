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

package scope_test

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/corestoreio/csfw/store/scope"
	"github.com/corestoreio/csfw/util/errors"
)

const hashCount = 10

func generateHashes() (hashes [hashCount + 1]scope.Hash) {
	const scopeCount = 2
	scopes := [scopeCount]scope.Scope{
		scope.Website, scope.Store,
	}
	for i := 0; i < hashCount; i++ {
		s := scopes[i%2]
		hashes[i] = scope.NewHash(s, int64(i))
	}
	hashes[hashCount] = scope.DefaultHash
	return hashes
}

func TestHashState_Initialized(t *testing.T) {
	var hs scope.HashState
	if hs.Initialized() {
		t.Error("Should not be initialized")
	}
	if err := hs.Done(scope.DefaultHash); !errors.IsFatal(err) {
		t.Errorf("Expecting a Fatal error: %s", err)
	}
	if hs.IsRunning(scope.DefaultHash) {
		t.Error("Should not be running")
	}

	hs = scope.NewHashState()
	if !hs.Initialized() {
		t.Error("Should be initialized")
	}
	if err := hs.Done(scope.DefaultHash); !errors.IsFatal(err) {
		t.Errorf("Expecting a Fatal error: %s", err)
	}
}

func TestHashState_CanRun(t *testing.T) {
	var hs scope.HashState
	if hs.ShouldStart(scope.DefaultHash) {
		t.Fatal("Should not return true because not yet initialized")
	}
	hs = scope.NewHashState()

	hashes := generateHashes()

	t.Run("Start", func(t *testing.T) {
		var wg sync.WaitGroup
		for i := 0; i < len(hashes); i++ {
			wg.Add(1)
			go func(wg *sync.WaitGroup, i int) {
				defer wg.Done()
				h := hashes[i]
				if hs.IsRunning(h) {
					t.Fatal("Should not run")
				}
				if !hs.ShouldStart(h) {
					t.Fatal("Should return true that it has been started")
				}
			}(&wg, i)
		}
		wg.Wait()
	})

	t.Run("NoStart", func(t *testing.T) {
		var wg sync.WaitGroup
		for i := 0; i < len(hashes); i++ {
			wg.Add(1)
			go func(wg *sync.WaitGroup, i int) {
				defer wg.Done()
				h := hashes[i]
				if hs.ShouldStart(h) {
					t.Fatal("Multiple starts not allowed")
				}
			}(&wg, i)
		}
		wg.Wait()
	})

	if have, want := hs.Len(), len(hashes); have != want {
		t.Errorf("Incorrect HashState Length: Have %d Want %d", have, want)
	}
}

func TestHashState_IsRunning(t *testing.T) {
	var hs scope.HashState
	if hs.ShouldStart(scope.DefaultHash) {
		t.Fatal("Should not return true because not yet initialized")
	}

	hs = scope.NewHashState()

	hashes := generateHashes()

	const iterations = 23 // just a random number picked, this would be the expect # of http.requests
	var countWaiter = new(uint32)
	var countRun = new(uint32)
	var wg sync.WaitGroup

	for i := 0; i < len(hashes); i++ {
		for j := 0; j < iterations; j++ {
			wg.Add(1)
			go func(wg *sync.WaitGroup, i, j int) {
				defer wg.Done()
				h := hashes[i]

				switch {
				case hs.IsRunning(h):
					// wait here
					//t.Logf("IsRunning Iteration %d Hash %s", j, h)
					atomic.AddUint32(countWaiter, 1)
				case hs.ShouldStart(h):
					atomic.AddUint32(countRun, 1)
					time.Sleep(time.Millisecond) // simulate race detector ;-)
					if err := hs.Done(h); err != nil {
						t.Fatal(errors.PrintLoc(err))
					}
				}

				if hs.IsRunning(h) {
					// This is weird because sometimes it prints out the log and
					// sometimes nothing gets printed. The race detector detects
					// nothing. AFAIK there must be delay between the first
					// switch case and the last switch case. but if we comment
					// out the whole switch everything still works as expected.
					t.Logf("Should not be running, because it should have ran "+
						"already, so now we're waiting. Iteration %d Hash %s", j, h)
				}
			}(&wg, i, j)
		}
	}
	wg.Wait()

	if have, want := atomic.LoadUint32(countRun), uint32(len(hashes)); have != want {
		t.Errorf("Runner: Have %d Want %d", have, want)
	}
	if have, want := atomic.LoadUint32(countWaiter), uint32(10); have < want {
		// we have around ~50-190 waiting goroutines. the number of waiting goroutines
		// is totally random and depends on how fast the case CanRun: can process the
		// work or if the race detector has been enabled. So we check here for at least
		// 10 sleeper.
		t.Errorf("Waiter: Have %d > Want %d", have, want)
	} else {
		t.Logf("INFO Waiter: %d", have)
	}

	hs.Reset()
	if have, want := hs.Len(), 0; have != want {
		t.Errorf("After calling Reset the internal map must be cleared. Have %d", have)
	}
}
