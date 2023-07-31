package store

import (
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
)

func TestCreateTables(t *testing.T) {
	db, err := initTestDB("/Users/ekam/Desktop/AuctoCode/BackendAPI/.env")
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
	var buyersExist bool
	db, err := initTestDB("/Users/ekam/Desktop/AuctoCode/BackendAPI/.env")
	assert.NoError(t, err)

	dropDB(db)
	//Test 1: No errors in creating buyers table
	err = db.QueryRowContext(context.Background(), queryCheckTableBuyers).Scan(&buyersExist)
	assert.Equal(t, false, buyersExist)
	err = createBuyersTable(db)
	assert.NoError(t, err)

	//Test 2: Check if neccessary buyers tables exists
	err = db.QueryRowContext(context.Background(), queryCheckTableBuyers).Scan(&buyersExist)
	assert.NoError(t, err)
	assert.Equal(t, true, buyersExist)

	CloseDB(db)

}

func TestCreateSellersTable(t *testing.T) {
	db, err := initTestDB("/Users/ekam/Desktop/AuctoCode/BackendAPI/.env")
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
	db, err := initTestDB("/Users/ekam/Desktop/AuctoCode/BackendAPI/.env")
	assert.NoError(t, err)

	dropDB(db)

	//Test 2: No Error in creating products table
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
