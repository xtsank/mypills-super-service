.PHONY: build run swag

login=admin111

DB_USER = user
DB_PASSWORD = password
DB_NAME = mypills

build: swag
	@export DB_HOST=localhost; go build -o build/mypills-app ./src/cmd/app/main.go ./src/cmd/app/app.go

swag:
	swag init -q -g src/cmd/app/main.go -o docs/swagger --parseDependency --parseInternal --useStructName

run: swag build
	./build/mypills-app

test_bl:
	go test -v -cover -count=1 ./src/internal/service

test_repo:
	@export DB_HOST=localhost; go test -v -cover -count=1 ./src/internal/infra/postgres/repository

test: test_bl test_repo

docker-db:
	sudo docker exec -it mypills-db psql -U user -d mypills

log:
	tail -n 200 -f "logs/app.log"

make_admin:
	./scripts/make_admin.sh $(login)

test_notification:
	sudo docker compose exec -T db bash -lc "/docker-entrypoint-initdb.d/seed_notification_demo.sh timofey1 etima16122005@gmail.com 12345678"

migration:
	MIGRATION_DIRECTION=pg2mongo MIGRATION_CLEANUP=true ./scripts/migration.sh

migration-back:
	MIGRATION_DIRECTION=mongo2pg MIGRATION_CLEANUP=true ./scripts/migration.sh

psql-users:
	sudo docker compose exec -T db psql -U $(DB_USER) -d $(DB_NAME) -c "select id, login, email from users limit 20;"

mongo-users:
	sudo docker compose exec -T mongo mongosh --quiet --username $(DB_USER) --password $(DB_PASSWORD) --authenticationDatabase admin --eval "db=db.getSiblingDB('$(DB_NAME)'); db.users.find({}, {login:1,email:1}).limit(20).pretty()"
