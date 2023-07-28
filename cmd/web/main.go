package main

import (
<<<<<<< HEAD
<<<<<<< HEAD
	_ "BackendAPI/docs"
	"BackendAPI/store"

=======
	"BackendAPI/store"
>>>>>>> 005bc68 (Add login and signup API)
=======
	_ "BackendAPI/docs"
	"BackendAPI/store"

>>>>>>> d26034a (Add Swagger Documentation for API Endpoints)
	"database/sql"
	"log"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/lib/pq"
)

<<<<<<< HEAD
<<<<<<< HEAD
var db *sql.DB

// @title           AUCTO Backend API
// @version         1.0
// @description     This is the backend REST API for Aucto's marketplace, it is currently in v1.

// @host      localhost:8080
// @BasePath  /api/v1
<<<<<<< HEAD
=======
var Db *sql.DB
=======
var db *sql.DB
>>>>>>> e5d2750 (Add Tests for Login/Signup)

>>>>>>> 005bc68 (Add login and signup API)
=======
>>>>>>> d26034a (Add Swagger Documentation for API Endpoints)
func main() {

	//Setup Router and Database Connection
	router := gin.Default()
<<<<<<< HEAD
<<<<<<< HEAD
	var err error

	db, err = store.SetupDB()
=======
	Db, err := store.SetupDB()
>>>>>>> 005bc68 (Add login and signup API)
=======
	var err error

	db, err = store.SetupDB()
>>>>>>> e5d2750 (Add Tests for Login/Signup)

	if err != nil {
		log.Println("Could not connect to the database:", err)
	}

<<<<<<< HEAD
<<<<<<< HEAD
=======
	defer store.CloseDB(Db)

>>>>>>> 005bc68 (Add login and signup API)
=======
>>>>>>> e5d2750 (Add Tests for Login/Signup)
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
<<<<<<< HEAD
	defer store.CloseDB(db)
=======
>>>>>>> 005bc68 (Add login and signup API)
=======
	defer store.CloseDB(db)
>>>>>>> e5d2750 (Add Tests for Login/Signup)
}
