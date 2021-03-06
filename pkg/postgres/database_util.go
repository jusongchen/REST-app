package postgres

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"testing"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/ory/dockertest"
	"github.com/sethvargo/go-retry"

	// imported to register the postgres migration driver
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	// imported to register the "file" source migration driver
	_ "github.com/golang-migrate/migrate/v4/source/file"
	// imported to register the "postgres" database driver for migrate
)

const (
	//this is for testing only
	dbMigrationScriptLocation = "../../migration"
)

// NewTestDatabaseWithConfig creates a new database suitable for use in testing.
// This should not be used outside of testing, but it is exposed in the main
// package so it can be shared with other packages.
//
// All database tests can be skipped by running `go test -short` or by setting
// the `SKIP_DATABASE_TESTS` environment variable.
func NewTestDatabaseWithConfig(tb testing.TB) (*DB, *Config) {
	tb.Helper()

	if testing.Short() {
		tb.Skipf("🚧 Skipping database tests (short!")
	}

	if skip, _ := strconv.ParseBool(os.Getenv("SKIP_DATABASE_TESTS")); skip {
		tb.Skipf("🚧 Skipping database tests (SKIP_DATABASE_TESTS is set)!")
	}

	// Context.
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Create the pool (docker instance).
	pool, err := dockertest.NewPool("")
	if err != nil {
		tb.Fatalf("failed to create Docker pool: %s", err)
	}

	// Start the container.
	dbname, username, password := "en-server", "my-username", "abcd1234"
	container, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "12-alpine",
		Env: []string{
			"LANG=C",
			"POSTGRES_DB=" + dbname,
			"POSTGRES_USER=" + username,
			"POSTGRES_PASSWORD=" + password,
		},
	})
	if err != nil {
		tb.Fatalf("failed to start postgres container: %s", err)
	}

	// Ensure container is cleaned up.
	tb.Cleanup(func() {
		if err := pool.Purge(container); err != nil {
			tb.Fatalf("failed to cleanup postgres container: %s", err)
		}
	})

	// Get the host. On Mac, Docker runs in a VM.
	host := container.Container.NetworkSettings.IPAddress
	if runtime.GOOS == "darwin" {
		host = net.JoinHostPort(container.GetBoundIP("5432/tcp"), container.GetPort("5432/tcp"))
	}

	// Build the connection URL.
	connURL := &url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(username, password),
		Host:   host,
		Path:   dbname,
	}
	q := connURL.Query()
	q.Add("sslmode", "disable")
	connURL.RawQuery = q.Encode()

	// Wait for the container to start - we'll retry connections in a loop below,
	// but there's no point in trying immediately.
	time.Sleep(1 * time.Second)

	b, err := retry.NewFibonacci(500 * time.Millisecond)
	if err != nil {
		tb.Fatalf("failed to configure backoff: %v", err)
	}
	b = retry.WithMaxRetries(10, b)
	b = retry.WithCappedDuration(10*time.Second, b)

	// Establish a connection to the database. Use a Fibonacci backoff instead of
	// exponential so wait times scale appropriately.
	var dbpool *pgxpool.Pool
	if err := retry.Do(ctx, b, func(ctx context.Context) error {
		var err error
		dbpool, err = pgxpool.Connect(ctx, connURL.String())
		if err != nil {
			return retry.RetryableError(err)
		}
		return nil
	}); err != nil {
		tb.Fatalf("failed to start postgres: %s", err)
	}

	// Run the migrations.
	if err := dbMigrate(connURL.String(), dbMigrationsDir(dbMigrationScriptLocation)); err != nil {
		tb.Fatalf("failed to migrate database: %s", err)
	}

	// Create the db instance.
	db := &DB{Pool: dbpool}

	// Close db when done.
	tb.Cleanup(func() {
		db.Close(context.Background())
	})

	return db, &Config{
		Name:     dbname,
		User:     username,
		Host:     container.GetBoundIP("5432/tcp"),
		Port:     container.GetPort("5432/tcp"),
		SSLMode:  "disable",
		Password: password,
	}
}

//NewTestDatabase is a test helper func
func NewTestDatabase(tb testing.TB) *DB {
	tb.Helper()

	db, _ := NewTestDatabaseWithConfig(tb)
	return db
}

//DbMigrate runs DB migration scripts
func DbMigrate(dbConnURL, scriptDir string) error {
	return dbMigrate(dbConnURL, scriptDir)
}

// dbMigrate runs the migrations. u is the connection URL string (e.g.
// postgres://...).
func dbMigrate(u string, scriptDir string) error {
	// Run the migrations
	migrationsDir := fmt.Sprintf("file://%s", scriptDir)
	m, err := migrate.New(migrationsDir, u)
	if err != nil {
		return fmt.Errorf("failed create migrate: %w", err)
	}
	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("failed run migrate: %w", err)
	}
	srcErr, dbErr := m.Close()
	if srcErr != nil {
		return fmt.Errorf("migrate source error: %w", srcErr)
	}
	if dbErr != nil {
		return fmt.Errorf("migrate database error: %w", dbErr)
	}
	return nil
}

// dbMigrationsDir returns the path on disk to the migrations. It uses
// runtime.Caller() to get the path to the caller
func dbMigrationsDir(dir string) string {
	_, filename, _, ok := runtime.Caller(1)
	if !ok {
		return ""
	}
	return filepath.Join(filepath.Dir(filename), dir)
}
