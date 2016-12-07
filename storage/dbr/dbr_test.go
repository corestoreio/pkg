package dbr

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/corestoreio/csfw/util/null"
	_ "github.com/go-sql-driver/mysql"
)

//
// Test helpers
//

// Returns a session that's not backed by a database
func createFakeSession() *Session {
	cxn, err := NewConnection()
	if err != nil {
		panic(err)
	}
	return cxn.NewSession()
}

func createRealSession() *Session {
	_, dsn := realDb()
	cxn, err := NewConnection(
		WithDSN(dsn),
	)
	if err != nil {
		panic(err)
	}
	return cxn.NewSession()
}

func createRealSessionWithFixtures() *Session {
	sess := createRealSession()
	installFixtures(sess.cxn.DB)
	return sess
}

func realDb() (driver string, dsn string) {
	driver = os.Getenv("DBR_TEST_DRIVER")
	if driver == "" {
		driver = DefaultDriverName
	}

	dsn = os.Getenv("CS_DSN")
	if dsn == "" {
		dsn = "root:unprotected@unix(/tmp/mysql.sock)/uservoice_development?charset=utf8&parseTime=true"
	}
	return
}

type dbrPerson struct {
	ID    int64 `db:"id"`
	Name  string
	Email null.String
	Key   null.String
}

type nullTypedRecord struct {
	ID         int64 `db:"id"`
	StringVal  null.String
	Int64Val   null.Int64
	Float64Val null.Float64
	TimeVal    null.Time
	BoolVal    null.Bool
}

func installFixtures(db *sql.DB) {
	createPeopleTable := fmt.Sprintf(`
		CREATE TABLE dbr_people (
			id int(11) NOT NULL auto_increment PRIMARY KEY,
			name varchar(255) NOT NULL,
			email varchar(255),
			%s varchar(255)
		)
	`, "`key`")

	createNullTypesTable := `
		CREATE TABLE null_types (
			id int(11) NOT NULL auto_increment PRIMARY KEY,
			string_val varchar(255) NULL,
			int64_val int(11) NULL,
			float64_val float NULL,
			time_val datetime NULL,
			bool_val bool NULL
		)
	`

	sqlToRun := []string{
		"DROP TABLE IF EXISTS dbr_people",
		createPeopleTable,
		"INSERT INTO dbr_people (name,email) VALUES ('Jonathan', 'jonathan@uservoice.com')",
		"INSERT INTO dbr_people (name,email) VALUES ('Dmitri', 'zavorotni@jadius.com')",

		"DROP TABLE IF EXISTS null_types",
		createNullTypesTable,
	}

	for _, v := range sqlToRun {
		_, err := db.Exec(v)
		if err != nil {
			log.Fatalln("Failed to execute statement: ", v, " Got error: ", err)
		}
	}
}

var _ Preparer = (*dbMock)(nil)
var _ Querier = (*dbMock)(nil)
var _ Execer = (*dbMock)(nil)

type dbMock struct {
	error
	prepareFn func(query string) (*sql.Stmt, error)
}

func (pm dbMock) Prepare(query string) (*sql.Stmt, error) {
	if pm.error != nil {
		return nil, pm.error
	}
	return pm.prepareFn(query)
}

func (pm dbMock) Query(query string, args ...interface{}) (*sql.Rows, error) {
	if pm.error != nil {
		return nil, pm.error
	}
	return nil, nil
}

func (pm dbMock) Exec(query string, args ...interface{}) (sql.Result, error) {
	if pm.error != nil {
		return nil, pm.error
	}
	return nil, nil
}
