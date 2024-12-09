package main

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	"github.com/fajarherdian22/credit_bank/controller"
	"github.com/fajarherdian22/credit_bank/db"
	"github.com/fajarherdian22/credit_bank/repository"
	"github.com/fajarherdian22/credit_bank/service"
	"github.com/fajarherdian22/credit_bank/util"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal().AnErr("Cannot load config :", err)
	}

	dbCon := db.ConDB(config.DBSource)
	repo := repository.New(dbCon)

	custService := service.NewCustomerService(repo)
	custController := controller.NewCustomerController(custService)

	transactionService := service.NewTransactionService(repo)
	transactionController := controller.NewTransactionController(transactionService)

	router := gin.New()

	r := router.Group("/api/")
	r.POST("/customers/create", custController.CreateCustomersUser)
	r.POST("/customers/findcust", custController.GetCustomerId)
	r.POST("/customers/transaction", transactionController.CreateTransaction)

	err = router.Run(":8080")
	if err != nil {
		log.Fatal().AnErr("Failed to start server:", err)
	}
	log.Info().Msg("Running Server")

}
