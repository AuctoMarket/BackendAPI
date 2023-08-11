package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

/*
Test Ping as a sanity check
*/
func handlePing(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "pong"})
}
