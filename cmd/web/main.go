package main

import (
	_ "BackendAPI/docs"
	"BackendAPI/store"

	"database/sql"
	"log"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/lib/pq"
)

var db *sql.DB

// @title           AUCTO Backend API
// @version         1.0
// @description     This is the backend REST API for Aucto's marketplace, it is currently in v1.

// @host      localhost:8080
// @BasePath  /api/v1
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

		swaggerGroup := apiGroup.Group("/swagger")
		{
			swaggerGroup.GET("/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
		}
	}

	router.Run(":8080")

	defer store.CloseDB(db)
}
