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

package bgwork

import (
	"context"
	"sync"

	"golang.org/x/sync/errgroup"
)

// Wait is a better pattern around using sync.WaitGroup. For returning errors
// while working on a subtask of a common task, consider using package
// x/errgroup.
func Wait(length int, block func(index int)) {
	if length == 1 {
		block(0)
		return
	}
	var w sync.WaitGroup
	w.Add(length)
	for i := 0; i < length; i++ {
		go func(w *sync.WaitGroup, index int) {
			block(index)
			w.Done()
		}(&w, i)
	}
	w.Wait()
}

// WaitContext is a better pattern around using errgroup.Group. See function
// Wait. The returned error of the block function cancels the context.
func WaitContext(length int, block func(ctx context.Context, index int) error) error {
	ctx := context.Background()
	if length == 1 {
		return block(ctx, 0)
	}

	g, ctx := errgroup.WithContext(ctx)

	for i := 0; i < length; i++ {
		i := i
		g.Go(func() error {
			return block(ctx, i)
		})
	}
	return g.Wait()
}
