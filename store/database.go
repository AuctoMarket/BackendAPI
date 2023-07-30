package store

import (
	"BackendAPI/utils"
	"database/sql"
	"fmt"
	"log"
)

/*
Function to setup the DB connections, create all the tables and return
the db connection
*/
func SetupDB() (*sql.DB, error) {
	db, err := initDB(".env", false)

	if err != nil {
		return db, err
	}

	log.Println("Established a successful connection!")

	err = createTables(db)

	if err != nil {
		return db, err
	}

	log.Println("Tables Created Successfully!")

	return db, nil
}

/*
Function to initiate the DB connection and returns the DB connection
*/
func initDB(path string, isTest bool) (*sql.DB, error) {
	var host string
	port, err := utils.GetDotEnvInt("POSTGRES_PORT", path)
	user, err := utils.GetDotEnv("POSTGRES_USER", path)
	password, err := utils.GetDotEnv("POSTGRES_PASSWORD", path)
	dbname, err := utils.GetDotEnv("POSTGRES_DBNAME", path)

	if isTest {
		host, err = utils.GetDotEnv("POSTGRES_HOST_TEST", path)
	} else {
		host, err = utils.GetDotEnv("POSTGRES_HOST", path)
	}

	if err != nil {
		return nil, err
	}

	postgresqlDbInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", postgresqlDbInfo)
	if err != nil {
		return db, err
	}

	err = db.Ping()
	if err != nil {
		return db, err
	}

	log.Println("Established a successful connection!")

	return db, nil
}

/*
Function to initiate the DB connection and returns the DB connection
*/
func initTestDB(path string) (*sql.DB, error) {
	port, err := utils.GetDotEnvInt("POSTGRES_PORT_TEST", path)
	user, err := utils.GetDotEnv("POSTGRES_USER", path)
	password, err := utils.GetDotEnv("POSTGRES_PASSWORD", path)
	dbname, err := utils.GetDotEnv("POSTGRES_DBNAME_TEST", path)
	host, err := utils.GetDotEnv("POSTGRES_HOST_TEST", path)

	if err != nil {
		return nil, err
	}

	postgresqlDbInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", postgresqlDbInfo)
	if err != nil {
		return db, err
	}

	err = db.Ping()
	if err != nil {
		return db, err
	}

	return db, nil
}

/*
Function to reset all the tables in the DB, used mainly during testing
*/
func resetDB(db *sql.DB) {
	queryResetBuyers := `TRUNCATE buyers CASCADE;`
	queryResetSellers := `TRUNCATE sellers CASCADE;`
	queryResetProducts := `TRUNCATE products CASCADE;`

	db.Exec(queryResetBuyers)
	db.Exec(queryResetSellers)
	db.Exec(queryResetProducts)
}

/*
Function to Close the DB connection
*/
func CloseDB(db *sql.DB) {
	db.Close()
}

/*
Function to setup the DB connections for tests, create all the tables and
return the db connection
*/
func SetupTestDB() (*sql.DB, error) {
	db, err := initTestDB("/Users/ekam/Desktop/AuctoCode/BackendAPI/.env")

	if err != nil {
		return db, err
	}

	err = createTables(db)

	if err != nil {
		return db, err
	}

	resetDB(db)

	return db, nil
}
