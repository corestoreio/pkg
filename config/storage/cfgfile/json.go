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
	"encoding/json"
	"io/ioutil"
	"path/filepath"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/pkg/config"
	"github.com/corestoreio/pkg/util/conv"
)

// WithLoadJSON reads the configuration values from a JSON file and applies it
// to the config.service. "testdata/example.json" provides an example JSON file.
func WithLoadJSON(pathToFile ...string) config.LoadDataFn {
	return func(s *config.Service) error {

		fp := filepath.Join(pathToFile...)
		fpd, err := ioutil.ReadFile(fp)
		if err != nil {
			return errors.WithStack(err)
		}

		jd := make(map[string]interface{})
		if err := json.Unmarshal(fpd, &jd); err != nil {
			return errors.Wrapf(err, "[cfgjson] JSON loading error in file: %q", fp)
		}

		for route, v1 := range jd {
			k2, ok := v1.(map[string]interface{})
			if !ok {
				return errors.CorruptData.Newf("[cfgjson] Unexpected data in %#v", v1)
			}

			for scp, v2 := range k2 {

				var p = new(config.Path)
				switch v2t := v2.(type) {
				case map[string]interface{}:
					for scpID, dataIF := range v2t {

						data, err := conv.ToByteE(dataIF)
						if err != nil {
							return errors.CorruptData.New(err, "[cfgjson] Failed to convert %v into a byte slice for path: %q %q %q", dataIF, route, scp, scpID)
						}

						if err := p.ParseStrings(scp, scpID, route); err != nil {
							return errors.CorruptData.New(err, "[cfgjson] Failed to create path: %q %q %q", route, scp, scpID)
						}

						if err := s.Set(p, data); err != nil {
							return errors.Fatal.New(err, "[cfgjson] Service.Set failed with %q", p.String())
						}
					}

				case string, int, float64, bool:
					data, err := conv.ToByteE(v2t)
					if err != nil {
						return errors.CorruptData.New(err, "[cfgjson] Failed to convert %v into a byte slice for path: %q %q", v2t, route, scp)
					}

					if err := p.ParseStrings(scp, "0", route); err != nil {
						return errors.CorruptData.New(err, "[cfgjson] Failed to create path: %q %q", scp, route)
					}

					if err := s.Set(p, data); err != nil {
						return errors.Fatal.New(err, "[cfgjson] Service.Set failed with %q", p.String())
					}

				default:
					return errors.CorruptData.Newf("[cfgjson] Unexpected data in %#v", v2)
				}
			}
		}

		return nil
	}
}
