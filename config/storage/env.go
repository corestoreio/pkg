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
	"os"
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

const strSep = string(config.PathSeparator)

var (
	replTo   = strings.NewReplacer(string(config.PathSeparator), SlashSeparator, scope.StrDefault.String()+strSep+"0"+strSep, "")
	replFrom = strings.NewReplacer(SlashSeparator, string(config.PathSeparator))
)

// ToEnvVar converts a Path to a valid environment key. Returns an empty string
// in case of an error.
//	scope.DefaultTypeID, etc/credentials/user_name => CONFIG__ETCD__CREDENTIALS__USER_NAME
func ToEnvVar(p *config.Path) string {
	buf := bufferpool.Get()
	defer bufferpool.Put(buf)
	buf.WriteString(Prefix)

	if err := p.AppendFQ(buf); err != nil {
		return ""
	}
	return strings.ToUpper(replTo.Replace(buf.String()))
}

// FromEnvVar parses an environment key into a config.Path.
//	CONFIG__ETCD__CREDENTIALS__USER_NAME => scope.DefaultTypeID, etc/credentials/user_name
func FromEnvVar(prefix, envVar string) (*config.Path, error) {
	lp := len(prefix)
	if len(envVar) > lp {
		envVar = envVar[lp:]
	}
	envVar = strings.ToLower(replFrom.Replace(envVar))

	slashes := strings.Count(envVar, string(config.PathSeparator))
	if slashes == 2 {
		// Looks like: etcd/credentials/user_name for default scope
		return config.NewPathWithScope(scope.DefaultTypeID, envVar)
	}
	// looks like: stores/2/carrier/dhl/title

	var p config.Path
	if err := p.ParseFQ(envVar); err != nil {
		return nil, errors.WithStack(err)
	}
	return &p, nil
}

// EnvOp allows to set options when creating a new Environment storage service.
type EnvOp struct {
	// Preload loads all applicable environment variables during instantiation of
	// NewEnvironment and caches the variables internally. No other access to the
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
	CacheVariableFn func(*config.Path) bool
}

// Storage reads configuration values from environment variables. A valid env var is
// always uppercase and [A-Z0-9_].
// CONFIG__ETCD__CREDENTIALS__USER_NAME -> etc/credentials/user_name
// Package cfgenv reads config paths from special crafted environment variables.
//
// Env vars (environment variables) can be unset after reading and optionally
// cached. They can be read during initialization of the storage backend or on
// each request.
type Environment struct {
	options EnvOp
	mu      sync.Mutex
	cache   map[cacheKey]string
}

// NewEnvironment creates a new storage service which reads from the environment
// variables. EnvOp provides various options to configure the behaviour.
func NewEnvironment(o EnvOp) (config.Storager, error) {
	s := &Environment{
		options: o,
		cache:   make(map[cacheKey]string),
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

			p, err := FromEnvVar(s.options.Prefix, key)
			if err != nil {
				return nil, errors.WithStack(err)
			}
			if s.options.CacheVariableFn == nil || s.options.CacheVariableFn(p) {
				s.cache[makeCacheKey(p.ScopeRoute())] = ev[equalPos+1:]
				if s.options.UnsetEnvAfterRead {
					os.Unsetenv(key)
				}
			}
		}
	}
	if s.options.CacheVariableFn == nil {
		s.options.CacheVariableFn = func(_ *config.Path) bool { return false }
	}
	return s, nil
}

// Set is a no-op function.
func (s *Environment) Set(_ *config.Path, _ []byte) error {
	return nil
}

// Get returns the value from the environment by match scope and route.
func (s *Environment) Get(p *config.Path) (v []byte, found bool, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	key := makeCacheKey(p.ScopeRoute())
	if s.cache != nil {
		if val, ok := s.cache[key]; ok {
			return []byte(val), true, nil
		}
	}

	envKey := ToEnvVar(p)
	if !s.isAllowed(envKey) {
		return nil, false, nil
	}

	envVal, ok := os.LookupEnv(envKey)
	if !ok {
		return nil, false, nil
	}

	if s.options.CacheVariableFn(p) {
		s.cache[key] = envVal
	}
	if s.options.UnsetEnvAfterRead {
		if err := os.Unsetenv(envKey); err != nil {
			return nil, false, errors.Wrapf(err, "[cfgenv] Unset failed with env variable key: %q", envKey)
		}
	}
	return []byte(envVal), true, nil
}

func (s *Environment) isAllowed(varName string) bool {
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
