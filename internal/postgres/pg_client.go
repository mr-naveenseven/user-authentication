package postgres

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"
	logger "user-authentication/pkg"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/jackc/pgx/v5/pgxpool"
)

// PGClient represents the Postgres client connection.
type PGClient struct {
	Client       *pgxpool.Pool
	Timeout      int
	MigrationDir string
}

// NewPGClient creates a new instance of PGClient.
func NewPGClient(timeout int, migrationDir string) *PGClient {

	return &PGClient{
		Client:       nil,
		Timeout:      timeout,
		MigrationDir: migrationDir,
	}
}

// Connect establishes a connection to the Postgres database.
func (pgClient *PGClient) Connect(dbURI string) error {
	// Create a context with a timeout for the connection attempt
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(pgClient.Timeout)*time.Second)
	defer cancel()

	// Create a connection pool
	pool, err := pgxpool.New(ctx, dbURI)
	if err != nil {

		return fmt.Errorf("unable to connect to database: %w", err)
	}

	pgClient.Client = pool
	logger.Info("Connected to postgres database.....")

	return nil
}

// Initialize initializes the Postgres client by migrating the database tables.
func (pgClient *PGClient) Initialize() error {

	// Migrate the database tables
	if err := pgClient.migrateTables(); err != nil {
		log.Println("Error migrating tables:", err)

		return err
	}
	log.Println("Postgres client initialized successfully.")

	return nil
}

// nugrateTables applies database migrations using the golang-migrate library.
func (pgClient *PGClient) migrateTables() error {

	pgConnString := pgClient.Client.Config().ConnString()
	m, err := migrate.New(
		"file://"+pgClient.MigrationDir,
		pgConnString,
	)
	if err != nil {
		log.Println("Error creating migration instance:", err)

		return fmt.Errorf("error creating migration instance: %w", err)
	}
	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Println("Error applying migrations:", err)

		return fmt.Errorf("error applying migrations: %w", err)
	} else if errors.Is(err, migrate.ErrNoChange) {
		log.Println("No migrations to apply.")

		return nil
	} else if err == nil {
		log.Println("Migrations applied successfully.")
	}

	return nil

}

// Disconnect closes the connection to the Postgres database.
func (pgClient *PGClient) Disconnect() {
	if pgClient.Client != nil {
		pgClient.Client.Close()
		log.Println("Disconnected from Postgres database successfully.")
	} else {
		log.Println("Postgres client is already disconnected or was never connected.")
	}
}
