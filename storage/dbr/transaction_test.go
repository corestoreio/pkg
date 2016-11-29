package dbr

import (
	// "database/sql"
	"testing"

	"context"

	"github.com/stretchr/testify/assert"
)

func TestTransactionReal(t *testing.T) {
	s := createRealSessionWithFixtures()

	tx, err := s.Begin()
	assert.NoError(t, err)

	res, err := tx.InsertInto("dbr_people").Columns("name", "email").Values("Barack", "obama@whitehouse.gov").Exec()

	assert.NoError(t, err)
	id, err := res.LastInsertId()
	assert.NoError(t, err)
	rowsAff, err := res.RowsAffected()
	assert.NoError(t, err)

	assert.True(t, id > 0)
	assert.Equal(t, rowsAff, int64(1))

	var person dbrPerson
	err = tx.Select("*").From("dbr_people").Where(ConditionRaw("id = ?", id)).LoadStruct(context.TODO(), &person)
	assert.NoError(t, err)

	assert.Equal(t, person.Id, id)
	assert.Equal(t, person.Name, "Barack")
	assert.Equal(t, person.Email.Valid, true)
	assert.Equal(t, person.Email.String, "obama@whitehouse.gov")

	err = tx.Commit()
	assert.NoError(t, err)
}

func TestTransactionRollbackReal(t *testing.T) {
	// Insert by specifying values
	s := createRealSessionWithFixtures()

	tx, err := s.Begin()
	assert.NoError(t, err)

	var person dbrPerson
	err = tx.Select("*").From("dbr_people").Where(ConditionRaw("email = ?", "jonathan@uservoice.com")).LoadStruct(context.TODO(), &person)
	assert.NoError(t, err)
	assert.Equal(t, person.Name, "Jonathan")

	err = tx.Rollback()
	assert.NoError(t, err)
}
