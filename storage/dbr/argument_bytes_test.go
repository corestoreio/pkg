// Copyright 2015-2017, Cyrill @ Schumacher.fm and the CoreStore contributors
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

package dbr

import (
	"github.com/stretchr/testify/assert"
	"runtime"
	"testing"
)

func TestArgBytes(t *testing.T) {
	t.Parallel()

	var ac argBytes
	t.Run("simple case, immutable", func(t *testing.T) {
		ac = makeArgBytes().
			appendInt64s(123, 456, 89, 12, 13, 14, 15, 16, 17, 18, 19).
			appendNull().
			appendInt64(567).
			appendNull()
		// Call it twice to shoe immutability.
		assert.Exactly(t, `0:{123,456,89,12,13,14,15,16,17,18,19} 1:{NULL} 2:{567} 3:{NULL}`, ac.DebugBytes())
		assert.Exactly(t, `0:{123,456,89,12,13,14,15,16,17,18,19} 1:{NULL} 2:{567} 3:{NULL}`, ac.DebugBytes())
		assert.Len(t, ac, 4)
		assert.Exactly(t, argBytesCap, cap(ac), "cap should be unchanged")
	})
	t.Run("reset", func(t *testing.T) {
		ac = ac.reset()
		assert.Len(t, ac, 0)
		assert.Exactly(t, argBytesCap, cap(ac), "cap should be unchanged")
	})
	t.Run("allocate no new memory but use different args", func(t *testing.T) {
		msBefore := new(runtime.MemStats)
		runtime.ReadMemStats(msBefore)

		ac = ac.
			appendInt64s(12, 13, 14, 15, 16, 17, 18, 19, 123, 456, 89).
			appendNull().
			appendInt64(765).
			appendNull()

		msAfter := new(runtime.MemStats)
		runtime.ReadMemStats(msAfter)
		assert.Exactly(t, "0:{12,13,14,15,16,17,18,19,123,456,89} 1:{NULL} 2:{765} 3:{NULL}", ac.DebugBytes())
		assert.Exactly(t, "0:{12,13,14,15,16,17,18,19,123,456,89} 1:{NULL} 2:{765} 3:{NULL}", ac.DebugBytes())
		assert.Len(t, ac, 4)
		assert.Exactly(t, argBytesCap, cap(ac), "cap should be unchanged")

		assert.Empty(t, msAfter.Alloc-msBefore.Alloc, "Alloc should be zero")
		assert.Empty(t, msAfter.TotalAlloc-msBefore.TotalAlloc, "TotalAlloc should be zero")
		assert.Empty(t, msAfter.Mallocs-msBefore.Mallocs, "Mallocs should be zero")
		//t.Logf("Alloc %d", msAfter.Alloc-msBefore.Alloc)
		//t.Logf("TotalAlloc %d", msAfter.TotalAlloc-msBefore.TotalAlloc)
		//t.Logf("Mallocs %d", msAfter.Mallocs-msBefore.Mallocs)
	})

	t.Run("allocate 5 new blocks but use different args", func(t *testing.T) {
		ac = ac.reset()
		msBefore := new(runtime.MemStats)
		runtime.ReadMemStats(msBefore)

		ac = ac.
			appendInt64s(12, 13, 14, 15, 16, 17, 18, 19, 123, 456, 89).
			appendNull().
			appendInt64(765).
			appendInt64(34).
			appendInt64(35).
			appendInt64(36).
			appendInt64(37).
			appendInt64(38).
			appendInt64(39).
			appendNull()

		msAfter := new(runtime.MemStats)
		runtime.ReadMemStats(msAfter)
		assert.Exactly(t, "0:{12,13,14,15,16,17,18,19,123,456,89} 1:{NULL} 2:{765} 3:{34} 4:{35} 5:{36} 6:{37} 7:{38} 8:{39} 9:{NULL}", ac.DebugBytes())
		assert.Exactly(t, "0:{12,13,14,15,16,17,18,19,123,456,89} 1:{NULL} 2:{765} 3:{34} 4:{35} 5:{36} 6:{37} 7:{38} 8:{39} 9:{NULL}", ac.DebugBytes())
		assert.Len(t, ac, 10)
		assert.Exactly(t, 2*argBytesCap, cap(ac), "cap should be doubled")

		assert.Exactly(t, uint64(464), msAfter.Alloc-msBefore.Alloc, "Alloc should be 464")
		assert.Exactly(t, uint64(464), msAfter.TotalAlloc-msBefore.TotalAlloc, "TotalAlloc should be 464")
		assert.Exactly(t, uint64(5), msAfter.Mallocs-msBefore.Mallocs, "Mallocs should be 5")
		//t.Logf("Alloc %d", msAfter.Alloc-msBefore.Alloc)
		//t.Logf("TotalAlloc %d", msAfter.TotalAlloc-msBefore.TotalAlloc)
		//t.Logf("Mallocs %d", msAfter.Mallocs-msBefore.Mallocs)
	})

}
