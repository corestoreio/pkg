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

package cfgfile

import (
	"io"
	"os"
	"path/filepath"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/log"
	"github.com/corestoreio/pkg/config"
)

type option func(*config.Service, func(config.Setter, io.Reader) error) error

func processFile(file string, s *config.Service, cb func(config.Setter, io.Reader) error) (err error) {
	var f io.ReadCloser
	f, err = os.Open(file)
	if err != nil {
		return errors.NotFound.New(err, "[cfgfile] os.Open")
	}
	defer func() {
		if errC := f.Close(); err == nil && errC != nil {
			err = errC
		}
	}()

	if err = cb(s, f); err != nil {
		err = errors.WithStack(err)
	}
	return
}

// If someone needs it, uncomment and add a test
//func WithIOReader(r io.Reader) option {
//	return func(s *config.Service, cb func(config.Setter, io.Reader) error) error {
//		return errors.WithStack(cb(s, r))
//	}
//}

// WithGlob uses a glob pattern to search for configurations files. If the
// pattern contains the variable from constant EnvNamePlaceHolder it gets
// replaced with the current environment name.
//		/var/www/site/config/{ENV}/*.yaml
func WithGlob(pattern string) option {
	return func(s *config.Service, cb func(config.Setter, io.Reader) error) error {
		p2 := s.ReplaceEnvName(pattern)
		matches, err := filepath.Glob(p2)
		if s.Log != nil && s.Log.IsDebug() {
			s.Log.Debug("cfgfile.WithGlob", log.String("pattern", pattern), log.String("replaced", p2),
				log.Strings("matched_files", matches...), log.String("env_name", s.EnvName()))
		}
		if err != nil {
			return errors.WithStack(err)
		}
		for _, file := range matches {
			if err := processFile(file, s, cb); err != nil {
				return errors.Wrapf(err, "[cfgfile] WithGlob for file %q", file)
			}
		}
		return nil
	}
}

// WithFiles loads the given files into the config.Service. If a path contains
// the variable from constant EnvNamePlaceHolder it gets replaced with the
// current environment name.
//		WithFiles([]string{"/var","www","site","config","db_{ENV}.yaml"})
func WithFiles(filePathParts ...[]string) option {
	return func(s *config.Service, cb func(config.Setter, io.Reader) error) error {
		for _, fileNameParts := range filePathParts {
			fp := filepath.Join(fileNameParts...)
			if err := processFile(s.ReplaceEnvName(fp), s, cb); err != nil {
				return errors.WithStack(err)
			}
		}
		return nil
	}
}

// WithFile loads the given files into the config.Service. If a path contains
// the variable from constant EnvNamePlaceHolder it gets replaced with the
// current environment name.
//		WithFile("/var","www","site","config/db_{ENV}.yaml")
func WithFile(filePath ...string) option {
	return func(s *config.Service, cb func(config.Setter, io.Reader) error) error {
		fp := filepath.Join(filePath...)
		err := processFile(s.ReplaceEnvName(fp), s, cb)
		return errors.WithStack(err)
	}
}
