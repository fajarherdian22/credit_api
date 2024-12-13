package test

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/fajarherdian22/credit_bank/db"
	"github.com/fajarherdian22/credit_bank/repository"
	"github.com/fajarherdian22/credit_bank/util"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestDb(t *testing.T) {
	config, err := util.LoadConfig("../")
	require.NoError(t, err)
	if err != nil {
		log.Fatal("Cannot load config :", err)
	}

	dbCon := db.ConDB(config.DBSource)
	repo := repository.New(dbCon)

	payload, err := repo.GetCustomers(context.Background(), "budi@gmail.com")
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	for i := 1; i <= 4; i++ {
		salary := float64(0)
		switch i {
		case 1:
			salary = 100000
		case 2:
			salary = 200000
		case 3:
			salary = 500000
		case 4:
			salary = 700000
		}

		arg := repository.GenerateLimitParams{
			ID:         uuid.NewString(),
			CustomerID: payload.ID,
			Tenor:      int32(i),
			Limit:      salary,
		}

		err := repo.GenerateLimit(context.Background(), arg)
		if err != nil {
			fmt.Println("Error inserting limit:", err.Error())
			require.NoError(t, err)
		}
	}

	payload2, err := repo.GetCustomers(context.Background(), "anisa@gmail.com")
	require.NoError(t, err)
	require.NotEmpty(t, payload2)

	for i := 1; i <= 4; i++ {
		salary := float64(0)
		switch i {
		case 1:
			salary = 1000000
		case 2:
			salary = 1200000
		case 3:
			salary = 1500000
		case 4:
			salary = 2000000
		}

		arg := repository.GenerateLimitParams{
			ID:         uuid.NewString(),
			CustomerID: payload2.ID,
			Tenor:      int32(i),
			Limit:      salary,
		}

		err := repo.GenerateLimit(context.Background(), arg)
		if err != nil {
			fmt.Println("Error inserting limit:", err.Error())
			require.NoError(t, err)
		}
	}
}
