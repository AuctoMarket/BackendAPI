package store

import (
	"context"
	"database/sql"
)

func CreateTables(db *sql.DB) error {
	err := installExtensions(db)

	if err != nil {
		return err
	}

	err = createBuyersTable(db)

	if err != nil {
		return err
	}

	err = createSellersTable(db)

	if err != nil {
		return err
	}

	return nil
}

func installExtensions(db *sql.DB) error {
	uuidExtensionQuery := `CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`
	_, err := db.ExecContext(context.Background(), uuidExtensionQuery)
	return err
}

func createBuyersTable(db *sql.DB) error {
	query := `CREATE TABLE IF NOT EXISTS buyers(
		buid uuid DEFAULT uuid_generate_v1() NOT NULL,
		email VARCHAR NOT NULL, 
		password VARCHAR NOT NULL,
		PRIMARY KEY(buid));`

	_, err := db.ExecContext(context.Background(), query)
	return err
}

func createSellersTable(db *sql.DB) error {
	query := `CREATE TABLE IF NOT EXISTS sellers(
		suid uuid DEFAULT uuid_generate_v1() NOT NULL,
		email VARCHAR NOT NULL, 
		password VARCHAR NOT NULL,
		PRIMARY KEY(suid));`

	_, err := db.ExecContext(context.Background(), query)
	return err
}
