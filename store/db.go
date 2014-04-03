package store

import (
	"fmt"
	"os"
	"path"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

var (
	db     *sqlx.DB
	dbPath string
)

func init() {
	if dbPath = os.Getenv("OPENSHIFT_DATA"); dbPath == "" {
		dbPath = "."
	}
}

var schema = `
CREATE TABLE IF NOT EXISTS scripts (
	n	varchar(255) NOT NULL,
	o   varchar(255) NOT NULL,
	d	BLOB,
	PRIMARY KEY(n)
);
`

func GetSession() (db *sqlx.DB, err error) {
	if db == nil {
		db, err = sqlx.Open("sqlite3", path.Join(dbPath, "noye.db"))
		if err != nil {
			return
		}
		db.Execf(schema)
	}

	return
}

func Get(table, key string) (string, error) {
	sess, err := GetSession()
	if err != nil {
		return "", err
	}

	temp := map[string]interface{}{}
	row := sess.QueryRowx("SELECT d FROM scripts WHERE n = ? AND o = ?", key, table)
	if err := row.MapScan(temp); err != nil {
		return "", err
	}

	if res, ok := temp["d"].(string); ok {
		return res, nil
	}

	return "", fmt.Errorf("couldn't find '%s'/'%s' on scripts", key, table)
}

func Set(table, key, data string) (err error) {
	sess, err := GetSession()
	if err != nil {
		return err
	}

	tx, err := sess.Beginx()
	if err != nil {
		return err
	}
	input := []byte(data)

	// try update
	res, err := tx.Exec("UPDATE scripts SET d = ? WHERE n = ? AND o = ?", input, key, table)
	if err != nil {
		return err
	}

	// it touched a row
	if n, _ := res.RowsAffected(); n > 0 {
		tx.Commit()
		return nil
	}

	// try insert
	_, err = tx.Exec("INSERT INTO scripts (d, n, o) VALUES (?, ?, ?)", input, key, table)
	if err == nil {
		tx.Commit()
	}

	return
}
