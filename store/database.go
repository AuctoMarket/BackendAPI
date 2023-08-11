package store

import (
	"BackendAPI/utils"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
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
	var (
		env      string
		host     string
		port     int
		user     string
		password string
		dbname   string
		sslmode  string
		err      error
	)

	env = os.Getenv("DB_ENV")

	if env == "rds" {
		port, err = utils.GetDotEnvInt("POSTGRES_PORT_RDS")
		user = os.Getenv("POSTGRES_USER_RDS")
		password = os.Getenv("POSTGRES_PASSWORD_RDS")
		dbname = os.Getenv("POSTGRES_DBNAME_RDS")
		host = os.Getenv("POSTGRES_HOST_RDS")
		sslmode = "require"
	} else {
		port, err = utils.GetDotEnvInt("POSTGRES_PORT_LOCAL")
		user = os.Getenv("POSTGRES_USER_LOCAL")
		password = os.Getenv("POSTGRES_PASSWORD_LOCAL")
		dbname = os.Getenv("POSTGRES_DBNAME_LOCAL")
		host = os.Getenv("POSTGRES_HOST_LOCAL")
		sslmode = "disable"

		if isTest {
			host = os.Getenv("POSTGRES_HOST_TEST")
		}
	}

	if err != nil {
		return nil, errors.New("The .env environments could not be loaded:" + err.Error())
	}

	postgresqlDbInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode)

	db, err := sql.Open("postgres", postgresqlDbInfo)
	if err != nil {
		return db, errors.New("The database connection could not be reated:" + err.Error())
	}

	err = db.Ping()
	if err != nil {
		return db, errors.New("Could not ping the database:" + err.Error())
	}

	log.Println("Established a successful connection!")

	return db, nil
}

/*
Function to initiate the DB connection and returns the DB connection
*/
func initTestDB() (*sql.DB, error) {
	port, err := utils.GetDotEnvInt("POSTGRES_PORT_TEST")
	user := os.Getenv("POSTGRES_USER_TEST")
	password := os.Getenv("POSTGRES_PASSWORD_TEST")
	dbname := os.Getenv("POSTGRES_DBNAME_TEST")
	host := os.Getenv("POSTGRES_HOST_TEST")

	if err != nil {
		return nil, errors.New("The .env environments could not be loaded:" + err.Error())
	}

	postgresqlDbInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", postgresqlDbInfo)
	if err != nil {
		return db, errors.New("The database connection could not be reated:" + err.Error())
	}

	err = db.Ping()
	if err != nil {
		return db, errors.New("Could not ping the database:" + err.Error())
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
	queryResetProductImages := `TRUNCATE product_images CASCADE;`

	db.Exec(queryResetBuyers)
	db.Exec(queryResetSellers)
	db.Exec(queryResetProducts)
	db.Exec(queryResetProductImages)
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
func SetupTestDB(path string) (*sql.DB, error) {

	loadErr := utils.LoadDotEnv(path)

	if loadErr != nil {
		utils.LogError(loadErr, "Cannot fetch .env, no .env file")
		return nil, loadErr
	}

	db, err := initTestDB()

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
