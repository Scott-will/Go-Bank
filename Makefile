postgres:
	docker run --name postgres17_5 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=Welcome1234 -d postgres:17.5-alpine
createdb:
	docker exec -it postgres17_5 createdb --username=root --owner=root simple_bank

dropdb:
	docker exec -it postgres17_5 dropdb simple_bank

migrateup:
	migrate -path db/migration -database "postgresql://root:Welcome1234@localhost:5432/simple_bank?sslmode=disable" -verbose up

migrateup1:
	migrate -path db/migration -database "postgresql://root:Welcome1234@localhost:5432/simple_bank?sslmode=disable" -verbose up 1


migratedown:
	migrate -path db/migration -database "postgresql://root:Welcome1234@localhost:5432/simple_bank?sslmode=disable" -verbose down


migratedown1:
	migrate -path db/migration -database "postgresql://root:Welcome1234@localhost:5432/simple_bank?sslmode=disable" -verbose down 1

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go

mock:
	mockgen -package mockdb -destination db/mock/store.go go-bank/db/sqlc Store

.PHONY: createdb dropdb postgres migrateup migrateup1 migratedown migratedown1 server mock