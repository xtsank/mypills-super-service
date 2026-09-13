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
- MongoDB: MongoDB 7 в Docker Compose, официальный `go.mongodb.org/mongo-driver/v2`.
- Авторизация: JWT через `github.com/golang-jwt/jwt/v5`, хеширование паролей через `golang.org/x/crypto`.
- Идентификаторы: `github.com/google/uuid`.
- Конфигурация: переменные окружения и `.env` через `github.com/joho/godotenv`.
- OpenAPI/Swagger: `swaggo/swag`, `gin-swagger`, `swaggo/files`.
- Логирование: стандартный `log/slog`.
- Email: SMTP; для локальной среды в Compose предусмотрен Mailpit.
- Тесты: стандартный `testing`, `github.com/stretchr/testify`, `go-sqlmock` и адаптер Allure Go.

### Frontend

- React 18 + React DOM.
- TypeScript 5.6 со включённым `strict`.
- Vite 5 и `@vitejs/plugin-react`.
- Архитектура GUI прямо описана в `gui/README.md` как MVVM.
- Роутер, библиотека управления состоянием, UI-kit и тестовый фреймворк: **не найдено**.

### Инфраструктура

- Docker и Docker Compose.
- Два профиля Compose: `postgres` и `mongo`; каждому соответствует свой backend- и GUI-сервис.
- Bash-скрипты для миграции данных, назначения администратора и демонстрационного наполнения уведомлений.

## 2. Сборка, запуск и тесты

### Backend

```bash
# Генерация Swagger
make swag

# Генерация Swagger и сборка build/mypills-app
make build

# Генерация Swagger, сборка и запуск бинарника
make run

# Пять классических или пять лондонских тестов
make test-classic
make test-london

# Все 10 тестов, случайный порядок, Allure results и line coverage
make test

# Тот же набор без доступа к сети
make test-offline

# Установка gocove и отчёты branch/condition coverage
make install-gocove
make coverage-branch

# Генерация и открытие HTML-отчёта Allure
make allure-report
make allure-open

# Только тестовые артефакты / все артефакты проекта
make test-clean
make clean
```

Особенности найденных команд:

- `make swag`, `make build` и `make run` требуют доступную в `PATH` CLI-команду `swag`; команда её установки в проекте **не найдена**.
- `make build` пишет бинарник в `build/mypills-app`, но сам `Makefile` не создаёт каталог `build/`.
- `make test`, `make test-classic` и `make test-london` запускают пакет `src/internal/service`; отдельная цель для всех Go-пакетов: **не найдена**.
- `make test-clean` удаляет результаты тестов, Allure и покрытия. `make clean` дополнительно удаляет backend/frontend build-артефакты.
- Для `make allure-report` и `make allure-open` требуется `npx`; для `make coverage-branch` требуется предварительный `make install-gocove`.

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

- Готовая документированная команда запуска всего Compose-стека: **не найдена**.
- `docker-compose.yaml` требует значения переменных из `.env`, включая параметры БД и порты. Пример `.env`: **не найден**.
- Миграции между PostgreSQL и MongoDB:

```bash
make migration       # pg2mongo, затем очистка PostgreSQL
make migration-back  # mongo2pg, затем очистка MongoDB
```

- Другие цели `Makefile`: `make docker-db`, `make log`, `make make_admin login=<login>`, `make test_notification`, `make psql-users`, `make mongo-users`.

## 3. Структура директорий

```text
.
├── AGENTS.md
├── README.md              # описание предметной области и диаграммы
├── Makefile               # backend, тесты, миграции и служебные цели
├── Dockerfile             # production-like backend image
├── docker-compose.yaml    # PostgreSQL/MongoDB, backend, GUI, Mailpit
├── go.mod / go.sum        # Go-модуль и зависимости
├── docs/
│   ├── swagger/           # сгенерированная OpenAPI-документация
│   ├── bpmn/              # BPMN-исходник и SVG
│   ├── c4/                # C4-модель и SVG
│   ├── class/             # диаграмма классов
│   ├── dm_schema/         # DBML-схема и SVG
│   ├── erd/               # ERD в GraphML/PNG/SVG
│   ├── seq/               # sequence-диаграммы
│   └── use-case/          # use-case диаграмма
├── scripts/               # migration, make_admin, notification seed
├── src/
│   ├── cmd/
│   │   ├── app/           # точка входа HTTP-приложения и DI wiring
│   │   └── migrator/      # перенос данных PostgreSQL ↔ MongoDB
│   └── internal/
│       ├── config/        # чтение env-конфигурации
│       ├── domain/        # доменные сущности и интерфейсы репозиториев
│       ├── dto/           # внутренние DTO
│       ├── errors/        # ошибки приложения
│       ├── infra/         # PostgreSQL, MongoDB и SMTP
│       ├── service/       # бизнес-логика и command-объекты
│       ├── test_util/     # builders, Object Mother, fixtures, helpers и моки
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
- HTTP request/response DTO разделены на `transport/dto/req` и `transport/dto/res`.
- Swagger-аннотации находятся над Gin handlers; после изменения публичных endpoint/DTO документация обновляется целью `make swag`.
- SQL-запросы PostgreSQL вынесены в `src/internal/infra/postgres/queries`, а структуры отображения БД — в `.../entity`.
- Формальная конфигурация Go-линтера/форматтера: **не найдена**. `gofmt -l` выводит ряд текущих `.go`-файлов, поэтому утверждать, что весь репозиторий уже отформатирован `gofmt`, нельзя.
- `.editorconfig`, ESLint и Prettier: **не найдено**.
- Документированные требования к именованию Go-сущностей, покрытию тестами и порядку импортов: **не найдено**.

## 5. Что НЕ трогать

- Не редактировать вручную `docs/swagger/docs.go`, `swagger.json` и `swagger.yaml`: они создаются `swag init` через `make swag`; `docs.go` явно помечен `Code generated ... DO NOT EDIT`.
- Не добавлять в Git `.env`/`.env.*`, логи, `node_modules/`, `dist/`, `.vite/` и `*.tsbuildinfo`: эти пути исключены в `.gitignore`. Исключение для примера окружения в `.gitignore` **не найдено**.
- Не запускать `make migration` или `make migration-back` без явной необходимости и проверки данных: обе цели задают `MIGRATION_CLEANUP=true`, а migrator после успешного копирования очищает исходную БД.
- Не перезаписывать несвязанные пользовательские изменения в рабочем дереве. На момент создания файла уже изменены `go.mod`, `go.sum` и добавлены тестовые файлы; они не изменялись при подготовке этого документа.
- Дополнительный явно документированный список защищённых файлов или директорий: **не найдено**.
