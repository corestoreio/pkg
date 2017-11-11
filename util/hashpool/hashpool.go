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

package hashpool

import (
	"crypto/subtle"
	"encoding/hex"
	"hash"
	"io"
	"sync"

	"github.com/corestoreio/cspkg/util/bufferpool"
	"github.com/corestoreio/errors"
)

// Hash64Mock allows to use a hash.Hash as an argument to the Hash64 Tank.
type Hash64Mock struct {
	hash.Hash
}

// Sum64 returns always zero.
func (d Hash64Mock) Sum64() uint64 {
	return 0
}

// Tank implements a sync.Pool for hash.Hash
type Tank struct {
	p *sync.Pool
	// BufferSize used in SumBase64() to append the hashed data to. Default 1024.
	BufferSize int
}

// Get returns type safe a hash.
func (t Tank) Get() hash.Hash64 {
	return t.p.Get().(hash.Hash64)
}

// Sum calculates the hash of data and appends the current hash to appendTo and
// returns the resulting slice. It does not change the underlying hash state. It
// fetches a hash from the pool and returns it after writing the sum. No need to
// call Get() and Put().
func (t Tank) Sum(data, appendTo []byte) []byte {
	h := t.Get()
	defer t.Put(h)
	_, _ = h.Write(data)
	return h.Sum(appendTo)
}

// SumHex writes the hashed data into the hex encoder.
func (t Tank) SumHex(data []byte) string {
	buf := bufferpool.Get()
	defer bufferpool.Put(buf)
	bs := 1024
	if t.BufferSize > 0 {
		bs = t.BufferSize
	}
	buf.Grow(bs)
	tmpBuf := t.Sum(data, buf.Bytes())
	buf.Reset()
	_, _ = buf.Write(tmpBuf)
	return hex.EncodeToString(buf.Bytes())
}

// Equal hashes data and compares it with MAC for equality without leaking
// timing information.
func (t Tank) Equal(data []byte, mac []byte) bool {
	buf := bufferpool.Get()
	defer bufferpool.Put(buf)
	// We don't have to be constant time if the lengths of the MACs are
	// different as that suggests that a completely different hash function
	// was used.
	return subtle.ConstantTimeCompare(t.Sum(data, buf.Bytes()), mac) == 1
}

// EqualPairs compares data pairs for equality via constant time. Returns true
// only if all pairs are equal and one data item length is at least greater than
// zero. Use case: Compare username and password from database with a username
// and password from an outside input form.
func (t Tank) EqualPairs(dataPairs ...[]byte) bool {
	ldp := len(dataPairs)
	if ldp == 0 || ldp%2 == 1 {
		return false
	}

	buf1 := bufferpool.Get()
	defer bufferpool.Put(buf1)
	buf2 := bufferpool.Get()
	defer bufferpool.Put(buf2)

	eq := 0
	for i := 0; i <= ldp/2; i = i + 2 {
		if len(dataPairs[i]) == 0 {
			return false
		}
		if subtle.ConstantTimeCompare(t.Sum(dataPairs[i], buf1.Bytes()), t.Sum(dataPairs[i+1], buf2.Bytes())) == 1 {
			eq++
		}
		buf1.Reset()
		buf2.Reset()
	}
	return eq == ldp/2
}

// EqualReader hashes io.Reader and compares it with MAC for equality without
// leaking timing information. The internal buffer to read into data from
// io.Reader can be adjusted via field BufferSize.
func (t Tank) EqualReader(r io.Reader, mac []byte) (bool, error) {
	h := t.Get()
	defer t.Put(h)

	bs := 4096
	if t.BufferSize > 0 {
		bs = t.BufferSize
	}
	buf := make([]byte, bs)
	for {
		n, err := r.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			return false, errors.Wrap(err, "[hashpool] r.Read")
		}
		if _, err := h.Write(buf[:n]); err != nil {
			return false, errors.Wrap(err, "[hashpool] Hash.Write")
		}
	}

	// We don't have to be constant time if the lengths of the MACs are
	// different as that suggests that a completely different hash function
	// was used.
	return subtle.ConstantTimeCompare(h.Sum(buf[:0]), mac) == 1, nil
}

// Put empties the hash and returns it back to the pool.
//
//		hp := New(fnv.New64)
//		hsh := hp.Get()
//		defer hp.Put(hsh)
//		// your code
//		return hsh.Sum([]byte{})
//
func (t Tank) Put(h hash.Hash64) {
	h.Reset()
	t.p.Put(h)
}

// New64 instantiates a new hash pool with a custom pre-allocated hash.Hash64.
func New64(h func() hash.Hash64) Tank {
	return Tank{
		p: &sync.Pool{
			New: func() interface{} {
				nh := h()
				nh.Reset()
				return nh
			},
		},
	}
}

// New instantiates a new hash pool with a custom pre-allocated hash.Hash.
func New(h func() hash.Hash) Tank {
	return Tank{
		p: &sync.Pool{
			New: func() interface{} {
				nh := h()
				nh.Reset()
				return Hash64Mock{Hash: nh}
			},
		},
	}
}
