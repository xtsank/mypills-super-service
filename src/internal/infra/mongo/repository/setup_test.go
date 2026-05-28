package repository

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"github.com/samber/do/v2"
	cfgpkg "github.com/xtsank/mypills-super-service/src/internal/config"
	"github.com/xtsank/mypills-super-service/src/internal/infra/mongo/db"
	"github.com/xtsank/mypills-super-service/src/internal/transport/middleware"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func setupTestMongo(t *testing.T) *mongo.Client {
	t.Helper()

	loadTestEnv(t)

	i := do.New()
	do.Provide(i, middleware.NewLogger)
	do.Provide(i, cfgpkg.NewConfig)

	var client *mongo.Client
	var err error
	for attempt := 0; attempt < 10; attempt++ {
		client, err = db.NewMongoDB(i)
		if err == nil {
			break
		}
		time.Sleep(500 * time.Millisecond)
	}
	if err != nil {
		t.Skipf("integration tests skipped: cannot connect to mongo: %v", err)
	}

	return client
}

func loadTestEnv(t *testing.T) {
	t.Helper()

	envPath := filepath.Join(repoRootPath(t), ".env")
	if _, err := os.Stat(envPath); err != nil {
		t.Logf(".env not found at %s: %v", envPath, err)
		return
	}

	if err := godotenv.Load(envPath); err != nil {
		t.Logf("failed to load .env from %s: %v", envPath, err)
	}
}

func repoRootPath(t *testing.T) string {
	t.Helper()

	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatalf("cannot resolve setup_test.go path")
	}

	return filepath.Clean(filepath.Join(filepath.Dir(file), "..", "..", "..", "..", ".."))
}
