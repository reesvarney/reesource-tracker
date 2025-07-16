package database

import (
	sqlite_driver "reesource-tracker/lib/database/drivers/sqlite"
	"context"
)

var Connection *Queries

func Connect(ctx context.Context, schema string) {
	err, db := sqlite_driver.Connect(ctx, schema)
	if err == nil {
		Connection = New(db)
	} else {
		println("Got error", err.Error())

	}
}
