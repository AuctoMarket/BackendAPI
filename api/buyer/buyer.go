package buyer

import (
	"BackendAPI/data"
	"BackendAPI/utils"
	"context"
	"database/sql"
)

/*
Logic for Buyer login, checks if the user exists in the database and checks if the stored
password matches the plaintext password.
*/
func BuyerLogin(db *sql.DB, loginData data.BuyerLoginData) (data.BuyerLoginResponseData, *utils.ErrorHandler) {
	var response data.BuyerLoginResponseData
	var hashedPwd string
	buyerExists := doesBuyerEmailExist(db, loginData.Email)

	if !buyerExists {
		return response, utils.UnauthorizedError("Incorrect user email or password!")
	}

	query := `SELECT email, buyer_id, password from buyers WHERE email = $1;`
	err := db.QueryRowContext(context.Background(), query, loginData.Email).Scan(&response.Email, &response.BuyerId, &hashedPwd)

	if err != nil {
		errResp := utils.InternalServerError()
		utils.LogError(err, "Error in Selecting Buyer rows")
		return response, errResp
	}

	if !utils.ComparePasswords(hashedPwd, loginData.Password) {
		return response, utils.UnauthorizedError("Incorrect user email or password!")
	}

	return response, nil
}

/*
Logic for Buyer signup, adds the new user to the database if the email is not already in use.
Returned response is similar to as if the user logged in.
*/
func BuyerSignUp(db *sql.DB, signupData data.BuyerSignUpData) (data.BuyerLoginResponseData, *utils.ErrorHandler) {
	var response data.BuyerLoginResponseData

	buyerExists := doesBuyerEmailExist(db, signupData.Email)

	if buyerExists {
		return response, utils.BadRequestError("Email is already in use")
	}

	hashPassword, err := utils.HashAndSalt([]byte(signupData.Password))

	if err != nil {
		errResp := utils.InternalServerError()
		utils.LogError(err, "Error in hash function!")
		return response, errResp
	}

	query := `INSERT INTO buyers(email, password) VALUES ($1,$2) RETURNING email, buyer_id;`
	err = db.QueryRowContext(context.Background(), query, signupData.Email, hashPassword).Scan(&response.Email, &response.BuyerId)

	if err != nil {
		errResp := utils.InternalServerError()
		utils.LogError(err, "Error in Inserting Rows into Buyers table")
		return response, errResp
	}

	return response, nil
}

/*
Checks wether a Buyer with a given email address already exists in the database
and returns true if it does false otherwise.
*/
func doesBuyerEmailExist(db *sql.DB, email string) bool {
	var buyerExists bool
	query := `SELECT EXISTS(SELECT * FROM buyers WHERE email = $1);`
	err := db.QueryRowContext(context.Background(), query, email).Scan(&buyerExists)

	if err != nil {
		return false
	}

	return buyerExists
}
