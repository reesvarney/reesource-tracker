package database

import (
	"context"
	"os"
	sqlite_driver "reesource-tracker/lib/database/drivers/sqlite"
)

var Connection *Queries

func Connect(ctx context.Context) {
	migration_dir := "file://database/migrations"
	if _, err := os.Stat("migrations"); err == nil {
		migration_dir = "file://migrations"
	}
	db, m, err := sqlite_driver.Connect(ctx, migration_dir)
	if err != nil {
		println("Got error", err.Error())
		return
	}
	Connection = New(db)

	err = m.Up()
	if err != nil && err.Error() != "no change" {
		println("Error applying migration", err.Error())
		return
	}
}
