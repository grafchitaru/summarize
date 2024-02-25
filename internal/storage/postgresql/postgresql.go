package postgresql

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

type Storage struct {
	pool *pgxpool.Pool
}

func New(connString string) (*Storage, error) {
	const op = "storage.postgresql.NewStorage"

	config, err := pgxpool.ParseConfig(connString)

	if err != nil {
		return nil, fmt.Errorf("%s: unable to parse config: %w", op, err)
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, fmt.Errorf("%s: unable to connect: %w", op, err)
	}

	return &Storage{pool: pool}, nil
}

func (s *Storage) Ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return s.pool.Ping(ctx)
}

func (s *Storage) Close() {
	s.pool.Close()
}

func (s *Storage) GetUser(login string) (string, error) {
	const op = "storage.postgresql.GetUser"

	var id string
	err := s.pool.QueryRow(context.Background(), "SELECT id FROM users WHERE login = $1", login).Scan(&id)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (s *Storage) GetUserPassword(login string) (string, error) {
	const op = "storage.postgresql.GetUserPassword"

	var password string
	err := s.pool.QueryRow(context.Background(), "SELECT password FROM users WHERE login = $1", login).Scan(&password)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return password, nil
}

func (s *Storage) Registration(id string, login string, password string) (string, error) {
	const op = "storage.postgresql.Registration"

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	tx, err := s.pool.Begin(ctx)

	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}
	defer tx.Rollback(context.Background())

	now := time.Now()

	_, err = tx.Exec(ctx, `
        INSERT INTO users(id, login, password, created_at)   
        VALUES($1, $2, $3, $4);
    `, id, login, password, now.Format("2006-01-02  15:04:05"))
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	if err := tx.Commit(context.Background()); err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}
