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

package dbr

import (
	"database/sql"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/log"
)

// Tx at a transaction for the given Session
type Tx struct {
	log.Logger
	*sql.Tx
}

// Begin creates a transaction for the given session
func (c *Connection) Begin() (*Tx, error) {
	dbTx, err := c.DB.Begin()
	if err != nil {
		return nil, errors.Wrap(err, "[dbr] transaction.begin.error")
	}
	tx := &Tx{
		Tx: dbTx,
	}
	if c.Log != nil {
		tx.Logger = c.Log.With(log.Bool("transaction", true))
	}
	return tx, nil
}

// Commit finishes the transaction
func (tx *Tx) Commit() error {
	return errors.Wrap(tx.Tx.Commit(), "[dbr] transaction.commit.error")
}

// Rollback cancels the transaction
func (tx *Tx) Rollback() error {
	return errors.Wrap(tx.Tx.Rollback(), "[dbr] transaction.rollback.error")
}

// Wrap is a helper method that will automatically COMMIT or ROLLBACK once the
// supplied functions are done executing.
//
//      tx, err := db.Begin()
//      if err != nil{
//           panic(err.Error()) // you could gracefully handle the error also
//      }
//      if err := tx.Wrap(func() error {
//          // SQL
//          return nil
//      }); err != nil{
//           panic(err.Error()) // you could gracefully handle the error also
//      }
func (tx *Tx) Wrap(fns ...func() error) error {
	for i, f := range fns {
		if err := f(); err != nil {
			if rErr := tx.Rollback(); rErr != nil {
				return errors.Wrapf(err, "[dbr] transaction.wrap.Rollback.error at index %d", i)
			}
			return errors.Wrapf(err, "[dbr] transaction.wrap.error at index %d", i)
		}
	}
	return errors.Wrap(tx.Commit(), "[dbr] transaction.wrap.commit")
}
