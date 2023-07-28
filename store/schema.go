package store

import (
	"context"
	"database/sql"
)

<<<<<<< HEAD
<<<<<<< HEAD
=======
>>>>>>> 6690fa6 (Add database connection tests)
/*
Install any extensions and create all tables for the database
*/
func createTables(db *sql.DB) error {
<<<<<<< HEAD
=======
func CreateTables(db *sql.DB) error {
>>>>>>> 4a39705 (Add .env file & Read .env code)
=======
>>>>>>> 6690fa6 (Add database connection tests)
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

<<<<<<< HEAD
<<<<<<< HEAD
/*
Install and Postgres Extensions
*/
=======
>>>>>>> 4a39705 (Add .env file & Read .env code)
=======
/*
Install and Postgres Extensions
*/
>>>>>>> 6690fa6 (Add database connection tests)
func installExtensions(db *sql.DB) error {
	uuidExtensionQuery := `CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`
	_, err := db.ExecContext(context.Background(), uuidExtensionQuery)
	return err
}

<<<<<<< HEAD
<<<<<<< HEAD
=======
>>>>>>> 6690fa6 (Add database connection tests)
/*
Create the table for Buyers
*/
func createBuyersTable(db *sql.DB) error {
	query := `CREATE TABLE IF NOT EXISTS buyers(
		buid uuid DEFAULT uuid_generate_v1() NOT NULL,
		email VARCHAR NOT NULL, 
		password VARCHAR NOT NULL,
=======
func createBuyersTable(db *sql.DB) error {
	query := `CREATE TABLE IF NOT EXISTS buyers(
		buid uuid DEFAULT uuid_generate_v1() NOT NULL,
		email VARCHAR NOT NULL, 
		password VARCHAR NOT NULL,
<<<<<<< HEAD
		salt VARCHAR NOT NULL,
>>>>>>> 4a39705 (Add .env file & Read .env code)
=======
>>>>>>> 005bc68 (Add login and signup API)
		PRIMARY KEY(buid));`

	_, err := db.ExecContext(context.Background(), query)
	return err
}

<<<<<<< HEAD
<<<<<<< HEAD
=======
>>>>>>> 6690fa6 (Add database connection tests)
/*
Create the table for Sellers
*/
func createSellersTable(db *sql.DB) error {
	query := `CREATE TABLE IF NOT EXISTS sellers(
		suid uuid DEFAULT uuid_generate_v1() NOT NULL,
		email VARCHAR NOT NULL, 
		password VARCHAR NOT NULL,
=======
func createSellersTable(db *sql.DB) error {
	query := `CREATE TABLE IF NOT EXISTS sellers(
		suid uuid DEFAULT uuid_generate_v1() NOT NULL,
		email VARCHAR NOT NULL, 
		password VARCHAR NOT NULL,
<<<<<<< HEAD
		salt VARCHAR NOT NULL,
>>>>>>> 4a39705 (Add .env file & Read .env code)
=======
>>>>>>> 005bc68 (Add login and signup API)
		PRIMARY KEY(suid));`

	_, err := db.ExecContext(context.Background(), query)
	return err
}
