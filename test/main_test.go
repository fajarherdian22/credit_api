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
	payload, err := repo.GetCustomers(context.Background(), "Budi@gmail.com")
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	// arg := repository.GenerateLimitParams{
	// 	ID:         uuid.NewString(),
	// 	CustomerID: "3c5b6f7b-60ab-43cc-ae3f-f05a2fba57e5",
	// 	Tenor:      4,
	// 	Limit:      700000,
	// }
	// err = repo.GenerateLimit(context.Background(), arg)
	// require.NoError(t, err)

	for i := 1; i <= 4; i++ {
		salary := float64(10000000)
		switch i {
		case 1:
			salary = salary * 0.15
		case 2:
			salary = salary * 0.2
		case 3:
			salary = salary * 0.25
		case 4:
			salary = salary * 0.3
		}
		arg := repository.GenerateLimitParams{
			ID:         uuid.NewString(),
			CustomerID: "d81e05f9-279d-493c-9a12-a6c331436b91",
			Tenor:      int32(i),
			Limit:      salary,
		}
		err = repo.GenerateLimit(context.Background(), arg)
		fmt.Print(err.Error())
		require.NoError(t, err)
	}

}
