package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upInitUsersDataTable, downInitUsersDataTable)
}

func upInitUsersDataTable(_ context.Context, tx *sql.Tx) error {
	_, err := tx.Exec(`CREATE TABLE IF NOT EXISTS users_data (
    id SERIAL PRIMARY KEY,
    user_id INTEGER UNIQUE NOT NULL,
    name VARCHAR(255),
    last_name VARCHAR(255),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP NULL,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
)`)
	if err != nil {
		return err
	}

	return nil
}

func downInitUsersDataTable(_ context.Context, tx *sql.Tx) error {
	_, err := tx.Exec(`DROP TABLE IF EXISTS users_data`)
	if err != nil {
		return err
	}

	return nil
}
