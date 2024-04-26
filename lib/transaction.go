package lib

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
)

func (sdmt *Sedimentary) transaction(fn func(tx *sql.Tx) error) error {
	tx, err := sdmt.db.Begin()
	if err != nil {
		return err
	}

	done := false

	defer func() {
		if tx == nil || done {
			return
		}
		r := recover()

		if err := tx.Rollback(); err != nil {
			fmt.Fprintf(os.Stderr, "db transaction rollback failed: %s", err.Error())
		}

		if r != nil {
			fmt.Fprintf(os.Stderr, "panic: %v", r)
		}
	}()

	if txErr := fn(tx); txErr != nil {
		done = true
		if rbErr := tx.Rollback(); rbErr != nil {
			return errors.Join(rbErr, txErr)
		}
		return txErr
	}

	done = true
	if err := tx.Commit(); err != nil {
		return errors.Join(errors.New("failed to commit"), err)
	}

	return nil
}
