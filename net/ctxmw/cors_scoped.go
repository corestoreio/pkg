// Copyright 2015, Cyrill @ Schumacher.fm and the CoreStore contributors
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

package ctxmw

import (
	"net/http"

	"sync"

	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/net/ctxhttp"
	"github.com/corestoreio/csfw/store"
	"github.com/juju/errgo"
	"golang.org/x/net/context"
)

const (
	PathCorsExposedHeaders = "web/cors/exposed_headers"
	PathCorsAllowedOrigins = "web/cors/allowed_origins"
)

// corsService
type corsService struct {
	config config.Getter

	// rwmu protects the map
	rwmu sync.RWMutex
	// storage key is the website ID and value the current cors config
	storage map[int64]*Cors
}

// get uses a read lock to check if a Cors exists for a website ID. returns nil
// if there is no Cors pointer. aim is: multiple goroutines can read from the
// map while adding new Cors pointers can only be done by one goroutine.
func (cs *corsService) get(websiteID int64) *Cors {
	cs.rwmu.RLock()
	defer cs.rwmu.RUnlock()
	if c, ok := cs.storage[websiteID]; ok {
		return c
	}
	return nil
}

// create creates a new Cors type and returns it.
func (cs *corsService) insert(websiteID int64) *Cors {
	cs.rwmu.Lock()
	defer cs.rwmu.Unlock()

	// pulls the options from the scoped reader

	c := NewCors()
	cs.storage[websiteID] = c
	return c
}

// WithCORSScoped allows Website based cors configuration.
func WithCORSScoped(cg config.Getter) ctxhttp.Middleware {

	cs := &corsService{
		config:  cg,
		storage: make(map[int64]*Cors),
	}

	return func(hf ctxhttp.HandlerFunc) ctxhttp.HandlerFunc {
		return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			_, st, err := store.FromContextReader(ctx)
			if err != nil {
				return errgo.Mask(err)
			}

			var cc *Cors // cc == current CORS config the current request
			if cc = cs.get(st.WebsiteID()); cc == nil {
				cc = cs.insert(st.WebsiteID())
			}
			// todo: run a defer or goroutine to check if config changes
			// and if so delete the entry from the map

			if r.Method == "OPTIONS" {
				if cc.Log.IsDebug() {
					cc.Log.Debug("ctxmw.Cors.WithCORS.handlePreflight", "method", r.Method, "OptionsPassthrough", cc.OptionsPassthrough)
				}
				cc.handlePreflight(w, r)
				// Preflight requests are standalone and should stop the chain as some other
				// middleware may not handle OPTIONS requests correctly. One typical example
				// is authentication middleware ; OPTIONS requests won't carry authentication
				// headers (see #1)
				if cc.OptionsPassthrough {
					return hf(ctx, w, r)
				}
				return nil
			}
			if cc.Log.IsDebug() {
				cc.Log.Debug("ctxmw.Cors.WithCORS.handleActualRequest", "method", r.Method)
			}
			cc.handleActualRequest(w, r)
			return hf(ctx, w, r)
		}
	}
}
