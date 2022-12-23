package handlers

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"workshop/internal/models"
	"workshop/internal/users"

	"github.com/stretchr/testify/assert"
)

func TestUsers_Create(t *testing.T) {
	type fields struct {
		user UsersService
	}
	type args struct {
		w *httptest.ResponseRecorder
		r *http.Request
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		wantCode int
		wantBody []byte
	}{
		{
			name: "success",
			fields: fields{
				user: &UsersServiceMock{
					CreateFunc: func(ctx context.Context, name string) (models.User, error) {
						assert.Equal(t, "mike", name)
						return models.User{Name: name, ID: "1"}, nil
					},
				},
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"name": "mike"}`)),
			},
			wantCode: http.StatusOK,
			wantBody: []byte(`{"id":"1","name":"mike"}` + "\n"),
		},
		{
			name: "error",
			fields: fields{
				user: &UsersServiceMock{
					CreateFunc: func(ctx context.Context, name string) (models.User, error) {
						return models.User{}, fmt.Errorf("invalid name argument: %w", models.UserCreateParamInvalidNameErr)
					},
				},
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest("POST", "/", strings.NewReader(`{"name": ""}`)),
			},
			wantCode: http.StatusBadRequest,
			wantBody: []byte("Bad Request" + "\n"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := Users{
				user: tt.fields.user,
			}
			u.Create(tt.args.w, tt.args.r)

			assert.Equal(t, tt.wantCode, tt.args.w.Code)
			assert.Equal(t, tt.wantBody, tt.args.w.Body.Bytes(), "unexpected body")
		})
	}
}

func TestUsers_Get(t *testing.T) {
	type fields struct {
		repo users.Repository
	}
	type args struct {
		w *httptest.ResponseRecorder
		r *http.Request
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		wantCode int
		wantBody []byte
	}{
		{
			name: "success",
			fields: fields{
				repo: &users.RepositoryMock{
					GetByIDFunc: func(ctx context.Context, ID string) (models.User, error) {
						return models.User{Name: "mike", ID: "1"}, nil
					},
				},
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest("POST", "/users/1", strings.NewReader(``)),
			},
			wantCode: http.StatusOK,
			wantBody: []byte(`{"id":"1","name":"mike"}` + "\n"),
		},
		{
			name: "not found",
			fields: fields{
				repo: &users.RepositoryMock{
					GetByIDFunc: func(ctx context.Context, ID string) (models.User, error) {
						return models.User{}, models.NotFoundErr
					},
				},
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest("POST", "/users/1", strings.NewReader(``)),
			},
			wantCode: http.StatusNotFound,
			wantBody: []byte("Not Found" + "\n"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := Users{
				repo: tt.fields.repo,
			}
			u.Get(tt.args.w, tt.args.r)

			assert.Equal(t, tt.wantCode, tt.args.w.Code)
			assert.Equalf(t, tt.wantBody, tt.args.w.Body.Bytes(), "unexpected body want %s got %s", string(tt.wantBody), string(tt.args.w.Body.Bytes()))
		})
	}
}
