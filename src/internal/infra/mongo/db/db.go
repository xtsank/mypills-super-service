package db

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/samber/do/v2"
	"github.com/xtsank/mypills-super-service/src/internal/config"
	apperrors "github.com/xtsank/mypills-super-service/src/internal/errors"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

func NewMongoDB(i do.Injector) (*mongo.Client, error) {
	cfg := do.MustInvoke[*config.Config](i)
	logger := do.MustInvoke[*slog.Logger](i)

	connOptions := options.Client().ApplyURI(ConnectionString(cfg))
	connOptions.SetMaxPoolSize(uint64(cfg.DBMaxOpenConns))
	connOptions.SetMinPoolSize(uint64(cfg.DBMaxIdleConns))
	connOptions.SetMaxConnIdleTime(cfg.DBConnMaxIdleTime)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db, err := mongo.Connect(connOptions)
	if err != nil {
		logger.Error("failed to connect to mongo", slog.Any("error", err))
		return nil, apperrors.ErrInternal.WithError(err)
	}

	if err := db.Ping(ctx, readpref.Primary()); err != nil {
		logger.Error("failed to ping mongo", slog.Any("error", err))
		return nil, apperrors.ErrInternal.WithError(err)
	}

	return db, nil
}

func ConnectionString(c *config.Config) string {
	return fmt.Sprintf("mongodb://%s:%s@%s:%s/%s?authSource=admin", c.DBUser, c.DBPassword, c.DBHost, c.DBPort, c.DBName)
}
