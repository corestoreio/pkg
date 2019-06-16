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

package store_test

import (
	"bytes"
	"encoding/json"

	"github.com/corestoreio/pkg/storage/null"
	"github.com/corestoreio/pkg/store"
)

func init() {
	null.MustSetJSONMarshaler(json.Marshal, json.Unmarshal)
}

/*
std library vs easyjson with json.Encoder
name                     old time/op    new time/op    delta
Service_Json_Encoding-4     251µs ± 1%     272µs ± 1%   +8.34%  (p=0.016 n=5+4)

name                     old alloc/op   new alloc/op   delta
Service_Json_Encoding-4     149kB ± 0%     144kB ± 0%   -3.26%  (p=0.008 n=5+5)

name                     old allocs/op  new allocs/op  delta
Service_Json_Encoding-4       195 ± 0%       149 ± 0%  -23.59%  (p=0.008 n=5+5)


std library vs pure easyjson
name                     old time/op    new time/op    delta
Service_Json_Encoding-4     251µs ± 1%      30µs ± 2%  -88.11%  (p=0.008 n=5+5)

name                     old alloc/op   new alloc/op   delta
Service_Json_Encoding-4     149kB ± 0%      15kB ± 0%  -89.66%  (p=0.008 n=5+5)

name                     old allocs/op  new allocs/op  delta
Service_Json_Encoding-4       195 ± 0%        75 ± 0%  -61.54%  (p=0.008 n=5+5)
*/

// func toJSON2(srv *store.Service) []byte {
// 	jw := &jwriter.Writer{}
//
// 	jw.Buffer.AppendString("[\n")
// 	srv.Websites().MarshalEasyJSON(jw)
// 	jw.Buffer.AppendString(",\n")
// 	srv.Websites().MarshalEasyJSON(jw)
// 	jw.Buffer.AppendString(",\n")
// 	srv.Websites().MarshalEasyJSON(jw)
// 	jw.Buffer.AppendString("]\n")
// 	data, err := jw.BuildBytes()
// 	if err != nil {
// 		panic(err)
// 	}
// 	return data
// }

func toJSON(srv *store.Service) []byte {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetIndent("", "  ")
	buf.WriteString("[\n")
	if err := enc.Encode(srv.Websites()); err != nil {
		panic(err)
	}
	buf.WriteString(",\n")
	if err := enc.Encode(srv.Groups()); err != nil {
		panic(err)
	}
	buf.WriteString(",\n")
	if err := enc.Encode(srv.Stores()); err != nil {
		panic(err)
	}
	buf.WriteString("]\n")
	return buf.Bytes()
}
