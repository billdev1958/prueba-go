package transaction

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	transaction "prueba-go/internal/domain/transaction"
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
var testMerchantID = uuid.NewUUID()

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

	// Insert a default merchant for tests
	_, err = testPool.Exec(ctx, "INSERT INTO comercios (id, name, comission_rate) VALUES ($1, $2, $3)", testMerchantID, "Test Merchant", "0.0500")
	if err != nil {
		fmt.Printf("failed to insert default merchant: %s", err)
		testPool.Close()
		_ = postgresContainer.Terminate(ctx)
		os.Exit(1)
	}

	// Run tests
	code := m.Run()

	// Cleanup
	testPool.Close()
	if err := postgresContainer.Terminate(ctx); err != nil {
		fmt.Printf("failed to terminate container: %s", err)
	}

	os.Exit(code)
}

func cleanTable(t *testing.T) {
	_, err := testPool.Exec(context.Background(), "DELETE FROM transactions")
	require.NoError(t, err)
}

func TestTransactionRepository_CRUD(t *testing.T) {
	repo := NewTransactionRepository(testPool)
	ctx := context.Background()

	t.Run("Create and GetByID", func(t *testing.T) {
		cleanTable(t)

		amt, _ := money.NewAmmount("100.00")
		rate, _ := money.NewRate("0.05")
		comm, _ := money.NewAmmount("5.00")
		net, _ := money.NewAmmount("95.00")

		tr := &transaction.Transaction{
			ID:          uuid.NewUUID(),
			CommercioID: testMerchantID,
			Amount:      amt,
			AppliedRate: rate,
			Commission:  comm,
			NetAmount:   net,
		}

		// Create
		res, err := repo.Create(ctx, tr)
		assert.NoError(t, err)
		assert.Equal(t, tr.ID, res.ID)
		assert.NotZero(t, res.CreatedAt)

		// GetByID
		found, err := repo.GetByID(ctx, tr.ID)
		assert.NoError(t, err)
		assert.Equal(t, tr.ID, found.ID)
		assert.Equal(t, "100.0000", found.Amount.AmmountToString())
		assert.Equal(t, "0.0500", found.AppliedRate.RateToString())
	})

	t.Run("GetAll", func(t *testing.T) {
		cleanTable(t)
		amt, _ := money.NewAmmount("100.00")
		rate, _ := money.NewRate("0.05")
		comm, _ := money.NewAmmount("5.00")
		net, _ := money.NewAmmount("95.00")

		_, _ = repo.Create(ctx, &transaction.Transaction{ID: uuid.NewUUID(), CommercioID: testMerchantID, Amount: amt, AppliedRate: rate, Commission: comm, NetAmount: net})
		_, _ = repo.Create(ctx, &transaction.Transaction{ID: uuid.NewUUID(), CommercioID: testMerchantID, Amount: amt, AppliedRate: rate, Commission: comm, NetAmount: net})

		list, err := repo.GetAll(ctx)
		assert.NoError(t, err)
		assert.Len(t, list, 2)
	})
}
