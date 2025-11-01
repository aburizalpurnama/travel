package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upInvalidCommand, downInvalidCommand)
}

func upInvalidCommand(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is applied.
	_, err := tx.ExecContext(ctx, "ALTER core.invalid_table ADD IF NOT EXIST invalid varchar(10) NULL;")
	if err != nil {
		return err
	}

	return nil
}

func downInvalidCommand(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is rolled back.
	return nil
}
