package store

import (
	"context"
	"database/sql"
)

/*
Install any extensions and create all tables for the database
*/
func createTables(db *sql.DB) error {
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

	err = createProductsTable(db)

	if err != nil {
		return err
	}

	err = createProductImagesTable(db)

	if err != nil {
		return err
	}

	return nil
}

/*
Install and Postgres Extensions
*/
func installExtensions(db *sql.DB) error {
	uuidExtensionQuery := `CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`
	_, err := db.ExecContext(context.Background(), uuidExtensionQuery)
	return err
}

/*
Create the table for Buyers
*/
func createBuyersTable(db *sql.DB) error {
	query := `CREATE TABLE IF NOT EXISTS buyers(
		buyer_id uuid DEFAULT uuid_generate_v1() NOT NULL,
		email VARCHAR NOT NULL, 
		password VARCHAR NOT NULL,
		PRIMARY KEY(buyer_id));`

	_, err := db.ExecContext(context.Background(), query)
	return err
}

/*
Create the table for Sellers
*/
func createSellersTable(db *sql.DB) error {
	query := `CREATE TABLE IF NOT EXISTS sellers(
		seller_id uuid DEFAULT uuid_generate_v1() NOT NULL,
		email VARCHAR NOT NULL, 
		password VARCHAR NOT NULL,
		seller_name VARCHAR NOT NULL,
		PRIMARY KEY(seller_id));`

	_, err := db.ExecContext(context.Background(), query)
	return err
}

/*
Create the table for Products
*/
func createProductsTable(db *sql.DB) error {
	query := `CREATE TABLE IF NOT EXISTS products(
		product_id uuid DEFAULT uuid_generate_v1() NOT NULL,
		seller_id uuid REFERENCES sellers(seller_id), 
		title TEXT NOT NULL,
		description TEXT NOT NULL,
		condition INT NOT NULL CONSTRAINT isOutOfFive CHECK (condition >= 0 AND condition <= 5),
		price INT NOT NULL CONSTRAINT isPositive CHECK (price >= 0),
		product_type VARCHAR NOT NULL,
		posted_date TIMESTAMPTZ NOT NULL,
		PRIMARY KEY(product_id));`

	_, err := db.ExecContext(context.Background(), query)
	return err
}

/*
Create the table for Product Images
*/
func createProductImagesTable(db *sql.DB) error {
	query := `CREATE TABLE IF NOT EXISTS product_images(
		product_image_id uuid DEFAULT uuid_generate_v1() NOT NULL,
		product_id uuid REFERENCES products(product_id), 
		path VARCHAR NOT NULL,
		PRIMARY KEY(product_image_id));`

	_, err := db.ExecContext(context.Background(), query)
	return err
}
