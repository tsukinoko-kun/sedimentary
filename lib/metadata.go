package lib

import (
	"database/sql"
	"errors"

	"github.com/tsukinoko-kun/sedimentary/build"
)

func (sdmt *Sedimentary) ReadVersion() (string, error) {
	row := sdmt.db.QueryRow("SELECT value FROM meta_data WHERE key = 'version'")
	if row == nil {
		return build.Version, errors.New("no version stored")
	}

	var versionStr string
	err := row.Scan(&versionStr)
	return versionStr, err
}

func (sdmt *Sedimentary) WriteVersion() error {
	return sdmt.transaction(func(tx *sql.Tx) error {
		stmt, err := tx.Prepare("INSERT OR REPLACE INTO meta_data (key, value) VALUES (?, ?)")
		if err != nil {
			return errors.Join(errors.New("failed to prepare insert statement for meta_data"), err)
		}
		defer stmt.Close()

		if _, err := stmt.Exec("version", build.Version); err != nil {
			return errors.Join(errors.New("failed to insert version into meta_data"), err)
		}
		return nil
	})
}
