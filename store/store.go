package store

import (
	"context"
	"errors"
	"fmt"
	"log"

	//_ "github.com/jackc/pgx/stdlib"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/jackc/tern/migrate"
)

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

	//defer pool.Close()
	return pool, nil
}

func (s *Store) MigrateDatabase(migrationOutput string) error {
	pool, err := s.Acquire(context.Background())
	if err != nil {
		return errors.New(fmt.Sprintf("Unable to acquire db connection: %v\n", err))
	}
	migrator, err := migrate.NewMigrator(context.Background(), pool.Conn(), "schema_version")
	if err != nil {
		return errors.New(fmt.Sprintf("Unable to create a migrator: %v\n", err))
	}

	err = migrator.LoadMigrations(migrationOutput)
	if err != nil {
		return errors.New(fmt.Sprintf("Unable to load migrations: %v\n", err))
	}

	err = migrator.Migrate(context.Background())
	if err != nil {
		return errors.New(fmt.Sprintf("Unable to migrate: %v\n", err))
	}

	ver, err := migrator.GetCurrentVersion(context.Background())
	if err != nil {
		return errors.New(fmt.Sprintf("Unable to get current schema version: %v\n", err))
	}
	pool.Release()
	log.Printf("Migrations are done. Current schema version: %v\n", ver)
	return nil
}

func (s *Store) Query(sql string, args ...interface{}) (pgx.Rows, error) {
	ctx := context.Background()
	conn, err := s.Acquire(ctx)
	if err != nil {
		return nil, err
	}

	rows, err := conn.Query(ctx, sql, args...)
	conn.Release()
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

func (s *Store) Execute(sql string, args ...interface{}) error {
	ctx := context.Background()
	conn, err := s.Acquire(ctx)
	if err != nil {
		return err
	}
	_, err = conn.Exec(ctx, sql, args...)
	conn.Release()
	if err != nil {
		return err
	}
	return nil
}
