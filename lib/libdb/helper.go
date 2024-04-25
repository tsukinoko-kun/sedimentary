package libdb

import (
	"context"
	"database/sql"
	"errors"
)

func exec(tx *sql.Tx, query string, args ...any) (sql.Result, error) {
	ctx := context.Background()
	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return nil, errors.Join(errors.New("failed to prepare statement"), err)
	}
	defer stmt.Close()

	return stmt.Exec(args...)
}
