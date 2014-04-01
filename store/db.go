package store

import (
	"database/sql"
	"log"
	"os"
	"path"

	"github.com/coopernurse/gorp"
)

var (
	db         *gorp.DbMap
	sqlitePath string
)

func init() {
	if sqlitePath = os.Getenv("OPENSHIFT_DATA"); sqlitePath == "" {
		sqlitePath = "."
	}
}

func GetSession() (*gorp.DbMap, error) {
	var (
		err    error
		sqlite *sql.DB
	)

	if sqlite, err = sql.Open("sqlite", path.Join(sqlitePath, "noye.db")); err != nil {
		return nil, err
	}

	if db == nil {
		db = &gorp.DbMap{Db: sqlite, Dialect: gorp.SqliteDialect{}}
		db.TraceOn("[gorp]", log.New(os.Stdout, "noye:", log.Lmicroseconds))
		db.CreateTablesIfNotExists()
	}

	return db, err
}
