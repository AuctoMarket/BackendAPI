package main

import (
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

	// ping
	router.GET("/ping", handlePing)

	router.Run(":8080")

}

func handlePing(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"message": "pong"})
}
