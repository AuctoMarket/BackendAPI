package store

import (
	"context"
	"testing"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestSetupTestDB(t *testing.T) {
	queryCheckTableBuyers := `SELECT EXISTS(
		SELECT * 
		FROM information_schema.tables 
		WHERE 
		  table_schema = 'public' AND 
		  table_name = 'buyers'
	);`
	queryCheckTableSellers := `SELECT EXISTS(
		SELECT * 
		FROM information_schema.tables 
		WHERE 
		  table_schema = 'public' AND 
		  table_name = 'sellers'
	);`

	queryCheckTableImaginary := `SELECT EXISTS(
		SELECT * 
		FROM information_schema.tables 
		WHERE 
		  table_schema = 'public' AND 
		  table_name = 'imaginary'
	);`

	//Test 1: Setup a DB connection
	db, err := SetupTestDB()
	assert.NoError(t, err)
	assert.NotEmpty(t, db)

	//Test 2: Check if neccessary buyers tables exists
	var buyersExist bool
	err = db.QueryRowContext(context.Background(), queryCheckTableBuyers).Scan(&buyersExist)
	assert.NoError(t, err)
	assert.Equal(t, buyersExist, true)

	//Test 3: Check if neccessary sellers tables exists
	var sellersExist bool
	err = db.QueryRowContext(context.Background(), queryCheckTableSellers).Scan(&sellersExist)
	assert.NoError(t, err)
	assert.Equal(t, sellersExist, true)

	//Test 3: Check no uneccessary tables exist
	var imaginaryExist bool
	db.QueryRowContext(context.Background(), queryCheckTableImaginary).Scan(&imaginaryExist)
	assert.NoError(t, err)

	assert.Equal(t, imaginaryExist, false)

	CleaupTestDB(db)
}
