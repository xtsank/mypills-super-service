.PHONY: build run swag

login=admin111

build: swag
	go build -o build/mypills-app ./src/cmd/app/main.go ./src/cmd/app/app.go

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
