// Copyright 2015-2016, Cyrill @ Schumacher.fm and the CoreStore contributors
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

package store_test

import (
	"github.com/corestoreio/csfw/storage/csdb"
	"github.com/corestoreio/csfw/storage/dbr"
	"github.com/corestoreio/csfw/store"
)

func init() {
	store.TableCollection = csdb.MustNewTableService(
		csdb.WithTable(
			store.TableIndexStore,
			"store",
			csdb.Column{Field: dbr.NewNullString(`store_id`), Type: dbr.NewNullString(`smallint(5) unsigned`), Null: dbr.NewNullString(`NO`), Key: dbr.NewNullString(`PRI`), Default: dbr.NewNullString(nil), Extra: dbr.NewNullString(`auto_increment`)},
			csdb.Column{Field: dbr.NewNullString(`code`), Type: dbr.NewNullString(`varchar(32)`), Null: dbr.NewNullString(`YES`), Key: dbr.NewNullString(`UNI`), Default: dbr.NewNullString(nil), Extra: dbr.NewNullString(``)},
			csdb.Column{Field: dbr.NewNullString(`website_id`), Type: dbr.NewNullString(`smallint(5) unsigned`), Null: dbr.NewNullString(`NO`), Key: dbr.NewNullString(`MUL`), Default: dbr.NewNullString(`0`), Extra: dbr.NewNullString(``)},
			csdb.Column{Field: dbr.NewNullString(`group_id`), Type: dbr.NewNullString(`smallint(5) unsigned`), Null: dbr.NewNullString(`NO`), Key: dbr.NewNullString(`MUL`), Default: dbr.NewNullString(`0`), Extra: dbr.NewNullString(``)},
			csdb.Column{Field: dbr.NewNullString(`name`), Type: dbr.NewNullString(`varchar(255)`), Null: dbr.NewNullString(`NO`), Key: dbr.NewNullString(``), Default: dbr.NewNullString(nil), Extra: dbr.NewNullString(``)},
			csdb.Column{Field: dbr.NewNullString(`sort_order`), Type: dbr.NewNullString(`smallint(5) unsigned`), Null: dbr.NewNullString(`NO`), Key: dbr.NewNullString(``), Default: dbr.NewNullString(`0`), Extra: dbr.NewNullString(``)},
			csdb.Column{Field: dbr.NewNullString(`is_active`), Type: dbr.NewNullString(`smallint(5) unsigned`), Null: dbr.NewNullString(`NO`), Key: dbr.NewNullString(`MUL`), Default: dbr.NewNullString(`0`), Extra: dbr.NewNullString(``)},
		),
		csdb.WithTable(
			store.TableIndexGroup,
			"store_group",
			csdb.Column{Field: dbr.NewNullString(`group_id`), Type: dbr.NewNullString(`smallint(5) unsigned`), Null: dbr.NewNullString(`NO`), Key: dbr.NewNullString(`PRI`), Default: dbr.NewNullString(nil), Extra: dbr.NewNullString(`auto_increment`)},
			csdb.Column{Field: dbr.NewNullString(`website_id`), Type: dbr.NewNullString(`smallint(5) unsigned`), Null: dbr.NewNullString(`NO`), Key: dbr.NewNullString(`MUL`), Default: dbr.NewNullString(`0`), Extra: dbr.NewNullString(``)},
			csdb.Column{Field: dbr.NewNullString(`name`), Type: dbr.NewNullString(`varchar(255)`), Null: dbr.NewNullString(`NO`), Key: dbr.NewNullString(``), Default: dbr.NewNullString(nil), Extra: dbr.NewNullString(``)},
			csdb.Column{Field: dbr.NewNullString(`root_category_id`), Type: dbr.NewNullString(`int(10) unsigned`), Null: dbr.NewNullString(`NO`), Key: dbr.NewNullString(``), Default: dbr.NewNullString(`0`), Extra: dbr.NewNullString(``)},
			csdb.Column{Field: dbr.NewNullString(`default_store_id`), Type: dbr.NewNullString(`smallint(5) unsigned`), Null: dbr.NewNullString(`NO`), Key: dbr.NewNullString(`MUL`), Default: dbr.NewNullString(`0`), Extra: dbr.NewNullString(``)},
		),
		csdb.WithTable(
			store.TableIndexWebsite,
			"store_website",
			csdb.Column{Field: dbr.NewNullString(`website_id`), Type: dbr.NewNullString(`smallint(5) unsigned`), Null: dbr.NewNullString(`NO`), Key: dbr.NewNullString(`PRI`), Default: dbr.NewNullString(nil), Extra: dbr.NewNullString(`auto_increment`)},
			csdb.Column{Field: dbr.NewNullString(`code`), Type: dbr.NewNullString(`varchar(32)`), Null: dbr.NewNullString(`YES`), Key: dbr.NewNullString(`UNI`), Default: dbr.NewNullString(nil), Extra: dbr.NewNullString(``)},
			csdb.Column{Field: dbr.NewNullString(`name`), Type: dbr.NewNullString(`varchar(64)`), Null: dbr.NewNullString(`YES`), Key: dbr.NewNullString(``), Default: dbr.NewNullString(nil), Extra: dbr.NewNullString(``)},
			csdb.Column{Field: dbr.NewNullString(`sort_order`), Type: dbr.NewNullString(`smallint(5) unsigned`), Null: dbr.NewNullString(`NO`), Key: dbr.NewNullString(`MUL`), Default: dbr.NewNullString(`0`), Extra: dbr.NewNullString(``)},
			csdb.Column{Field: dbr.NewNullString(`default_group_id`), Type: dbr.NewNullString(`smallint(5) unsigned`), Null: dbr.NewNullString(`NO`), Key: dbr.NewNullString(`MUL`), Default: dbr.NewNullString(`0`), Extra: dbr.NewNullString(``)},
			csdb.Column{Field: dbr.NewNullString(`is_default`), Type: dbr.NewNullString(`smallint(5) unsigned`), Null: dbr.NewNullString(`YES`), Key: dbr.NewNullString(``), Default: dbr.NewNullString(`0`), Extra: dbr.NewNullString(``)},
		),
	)
}
