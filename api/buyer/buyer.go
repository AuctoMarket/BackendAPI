package buyer

import (
	"BackendAPI/data"
	"BackendAPI/utils"
	"context"
	"database/sql"
	"errors"
)

/*
Logic for Buyer login, checks if the user exists in the database and checks if the stored
password matches the plaintext password.
*/
func BuyerLogin(db *sql.DB, loginData data.LoginData) (data.LoginResponseData, error) {
	var response data.LoginResponseData
	var hashedPwd string
	buyerExists, err := doesBuyerEmailExist(db, loginData.Email)

	if err != nil {
		return response, utils.LogError("Error in checking if buyer email exists", err)
	}

	if !buyerExists {
		return response, errors.New("Incorrect user email or password!")
	}

	query := `SELECT email, buid, password from buyers WHERE email = $1;`
	err = db.QueryRowContext(context.Background(), query, loginData.Email).Scan(&response.Email, &response.BUID, &hashedPwd)

	if err != nil {
		return response, utils.LogError("Error in Selecting Buyer rows", err)
	}

	if !utils.ComparePasswords(hashedPwd, loginData.Password) {
		return response, errors.New("Incorrect user email or password")
	}

	return response, nil
}

/*
Logic for Buyer signup, adds the new user to the database if the email is not already in use.
Returned response is similar to as if the user logged in.
*/
func BuyerSignUp(db *sql.DB, signupData data.SignUpData) (data.LoginResponseData, error) {
	var response data.LoginResponseData

	buyerExists, err := doesBuyerEmailExist(db, signupData.Email)
	if err != nil {
		return response, utils.LogError("Error in checking if buyer email exists", err)
	}

	if buyerExists {
		return response, errors.New("Email is already in use")
	}

	hashPassword, err := utils.HashAndSalt([]byte(signupData.Password))

	if err != nil {
		return response, utils.LogError("Error in hash function!", err)
	}

	query := `INSERT INTO buyers(email, password) VALUES ($1,$2) RETURNING email, buid;`
	err = db.QueryRowContext(context.Background(), query, signupData.Email, hashPassword).Scan(&response.Email, &response.BUID)

	if err != nil {
		return response, utils.LogError("Error in Inserting Rows into Buyers table", err)
	}
	return response, nil
}

/*
Checks wether a Buyer with a given email address already exists in the database
and returns true if it does false otherwise.
*/
func doesBuyerEmailExist(db *sql.DB, email string) (bool, error) {
	var buyerExists bool
	query := `SELECT EXISTS(SELECT * FROM buyers WHERE email = $1);`
	err := db.QueryRowContext(context.Background(), query, email).Scan(&buyerExists)

	if err != nil {
		return false, err
	}

	return buyerExists, nil
}
