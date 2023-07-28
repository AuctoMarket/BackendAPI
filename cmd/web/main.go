package main

import (
<<<<<<< HEAD
	_ "BackendAPI/docs"
	"BackendAPI/store"

=======
	"BackendAPI/store"
>>>>>>> 005bc68 (Add login and signup API)
	"database/sql"
	"log"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/lib/pq"
)

<<<<<<< HEAD
var db *sql.DB

// @title           AUCTO Backend API
// @version         1.0
// @description     This is the backend REST API for Aucto's marketplace, it is currently in v1.

// @host      localhost:8080
// @BasePath  /api/v1
=======
var Db *sql.DB

>>>>>>> 005bc68 (Add login and signup API)
func main() {

	//Setup Router and Database Connection
	router := gin.Default()
<<<<<<< HEAD
	var err error

	db, err = store.SetupDB()
=======
	Db, err := store.SetupDB()
>>>>>>> 005bc68 (Add login and signup API)

	if err != nil {
		log.Println("Could not connect to the database:", err)
	}

<<<<<<< HEAD
=======
	defer store.CloseDB(Db)

>>>>>>> 005bc68 (Add login and signup API)
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

<<<<<<< HEAD
	defer store.CloseDB(db)
=======
>>>>>>> 005bc68 (Add login and signup API)
}
