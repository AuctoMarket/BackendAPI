package main

import (
	"BackendAPI/api/buyer"
	"BackendAPI/store"

	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {

	//Setup Router and Database Connection
	router := gin.Default()
	db, err := store.SetupDB()

	if err != nil {
		fmt.Println("Could not connect to the database:", err)
	}

	defer store.CloseDB(db)

	apiGroup := router.Group("/api/v1")
	{
		buyerGroup := apiGroup.Group("/buyer")
		{
			buyerGroup.POST("/login", buyer.HandleBuyerLogin)
			buyerGroup.POST("/signup", buyer.HandleBuyerSignUp)
		}

		testGroup := apiGroup.Group("/test")
		{
			testGroup.GET("/ping", handlePing)
		}
	}

	router.Run(":8080")

}

func handlePing(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"message": "pong"})
}
