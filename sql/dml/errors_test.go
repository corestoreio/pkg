package dml

import (
	"testing"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/pkg/util/assert"
	"github.com/go-sql-driver/mysql"
)

func TestMySQLFromError(t *testing.T) {
	t.Parallel()

	myErr := &mysql.MySQLError{
		Number:  1062,
		Message: "Duplicate Key",
	}
	haveM := MySQLMessageFromError(errors.Fatal.New(myErr, "Outer fatal error"))
	assert.Exactly(t, "Duplicate Key", haveM)

	haveN := MySQLNumberFromError(errors.Fatal.New(myErr, "Outer fatal error"))
	assert.Exactly(t, uint16(1062), haveN)
}
