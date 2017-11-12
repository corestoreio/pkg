// Copyright 2015-2017, Cyrill @ Schumacher.fm and the CoreStore contributors
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

package migration

import (
	"context"

	"github.com/corestoreio/pkg/sql/dml"
)

// ToUTF8MB4 converts MySQL compatible databases from utf8 to utf8mb4. What’s
// the difference between utf8 and utf8mb4? MySQL decided that UTF-8 can only
// hold 3 bytes per character. Why? No good reason can be found documented
// anywhere. Few years later, when MySQL 5.5.3 was released, they introduced a
// new encoding called utf8mb4, which is actually the real 4-byte utf8 encoding
// that you know and love.
//
// # Run this once on each schema you have (Replace database_name with your schema name)
// ALTER DATABASE database_name CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci;
//
// # Run this once for each table you have (replace table_name with the table name)
// ALTER TABLE table_name CONVERT TO CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
//
// # Run this for each column (replace table name, column_name, the column type, maximum length, etc.)
// ALTER TABLE table_name CHANGE column_name column_name VARCHAR(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
func ToUTF8MB4(ctx context.Context, db interface {
	dml.Querier
	dml.Execer
	dml.Preparer
}) error {
	// TODO: Implement utf8mb4 conversion
	return nil
}
