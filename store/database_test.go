package store

import (
	"testing"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestSetupTestDB(t *testing.T) {

	//Test 1: Setup a DB connection and check if no errors
	db, err := SetupTestDB("../.env")
	assert.NoError(t, err)
	assert.NotEmpty(t, db)

	CloseDB(db)

	//Test 2: Setup a DB connection with invalid env path
	db, err = SetupTestDB(".env")
	assert.Error(t, err)

}

func TestInitDB(t *testing.T) {
	//Test 1: Init DB and ensure the db exists and there is no error for main db
	db, err := initDB("../.env", true)
	assert.NotEmpty(t, db)
	assert.NoError(t, err)

	//Test 2: Ping the db and ensure the connection exists
	err = db.Ping()
	assert.NoError(t, err)

	CloseDB(db)

	//Test 2: Setup a DB connection with invalid env path
	db, err = initDB(".env", true)
	assert.Error(t, err)

}
