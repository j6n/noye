package store

import (
	"fmt"
	"log"
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

const tableSchema = `
CREATE TABLE IF NOT EXISTS %s (
	k	varchar(255) NOT NULL,
	v	BLOB,
	PRIMARY KEY(k)
);
`

func GetSession() (db *sqlx.DB, err error) {
	if db == nil {
		db, err = sqlx.Open("sqlite3", path.Join(dbPath, "noye.db"))
		if err != nil {
			return
		}
	}

	return
}

func checkTable(table string) (*sqlx.DB, error, bool) {
	sess, err := GetSession()
	if err != nil {
		log.Println("err get sess:", err)
		return sess, err, false
	}

	res, err := sess.Exec("SELECT name FROM sqlite_master WHERE type = 'table' AND name = ?;", table)
	if err != nil {
		return sess, err, false
	}

	if i, _ := res.RowsAffected(); i == 0 {
		sess.Execf(fmt.Sprintf(tableSchema, table))
	}

	return sess, err, true
}

func Get(table, key string) (string, error) {
	table = table + "_script"
	sess, err, ok := checkTable(table)
	if err != nil || !ok {
		return "", err
	}

	temp := map[string]interface{}{}
	row := sess.QueryRowx(fmt.Sprintf("SELECT v FROM %s WHERE k = ?", table), key)
	if err := row.MapScan(temp); err != nil {
		return "", err
	}

	if res, ok := temp["v"].(string); ok {
		return res, nil
	}

	return "", fmt.Errorf("couldn't find '%s' on '%s'", key, table)
}

func Set(table, key, data string) (err error) {
	table = table + "_script"
	sess, err, ok := checkTable(table)
	if err != nil || !ok {
		log.Println("err check table:", err)
		return err
	}

	tx, err := sess.Beginx()
	if err != nil {
		log.Println("err begin:", err)
		return err
	}
	input := []byte(data)

	// try update, this is awful
	res, err := tx.Exec(fmt.Sprintf("UPDATE %s SET v = ? WHERE k = ?", table), input, key)
	if err != nil {
		log.Println("err exec:", err)
		return err
	}

	// it touched a row
	if n, _ := res.RowsAffected(); n > 0 {
		tx.Commit()
		return nil
	}

	// try insert
	_, err = tx.Exec(fmt.Sprintf("INSERT INTO %s (k, v) VALUES (?, ?);", table), key, input)
	if err == nil {
		tx.Commit()
	}

	return
}
