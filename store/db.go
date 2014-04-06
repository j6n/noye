package store

import (
	"fmt"
	"path"
	"strings"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

var (
	db     *DB
	dbPath string
)

func init() {
	dbPath = "."

	var err error
	if db = NewDB(); db == nil {
		panic(fmt.Errorf("loading db: %s\n", err))
	}
}

const KvSchema = `
CREATE TABLE IF NOT EXISTS %s (
	k	varchar(255) NOT NULL,
	v	BLOB,
	PRIMARY KEY(k)
);
`

func fixScriptName(script string) string {
	script = script + "_script"
	script = strings.Replace(script, ".", "_", -1)
	return script
}

func Get(table, key string) (string, error) {
	table = fixScriptName(table)
	return db.Get(table, key)
}

func Set(table, key, data string) (err error) {
	table = fixScriptName(table)
	return db.Set(table, key, data)
}

type DB struct{ *sqlx.DB }

func NewDB() (db *DB) {
	if db == nil {
		temp, err := sqlx.Open("sqlite3", path.Join(dbPath, "noye.db"))
		if err != nil {
			return nil
		}
		db = &DB{temp}
	}
	return
}

func (d *DB) Close() {
	d.Close()
}

func (d *DB) Set(table, key, data string) (err error) {
	if err := d.CheckTable(table, KvSchema); err != nil {
		return err
	}

	tx, err := d.Beginx()
	if err != nil {
		return err
	}
	input := []byte(data)

	// try update, this is awful
	res, err := tx.Execv(fmt.Sprintf("UPDATE %s SET v = ? WHERE k = ?", table), input, key)
	if err != nil {
		return err
	}

	// it touched a row
	if n, _ := res.RowsAffected(); n > 0 {
		tx.Commit()
		return nil
	}

	// try insert
	_, err = tx.Execv(fmt.Sprintf("INSERT INTO %s (k, v) VALUES (?, ?);", table), key, input)
	if err == nil {
		tx.Commit()
	}

	return
}

func (d *DB) Get(table, key string) (string, error) {
	err := d.CheckTable(table, KvSchema)
	if err != nil {
		return "", err
	}

	temp := map[string]interface{}{}
	row := d.QueryRowx(fmt.Sprintf("SELECT v FROM %s WHERE k = ?", table), key)
	if err := row.MapScan(temp); err != nil {
		return "", err
	}

	if res, ok := temp["v"].(string); ok {
		return res, nil
	}

	return "", fmt.Errorf("couldn't find '%s' on '%s'", key, table)
}

func (d *DB) CheckTable(table, schema string) error {
	res, err := d.Execv("SELECT name FROM sqlite_master WHERE type = 'table' AND name = ?;", table)
	if err != nil {
		return err
	}

	if i, _ := res.RowsAffected(); i == 0 {
		if _, err := d.Exec(fmt.Sprintf(schema, table)); err != nil {
			return err
		}
	}

	return nil
}
