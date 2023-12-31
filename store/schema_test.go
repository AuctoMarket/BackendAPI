package store

import (
	"BackendAPI/utils"
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	queryCheckTableBuyers = `SELECT EXISTS(
		SELECT * 
		FROM information_schema.tables 
		WHERE 
		  table_schema = 'public' AND 
		  table_name = 'buyers'
	);`

	queryCheckTableBuyerOtps = `SELECT EXISTS(
		SELECT * 
		FROM information_schema.tables 
		WHERE 
		  table_schema = 'public' AND 
		  table_name = 'buyer_otps'
	);`
	queryCheckTableSellers = `SELECT EXISTS(
		SELECT * 
		FROM information_schema.tables 
		WHERE 
		  table_schema = 'public' AND 
		  table_name = 'sellers'
	);`

	queryCheckTableProducts = `SELECT EXISTS(
		SELECT * 
		FROM information_schema.tables 
		WHERE 
		  table_schema = 'public' AND 
		  table_name = 'products'
	);`

	queryCheckTableProductImages = `SELECT EXISTS(
		SELECT * 
		FROM information_schema.tables 
		WHERE 
		  table_schema = 'public' AND 
		  table_name = 'product_images'
	);`

	queryCheckTablePreorderInformation = `SELECT EXISTS(
		SELECT * 
		FROM information_schema.tables 
		WHERE 
		  table_schema = 'public' AND 
		  table_name = 'preorder_information'
	);`

	queryCheckTableProductDiscounts = `SELECT EXISTS(
		SELECT * 
		FROM information_schema.tables 
		WHERE 
		  table_schema = 'public' AND 
		  table_name = 'product_discounts'
	);`

	queryCheckTableOrders = `SELECT EXISTS(
		SELECT * 
		FROM information_schema.tables 
		WHERE 
		  table_schema = 'public' AND 
		  table_name = 'orders'
	);`
	queryCheckTableGuestOrders = `SELECT EXISTS(
		SELECT * 
		FROM information_schema.tables 
		WHERE 
		  table_schema = 'public' AND 
		  table_name = 'guest_orders'
	);`

	queryCheckTableOrderProducts = `SELECT EXISTS(
		SELECT * 
		FROM information_schema.tables 
		WHERE 
		  table_schema = 'public' AND 
		  table_name = 'order_products'
	);`

	queryCheckTableGuestOrderProducts = `SELECT EXISTS(
		SELECT * 
		FROM information_schema.tables 
		WHERE 
		  table_schema = 'public' AND 
		  table_name = 'guest_order_products'
	);`
)

func TestCreateTables(t *testing.T) {
	err := utils.LoadDotEnv("../.env")
	assert.NoError(t, err)
	db, err := initTestDB()
	assert.NoError(t, err)

	dropDB(db)

	//Test 1: No errors in creating tables
	err = createTables(db)
	assert.NoError(t, err)

	//Test 2: Check if neccessary buyers tables exists
	var buyersExist bool
	err = db.QueryRowContext(context.Background(), queryCheckTableBuyers).Scan(&buyersExist)
	assert.NoError(t, err)
	assert.Equal(t, true, buyersExist)

	//Test 3: Check if neccessary sellers tables exists
	var sellersExist bool
	err = db.QueryRowContext(context.Background(), queryCheckTableSellers).Scan(&sellersExist)
	assert.NoError(t, err)
	assert.Equal(t, true, sellersExist)

	//Test 4: Check if neccessary sellers tables exists
	var productsExist bool
	err = db.QueryRowContext(context.Background(), queryCheckTableProducts).Scan(&productsExist)
	assert.NoError(t, err)
	assert.Equal(t, true, productsExist)

	CloseDB(db)
}

func TestCreateBuyersTable(t *testing.T) {
	err := utils.LoadDotEnv("../.env")
	assert.NoError(t, err)
	db, err := initTestDB()
	assert.NoError(t, err)

	dropDB(db)
	//Test 1: No errors in creating buyers table
	var buyersExist bool
	err = createBuyersTable(db)
	assert.NoError(t, err)

	//Test 2: Check if neccessary buyers tables exists
	err = db.QueryRowContext(context.Background(), queryCheckTableBuyers).Scan(&buyersExist)
	assert.NoError(t, err)
	assert.Equal(t, true, buyersExist)

	CloseDB(db)

}

func TestCreateBuyerOtpsTable(t *testing.T) {
	err := utils.LoadDotEnv("../.env")
	assert.NoError(t, err)
	db, err := initTestDB()
	assert.NoError(t, err)

	dropDB(db)
	//Test 1: No errors in creating buyers table
	var buyerOtpsExist bool
	err = createBuyerOtpsTable(db)
	assert.NoError(t, err)

	//Test 2: Check if neccessary buyers tables exists
	err = db.QueryRowContext(context.Background(), queryCheckTableBuyerOtps).Scan(&buyerOtpsExist)
	assert.NoError(t, err)
	assert.Equal(t, true, buyerOtpsExist)

	CloseDB(db)
}

func TestCreateSellersTable(t *testing.T) {
	err := utils.LoadDotEnv("../.env")
	assert.NoError(t, err)
	db, err := initTestDB()
	assert.NoError(t, err)

	//Test 1: No errors in creating sellers table
	err = createSellersTable(db)
	assert.NoError(t, err)

	//Test 2: Check if neccessary buyers tables exists
	var sellersExist bool
	err = db.QueryRowContext(context.Background(), queryCheckTableSellers).Scan(&sellersExist)
	assert.NoError(t, err)
	assert.Equal(t, true, sellersExist)

	CloseDB(db)
}

func TestCreateProductsTable(t *testing.T) {
	err := utils.LoadDotEnv("../.env")
	assert.NoError(t, err)
	db, err := initTestDB()
	assert.NoError(t, err)

	dropDB(db)

	//Test 1: Error in creating products table
	err = createProductsTable(db)
	assert.Error(t, err)

	//Test 2: No Error in creating products table
	createSellersTable(db)
	err = createProductsTable(db)
	assert.NoError(t, err)

	//Test 3: Check if neccessary products tables exists
	var productsExist bool
	err = db.QueryRowContext(context.Background(), queryCheckTableProducts).Scan(&productsExist)
	assert.NoError(t, err)
	assert.Equal(t, true, productsExist)

	CloseDB(db)
}

func TestCreateProductImagesTable(t *testing.T) {
	err := utils.LoadDotEnv("../.env")
	assert.NoError(t, err)
	db, err := initTestDB()
	assert.NoError(t, err)

	dropDB(db)

	//Test 1: No Error in creating products table
	createSellersTable(db)
	createProductsTable(db)
	err = createProductImagesTable(db)
	assert.NoError(t, err)

	//Test 2: Check if neccessary products tables exists
	var productImagesExist bool
	err = db.QueryRowContext(context.Background(), queryCheckTableProductImages).Scan(&productImagesExist)
	assert.NoError(t, err)
	assert.Equal(t, true, productImagesExist)

	CloseDB(db)
}

func TestCreatePreorderInformationTable(t *testing.T) {
	err := utils.LoadDotEnv("../.env")
	assert.NoError(t, err)
	db, err := initTestDB()
	assert.NoError(t, err)

	dropDB(db)

	//Test 1: No Error in creating products table
	createSellersTable(db)
	createProductsTable(db)
	err = createPreorderInformationTable(db)
	assert.NoError(t, err)

	//Test 2: Check if neccessary products tables exists
	var preorderInformationExists bool
	err = db.QueryRowContext(context.Background(), queryCheckTablePreorderInformation).Scan(&preorderInformationExists)
	assert.NoError(t, err)
	assert.Equal(t, true, preorderInformationExists)

	CloseDB(db)
}

func TestCreateProductDiscountsTable(t *testing.T) {
	err := utils.LoadDotEnv("../.env")
	assert.NoError(t, err)
	db, err := initTestDB()
	assert.NoError(t, err)

	dropDB(db)

	//Test 1: No Error in creating products table
	createSellersTable(db)
	createProductsTable(db)
	err = createProductDiscountsTable(db)
	assert.NoError(t, err)

	//Test 2: Check if neccessary products tables exists
	var productDiscountsExists bool
	err = db.QueryRowContext(context.Background(), queryCheckTableProductDiscounts).Scan(&productDiscountsExists)
	assert.NoError(t, err)
	assert.Equal(t, true, productDiscountsExists)

	CloseDB(db)
}

func TestCreateOrdersTable(t *testing.T) {
	err := utils.LoadDotEnv("../.env")
	assert.NoError(t, err)
	db, err := initTestDB()
	assert.NoError(t, err)

	dropDB(db)

	//Test 1: No Error in creating orders table
	createSellersTable(db)
	createProductsTable(db)
	createBuyersTable(db)
	err = createOrdersTable(db)
	assert.NoError(t, err)

	//Test 2: Check if neccessary orders tables exists
	var ordersExist bool
	err = db.QueryRowContext(context.Background(), queryCheckTableOrders).Scan(&ordersExist)
	assert.NoError(t, err)
	assert.Equal(t, true, ordersExist)

	CloseDB(db)
}

func TestCreateGuestOrdersTable(t *testing.T) {
	err := utils.LoadDotEnv("../.env")
	assert.NoError(t, err)
	db, err := initTestDB()
	assert.NoError(t, err)

	dropDB(db)

	//Test 1: No Error in creating guest orders table
	createSellersTable(db)
	createProductsTable(db)
	err = createGuestOrdersTable(db)
	assert.NoError(t, err)

	//Test 2: Check if neccessary guest orders tables exists
	var guestOrdersExist bool
	err = db.QueryRowContext(context.Background(), queryCheckTableGuestOrders).Scan(&guestOrdersExist)
	assert.NoError(t, err)
	assert.Equal(t, true, guestOrdersExist)

	CloseDB(db)
}

func TestCreateOrderProductsTable(t *testing.T) {
	err := utils.LoadDotEnv("../.env")
	assert.NoError(t, err)
	db, err := initTestDB()
	assert.NoError(t, err)

	dropDB(db)

	//Test 1: No Error in creating order products table
	createSellersTable(db)
	createProductsTable(db)
	createOrdersTable(db)
	err = createOrderProductsTable(db)
	assert.NoError(t, err)

	//Test 2: Check if neccessary order products tables exists
	var orderProductsExist bool
	err = db.QueryRowContext(context.Background(), queryCheckTableOrderProducts).Scan(&orderProductsExist)
	assert.NoError(t, err)
	assert.Equal(t, true, orderProductsExist)

	CloseDB(db)
}

func TestCreateGuestOrderProductsTable(t *testing.T) {
	err := utils.LoadDotEnv("../.env")
	assert.NoError(t, err)
	db, err := initTestDB()
	assert.NoError(t, err)

	dropDB(db)

	//Test 1: No Error in creating guest order products table
	createSellersTable(db)
	createProductsTable(db)
	createOrdersTable(db)
	err = createGuestOrderProductsTable(db)
	assert.NoError(t, err)

	//Test 2: Check if neccessary guest order products tables exists
	var guestOrderProductsExist bool
	err = db.QueryRowContext(context.Background(), queryCheckTableGuestOrderProducts).Scan(&guestOrderProductsExist)
	assert.NoError(t, err)
	assert.Equal(t, true, guestOrderProductsExist)

	CloseDB(db)
}

/*
Function to reset all the tables in the DB, used mainly during testing
*/
func dropDB(db *sql.DB) {
	queryDropBuyers := `DROP TABLE buyers CASCADE;`
	queryDropSellers := `DROP TABLE sellers CASCADE;`
	queryDropProducts := `DROP TABLE products CASCADE;`

	db.Exec(queryDropBuyers)
	db.Exec(queryDropSellers)
	db.Exec(queryDropProducts)
}
