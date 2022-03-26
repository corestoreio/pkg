package dml

import (
	"fmt"
	"testing"

	"github.com/corestoreio/pkg/util/assert"
	"github.com/go-sql-driver/mysql"
)

var _ error = (*Error)(nil)

func TestMySQLFromError(t *testing.T) {
	myErr := &mysql.MySQLError{
		Number:  1062,
		Message: "Duplicate Key",
	}
	haveM := MySQLMessageFromError(fmt.Errorf("outer fatal error: %w", myErr))
	assert.Exactly(t, "Duplicate Key", haveM)

	haveN := MySQLNumberFromError(fmt.Errorf("outer fatal error: %w", myErr))
	assert.Exactly(t, uint16(1062), haveN)
}
