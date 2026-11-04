# AGENTS.md

## Область проекта

- Каталог, содержащий этот файл, — корень приложения и отдельного Git-репозитория.
- Соседние с ним каталоги `../testing_docs/` и `../web_docs/` содержат учебные PDF и не входят в исходный код приложения.
- Все команды ниже, если не указано иное, запускаются из каталога, содержащего этот файл.

## 1. Стек и ключевые библиотеки

### Backend

- Go `1.26.2` (значение из `go.mod` и образа в `Dockerfile`).
- HTTP API: `github.com/gin-gonic/gin`.
- Dependency injection: `github.com/samber/do/v2`.
- PostgreSQL: PostgreSQL 16 в Docker Compose, `pgx/v5` и `sqlx` в Go-коде.
- Авторизация: JWT через `github.com/golang-jwt/jwt/v5`, хеширование паролей через `golang.org/x/crypto`.
- Идентификаторы: `github.com/google/uuid`.
- Конфигурация: переменные окружения и `.env` через `github.com/joho/godotenv`.
- OpenAPI/Swagger: `swaggo/swag`, `gin-swagger`, `swaggo/files`.
- Логирование: стандартный `log/slog`.
- Тесты: стандартный `testing`, `github.com/stretchr/testify`, `github.com/mattn/go-sqlite3` для SQLite `:memory:` и адаптер Allure Go. SQLite-драйвер использует CGO.

### Frontend

- React 18 + React DOM.
- TypeScript 5.6 со включённым `strict`.
- Vite 5 и `@vitejs/plugin-react`.
- Архитектура GUI прямо описана в `gui/README.md` как MVVM.
- Роутер, библиотека управления состоянием, UI-kit и тестовый фреймворк: **не найдено**.

### Инфраструктура

- Docker и Docker Compose.
- Compose запускает PostgreSQL, backend и GUI без профилей.

## 2. Сборка, запуск и тесты

### Backend

```bash
# Генерация Swagger
make swag

# Генерация Swagger и сборка build/mypills-app
make build

# Генерация Swagger, сборка и запуск бинарника
make run

# Шесть классических или пять лондонских тестов двух сервисов
make test-classic
make test-london

# Все 11 тестов ЛР №1, случайный порядок, Allure results и coverage
make test

# Тот же набор без доступа к сети
make test-offline

# Установка gocove и отчёты branch/condition coverage
make install-gocove
make coverage-branch

# Запуск ZIP-установки SonarQube и анализ после make test
make sonar-up SONAR_HOME=/путь/к/sonarqube
make sonar-status
SONAR_TOKEN=... make sonar-scan

# Проверка и открытие готового HTML-отчёта Allure
make allure-report
make allure-open

# Только тестовые артефакты / все артефакты проекта
make test-clean
make clean
```

Особенности найденных команд:

- `make swag`, `make build` и `make run` требуют доступную в `PATH` CLI-команду `swag`; команда её установки в проекте **не найдена**.
- `make build` пишет бинарник в `build/mypills-app`, но сам `Makefile` не создаёт каталог `build/`.
- `make test` запускает пакет `src/internal/service`: 6 классических и 5 лондонских тестов `AdminService` и `CabinetService`. Сценарии, SQLite-границы и подготовка данных описаны в `docs/lr-testing/lab1.md`.
- `make test` выводит statement coverage через Go и число инструментированных строк через `scripts/line-coverage.go`.
- `make sonar-scan` импортирует `coverage.out` в доступный сервер SonarQube; нужны `sonar-scanner`, `SONAR_TOKEN` и ZIP-установка сервера для `make sonar-up`. Адрес по умолчанию `http://localhost:9000`.
- `make test-clean` удаляет результаты и отчёт, сохраняя `allure-history-v2/` для трендов. `make clean` дополнительно удаляет backend/frontend build-артефакты.
- `make test` и `make test-offline` запускают тесты и создают отчёт официальным Allure 2.46.1. `make install-allure` устанавливает локальную CLI из `package-lock.json` через `npm ci`. Перед генерацией `Makefile` записывает версию Go, ОС и архитектуру в `allure-results/environment.properties`, а сведения о запуске через `go test` — в `allure-results/executor.json`; переносит созданную Allure папку истории из предыдущего прогона в `allure-results/history`, после генерации сохраняет новую историю в `allure-history-v2/`. `make allure-report` проверяет наличие отчёта, `make allure-open` открывает его. Для `make coverage-branch` требуется предварительный `make install-gocove`.

### Frontend

Команды запускаются из `gui/`:

```bash
npm install
npm run dev
npm run build
npm run preview
```

- `npm run dev` запускает Vite на порту `5173`.
- `npm run build` выполняет `tsc -b && vite build`.
- `npm run preview` запускает предпросмотр Vite на порту `5173`.
- Команда frontend-тестов и lint-команда: **не найдено**.

### Docker Compose и служебные команды

- Весь стек запускается командой `docker compose up -d --build`.
- `docker-compose.yaml` требует значения переменных из `.env`, включая параметры БД и порты. Пример `.env`: **не найден**.
- Другие цели `Makefile`: `make docker-db`, `make log`, `make make_admin login=<login>`, `make psql-users`.

## 3. Структура директорий

```text
.
├── AGENTS.md
├── README.md              # описание предметной области и диаграммы
├── Makefile               # backend, тесты и служебные цели
├── Dockerfile             # production-like backend image
├── docker-compose.yaml    # PostgreSQL, backend, GUI
├── go.mod / go.sum        # Go-модуль и зависимости
├── docs/
│   ├── swagger/           # сгенерированная OpenAPI-документация
│   ├── bpmn/              # BPMN-исходник и SVG
│   ├── c4/                # C4-модель и SVG
│   ├── dm_schema/         # DBML-схема и SVG
│   ├── erd/               # ERD в GraphML/PNG/SVG
│   ├── seq/               # sequence-диаграммы
│   └── use-case/          # use-case диаграмма
├── scripts/               # проверка API и служебные команды
├── src/
│   ├── cmd/
│   │   └── app/           # точка входа HTTP-приложения и DI wiring
│   └── internal/
│       ├── config/        # чтение env-конфигурации
│       ├── domain/        # доменные сущности и интерфейсы репозиториев
│       ├── dto/           # внутренние DTO
│       ├── errors/        # ошибки приложения
│       ├── infra/         # PostgreSQL
│       ├── service/       # бизнес-логика и command-объекты
│       ├── test_util/     # Builder, Object Mother, SQLite-хелперы, ID и моки
│       └── transport/     # Gin handlers, middleware, request/response DTO
└── gui/
    ├── src/
    │   ├── api/           # HTTP-клиент и TypeScript DTO
    │   ├── app/           # корневой React-компонент и CSS
    │   ├── components/    # переиспользуемые UI-компоненты
    │   ├── store/         # React Context providers
    │   ├── viewmodels/    # MVVM hooks с логикой экранов
    │   └── views/         # компоненты экранов
    ├── package.json / package-lock.json
    ├── tsconfig.json / tsconfig.node.json
    └── vite.config.ts
```

## 4. Стиль кода

Подтверждённые проектом правила и соглашения:

- Frontend должен сохранять MVVM-разделение, указанное в `gui/README.md`: экраны находятся в `views/`, логика экранов — в hooks из `viewmodels/`, обращения к backend и типы — в `api/`, общее состояние — в `store/`.
- TypeScript проверяется в строгом режиме (`strict: true`), с `noEmit`, `isolatedModules` и JSX runtime `react-jsx`.
- В существующем frontend-коде используются именованные экспорты, двойные кавычки, точки с запятой, отступ в два пробела; React-компоненты названы в `PascalCase`, hooks — с префиксом `use`. Это наблюдаемые соглашения, отдельным formatter/linter они не закреплены.
- Backend разделён на `domain` → `service` → `transport` и реализации в `infra`; сборка зависимостей выполняется в `src/cmd/app/app.go` через `do/v2`.
- Для маршрутов v2 обработчик разбирает входные данные в `req DTO`, без проверок предметной области создаёт `command` и вызывает сервис. Сервис работает с доменными сущностями через конструкторы и методы, вызывает репозитории и возвращает `res.New...Dto`. Обработчик передаёт результат и HTTP-статус.
- HTTP request/response DTO разделены на `transport/dto/req` и `transport/dto/res`.
- Swagger-аннотации находятся над Gin handlers; после изменения публичных endpoint/DTO документация обновляется целью `make swag`.
- SQL-запросы PostgreSQL вынесены в `src/internal/infra/postgres/queries`, а структуры отображения БД — в `.../entity`.
- Формальная конфигурация Go-линтера/форматтера: **не найдена**. `gofmt -l` выводит ряд текущих `.go`-файлов, поэтому утверждать, что весь репозиторий уже отформатирован `gofmt`, нельзя.
- `.editorconfig`, ESLint и Prettier: **не найдено**.
- Документированные требования к именованию Go-сущностей, покрытию тестами и порядку импортов: **не найдено**.

## 5. Что НЕ трогать

- Не редактировать вручную `docs/swagger/docs.go`, `swagger.json` и `swagger.yaml`: они создаются `swag init` через `make swag`; `docs.go` явно помечен `Code generated ... DO NOT EDIT`.
- Не добавлять в Git `.env`/`.env.*`, логи, `node_modules/`, `dist/`, `.vite/` и `*.tsbuildinfo`: эти пути исключены в `.gitignore`. Исключение для примера окружения в `.gitignore` **не найдено**.
- Не перезаписывать несвязанные пользовательские изменения в рабочем дереве. На момент создания файла уже изменены `go.mod`, `go.sum` и добавлены тестовые файлы; они не изменялись при подготовке этого документа.
- Дополнительный явно документированный список защищённых файлов или директорий: **не найдено**.

## 6. Контракт API

- Основной контракт — `docs/openapi/v2.yaml`; GUI использует `/api/v2`.
- `/api/v1` и прежние корневые маршруты поддерживаются обработчиками `legacy_*`.
- PostgreSQL — единственное хранилище приложения.
