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

package path

import "errors"

// ErrPathTooLong ...
var ErrPathTooLong = errors.New("path.pathBuf: buffer too small and path too long")

const pathBufLen = 256

// pathBuf defines the maximum length of a path. This type is equal to database
// table core_config_data column path varchar(255)
// A valid is defined in Path.IsValid
type pathBuf struct {
	pos  int
	data [pathBufLen]byte
}

func newPathBuf() *pathBuf {
	return &pathBuf{}
}

func (pb *pathBuf) Bytes() []byte {
	return pb.data[:pb.pos]
}

func (pb *pathBuf) Reset() {
	pb.pos = 0
}

func (pb *pathBuf) WriteByte(c byte) error {
	if c < 1 {
		return nil
	}
	if pb.pos+1 > pathBufLen {
		return ErrPathTooLong
	}
	pb.data[pb.pos] = c
	pb.pos++
	return nil
}

func (pb *pathBuf) WriteString(s string) (n int, err error) {
	if len(s) < 1 {
		return
	}
	if pb.pos+len(s) > pathBufLen {
		err = ErrPathTooLong
		return
	}
	n = copy(pb.data[pb.pos:], s)
	pb.pos += n
	return n, nil
}

func (pb *pathBuf) Write(p []byte) (n int, err error) {
	if len(p) < 1 {
		return
	}
	if pb.pos+len(p) > pathBufLen {
		err = ErrPathTooLong
		return
	}
	n = copy(pb.data[pb.pos:], p)
	pb.pos += n
	return n, nil
}
