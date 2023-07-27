package main

import (
	"BackendAPI/store"
	"database/sql"
	"log"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

var db *sql.DB

func main() {

	//Setup Router and Database Connection
	router := gin.Default()
	var err error

	db, err = store.SetupDB()

	if err != nil {
		log.Println("Could not connect to the database:", err)
	}

	apiGroup := router.Group("/api/v1")
	{
		buyerGroup := apiGroup.Group("/buyer")
		{
			buyerGroup.POST("/login", handleBuyerLogin)
			buyerGroup.POST("/signup", handleBuyerSignUp)
		}

		testGroup := apiGroup.Group("/test")
		{
			testGroup.GET("/ping", handlePing)
		}
	}

	router.Run(":8080")

	defer store.CloseDB(db)
}
