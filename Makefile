run:
	go run main.go

migrate-up:
	migrate -path database/migration/ -database "postgresql://root:root@localhost:5434/postgres?sslmode=disable" -verbose up

migrate-down:
	migrate -path database/migration/ -database "postgresql://root:root@localhost:5434/postgres?sslmode=disable" -verbose down