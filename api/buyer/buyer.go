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
func BuyerLogin(db *sql.DB, loginData data.UserLoginData) (data.BuyerLoginResponseData, *utils.ErrorHandler) {
	var response data.BuyerLoginResponseData
	var hashedPwd string
	buyerExists := doesBuyerEmailExist(db, loginData.Email)

	if !buyerExists {
		return response, utils.UnauthorizedError("Incorrect user email or password!")
	}

	query := `SELECT email, buyer_id, verification, password from buyers WHERE email = $1;`
	err := db.QueryRowContext(context.Background(), query, loginData.Email).Scan(
		&response.Email, &response.BuyerId, &response.Verification, &hashedPwd)

	if err != nil {
		errResp := utils.InternalServerError(err)
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
		errResp := utils.InternalServerError(err)
		utils.LogError(err, "Error in hash function!")
		return response, errResp
	}

	otp := utils.GetOtp(6)

	query := `INSERT INTO buyers(email, password, email_otp) VALUES ($1,$2,$3) RETURNING email, buyer_id, verification;`
	err = db.QueryRowContext(context.Background(), query, signupData.Email, hashPassword, otp).Scan(
		&response.Email, &response.BuyerId, &response.Verification)

	if err != nil {
		errResp := utils.InternalServerError(err)
		utils.LogError(err, "Error in Inserting Rows into Buyers table")
		return response, errResp
	}

	err = utils.SendOtpMail(response.Email, otp)

	return response, nil
}

func ResendOtp(db *sql.DB, resendOtpReq data.BuyerResendOtpData) (data.Message, *utils.ErrorHandler) {
	var response data.Message

	if !DoesBuyerExist(db, resendOtpReq.BuyerId) {
		return response, utils.BadRequestError("The buyer_id provided is invalid")
	}

	newOtp := utils.GetOtp(6)
	query := `UPDATE buyers SET email_otp = $1 WHERE buyer_id = $2 RETURNING email`
	var email string
	err := db.QueryRowContext(context.Background(), query, newOtp, resendOtpReq.BuyerId).Scan(&email)

	if err != nil {
		errResp := utils.InternalServerError(nil)
		utils.LogError(err, "Error in Updating buyer rows")
		return response, errResp
	}

	utils.SendOtpMail(email, newOtp)

	response.Message = "Resent otp to provided email address"
	return response, nil
}

func ValidateOtp(db *sql.DB, validateOtpReq data.BuyerValidateOtpData) (data.BuyerLoginResponseData, *utils.ErrorHandler) {
	var response data.BuyerLoginResponseData

	if !DoesBuyerExist(db, validateOtpReq.BuyerId) {
		return response, utils.BadRequestError("The buyer_id provided is invalid")
	}

	query := `SELECT email_otp, email FROM buyers WHERE buyer_id = $1`
	var otp string
	err := db.QueryRowContext(context.Background(), query, validateOtpReq.BuyerId).Scan(&otp, &response.Email)

	if err != nil {
		errResp := utils.InternalServerError(nil)
		utils.LogError(err, "Error in Selecting buyer rows")
		return response, errResp
	}

	if otp != validateOtpReq.Otp {
		return response, utils.UnauthorizedError("Incorrect Otp!")
	}

	query = `UPDATE buyers SET verification = 'verified' WHERE buyer_id = $1`
	_, err = db.ExecContext(context.Background(), query, validateOtpReq.BuyerId)

	response.BuyerId = validateOtpReq.BuyerId
	response.Verification = "verified"
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

/*
Checks wether a Buyer with a given email address already exists in the database
and returns true if it does false otherwise.
*/
func DoesBuyerExist(db *sql.DB, id string) bool {
	var buyerExists bool
	query := `SELECT EXISTS(SELECT * FROM buyers WHERE buyer_id = $1);`
	err := db.QueryRowContext(context.Background(), query, id).Scan(&buyerExists)

	if err != nil {
		return false
	}

	return buyerExists
}
