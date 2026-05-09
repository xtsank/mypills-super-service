package repository

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"github.com/jmoiron/sqlx"
	"github.com/samber/do/v2"
	cfgpkg "github.com/xtsank/mypills-super-service/src/internal/infra/postgres/config"
	"github.com/xtsank/mypills-super-service/src/internal/infra/postgres/db"
)

// setupTestDB connects to the database using the same providers as main (config + db.NewDB).
// If connection fails (no docker-compose / env), the test will be skipped.
func setupTestDB(t *testing.T) *sqlx.DB {
	t.Helper()

	loadTestEnv(t)

	i := do.New()
	do.Provide(i, cfgpkg.NewConfig)

	var d *sqlx.DB
	var err error
	for attempt := 0; attempt < 10; attempt++ {
		d, err = db.NewDB(i)
		if err == nil {
			break
		}
		time.Sleep(500 * time.Millisecond)
	}
	if err != nil {
		t.Skipf("integration tests skipped: cannot connect to DB: %v", err)
	}

	return d
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
