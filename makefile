create_migrate:
	migrate create -ext .sql -dir db/migration -seq name_migrate
migrateup:
	migrate -path db/migration -database "mysql://admin:admin1234@tcp(localhost:3306)/credit_db" --verbose up
migratedown:
	migrate -path db/migration -database "mysql://admin:admin1234@tcp(localhost:3306)/credit_db" --verbose down
server:
	go run main.go
sqlc:
	sqlc generate
