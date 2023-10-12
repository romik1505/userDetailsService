package migrations

import (
	"context"
	"database/sql"
	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upFirstPersonMigraion, downFirstPersonMigraion)
}

func upFirstPersonMigraion(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is applied.
	_, err := tx.ExecContext(ctx, `
		CREATE TYPE gender_type as ENUM ('male', 'female', '');
		
		CREATE TABLE persons (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255) NOT NULL ,
			surname VARCHAR(255) NOT NULL ,
			patronymic VARCHAR(255) NOT NULL ,
			age INTEGER NOT NULL ,
			gender gender_type,
			nationality VARCHAR(2) NOT NULL ,
		
			created_at timestamptz NOT NULL DEFAULT now(),
			updated_at timestamptz,
			surname_index_col tsvector
		);
		
		CREATE INDEX ix_person_surname ON persons USING gin(surname_index_col);
	`)

	return err
}

func downFirstPersonMigraion(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is rolled back.
	_, err := tx.ExecContext(ctx, `
		
		DROP TABLE persons;		
		DROP TYPE gender_type;
`)
	if err != nil {
		return err
	}
	return nil
}
