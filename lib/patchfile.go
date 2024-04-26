package lib

import (
	"errors"
	"fmt"

	"github.com/sergi/go-diff/diffmatchpatch"
)

func patchTextFile(db dbCommon, path string, diff []diffmatchpatch.Diff) error {
	textStmt, err := db.Prepare("WITH text_values AS (SELECT ? AS text) INSERT INTO text_content (text) SELECT text FROM text_values WHERE NOT EXISTS (SELECT 1 FROM text_content WHERE text = (SELECT text FROM text_values))")
	if err != nil {
		return errors.Join(errors.New("failed to prepare text statement"), err)
	}
	patchStmt, err := db.Prepare("INSERT INTO patch (path, text, type) SELECT ?, id, ? FROM text_content WHERE text = ?;")
	if err != nil {
		return errors.Join(errors.New("failed to prepare patch statement"), err)
	}
	for _, d := range diff {
		if _, err := textStmt.Exec(d.Text); err != nil {
			return errors.Join(fmt.Errorf("failed to execute prepared statement with (%q, %q, %d)", path, d.Text, d.Type), err)
		}
		if _, err := patchStmt.Exec(path, d.Type, d.Text); err != nil {
			return errors.Join(fmt.Errorf("failed to execute prepared statement with (%q, %q, %d)", path, d.Text, d.Type), err)
		}
	}
	if err := textStmt.Close(); err != nil {
		return errors.Join(errors.New("failed to close text statement"), err)
	}
	if err := textStmt.Close(); err != nil {
		return errors.Join(errors.New("failed to close patch statement"), err)
	}
	return nil
}
