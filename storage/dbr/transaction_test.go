package dbr

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTransactionReal(t *testing.T) {
	s := createRealSessionWithFixtures()

	tx, err := s.Begin()
	assert.NoError(t, err)

	res, err := tx.InsertInto("dbr_people").Columns("name", "email").Values(
		ArgString("Barack"), ArgString("obama@whitehouse.gov"),
		ArgString("Obama"), ArgString("barack@whitehouse.gov"),
	).Exec(context.TODO())

	assert.NoError(t, err)
	id, err := res.LastInsertId()
	assert.NoError(t, err)
	rowsAff, err := res.RowsAffected()
	assert.NoError(t, err)

	assert.True(t, id > 0)
	assert.Equal(t, int64(2), rowsAff)

	var person dbrPerson
	err = tx.Select("*").From("dbr_people").Where(Condition("id = ?", ArgInt64(id))).LoadStruct(context.TODO(), &person)
	assert.NoError(t, err)

	assert.Equal(t, id, person.ID)
	assert.Equal(t, "Barack", person.Name)
	assert.Equal(t, true, person.Email.Valid)
	assert.Equal(t, "obama@whitehouse.gov", person.Email.String)

	err = tx.Commit()
	assert.NoError(t, err)
}

func TestTransactionRollbackReal(t *testing.T) {
	// Insert by specifying values
	s := createRealSessionWithFixtures()

	tx, err := s.Begin()
	assert.NoError(t, err)

	var person dbrPerson
	err = tx.Select("*").From("dbr_people").Where(Condition("email = ?", ArgString("jonathan@uservoice.com"))).LoadStruct(context.TODO(), &person)
	assert.NoError(t, err)
	assert.Equal(t, "Jonathan", person.Name)

	err = tx.Rollback()
	assert.NoError(t, err)
}
