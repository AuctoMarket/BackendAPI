package store

import (
	"testing"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestSetupTestDB(t *testing.T) {

	//Test 1: Setup a DB connection and check if no errors
	db, err := SetupTestDB()
	assert.NoError(t, err)
	assert.NotEmpty(t, db)

	CleaupTestDB(db)
}

func TestInitDB(t *testing.T) {
	//Test 1: Init DB and ensure the db exists and there ius no error
	db, err := initDB("/Users/ekam/Desktop/AuctoCode/BackendAPI/.env")
	assert.NotEmpty(t, db)
	assert.NoError(t, err)

	//Test 2: Ping the db and ensure the connection exists
	err = db.Ping()
	assert.NoError(t, err)

	CleaupTestDB(db)
}
