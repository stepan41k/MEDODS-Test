package postgres

import (
	"context"
	"fmt"
)

func (s *Storage) CreateUser(ctx context.Context, guid []byte) ([]byte, error) {
	const op = "storage.postgres.auth.Create"

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
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
		INSERT INTO users(guid)
		VALUES($1)
		RETURNING guid;
	`, guid)

	var nguid []byte

	err = row.Scan(&nguid)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return nguid, nil
}