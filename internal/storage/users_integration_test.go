package storage

import (
	"context"
	"database/sql"
	"testing"
	"workshop/internal/models"
	"workshop/pkg/integration_testing"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const databaseDSN = "host=localhost port=5435 user=db_user password=db_pass dbname=test_workshop sslmode=disable"

func TestIntegrationStorage_Create(t *testing.T) {
	integration_testing.ShouldSkip(t)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	db, err := sql.Open("postgres", databaseDSN)
	if err != nil {
		t.Error("fail to connect to database", err)
	}

	t.Run("success", func(t *testing.T) {
		s := Storage{db: db}
		usr, err := s.Create(ctx, "mike")

		require.NoError(t, err)

		dbUser := models.User{}
		row := db.QueryRowContext(ctx, "SELECT id, name FROM users WHERE id = $1", usr.ID)
		err = row.Scan(&dbUser.ID, &dbUser.Name)
		require.NoError(t, err)

		assert.Equal(t, usr, dbUser)
	})
}
