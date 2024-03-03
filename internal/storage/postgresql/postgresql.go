package postgresql

import (
	"context"
	"errors"
	"fmt"
	"github.com/grafchitaru/summarize/internal/models"
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

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var id string
	err := s.pool.QueryRow(ctx, "SELECT id FROM users WHERE login = $1", login).Scan(&id)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return "", fmt.Errorf("%s: operation timed out: %w", op, err)
		}
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (s *Storage) GetUserPassword(login string) (string, error) {
	const op = "storage.postgresql.GetUserPassword"

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	var password string
	err := s.pool.QueryRow(ctx, "SELECT password FROM users WHERE login = $1", login).Scan(&password)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return "", fmt.Errorf("%s: operation timed out: %w", op, err)
		}
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
		return "", fmt.Errorf("%s begin: %w", op, err)
	}
	defer tx.Rollback(ctx)

	now := time.Now()

	_, err = tx.Exec(ctx, `
        INSERT INTO users(id, login, password, created_at, updated_at)   
        VALUES($1, $2, $3, $4, $5);
    `, id, login, password, now.Format("2006-01-02 15:04:05"), now.Format("2006-01-02 15:04:05"))
	if err != nil {
		return "", fmt.Errorf("%s exec: %w", op, err)
	}

	if err := tx.Commit(ctx); err != nil {
		return "", fmt.Errorf("%s commit: %w", op, err)
	}

	return id, nil
}

func (s *Storage) CreateSummarize(id string, user_id string, text string, status string, tokens int) error {
	const op = "storage.postgresql.CreateSummarize"

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	tx, err := s.pool.Begin(ctx)

	if err != nil {
		return fmt.Errorf("%s begin: %w", op, err)
	}
	defer tx.Rollback(context.Background())

	now := time.Now()

	_, err = tx.Exec(ctx, `
        INSERT INTO summarize(id, user_id, created_at, updated_at, text, status, tokens)   
        VALUES($1, $2, $3, $4, $5, $6, $7);
    `, id, user_id, now.Format("2006-01-02  15:04:05"), now.Format("2006-01-02  15:04:05"), text, status, tokens)
	if err != nil {
		return fmt.Errorf("%s exec: %w", op, err)
	}

	if err := tx.Commit(context.Background()); err != nil {
		return fmt.Errorf("%s commit: %w", op, err)
	}

	return nil
}
func (s *Storage) UpdateSummarizeStatus(id string, status string) error {
	const op = "storage.postgresql.UpdateSummarizeStatus"

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	tx, err := s.pool.Begin(ctx)

	if err != nil {
		return fmt.Errorf("%s begin: %w", op, err)
	}
	defer tx.Rollback(context.Background())

	now := time.Now()

	_, err = tx.Exec(ctx, `
        UPDATE summarize 
        SET status = $1, updated_at = $2
        WHERE id = $3;
    `, status, now.Format("2006-01-02  15:04:05"), id)
	if err != nil {
		return fmt.Errorf("%s exec: %w", op, err)
	}

	if err := tx.Commit(context.Background()); err != nil {
		return fmt.Errorf("%s commit: %w", op, err)
	}

	return nil
}
func (s *Storage) UpdateSummarizeResult(id string, status string, result string) error {
	const op = "storage.postgresql.UpdateSummarizeResult"

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	tx, err := s.pool.Begin(ctx)

	if err != nil {
		return fmt.Errorf("%s begin: %w", op, err)
	}
	defer tx.Rollback(context.Background())

	now := time.Now()

	_, err = tx.Exec(ctx, `
        UPDATE summarize 
        SET status = $1, result = $2, updated_at = $3
        WHERE id = $4;
    `, status, result, now.Format("2006-01-02  15:04:05"), id)
	if err != nil {
		return fmt.Errorf("%s exec: %w", op, err)
	}

	if err := tx.Commit(context.Background()); err != nil {
		return fmt.Errorf("%s commit: %w", op, err)
	}

	return nil
}

func (s *Storage) GetSummarize(id string, user_id string) (models.Summarize, error) {
	const op = "storage.postgresql.GetSummarize"

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	var summarize models.Summarize

	err := s.pool.QueryRow(ctx, "SELECT * FROM summarize WHERE id = $1 AND user_id = $2", id, user_id).Scan(&summarize.Id, &summarize.UserId, &summarize.CreatedAt, &summarize.UpdatedAt, &summarize.Text, &summarize.Result, &summarize.Status, &summarize.Tokens)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return models.Summarize{}, fmt.Errorf("%s: operation timed out: %w", op, err)
		}
		return models.Summarize{}, fmt.Errorf("%s: %w", op, err)
	}

	return summarize, nil
}

func (s *Storage) GetSummarizeByText(text string) (string, error) {
	const op = "storage.postgresql.GetSummarizeByText"

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	var id string
	err := s.pool.QueryRow(ctx, "SELECT id FROM summarize WHERE text = $1", text).Scan(&id)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return "", fmt.Errorf("%s: operation timed out: %w", op, err)
		}
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (s *Storage) GetStat(user_id string) ([]models.Stat, error) {
	const op = "storage.postgresql.GetStat"

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	rows, err := s.pool.Query(ctx, "SELECT user_id, status, count(id), sum(tokens) AS tokens FROM summarize WHERE user_id = $1 GROUP BY user_id, status;", user_id)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var stats []models.Stat
	for rows.Next() {
		var stat models.Stat
		if err := rows.Scan(&stat.UserId, &stat.Status, &stat.Count, &stat.Tokens); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		stats = append(stats, stat)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return stats, nil
}

func (s *Storage) GetStatus(user_id string, AiMaxLimitCount int, AiMaxLimitTokens int) (models.Status, error) {
	const op = "storage.postgresql.GetStatus"

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	var status models.Status
	err := s.pool.QueryRow(ctx, "SELECT count(id) AS count, sum(tokens) AS tokens FROM summarize WHERE user_id = $1", user_id).Scan(&status.Count, &status.Tokens)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return models.Status{}, fmt.Errorf("%s: operation timed out: %w", op, err)
		}
		return models.Status{}, fmt.Errorf("%s: %w", op, err)
	}
	status.Count = AiMaxLimitCount - status.Count
	status.Tokens = AiMaxLimitTokens - status.Tokens

	return status, nil
}
