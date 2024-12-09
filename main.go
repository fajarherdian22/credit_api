package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
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

	custService := service.NewCustomerService(repo)
	custController := controller.NewCustomerController(custService, tokenMaker)

	transactionService := service.NewTransactionService(repo)
	transactionController := controller.NewTransactionController(transactionService, tokenMaker)

	router := gin.New()

	r := router.Group("/api/")
	r.POST("/customers/create", custController.CreateCustomersUser)
	r.POST("/customers/login", custController.LoginCustomers)
	authGroup := r.Group("/")

	authGroup.Use(middleware.AuthMiddleware(tokenMaker))
	authGroup.POST("/customers/transaction", transactionController.CreateTransaction)
	authGroup.GET("/customers/listtx", transactionController.ListTx)

	err = router.Run(":8080")
	if err != nil {
		log.Fatal().AnErr("Failed to start server:", err)
	}
	log.Info().Msg("Running Server")

}
