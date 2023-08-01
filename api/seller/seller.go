package seller

import (
	"BackendAPI/data"
	"BackendAPI/utils"
	"context"
	"database/sql"
)

/*
Logic for Seller login, checks if the user exists in the database and checks if the stored
password matches the plaintext password.
*/
func SellerLogin(db *sql.DB, loginData data.UserLoginData) (data.SellerResponseData, *utils.ErrorHandler) {
	var response data.SellerResponseData
	var hashedPwd string
	sellerExists := doesSellerEmailExist(db, loginData.Email)

	if !sellerExists {
		return response, utils.UnauthorizedError("Incorrect user email or password!")
	}

	query := `SELECT email, seller_id, seller_name, password from sellers WHERE email = $1;`
	err := db.QueryRowContext(context.Background(), query, loginData.Email).Scan(
		&response.Email, &response.SellerId, &response.SellerName, &hashedPwd)

	if err != nil {
		errResp := utils.InternalServerError()
		utils.LogError(err, "Error in Selecting Seller rows")
		return response, errResp
	}

	if !utils.ComparePasswords(hashedPwd, loginData.Password) {
		return response, utils.UnauthorizedError("Incorrect user email or password!")
	}

	return response, nil
}

/*
Logic for Seller signup, adds the new seller to the database if the email is not already in use.
Returned response is similar to as if the seller logged in.
*/
func SellerSignUp(db *sql.DB, signupData data.SellerSignUpData) (data.SellerResponseData, *utils.ErrorHandler) {
	var response data.SellerResponseData

	sellerEmailExists := doesSellerEmailExist(db, signupData.Email)

	if sellerEmailExists {
		return response, utils.BadRequestError("Email is already in use")
	}

	sellerNameExists := doesSellerNameExist(db, signupData.SellerName)

	if sellerNameExists {
		return response, utils.BadRequestError("Seller name is already in use")
	}

	hashPassword, err := utils.HashAndSalt([]byte(signupData.Password))

	if err != nil {
		errResp := utils.InternalServerError()
		utils.LogError(err, "Error in hash function!")
		return response, errResp
	}

	query := `INSERT INTO sellers(email, password, seller_name) VALUES ($1,$2,$3) RETURNING email, seller_id, seller_name;`
	err = db.QueryRowContext(context.Background(), query,
		signupData.Email, hashPassword, signupData.SellerName).Scan(
		&response.Email, &response.SellerId, &response.SellerName)

	if err != nil {
		errResp := utils.InternalServerError()
		utils.LogError(err, "Error in Inserting Rows into Sellers table")
		return response, errResp
	}

	return response, nil
}

/*
Checks wether a seller with a given email address already exists in the database
and returns true if it does false otherwise.
*/
func doesSellerEmailExist(db *sql.DB, email string) bool {
	var sellersExists bool
	query := `SELECT EXISTS(SELECT * FROM sellers WHERE email = $1);`
	err := db.QueryRowContext(context.Background(), query, email).Scan(&sellersExists)

	if err != nil {
		return false
	}

	return sellersExists
}

/*
Checks wether a seller with a given name already exists in the database
and returns true if it does false otherwise.
*/
func doesSellerNameExist(db *sql.DB, sellerName string) bool {
	var sellersExists bool
	query := `SELECT EXISTS(SELECT * FROM sellers WHERE seller_name = $1);`
	err := db.QueryRowContext(context.Background(), query, sellerName).Scan(&sellersExists)

	if err != nil {
		return false
	}

	return sellersExists
}

/*
Checks wether a Buyer with a given seller_id already exists in the database
and returns true if it does false otherwise.
*/
func DoesSellerExist(db *sql.DB, id string) bool {
	var sellersExists bool
	query := `SELECT EXISTS(SELECT * FROM sellers WHERE seller_id = $1);`
	err := db.QueryRowContext(context.Background(), query, id).Scan(&sellersExists)

	if err != nil {
		return false
	}

	return sellersExists
}
