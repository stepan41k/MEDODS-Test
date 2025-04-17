package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v4"
	"github.com/stepan41k/MEDODS-Test/internal/storage"
)


func (s *Storage) Create(ctx context.Context, guid []byte, refreshToken []byte) error {
	const op = "storage.postgres.auth.Create"

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
			return
		}

		commitErr := tx.Commit(ctx)
		if commitErr != nil {
			err = fmt.Errorf("%s: %w", op, err)
		}
	}()

	row := tx.QueryRow(ctx, `
		UPDATE users
		SET refresh_token = $1
		WHERE guid = $2
		RETURNING guid;
	`, refreshToken, guid)

	var rguid []byte

	err = row.Scan(&rguid)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return fmt.Errorf("%s: %w", op, storage.ErrUserNotFound)
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}


func (s *Storage) GetRefresh(ctx context.Context, guid []byte) (string, error) {
	const op = "storage.postgres.auth.Refresh"

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
			return
		}

		commitErr := tx.Commit(ctx)
		if commitErr != nil {
			err = fmt.Errorf("%s: %w", op, err)
		}
	}()

	row := tx.QueryRow(ctx, `
		SELECT refresh_token
		FROM users
		WHERE guid = $1;
	`, guid)

	var refreshToken []byte

	err = row.Scan(&refreshToken)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return string(refreshToken), nil
}	