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

package gzippool

import (
	"compress/gzip"
	"io"
	"sync"
)

var readerPool = &tankReader{}

// GetReader returns a new gzip reader for decompression.
func GetReader(r io.Reader) *gzip.Reader {
	return readerPool.Get(r)
}

// PutReader returns a reader to the pool.
func PutReader(zr *gzip.Reader) {
	readerPool.Put(zr)
}

type tankReader struct {
	p sync.Pool
}

func (t *tankReader) Get(r io.Reader) (zr *gzip.Reader) {
	if zrr := t.p.Get(); zrr != nil {
		zr = zrr.(*gzip.Reader)
		zr.Reset(r)
	} else {
		zr, _ = gzip.NewReader(r)
	}
	return zr
}

func (t *tankReader) Put(zr *gzip.Reader) {
	zr.Close()
	t.p.Put(zr)
}
