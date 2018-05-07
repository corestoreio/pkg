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

package storage

import (
	"context"
	"time"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/pkg/config"
	"github.com/corestoreio/pkg/store/scope"
	"github.com/golang/sync/errgroup"
)

type MultiOptions struct {
	// ContextTimeout if greater than zero a timeout will
	ContextTimeout time.Duration
}

// Multi wraps multiple backends into one. Writing to the backend
// implementations occur concurrent and in parallel. Even a timeout can be set
// to cancel the writing. Reading a value processes the backends in serial
// order. The backend which returns the first found value wins. Subsequent calls
// to other backends are getting skipped.
type multi struct {
	options  MultiOptions
	backends []config.Storager
}

// MakeMulti creates a new Multi backend wrapper. Supports other Multi backend
// wrappers.
func MakeMulti(o MultiOptions, ss ...config.Storager) config.Storager {
	allStorages := make([]config.Storager, 0, len(ss))
	for _, s := range ss {
		if mw, ok := s.(*multi); ok {
			allStorages = append(allStorages, mw.backends...)
		} else {
			allStorages = append(allStorages, s)
		}
	}
	return &multi{options: o, backends: allStorages}
}

// Set writes concurrently to the backends. A ContextTimeout can be defined to
// cancel the internal goroutine. It returns the first error.
func (ms *multi) Set(scp scope.TypeID, path string, value []byte) error {
	// investigate if that concept of timeout and cancellation is good enough
	ctx := context.Background()
	if ms.options.ContextTimeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, ms.options.ContextTimeout)
		defer cancel()
	}

	g, ctx := errgroup.WithContext(ctx)

	for _, s := range ms.backends {
		s := s
		g.Go(func() error {
			errChan := make(chan error)
			stopChan := make(chan struct{})

			go func() {
				select {
				case <-stopChan:
					return
				case errChan <- errors.WithStack(s.Set(scp, path, value)):
				}
			}()

			select {
			case <-ctx.Done():
				close(stopChan)
				return ctx.Err()
			case err := <-errChan:
				close(stopChan)
				close(errChan)
				return err
			}
		})
	}

	return g.Wait()
}

// Get returns the first found value from the backend storage.
func (ms *multi) Get(scp scope.TypeID, path string) (v []byte, found bool, err error) {
	for idx, s := range ms.backends {
		v, found, err = s.Get(scp, path)
		if err != nil {
			return nil, false, errors.Wrapf(err, "[config] Multi.Value failed at backend index %d with path %q", idx, path)
		}
		if found {
			return
		}
	}
	return nil, false, nil
}
