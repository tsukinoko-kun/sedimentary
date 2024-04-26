package lib

import (
	"errors"
	"fmt"

	"github.com/sergi/go-diff/diffmatchpatch"
)

func (sdmt *Sedimentary) ReadTextFile(path string) (string, error) {
	rows, err := sdmt.db.Query("SELECT text_content.text, patch.type FROM patch, text_content WHERE patch.path = ?", path)
	if err != nil || rows == nil {
		return "", errors.Join(fmt.Errorf("failed to query patched.path=%q", path), err)
	}
	defer rows.Close()

	diff := []diffmatchpatch.Diff{}

	for rows.Next() {
		d := diffmatchpatch.Diff{}
		if err := rows.Scan(&d.Text, &d.Type); err != nil {
			return "", errors.Join(errors.New("failed to scan selected rows from patch"), err)
		}
		diff = append(diff, d)
	}

	return dmp.DiffText2(diff), nil
}
