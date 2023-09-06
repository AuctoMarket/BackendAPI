package main

import (
	"BackendAPI/api/buyer"
	"BackendAPI/data"
	"net/http"

	"github.com/gin-gonic/gin"
)

// handleBuyerSignup godoc
// @Summary      Signs a new buyer up
// @Description  Checks to see if a buyer email exists and if not creates a new account with supplied email and password
// if not returns a bad request error (400).
// @Accept       json
// @Produce      json
// @Param 		 email body string true "Buyers email [UNIQUE]"
// @Param 		 password body string true "Buyers password as plaintext"
// @Success      200  {object}  data.BuyerLoginResponseData
// @Failure      400  {object}  data.Message
// @Failure      500  {object}  data.Message
// @Router       /buyers/signup [post]
func handleBuyerSignUp(c *gin.Context) {
	var signUpData data.BuyerSignUpData
	bindErr := c.ShouldBindJSON(&signUpData)

	if bindErr != nil {
		r := data.Message{Message: "Bad Request Body"}
		c.JSON(http.StatusBadRequest, r)
		return
	}

	response, err := buyer.BuyerSignUp(db, signUpData)

	if err != nil {
		r := data.Message{Message: err.Error()}
		c.JSON(err.ErrorCode(), r)
		return
	}

	c.JSON(http.StatusCreated, &response)
}

// handleBuyerLogin godoc
// @Summary      Logs a buyer into their account
// @Description  Checks to see if a buyer email exists and if supplied password matches the stored password
// if not returns a unauthorized error (401).
// @Accept       json
// @Produce      json
// @Param 		 email body string true "Buyers email"
// @Param 		 password body string true "Buyers password as plaintext"
// @Success      200  {object}  data.BuyerLoginResponseData
// @Failure      400  {object}  data.Message
// @Failure      401  {object}  data.Message
// @Failure      500  {object}  data.Message
// @Router       /buyers/login [post]
func handleBuyerLogin(c *gin.Context) {
	var loginData data.UserLoginData
	bindErr := c.ShouldBindJSON(&loginData)

	if bindErr != nil {
		r := data.Message{Message: "Bad Request Body"}
		c.JSON(http.StatusBadRequest, r)
		return
	}

	response, err := buyer.BuyerLogin(db, loginData)

	if err != nil {
		r := data.Message{Message: err.Error()}
		c.JSON(err.ErrorCode(), r)
		return
	}

	c.JSON(http.StatusOK, &response)
}

// handleResendOtp godoc
// @Summary      Sends a new Otp to the provided email
// @Description  Checks to see if the provided buyer_id exists and sends a email to the specific buy_ids email with a newly
// generated Otp. If buyer_id does not exist, then it returns a 400 error.
// @Accept       json
// @Produce      json
// @Param 		 buyer_id body string true "Buyer Id"
// @Success      200  {object}  data.Message
// @Failure      400  {object}  data.Message
// @Router       /buyers/resend-otp [post]
func handleResendOtp(c *gin.Context) {
	var resendOtpReq data.BuyerResendOtpData
	bindErr := c.ShouldBindJSON(&resendOtpReq)

	if bindErr != nil {
		r := data.Message{Message: "Bad Request Body"}
		c.JSON(http.StatusBadRequest, r)
		return
	}

	response, err := buyer.ResendOtp(db, resendOtpReq)

	if err != nil {
		r := data.Message{Message: err.Error()}
		c.JSON(err.ErrorCode(), r)
		return
	}

	c.JSON(http.StatusOK, &response)
}

// handleValidateOtp godoc
// @Summary      Validates a given otp from a specific buyer
// @Description  Checks to see if the provided buyer exists, if not returns a 400. Otherwise it checks to see if the otps match. If not it
// returns a 401 unauthorized. If successful, it returns buyer login response data but with updated verification state.
// @Accept       json
// @Produce      json
// @Param 		 buyer_id body string true "Buyer Id"
// @Param 		 otp body string true "Otp"
// @Success      200  {object}  data.BuyerLoginResponseData
// @Failure      400  {object}  data.Message
// @Failure      401  {object}  data.Message
// @Router       /buyers/validate-otp [post]
func handleValidateOtp(c *gin.Context) {
	var validateOtpReq data.BuyerValidateOtpData
	bindErr := c.ShouldBindJSON(&validateOtpReq)

	if bindErr != nil {
		r := data.Message{Message: "Bad Request Body"}
		c.JSON(http.StatusBadRequest, r)
		return
	}

	response, err := buyer.ValidateOtp(db, validateOtpReq)

	if err != nil {
		r := data.Message{Message: err.Error()}
		c.JSON(err.ErrorCode(), r)
		return
	}

	c.JSON(http.StatusOK, &response)
}
