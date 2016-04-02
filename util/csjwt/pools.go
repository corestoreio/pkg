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

package csjwt

import (
	"crypto"
	"crypto/hmac"
	"hash"
	"sync"

	"github.com/corestoreio/csfw/util/bufferpool"
)

var bufPool = bufferpool.New(8192) // estimated *cough* average size of JWT 8kb

type hmacTank struct {
	p *sync.Pool
}

func (t hmacTank) get() hash.Hash {
	return t.p.Get().(hash.Hash)
}

func (t hmacTank) put(h hash.Hash) {
	h.Reset()
	t.p.Put(h)
}

func newHMACTank(ch crypto.Hash, key []byte) hmacTank {
	return hmacTank{
		p: &sync.Pool{
			New: func() interface{} {
				return hmac.New(ch.New, key)
			},
		},
	}
}
