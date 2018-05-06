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

package cfgenv

import (
	"os"
	"strconv"
	"strings"
	"sync"
	"unicode"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/pkg/config"
	"github.com/corestoreio/pkg/store/scope"
	"github.com/corestoreio/pkg/util/bufferpool"
)

// SlashSeparator acts as a slash separator for environment keys.
const SlashSeparator = "__"

// Prefix with which an environment variable must begin with to be considered as
// a configuration value.
const Prefix = "CONFIG" + SlashSeparator

var (
	replTo   = strings.NewReplacer(string(config.Separator), SlashSeparator)
	replFrom = strings.NewReplacer(SlashSeparator, string(config.Separator))
)

// ToEnvVar converts a scope and its ID together with a route to a valid
// environment key.
//	scope.DefaultTypeID, etc/credentials/user_name => CONFIG__ETCD__CREDENTIALS__USER_NAME
func ToEnvVar(scpID scope.TypeID, route string) string {
	buf := bufferpool.Get()
	defer bufferpool.Put(buf)
	buf.WriteString(Prefix)

	if scp, id := scpID.Unpack(); scp > scope.Default {
		buf.WriteString(scp.StrType())
		buf.WriteString(SlashSeparator)
		buf.WriteString(strconv.FormatInt(id, 10))
		buf.WriteString(SlashSeparator)
	}
	buf.WriteString(route)
	return strings.ToUpper(replTo.Replace(buf.String()))
}

// FromEnvVar parses an environment key into the scope type ID and the route.
// Returns an empty route on error.
//	CONFIG__ETCD__CREDENTIALS__USER_NAME => scope.DefaultTypeID, etc/credentials/user_name
func FromEnvVar(prefix, envVar string) (scpID scope.TypeID, route string) {
	lp := len(prefix)
	if len(envVar) > lp {
		envVar = envVar[lp:]
	}
	envVar = strings.ToLower(replFrom.Replace(envVar))

	slashes := strings.Count(envVar, string(config.Separator))
	if slashes == 2 {
		// Looks like: etcd/credentials/user_name for default scope
		return scope.DefaultTypeID, envVar
	}
	// looks like: stores/2/carrier/dhl/title

	p, err := config.NewPathFromFQ(envVar)
	if err != nil {
		return
	}
	return p.ScopeRoute()
}

// Options used when creating a new Storage service.
type Options struct {
	// Preload loads all appicable environment variables during instantiation of
	// NewStorage and caches the variables internally. No other access to the
	// environment will be made.
	Preload bool
	// Prefix can be set to as a custom optional Prefix in case where default
	// Prefix does not work.
	Prefix string
	// UnsetEnvAfterRead calls os.Unsetenv after retrieving the value.
	UnsetEnvAfterRead bool
	// CacheVariableFn custom call back function which can be used to
	// distinguish if a variable should be cached internally or not. Returning
	// always true, caches all found environment variables. Only applicable when
	// Preload is false.
	CacheVariableFn func(scp scope.TypeID, route string) bool
}

type cacheKey struct {
	scope.TypeID
	string // route
}

// Storage reads configuration values from environment variables. A valid env var is
// always uppercase and [A-Z0-9_].
// CONFIG__ETCD__CREDENTIALS__USER_NAME -> etc/credentials/user_name
type Storage struct {
	options Options
	mu      sync.Mutex
	cache   map[cacheKey]string
}

// NewStorage creates a new storage service which reads from the environment
// variables.
func NewStorage(o Options) (*Storage, error) {
	s := &Storage{
		options: o,
		cache:   make(map[cacheKey]string),
	}
	if s.options.CacheVariableFn == nil {
		s.options.CacheVariableFn = func(scope.TypeID, string) bool { return false }
	}
	if s.options.Prefix == "" {
		s.options.Prefix = Prefix
	}

	if o.Preload {
		for _, ev := range os.Environ() {
			equalPos := strings.IndexRune(ev, '=')
			if equalPos < 0 {
				continue
			}
			key := ev[:equalPos]
			if !s.isAllowed(key) {
				continue
			}

			envScp, envRoute := FromEnvVar(s.options.Prefix, key)
			if envRoute != "" {
				s.cache[cacheKey{envScp, envRoute}] = ev[equalPos+1:]
				if s.options.UnsetEnvAfterRead {
					os.Unsetenv(key)
				}
			}
		}
	}

	return s, nil
}

// Set is a no-op function.
func (s *Storage) Set(scp scope.TypeID, route string, value []byte) error {
	return nil
}

// Value returns the value from the environment by match scope and route.
func (s *Storage) Value(scp scope.TypeID, route string) (v []byte, found bool, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.cache != nil {
		if val, ok := s.cache[cacheKey{scp, route}]; ok {
			return []byte(val), true, nil
		}
	}

	envKey := ToEnvVar(scp, route)
	if !s.isAllowed(envKey) {
		return nil, false, nil
	}

	envVal, ok := os.LookupEnv(envKey)
	if !ok {
		return nil, false, nil
	}

	if s.options.CacheVariableFn(scp, route) {
		s.cache[cacheKey{scp, route}] = envVal
	}
	if s.options.UnsetEnvAfterRead {
		if err := os.Unsetenv(envKey); err != nil {
			return nil, false, errors.Wrapf(err, "[cfgenv] Unset failed with env variable key: %q", envKey)
		}
	}
	return []byte(envVal), true, nil
}

func (s *Storage) isAllowed(varName string) bool {
	if !strings.HasPrefix(varName, s.options.Prefix) {
		return false
	}

	for _, r := range varName {
		if ok := unicode.IsUpper(r) || unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_'; !ok {
			return false
		}
	}
	return true
}
