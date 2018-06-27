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

	"github.com/corestoreio/errors"
	"github.com/corestoreio/pkg/config"
	"github.com/corestoreio/pkg/store/scope"
	"gopkg.in/yaml.v2"
)

// WithLoadYAML reads the configuration values from a YAML file and applies it
// to the config.service. "testdata/example.yaml" provides an example YAML file.
// This function processes the YAML file stream based.
func WithLoadYAML(opts ...option) config.LoadDataOption {
	return config.MakeLoadDataOption(func(s *config.Service) (err error) {
		for i := 0; i < len(opts) && err == nil; i++ {
			err = opts[i](s, loadYAML)
		}
		return
	}).WithUseStorageLevel(1)
}

func loadYAML(s config.Setter, r io.Reader) error {

	d := yaml.NewDecoder(r)
	d.SetStrict(true)
	p := new(config.Path)

	for {
		var yd map[string]map[string]map[string]string
		if err := d.Decode(&yd); err == io.EOF {
			break
		} else if err != nil {
			return errors.WithStack(err)
		}
		for route, v1 := range yd {
			for scp, v2 := range v1 {
				for scpID, data := range v2 {
					if err := p.ParseStrings(scp, scpID, route); err != nil {
						return errors.WithStack(err)
					}
					if err := s.Set(p, []byte(data)); err != nil {
						return errors.WithStack(err)
					}
				}
			}
		}
	}
	return nil
}

// WithLoadFieldMetaYAML reads the immutable default values and permissions from
// a YAML file and applies it to the config.Service. The data gets loaded only
// once. "testdata/example_field_meta.yaml" provides an example YAML file.
// WithLoadFieldMetaYAML has been optimized for stream operation and does not load
// the whole file into memory.
func WithLoadFieldMetaYAML(opts ...option) config.LoadDataOption {
	return config.WithFieldMetaGenerator(func(s *config.Service) (<-chan *config.FieldMeta, <-chan error) {

		fmC := make(chan *config.FieldMeta)
		errC := make(chan error)

		go func() {
			defer func() { close(fmC); close(errC) }()

			loadYaml := func(_ config.Setter, r io.Reader) error {

				d := yaml.NewDecoder(r)
				d.SetStrict(true)

				var yd map[string]struct {
					Default  string
					Perm     string
					Websites map[int64]string
					Stores   map[int64]string
				}

				for {
					if err := d.Decode(&yd); err == io.EOF {
						break
					} else if err != nil {
						return errors.Fatal.New(err, "[cfgfile] YAML.Decode")
					}

					for route, meta := range yd {

						var wsp scope.Perm
						if meta.Perm != "" {
							var err error
							wsp, err = scope.MakePerm(meta.Perm)
							if err != nil {
								return errors.WithStack(err)
							}
						}

						fmC <- &config.FieldMeta{
							Route:          route,
							WriteScopePerm: wsp,
							ScopeID:        0,
							DefaultValid:   meta.Default != "",
							Default:        meta.Default,
						}

						sendFM(fmC, meta.Websites, route, scope.Website)
						sendFM(fmC, meta.Stores, route, scope.Store)
					}
				}
				return nil
			}

			for _, opt := range opts {
				if err := opt(s, loadYaml); err != nil {
					errC <- errors.WithStack(err)
					return
				}
			}
		}()

		return fmC, errC
	})

}

func sendFM(fmC chan *config.FieldMeta, data map[int64]string, route string, scp scope.Type) {
	for id, defaultVal := range data {
		// prevent race condition
		fmC <- &config.FieldMeta{
			Route:        route,
			ScopeID:      scp.WithID(id),
			DefaultValid: defaultVal != "",
			Default:      defaultVal,
		}
	}
}
