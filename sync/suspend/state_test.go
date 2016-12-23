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

package suspend_test

import (
	"crypto/rand"
	"hash"
	"hash/fnv"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/corestoreio/csfw/store/scope"
	"github.com/corestoreio/csfw/sync/suspend"
	"github.com/corestoreio/errors"
)

const hashCount = 10

func generateHashes() (hashes [hashCount + 1]scope.TypeID) {
	const scopeCount = 2
	scopes := [scopeCount]scope.Type{
		scope.Website, scope.Store,
	}
	for i := 0; i < hashCount; i++ {
		s := scopes[i%2]
		hashes[i] = scope.MakeTypeID(s, int64(i))
	}
	hashes[hashCount] = scope.DefaultTypeID
	return hashes
}

func TestHashState_Initialized(t *testing.T) {
	var hs suspend.State
	if hs.Initialized() {
		t.Error("Should not be initialized")
	}
	if err := hs.Done(scope.DefaultTypeID.ToUint64()); !errors.IsFatal(err) {
		t.Errorf("Expecting a Fatal error: %s", err)
	}
	if hs.ShouldWait(scope.DefaultTypeID.ToUint64()) {
		t.Error("Should not be running")
	}

	hs = suspend.NewState()
	if !hs.Initialized() {
		t.Error("Should be initialized")
	}
	if err := hs.Done(scope.DefaultTypeID.ToUint64()); !errors.IsFatal(err) {
		t.Errorf("Expecting a Fatal error: %s", err)
	}
}

func TestHashState_CanRun(t *testing.T) {
	var hs suspend.State

	if hs.ShouldStart(scope.DefaultTypeID.ToUint64()) {
		t.Fatal("Should not return true because not yet initialized")
	}
	hs = suspend.NewState()

	hashes := generateHashes()

	t.Run("Start", func(t *testing.T) {
		var wg sync.WaitGroup
		for i := 0; i < len(hashes); i++ {
			wg.Add(1)
			go func(wg *sync.WaitGroup, i int) {
				defer wg.Done()
				h := hashes[i]
				if hs.ShouldWait(h.ToUint64()) {
					t.Fatal("Should not run")
				}
				if !hs.ShouldStart(h.ToUint64()) {
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
				if hs.ShouldStart(h.ToUint64()) {
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

func TestHashState_ShouldWait(t *testing.T) {
	var hs suspend.State
	if hs.ShouldStart(scope.DefaultTypeID.ToUint64()) {
		t.Fatal("Should not return true because not yet initialized")
	}

	hs = suspend.NewState()

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
				// You should stick to the sequence of the cases: 1. Start 2. Running
				case hs.ShouldStart(h.ToUint64()):
					atomic.AddUint32(countRun, 1)
					time.Sleep(time.Millisecond) // simulate race detector ;-)
					if err := hs.Done(h.ToUint64()); err != nil {
						t.Fatalf("%+v", err)
					}
				case hs.ShouldWait(h.ToUint64()): // this case is normally not needed
					atomic.AddUint32(countWaiter, 1)
				}

				if hs.ShouldWait(h.ToUint64()) {
					t.Fatalf("Should not be running, because it should have ran "+
						"already, so now we're waiting. Iteration %d Hash %s", j, h)
				}
			}(&wg, i, j)
		}
	}
	wg.Wait()

	if have, want := atomic.LoadUint32(countRun), uint32(len(hashes)); have != want {
		t.Errorf("Runner: Have %d Want %d", have, want)
	}
	if have, want := atomic.LoadUint32(countWaiter), uint32(140); have < want {
		// we have around ~50-190 waiting goroutines. the number of waiting goroutines
		// is totally random and depends on how fast the case CanRun: can process the
		// work or if the race detector has been enabled. So we check here for at least
		// 100 sleeper.
		t.Errorf("Waiter: Have %d > Want %d", have, want)
	} else {
		t.Logf("INFO Waiter: %d", have)
	}

	hs.Reset()
	if have, want := hs.Len(), 0; have != want {
		t.Errorf("After calling Reset the internal map must be cleared. Have %d", have)
	}
}

func TestHashStateBytes_ShouldWait(t *testing.T) {

	hs := suspend.NewStateWithHash(fnv.New64a())

	const randBytesGenerate = 11

	const iterations = 5 // just a random number picked, this would be the expect # of http.requests
	var countWaiter = new(uint32)
	var countRun = new(uint32)
	var wg sync.WaitGroup

	for i := 0; i < randBytesGenerate; i++ {
		var key = make([]byte, 32)
		if n, err := rand.Read(key); err != nil || n == 0 {
			t.Fatalf("Failed to random Bytes: %s Read Count: %d", err, n)
		}

		for j := 0; j < iterations; j++ {
			wg.Add(1)
			go func(wg *sync.WaitGroup, i, j int, key []byte) {
				defer wg.Done()
				switch {
				// You should stick to the sequence of the cases: 1. Start 2. Running
				case hs.ShouldStartBytes(key):
					atomic.AddUint32(countRun, 1)
					time.Sleep(time.Millisecond) // simulate race detector ;-)
					if err := hs.DoneBytes(key); err != nil {
						t.Fatalf("%+v", err)
					}
				case hs.ShouldWaitBytes(key): // this case is normally not needed
					atomic.AddUint32(countWaiter, 1)
				}

				if hs.ShouldWaitBytes(key) {
					t.Fatalf("Should not be running, because it should have ran "+
						"already, so now we're waiting. Iteration %d Hash %x", j, key)
				}
			}(&wg, i, j, key)
		}
	}
	wg.Wait()

	if have, want := atomic.LoadUint32(countRun), uint32(randBytesGenerate); have != want {
		t.Errorf("Runner: Have %d Want %d", have, want)
	}
	if have, want := atomic.LoadUint32(countWaiter), uint32(iterations*4); have < want {
		t.Errorf("Waiter: Have %d > Want %d", have, want)
	} else {
		t.Logf("INFO Waiter: %d", have)
	}

	hs.Reset()
	if have, want := hs.Len(), 0; have != want {
		t.Errorf("After calling Reset the internal map must be cleared. Have %d", have)
	}
}

var _ hash.Hash64 = (*hashMock)(nil)

type hashMock struct {
	err error
	hash.Hash64
}

func (hm hashMock) Write(p []byte) (n int, err error) {
	return 0, hm.err
}
func (hm hashMock) Reset()        {}
func (hm hashMock) Sum64() uint64 { return 0 }

func TestHashStateBytes_Error(t *testing.T) {

	hs := suspend.NewStateWithHash(hashMock{
		err: errors.NewWriteFailedf("Write failed"),
	})
	key := []byte(`I'm writing Go code`)
	if hs.ShouldWaitBytes(key) {
		t.Error("Should not wait")
	}
	if have, want := hs.Len(), 0; have != want {
		t.Errorf("Have: %v Want: %v", have, want)
	}

}
