timescaledb :
	docker run -d --name timescaledb --network simplebankPostgres -p 5432:5432 \
-e POSTGRES_USER=dbusername -e POSTGRES_PASSWORD=dbpassword timescale/timescaledb:latest-pg14

createdb:
	docker exec -it timescaledb createdb --username=dbusername simpleBank

dropdb:
	docker exec -it timescaledb dropdb --username=dbusername simpleBank

migrateup:
	migrate -path db/migration -database "postgresql://dbusername:dbpassword@localhost:5432/simpleBank?sslmode=disable" -verbose up

migrateup1:
	migrate -path db/migration -database "postgresql://dbusername:dbpassword@localhost:5432/simpleBank?sslmode=disable" -verbose up 1

sqlc:
	sqlc generate

migratedown:
	migrate -path db/migration -database "postgresql://dbusername:dbpassword@localhost:5432/simpleBank?sslmode=disable" -verbose down

migratedown1:
	migrate -path db/migration -database "postgresql://dbusername:dbpassword@localhost:5432/simpleBank?sslmode=disable" -verbose down 1

server:
	go run main.go

test:
	go test -v -cover ./...

mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/hyson007/simpleBank/db/sqlc Store

.PHONY: timescaledb, createdb, dropdb, migrateup, migrateup1, migratedown, migratedown1, sqlc, test, server, mock