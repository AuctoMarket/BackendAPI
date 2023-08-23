package main

import (
	_ "BackendAPI/docs"
	"BackendAPI/store"
	"BackendAPI/utils"
	"context"
	"os"

	"database/sql"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware

	_ "github.com/lib/pq"
)

var ginLambda *ginadapter.GinLambda
var db *sql.DB
var s3Client *s3.Client

// @title           AUCTO Backend API
// @version         1.0
// @description     This is the REST API for Aucto's marketplace, it is currently in v1.

// @host      *
// @BasePath  /api/v1
func main() {

	//Setup Router
	router := gin.Default()
	//Setup CORS Middlware
	router.Use(cors.Default())
	var err error

	// Load .env variables from environment file
	loadErr := utils.LoadDotEnv(".env")

	if loadErr != nil {
		utils.LogError(loadErr, "Cannot fetch .env, no .env file")
		return
	}

	//Setup DB connection
	db, err = store.SetupDB()
	if err != nil {
		log.Println("Could not connect to the database:", err)
	}
	//Setup S3 connection
	s3Client, err = store.CreateNewS3()
	if err != nil {
		log.Println("Could not connect to the S3 Instance:", err)
	}

	apiGroup := router.Group("/api/v1")
	{
		buyerGroup := apiGroup.Group("/buyers")
		{
			buyerGroup.POST("/login", handleBuyerLogin)
			buyerGroup.POST("/signup", handleBuyerSignUp)
		}

		productGroup := apiGroup.Group("/products")
		{
			productGroup.GET("/:id", handleGetProductById)
			productGroup.POST("", handleCreateProduct)
			productGroup.POST("/:id/images", handleCreateProductImages)
			productGroup.GET("", handleGetProductList)
		}

		sellerGroup := apiGroup.Group("/sellers")
		{
			sellerGroup.POST("/signup", handleSellerSignUp)
			sellerGroup.POST("/login", handleSellerLogin)
			sellerGroup.GET("/:id", handleGetSellerById)

		}

		orderGroup := apiGroup.Group("/orders")
		{
			orderGroup.POST("", handleCreateOrder)
			orderGroup.POST("/guest", handleCreateGuestOrder)
			orderGroup.GET("/:id", handleGetOrderById)
			orderGroup.GET("/:id/guest", handleGetGuestOrderById)
			// orderGroup.POST("/:id/payment-complete", handlePaymentComplete)
			// orderGroup.POST("guest/:id/payment-complete", handleGuestPaymentComplete)
		}

		testGroup := apiGroup.Group("/tests")
		{
			testGroup.GET("/ping", handlePing)
		}

		docGroup := apiGroup.Group("/docs")
		{
			docGroup.GET("/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
		}
	}

	env := os.Getenv("API_ENV")

	if env != "local" {
		ginLambda = ginadapter.New(router)

		lambda.Start(Handler)
	} else {
		router.Run(":8080")
	}

	defer store.CloseDB(db)
}

func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return ginLambda.ProxyWithContext(ctx, request)
}
