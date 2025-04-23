package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"

	"github.com/knstch/subtrack-kafka/outbox"
)

func init() {
	goose.AddMigrationContext(upInitOutbox, downInitOutbox)
}

func upInitOutbox(_ context.Context, tx *sql.Tx) error {
	_, err := tx.Exec(outbox.OutboxMigrationUp)
	if err != nil {
		return err
	}

	return nil
}

func downInitOutbox(_ context.Context, tx *sql.Tx) error {
	_, err := tx.Exec(outbox.OutboxMigrationDown)
	if err != nil {
		return err
	}
	return nil
}
