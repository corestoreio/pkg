package dbr

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/corestoreio/errors"
	_ "github.com/go-sql-driver/mysql"
)

//
// Test helpers
//

// Returns a session that's not backed by a database
func createFakeSession() *Connection {
	cxn, err := NewConnection()
	if err != nil {
		panic(err)
	}
	return cxn
}

func createRealSession() *Connection {
	_, dsn := realDb()
	cxn, err := NewConnection(
		WithDSN(dsn),
	)
	if err != nil {
		panic(err)
	}
	return cxn
}

func createRealSessionWithFixtures() *Connection {
	sess := createRealSession()
	installFixtures(sess.DB)
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

var _ InsertArgProducer = (*dbrPerson)(nil)
var _ UpdateArgProducer = (*dbrPerson)(nil)
var _ InsertArgProducer = (*nullTypedRecord)(nil)

//var _ UpdateArgProducer = (*nullTypedRecord)(nil)

type dbrPerson struct {
	ID    int64 `db:"id"`
	Name  string
	Email NullString
	Key   NullString
}

func (p *dbrPerson) columnToArg(t byte, args Arguments, columns []string) (Arguments, error) {
	for _, c := range columns {
		switch c {
		case "id":
			if t == 'i' {
				args = append(args, ArgInt64(p.ID))
			}
		case "name":
			args = append(args, ArgString(p.Name))
		case "email":
			args = append(args, ArgNullString(p.Email))
		//case "key":
		//	args = append(args, ArgNullString(p.Key))
		default:
			return nil, errors.NewNotFoundf("[dbr_test] Column %q not found", c)
		}
	}
	return args, nil
}

func (p *dbrPerson) ProduceInsertArgs(args Arguments, columns []string) (Arguments, error) {
	return p.columnToArg('i', args, columns)
}

func (p *dbrPerson) ProduceUpdateArgs(args Arguments, columns, condition []string) (_ Arguments, err error) {
	args, err = p.columnToArg('u', args, columns)
	for _, c := range condition {
		switch c {
		case "id":
			args = append(args, ArgInt64(p.ID))
		default:
			return nil, errors.NewNotFoundf("[dbr_test] Column %q not found", c)
		}
	}
	return args, err
}

type nullTypedRecord struct {
	ID         int64 `db:"id"`
	StringVal  NullString
	Int64Val   NullInt64
	Float64Val NullFloat64
	TimeVal    NullTime
	BoolVal    NullBool
}

func (p *nullTypedRecord) ProduceInsertArgs(args Arguments, columns []string) (Arguments, error) {
	for _, c := range columns {
		switch c {
		case "id":
			args = append(args, ArgInt64(p.ID))
		case "string_val":
			args = append(args, ArgNullString(p.StringVal))
		case "int64_val":
			if p.Int64Val.Valid {
				args = append(args, ArgInt64(p.Int64Val.Int64))
			} else {
				args = append(args, ArgNull())
			}
		case "float64_val":
			if p.Float64Val.Valid {
				args = append(args, ArgFloat64(p.Float64Val.Float64))
			} else {
				args = append(args, ArgNull())
			}
		case "time_val":
			if p.TimeVal.Valid {
				args = append(args, ArgTime(p.TimeVal.Time))
			} else {
				args = append(args, ArgNull())
			}
		case "bool_val":
			if p.BoolVal.Valid {
				args = append(args, ArgBool(p.BoolVal.Bool))
			} else {
				args = append(args, ArgNull())
			}
		default:
			return nil, errors.NewNotFoundf("[dbr_test] Column %q not found", c)
		}
	}
	return args, nil
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

func (pm dbMock) PrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {
	if pm.error != nil {
		return nil, pm.error
	}
	return pm.prepareFn(query)
}

func (pm dbMock) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	if pm.error != nil {
		return nil, pm.error
	}
	return nil, nil
}

func (pm dbMock) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	if pm.error != nil {
		return nil, pm.error
	}
	return nil, nil
}
