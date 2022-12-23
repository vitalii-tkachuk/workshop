package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"workshop/internal/models"
	"workshop/pkg/logger"
	"workshop/pkg/response"

	"github.com/go-chi/chi/v5"
)

type CreateUserParams struct {
	Name string `json:"name"`
}

//go:generate moq -rm -out users_mock.go . UsersService
type UsersService interface {
	Create(ctx context.Context, name string) (models.User, error)
	Get(ctx context.Context, id string) (models.User, error)
}

type Users struct {
	user UsersService
}

func NewUsers(us UsersService) Users {
	return Users{us}
}

func (u Users) Routes() http.Handler {
	r := chi.NewRouter()

	r.Post("/", u.Create)
	r.Get("/{userId}", u.Get)

	return r
}

func (u Users) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)

	var userParams CreateUserParams
	if err := json.NewDecoder(r.Body).Decode(&userParams); err != nil {
		log.Error().Err(err).Msg("failed to parse params")
		response.InternalError(w)
		return
	}

	user, err := u.user.Create(ctx, userParams.Name)
	if err != nil {
		if errors.Is(err, models.UserCreateParamInvalidNameErr) {
			response.BadRequest(w)
			return
		}

		log.Error().Err(err).Msg("failed to create user")
		response.InternalError(w)
		return
	}

	if err := json.NewEncoder(w).Encode(user); err != nil {
		log.Error().Err(err).Msg("failed to encode response")
		response.InternalError(w)
		return
	}
}

func (u Users) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)

	userID := chi.URLParam(r, "userId")

	user, err := u.user.Get(ctx, userID)
	if err != nil {
		if errors.Is(err, models.NotFoundErr) {
			response.NotFound(w)
			return
		}

		log.Error().Err(err).Msg("failed to get user")
		response.InternalError(w)
		return
	}

	if err := json.NewEncoder(w).Encode(user); err != nil {
		log.Error().Err(err).Msg("failed to encode response")
		response.InternalError(w)
		return
	}
}
