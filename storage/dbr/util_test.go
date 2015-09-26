package dbr

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStmtChecker(t *testing.T) {
	tests := []struct {
		sel   string
		selok bool
		upd   string
		updok bool
		del   string
		delok bool
		ins   string
		insok bool
	}{
		{
			"SELECT ...",
			false,
			"UPDATE ...",
			false,
			"DELETE ...",
			false,
			"INSERT",
			false,
		},
		{
			"SELECT ... From ",
			true,
			"UPDATE ... From ",
			true,
			"DELETE ...From ",
			true,
			"INSERT ",
			true,
		},
	}

	for _, test := range tests {
		assert.Equal(t, test.selok, Stmt.IsSelect(test.sel), "%#v", test)
		assert.Equal(t, test.updok, Stmt.IsUpdate(test.upd), "%#v", test)
		assert.Equal(t, test.delok, Stmt.IsDelete(test.del), "%#v", test)
		assert.Equal(t, test.insok, Stmt.IsInsert(test.ins), "%#v", test)
	}
}
