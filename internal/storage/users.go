package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"workshop/internal/models"
)

type Storage struct {
	db *sql.DB
}

func NewStorage(db *sql.DB) Storage {
	return Storage{db: db}
}

func (s Storage) Create(ctx context.Context, name string) (models.User, error) {
	id := uuid.NewString()

	_, err := s.db.ExecContext(ctx, "INSERT INTO users (id, name) VALUES ($1, $2)", id, name)
	if err != nil {
		return models.User{}, fmt.Errorf("failed to execute insert: %w", err)
	}

	return models.User{ID: id, Name: name}, nil
}

func (s Storage) GetByID(ctx context.Context, ID string) (models.User, error) {
	var usr models.User

	err := s.db.QueryRowContext(ctx, "SELECT * FROM users WHERE id = $1", ID).Scan(&usr.ID, &usr.Name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, fmt.Errorf("failed to query user: %w", models.NotFoundErr)
		}

		return models.User{}, fmt.Errorf("failed to query user: %w", err)
	}

	return usr, nil
}
