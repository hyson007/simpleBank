timescaledb :
	docker run -d --name timescaledb -p 5432:5432 \
-e POSTGRES_USER=dbusername -e POSTGRES_PASSWORD=dbpassword timescale/timescaledb:latest-pg14

createdb:
	docker exec -it timescaledb createdb --username=dbusername simpleBank

dropdb:
	docker exec -it timescaledb dropdb --username=dbusername simpleBank

migrateup:
	migrate -path db/migration -database "postgresql://dbusername:dbpassword@localhost:5432/simpleBank?sslmode=disable" -verbose up

sqlc:
	sqlc generate

migratedown:
	migrate -path db/migration -database "postgresql://dbusername:dbpassword@localhost:5432/simpleBank?sslmode=disable" -verbose down

test:
	go test -v -cover ./...

.PHONY: timescaledb, createdb, dropdb, migrateup, migratedown, sqlc, test