package cstime

import (
	"sync"
)

// lockedSource is an implementation of Source that is concurrency-safe.
// It is just a standard Source with its operations protected by a sync.Mutex.
type lockedSource struct {
	lk  sync.Mutex
	src interface {
		Seed(seed int64)
		Int63n(n int64) int64
	}
}

func (s *lockedSource) Int63n(n int64) int64 {
	s.lk.Lock()
	n = s.src.Int63n(n)
	s.lk.Unlock()
	return n
}

func (s *lockedSource) Seed(seed int64) {
	s.lk.Lock()
	s.src.Seed(seed)
	s.lk.Unlock()
}
