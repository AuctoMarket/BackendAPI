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
	var (
		host     string
		port     int
		user     string
		password string
		dbname   string
		sslmode  string
		err      error
	)
	env, err := utils.GetDotEnv("DB_ENV", path)
	if env == "rds" {
		port, err = utils.GetDotEnvInt("POSTGRES_PORT_RDS", path)
		user, err = utils.GetDotEnv("POSTGRES_USER_RDS", path)
		password, err = utils.GetDotEnv("POSTGRES_PASSWORD_RDS", path)
		dbname, err = utils.GetDotEnv("POSTGRES_DBNAME_RDS", path)
		host, err = utils.GetDotEnv("POSTGRES_HOST_RDS", path)
		sslmode = "require"
	} else {
		port, err = utils.GetDotEnvInt("POSTGRES_PORT_LOCAL", path)
		user, err = utils.GetDotEnv("POSTGRES_USER_LOCAL", path)
		password, err = utils.GetDotEnv("POSTGRES_PASSWORD_LOCAL", path)
		dbname, err = utils.GetDotEnv("POSTGRES_DBNAME_LOCAL", path)
		sslmode = "disable"

		if isTest {
			host, err = utils.GetDotEnv("POSTGRES_HOST_TEST", path)
		} else {
			host, err = utils.GetDotEnv("POSTGRES_HOST_LOCAL", path)
		}
	}

	if err != nil {
		return nil, err
	}

	postgresqlDbInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode)

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
	user, err := utils.GetDotEnv("POSTGRES_USER_TEST", path)
	password, err := utils.GetDotEnv("POSTGRES_PASSWORD_TEST", path)
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
func SetupTestDB(path string) (*sql.DB, error) {
	db, err := initTestDB(path)

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
