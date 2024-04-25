package libdb

import (
	"database/sql"
	"errors"
	"os"

	_ "modernc.org/sqlite"
)

func Open() (*sql.DB, error) {
	db, err := sql.Open("sqlite", ".sedimentary")
	if err != nil {
		return nil, errors.Join(errors.New("failed to open newly created db"), err)
	}

	return db, nil
}

func Init() (*sql.DB, error) {
	_, err := os.Create(".sedimentary")
	if err != nil {
		return nil, errors.Join(errors.New("failed to create db file"), err)
	}

	db, err := sql.Open("sqlite", ".sedimentary")
	if err != nil {
		return nil, errors.Join(errors.New("failed to open newly created db"), err)
	}

	if err := WriteVersion(db); err != nil {
		return db, errors.Join(errors.New("failed to write version"), err)
	}

	return db, nil
}
