package db

import (
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/samber/do/v2"
	"github.com/xtsank/mypills-super-service/src/internal/infra/postgres/config"
)

func NewDB(i do.Injector) (*sqlx.DB, error) {
	cfg := do.MustInvoke[*config.Config](i)

	db, err := sqlx.Open("pgx", cfg.ConnectionString())
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(cfg.DBMaxOpenConns)
	db.SetMaxIdleConns(cfg.DBMaxIdleConns)
	db.SetConnMaxLifetime(cfg.DBConnMaxLifetime)
	db.SetConnMaxIdleTime(cfg.DBConnMaxIdleTime)

	return db, nil
}
