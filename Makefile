.PHONY: build migrate-up migrate-down

build:
	docker build -t halosuster .

migrate-up:
	migrate -path database/migrations -database "postgres://digiboyz:digiboyz@localhost:5433/halosuster_db?sslmode=disable" -verbose up

migrate-down:
	migrate -path database/migrations -database "postgres://digiboyz:digiboyz@localhost:5433/halosuster_db?sslmode=disable" -verbose down

drop-db:
	migrate -path database/migrations -database "postgres://digiboyz:digiboyz@localhost:5433/halosuster_db?sslmode=disable" -verbose drop
