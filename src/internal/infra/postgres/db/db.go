package db

import (
	"log/slog"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/samber/do/v2"
	apperrors "github.com/xtsank/mypills-super-service/src/internal/errors"
	"github.com/xtsank/mypills-super-service/src/internal/infra/postgres/config"
)

func NewDB(i do.Injector) (*sqlx.DB, error) {
	cfg := do.MustInvoke[*config.Config](i)
	logger := do.MustInvoke[*slog.Logger](i)

	db, err := sqlx.Open("pgx", cfg.ConnectionString())
	if err != nil {
		logger.Error("failed to open db", slog.Any("error", err))
		return nil, apperrors.ErrInternal.WithError(err)
	}

	if err := db.Ping(); err != nil {
		logger.Error("failed to ping db", slog.Any("error", err))
		return nil, apperrors.ErrInternal.WithError(err)
	}

	db.SetMaxOpenConns(cfg.DBMaxOpenConns)
	db.SetMaxIdleConns(cfg.DBMaxIdleConns)
	db.SetConnMaxLifetime(cfg.DBConnMaxLifetime)
	db.SetConnMaxIdleTime(cfg.DBConnMaxIdleTime)

	return db, nil
}
