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

/*
Package codegen generates Go code and is only used in development.

Configuration

If you would like to configure the generation process because of custom EAV and attribute models
please create a new file in this folder called `config_*.go` where `*` can be any name.
@todo rethink that process and provide a better solution.

You must then use the `init()` function to append new values or change existing values
of the configuration variables.

All defined variables in the file `config.go` can be changed. File is documented.

Why the init() function? https://golang.org/doc/effective_go.html#initialization

TODO(CS) Use https://github.com/AOEpeople/MageTestStand for TravisCI
*/
package codegen
