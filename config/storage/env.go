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
	"unicode"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/log"
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
	if err := p.Parse(envVar); err != nil {
		return nil, errors.WithStack(err)
	}
	return &p, nil
}

// EnvOp allows to set options when creating a new Environment storage service.
type EnvOp struct {
	// Prefix can be set to as a custom optional Prefix in case where default
	// Prefix does not work.
	Prefix string
}

func isEnvVarAllowed(prefix, varName string) bool {
	if !strings.HasPrefix(varName, prefix) {
		return false
	}
	for _, r := range varName {
		if ok := unicode.IsUpper(r) || unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_'; !ok {
			return false
		}
	}
	return true
}

// WithLoadEnvironmentVariables reads configuration values from environment
// variables. A valid env var is always uppercase and [A-Z0-9_].
//		CONFIG__ETCD__CREDENTIALS__USER_NAME -> etc/credentials/user_name
//
// Env vars (environment variables) can be unset after reading and optionally
// cached. They can be read during initialization of the storage backend or on
// each request.
func WithLoadEnvironmentVariables(op EnvOp) config.LoadDataOption {
	if op.Prefix == "" {
		op.Prefix = Prefix
	}

	return config.MakeLoadDataOption(func(s *config.Service) (err error) {
		for _, ev := range os.Environ() {
			equalPos := strings.IndexRune(ev, '=')
			if equalPos < 0 {
				continue
			}
			key := ev[:equalPos]
			if !isEnvVarAllowed(op.Prefix, key) {
				continue
			}
			envVal, ok := os.LookupEnv(key)
			if ok {
				p, err := FromEnvVar(op.Prefix, key)
				if err != nil {
					return errors.WithStack(err)
				}
				if err := s.Set(p, []byte(envVal)); err != nil {
					return errors.WithStack(err)
				}
			}
			if s.Log != nil && s.Log.IsDebug() {
				s.Log.Debug("config.storage.WithLoadEnvironmentVariables", log.String("name", key), log.Bool("found", ok), log.Int("value_length", len(envVal)))
			}
		}

		return
	}).WithUseStorageLevel(1)
}
