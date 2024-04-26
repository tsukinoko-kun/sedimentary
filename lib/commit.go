package lib

import (
	"database/sql"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/sergi/go-diff/diffmatchpatch"
)

var dmp = diffmatchpatch.New()

// Commit creates patches for all files in she paths parameter and labels it with the message parameter.
// All paths must be specified absolutely.
func (sdmt *Sedimentary) Commit(paths []string, message string) error {
	err := sdmt.transaction(func(tx *sql.Tx) error {
		for _, path := range paths {
			path := filepath.Join(sdmt.root, path)
			if fi, err := os.Stat(path); err != nil {
				return err
			} else if fi.IsDir() {
				return sdmt.commitDir(tx, path)
			} else {
				return sdmt.commitFile(tx, path)
			}
		}
		return nil
	})

	return err
}

func (sdmt *Sedimentary) commitDir(tx dbCommon, path string) error {
	err := filepath.WalkDir(path, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if err := sdmt.commitFile(tx, path); err != nil {
			return errors.Join(fmt.Errorf("failed to commit file %q", path), err)
		}
		return nil
	})
	if err != nil {
		return errors.Join(fmt.Errorf("failed to walk dir %q", path), err)
	}
	return nil
}

func (sdmt *Sedimentary) commitFile(tx dbCommon, path string) error {
	{
		fileName := filepath.Base(path)
		if fileName == ".DS_Store" || fileName == ".sedimentary" {
			return nil
		}

		pathParts := strings.Split(path, string(filepath.Separator))
		for _, pathPart := range pathParts {
			if pathPart == ".git" {
				return nil
			}
		}
	}

	relPath := sdmt.relativePath(path)

	tracked, err := sdmt.ReadTextFile(relPath)
	if err != nil {
		return errors.Join(fmt.Errorf("failed to get tracked file %q", relPath), err)
	}

	updated, err := os.ReadFile(path)
	if err != nil {
		return errors.Join(fmt.Errorf("failed to read local file %q", path), err)
	}
	diff := dmp.DiffMain(tracked, string(updated), true)
	if len(diff) == 0 {
		return nil
	}
	fmt.Printf("changes in %q\n", relPath)
	if err := patchTextFile(tx, relPath, diff); err != nil {
		return errors.Join(fmt.Errorf("failed to commit patch for file %q", relPath), err)
	}
	return nil
}
