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

// +build csall yaml cue

package store

import (
	"os"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/pkg/storage/null"
	"github.com/corestoreio/pkg/util/conv"
	"gopkg.in/yaml.v2"
)

// YAMLOptions sets options to WithLoadFromYAML.
type YAMLOptions struct {
	// NextAutoIncWebsite defines the MySQL/MariaDB table next auto increment
	// value.
	NextAutoIncWebsite uint32
	// NextAutoIncGroup defines the MySQL/MariaDB table next auto increment
	// value.
	NextAutoIncGroup uint32
	// NextAutoIncStore defines the MySQL/MariaDB table next auto increment
	// value.
	NextAutoIncStore uint32
	// MapRootCategoryIDByName defines a custom function to map the retrieved
	// name from the YAML file to an ID (usually the category_id from the DB).
	MapRootCategoryIDByName func(rootCategoryName string) (rootCategoryID uint32, _ error)
}

// WithLoadFromYAML loads a file into memory. Supports two different version of
// the yaml file. See testdata directory "store-structureX.yaml".
func WithLoadFromYAML(pathToFile string, yo YAMLOptions) Option {
	if yo.MapRootCategoryIDByName == nil {
		yo.MapRootCategoryIDByName = func(_ string) (uint32, error) {
			return 0, nil
		}
	}
	return Option{
		sortOrder: 210,
		fn: func(s *Service) error {
			var unlocked bool
			s.mu.Lock()
			defer func() {
				if !unlocked {
					s.mu.Unlock()
				}
			}()

			f, err := os.Open(pathToFile)
			if err != nil {
				return errors.ReadFailed.New(err, "[store] Can't read file")
			}
			defer f.Close()

			dec := yaml.NewDecoder(f)
			dec.SetStrict(true)
			data := map[string]map[string]map[string]map[string]interface{}{}
			if err := dec.Decode(data); err != nil {
				return errors.CorruptData.New(err, "[store] Can't decode file %q", pathToFile)
			}

			var websites StoreWebsites
			var groups StoreGroups
			var stores Stores
			for _, wData := range data["store-structure"]["websites"] {
				defaultGroupID := uint32(conv.ToUint(wData["default_group_id"]))
				if defaultGroupID == 0 {
					defaultGroupID = yo.NextAutoIncGroup
				}
				websites.Append(&StoreWebsite{
					WebsiteID:      yo.NextAutoIncWebsite,
					Code:           conv.ToString(wData["code"]),
					Name:           null.MakeString(conv.ToString(wData["name"])),
					SortOrder:      uint32(conv.ToUint(wData["sort_order"])),
					DefaultGroupID: defaultGroupID,
					IsDefault:      conv.ToBool(wData["is_default"]),
				})

				yGroups := wData["stores"]
				if yGroups == nil {
					yGroups = wData["groups"]
				}
				wdm, err := conv.ToStringMapE(yGroups)
				if err != nil {
					return errors.WithStack(err)
				}

				for _, gDataIF := range wdm {
					gData := gDataIF.(map[interface{}]interface{})
					defaultStoreID := uint32(conv.ToUint(wData["default_store_id"]))
					if defaultStoreID == 0 {
						defaultStoreID = yo.NextAutoIncStore
					}

					rcID, err := yo.MapRootCategoryIDByName(conv.ToString(gData["root-category"]))
					if err != nil {
						return errors.WithStack(err)
					}

					groups.Append(&StoreGroup{
						GroupID:        yo.NextAutoIncGroup,
						WebsiteID:      yo.NextAutoIncWebsite,
						Code:           conv.ToString(gData["code"]),
						Name:           conv.ToString(gData["name"]),
						DefaultStoreID: defaultStoreID,
						RootCategoryID: rcID,
					})

					yStores := gData["store-views"]
					if yStores == nil {
						yStores = gData["stores"]
					}
					svm, err := conv.ToStringMapE(yStores)
					if err != nil {
						return errors.WithStack(err)
					}

					for _, sDataIF := range svm {
						sData, err := conv.ToStringMapStringE(sDataIF)
						if err != nil {
							return errors.WithStack(err)
						}
						stores.Append(&Store{
							StoreID:   yo.NextAutoIncStore,
							GroupID:   yo.NextAutoIncGroup,
							WebsiteID: yo.NextAutoIncWebsite,
							Code:      sData["code"],
							Name:      sData["name"],
							SortOrder: uint32(conv.ToUint(sData["sort_order"])),
							IsActive:  conv.ToBool(sData["is_active"]),
						})
						yo.NextAutoIncStore++
					}
					yo.NextAutoIncGroup++
				}
				yo.NextAutoIncWebsite++
			}
			unlocked = true
			s.mu.Unlock()
			return s.Options(
				WithWebsites(websites.Data...),
				WithGroups(groups.Data...),
				WithStores(stores.Data...),
			)
		},
	}
}
