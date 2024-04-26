package lib

import (
	"database/sql"
	"errors"
)

type (
	dbCommon interface {
		Exec(query string, args ...any) (sql.Result, error)
		Prepare(query string) (*sql.Stmt, error)
	}
)

func migrate(db dbCommon) error {
	if _, err := db.Exec("PRAGMA foreign_keys = ON"); err != nil {
		return errors.Join(errors.New("failed to enable foreign keys"), err)
	}

	// meta_data

	if _, err := db.Exec(
		"CREATE TABLE IF NOT EXISTS meta_data (key TEXT PRIMARY KEY, value TEXT NOT NULL)",
	); err != nil {
		return errors.Join(errors.New("failed to create metadata table"), err)
	}

	// text_content

	if _, err := db.Exec(
		"CREATE TABLE IF NOT EXISTS text_content (id INTEGER PRIMARY KEY AUTOINCREMENT, text TEXT NOT NULL)",
	); err != nil {
		return errors.Join(errors.New("failed to create text_content table"), err)
	}

	// patch

	if _, err := db.Exec(
		"CREATE TABLE IF NOT EXISTS patch (id INTEGER PRIMARY KEY AUTOINCREMENT, path TEXT NOT NULL, text INTEGER NOT NULL, type INTEGER NOT NULL, FOREIGN KEY(text) REFERENCES text_content(id))",
	); err != nil {
		return errors.Join(errors.New("failed to create patch table"), err)
	}

	if _, err := db.Exec(
		"CREATE INDEX IF NOT EXISTS idx_patch_path ON patch (path)",
	); err != nil {
		return errors.Join(errors.New("failed to create index on patch.path"), err)
	}

	return nil
}
