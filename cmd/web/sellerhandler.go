package main

import (
	"BackendAPI/api/seller"
	"BackendAPI/data"
	"net/http"

	"github.com/gin-gonic/gin"
)

// handleSellerSignup godoc
// @Summary      Signs a new seller up
// @Description  Checks to see if a seller email does not already exists if so creates a new
// seller account with supplied email, password and seller_name
// if not returns a bad request error (400).
// @Accept       json
// @Produce      json
// @Param 		 email body string true "Sellers email [UNIQUE]"
// @Param 		 password body string true "Sellers password as plaintext"
// @Param 		 seller_name body string true "Sellers seller alias that is displayed as their seller name [UNIQUE]"
// @Success      200  {object}  data.SellerResponseData
// @Failure      400  {object}  data.Message
// @Failure      500  {object}  data.Message
// @Router       /sellers/signup [post]
func handleSellerSignUp(c *gin.Context) {
	var signUpData data.SellerSignUpData
	bindErr := c.ShouldBindJSON(&signUpData)

	if bindErr != nil {
		r := data.Message{Message: "Bad Request Body"}
		c.JSON(http.StatusBadRequest, r)
		return
	}

	signUpResponse, err := seller.SellerSignUp(db, signUpData)

	if err != nil {
		r := data.Message{Message: err.Error()}
		c.JSON(err.ErrorCode(), r)
		return
	}

	c.JSON(http.StatusCreated, &signUpResponse)
}

// handleSellerLogin godoc
// @Summary      Logs a seller into their account
// @Description  Checks to see if a sellers email exists and if supplied password matches the stored password
// if not returns a unauthorized error (401).
// @Accept       json
// @Produce      json
// @Param 		 email body string true "Sellers email"
// @Param 		 password body string true "Sellers password as plaintext"
// @Success      200  {object}  data.SellerResponseData
// @Failure      400  {object}  data.Message
// @Failure      401  {object}  data.Message
// @Failure      500  {object}  data.Message
// @Router       /sellers/login [post]
func handleSellerLogin(c *gin.Context) {
	var loginData data.UserLoginData
	bindErr := c.ShouldBindJSON(&loginData)

	if bindErr != nil {
		r := data.Message{Message: "Bad Request Body"}
		c.JSON(http.StatusBadRequest, r)
		return
	}

	loginResponse, err := seller.SellerLogin(db, loginData)

	if err != nil {
		r := data.Message{Message: err.Error()}
		c.JSON(err.ErrorCode(), r)
		return
	}

	c.JSON(http.StatusOK, &loginResponse)
}
