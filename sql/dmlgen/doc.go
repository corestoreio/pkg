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

// Package dmlgen provides code generation templates and library code for
// sql/dml.
//
// To generated the protocol buffer file
// $ protoc --gogo_out=Mgoogle/protobuf/timestamp.proto=github.com/gogo/protobuf/types:. --proto_path=/Users/kiri/GoPro/src/:/Users/kiri/GoPro/src/github.com/gogo/protobuf/protobuf/:. *.proto
//
// TODO: Generate also protobuf code for https://github.com/twitchtv/twirp/wiki
package dmlgen
