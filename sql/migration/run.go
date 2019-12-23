package migration

import (
	"github.com/corestoreio/errors"
	"github.com/corestoreio/pkg/sql/dml"
	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/mysql"
	_ "github.com/golang-migrate/migrate/source/file"
)

// Run has the idea of gathering all eligible .sql files in the folder
// `_dbmigrate` and put them in a zip and execute that zip file with this
// program on the server. Alternatively this program reads the .sql files from
// the go-bindata archive so that you only need to deploy one runable file.
func Run() error {
	// Idea
	dbPool, err := dml.NewConnPool(dml.WithDSNFromEnv(""))
	if err != nil {
		return errors.WithStack(err)
	}

	driver, err := mysql.WithInstance(dbPool.DB, &mysql.Config{})
	m, err := migrate.NewWithDatabaseInstance(
		"file:///migrations",
		"mysql", driver)
	if err != nil {
		return errors.WithStack(err)
	}

	m.Steps(2)

	return nil
}
