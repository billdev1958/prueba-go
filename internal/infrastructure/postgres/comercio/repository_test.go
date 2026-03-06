package comercio

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	comercio "prueba-go/internal/domain/commerce"
	"prueba-go/pkg/util/money"
	"prueba-go/pkg/util/uuid"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

var testPool *pgxpool.Pool

func TestMain(m *testing.M) {
	ctx := context.Background()

	dbName := "testdb"
	dbUser := "user"
	dbPassword := "password"

	postgresContainer, err := postgres.Run(ctx,
		"postgres:16-alpine",
		postgres.WithDatabase(dbName),
		postgres.WithUsername(dbUser),
		postgres.WithPassword(dbPassword),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second)),
	)
	if err != nil {
		fmt.Printf("failed to start container: %s", err)
		os.Exit(1)
	}

	connStr, err := postgresContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		fmt.Printf("failed to get connection string: %s", err)
		_ = postgresContainer.Terminate(ctx)
		os.Exit(1)
	}

	testPool, err = pgxpool.New(ctx, connStr)
	if err != nil {
		fmt.Printf("failed to create pool: %s", err)
		_ = postgresContainer.Terminate(ctx)
		os.Exit(1)
	}

	schemaPath := filepath.Join("..", "..", "..", "..", "db", "schema.sql")
	schema, err := os.ReadFile(schemaPath)
	if err != nil {
		fmt.Printf("failed to read schema file: %s", err)
		testPool.Close()
		_ = postgresContainer.Terminate(ctx)
		os.Exit(1)
	}

	_, err = testPool.Exec(ctx, string(schema))
	if err != nil {
		fmt.Printf("failed to apply schema: %s", err)
		testPool.Close()
		_ = postgresContainer.Terminate(ctx)
		os.Exit(1)
	}

	// Run tests
	code := m.Run()

	// 4. Cleanup
	testPool.Close()
	if err := postgresContainer.Terminate(ctx); err != nil {
		fmt.Printf("failed to terminate container: %s", err)
	}

	os.Exit(code)
}

func cleanTable(t *testing.T) {
	_, err := testPool.Exec(context.Background(), "DELETE FROM comercios")
	require.NoError(t, err)
}

func TestComercioRepository_CRUD(t *testing.T) {
	repo := NewComercioRepository(testPool)
	ctx := context.Background()

	t.Run("Create and GetByID", func(t *testing.T) {
		cleanTable(t)
		rate, _ := money.NewRate("0.05")
		c := &comercio.Comercio{
			ID:            uuid.NewUUID(),
			Name:          "Test Store",
			ComissionRate: rate,
		}

		// Create
		res, err := repo.Create(ctx, c)
		assert.NoError(t, err)
		assert.Equal(t, c.ID, res.ID)
		assert.NotZero(t, res.CreatedAt)

		// GetByID
		found, err := repo.GetByID(ctx, c.ID)
		assert.NoError(t, err)
		assert.Equal(t, c.Name, found.Name)
		assert.Equal(t, "0.0500", found.ComissionRate.RateToString())
	})

	t.Run("Update", func(t *testing.T) {
		cleanTable(t)
		rate, _ := money.NewRate("0.05")
		c := &comercio.Comercio{
			ID:            uuid.NewUUID(),
			Name:          "Original Name",
			ComissionRate: rate,
		}
		_, _ = repo.Create(ctx, c)

		newRate, _ := money.NewRate("0.10")
		c.Name = "Updated Name"
		c.ComissionRate = newRate

		err := repo.Update(ctx, c)
		assert.NoError(t, err)

		found, _ := repo.GetByID(ctx, c.ID)
		assert.Equal(t, "Updated Name", found.Name)
		assert.Equal(t, "0.1000", found.ComissionRate.RateToString())
	})

	t.Run("Delete", func(t *testing.T) {
		cleanTable(t)
		id := uuid.NewUUID()
		rate, _ := money.NewRate("0.05")
		_, _ = repo.Create(ctx, &comercio.Comercio{
			ID:            id,
			Name:          "To Delete",
			ComissionRate: rate,
		})

		err := repo.Delete(ctx, id)
		assert.NoError(t, err)

		_, err = repo.GetByID(ctx, id)
		assert.ErrorIs(t, err, comercio.ErrComercioNotFound)
	})

	t.Run("GetAll", func(t *testing.T) {
		cleanTable(t)
		rate, _ := money.NewRate("0.05")
		_, _ = repo.Create(ctx, &comercio.Comercio{ID: uuid.NewUUID(), Name: "Store 1", ComissionRate: rate})
		_, _ = repo.Create(ctx, &comercio.Comercio{ID: uuid.NewUUID(), Name: "Store 2", ComissionRate: rate})

		list, err := repo.GetAll(ctx)
		assert.NoError(t, err)
		assert.Len(t, list, 2)
	})

	t.Run("NotFound Errors", func(t *testing.T) {
		cleanTable(t)
		nonExistentID := uuid.NewUUID()

		_, err := repo.GetByID(ctx, nonExistentID)
		assert.ErrorIs(t, err, comercio.ErrComercioNotFound)

		err = repo.Update(ctx, &comercio.Comercio{ID: nonExistentID, Name: "X"})
		assert.ErrorIs(t, err, comercio.ErrComercioNotFound)

		err = repo.Delete(ctx, nonExistentID)
		assert.ErrorIs(t, err, comercio.ErrComercioNotFound)
	})
}
