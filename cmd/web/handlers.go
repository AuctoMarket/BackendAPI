package main

import (
	"BackendAPI/api/buyer"
	"BackendAPI/data"
	"net/http"

	"github.com/gin-gonic/gin"
)

/*
Handles error and API response for the Login API for buyers
*/
func handleBuyerLogin(c *gin.Context) {
	var loginData data.LoginData
	err := c.ShouldBindJSON(&loginData)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Bad request body"})
		return
	}

	loginResponse, err := buyer.BuyerLogin(db, loginData)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, &loginResponse)
}

/*
Handles error and API response for the Sign Up API for buyers
*/
func handleBuyerSignUp(c *gin.Context) {
	var signUpData data.SignUpData
	err := c.ShouldBindJSON(&signUpData)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Bad request body"})
		return
	}

	signUpResponse, err := buyer.BuyerSignUp(db, signUpData)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, &signUpResponse)
}

/*
Test Ping as a sanity check
*/
func handlePing(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "pong"})
}
