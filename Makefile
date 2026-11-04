.PHONY: build run swag test-classic test-london test test-integration test-integration-go test-e2e test-e2e-go test-offline coverage-branch install-gocove install-allure allure-generate allure-report allure-open test-clean clean docker-db log make_admin psql-users

login=admin111

DB_USER ?= user
DB_PASSWORD ?= password
DB_NAME ?= mypills

TEST_PACKAGES = ./src/internal/service
TEST_COVERPKG = ./src/internal/service,./src/internal/infra/postgres/repository
TEST_FLAGS = -v -count=1 -shuffle=on
ALLURE_RESULTS_DIR = $(CURDIR)/allure-results
ALLURE_REPORT_DIR = $(CURDIR)/allure-report
ALLURE_HISTORY_DIR = $(CURDIR)/allure-history-v2
LEGACY_ALLURE_RESULTS_DIR = $(CURDIR)/src/internal/service/allure-results
COVERAGE_PROFILE = $(CURDIR)/coverage.out
GOCOVE_META_DIR = $(CURDIR)/.gocove
GOCOVE_CACHE_BASE = $(CURDIR)/.tools/gocove-cache
SCANNER_WORK_DIR = $(CURDIR)/.scannerwork
BUILD_DIR = $(CURDIR)/build
BIN_DIR = $(CURDIR)/bin
FRONTEND_DIST_DIR = $(CURDIR)/gui/dist
FRONTEND_VITE_CACHE_DIR = $(CURDIR)/gui/.vite
FRONTEND_TSC_BUILD_INFO = $(CURDIR)/gui/tsconfig.tsbuildinfo
FRONTEND_NODE_TSC_BUILD_INFO = $(CURDIR)/gui/tsconfig.node.tsbuildinfo
GO_BIN = $(or $(shell go env GOBIN),$(shell go env GOPATH)/bin)
GOCOVE ?= $(GO_BIN)/gocove
GOCOVE_VERSION = v0.0.0-20260427173940-a0dceb9dadca
ALLURE_CLI ?= $(CURDIR)/node_modules/.bin/allure

build: swag
	@export DB_HOST=localhost; go build -o build/mypills-app ./src/cmd/app/main.go ./src/cmd/app/app.go

swag:
	swag init -q -g src/cmd/app/main.go -o docs/swagger --parseDependency --parseInternal --useStructName

run: swag build
	./build/mypills-app

test-classic: install-allure test-clean
	@mkdir -p "$(ALLURE_RESULTS_DIR)"
	@status=0; ALLURE_RESULTS_DIR="$(ALLURE_RESULTS_DIR)" go test $(TEST_FLAGS) -run '^TestClassic_' ./src/internal/service || status=$$?; $(MAKE) allure-generate || exit $$?; exit $$status

test-london: install-allure test-clean
	@mkdir -p "$(ALLURE_RESULTS_DIR)"
	@status=0; ALLURE_RESULTS_DIR="$(ALLURE_RESULTS_DIR)" go test $(TEST_FLAGS) -run '^TestLondon_' ./src/internal/service || status=$$?; $(MAKE) allure-generate || exit $$?; exit $$status

test: install-allure test-clean
	@mkdir -p "$(ALLURE_RESULTS_DIR)"
	@status=0; ALLURE_RESULTS_DIR="$(ALLURE_RESULTS_DIR)" go test $(TEST_FLAGS) -covermode=count -coverpkg=$(TEST_COVERPKG) -coverprofile="$(COVERAGE_PROFILE)" $(TEST_PACKAGES) || status=$$?; $(MAKE) allure-generate || exit $$?; exit $$status
	go tool cover -func="$(COVERAGE_PROFILE)"
	go run ./scripts/line-coverage.go "$(COVERAGE_PROFILE)"

test-integration:
	@status=0; ./scripts/test-ci-docker.sh ./scripts/test-stand.sh run docker compose -p "mypills-lab2-$${TEST_RUN_ID:-test}" -f docker-compose.test.yaml run --rm --no-deps test-runner make test-integration-go || status=$$?; ./scripts/test-ci-docker.sh ./scripts/test-stand.sh down || status=1; exit $$status

test-integration-go:
	@mkdir -p "$(ALLURE_RESULTS_DIR)"
	ALLURE_RESULTS_DIR="$(ALLURE_RESULTS_DIR)" go test -tags=integration -v -count=1 -shuffle=on -timeout=60s ./src/internal/integration

test-e2e:
	@status=0; TEST_SEED_FILE=e2e_seed.sql ./scripts/test-ci-docker.sh ./scripts/test-stand.sh run docker compose -p "mypills-lab2-$${TEST_RUN_ID:-test}" -f docker-compose.test.yaml run --rm --no-deps test-runner make test-e2e-go || status=$$?; TEST_SEED_FILE=e2e_seed.sql ./scripts/test-ci-docker.sh ./scripts/test-stand.sh down || status=1; exit $$status

test-e2e-go:
	@mkdir -p "$(ALLURE_RESULTS_DIR)"
	ALLURE_RESULTS_DIR="$(ALLURE_RESULTS_DIR)" go test -tags=e2e -v -count=1 -timeout=60s ./src/internal/e2e

test-offline: install-allure test-clean
	@mkdir -p "$(ALLURE_RESULTS_DIR)"
	@status=0; GOPROXY=off ALLURE_RESULTS_DIR="$(ALLURE_RESULTS_DIR)" go test $(TEST_FLAGS) -covermode=count -coverpkg=$(TEST_COVERPKG) -coverprofile="$(COVERAGE_PROFILE)" $(TEST_PACKAGES) || status=$$?; $(MAKE) allure-generate || exit $$?; exit $$status
	go tool cover -func="$(COVERAGE_PROFILE)"
	go run ./scripts/line-coverage.go "$(COVERAGE_PROFILE)"

coverage-branch:
	@command -v "$(GOCOVE)" >/dev/null || { echo "gocove не найден: выполните make install-gocove"; exit 1; }
	@rm -rf -- "$(GOCOVE_META_DIR)"
	@mkdir -p "$(GOCOVE_META_DIR)"
	GOCOVE_CACHE_BASE="$(GOCOVE_CACHE_BASE)" GOCOVE_META_DIR="$(GOCOVE_META_DIR)" ALLURE_RESULTS_DIR="$(GOCOVE_META_DIR)/allure-results" "$(GOCOVE)" test -a $(TEST_PACKAGES)
	"$(GOCOVE)" report --meta-dir="$(GOCOVE_META_DIR)"

install-gocove:
	go install github.com/srvgit/gocove/cmd/gocove@$(GOCOVE_VERSION)

install-allure:
	@test -x "$(ALLURE_CLI)" || npm ci

allure-generate:
	@test -d "$(ALLURE_RESULTS_DIR)" || { echo "Результаты Allure не найдены: сначала выполните тесты"; exit 1; }
	@printf 'Go=%s\nOS=%s\nArchitecture=%s\n' "$$(go env GOVERSION)" "$$(go env GOOS)" "$$(go env GOARCH)" > "$(ALLURE_RESULTS_DIR)/environment.properties"
	@printf '%s\n' '{"name":"go test","buildName":"MyPills unit tests"}' > "$(ALLURE_RESULTS_DIR)/executor.json"
	@if test -d "$(ALLURE_HISTORY_DIR)"; then cp -a "$(ALLURE_HISTORY_DIR)" "$(ALLURE_RESULTS_DIR)/history"; fi
	$(ALLURE_CLI) generate "$(ALLURE_RESULTS_DIR)" --clean --output "$(ALLURE_REPORT_DIR)" --lang ru --name "MyPills tests"
	@mkdir -p "$(ALLURE_HISTORY_DIR)"
	@cp -a "$(ALLURE_REPORT_DIR)/history/." "$(ALLURE_HISTORY_DIR)/"

allure-report:
	@test -f "$(ALLURE_REPORT_DIR)/index.html" || { echo "Отчёт Allure не найден: сначала выполните make test"; exit 1; }

allure-open: allure-report
	$(ALLURE_CLI) open "$(ALLURE_REPORT_DIR)"

test-clean:
	@rm -rf -- "$(ALLURE_RESULTS_DIR)" "$(ALLURE_REPORT_DIR)" "$(GOCOVE_META_DIR)" "$(LEGACY_ALLURE_RESULTS_DIR)" "$(SCANNER_WORK_DIR)"
	@rm -f -- "$(COVERAGE_PROFILE)" "$(CURDIR)/coverage.html"

clean: test-clean
	@rm -rf -- "$(BUILD_DIR)" "$(BIN_DIR)" "$(FRONTEND_DIST_DIR)" "$(FRONTEND_VITE_CACHE_DIR)"
	@rm -f -- "$(FRONTEND_TSC_BUILD_INFO)" "$(FRONTEND_NODE_TSC_BUILD_INFO)"

docker-db:
	sudo docker exec -it mypills-db psql -U user -d mypills

log:
	tail -n 200 -f "logs/app.log"

make_admin:
	./scripts/make_admin.sh $(login)

psql-users:
	sudo docker compose exec -T db psql -U $(DB_USER) -d $(DB_NAME) -c "select id, login, email from users limit 20;"
