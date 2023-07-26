package store

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
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
	port, err := getDotEnvInt("POSTGRES_PORT", path)
	user, err := getDotEnv("POSTGRES_USER", path)
	password, err := getDotEnv("POSTGRES_PASSWORD", path)
	dbname, err := getDotEnv("POSTGRES_DBNAME", path)
	if isTest {
		host, err = getDotEnv("POSTGRES_HOST_TEST", path)
	} else {
		host, err = getDotEnv("POSTGRES_HOST", path)
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

	err = createTables(db)

	if err != nil {
		return db, err
	}

	return db, nil
}

/*
Function to reset all the tables in the DB, used mainly during testing
*/
func ResetDB(db *sql.DB) {
	queryResetBuyers := `TRUNCATE buyers;`
	queryResetSellers := `TRUNCATE sellers;`

	db.Exec(queryResetBuyers)
	db.Exec(queryResetSellers)
}

/*
Function to Close the DB connection
*/
func CloseDB(db *sql.DB) {
	db.Close()
}

/*
Function to get environment variables from the .Env file if they are a
string
*/
func getDotEnv(key string, path string) (string, error) {
	err := godotenv.Load(path)

	if err != nil {
		return "", err
	}

	return os.Getenv(key), nil
}

/*
Function to get environment variables from the .Env file if they are a
int
*/
func getDotEnvInt(key string, path string) (int, error) {
	err := godotenv.Load(path)

	if err != nil {
		return 0, err
	}

	num, err := strconv.ParseInt(os.Getenv(key), 10, 32)

	if err != nil {
		return 0, err
	}

	return int(num), nil
}

/*
Function to setup the DB connections for tests, create all the tables and
return the db connection
*/
func SetupTestDB() (*sql.DB, error) {
	db, err := initDB("/Users/ekam/Desktop/AuctoCode/BackendAPI/.env", true)

	if err != nil {
		return db, err
	}

	log.Println("Established a successful connection!")

	err = createTables(db)

	if err != nil {
		return db, err
	}

	log.Println("Tables Created Successfully!")
	ResetDB(db)

	return db, nil
}

/*
Cleans up the test DB and clears all test data
*/
func CleaupTestDB(db *sql.DB) {
	ResetDB(db)
	CloseDB(db)
}
