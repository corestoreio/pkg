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
	"gopkg.in/yaml.v2"
)

// WithLoadYAML reads the configuration values from a JSON file and applies it
// to the config.service. "testdata/example.json" provides an example JSON file.
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
