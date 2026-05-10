.PHONY: swag run

swag:
	swag init -q -g src/cmd/app/main.go -o docs/swagger --parseDependency --parseInternal --useStructName

run: swag
	go mod tidy
	go run src/cmd/app/main.go src/cmd/app/app.go

test_bl:
	go test -v -cover -count=1 ./src/internal/service

test_repo:
	go test -v -cover ./src/internal/infra/postgres/repository

docker-db:
	sudo docker exec -it mypills-db psql -U user -d mypills