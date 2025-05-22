package drivers

import (
	"context"
	"fmt"
	"github.com/docker/go-connections/nat"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"quotes/internal/models"
	"slices"
	"strconv"
	"strings"
	"testing"
	"time"
)

func setupPostgresContainer(t *testing.T) (*pgxpool.Pool, func()) {
	ctx := context.Background()

	pgPort := "5432/tcp"
	dbName := "testdb"
	dbUser := "postgres"
	dbPassword := "postgres"

	req := testcontainers.ContainerRequest{
		Image:        "postgres:15",
		ExposedPorts: []string{pgPort},
		Env: map[string]string{
			"POSTGRES_DB":       dbName,
			"POSTGRES_USER":     dbUser,
			"POSTGRES_PASSWORD": dbPassword,
		},
		WaitingFor: wait.ForLog("database system is ready to accept connections").WithStartupTimeout(2 * time.Minute),
	}

	postgresContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		t.Fatalf("Failed to start postgres container: %v", err)
	}

	hostIP, err := postgresContainer.Host(ctx)
	require.NoError(t, err)

	mappedPort, err := postgresContainer.MappedPort(ctx, nat.Port(pgPort))
	require.NoError(t, err)

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		dbUser, dbPassword, hostIP, mappedPort.Port(), dbName)

	time.Sleep(2 * time.Second)

	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		postgresContainer.Terminate(ctx)
		t.Fatalf("Could not connect to database: %v", err)
	}

	err = setupSchema(ctx, pool)
	if err != nil {
		pool.Close()
		postgresContainer.Terminate(ctx)
		t.Fatalf("Could not set up schema: %v", err)
	}

	return pool, func() {
		pool.Close()
		if err := postgresContainer.Terminate(ctx); err != nil {
			t.Logf("Failed to terminate container: %s", err)
		}
	}
}

func setupSchema(ctx context.Context, pool *pgxpool.Pool) error {
	schema := createTestSchema

	_, err := pool.Exec(ctx, schema)
	return err
}

func createTestData(ctx context.Context, pool *pgxpool.Pool) ([]pgtype.UUID, error) {
	quoteIds := make([]pgtype.UUID, 5)
	for i := 0; i < 5; i++ {
		idBytes := uuid.New()
		id := pgtype.UUID{Bytes: idBytes, Valid: true}
		quoteIds[i] = id

		author := "author" + strconv.Itoa(i)
		text := "text" + strconv.Itoa(i)

		_, err := pool.Exec(ctx, queryCreateQuote,
			id,
			author,
			text,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to create person: %w", err)
		}
	}

	return quoteIds, nil
}

func TestGetQuoteById(t *testing.T) {
	pool, cleanup := setupPostgresContainer(t)
	defer cleanup()

	driver := NewQuoteDriver(pool)
	ctx := context.Background()

	quoteIds, err := createTestData(ctx, pool)
	require.NoError(t, err)
	require.NotEmpty(t, quoteIds)

	t.Run("GetQuoteById with existing id", func(t *testing.T) {
		expAuthor := "author0"
		expText := "text0"

		quote, err := driver.GetQuoteById(ctx, quoteIds[0])
		require.NoError(t, err)
		assert.Equal(t, expAuthor, quote.Author)
		assert.Equal(t, expText, quote.Text)
	})

	t.Run("GetQuoteById with invalid id", func(t *testing.T) {
		idBytes := uuid.New()
		id := pgtype.UUID{Bytes: idBytes, Valid: true}

		quote, err := driver.GetQuoteById(ctx, id)
		require.Error(t, err)
		require.Nil(t, quote)
		require.Equal(t, pgx.ErrNoRows, err)
	})
}

func TestCreateQuote(t *testing.T) {
	pool, cleanup := setupPostgresContainer(t)
	defer cleanup()

	driver := NewQuoteDriver(pool)
	ctx := context.Background()

	idBytes := uuid.New()
	id := pgtype.UUID{Bytes: idBytes, Valid: true}

	quote := &models.Quote{
		Id:     id,
		Author: "author",
		Text:   "text",
	}

	err := driver.CreateQuote(ctx, quote)
	require.NoError(t, err)

	expQuote, err := driver.GetQuoteById(ctx, quote.Id)

	require.NoError(t, err)
	assert.Equal(t, quote, expQuote)
}

func TestDeleteQuote(t *testing.T) {
	pool, cleanup := setupPostgresContainer(t)
	defer cleanup()

	driver := NewQuoteDriver(pool)
	ctx := context.Background()

	quoteIds, err := createTestData(ctx, pool)
	require.NoError(t, err)
	require.NotEmpty(t, quoteIds)

	err = driver.DeleteQuote(ctx, quoteIds[0])
	require.NoError(t, err)
}

func TestGetAllQuotes(t *testing.T) {
	pool, cleanup := setupPostgresContainer(t)
	defer cleanup()

	driver := NewQuoteDriver(pool)
	ctx := context.Background()

	quoteIds, err := createTestData(ctx, pool)
	require.NoError(t, err)
	require.NotEmpty(t, quoteIds)

	quotes, err := driver.GetAllQuotes(ctx)
	require.NoError(t, err)
	require.Len(t, quotes, len(quoteIds))
	require.Equal(t, quotes[0].Id, quoteIds[0])
	require.Equal(t, quotes[1].Author, "author1")
	require.Equal(t, quotes[2].Text, "text2")
}

func TestGetQuotesByAuthor(t *testing.T) {
	pool, cleanup := setupPostgresContainer(t)
	defer cleanup()

	driver := NewQuoteDriver(pool)
	ctx := context.Background()

	quoteIds, err := createTestData(ctx, pool)
	require.NoError(t, err)
	require.NotEmpty(t, quoteIds)

	author := "author0"

	quotes, err := driver.GetQuotesByAuthor(ctx, author)
	require.NoError(t, err)
	require.Equal(t, quotes[0].Id, quoteIds[0])
	require.Equal(t, quotes[0].Author, "author0")
	require.Equal(t, quotes[0].Text, "text0")
}

func TestGetRandomQuote(t *testing.T) {
	pool, cleanup := setupPostgresContainer(t)
	defer cleanup()

	driver := NewQuoteDriver(pool)
	ctx := context.Background()

	quoteIds, err := createTestData(ctx, pool)
	require.NoError(t, err)
	require.NotEmpty(t, quoteIds)

	quote, err := driver.GetRandomQuote(ctx)
	require.NoError(t, err)
	require.True(t, slices.Contains(quoteIds, quote.Id))
	require.True(t, strings.HasPrefix(quote.Author, "author"))
	require.True(t, strings.HasPrefix(quote.Text, "text"))
}
