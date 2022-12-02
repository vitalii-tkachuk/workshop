package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"net/http"
	"workshop/internal/users"

	"workshop/internal/models"
)

type CreateUserParams struct {
	Name string `json:"name"`
}

//go:generate moq -rm -out users_mock.go . UsersService
type UsersService interface {
	Create(ctx context.Context, name string) (models.User, error)
}

type Users struct {
	user UsersService
	repo users.Repository
}

func NewUsers(us UsersService, repo users.Repository) Users {
	return Users{us, repo}
}

func (u Users) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var userParams CreateUserParams
	if err := json.NewDecoder(r.Body).Decode(&userParams); err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	user, err := u.user.Create(ctx, userParams.Name)
	if err != nil {
		if errors.Is(err, models.UserCreateParamInvalidNameErr) {
			http.Error(w, "", http.StatusBadRequest)
		}

		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(user); err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
}

func (u Users) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID := chi.URLParam(r, "userId")

	user, err := u.repo.GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "", http.StatusBadRequest)
		}

		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(user); err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
}
