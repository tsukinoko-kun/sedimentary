package libdb

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/tsukinoko-kun/sedimentary/build"
)

func ReadVersion(db *sql.DB) (string, error) {
	row := db.QueryRow(fmt.Sprintf("SELECT value FROM %s WHERE key = 'version'", tableMetaData))
	if row == nil {
		return build.Version, errors.New("no version stored")
	}

	var versionStr string
	err := row.Scan(&versionStr)
	return versionStr, err
}

func WriteVersion(db *sql.DB) error {
	tx, err := db.Begin()
	if err != nil {
		_ = tx.Rollback()
		return errors.Join(errors.New("failed to create transaction"), err)
	}

	if _, err := tx.Exec(fmt.Sprintf(
		"CREATE TABLE IF NOT EXISTS %s (key TEXT PRIMARY KEY, value TEXT)",
		tableMetaData)); err != nil {
		_ = tx.Rollback()
		return errors.Join(fmt.Errorf("failed to create %s table", tableMetaData), err)
	}

	stmt, err := tx.Prepare(fmt.Sprintf("INSERT OR REPLACE INTO %s (key, value) VALUES (?, ?)", tableMetaData))
	if err != nil {
		_ = tx.Rollback()
		return errors.Join(errors.New("failed to prepare insert statement"), err)
	}
	defer stmt.Close()

	if _, err := stmt.Exec("version", build.Version); err != nil {
		_ = tx.Rollback()
		return errors.Join(errors.New("failed to insert version"), err)
	}

	if err := tx.Commit(); err != nil {
		return errors.Join(errors.New("failed to commit transaction"), err)
	}
	return nil
}
