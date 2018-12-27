package pseudo

import (
	"sync"

	"golang.org/x/exp/rand"
)

// lockedSource is an implementation of Source that is concurrency-safe.
// It is just a standard Source with its operations protected by a sync.Mutex.
// CSC: That is stupid to copy the locked source ...
type lockedSource struct {
	lk  sync.Mutex
	src rand.Source
}

func (s *lockedSource) Uint64() (n uint64) {
	s.lk.Lock()
	n = s.src.Uint64()
	s.lk.Unlock()
	return
}

func (s *lockedSource) Seed(seed uint64) {
	s.lk.Lock()
	s.src.Seed(seed)
	s.lk.Unlock()
}

// seedPos implements Seed for a lockedSource without a race condiiton.
func (s *lockedSource) seedPos(seed uint64, readPos *int8) {
	s.lk.Lock()
	s.src.Seed(seed)
	*readPos = 0
	s.lk.Unlock()
}

// Read implements Read for a lockedSource.
func (s *lockedSource) Read(p []byte, readVal *uint64, readPos *int8) (n int, err error) {
	s.lk.Lock()
	n, err = read(p, s.src.Uint64, readVal, readPos)
	s.lk.Unlock()
	return
}

func read(p []byte, uint64 func() uint64, readVal *uint64, readPos *int8) (n int, err error) {
	pos := *readPos
	val := *readVal
	for n = 0; n < len(p); n++ {
		if pos == 0 {
			val = uint64()
			pos = 8
		}
		p[n] = byte(val)
		val >>= 8
		pos--
	}
	*readPos = pos
	*readVal = val
	return
}
