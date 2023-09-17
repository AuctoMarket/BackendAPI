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

	err = createBuyerOtpsTable(db)

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

	err = createPreorderInformationTable(db)

	if err != nil {
		return err
	}

	err = createProductDiscountsTable(db)

	if err != nil {
		return err
	}

	err = createProductImagesTable(db)

	if err != nil {
		return err
	}

	err = createOrdersTable(db)

	if err != nil {
		return err
	}

	err = createGuestOrdersTable(db)

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
		email VARCHAR NOT NULL UNIQUE, 
		password VARCHAR NOT NULL,
		verification VARCHAR NOT NULL DEFAULT 'pending',
		PRIMARY KEY(buyer_id));`

	_, err := db.ExecContext(context.Background(), query)
	return err
}

/*
Create the table for Buyer Otps
*/
func createBuyerOtpsTable(db *sql.DB) error {
	query := `CREATE TABLE IF NOT EXISTS buyer_otps(
		buyer_id uuid REFERENCES buyers(buyer_id) NOT NULL,
		email_otp VARCHAR NOT NULL,
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
		email VARCHAR NOT NULL UNIQUE, 
		password VARCHAR NOT NULL,
		seller_name VARCHAR NOT NULL UNIQUE,
		followers INT DEFAULT 0 NOT NULL,
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
		image_count INT NOT NULL DEFAULT 0,
		condition INT NOT NULL CONSTRAINT isOutOfFive CHECK (condition >= 0 AND condition <= 5),
		price INT NOT NULL CONSTRAINT isPositive CHECK (price >= 0),
		product_type VARCHAR NOT NULL,
		posted_date TIMESTAMPTZ NOT NULL,
		product_quantity INT NOT NULL,
		sold_quantity INT DEFAULT 0 NOT NULL,
		language VARCHAR DEFAULT 'ENG' NOT NULL,
		PRIMARY KEY(product_id));`

	_, err := db.ExecContext(context.Background(), query)
	return err
}

/*
Create Pre Order Information table
*/
func createPreorderInformationTable(db *sql.DB) error {
	query := `CREATE TABLE IF NOT EXISTS preorder_information(
		product_id uuid REFERENCES products(product_id) NOT NULL,
		order_by TIMESTAMPTZ NOT NULL,
		releases_on TIMESTAMPTZ NOT NULL,
		PRIMARY KEY(product_id));`

	_, err := db.ExecContext(context.Background(), query)
	return err
}

/*
Create Product Discounts table
*/
func createProductDiscountsTable(db *sql.DB) error {
	query := `CREATE TABLE IF NOT EXISTS product_discounts(
		product_id uuid REFERENCES products(product_id) NOT NULL,
		discount INT NOT NULL,
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
		image_no INT NOT NULL,
		PRIMARY KEY(product_image_id));`

	_, err := db.ExecContext(context.Background(), query)
	return err
}

/*
Create the table for Orders
*/
func createOrdersTable(db *sql.DB) error {
	query := `CREATE TABLE IF NOT EXISTS orders(
		order_id uuid DEFAULT uuid_generate_v1() NOT NULL,
		product_id uuid REFERENCES products(product_id) NOT NULL,
		buyer_id uuid REFERENCES buyers(buyer_id) NOT NULL,
		delivery_type VARCHAR NOT NULL,
		order_quantity INT NOT NULL, 
		payment_type VARCHAR NOT NULL,
		payment_status VARCHAR DEFAULT 'pending' NOT NULL,
		phone_number VARCHAR NOT NULL,
		order_date TIMESTAMPTZ NOT NULL,
		address_line_1 VARCHAR NOT NULL,
		address_line_2 VARCHAR,
		postal_code VARCHAR NOT NULL,
		telegram_handle VARCHAR,
		product_price INT NOT NULL,
		delivery_fee INT NOT NULL,
		payment_fee INT NOT NULL,
		small_order_fee INT NOT NULL,
		total_paid INT NOT NULL,
		PRIMARY KEY(order_id));`

	_, err := db.ExecContext(context.Background(), query)
	return err
}

/*
Create the table for Guest Orders
*/
func createGuestOrdersTable(db *sql.DB) error {
	query := `CREATE TABLE IF NOT EXISTS guest_orders(
		guest_order_id uuid DEFAULT uuid_generate_v1() NOT NULL,
		product_id uuid REFERENCES products(product_id) NOT NULL,
		delivery_type VARCHAR NOT NULL,
		order_quantity INT NOT NULL, 
		payment_type VARCHAR NOT NULL,
		payment_status VARCHAR DEFAULT 'pending' NOT NULL,
		phone_number VARCHAR NOT NULL,
		email VARCHAR NOT NULL,
		order_date TIMESTAMPTZ NOT NULL,
		address_line_1 VARCHAR NOT NULL,
		address_line_2 VARCHAR,
		postal_code VARCHAR NOT NULL,
		telegram_handle VARCHAR,
		product_price INT NOT NULL,
		delivery_fee INT NOT NULL,
		payment_fee INT NOT NULL,
		small_order_fee INT NOT NULL,
		total_paid INT NOT NULL,
		PRIMARY KEY(guest_order_id));`

	_, err := db.ExecContext(context.Background(), query)
	return err
}
