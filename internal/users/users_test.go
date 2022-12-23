package users

import (
	"context"
	"database/sql/driver"
	"testing"
	"workshop/internal/models"

	"github.com/stretchr/testify/assert"
)

const ID = "bacfa697-c25c-48fa-b603-beb4fa53f8eb"

func TestService_Create(t *testing.T) {
	type fields struct {
		repo Repository
	}
	type args struct {
		ctx  context.Context
		name string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    models.User
		wantErr error
	}{
		{
			name: "success",
			fields: fields{
				repo: &RepositoryMock{
					CreateFunc: func(ctx context.Context, name string) (models.User, error) {
						return models.User{ID: ID, Name: "mike"}, nil
					},
				},
			},
			args: args{
				ctx:  context.TODO(),
				name: "mike",
			},
			want: models.User{
				ID:   ID,
				Name: "mike",
			},
			wantErr: nil,
		},
		{
			name: "success",
			fields: fields{
				repo: &RepositoryMock{
					CreateFunc: func(ctx context.Context, name string) (models.User, error) {
						return models.User{}, driver.ErrBadConn
					},
				},
			},
			args: args{
				ctx:  context.TODO(),
				name: "mike",
			},
			want:    models.User{},
			wantErr: driver.ErrBadConn,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Service{
				repo: tt.fields.repo,
			}

			got, err := s.Create(tt.args.ctx, tt.args.name)

			if nil != tt.wantErr {
				assert.ErrorIs(t, err, tt.wantErr, "unexpected error")
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.want, got, "unexpected result")
		})
	}
}
