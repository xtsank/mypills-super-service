.PHONY: build run swag test-classic test-london test test-offline coverage-branch install-gocove allure-report allure-open test-clean clean docker-db log make_admin test_notification migration migration-back psql-users mongo-users

login=admin111

DB_USER = user
DB_PASSWORD = password
DB_NAME = mypills

TEST_PACKAGE = ./src/internal/service
TEST_FLAGS = -v -count=1 -shuffle=on
ALLURE_RESULTS_DIR = $(CURDIR)/allure-results
ALLURE_REPORT_DIR = $(CURDIR)/allure-report
LEGACY_ALLURE_RESULTS_DIR = $(CURDIR)/src/internal/service/allure-results
COVERAGE_PROFILE = $(CURDIR)/coverage.out
COVERAGE_HTML = $(CURDIR)/coverage.html
COVERAGE_JSON = $(CURDIR)/coverage.json
GOCOVE_META_DIR = $(CURDIR)/.gocove
BUILD_DIR = $(CURDIR)/build
BIN_DIR = $(CURDIR)/bin
FRONTEND_DIST_DIR = $(CURDIR)/gui/dist
FRONTEND_VITE_CACHE_DIR = $(CURDIR)/gui/.vite
FRONTEND_TSC_BUILD_INFO = $(CURDIR)/gui/tsconfig.tsbuildinfo
FRONTEND_NODE_TSC_BUILD_INFO = $(CURDIR)/gui/tsconfig.node.tsbuildinfo
GO_BIN = $(or $(shell go env GOBIN),$(shell go env GOPATH)/bin)
GOCOVE ?= $(GO_BIN)/gocove
GOCOVE_VERSION = v0.0.0-20260427173940-a0dceb9dadca
NPX ?= npx
ALLURE_VERSION = 3.17.0

build: swag
	@export DB_HOST=localhost; go build -o build/mypills-app ./src/cmd/app/main.go ./src/cmd/app/app.go

swag:
	swag init -q -g src/cmd/app/main.go -o docs/swagger --parseDependency --parseInternal --useStructName

run: swag build
	./build/mypills-app

test-classic: test-clean
	@mkdir -p "$(ALLURE_RESULTS_DIR)"
	ALLURE_RESULTS_DIR="$(ALLURE_RESULTS_DIR)" go test $(TEST_FLAGS) -run '^TestClassic_' $(TEST_PACKAGE)

test-london: test-clean
	@mkdir -p "$(ALLURE_RESULTS_DIR)"
	ALLURE_RESULTS_DIR="$(ALLURE_RESULTS_DIR)" go test $(TEST_FLAGS) -run '^TestLondon_' $(TEST_PACKAGE)

test: test-clean
	@mkdir -p "$(ALLURE_RESULTS_DIR)"
	ALLURE_RESULTS_DIR="$(ALLURE_RESULTS_DIR)" go test $(TEST_FLAGS) -covermode=count -coverprofile="$(COVERAGE_PROFILE)" $(TEST_PACKAGE)
	go tool cover -func="$(COVERAGE_PROFILE)"
	go tool cover -html="$(COVERAGE_PROFILE)" -o "$(COVERAGE_HTML)"

test-offline: test-clean
	@mkdir -p "$(ALLURE_RESULTS_DIR)"
	GOPROXY=off ALLURE_RESULTS_DIR="$(ALLURE_RESULTS_DIR)" go test $(TEST_FLAGS) -covermode=count -coverprofile="$(COVERAGE_PROFILE)" $(TEST_PACKAGE)
	go tool cover -func="$(COVERAGE_PROFILE)"

coverage-branch:
	@command -v "$(GOCOVE)" >/dev/null || { echo "gocove не найден: выполните make install-gocove"; exit 1; }
	@rm -rf -- "$(GOCOVE_META_DIR)"
	@mkdir -p "$(GOCOVE_META_DIR)" "$(ALLURE_RESULTS_DIR)"
	GOCOVE_META_DIR="$(GOCOVE_META_DIR)" ALLURE_RESULTS_DIR="$(ALLURE_RESULTS_DIR)" "$(GOCOVE)" test $(TEST_PACKAGE)
	"$(GOCOVE)" report --meta-dir="$(GOCOVE_META_DIR)"
	"$(GOCOVE)" report --meta-dir="$(GOCOVE_META_DIR)" --format=json > "$(COVERAGE_JSON)"
	"$(GOCOVE)" report --meta-dir="$(GOCOVE_META_DIR)" --format=html --output="$(COVERAGE_HTML)"

install-gocove:
	go install github.com/srvgit/gocove/cmd/gocove@$(GOCOVE_VERSION)

allure-report:
	@test -d "$(ALLURE_RESULTS_DIR)" || { echo "Результаты Allure не найдены: сначала выполните make test"; exit 1; }
	@rm -rf -- "$(ALLURE_REPORT_DIR)"
	$(NPX) --yes allure@$(ALLURE_VERSION) generate "$(ALLURE_RESULTS_DIR)" --output "$(ALLURE_REPORT_DIR)"

allure-open: allure-report
	$(NPX) --yes allure@$(ALLURE_VERSION) open "$(ALLURE_REPORT_DIR)"

test-clean:
	@rm -rf -- "$(ALLURE_RESULTS_DIR)" "$(ALLURE_REPORT_DIR)" "$(GOCOVE_META_DIR)" "$(LEGACY_ALLURE_RESULTS_DIR)"
	@rm -f -- "$(COVERAGE_PROFILE)" "$(COVERAGE_HTML)" "$(COVERAGE_JSON)"

clean: test-clean
	@rm -rf -- "$(BUILD_DIR)" "$(BIN_DIR)" "$(FRONTEND_DIST_DIR)" "$(FRONTEND_VITE_CACHE_DIR)"
	@rm -f -- "$(FRONTEND_TSC_BUILD_INFO)" "$(FRONTEND_NODE_TSC_BUILD_INFO)"

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
