package main

import (
	"fmt"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog/log"

	"github.com/fajarherdian22/credit_bank/controller"
	"github.com/fajarherdian22/credit_bank/db"
	"github.com/fajarherdian22/credit_bank/middleware"
	"github.com/fajarherdian22/credit_bank/repository"
	"github.com/fajarherdian22/credit_bank/service"
	"github.com/fajarherdian22/credit_bank/token"
	"github.com/fajarherdian22/credit_bank/util"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal().AnErr("Cannot load config :", err)
	}

	dbCon := db.ConDB(config.DBSource)
	repo := repository.New(dbCon)

	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		fmt.Printf("cannot generate token %s", err.Error())
	}

	validate := validator.New()
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("ProductName", util.ProductNameValidator)
	}

	custService := service.NewCustomerService(repo)
	custController := controller.NewCustomerController(custService, tokenMaker, validate)

	transactionService := service.NewTransactionService(dbCon)
	transactionController := controller.NewTransactionController(transactionService, tokenMaker, validate)

	loanService := service.NewLoanService(repo)
	loanController := controller.NewLoanController(loanService, tokenMaker)

	router := gin.New()
	router.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
		AllowMethods:    []string{"POST", "GET"},
		AllowHeaders:    []string{"Content-Type", "Origin"},
		ExposeHeaders:   []string{"Content-Length"},
		MaxAge:          12 * time.Hour,
	}))
	router.Use(gin.Recovery())

	r := router.Group("/api/")
	r.POST("/customers/create", custController.CreateCustomersUser)
	r.POST("/customers/login", middleware.RateLimiter(), custController.LoginCustomers)
	r.POST("/token/refresh", custController.RenewAccessToken)

	authGroup := r.Group("/")

	authGroup.Use(middleware.AuthMiddleware(tokenMaker))
	authGroup.POST("/customers/transaction", transactionController.CreateTransaction)
	authGroup.GET("/customers/listtx", transactionController.ListTx)
	authGroup.GET("request/loan", loanController.GenerateLimit)

	err = router.Run(":8080")
	if err != nil {
		log.Fatal().AnErr("Failed to start server:", err)
	}
	log.Info().Msg("Running Server")

}
