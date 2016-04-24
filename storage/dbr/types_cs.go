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

package dbr

import (
	"fmt"
	"strings"
	// "github.com/ugorji/go/codec"
)

// removed because the binary grows for 5MB. will be added once needed.
// var _ codec.Selfer = (*NullString)(nil)

// GoString satisfies the interface fmt.GoStringer when using %#v in Printf methods.
// Returns
// 		dbr.NewNullString(`...`,bool)
func (ns NullString) GoString() string {
	if ns.Valid && strings.ContainsRune(ns.String, '`') {
		// `This is my`string`
		ns.String = strings.Join(strings.Split(ns.String, "`"), "`+\"`\"+`")
		// `This is my`+"`"+`string`
	}

	ns.String = "`" + ns.String + "`"
	if !ns.Valid {
		ns.String = "nil"
	}
	return fmt.Sprintf("dbr.NewNullString(%s)", ns.String)
}

//// CodecEncodeSelf for ugorji.go codec package
//func (n NullString) CodecEncodeSelf(e *codec.Encoder) {
//	if err := e.Encode(n.String); err != nil {
//		PkgLog.Debug("dbr.NullString.CodecEncodeSelf", "err", err, "n", n)
//	}
//}
//
//// CodecDecodeSelf  for ugorji.go codec package @todo write test ... not sure if ok
//func (n *NullString) CodecDecodeSelf(d *codec.Decoder) {
//	if err := d.Decode(&n.String); err != nil {
//		PkgLog.Debug("dbr.NullString.CodecDecodeSelf", "err", err, "n", n)
//	}
//	// think about empty string and Valid value ...
//}
//
//// CodecEncodeSelf for ugorji.go codec package
//func (n *NullInt64) CodecEncodeSelf(e *codec.Encoder) {
//	if err := e.Encode(n.Int64); err != nil {
//		PkgLog.Debug("dbr.NullInt64.CodecEncodeSelf", "err", err, "n", n)
//	}
//}
//
//// CodecDecodeSelf  for ugorji.go codec package @todo write test ... not sure if ok
//func (n *NullInt64) CodecDecodeSelf(d *codec.Decoder) {
//	if err := d.Decode(&n.Int64); err != nil {
//		PkgLog.Debug("dbr.NullInt64.CodecDecodeSelf", "err", err, "n", n)
//	}
//}
//
//// CodecEncodeSelf for ugorji.go codec package
//func (n NullFloat64) CodecEncodeSelf(e *codec.Encoder) {
//	if err := e.Encode(n.Float64); err != nil {
//		PkgLog.Debug("dbr.NullFloat64.CodecEncodeSelf", "err", err, "n", n)
//	}
//}
//
//// CodecDecodeSelf  for ugorji.go codec package @todo write test ... not sure if ok
//func (n *NullFloat64) CodecDecodeSelf(d *codec.Decoder) {
//	if err := d.Decode(&n.Float64); err != nil {
//		PkgLog.Debug("dbr.NullFloat64.CodecDecodeSelf", "err", err, "n", n)
//	}
//}
//
//// CodecEncodeSelf for ugorji.go codec package
//func (n NullTime) CodecEncodeSelf(e *codec.Encoder) {
//	if err := e.Encode(n.Time); err != nil {
//		PkgLog.Debug("dbr.NullTime.CodecEncodeSelf", "err", err, "n", n)
//	}
//}
//
//// CodecDecodeSelf  for ugorji.go codec package @todo write test ... not sure if ok
//func (n *NullTime) CodecDecodeSelf(d *codec.Decoder) {
//	if err := d.Decode(&n.Time); err != nil {
//		PkgLog.Debug("dbr.NullTime.CodecDecodeSelf", "err", err, "n", n)
//	}
//}
//
//// CodecEncodeSelf for ugorji.go codec package
//func (n NullBool) CodecEncodeSelf(e *codec.Encoder) {
//	if err := e.Encode(n.Bool); err != nil {
//		PkgLog.Debug("dbr.NullBool.CodecEncodeSelf", "err", err, "n", n)
//	}
//}
//
//// CodecDecodeSelf  for ugorji.go codec package @todo write test ... not sure if ok
//func (n *NullBool) CodecDecodeSelf(d *codec.Decoder) {
//	if err := d.Decode(&n.Bool); err != nil {
//		PkgLog.Debug("dbr.NullBool.CodecDecodeSelf", "err", err, "n", n)
//	}
//}
