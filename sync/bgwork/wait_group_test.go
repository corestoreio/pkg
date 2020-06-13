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

package bgwork_test

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"

	"github.com/corestoreio/pkg/sync/bgwork"
	"github.com/corestoreio/pkg/util/assert"
)

func TestWait(t *testing.T) {
	haveIdx := new(int32)
	const goroutines = 3
	bgwork.Wait(goroutines, func(index int) {
		atomic.AddInt32(haveIdx, 1)
	})
	if *haveIdx != goroutines {
		t.Errorf("Have %d Want %d", haveIdx, goroutines)
	}
}

func TestWaitContext(t *testing.T) {
	type args struct {
		length int
		block  func(ctx context.Context, called []bool, index int) error
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "err cancelation", args: args{length: 2, block: func(ctx context.Context, called []bool, i int) error {
			called[i] = true
			if i == 0 {
				return errors.New("err")
			}
			<-ctx.Done()
			return nil
		}}, wantErr: true},

		{name: "all", args: args{length: 200, block: func(ctx context.Context, called []bool, i int) error {
			called[i] = true
			return nil
		}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			called := make([]bool, tt.args.length)
			if err := bgwork.WaitContext(tt.args.length, func(ctx context.Context, i int) error {
				return tt.args.block(ctx, called, i)
			}); (err != nil) != tt.wantErr {
				t.Errorf("WaitContext() error = %v, wantErr %v", err, tt.wantErr)
			}

			for i, call := range called {
				assert.True(t, call, "Call to goroutine %d did not occur", i)
			}
		})
	}
}
