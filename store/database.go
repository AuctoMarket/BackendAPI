package store

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

func SetupDB() (*sql.DB, error) {

	host, err := getDotEnv("POSTGRES_HOST")
	port, err := getDotEnvInt("POSTGRES_PORT")
	user, err := getDotEnv("POSTGRES_USER")
	password, err := getDotEnv("POSTGRES_PASSWORD")
	dbname, err := getDotEnv("POSTGRES_DBNAME")

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

	fmt.Println("Established a successful connection!")

	err = CreateTables(db)

	if err != nil {
		return db, err
	}

	fmt.Println("Tables Created Successfully!")

	return db, nil
}

func CloseDB(db *sql.DB) {
	db.Close()
}

func getDotEnv(key string) (string, error) {
	err := godotenv.Load(".env")

	if err != nil {
		return "", err
	}

	return os.Getenv(key), nil
}

func getDotEnvInt(key string) (int, error) {
	err := godotenv.Load(".env")

	if err != nil {
		return 0, err
	}

	num, err := strconv.ParseInt(os.Getenv(key), 10, 32)

	if err != nil {
		return 0, err
	}

	return int(num), nil
}
