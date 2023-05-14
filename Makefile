postgres12:
	docker run --name postgres12 -e POSTGRES_PASSWORD=123456 -e POSTGRES_USER=root -p 5432:5432 -v $(pwd)/data:/var/lib/postgresql/data -d postgres:12-alpine

createdb:
	docker exec -it postgres12 createdb --username=root --owner=root simple_bank

dropdb:
	docker exec -it postgres12 dropdb simple_bank

migrateup:
	migrate -path db/migration -database "postgresql://root:123456@localhost:5432/simple_bank?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://root:123456@localhost:5432/simple_bank?sslmode=disable" -verbose down

sqlc:
	sqlc generate

.PHONY: createdb dropdb postgres12 migrateup migratedown sqlc