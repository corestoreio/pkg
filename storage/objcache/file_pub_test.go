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

// +build file csall

package objcache_test

import (
	"testing"
	"time"

	"github.com/corestoreio/pkg/storage/objcache"
)

func TestNewFileSystemClient_Delete(t *testing.T) {
	newTestServiceDelete(t, objcache.NewFileSystemClient(nil))
}

func TestNewFileSystemClient_Expires(t *testing.T) {
	testExpiration(t, func() {
		time.Sleep(time.Second * 2)
	}, objcache.NewFileSystemClient(nil), newSrvOpt(JSONCodec{}))
}

func TestNewFileSystemClient_ComplexParallel(t *testing.T) {

	t.Run("gob", func(t *testing.T) {
		newServiceComplexParallelTest(t, objcache.NewFileSystemClient(&objcache.FileSystemConfig{
			DirectoryLevel: -1,
			CleanOnClose:   true,
		}), nil)
	})

	t.Run("json", func(t *testing.T) {
		newServiceComplexParallelTest(t, objcache.NewFileSystemClient(&objcache.FileSystemConfig{
			CleanOnClose: true,
		}), &objcache.ServiceOptions{
			Codec: JSONCodec{},
		})
	})

}
