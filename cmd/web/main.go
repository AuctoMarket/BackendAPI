package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {

	router := gin.Default()

	// ping
	router.GET("/ping", handlePing)

	router.Run(":8080")
}

func handlePing(ctx *gin.Context) {
	ctx.JSON(http.StatusNotFound, gin.H{"message": "pong"})
}
