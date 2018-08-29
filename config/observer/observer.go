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

package observer

import (
	"sync"

	"github.com/corestoreio/pkg/config"
)

// FactoryFunc allows to implement a custom observer which gets created based on
// the input data. The function gets called in Configuration.MakeObserver or in
// JSONRegisterObservers. Input data can be raw JSON or YAML or XML.
type FactoryFunc func(data []byte) (config.Observer, error)

type obsReg struct {
	sync.RWMutex
	pool map[string]FactoryFunc
}

var observerRegistry = &obsReg{
	pool: make(map[string]FactoryFunc),
}

// RegisterFactory adds a custom observer factory to the global registry. A
// custom observer can be accessed via Configuration.MakeObserver or via
// JSONRegisterObservers.
func RegisterFactory(typeName string, fn FactoryFunc) {
	observerRegistry.Lock()
	defer observerRegistry.Unlock()
	observerRegistry.pool[typeName] = fn
}

func lookupFactory(typeName string) (FactoryFunc, bool) {
	observerRegistry.RLock()
	defer observerRegistry.RUnlock()
	fn, ok := observerRegistry.pool[typeName]
	return fn, ok
}

func availableFactories(ret ...string) []string {
	observerRegistry.RLock()
	defer observerRegistry.RUnlock()
	for n := range observerRegistry.pool {
		ret = append(ret, n)
	}
	return ret
}
