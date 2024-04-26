package lib

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

type Sedimentary struct {
	root string
	db   *sql.DB
}

// New initializes a new sedimentary at the given root path
func New(root string) (*Sedimentary, error) {
	filePath := filepath.Join(root, ".sedimentary")

	if _, err := os.Create(filePath); err != nil {
		// db file creation failed
		return nil, err
	}

	db, err := sql.Open("sqlite", filePath)
	if err != nil {
		return nil, errors.Join(fmt.Errorf("failed to open db %q", filePath), err)
	}

	if err := migrate(db); err != nil {
		return nil, err
	}

	return &Sedimentary{root, db}, nil
}

// Open opens an existing sedimentary at the given root path or creates a new if there is none
func Open(root string) (*Sedimentary, error) {
	filePath := filepath.Join(root, ".sedimentary")

	// does db file exist?
	if fi, err := os.Stat(filePath); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			// file does not exist
			// create
			if _, err := os.Create(filePath); err != nil {
				// db file creation failed
				return nil, err
			}
		} else {
			// failed to look at db file
			return nil, err
		}
	} else if !fi.Mode().IsRegular() {
		// db file is not a regular file
		return nil, fmt.Errorf("path %q is not a regular file", filePath)
	}

	db, err := sql.Open("sqlite", filePath)
	if err != nil {
		return nil, errors.Join(fmt.Errorf("failed to open db %q", filePath), err)
	}

	if err := migrate(db); err != nil {
		return nil, err
	}

	return &Sedimentary{root, db}, nil
}

func (sdmt *Sedimentary) Close() error {
	return sdmt.db.Close()
}

func (sdmt *Sedimentary) relativePath(path string) string {
	if filepath.IsAbs(path) {
		rel, err := filepath.Rel(sdmt.root, path)
		if err == nil {
			return rel
		}
	}
	return path
}
