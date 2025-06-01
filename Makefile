postgres:
	docker run --name postgres17_5 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=Welcome1234 -d postgres:17.5-alpine
createdb:
	docker exec -it postgres17_5 createdb --username=root --owner=root simple_bank

dropdb:
	docker exec -it postgres17_5 dropdb simple_bank

migrateup:
	migrate -path db/migration -database "postgresql://root:Welcome1234@localhost:5432/simple_bank?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://root:Welcome1234@localhost:5432/simple_bank?sslmode=disable" -verbose down

.PHONY: createdb dropdb postgres migrateup migratedown