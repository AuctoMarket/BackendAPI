package main

import (
	"BackendAPI/api/buyer"
	"BackendAPI/data"
	"net/http"

	"github.com/gin-gonic/gin"
)

<<<<<<< HEAD

// handleBuyerLogin godoc
// @Summary      Logs a buyer into their account
// @Description  Checks to see if a user email exists and if supplied password matches the stored password
// @Accept       json
// @Produce      json
// @Param 		 email body string true "email"
// @Param 		 password body string true "password"
// @Success      200  {object}  data.LoginResponseData
// @Failure      400  {object}  data.Message
// @Failure      401  {object}  data.Message
// @Failure      500  {object}  data.Message
// @Router       /buyer/login/ [post]
func handleBuyerLogin(c *gin.Context) {
	var loginData data.LoginData
	bindErr := c.ShouldBindJSON(&loginData)

	if bindErr != nil {
		r := data.Message{Message: "Bad Request Body"}
		c.JSON(http.StatusBadRequest, r)
=======
/*
Handles error and API response for the Login API for buyers
*/
func handleBuyerLogin(c *gin.Context) {
	var loginData data.LoginData
	bindErr := c.ShouldBindJSON(&loginData)

<<<<<<< HEAD
	if err != nil {
<<<<<<< HEAD
		c.JSON(http.StatusBadRequest, gin.H{"message": "Bad request data"})
>>>>>>> 005bc68 (Add login and signup API)
=======
		c.JSON(http.StatusBadRequest, gin.H{"message": "Bad request body"})
>>>>>>> 0b4235d (Changer Status to 201, Update error msg)
=======
	if bindErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Bad Request Body"})
>>>>>>> c9cf9e0 (Update Error Handling)
		return
	}

	loginResponse, err := buyer.BuyerLogin(db, loginData)

	if err != nil && err.Error() == "Something went wrong" {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	if err != nil {
<<<<<<< HEAD
<<<<<<< HEAD
		r := data.Message{Message: err.Error()}
		c.JSON(err.ErrorCode(), r)
=======
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
>>>>>>> 005bc68 (Add login and signup API)
=======
		c.JSON(err.ErrorCode(), gin.H{"message": err.Error()})
>>>>>>> c9cf9e0 (Update Error Handling)
		return
	}

	c.JSON(http.StatusOK, &loginResponse)
}

<<<<<<< HEAD
// handleBuyerSignup godoc
// @Summary      Signs a new buyer up
// @Description  Checks to see if a user email exists and if not creates a new account with supplied email and password
// @Accept       json
// @Produce      json
// @Param 		 email body string true "email"
// @Param 		 password body string true "password"
// @Success      200  {object}  data.LoginResponseData
// @Failure      400  {object}  data.Message
// @Failure      401  {object}  data.Message
// @Failure      500  {object}  data.Message
// @Router       /buyer/signup/ [post]
func handleBuyerSignUp(c *gin.Context) {
	var signUpData data.SignUpData
	bindErr := c.ShouldBindJSON(&signUpData)

	if bindErr != nil {
		r := data.Message{Message: "Bad Request Body"}
		c.JSON(http.StatusBadRequest, r)
		return
	}

	signUpResponse, err := buyer.BuyerSignUp(db, signUpData)

	if err != nil {
		r := data.Message{Message: err.Error()}
		c.JSON(err.ErrorCode(), r)
		return
	}

	c.JSON(http.StatusCreated, &signUpResponse)
=======
/*
Handles error and API response for the Sign Up API for buyers
*/
func handleBuyerSignUp(c *gin.Context) {
	var signUpData data.SignUpData
	bindErr := c.ShouldBindJSON(&signUpData)

<<<<<<< HEAD
<<<<<<< HEAD
>>>>>>> 005bc68 (Add login and signup API)
=======
	if err != nil {
=======
	if bindErr != nil {
>>>>>>> c9cf9e0 (Update Error Handling)
		c.JSON(http.StatusBadRequest, gin.H{"message": "Bad request body"})
		return
	}

	signUpResponse, err := buyer.BuyerSignUp(db, signUpData)

	if err != nil && err.Error() == "Something went wrong" {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	if err != nil {
		c.JSON(err.ErrorCode(), gin.H{"message": err.Error()})
		return
	}

<<<<<<< HEAD
	c.JSON(http.StatusOK, &signUpResponse)
>>>>>>> e5d2750 (Add Tests for Login/Signup)
=======
	c.JSON(http.StatusCreated, &signUpResponse)
>>>>>>> 0b4235d (Changer Status to 201, Update error msg)
}

/*
Test Ping as a sanity check
*/
func handlePing(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "pong"})
}
