package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upInitUsersTable, downInitUsersTable)
}

func upInitUsersTable(_ context.Context, tx *sql.Tx) error {
	_, err := tx.Exec(`
		CREATE TABLE IF NOT EXISTS users (
		    id SERIAL PRIMARY KEY,
		    email VARCHAR(255) UNIQUE NOT NULL,
		    password VARCHAR(255) NOT NULL,
		    role VARCHAR(255) NOT NULL,
		    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
		    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
		)
`)
	if err != nil {
		return err
	}

	return nil
}

func downInitUsersTable(_ context.Context, tx *sql.Tx) error {
	_, err := tx.Exec(`DROP TABLE IF EXISTS users`)
	if err != nil {
		return err
	}

	return nil
}
