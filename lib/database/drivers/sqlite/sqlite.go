package sqlite_driver

import (
	"context"
	"database/sql"

	_ "modernc.org/sqlite"
)

func Connect(ctx context.Context, schema string) (error, *sql.DB) {
	db, err := sql.Open("sqlite", "database/db.sqlite")
	if err != nil {
		return err, nil
	}

	// create tables
	if _, err := db.ExecContext(ctx, schema); err != nil {
		return err, nil
	}

	println("Connected to sqlite database")
	return nil, db
}
