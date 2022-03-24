package store

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/jackc/tern/migrate"
)

const ErrorNotFound = "Entity not found"

var ctx context.Context = context.Background()

type Store struct {
	*pgxpool.Pool
}

func New(url string) (*Store, error) {
	pool, err := new(url)
	return &Store{Pool: pool}, err
}

func new(url string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.Connect(context.Background(), url)
	if err != nil {
		return nil, err
	}

	return pool, nil
}

func (s *Store) MigrateDatabase(migrationOutput string) (string, error) {
	connectionPool, err := s.Acquire(ctx)
	if err != nil {
		return "", errors.New(fmt.Sprintf("Unable to acquire db connection: %v\n", err))
	}
	defer connectionPool.Release()
	conn := connectionPool.Conn()
	defer conn.Close(ctx)
	migrator, err := migrate.NewMigrator(ctx, conn, "schema_version")
	if err != nil {
		return "", errors.New(fmt.Sprintf("Unable to create a migrator: %v\n", err))
	}

	err = migrator.LoadMigrations(migrationOutput)
	if err != nil {
		return "", errors.New(fmt.Sprintf("Unable to load migrations: %v\n", err))
	}

	err = migrator.Migrate(ctx)
	if err != nil {
		return "", errors.New(fmt.Sprintf("Unable to migrate: %v\n", err))
	}

	ver, err := migrator.GetCurrentVersion(ctx)
	if err != nil {
		return "", errors.New(fmt.Sprintf("Unable to get current schema version: %v\n", err))
	}

	migrationCount := len(migrator.Migrations)
	diff := migrationCount - int(ver)
	if diff == 0 {
		return "Found no migrations to apply.", nil
	}
	return fmt.Sprintf("Migrations are done. Current schema version: %v\n", ver), nil
}

func (s *Store) Query(sql string, args ...interface{}) (pgx.Rows, error) {
	ctx := context.Background()
	conn, err := s.Acquire(ctx)
	if err != nil {
		return nil, err
	}

	rows, err := conn.Query(ctx, sql, args...)
	defer conn.Release()
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func (s *Store) QueryRow(sql string, args ...interface{}) (pgx.Row, error) {
	ctx := context.Background()
	conn, err := s.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	row := conn.QueryRow(ctx, sql, args...)
	defer conn.Release()
	if err != nil {
		return nil, err
	}
	return row, nil
}

func (s *Store) Execute(sql string, args ...interface{}) (bool, error) {
	ctx := context.Background()
	conn, err := s.Acquire(ctx)
	if err != nil {
		return false, err
	}

	tag, err := conn.Exec(ctx, sql, args...)
	defer conn.Release()
	if err != nil {
		return false, err
	}

	r := tag.RowsAffected()
	if r == 0 {
		return false, errors.New(ErrorNotFound)
	}

	return true, nil
}
